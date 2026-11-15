package speech

import (
	"errors"
	"testing"
)

func TestClean(t *testing.T) {
	cases := map[string]string{
		"¶ For God so loved the world": "For God so loved the world",
		"his wife’s name  Eve":         "his wife's name Eve",
		"and the earth [was] without":  "and the earth was without",
		"  spaced\tout \n text  ":      "spaced out text",
		"":                             "",
	}

	for in, want := range cases {
		if got := Clean(in); got != want {
			t.Errorf("Clean(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassage(t *testing.T) {
	got := Passage("John, chapter 3", []string{"First verse.", "", "¶ Second verse."})
	want := "John, chapter 3.\nFirst verse.\nSecond verse."

	if got != want {
		t.Errorf("Passage = %q, want %q", got, want)
	}

	if got := Passage("", []string{"Only."}); got != "Only." {
		t.Errorf("Passage without announcement = %q", got)
	}
}

func TestEngine(t *testing.T) {
	found := func(names ...string) func(string) (string, error) {
		return func(name string) (string, error) {
			for _, n := range names {
				if n == name {
					return "/bin/" + name, nil
				}
			}
			return "", errors.New("not found")
		}
	}

	if path, args, ok := engine("windows", found("powershell.exe")); !ok ||
		path != "/bin/powershell.exe" || args[len(args)-1] != windowsScript {
		t.Errorf("windows engine = %q %q %v", path, args, ok)
	}

	if path, _, ok := engine("darwin", found("say")); !ok || path != "/bin/say" {
		t.Errorf("darwin engine = %q %v", path, ok)
	}

	if path, _, ok := engine("linux", found("espeak")); !ok || path != "/bin/espeak" {
		t.Errorf("linux should fall back to espeak, got %q %v", path, ok)
	}

	if path, _, ok := engine("linux", found("espeak", "espeak-ng")); !ok || path != "/bin/espeak-ng" {
		t.Errorf("linux should prefer espeak-ng, got %q %v", path, ok)
	}

	if _, _, ok := engine("linux", found()); ok {
		t.Error("expected no engine when nothing is installed")
	}
}
