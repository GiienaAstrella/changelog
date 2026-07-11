package keepachangelog

import (
	"bytes"
	"regexp"
	"slices"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// LayoutChangelog is a [time.Layout] string.
//
// [Keep a Changelog] dates are formatted as [ISO 8601] (YYYY-MM-DD).
//
// [Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
// [ISO 8601]: https://www.iso.org/iso-8601-date-and-time-format.html
const LayoutChangelog = "2006-01-02"

var (
	titlePattern                = regexp.MustCompile(`(?m)^#[\t ]+Changelog[\t ]*$`)
	verPattern                  = regexp.MustCompile(`(?m)^#{2}[\t ]+\[([^\[\]\s]+)\](?:[\t ]+\-){0,1}(?:[\t ]+([0-9]{4}\-[0-9]{2}\-[0-9]{2})){0,1}(?:[\t ]+\[(YANKED){1}\]){0,1}[\t ]*`)
	np_verPattern               = regexp.MustCompile(`(?m)^(?P<version>[^\s]+)(?:(?:[\t ]+\-[\t ]+)(?P<date>\d{4}\-\d{2}\-\d{2})(?:(?:[\t ]+)(?:\[(?P<yanked>.+)\]))?)?`)
	secPattern                  = regexp.MustCompile(`(?m)^#{3}[\t ]+(.+)[\t ]*$`)
	globalLintDisablePattern    = regexp.MustCompile(`(?m)^<!--[\t ]+markdownlint-disable[\t ]+(.*)[\t ]+-->$`)
	np_globalLintDisablePattern = regexp.MustCompile(`(?m)^<!--[\t ]+markdownlint-disable[\t ]+(?P<rules>.*)[\t ]+-->$`)
	lintRulePattern             = regexp.MustCompile(`MD[0-9]{3}`)
)

// normalize normalizes str, replacing `\r\n` and `\r` to `\n`.
// Effectively, normalize converts all line endings to LF.
func normalize(str string) string {
	normalized := strings.ReplaceAll(string(str), "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	return normalized
}

// extractText extracts text from node.
func extractText(node ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n == node {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			if entering {
				buf.Write(n.Segment.Value(source))
				if n.SoftLineBreak() || n.HardLineBreak() {
					buf.WriteRune(' ')
				}
			}
		}
		return ast.WalkContinue, nil
	})
	return buf.Bytes()
}

// extractMarkdown extracts the raw Markdown string from node.
func extractMarkdown(node ast.Node, source []byte, _ any, preserveLineBreak bool) []byte {
	var buf bytes.Buffer
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n == node {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			if entering {
				buf.Write(n.Segment.Value(source))
				if !preserveLineBreak {
					if n.SoftLineBreak() || n.HardLineBreak() {
						buf.WriteRune(' ')
					}
				} else {
					if n.SoftLineBreak() {
						buf.WriteRune('\n')
					} else if n.HardLineBreak() {
						buf.WriteString("  \n")
					}
				}
			}

		case *ast.CodeSpan:
			if entering {
				buf.WriteRune('`')
			} else {
				buf.WriteRune('`')
			}

		case *ast.CodeBlock, *ast.FencedCodeBlock:
			return ast.WalkSkipChildren, nil

		case *ast.Emphasis:
			marker := "*"
			if n.Level == 2 {
				marker = "**"
			}
			buf.WriteString(marker)
			if !entering {
				return ast.WalkContinue, nil
			}

		case *ast.AutoLink:
			if entering {
				buf.WriteRune('<')
				buf.Write(n.URL(source))
				buf.WriteRune('>')
			}

		case *ast.Link:
			if entering {
				buf.WriteRune('[')
				return ast.WalkContinue, nil
			}

			buf.WriteRune(']')

			buf.WriteRune('(')
			buf.Write(n.Destination)
			if len(n.Title) > 0 {
				buf.WriteString(` "`)
				buf.Write(n.Title)
				buf.WriteRune('"')
			}
			buf.WriteRune(')')
			return ast.WalkContinue, nil
		}
		return ast.WalkContinue, nil
	})
	return buf.Bytes()
}

// extractLintRules extracts markdownlint rules from an HTML block.
func extractLintRules(node *ast.HTMLBlock, source []byte) []string {
	raw := extractHTMLBlock(node, source)

	groups := findNamedSubmatch(np_globalLintDisablePattern, raw)
	if groups == nil {
		return []string{}
	}

	rules := strings.Split(string(groups["rules"]), " ")
	validRules := make([]string, 0, len(rules))
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if lintRulePattern.MatchString(rule) && !slices.Contains(validRules, rule) {
			validRules = append(validRules, rule)
		}
	}
	return validRules
}

// extractHTMLBlock extracts the raw text from an HTML block.
func extractHTMLBlock(node *ast.HTMLBlock, source []byte) []byte {
	var buf bytes.Buffer
	for i := 0; i < node.Lines().Len(); i++ {
		seg := node.Lines().At(i)
		buf.Write(seg.Value(source))
	}
	if node.HasClosure() {
		buf.Write(node.ClosureLine.Value(source))
	}
	return buf.Bytes()
}

// findNameSubmatch maps the result of pattern.FindSubmatch to their named capture groups.
func findNamedSubmatch(pattern *regexp.Regexp, b []byte) map[string][]byte {
	matches := pattern.FindSubmatch(b)
	if matches == nil {
		return nil
	}

	groups := make(map[string][]byte, len(matches))
	for i, name := range pattern.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		groups[name] = matches[i]
	}
	return groups
}
