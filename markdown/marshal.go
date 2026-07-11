package markdown

import "fmt"

// Marshaler is the interface implemented by types that can marshal themselves into a Markdown
// formatted representation.
//
// Deprecated: implement [fmt.Stringer] instead.
// Markdown strings are just strings.
type Marshaler interface {
	MarshalMarkdown() ([]byte, error)
}

// Marshal returns the Markdown formatted representation of v.
//
// If v implements [Marshaler], Marshal calls [Marshaler.MarshalMarkdown] to produce the Markdown.
// If v implements [fmt.Stringer] instead, Marshal calls [fmt.Stringer.String].
// If v implements [fmt.GoStringer] instead, Marshal calls [fmt.GoStringer.GoString].
// Otherwise, Marshal calls [fmt.Sprintf] to format v as is.
//
// Deprecated: use the type's String function instead.
func Marshal(v any) ([]byte, error) {
	if marshaler, ok := v.(Marshaler); ok {
		return marshaler.MarshalMarkdown()
	} else if stringer, ok := v.(fmt.Stringer); ok {
		return []byte(stringer.String()), nil
	} else if gostringer, ok := v.(fmt.GoStringer); ok {
		return []byte(gostringer.GoString()), nil
	} else {
		return fmt.Appendf(nil, "%v", v), nil
	}
}
