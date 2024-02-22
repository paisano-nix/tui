package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/oriser/regroup"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"

	"github.com/paisano-nix/paisano/data"
	"github.com/paisano-nix/paisano/flake"
)

type Spec struct {
	Cell   string `regroup:"cell,required"`
	Block  string `regroup:"block,required"`
	Target string `regroup:"target,required"`
	Action string `regroup:"action,required"`
}

var re = regroup.MustCompile(`^//(?P<cell>[^/]+)/(?P<block>[^/]+)/(?P<target>.+):(?P<action>[^:]+)`)

var forSystem string

var rootCmd = &cobra.Command{
	Use:                   fmt.Sprintf("%[1]s //[cell]/[block]/[target]:[action] [args...]", argv0),
	DisableFlagsInUseLine: true,
	Version:               fmt.Sprintf("%s (%s)", buildVersion, buildCommit),
	Short:                 fmt.Sprintf("%[1]s is the CLI / TUI companion for %[2]s", argv0, project),
	Long: fmt.Sprintf(`%[1]s is the CLI / TUI companion for %[2]s.

- Invoke without any arguments to start the TUI.
- Invoke with a target spec and action to run a known target's action directly.

Enable autocompletion via '%[1]s _carapace <shell>'.
For more instructions, see: https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html
`, argv0, project),
	Args: func(cmd *cobra.Command, args []string) error {
		s := &Spec{}
		if err := re.MatchToTarget(args[0], s); err != nil {
			return fmt.Errorf("invalid argument format: %s", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s := &Spec{}
		if err := re.MatchToTarget(args[0], s); err != nil {
			return err
		}
		command := flake.RunActionCmd{
			ShowCmdStr: false,
			CmdStr:     strings.Join(args[:], " "),
			System:     forSystem,
			Cell:       s.Cell,
			Block:      s.Block,
			Target:     s.Target,
			Action:     s.Action}
		if err := command.Exec(args[1:]); err != nil {
			return err
		}
		return nil

	},
}
var reCacheCmd = &cobra.Command{
	Use:   "re-cache",
	Short: "Refresh the CLI cache.",
	Long: `Refresh the CLI cache.
Use this command to cold-start or refresh the CLI cache.
The TUI does this automatically, but the command completion needs manual initialization of the CLI cache.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, key, loadCmd, buf, err := flake.LoadFlakeCmd()
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		loadCmd.Run()
		c.PutBytes(*key, buf.Bytes())
		return nil
	},
}
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate the repository.",
	Long: fmt.Sprintf(`Validates that the repository conforms to %[1]s.
Returns a non-zero exit code and an error message if the repository is not a valid %[1]s repository.
The TUI does this automatically.`, project),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, loadCmd, _, err := flake.LoadFlakeCmd()
		loadCmd.Args = append(loadCmd.Args, "--trace-verbose")
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		loadCmd.Stderr = os.Stderr
		if err := loadCmd.Run(); err != nil {
			os.Exit(1)
		}
		fmt.Printf("Valid %s repository ✓\n", project)

		return nil
	},
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available targets.",
	Long: `List available targets.
Shows a list of all available targets. Can be used as an alternative to the TUI.
Also loads the CLI cache, if no cache is found. Reads the cache, otherwise.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cache, key, loadCmd, buf, err := flake.LoadFlakeCmd()
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		cached, _, err := cache.GetBytes(*key)
		var root *data.Root
		if err == nil {
			root, err = LoadJson(bytes.NewReader(cached))
			if err != nil {
				return fmt.Errorf("while loading cached json: %w", err)
			}
		} else {
			loadCmd.Run()
			bufA := &bytes.Buffer{}
			r := io.TeeReader(buf, bufA)
			root, err = LoadJson(r)
			if err != nil {
				return fmt.Errorf("while loading json (cmd: '%v'): %w", loadCmd, err)
			}
			cache.PutBytes(*key, bufA.Bytes())
		}
		w := tabwriter.NewWriter(os.Stdout, 5, 2, 4, ' ', 0)
		for _, c := range root.Cells {
			for _, o := range c.Blocks {
				for _, t := range o.Targets {
					for _, a := range t.Actions {
						fmt.Fprintf(w, "//%s/%s/%s:%s\t--\t%s:  %s\n", c.Name, o.Name, t.Name, a.Name, t.Description(), a.Description())
					}
				}
			}
		}
		w.Flush()
		return nil
	},
}

func ExecuteCli() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&forSystem, "for", "", "system, for which the target will be built (e.g. 'x86_64-linux')")
	rootCmd.AddCommand(reCacheCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(checkCmd)
	carapace.Gen(rootCmd).Standalone()
	// completes: '//cell/block/target:action'
	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionCallback(func(ctx carapace.Context) carapace.Action {
			cache, key, _, _, err := flake.LoadFlakeCmd()
			if err != nil {
				return carapace.ActionMessage(fmt.Sprintf("%v\n", err))
			}
			cached, _, err := cache.GetBytes(*key)
			var root *data.Root
			if err == nil {
				root, err = LoadJson(bytes.NewReader(cached))
				if err != nil {
					return carapace.ActionMessage(fmt.Sprintf("%v\n", err))
				}
			} else {
				return carapace.ActionMessage(fmt.Sprintf("No completion cache: please initialize by running '%[1]s re-cache'.", argv0))
			}
			var cells = []string{}
			var blocks = map[string][]string{}
			var targets = map[string]map[string][]string{}
			var actions = map[string]map[string]map[string][]string{}
			for _, c := range root.Cells {
				blocks[c.Name] = []string{}
				targets[c.Name] = map[string][]string{}
				actions[c.Name] = map[string]map[string][]string{}
				cells = append(cells, c.Name, "cell")
				for _, b := range c.Blocks {
					targets[c.Name][b.Name] = []string{}
					actions[c.Name][b.Name] = map[string][]string{}
					blocks[c.Name] = append(blocks[c.Name], b.Name, "block")
					for _, t := range b.Targets {
						actions[c.Name][b.Name][t.Name] = []string{}
						targets[c.Name][b.Name] = append(targets[c.Name][b.Name], t.Name, t.Description())
						for _, a := range t.Actions {
							actions[c.Name][b.Name][t.Name] = append(
								actions[c.Name][b.Name][t.Name],
								a.Name,
								a.Description(),
							)
						}
					}
				}
			}
			return carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				// start with <tab>; no typing
				case 0:
					return carapace.ActionValuesDescribed(
						cells...,
					).Invoke(c).Prefix("//").Suffix("/").ToA().Style(
						style.Of(style.Bold, style.Carapace.Highlight(1)))
				// only a single / typed
				case 1:
					return carapace.ActionValuesDescribed(
						cells...,
					).Invoke(c).Prefix("/").Suffix("/").ToA()
				// start typing cell
				case 2:
					return carapace.ActionValuesDescribed(
						cells...,
					).Invoke(c).Suffix("/").ToA().Style(
						style.Carapace.Highlight(1))
				// start typing block
				case 3:
					return carapace.ActionValuesDescribed(
						blocks[c.Parts[2]]...,
					).Invoke(c).Suffix("/").ToA().Style(
						style.Carapace.Highlight(2))
				// start typing target
				case 4:
					return carapace.ActionMultiParts(":", func(d carapace.Context) carapace.Action {
						switch len(d.Parts) {
						// start typing target
						case 0:
							return carapace.ActionValuesDescribed(
								targets[c.Parts[2]][c.Parts[3]]...,
							).Invoke(c).Suffix(":").ToA().Style(
								style.Carapace.Highlight(3))
							// start typing action
						case 1:
							return carapace.ActionValuesDescribed(
								actions[c.Parts[2]][c.Parts[3]][d.Parts[0]]...,
							).Invoke(c).ToA()
						default:
							return carapace.ActionValues()
						}
					})
				default:
					return carapace.ActionValues()
				}
			})

		}),
	)
}
