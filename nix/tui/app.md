# `paisano` CLI / TUI

```console
❯ paisano --help
paisano is the CLI / TUI companion for Paisano.

- Invoke without any arguments to start the TUI.
- Invoke with a target spec and action to run a known target's action directly.

Enable autocompletion via 'paisano _carapace <shell>'.
For more instructions, see: https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html

Usage:
  paisano //[cell]/[block]/[target]:[action] [args...]
  paisano [command]

Available Commands:
  check       Validate the repository.
  list        List available targets.
  re-cache    Refresh the CLI cache.

Flags:
      --for string   system, for which the target will be built (e.g. 'x86_64-linux')
  -h, --help         help for paisano
  -v, --version      version for paisano

Use "paisano [command] --help" for more information about a command.

```
