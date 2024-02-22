package flake

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/paisano-nix/paisano/env"
)

type RunActionCmd struct {
	ShowCmdStr   bool
	CmdStr       string
	System       string
	Cell         string
	Block        string
	Target       string
	Action       string
	RequiresArgs *bool
}

func (c *RunActionCmd) Assemble(extraArgs []string) (string, []string, error) {
	nix, err := getNix()
	if err != nil {
		return "", nil, err
	}

	currentSystem, err := getCurrentSystem()
	if err != nil {
		return "", nil, err
	}

	args, err := c.getArgs(currentSystem)
	if err != nil {
		return "", nil, err
	}

	if extraArgs != nil && len(extraArgs) > 0 {
		args = append(args, "--")
		args = append(args, extraArgs...)
	}
	return nix, args, nil
}

func (c *RunActionCmd) Build(nix string, args, extraArgs []string) (string, []string, error) {
	if c.RequiresArgs != nil && *c.RequiresArgs == true && len(extraArgs) == 0 {
		return "", nil, errors.New(c.CmdStr + " - requires on or more arguments; run from command line")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		return "", nil, err
	}
	// grep, err := exec.LookPath("grep")
	// if err != nil {
	// 	return "", nil, err
	// }
	nom, err := exec.LookPath("nom")
	if err == nil {
		nix = nom
	}
	actionPath, err := env.GetStateActionPath()
	if err != nil {
		return "", nil, err
	}
	printout := ""
	if c.ShowCmdStr {
		printout += "echo -e \"\x1b[1;37m------------" + strings.Repeat("-", len(c.CmdStr)) + "-\x1b[0m\";"
		printout += "echo -e \"\x1b[1;37m Executing: \x1b[1;32m" + c.CmdStr + "\x1b[0m\";"
		printout += "echo -e \"\x1b[1;37m------------" + strings.Repeat("-", len(c.CmdStr)) + "-\x1b[0m\";"
	}
	args = append(args, "--out-link", actionPath)
	args = append(args,
		"--no-update-lock-file",
		"--no-write-lock-file",
		"--no-warn-dirty",
		"--accept-flake-config",
		"--builders-use-substitutes",
		"|| exit 1;",
		printout,
		"exec", actionPath, "\"${@}\"",
	)
	cmd := []string{bash, "-c", nix + " build " + strings.Join(args, " ")}
	if extraArgs != nil && len(extraArgs) > 0 {
		cmd = append(cmd, "--")
		cmd = append(cmd, extraArgs...)
	}
	// fmt.Printf("%+v\n", cmd)
	return bash, cmd, nil
}

func (c *RunActionCmd) Exec(extraArgs []string) error {

	nix, args, err := c.Assemble(nil)
	if err != nil {
		return err
	}

	bash, cmd, err := c.Build(nix, args, extraArgs)
	if err != nil {
		return err
	}

	env.SetEnv() // PRJ_* + NIX_CONFIG
	if err := syscall.Exec(bash, cmd, os.Environ()); err != nil {
		return err
	}
	return nil
}

func (c *RunActionCmd) getArgs(currentSystem string) ([]string, error) {

	if c.System == currentSystem {
		return nil, fmt.Errorf("set the --for flag to a different system than the current one ('%s')", currentSystem)
	}

	if c.System != "" {
		// if system is set, the impure flag provides a "hack" so that we
		// can transport this information to the action evaluation without
		// incurring in a prohibitively complex (m*n) data structure in
		// which we would have to account for _all_ combinations of current
		// and build system
		return []string{"--impure", c.renderFragmentFor(c.System)}, nil
	}
	return []string{c.renderFragmentFor(currentSystem)}, nil
}

func (c *RunActionCmd) renderFragmentFor(system string) string {
	return tprintf(c, "'"+flakeRegistry(".")+".actions."+system+".\"{{.Cell}}\".\"{{.Block}}\".\"{{.Target}}\".\"{{.Action}}\"'")
}
