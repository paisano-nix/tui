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

    vendorHash = "sha256-ja0nFWdWqieq8m6cSKAhE1ibeN0fODDCC837jw0eCnE=";

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
