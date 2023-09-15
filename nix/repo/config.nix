/*
This file holds configuration data for repo dotfiles.

Q: Why not just put the put the file there?

A: (1) dotfile proliferation
   (2) have all the things in one place / fromat
   (3) potentially share / re-use configuration data - keeping it in sync
*/
let
  inherit (inputs) nixpkgs;
  inherit (inputs.std.inputs) dmerge;
  inherit (inputs.std.data) configs;
  inherit (inputs.std.lib.dev) mkNixago;
in {
  # Tool Homepage: https://numtide.github.io/treefmt/
  treefmt = (mkNixago configs.treefmt) {
    packages = [nixpkgs.go];
    data = {
      formatter = {
        go = {
          command = "gofmt";
          options = ["-w"];
          includes = ["*.go"];
        };
      };
    };
  };

  # Tool Homepage: https://editorconfig.org/
  editorconfig = (mkNixago configs.editorconfig) {};
  conform = (mkNixago configs.conform) {};

  # Tool Homepage: https://github.com/evilmartians/lefthook
  lefthook = (mkNixago configs.lefthook) {};

  cog = (mkNixago configs.cog) {
    data = {
      changelog = {
        remote = "github.com";
        repository = "tui";
        owner = "paisano-nix";
      };
      post_bump_hooks = dmerge.append [
        "echo Go to and post: https://discourse.nixos.org/t/paisano-tui-cli/27351"
      ];
    };
  };

  # Tool Hompeage: https://github.com/apps/settings
  # Install Setting App in your repo to enable it
  githubsettings = (mkNixago configs.githubsettings) {
    data = {
      repository = {
        name = "tui";
        inherit (import (inputs.self + /flake.nix)) description;
        homepage = "https://paisano-nix.github.io/tui";
        topics = "nix, nix-flakes, flake, ux, tui, cli";
        default_branch = "main";
        allow_squash_merge = true;
        allow_merge_commit = false;
        allow_rebase_merge = true;
        delete_branch_on_merge = true;
        private = false;
        has_issues = true;
        has_projects = false;
        has_wiki = false;
        has_downloads = false;
      };
    };
  };

  # Tool Homepage: https://rust-lang.github.io/mdBook/
  mdbook = (mkNixago configs.mdbook) {
    # add preprocessor packages here
    packages = [nixpkgs.mdbook-linkcheck];
    data = {
      # Configuration Reference: https://rust-lang.github.io/mdBook/format/configuration/index.html
      book.title = "Paisano TUI Book";
      preprocessor.paisano-preprocessor = {
        multi = [
          {
            chapter = "TUI Reference";
            cell = "tui";
          }
        ];
      };
      output = {
        html = {
          additional-css = ["./mdbook-paisano-preprocessor.css"];
        };
        # Tool Homepage: https://github.com/Michael-F-Bryan/mdbook-linkcheck
        linkcheck = {};
      };
    };
  };
}
