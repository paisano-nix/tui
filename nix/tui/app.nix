let
  version = "0.15.0+dev";

  inherit (inputs) nixpkgs;
  inherit (nixpkgs.lib) licenses;
in {
  default = cell.app.paisano;

  paisano = nixpkgs.buildGoModule rec {
    inherit version;
    pname = "paisano";
    meta = {
      inherit (import (inputs.self + /flake.nix)) description;
      license = licenses.unlicense;
      homepage = "https://github.com/paisano-nix/tui";
    };

    src = inputs.self + /src;

    vendorHash = "sha256-1le14dcr2b8TDUNdhIFbZGX3khQoCcEZRH86eqlZaQE=";

    nativeBuildInputs = [nixpkgs.installShellFiles];

    postInstall = ''
      installShellCompletion --cmd paisano \
        --bash <($out/bin/paisano _carapace bash) \
        --fish <($out/bin/paisano _carapace fish) \
        --zsh <($out/bin/paisano _carapace zsh)
    '';

    ldflags = [
      "-s"
      "-w"
      "-X main.buildVersion=${version}"
    ];
  };
}
