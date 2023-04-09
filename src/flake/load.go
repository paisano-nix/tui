package flake

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/paisano-nix/paisano/cache"
	"github.com/paisano-nix/paisano/env"
)

func LoadFlakeCmd() (*cache.Cache, *cache.ActionID, *exec.Cmd, *bytes.Buffer, error) {

	nix, err := getNix()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	currentSystem, err := getCurrentSystem()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// load the paisano metadata from the flake
	buf := new(bytes.Buffer)
	args := []string{
		"eval",
		"--json",
		"--no-update-lock-file",
		"--no-write-lock-file",
		"--no-warn-dirty",
		"--accept-flake-config",
		flakeRegistry(".") + ".init." + currentSystem}
	cmd := exec.Command(nix, args...)
	cmd.Stdin = devNull
	cmd.Stdout = buf

	// initialize cache
	metadataCacheDir, err := env.GetProjectMetadataCacheDir()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	c, err := cache.Open(metadataCacheDir)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	key := cache.NewActionID([]byte(strings.Join(args, "")))

	return c, &key, cmd, buf, nil
}
