package keepachangelog

import (
	"fmt"
	"strings"
)

// A Section is a group of similar changes.
//
// For example, a Section may be a group of `Added` for new features, `Changed` for changes,
// `Deprecated` for outdated to-be-removed features, `Removed` for removed features, etc.
type Section struct {
	Heading string   `json:"heading"`
	Changes []string `json:"changes"`
}

// String returns the Markdown string for s.
func (s Section) String() string {
	var sb strings.Builder
	s.string(&sb, false)
	return sb.String()
}

// SingleLineString returns the Markdown string for s.
// Unlike String, sentences are not broken into multiple lines.
func (s Section) SingleLineString() string {
	var sb strings.Builder
	s.string(&sb, true)
	return sb.String()
}

// MarshalMarkdown implements [markdown.Marshaler].
//
// Deprecated: use String.
func (s Section) MarshalMarkdown() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalMarkdown implements [markdown.Unmarshaler].
//
// Deprecated: use Parse.
func (s *Section) UnmarshalMarkdown(data []byte) error {
	return s.unmarshalMarkdown(data)
}

// string encodes s to Markdown, writing into sb.
func (s Section) string(sb *strings.Builder, collapse bool) {
	fmt.Fprintf(sb, "### %s\n\n", s.Heading)

	for _, change := range s.Changes {
		if !collapse {
			lines := strings.Split(change, ". ")
			change = strings.Join(lines, ".\n"+prefix(change))
		}

		fmt.Fprintf(sb, "%s\n", change)
	}

	if len(s.Changes) > 0 {
		sb.WriteString("\n")
	}
}

// unmarshalMarkdown decodes a Section in Markdown representation from data, storing the parsed
// values in s.
//
// Deprecated: use Parse.
func (s *Section) unmarshalMarkdown(data []byte) error {
	normalized := normalize(string(data))

	header := secPattern.FindStringSubmatch(normalized)
	headerIndices := secPattern.FindIndex([]byte(normalized))

	s.Changes = make([]string, 0)

	for i, submatch := range header {
		switch i {
		case 0:
			continue

		case 1:
			s.Heading = submatch
		}
	}

	lines := strings.Split(normalized[headerIndices[1]:], "\n")

	var entry string
	for _, line := range lines {
		original := line
		line = strings.TrimSpace(line)

		if len(line) < 1 {
			continue
		}

		if line[:1] == "-" {
			if len(entry) > 0 {
				s.Changes = append(s.Changes, strings.TrimRight(entry, "\t "))
			}

			entry = original
		} else {
			entry += fmt.Sprintf(" %s", line)
		}
	}

	if len(entry) > 0 {
		s.Changes = append(s.Changes, strings.TrimRight(entry, "\t "))
	}

	return nil
}

// prefix returns the indentation prefix of line.
// This is useful when splitting sentences into multiple lines.
func prefix(line string) string {
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	leadingTabs := len(line) - len(strings.TrimLeft(line, "\t"))

	line = strings.TrimLeft(line, "\t ")
	if line[:1] == "-" {
		leadingSpaces += 2
	}

	return strings.Repeat("\t", leadingTabs) + strings.Repeat(" ", leadingSpaces)
}
