{
  # A flake-view into the top-level flake's contrib env
  # with `mdbook-paisano-preprocessor' available in inputs
  description = "Paisano's TUI/CLI extended docs env";

  inputs.super.url = "path:../.";
  inputs.mdbook-paisano-preprocessor.url = "github:paisano-nix/mdbook-paisano-preprocessor";
  inputs.std.follows = "super/std";
  inputs.nixpkgs.follows = "super/std/nixpkgs";

  outputs = {
    std,
    self,
    super,
    ...
  } @ inputs:
    std.growOn {
      inherit inputs;
      cellsFrom = std.incl ../nix [
        "repo"
      ];
      cellBlocks = with std.blockTypes; [
        # Development Environments
        (nixago "config")
        (devshells "shells")
      ];
    }
    {
      devShells = std.harvest self ["repo" "shells"];
    };
}
