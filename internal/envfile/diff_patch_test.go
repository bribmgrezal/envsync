package envfile

import (
	"testing"
)

func diffPatchSrc() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DEBUG", Value: "true"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func diffPatchDst() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DEBUG", Value: "false"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestBuildPatchFromDiff_AddsNewKeys(t *testing.T) {
	patch, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), DefaultDiffPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patch.Add) != 1 || patch.Add[0].Key != "DB_HOST" {
		t.Errorf("expected Add=[DB_HOST], got %+v", patch.Add)
	}
}

func TestBuildPatchFromDiff_UpdatesChangedKeys(t *testing.T) {
	patch, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), DefaultDiffPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patch.Update) != 1 || patch.Update[0].Key != "DEBUG" || patch.Update[0].Value != "true" {
		t.Errorf("expected Update=[DEBUG=true], got %+v", patch.Update)
	}
}

func TestBuildPatchFromDiff_NoRemovals_WhenDisabled(t *testing.T) {
	patch, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), DefaultDiffPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patch.Remove) != 0 {
		t.Errorf("expected no removals, got %v", patch.Remove)
	}
}

func TestBuildPatchFromDiff_RemovalsWhenEnabled(t *testing.T) {
	opts := DefaultDiffPatchOptions()
	opts.IncludeRemovals = true
	patch, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patch.Remove) != 1 || patch.Remove[0] != "LOG_LEVEL" {
		t.Errorf("expected Remove=[LOG_LEVEL], got %v", patch.Remove)
	}
}

func TestBuildPatchFromDiff_IgnoreKeys_SkipsKey(t *testing.T) {
	opts := DefaultDiffPatchOptions()
	opts.IgnoreKeys = []string{"DEBUG"}
	patch, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range patch.Update {
		if e.Key == "DEBUG" {
			t.Error("DEBUG should have been ignored but appeared in Update")
		}
	}
}

func TestBuildPatchFromDiff_BlankIgnoreKey_ReturnsError(t *testing.T) {
	opts := DefaultDiffPatchOptions()
	opts.IgnoreKeys = []string{"  "}
	_, err := BuildPatchFromDiff(diffPatchSrc(), diffPatchDst(), opts)
	if err == nil {
		t.Error("expected error for blank ignore key, got nil")
	}
}

func TestBuildPatchFromDiff_BothEmpty_ReturnsEmpty(t *testing.T) {
	patch, err := BuildPatchFromDiff(nil, nil, DefaultDiffPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patch.Add)+len(patch.Update)+len(patch.Remove) != 0 {
		t.Errorf("expected empty patch, got %+v", patch)
	}
}
