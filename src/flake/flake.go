package flake

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/hymkor/go-lazy"
)

var (
	registry      = "__std" // keep for now for historic reasons
	flakeRegistry = func(flake string) string { return fmt.Sprintf("%[2]s#%[1]s", registry, flake) }
)

type outt struct {
	drvPath string            `json:"drvPath"`
	outputs map[string]string `json:"outputs"`
}

var CellsFrom = lazy.Of[string]{
	New: func() string {
		if s, err := getCellsFrom(); err != nil {
			return "${cellsFrom}"
		} else {
			return s
		}
	},
}

// tprintf passed template string is formatted usign its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func tprintf(data interface{}, tmpl string) string {
	t := template.Must(template.New("tmp").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func getNix() (string, error) {
	nix, err := exec.LookPath("nix")
	if err != nil {
		return "", errors.New("You need to install 'nix' in order to use this tool")
	}
	return nix, nil
}

func getCurrentSystem() (string, error) {
	// detect the current system
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	currentSystem, err := exec.Command(
		nix, "eval", "--raw", "--impure", "--expr", "builtins.currentSystem",
	).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", err
	}
	currentSystemStr := string(currentSystem)
	return currentSystemStr, nil
}

func getCellsFrom() (string, error) {
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	cellsFrom, err := exec.Command(
		nix, "eval", "--raw", flakeRegistry(".")+".cellsFrom",
	).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", err
	}
	return string(cellsFrom[:]), nil
}
