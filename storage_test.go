package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// setupStorage points the package-level paths at a temp dir, replaces the
// mount function with a fake that succeeds, and restores everything after
// the test. Tests using it must not run in parallel.
func setupStorage(t *testing.T) (local, backupPath string, mountCalls *int) {
	t.Helper()

	dir := t.TempDir()
	local = filepath.Join(dir, "local", "quotes.json")
	backupPath = filepath.Join(dir, "backup", "quotes.json")

	oldLocal, oldBackup, oldMount, oldNotices := LocalFile, BackupFile, mountBackup, notices

	calls := 0
	LocalFile = local
	BackupFile = backupPath
	mountBackup = func() error { calls++; return nil }
	notices = nil

	t.Cleanup(func() {
		LocalFile, BackupFile, mountBackup, notices = oldLocal, oldBackup, oldMount, oldNotices
	})

	return local, backupPath, &calls
}

func sampleQuotes() *Quotes {
	return &Quotes{QUOTES: []Quote{
		{ID: 1, QUOTE: "Be yourself", AUTHOR: "Oscar Wilde", LANGUAGE: "English", DATE: "2026-01-02"},
		{ID: 2, QUOTE: "Tark nagu öö", AUTHOR: "Anonymous", LANGUAGE: "Estonian", DATE: ""},
	}}
}

func noTempFiles(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Errorf("leftover temp file in %s: %s", dir, e.Name())
		}
	}
}

/* Load */

func TestLoad_MissingFileStartsEmpty(t *testing.T) {
	setupStorage(t)

	q := Quotes{}
	if err := q.Load(filepath.Join(t.TempDir(), "nope.json")); err != nil {
		t.Fatalf("Load returned error for missing file: %v", err)
	}
	if len(q.QUOTES) != 0 {
		t.Errorf("expected empty quotes, got %d", len(q.QUOTES))
	}
	if len(notices) != 1 || !strings.Contains(notices[0], "No data file found") {
		t.Errorf("expected 'No data file found' notice, got %v", notices)
	}
}

