{
  description = "Paisano's TUI/CLI companion";

  inputs.std.url = "github:divnix/std";
  inputs.nixpkgs.follows = "std/nixpkgs";

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
