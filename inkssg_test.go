package inkssg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildCopiesStaticToOutputRoot(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("pages/index/content.md", "# hello")
	write("static/robots.txt", "robots")
	write("static/img/logo.txt", "logo")
	write("static/index.html", "static index")

	if err := Build(dir); err != nil {
		t.Fatal(err)
	}

	for rel, want := range map[string]string{
		"robots.txt":   "robots",
		"img/logo.txt": "logo",
	} {
		got, err := os.ReadFile(filepath.Join(dir, "public", rel))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", rel, got, want)
		}
	}

	index, err := os.ReadFile(filepath.Join(dir, "public", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(index) == "static index" {
		t.Error("static/index.html replaced the built page")
	}
}
