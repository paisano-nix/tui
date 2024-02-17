package env

import (
	"fmt"
	"os"
	"strings"

	spec "github.com/numtide/prj-spec/contrib/go"
)

// extraNixConfig implements quality of life flags for the nix command invocation
var extraNixConfig = strings.Join([]string{
	// can never occur: actions invoke store path copies of the flake
	// "warn-dirty = false",
	"accept-flake-config = true",
	"builders-use-substitutes = true",
	// TODO: these are unfortunately not available for setting as env flags
	// update-lock-file = false,
	// write-lock-file = false,
}, "\n")

func SetEnv() {
	spec.SetAll() // PRJ_*

	nixConfigEnv, present := os.LookupEnv("NIX_CONFIG")
	if !present {
		os.Setenv("NIX_CONFIG", extraNixConfig)
	} else {
		os.Setenv("NIX_CONFIG", fmt.Sprintf("%s\n%s", nixConfigEnv, extraNixConfig))
	}
}

func GetStateActionPath() (string, error) { return spec.DataFile("last-action") }
func GetProjectMetadataCacheDir() (string, error) {
	path, err := spec.CacheFile("metadata")
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}
