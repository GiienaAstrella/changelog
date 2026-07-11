package keepachangelog

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Changelog represents a chronologically ordered list of notable changes for each version of
// of a project.
type Changelog struct {
	Description      string               `json:"description"`
	DisableLintRules []string             `json:"disable_lint_rules,omitempty"`
	Versions         []Version            `json:"versions"`
	References       map[string]Reference `json:"-"`
}

// String returns the Markdown string for c.
func (c Changelog) String() string {
	var sb strings.Builder

	if len(c.DisableLintRules) > 0 {
		fmt.Fprintf(&sb, "<!-- markdownlint-disable %s -->", strings.Join(c.DisableLintRules, " "))
		sb.WriteString("\n\n")
	}

	sb.WriteString("# Changelog\n\n")
	sb.WriteString(c.Description)
	sb.WriteString("\n\n")

	for _, ver := range c.Versions {
		ver.string(&sb)
	}

	writeRefs(&sb, c.References)

	return strings.TrimSpace(sb.String()) + "\n"
}

// MarshalMarkdown implements [markdown.Marshaler].
//
// Deprecated: use String.
func (c Changelog) MarshalMarkdown() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalMarkdown implements [markdown.Unmarshaler].
//
// Deprecated: use Parse.
func (c *Changelog) UnmarshalMarkdown(data []byte) error {
	return c.unmarshalMarkdown(data)
}

// unmarshalMarkdown decodes a Changelog in Markdown representation from data, storing the parsed
// values in c.
//
// Deprecated: use Parse.
func (c *Changelog) unmarshalMarkdown(data []byte) error {
	normalized := normalize(string(data))

	c.DisableLintRules = make([]string, 0)

	lintRuleGroups := globalLintDisablePattern.FindAllStringSubmatch(normalized, -1)
	for _, ruleGroup := range lintRuleGroups {
		if len(ruleGroup) < 2 {
			continue
		}

		rules := strings.Split(ruleGroup[1], " ")

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)

			if lintRulePattern.MatchString(rule) && !slices.Contains(c.DisableLintRules, rule) {
				c.DisableLintRules = append(c.DisableLintRules, rule)
			}
		}
	}

	titleIndices := titlePattern.FindIndex([]byte(normalized))
	verIndices := verPattern.FindAllIndex([]byte(normalized), -1)

	c.Versions = make([]Version, 0)

	for i, index := range verIndices {
		start := index[0]
		var end int

		if i < len(verIndices)-1 {
			end = verIndices[i+1][0]
		} else {
			end = len(normalized)
		}

		if i == 0 {
			c.Description = strings.TrimSpace(normalized[titleIndices[1]:start])
		}

		ver := &Version{}
		err := ver.UnmarshalMarkdown([]byte(normalized[start:end]))
		if err != nil {
			return err
		}

		c.Versions = append(c.Versions, *ver)
	}

	return nil
}

// Parse parses Changelog from source using an AST-based parser.
func Parse(source []byte) (cl Changelog, err error) {
	md := goldmark.New()
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	cl.References = make(map[string]Reference)

	var version *Version
	var section *Section

	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := node.(type) {
		case *ast.Heading:
			text := extractText(node, source, cl.References)

			switch node.Level {
			case 1:
				return ast.WalkContinue, nil
			case 2:
				v, err := parseVersion(text)
				if err != nil {
					return ast.WalkStop, err
				}
				if version != nil {
					if section != nil {
						version.Sections = append(version.Sections, *section)
						section = nil
					}
					cl.Versions = append(cl.Versions, *version)
				}
				version = new(v)
				version.References = cl.References
			case 3:
				if section != nil && version != nil {
					version.Sections = append(version.Sections, *section)
				}
				section = new(Section)
				section.Heading = string(extractText(node, source, cl.References))
				section.Changes = make([]string, 0)
			}

		case *ast.Paragraph:
			if version == nil {
				p := string(extractMarkdown(node, source, cl.References, true))
				if cl.Description == "" {
					cl.Description = p
				} else {
					cl.Description += "\n\n" + p
				}
			}

		case *ast.List:
			if section == nil {
				return ast.WalkSkipChildren, nil
			}
			for item := node.FirstChild(); item != nil; item = item.NextSibling() {
				li, ok := item.(*ast.ListItem)
				if !ok {
					continue
				}
				var buf bytes.Buffer
				buf.WriteString("- ")
				buf.Write(extractMarkdown(li, source, cl.References, false))
				section.Changes = append(section.Changes, buf.String())
			}
			return ast.WalkSkipChildren, nil

		case *ast.HTMLBlock:
			rules := extractLintRules(node, source)
			if len(cl.DisableLintRules) > 0 {
				for _, rule := range rules {
					if !slices.Contains(cl.DisableLintRules, rule) {
						cl.DisableLintRules = append(cl.DisableLintRules, rule)
					}
				}
			} else {
				cl.DisableLintRules = rules
			}
		}
		return ast.WalkContinue, nil
	})

	if version != nil {
		if section != nil {
			version.Sections = append(version.Sections, *section)
			section = nil
		}
		cl.Versions = append(cl.Versions, *version)
		version = nil
	}

	if cl.DisableLintRules == nil {
		cl.DisableLintRules = []string{}
	}

	if len(cl.References) < 1 {
		cl.References = nil
		for i := range cl.Versions {
			cl.Versions[i].References = nil
		}
	}

	return
}
