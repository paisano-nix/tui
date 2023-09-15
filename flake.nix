{
  description = "Paisano's TUI/CLI companion";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  inputs.std = {
    url = "github:divnix/std";
    inputs.nixpkgs.follows = "nixpkgs";
    inputs.devshell.follows = "devshell";
    inputs.nixago.follows = "nixago";
  };

  inputs.devshell = {
    url = "github:numtide/devshell";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  inputs.nixago = {
    url = "github:nix-community/nixago";
    inputs.nixpkgs.follows = "nixpkgs";
    inputs.nixago-exts.follows = "";
  };

  outputs = {
    std,
    self,
    ...
  } @ inputs:
    std.growOn {
      inherit inputs;
      cellsFrom = ./nix;
      cellBlocks = with std.blockTypes; [
        # Development Environments
        (nixago "config")
        (devshells "shells")
        # Application Development
        (installables "app")
      ];
    }
    {
      packages = std.harvest self ["tui" "app"];
      devShells = std.harvest self ["repo" "shells"];
    };
}
