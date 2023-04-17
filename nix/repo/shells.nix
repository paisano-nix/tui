/*
This file holds reproducible shells with commands in them.

They conveniently also generate config files in their startup hook.
*/
{
  mdbook = inputs.std.lib.dev.mkShell {nixago = [(inputs.std.lib.cfg.mdbook cell.config.mdbook)];};
  # Tool Homepage: https://numtide.github.io/devshell/
  default = inputs.std.lib.dev.mkShell {
    name = "Paisano TUI";

    # Tool Homepage: https://nix-community.github.io/nixago/
    # This is Standard's devshell integration.
    # It runs the startup hook when entering the shell.
    nixago = [
      inputs.std.lib.cfg.adrgen
      inputs.std.lib.cfg.conform
      (inputs.std.lib.cfg.treefmt cell.config.treefmt)
      (inputs.std.lib.cfg.editorconfig cell.config.editorconfig)
      (inputs.std.lib.cfg.githubsettings cell.config.githubsettings)
      (inputs.std.lib.cfg.lefthook cell.config.lefthook)
      (inputs.std.lib.cfg.mdbook cell.config.mdbook)
    ];

    commands =
      [
        {
          category = "release";
          package = inputs.nixpkgs.cocogitto;
        }
        {
          category = "rendering";
          package = inputs.nixpkgs.mdbook;
        }
        {
          package = inputs.nixpkgs.delve;
          category = "dev";
          name = "dlv";
        }
        {
          package = inputs.nixpkgs.go;
          category = "dev";
        }
        {
          package = inputs.nixpkgs.gotools;
          category = "dev";
        }
        {
          package = inputs.nixpkgs.gopls;
          category = "dev";
        }
      ]
      ++ inputs.nixpkgs.lib.optionals inputs.nixpkgs.stdenv.isLinux [
        {
          package = inputs.nixpkgs.golangci-lint;
          category = "dev";
        }
      ];
  };
}
