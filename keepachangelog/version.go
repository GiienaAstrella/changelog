package keepachangelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

// A Version contains all changes for a given version.
type Version struct {
	ID          string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	Unreleased  bool      `json:"unreleased"`
	Yanked      bool      `json:"yanked"`
	Sections    []Section `json:"contents"`
}

// String returns the Markdown string for v.
func (v Version) String() string {
	var sb strings.Builder
	v.string(&sb)
	return sb.String()
}

// MarshalJSON implements [json.Marshaler].
//
// MarshalJSON skips empty Sections on released versions.
// Unreleased versions will *always* export all Sections.
func (v Version) MarshalJSON() ([]byte, error) {
	type shadowVersion Version
	export := shadowVersion(v)

	if !export.Unreleased {
		delIndices := make([][]int, 0)
		for i, section := range export.Sections {
			if len(section.Changes) < 1 {
				delIndices = append(delIndices, []int{i, i + 1})
			}
		}

		delCnt := 0
		for _, delIndex := range delIndices {
			export.Sections = slices.Delete(export.Sections,
				delIndex[0]-delCnt,
				delIndex[1]-delCnt)
			delCnt++
		}
	}

	return json.Marshal(export)
}

// MarshalMarkdown implements [markdown.Marshaler].
//
// Deprecated: use String.
func (v Version) MarshalMarkdown() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalMarkdown implements [markdown.Unmarshaler].
//
// Deprecated: use Parse.
func (v *Version) UnmarshalMarkdown(data []byte) error {
	return v.unmarshalMarkdown(data)
}

// string encodes v to Markdown, writing into sb.
func (v Version) string(sb *strings.Builder) {
	sb.WriteString("## ")

	if v.Unreleased {
		sb.WriteString("[UNRELEASED]")
	} else {
		fmt.Fprintf(sb, "[%s]", v.ID)
	}

	if !v.Unreleased && !v.ReleaseDate.IsZero() || v.Yanked {
		sb.WriteString(" -")
	}

	if !v.Unreleased && !v.ReleaseDate.IsZero() {
		fmt.Fprintf(sb, " %s", v.ReleaseDate.Format(LayoutChangelog))
	}

	if !v.Unreleased && v.Yanked {
		fmt.Fprint(sb, " [YANKED]")
	}

	sb.WriteString("\n\n")

	for _, content := range v.Sections {
		if !v.Unreleased && len(content.Changes) < 1 {
			continue
		}

		content.string(sb)
	}
}

// unmarshalMarkdown decodes a Version in Markdown representation from data, storing the parsed
// values in v.
//
// Deprecated: use Parse.
func (v *Version) unmarshalMarkdown(data []byte) error {
	var err error
	normalized := normalize(string(data))

	secIndices := secPattern.FindAllIndex([]byte(normalized), -1)

	header := verPattern.FindStringSubmatch(normalized)

	for i, submatch := range header {
		if len(submatch) < 1 {
			continue
		}

		switch i {
		case 0:
			continue

		case 1:
			v.ID = submatch
			if strings.ToLower(submatch) == "unreleased" {
				v.Unreleased = true
			}

		case 2:
			v.ReleaseDate, err = time.Parse(LayoutChangelog, submatch)

		case 3:
			if strings.ToLower(submatch) == "yanked" {
				v.Yanked = true
			}
		}
	}

	v.Sections = make([]Section, 0)

	for i, index := range secIndices {
		start := index[0]
		var end int

		if i < len(secIndices)-1 {
			end = secIndices[i+1][0]
		} else {
			end = len(normalized)
		}

		sec := &Section{}
		sec.UnmarshalMarkdown([]byte(normalized[start:end]))

		v.Sections = append(v.Sections, *sec)
	}

	return err
}

// parseVersion parses Version from heading string.
func parseVersion(heading []byte) (v Version, err error) {
	groups := findNamedSubmatch(np_verPattern, heading)
	if groups == nil {
		err = fmt.Errorf("invalid version heading %q", heading)
		return
	}

	v.ID = string(groups["version"])
	if v.ID[0:1] == "[" {
		v.ID = v.ID[1 : len(v.ID)-1]
	}

	v.Yanked = bytes.EqualFold(groups["yanked"], []byte("YANKED"))

	v.Unreleased = strings.EqualFold(v.ID, "UNRELEASED")
	if date := string(groups["date"]); date != "" {
		v.ReleaseDate, err = time.Parse(LayoutChangelog, string(date))
	}
	return
}
