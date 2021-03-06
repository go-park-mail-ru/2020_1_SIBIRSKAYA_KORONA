// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(in *jlexer.Lexer, out *Templates) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Templates, 0, 1)
			} else {
				*out = Templates{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Task
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(out *jwriter.Writer, in Templates) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Templates) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Templates) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Templates) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Templates) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(l, v)
}
func easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(in *jlexer.Lexer, out *Template) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "template":
			out.Variant = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(out *jwriter.Writer, in Template) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"template\":"
		out.RawString(prefix[1:])
		out.String(string(in.Variant))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Template) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Template) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDff6627eEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Template) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Template) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDff6627eDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(l, v)
}
