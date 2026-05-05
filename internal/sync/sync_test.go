package sync_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/sync"
)

func TestApply_AddsMissingKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "NEW_KEY", Status: diff.Missing, SourceValue: "new_val"},
	}
	target := map[string]string{"EXISTING": "val"}
	out, summary := sync.Apply(results, target, sync.DefaultOptions())

	if out["NEW_KEY"] != "new_val" {
		t.Errorf("expected NEW_KEY=new_val, got %q", out["NEW_KEY"])
	}
	if !strings.Contains(summary, "+ NEW_KEY") {
		t.Errorf("expected summary to mention NEW_KEY addition, got: %s", summary)
	}
}

func TestApply_UpdatesChangedKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Changed, SourceValue: "prod-db", TargetValue: "dev-db"},
	}
	target := map[string]string{"DB_HOST": "dev-db"}
	out, summary := sync.Apply(results, target, sync.DefaultOptions())

	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %q", out["DB_HOST"])
	}
	if !strings.Contains(summary, "~ DB_HOST") {
		t.Errorf("expected summary to mention DB_HOST update, got: %s", summary)
	}
}

func TestApply_RemovesExtraKeys_WhenEnabled(t *testing.T) {
	results := []diff.Result{
		{Key: "STALE_KEY", Status: diff.Extra, TargetValue: "old"},
	}
	target := map[string]string{"STALE_KEY": "old", "KEEP": "yes"}
	opts := sync.Options{RemoveExtra: true}
	out, summary := sync.Apply(results, target, opts)

	if _, ok := out["STALE_KEY"]; ok {
		t.Error("expected STALE_KEY to be removed")
	}
	if out["KEEP"] != "yes" {
		t.Error("expected KEEP to remain")
	}
	if !strings.Contains(summary, "- STALE_KEY") {
		t.Errorf("expected summary to mention STALE_KEY removal, got: %s", summary)
	}
}

func TestApply_RespectsDisabledOptions(t *testing.T) {
	results := []diff.Result{
		{Key: "MISSING", Status: diff.Missing, SourceValue: "x"},
		{Key: "CHANGED", Status: diff.Changed, SourceValue: "new", TargetValue: "old"},
	}
	target := map[string]string{"CHANGED": "old"}
	opts := sync.Options{AddMissing: false, UpdateChanged: false}
	out, summary := sync.Apply(results, target, opts)

	if _, ok := out["MISSING"]; ok {
		t.Error("expected MISSING to not be added")
	}
	if out["CHANGED"] != "old" {
		t.Error("expected CHANGED to remain old")
	}
	if summary != "no changes applied" {
		t.Errorf("expected no-change summary, got: %s", summary)
	}
}

func TestApply_NoResults_ReturnsNoChanges(t *testing.T) {
	target := map[string]string{"A": "1"}
	out, summary := sync.Apply(nil, target, sync.DefaultOptions())
	if out["A"] != "1" {
		t.Error("expected target to be preserved")
	}
	if summary != "no changes applied" {
		t.Errorf("unexpected summary: %s", summary)
	}
}