func TestLoad_ValidJSON(t *testing.T) {
	setupStorage(t)

	path := filepath.Join(t.TempDir(), "quotes.json")
	content := `{"quotes":[{"id":7,"quote":"Hello","author":"Me","language":"English","date":"2026-05-01"}]}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	q := Quotes{}
	if err := q.Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(q.QUOTES) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(q.QUOTES))
	}
	want := Quote{ID: 7, QUOTE: "Hello", AUTHOR: "Me", LANGUAGE: "English", DATE: "2026-05-01"}
	if q.QUOTES[0] != want {
		t.Errorf("got %+v, want %+v", q.QUOTES[0], want)
	}
	if len(notices) != 0 {
		t.Errorf("unexpected notices: %v", notices)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	setupStorage(t)

	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte(`{"quotes": [`), 0o600); err != nil {
		t.Fatal(err)
	}

	q := Quotes{}
	err := q.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing JSON") {
		t.Errorf("error should mention parsing JSON, got: %v", err)
	}
}

func TestLoad_ReadErrorIsNotTreatedAsMissing(t *testing.T) {
	setupStorage(t)

	// Reading a directory fails with an error that is not ErrNotExist.
	q := Quotes{}
	err := q.Load(t.TempDir())
	if err == nil {
		t.Fatal("expected error when source is a directory")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("should not be ErrNotExist: %v", err)
	}
	if len(notices) != 0 {
		t.Errorf("no 'starting empty' notice expected, got %v", notices)
	}
}

/* Save */

func TestSave_WritesLocalAndBackup(t *testing.T) {
	local, backupPath, mountCalls := setupStorage(t)

	q := sampleQuotes()
	if err := q.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	for _, p := range []string{local, backupPath} {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		var got Quotes
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("%s is not valid JSON: %v", p, err)
		}
		if len(got.QUOTES) != 2 || got.QUOTES[0] != q.QUOTES[0] || got.QUOTES[1] != q.QUOTES[1] {
			t.Errorf("%s content mismatch: %+v", p, got.QUOTES)
		}
	}

	if *mountCalls != 1 {
		t.Errorf("mount called %d times, want 1", *mountCalls)
	}
	if len(notices) != 0 {
		t.Errorf("unexpected notices: %v", notices)
	}
}

func TestSave_RoundTripWithLoad(t *testing.T) {
	local, _, _ := setupStorage(t)

	want := sampleQuotes()
	if err := want.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got := Quotes{}
	if err := got.Load(local); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.QUOTES) != len(want.QUOTES) {
		t.Fatalf("got %d quotes, want %d", len(got.QUOTES), len(want.QUOTES))
	}
	for i := range want.QUOTES {
		if got.QUOTES[i] != want.QUOTES[i] {
			t.Errorf("quote %d: got %+v, want %+v", i, got.QUOTES[i], want.QUOTES[i])
		}
	}
}

func TestSave_JSONFieldNames(t *testing.T) {
	local, _, _ := setupStorage(t)

	if err := sampleQuotes().Save(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(local)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"quotes"`, `"id"`, `"quote"`, `"author"`, `"language"`, `"date"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("saved JSON missing key %s", key)
		}
	}
	if !strings.Contains(string(data), "õ") && !strings.Contains(string(data), "ö") {
		t.Errorf("non-ASCII characters should be stored as is")
	}
}

func TestSave_BackupFailureDoesNotFailSave(t *testing.T) {
	local, backupPath, _ := setupStorage(t)

	mountBackup = func() error { return errors.New("drive not found") }

	if err := sampleQuotes().Save(); err != nil {
		t.Fatalf("Save should succeed when only the backup fails: %v", err)
	}

	if _, err := os.Stat(local); err != nil {
		t.Errorf("local file should exist: %v", err)
	}
	if _, err := os.Stat(backupPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("backup file should not exist, stat err: %v", err)
	}

	if len(notices) != 1 {
		t.Fatalf("expected 1 notice, got %v", notices)
	}
	if !strings.Contains(notices[0], "Backup failed") || !strings.Contains(notices[0], "drive not found") {
		t.Errorf("unexpected notice: %q", notices[0])
	}
}

func TestSave_LocalFailureReturnsError(t *testing.T) {
	_, _, mountCalls := setupStorage(t)

	// A regular file where a directory is needed makes MkdirAll fail.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	LocalFile = filepath.Join(blocker, "sub", "quotes.json")

	err := sampleQuotes().Save()
	if err == nil {
		t.Fatal("expected error when local save fails")
	}
	if !strings.Contains(err.Error(), "local save") {
		t.Errorf("error should mention local save, got: %v", err)
	}
	if *mountCalls != 0 {
		t.Errorf("backup should not be attempted after a local failure, mount called %d times", *mountCalls)
	}
}

func TestSave_OverwritesExistingFile(t *testing.T) {
	local, _, _ := setupStorage(t)

	q := sampleQuotes()
	if err := q.Save(); err != nil {
		t.Fatal(err)
	}

	q.QUOTES = q.QUOTES[:1]
	if err := q.Save(); err != nil {
		t.Fatal(err)
	}

	got := Quotes{}
	if err := got.Load(local); err != nil {
		t.Fatal(err)
	}
	if len(got.QUOTES) != 1 {
		t.Errorf("expected 1 quote after overwrite, got %d", len(got.QUOTES))
	}
}

func TestSave_EmptyCollection(t *testing.T) {
	local, _, _ := setupStorage(t)

	q := &Quotes{}
	if err := q.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got := Quotes{}
	if err := got.Load(local); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.QUOTES) != 0 {
		t.Errorf("expected 0 quotes, got %d", len(got.QUOTES))
	}
}

func TestSave_LeavesNoTempFiles(t *testing.T) {
	local, backupPath, _ := setupStorage(t)

	if err := sampleQuotes().Save(); err != nil {
		t.Fatal(err)
	}
	noTempFiles(t, filepath.Dir(local))
	noTempFiles(t, filepath.Dir(backupPath))
}

/* writeFileAtomic */

func TestWriteFileAtomic_CreatesNestedDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c", "file.txt")

	if err := writeFileAtomic(path, []byte("hello"), 0o600); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", data, "hello")
	}
}

func TestWriteFileAtomic_ReplacesContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")

	if err := writeFileAtomic(path, []byte("first, longer content"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(path, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "second" {
		t.Errorf("got %q, want %q (no leftover bytes from the old content)", data, "second")
	}
	noTempFiles(t, filepath.Dir(path))
}

func TestWriteFileAtomic_SetsPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file modes are not meaningful on Windows")
	}

	path := filepath.Join(t.TempDir(), "file.txt")
	if err := writeFileAtomic(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("permissions = %o, want 600", got)
	}
}

func TestWriteFileAtomic_FailureLeavesNoTempFileAndKeepsOldData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := writeFileAtomic(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Make the rename fail by turning the destination into a non-empty directory.
	bad := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(bad, "child"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(bad, []byte("new"), 0o600); err == nil {
		t.Fatal("expected error when destination is a non-empty directory")
	}

	noTempFiles(t, dir)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "original" {
		t.Errorf("unrelated file changed: %q", data)
	}
}

/* backup */

func TestBackup_MountErrorIsWrapped(t *testing.T) {
	setupStorage(t)

	sentinel := errors.New("no drive")
	mountBackup = func() error { return sentinel }

	err := backup([]byte("data"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("error should wrap the mount error: %v", err)
	}
	if !strings.Contains(err.Error(), "mount drive") {
		t.Errorf("error should mention mount drive: %v", err)
	}
}
