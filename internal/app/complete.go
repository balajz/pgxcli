package app

import (
	"sort"
	"strings"

	"github.com/balajz/bubbline"
	"github.com/balajz/bubbline/computil"
	"github.com/balajz/bubbline/editline"
	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxls/pkg/engine"
	"github.com/balajz/pgxls/pkg/types"
)

const maxCompletions = 8

func (p *pgxCLI) getCompletions() bubbline.AutoCompleteFn {
	go func() {
		err := p.client.Cache(p.compWorker)
		if err != nil {
			p.logger.Error("error caching database schema for completions", "error", err)
			p.logger.Debug("continuing with keyword completions")
		}
		p.logger.Debug("completion cache is ready")
	}()

	compEngine := engine.NewCompleter(p.compWorker.Cache())

	return func(v [][]rune, line, col int) (msg string, comps bubbline.Completions) {

		sql, _ := computil.Flatten(v, line, col)
		word, wstart, wend := computil.FindWord(v, line, col)

		if strings.HasPrefix(word, "\\") && strings.Index(sql, "\\") == 0 {
			return completeMetaCommand(word, col, wstart, wend, maxCompletions)
		}

		compEngine.DBCache = p.compWorker.Cache()

		items, err := compEngine.Complete(sql, line, col, true)
		if err != nil || len(items) == 0 {
			return "", nil
		}

		if lastDot := strings.LastIndex(word, "."); lastDot != -1 {
			wstart += lastDot + 1
		}

		compByCategory := make(map[string][]compCandidate)
		for _, row := range items {

			category := "others"

			switch row.Kind {
			case types.KeywordCompletion:
				category = "keywords"
			case types.FunctionCompletion:
				category = "functions"
			case types.ClassCompletion:
				category = "tables"
			case types.FieldCompletion:
				category = "columns"
			case types.ModuleCompletion:
				category = "schemas"
			case types.SnippetCompletion:
				category = "snippets"
			}

			docs := ""
			if row.Documentation != nil {
				docs = row.Documentation.Value
			}

			compByCategory[category] = append(compByCategory[category], compCandidate{
				completion: row.Label,
				desc:       row.Detail,
				moveRight:  wend - col,
				deleteLeft: wend - wstart,
				docs:       docs,
			})

		}

		if len(compByCategory) == 0 {
			return msg, comps
		}

		limitCandidate(compByCategory, maxCompletions)

		categories := make([]string, 0, len(compByCategory))
		for k := range compByCategory {
			categories = append(categories, k)
		}
		sort.Strings(categories)
		comps = &completions{
			categories:  categories,
			compEntries: compByCategory,
		}
		return msg, comps
	}
}

func completeMetaCommand(s string, col, wStart, wEnd, limit int) (string, bubbline.Completions) {
	cmds := pgxspecial.Export()

	var matches struct {
		cmds []string
		desc []string
	}

	for _, cmd := range cmds {
		if strings.HasPrefix(cmd.Cmd, s) {
			matches.cmds = append(matches.cmds, cmd.Cmd)
			matches.desc = append(matches.desc, cmd.Description)
		}

		if len(matches.cmds) >= limit {
			break
		}
	}

	if len(matches.cmds) == 0 {
		return "", nil
	}

	return "", editline.SimpleWordsCompletionWithDescriptions(
		matches.cmds, matches.desc, "commands", col, wStart, wEnd,
	)
}

// https://github.com/cockroachdb/cockroach/blob/refs/heads/master/pkg/cli/clisqlshell/complete.go

// completions is the interface between the shell and the bubbline
// completion infra.
type completions struct {
	categories  []string
	compEntries map[string][]compCandidate
}

var _ bubbline.Completions = (*completions)(nil)

// NumCategories is part of the bubbline.Completions interface.
func (c *completions) NumCategories() int { return len(c.categories) }

// CategoryTitle is part of the bubbline.Completions interface.
func (c *completions) CategoryTitle(cIdx int) string { return c.categories[cIdx] }

// NumEntries is part of the bubbline.Completions interface.
func (c *completions) NumEntries(cIdx int) int { return len(c.compEntries[c.categories[cIdx]]) }

// Entry is part of the bubbline.Completions interface.
func (c *completions) Entry(cIdx, eIdx int) bubbline.Entry {
	return &c.compEntries[c.categories[cIdx]][eIdx]
}

// Candidate is part of the bubbline.Completions interface.
func (c *completions) Candidate(e bubbline.Entry) bubbline.Candidate { return e.(*compCandidate) }

// compCandidate represents one completion candidate.
type compCandidate struct {
	completion string
	desc       string
	moveRight  int
	deleteLeft int
	docs       string
}

var _ bubbline.Entry = (*compCandidate)(nil)

// Title is part of the bubbline.Entry interface.
func (c *compCandidate) Title() string { return c.completion }

// Description is part of the bubbline.Entry interface.
func (c *compCandidate) Description() string { return c.desc }

// Replacement is part of the bubbline.Candidate interface.
func (c *compCandidate) Replacement() string { return c.completion }

// MoveRight is part of the bubbline.Candidate interface.
func (c *compCandidate) MoveRight() int { return c.moveRight }

// DeleteLeft is part of the bubbline.Candidate interface.
func (c *compCandidate) DeleteLeft() int { return c.deleteLeft }

func (c *compCandidate) SidePanel() string { return c.docs }

func limitCandidate(comp map[string][]compCandidate, limit int) {
	for category, comps := range comp {
		if len(comps) > limit {
			comp[category] = comps[:limit]
		}
	}
}
