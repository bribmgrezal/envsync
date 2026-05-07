package envfile

import (
	"testing"
)

func TestMerge_OverwritesExistingKeys(t *testing.T) {
	dst := []Entry{{Key: "HOST", Value: "localhost"}, {Key: "PORT", Value: "8080"}}
	src := []Entry{{Key: "PORT", Value: "9090"}}

	result := Merge(dst, src, DefaultMergeOptions())

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[1].Value != "9090" {
		t.Errorf("expected PORT=9090, got %s", result[1].Value)
	}
}

func TestMerge_AppendsMissingKeys(t *testing.T) {
	dst := []Entry{{Key: "HOST", Value: "localhost"}}
	src := []Entry{{Key: "DEBUG", Value: "true"}}

	result := Merge(dst, src, DefaultMergeOptions())

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[1].Key != "DEBUG" || result[1].Value != "true" {
		t.Errorf("unexpected last entry: %+v", result[1])
	}
}

func TestMerge_NoOverwrite_KeepsDstValue(t *testing.T) {
	dst := []Entry{{Key: "PORT", Value: "8080"}}
	src := []Entry{{Key: "PORT", Value: "9090"}}

	opts := MergeOptions{Overwrite: false, SkipEmpty: false}
	result := Merge(dst, src, opts)

	if result[0].Value != "8080" {
		t.Errorf("expected PORT=8080 (no overwrite), got %s", result[0].Value)
	}
}

func TestMerge_SkipEmpty_OmitsBlankSrcValues(t *testing.T) {
	dst := []Entry{{Key: "HOST", Value: "localhost"}}
	src := []Entry{{Key: "NEW_KEY", Value: ""}}

	opts := MergeOptions{Overwrite: true, SkipEmpty: true}
	result := Merge(dst, src, opts)

	if len(result) != 1 {
		t.Errorf("expected 1 entry (empty src skipped), got %d", len(result))
	}
}

func TestMerge_EmptyDst_ReturnsSrcEntries(t *testing.T) {
	dst := []Entry{}
	src := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}

	result := Merge(dst, src, DefaultMergeOptions())

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestMerge_PreservesDstOrder(t *testing.T) {
	dst := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}, {Key: "C", Value: "3"}}
	src := []Entry{{Key: "B", Value: "20"}, {Key: "D", Value: "4"}}

	result := Merge(dst, src, DefaultMergeOptions())

	keys := []string{"A", "B", "C", "D"}
	for i, want := range keys {
		if result[i].Key != want {
			t.Errorf("position %d: expected key %s, got %s", i, want, result[i].Key)
		}
	}
}
