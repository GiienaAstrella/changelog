package markdown

import "fmt"

// Unmarshaler is the interface implemented by types that can unmarshal data from a Markdown
// formatted representation.
//
// Deprecated: for [keepachangelog], use [keepachangelog.Parse].
// The new AST parser is more efficient and extendible than the previous text parser.
type Unmarshaler interface {
	UnmarshalMarkdown(data []byte) error
}

// Unmarshal parses a Markdown formatted text and stores the result in the value pointed by v.
//
// If v implements [Unmarshaler], Unmarshal calls [Unmarshaler.UnmarshalMarkdown].
// Otherwise, Unmarshal returns an error.
//
// Deprecated: for [keepachangelog], use [keepachangelog.Parse].
// The new AST parser is more efficient and extendible than the previous text parser.
func Unmarshal(data []byte, v any) error {
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalMarkdown(data)
	}

	return fmt.Errorf("%T does not implement markdown.Unmarshaler", v)
}
