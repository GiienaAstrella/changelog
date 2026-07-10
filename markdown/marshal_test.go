package markdown_test

import (
	"fmt"
	"testing"

	"giiena.me/changelog/markdown"
	"github.com/stretchr/testify/suite"
)

const payloadString = "Test payload"

type MarshalTestSuite struct {
	suite.Suite
}

type implCommon struct {
	Data string
}

type implMarshaler struct {
	implCommon
}

func (t implMarshaler) MarshalMarkdown() ([]byte, error) {
	return []byte(t.Data), nil
}

type implStringer struct {
	implCommon
}

func (t implStringer) String() string {
	return t.Data
}

type implGoStringer struct {
	implCommon
}

func (t implGoStringer) GoString() string {
	return t.Data
}

func (s *MarshalTestSuite) TestMarshaler() {
	marshaler := implMarshaler{}
	marshaler.Data = payloadString

	s.Require().Implements((*markdown.Marshaler)(nil), marshaler)
	s.Require().NotImplements((*fmt.Stringer)(nil), marshaler)
	s.Require().NotImplements((*fmt.GoStringer)(nil), marshaler)

	md, err := markdown.Marshal(marshaler)
	s.Require().NoError(err)

	s.Equal([]byte(payloadString), md)
}

func (s *MarshalTestSuite) TestStringer() {
	stringer := implStringer{}
	stringer.Data = payloadString

	s.Require().NotImplements((*markdown.Marshaler)(nil), stringer)
	s.Require().Implements((*fmt.Stringer)(nil), stringer)
	s.Require().NotImplements((*fmt.GoStringer)(nil), stringer)

	md, err := markdown.Marshal(stringer)
	s.Require().NoError(err)

	s.Equal([]byte(payloadString), md)
}

func (s *MarshalTestSuite) TestGoStringer() {
	gostringer := implGoStringer{}
	gostringer.Data = payloadString

	s.Require().NotImplements((*markdown.Marshaler)(nil), gostringer)
	s.Require().NotImplements((*fmt.Stringer)(nil), gostringer)
	s.Require().Implements((*fmt.GoStringer)(nil), gostringer)

	md, err := markdown.Marshal(gostringer)
	s.Require().NoError(err)

	s.Equal([]byte(payloadString), md)
}

func (s *MarshalTestSuite) TestFmt() {
	v := implCommon{}
	v.Data = payloadString

	s.Require().NotImplements((*markdown.Marshaler)(nil), v)
	s.Require().NotImplements((*fmt.Stringer)(nil), v)
	s.Require().NotImplements((*fmt.GoStringer)(nil), v)

	md, err := markdown.Marshal(v)
	s.Require().NoError(err)

	s.Equal(fmt.Appendf(nil, "{%s}", payloadString), md)
}

func (s *MarshalTestSuite) TestPointer() {
	ptr := new(implMarshaler)
	ptr.Data = payloadString

	s.Require().Implements((*markdown.Marshaler)(nil), ptr)

	md, err := markdown.Marshal(ptr)
	s.Require().NoError(err)

	s.Equal([]byte(payloadString), md)
}

func TestMarshal(t *testing.T) {
	s := new(MarshalTestSuite)

	suite.Run(t, s)
}
