// Package speech reads text aloud with the voice built into the
// operating system. Nothing is bundled and nothing leaves the machine.
package speech

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// ErrUnavailable means no speech engine was found on this system.
var ErrUnavailable = errors.New(unavailableMessage())

const windowsScript = `[Console]::InputEncoding = [Text.Encoding]::UTF8; ` +
	`Add-Type -AssemblyName System.Speech; ` +
	`$s = New-Object System.Speech.Synthesis.SpeechSynthesizer; ` +
	`$s.Speak([Console]::In.ReadToEnd())`

func unavailableMessage() string {
	switch runtime.GOOS {
	case "windows":
		return "no voice found: Windows PowerShell is required for --sound"
	case "darwin":
		return "no voice found: the macOS 'say' command is missing"
	default:
		return "no voice found: install espeak-ng to read aloud"
	}
}

// engine returns the program and arguments that speak text piped to stdin.
func engine(goos string, look func(string) (string, error)) (string, []string, bool) {
	var candidates [][]string

	switch goos {
	case "windows":
		candidates = [][]string{
			{"powershell.exe", "-NoProfile", "-NonInteractive", "-Command", windowsScript},
		}
	case "darwin":
		candidates = [][]string{
			{"say", "-f", "-"},
		}
	default:
		candidates = [][]string{
			{"espeak-ng", "--stdin"},
			{"espeak", "--stdin"},
		}
	}

	for _, c := range candidates {
		if path, err := look(c[0]); err == nil {
			return path, c[1:], true
		}
	}

	return "", nil, false
}

// Clean prepares verse text for a speech engine: typographic marks that
// engines read out or stumble on are removed or replaced.
func Clean(text string) string {
	text = strings.NewReplacer(
		"¶", " ",
		"[", "", "]", "",
		"’", "'", "‘", "'",
		"“", "\"", "”", "\"",
		"—", ", ", "–", ", ",
	).Replace(text)

	return strings.Join(strings.Fields(text), " ")
}

// Passage joins an optional announcement and verse texts into one
// utterance, each on its own line so the engine pauses between them.
func Passage(announce string, verses []string) string {
	var parts []string

	if announce = Clean(announce); announce != "" {
		parts = append(parts, announce+".")
	}

	for _, v := range verses {
		if v = Clean(v); v != "" {
			parts = append(parts, v)
		}
	}

	return strings.Join(parts, "\n")
}

// Speaker is one utterance in progress.
type Speaker struct {
	cmd  *exec.Cmd
	done chan struct{}
	once sync.Once
	err  error
}

// Start begins speaking text and returns immediately.
func Start(text string) (*Speaker, error) {
	path, args, ok := engine(runtime.GOOS, exec.LookPath)
	if !ok {
		return nil, ErrUnavailable
	}

	cmd := exec.Command(path, args...)
	cmd.Stdin = strings.NewReader(text)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	s := &Speaker{cmd: cmd, done: make(chan struct{})}

	go func() {
		s.err = cmd.Wait()
		close(s.done)
	}()

	return s, nil
}

// Done is closed when speaking has finished or been stopped.
func (s *Speaker) Done() <-chan struct{} {
	return s.done
}

// Wait blocks until speaking finishes.
func (s *Speaker) Wait() error {
	<-s.done
	return s.err
}

// Stop silences the speaker. It is safe to call more than once.
func (s *Speaker) Stop() {
	s.once.Do(func() {
		if s.cmd.Process != nil {
			s.cmd.Process.Kill()
		}
	})
	<-s.done
}

// Speak reads text aloud and blocks until it is finished.
func Speak(text string) error {
	s, err := Start(text)
	if err != nil {
		return err
	}

	return s.Wait()
}
