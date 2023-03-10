package flake

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/paisano-nix/paisano/cache"
	"github.com/paisano-nix/paisano/env"
)

var (
	registry = "__std" // keep for now for historic reasons
)

type outt struct {
	drvPath string            `json:"drvPath"`
	outputs map[string]string `json:"outputs"`
}

var (
	currentSystemArgs      = []string{"eval", "--raw", "--impure", "--expr", "builtins.currentSystem"}
	cellsFromArgs          = []string{"eval", "--raw"}
	flakeRegistry          = func(flake string) string { return fmt.Sprintf("%[2]s#%[1]s", registry, flake) }
	flakeCellsFromFragment = func(flake string) string { return fmt.Sprintf("%[1]s.cellsFrom", flakeRegistry(flake)) }
	flakeInitFragment      = func(flake, system string) string {
		return fmt.Sprintf("%[1]s.init.%[2]s", flakeRegistry(flake), system)
	}
	flakeActionsFragment = func(flake, system, c, b, t, a string) string {
		return fmt.Sprintf("%[1]s.actions.%[2]s.%[3]s.%[4]s.%[5]s.%[6]s", flakeRegistry(flake), system, c, b, t, a)
	}
	flakeEvalJson = []string{
		"eval",
		"--json",
		"--no-update-lock-file",
		"--no-write-lock-file",
		"--no-warn-dirty",
		"--accept-flake-config",
	}
	flakeBuild = func(out string) []string {
		return []string{
			"build",
			"--out-link", out,
			"--no-update-lock-file",
			"--no-write-lock-file",
			"--no-warn-dirty",
			"--accept-flake-config",
			"--builders-use-substitutes",
		}
	}
)

func getNix() (string, error) {
	nix, err := exec.LookPath("nix")
	if err != nil {
		return "", errors.New("You need to install 'nix' in order to use this tool")
	}
	return nix, nil
}

func getCurrentSystem() (*string, error) {
	// detect the current system
	nix, err := getNix()
	if err != nil {
		return nil, err
	}
	currentSystem, err := exec.Command(nix, currentSystemArgs...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return nil, err
	}
	currentSystemStr := string(currentSystem)
	return &currentSystemStr, nil
}

func GetCellsFrom() (string, error) {
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	cellsFrom, err := exec.Command(nix, append(cellsFromArgs, flakeCellsFromFragment("."))...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", err
	}
	return string(cellsFrom[:]), nil
}

func getInitEvalCmdArgs() (string, []string, error) {
	nix, err := getNix()
	if err != nil {
		return "", nil, err
	}

	currentSystem, err := getCurrentSystem()
	if err != nil {
		return "", nil, err
	}

	return nix, append(
		flakeEvalJson, flakeInitFragment(".", *currentSystem)), nil
}

func GetActionEvalCmdArgs(c, b, t, a string, system *string) (string, []string, error) {
	nix, err := getNix()
	if err != nil {
		return "", nil, err
	}

	_, _, _, actionPath, err := env.SetEnv()
	if err != nil {
		return "", nil, err
	}

	currentSystem, err := getCurrentSystem()
	if err != nil {
		return "", nil, err
	}

	if *system != "" {
		// if we specify the current system it could be used, in theory,
		// as a general hack to pass the impure flag, but we only use
		// the impure flag as a transport for the very specific use case
		// of conveying the current system to the action evaluation
		// as an action is always run on the local (i.e. "current") system
		// therefore, close this loophole and error if not for the impure flag
		// it would otherwise be and is intended to be a no-op
		// systemVal := *system
		if *system == *currentSystem {
			return "", nil, fmt.Errorf("set the --for flag to a different system than the current one ('%s')", *currentSystem)
		}
		// if system is set, the impure flag provides a "hack" so that we
		// can transport this information to the action evaluation without
		// incurring in a prohibitively complex (m*n) data structure in
		// which we would have to account for _all_ combinations of current
		// and build system
		return nix, append(
			flakeBuild(actionPath), "--impure", flakeActionsFragment(".", *system, c, b, t, a)), nil
	}

	return nix, append(
		flakeBuild(actionPath), flakeActionsFragment(".", *currentSystem, c, b, t, a)), nil
}

func LoadFlakeCmd() (*cache.Cache, *cache.ActionID, *exec.Cmd, *bytes.Buffer, error) {

	nix, args, err := getInitEvalCmdArgs()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// load the paisano metadata from the flake
	buf := new(bytes.Buffer)
	cmd := exec.Command(nix, args...)
	cmd.Stdin = devNull
	cmd.Stdout = buf

	// initialize cache
	_, _, prjCacheDir, _, err := env.SetEnv()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	path := prjCacheDir
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	c, err := cache.Open(path)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	key := cache.NewActionID([]byte(strings.Join(args, "")))

	return c, &key, cmd, buf, nil
}
