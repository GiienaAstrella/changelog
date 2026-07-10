package markdown_test

import (
	"testing"

	"giiena.me/changelog/markdown"
	"github.com/stretchr/testify/suite"
)

type UnmarshalTestSuite struct {
	suite.Suite
}

func (s *UnmarshalTestSuite) TestUnmarshaler() {
	var unmarshaler implUnmarshaler

	s.Require().Implements((*markdown.Unmarshaler)(nil), &unmarshaler)

	err := markdown.Unmarshal([]byte(payloadString), &unmarshaler)
	s.Require().NoError(err)

	s.Equal(payloadString, unmarshaler.Data)
}

func (s *UnmarshalTestSuite) TestNoImpl() {
	var common implCommon

	s.Require().NotImplements((*markdown.Unmarshaler)(nil), &common)

	err := markdown.Unmarshal([]byte(payloadString), &common)
	s.Error(err)
}

type implUnmarshaler struct {
	implCommon
}

func (t *implUnmarshaler) UnmarshalMarkdown(data []byte) error {
	t.Data = string(data)

	return nil
}

func TestUnmarshal(t *testing.T) {
	s := new(UnmarshalTestSuite)

	suite.Run(t, s)
}
