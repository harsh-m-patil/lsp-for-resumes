package analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/harsh-m-patil/ats-optimus-prime-lsp/lsp"
)

type State struct {
	// filename -> content
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

var verbs []string = []string{
	"achieved", "analyzed", "built", "communicated", "created", "designed", "developed", "executed", "helped", "improved", "increased", "led", "managed", "monitored", "organized", "reduced", "researched", "resolved", "tested", "worked on",
}

func getSynonymsForActionVerb(verb string) []string {
	actionSynonyms := map[string][]string{
		"developed": {
			"engineered",
			"implemented",
			"programmed",
			"designed",
			"created",
		},
		"built":     {"assembled", "constructed", "established", "produced", "devised"},
		"created":   {"innovated", "initiated", "authored", "originated", "designed"},
		"designed":  {"conceptualized", "crafted", "modeled", "planned", "drafted"},
		"led":       {"directed", "oversaw", "coordinated", "supervised", "mentored"},
		"managed":   {"orchestrated", "guided", "administered", "facilitated", "delegated"},
		"worked on": {"collaborated on", "contributed to", "participated in", "executed", "supported"},
		"helped":    {"assisted", "aided", "advised", "mentored", "supported"},
		"communicated": {
			"presented",
			"articulated",
			"conveyed",
			"documented",
			"reported",
		},
		"researched": {"investigated", "explored", "examined", "studied", "analyzed"},
		"analyzed":   {"assessed", "evaluated", "interpreted", "measured", "reviewed"},
		"tested":     {"validated", "verified", "benchmarked", "debugged", "inspected"},
		"improved":   {"optimized", "refined", "enhanced", "streamlined", "modernized"},
		"achieved":   {"accomplished", "delivered", "attained", "completed", "secured"},
		"increased":  {"boosted", "expanded", "elevated", "amplified", "augmented"},
		"reduced":    {"decreased", "cut", "minimized", "streamlined", "eliminated"},
		"organized": {
			"arranged",
			"structured",
			"scheduled",
			"systematized",
			"prioritized",
		},
		"executed":  {"implemented", "performed", "carried out", "fulfilled", "realized"},
		"monitored": {"tracked", "observed", "supervised", "evaluated", "checked"},
		"resolved":  {"fixed", "solved", "rectified", "addressed", "handled"},
	}

	return actionSynonyms[verb]
}

type FreqData struct {
	count      int
	ocurrences []OccurrenceData
}

type OccurrenceData struct {
	line  int
	start int
	end   int
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	freqMap := map[string]FreqData{}
	re := regexp.MustCompile(`^\s*-[^0-9]*$`)

	for row, line := range strings.Split(text, "\n") {
		lower := strings.ToLower(line)
		if re.MatchString(lower) {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, 0, len(line)),
				Severity: 2,
				Source:   "Resume Tips",
				Message:  "No impact",
			})
		}
		for _, verb := range verbs {
			idx := strings.Index(lower, verb)

			if idx >= 0 {
				data, exists := freqMap[verb]
				if !exists {
					data = FreqData{
						count:      0,
						ocurrences: []OccurrenceData{},
					}
				}

				data.count += 1
				data.ocurrences = append(data.ocurrences, OccurrenceData{
					line:  row,
					start: idx,
					end:   idx + len(verb),
				})

				freqMap[verb] = data
			}
		}
	}

	for verb, data := range freqMap {
		if data.count > 2 {
			for _, occ := range data.ocurrences {
				diagnostics = append(diagnostics, lsp.Diagnostic{
					Range:    LineRange(occ.line, occ.start, occ.end),
					Severity: 2,
					Source:   "Resume Tips",
					Message:  fmt.Sprintf("The verb '%s' is overused (%d times). Consider using synonyms.", verb, data.count),
				})
			}
		}
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// in real life, this would look up the type in our type analysis code

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// in real life, this would look up the definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")

		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS C*de",
				Edit: &lsp.WorkspaceEdit{
					Changes: censorChange,
				},
			})
		}
	}

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}

	return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "https://neovim.io",
		},
	}

	// ask your static analysis engine for good completion items
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
