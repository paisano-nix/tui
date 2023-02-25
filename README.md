# The `paisano` TUI / CLI

## Usage

- Install Paisano: `nix profile install github:paisano-nix/tui`
- Set up autocompletion (optional): `paisano _carapace [SHELL]` &mdash; see [carapace docs][carapace-docs]
- Enter a Paisano-based repository.
- Run `paisano` or `paisano list` and profit ✨!

[carapace-docs]: https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html

## Branding

To change the branding of this binary you can set these variables via `-X` compile flag:

```
main.buildVersion | default: dev
main.buildCommit  | default: dirty
main.argv0        | default: paisano
main.project      | default: Paisano
flake.registry    | default: __std   # temp kept, mainly for `std-action`
env.dotdir        | default: .std    # temp kept, for not rewriting many .gitignore
```

Example: `go build -o my-bin-name -ldflags="-X main.argv0=hive -X main.project=Hive"`

## Contributing

##### Prerequisites

You need [nix](https://nixos.org/download.html) and [direnv](https://direnv.net/).

##### Enter Contribution Environment

```console
direnv allow
```

##### Change Contribution Environment

```console
$EDITOR ./nix/repo/config.nix
direnv reload
```

##### Preview Documentation

<sub>You need to be <i>inside</i> the Contribution Environment.</sub>

```console
mdbook build -o
```
