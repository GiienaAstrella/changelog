package keepachangelog_test

import (
	"encoding/json"
	"slices"
	"testing"
	"time"

	"giiena.me/changelog/keepachangelog"
	"giiena.me/changelog/markdown"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const description = `This is a test Changelog.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0).`

type ChangelogTestSuite struct {
	suite.Suite
}

func (s *ChangelogTestSuite) TestUnreleased() {
	c := keepachangelog.Changelog{
		Description:      description,
		DisableLintRules: []string{},
		Versions: []keepachangelog.Version{
			{
				ID:         "UNRELEASED",
				Unreleased: true,
				Sections: []keepachangelog.Section{
					{
						Heading: "Added",
						Changes: []string{},
					},
					{
						Heading: "Changed",
						Changes: []string{
							"- Default preset now uses XXH3_64.",
						},
					},
				},
			},
		},
	}

	marshaledMD := `# Changelog

This is a test Changelog.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0).

## [UNRELEASED]

### Added

### Changed

- Default preset now uses XXH3_64.
`

	s.Run("MarshalMarkdown", func() {
		md, err := markdown.Marshal(c)
		s.Require().NoError(err)

		s.Equal(marshaledMD, string(md))
	})

	s.Run("UnmarshalMarkdown", func() {
		var cl keepachangelog.Changelog
		err := markdown.Unmarshal([]byte(marshaledMD), &cl)
		s.Require().NoError(err)

		s.Equal(c, cl)
	})
}

func (s *ChangelogTestSuite) TestReleased() {
	c := keepachangelog.Changelog{
		Description:      description,
		DisableLintRules: []string{},
		Versions: []keepachangelog.Version{
			{
				ID:          "0.1.0",
				ReleaseDate: dateMustParse(s.T(), keepachangelog.LayoutChangelog, "2023-08-24"),
				Sections: []keepachangelog.Section{
					{
						Heading: "Added",
						Changes: []string{
							"- Added support for multiple target directories.",
						},
					},
					{
						Heading: "Changed",
						Changes: []string{},
					},
					{
						Heading: "Removed",
						Changes: []string{
							"- Removed -v parameter.",
						},
					},
				},
			},
		},
	}

	marshaledMD := `# Changelog

This is a test Changelog.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0).

## [0.1.0] - 2023-08-24

### Added

- Added support for multiple target directories.

### Removed

- Removed -v parameter.
`
	marshaledJSON := `{"description":"This is a test Changelog.\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0).","versions":[{"name":"0.1.0","release_date":"2023-08-24T00:00:00Z","unreleased":false,"yanked":false,"contents":[{"heading":"Added","changes":["- Added support for multiple target directories."]},{"heading":"Removed","changes":["- Removed -v parameter."]}]}]}`

	s.Run("MarshalMarkdown", func() {
		md, err := markdown.Marshal(c)
		s.Require().NoError(err)

		s.Equal(marshaledMD, string(md))
	})

	s.Run("UnmarshalMarkdown", func() {
		var cl keepachangelog.Changelog
		err := markdown.Unmarshal([]byte(marshaledMD), &cl)
		s.Require().NoError(err)

		s.Equal(sanitizeCL(s.T(), c), cl)
	})

	s.Run("MarshalJSON", func() {
		j, err := json.Marshal(c)
		s.Require().NoError(err)

		s.Equal(marshaledJSON, string(j))
	})
}

func (s *ChangelogTestSuite) TestYanked() {
	c := keepachangelog.Changelog{
		Description:      description,
		DisableLintRules: []string{},
		Versions: []keepachangelog.Version{
			{
				ID:          "0.2.0",
				Yanked:      true,
				ReleaseDate: dateMustParse(s.T(), keepachangelog.LayoutChangelog, "2023-08-24"),
				Sections: []keepachangelog.Section{
					{
						Heading: "Added",
						Changes: []string{
							"- Added support for multiple target directories.",
						},
					},
				},
			},
		},
	}

	marshaledMD := `# Changelog

This is a test Changelog.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0).

## [0.2.0] - 2023-08-24 [YANKED]

### Added

- Added support for multiple target directories.
`

	s.Run("MarshalMarkdown", func() {
		md, err := markdown.Marshal(c)
		s.Require().NoError(err)

		s.Equal(marshaledMD, string(md))
	})

	s.Run("UnmarshalMarkdown", func() {
		var cl keepachangelog.Changelog
		err := markdown.Unmarshal([]byte(marshaledMD), &cl)
		s.Require().NoError(err)

		s.Equal(c, cl)
	})
}

func (s *ChangelogTestSuite) TestDisableLintRules() {
	c := keepachangelog.Changelog{
		Description: "",
		DisableLintRules: []string{
			"MD024",
			"MD009",
		},
		Versions: []keepachangelog.Version{},
	}

	marshaledMD := `<!-- markdownlint-disable MD024 MD009 -->

# Changelog
`

	s.Run("MarshalMarkdown", func() {
		md, err := markdown.Marshal(c)
		s.Require().NoError(err)

		s.Equal(marshaledMD, string(md))
	})

	s.Run("UnmarshalMarkdown", func() {
		var cl keepachangelog.Changelog
		err := markdown.Unmarshal([]byte(marshaledMD), &cl)
		s.Require().NoError(err)

		s.ElementsMatch(c.DisableLintRules, cl.DisableLintRules)
	})
}

func TestChangelog(t *testing.T) {
	s := new(ChangelogTestSuite)

	suite.Run(t, s)
}

func dateMustParse(t testing.TB, layout, values string) time.Time {
	t.Helper()

	d, err := time.Parse(layout, values)
	require.NoError(t, err)

	return d
}

func sanitizeCL(t testing.TB, cl keepachangelog.Changelog) keepachangelog.Changelog {
	t.Helper()

	for v, version := range cl.Versions {
		if version.Unreleased {
			continue
		}

		delIndices := make([][]int, 0)

		for s, section := range version.Sections {
			if len(section.Changes) < 1 {
				delIndices = append(delIndices, []int{s, s + 1})
			}
		}

		delCnt := 0
		for _, delIndex := range delIndices {
			version.Sections = slices.Delete(version.Sections,
				delIndex[0]-delCnt,
				delIndex[1]-delCnt)
			delCnt++
		}

		cl.Versions[v] = version
	}

	return cl
}
