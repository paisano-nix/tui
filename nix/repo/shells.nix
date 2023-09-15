/*
This file holds reproducible shells with commands in them.

They conveniently also generate config files in their startup hook.
*/
let
  inherit (inputs) nixpkgs;
  inherit (inputs.std.lib.dev) mkShell;
  inherit (inputs.nixpkgs.lib) mapAttrs optionals;
  # Tool Homepage: https://numtide.github.io/devshell/
in
  mapAttrs (_: mkShell) rec {
    mdbook.nixago = [cell.config.mdbook];
    default = {
      name = "Paisano TUI";

      # Tool Homepage: https://nix-community.github.io/nixago/
      # This is Standard's devshell integration.
      # It runs the startup hook when entering the shell.
      nixago = [
        cell.config.conform
        cell.config.treefmt
        cell.config.editorconfig
        cell.config.githubsettings
        cell.config.lefthook
        cell.config.mdbook
        cell.config.cog
      ];

      commands =
        [
          {
            package = nixpkgs.delve;
            category = "dev";
            name = "dlv";
          }
          {
            package = nixpkgs.go;
            category = "dev";
          }
          {
            package = nixpkgs.gotools;
            category = "dev";
          }
          {
            package = nixpkgs.gopls;
            category = "dev";
          }
        ]
        ++ optionals nixpkgs.stdenv.isLinux [
          {
            package = nixpkgs.golangci-lint;
            category = "dev";
          }
        ];
    };
  }
