package envfile

import (
	"os"
	"testing"
)

func interpolateEntries() []Entry {
	return []Entry{
		{Key: "BASE", Value: "/usr/local"},
		{Key: "BIN", Value: "${BASE}/bin"},
		{Key: "LIB", Value: "$BASE/lib"},
		{Key: "GREETING", Value: "hello world"},
		{Key: "MULTI", Value: "${BASE}/share:${BASE}/include"},
	}
}

func TestInterpolate_ResolvesBlockStyle(t *testing.T) {
	entries, err := Interpolate(interpolateEntries(), DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "BIN", "/usr/local/bin")
}

func TestInterpolate_ResolvesDollarStyle(t *testing.T) {
	entries, err := Interpolate(interpolateEntries(), DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "LIB", "/usr/local/lib")
}

func TestInterpolate_MultipleRefsInValue(t *testing.T) {
	entries, err := Interpolate(interpolateEntries(), DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "MULTI", "/usr/local/share:/usr/local/include")
}

func TestInterpolate_PlainValueUntouched(t *testing.T) {
	entries, err := Interpolate(interpolateEntries(), DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "GREETING", "hello world")
}

func TestInterpolate_FallsBackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	input := []Entry{{Key: "X", Value: "${OS_VAR}"}}
	opts := DefaultInterpolateOptions()
	opts.UseOS = true
	entries, err := Interpolate(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "X", "from-os")
}

func TestInterpolate_UseOS_False_LeavesOSRefIntact(t *testing.T) {
	os.Setenv("OS_VAR2", "should-not-appear")
	t.Cleanup(func() { os.Unsetenv("OS_VAR2") })
	input := []Entry{{Key: "X", Value: "${OS_VAR2}"}}
	opts := DefaultInterpolateOptions()
	opts.UseOS = false
	entries, err := Interpolate(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValue(t, entries, "X", "${OS_VAR2}")
}

func TestInterpolate_FailOnMissing_ReturnsError(t *testing.T) {
	input := []Entry{{Key: "X", Value: "${UNDEFINED_XYZ}"}}
	opts := DefaultInterpolateOptions()
	opts.UseOS = false
	opts.FailOnMissing = true
	_, err := Interpolate(input, opts)
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestInterpolate_DoesNotMutateOriginal(t *testing.T) {
	original := interpolateEntries()
	copy := make([]Entry, len(original))
	for i, e := range original {
		copy[i] = e
	}
	_, _ = Interpolate(original, DefaultInterpolateOptions())
	for i, e := range original {
		if e.Value != copy[i].Value {
			t.Errorf("entry %d mutated: got %q want %q", i, e.Value, copy[i].Value)
		}
	}
}

func assertValue(t *testing.T, entries []Entry, key, want string) {
	t.Helper()
	for _, e := range entries {
		if e.Key == key {
			if e.Value != want {
				t.Errorf("key %q: got %q, want %q", key, e.Value, want)
			}
			return
		}
	}
	t.Errorf("key %q not found in entries", key)
}
