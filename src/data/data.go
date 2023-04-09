package data

import (
	"fmt"

	"github.com/paisano-nix/paisano/flake"
)

var (
	targetTemplate = "//%s/%s/%s"
	actionTemplate = "//%s/%s/%s:%s"
	noReadme       = "🥺 No Readme available ...\n\n💡 But hey! You could create one ...\n\n💪 Start with: `$EDITOR %s`\n\n👉 It will also be rendered in the docs!"
	noDescription  = "🥺 Target has no 'meta.description' attribute"
)

type Root struct {
	Cells []Cell
}

type Cell struct {
	Name   string  `json:"cell"`
	Readme *string `json:"readme,omitempty"`
	Blocks []Block `json:"cellBlocks"`
}

type Block struct {
	Name      string   `json:"cellBlock"`
	Readme    *string  `json:"readme,omitempty"`
	Blocktype string   `json:"blockType"`
	Targets   []Target `json:"targets"`
}

type Action struct {
	Name  string `json:"name"`
	Descr string `json:"description"`
}

func (a Action) Title() string       { return a.Name }
func (a Action) Description() string { return a.Descr }
func (a Action) FilterValue() string { return a.Title() }

type Target struct {
	Name    string   `json:"name"`
	Readme  *string  `json:"readme,omitempty"`
	Deps    []string `json:"deps"`
	Descr   *string  `json:"description,omitempty"`
	Actions []Action `json:"actions"`
}

func (t Target) Description() string {
	if t.Descr != nil {
		return "💡 " + *t.Descr
	} else {
		return noDescription
	}
}

func (r *Root) Select(ci, bi, ti int) (Cell, Block, Target) {
	var (
		c = r.Cells[ci]
		b = c.Blocks[bi]
		t = b.Targets[ti]
	)
	return c, b, t
}

func (r *Root) ActionArg(ci, bi, ti, ai int) string {
	c, b, t := r.Select(ci, bi, ti)
	a := t.Actions[ai]
	return fmt.Sprintf(actionTemplate, c.Name, b.Name, t.Name, a.Name)
}

func (r *Root) ActionTitle(ci, bi, ti, ai int) string {
	_, _, t := r.Select(ci, bi, ti)
	a := t.Actions[ai]
	return a.Title()
}

func (r *Root) ActionDescription(ci, bi, ti, ai int) string {
	_, _, t := r.Select(ci, bi, ti)
	a := t.Actions[ai]
	return a.Description()
}

func (r *Root) TargetTitle(ci, bi, ti int) string {
	c, b, t := r.Select(ci, bi, ti)
	return fmt.Sprintf(targetTemplate, c.Name, b.Name, t.Name)
}

func (r *Root) TargetDescription(ci, bi, ti int) string {
	_, _, t := r.Select(ci, bi, ti)
	return t.Description()
}
func (r *Root) Cell(ci, bi, ti int) Cell       { c, _, _ := r.Select(ci, bi, ti); return c }
func (r *Root) CellName(ci, bi, ti int) string { return r.Cell(ci, bi, ti).Name }
func (r *Root) CellHelp(ci, bi, ti int) string {
	if r.HasCellHelp(ci, bi, ti) {
		return *r.Cell(ci, bi, ti).Readme
	} else {
		return fmt.Sprintf(noReadme, fmt.Sprintf("%s/%s/Readme.md", flake.CellsFrom.Value(), r.CellName(ci, bi, ti)))
	}
}
func (r *Root) HasCellHelp(ci, bi, ti int) bool {
	c := r.Cell(ci, bi, ti)
	return c.Readme != nil
}
func (r *Root) Block(ci, bi, ti int) Block      { _, o, _ := r.Select(ci, bi, ti); return o }
func (r *Root) BlockName(ci, bi, ti int) string { return r.Block(ci, bi, ti).Name }
func (r *Root) BlockHelp(ci, bi, ti int) string {
	if r.HasBlockHelp(ci, bi, ti) {
		return *r.Block(ci, bi, ti).Readme
	} else {
		return fmt.Sprintf(noReadme, fmt.Sprintf("%s/%s/%s/Readme.md", flake.CellsFrom.Value(), r.CellName(ci, bi, ti), r.BlockName(ci, bi, ti)))
	}
}
func (r *Root) HasBlockHelp(ci, bi, ti int) bool {
	b := r.Block(ci, bi, ti)
	return b.Readme != nil
}
func (r *Root) Target(ci, bi, ti int) Target     { _, _, t := r.Select(ci, bi, ti); return t }
func (r *Root) TargetName(ci, bi, ti int) string { return r.Target(ci, bi, ti).Name }
func (r *Root) TargetHelp(ci, bi, ti int) string {
	if r.HasTargetHelp(ci, bi, ti) {
		return *r.Target(ci, bi, ti).Readme
	} else {
		return fmt.Sprintf(noReadme, fmt.Sprintf("%s/%s/%s/%s.md", flake.CellsFrom.Value(), r.CellName(ci, bi, ti), r.BlockName(ci, bi, ti), r.TargetName(ci, bi, ti)))
	}
}
func (r *Root) HasTargetHelp(ci, bi, ti int) bool {
	t := r.Target(ci, bi, ti)
	return t.Readme != nil
}

func (r *Root) Len() int {
	sum := 0
	for _, c := range r.Cells {
		for _, o := range c.Blocks {
			sum += len(o.Targets)
		}
	}
	return sum
}
