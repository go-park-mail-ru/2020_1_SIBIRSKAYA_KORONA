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

func easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(in *jlexer.Lexer, out *Board) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "title":
			out.Name = string(in.String())
		case "columns":
			if in.IsNull() {
				in.Skip()
				out.Columns = nil
			} else {
				in.Delim('[')
				if out.Columns == nil {
					if !in.IsDelim(']') {
						out.Columns = make([]Column, 0, 1)
					} else {
						out.Columns = []Column{}
					}
				} else {
					out.Columns = (out.Columns)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Column
					easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(in, &v1)
					out.Columns = append(out.Columns, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "admins":
			if in.IsNull() {
				in.Skip()
				out.Admins = nil
			} else {
				in.Delim('[')
				if out.Admins == nil {
					if !in.IsDelim(']') {
						out.Admins = make([]User, 0, 1)
					} else {
						out.Admins = []User{}
					}
				} else {
					out.Admins = (out.Admins)[:0]
				}
				for !in.IsDelim(']') {
					var v2 User
					easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(in, &v2)
					out.Admins = append(out.Admins, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "members":
			if in.IsNull() {
				in.Skip()
				out.Members = nil
			} else {
				in.Delim('[')
				if out.Members == nil {
					if !in.IsDelim(']') {
						out.Members = make([]User, 0, 1)
					} else {
						out.Members = []User{}
					}
				} else {
					out.Members = (out.Members)[:0]
				}
				for !in.IsDelim(']') {
					var v3 User
					easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(in, &v3)
					out.Members = append(out.Members, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(out *jwriter.Writer, in Board) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	if len(in.Columns) != 0 {
		const prefix string = ",\"columns\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v4, v5 := range in.Columns {
				if v4 > 0 {
					out.RawByte(',')
				}
				easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(out, v5)
			}
			out.RawByte(']')
		}
	}
	if len(in.Admins) != 0 {
		const prefix string = ",\"admins\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v6, v7 := range in.Admins {
				if v6 > 0 {
					out.RawByte(',')
				}
				easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(out, v7)
			}
			out.RawByte(']')
		}
	}
	if len(in.Members) != 0 {
		const prefix string = ",\"members\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v8, v9 := range in.Members {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(out, v9)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Board) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Board) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Board) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Board) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels(l, v)
}
func easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(in *jlexer.Lexer, out *User) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "password":
			if in.IsNull() {
				in.Skip()
				out.Password = nil
			} else {
				out.Password = in.Bytes()
			}
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
func easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"surname\":"
		out.RawString(prefix)
		out.String(string(in.Surname))
	}
	{
		const prefix string = ",\"nickname\":"
		out.RawString(prefix)
		out.String(string(in.Nickname))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	if len(in.Password) != 0 {
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Password)
	}
	out.RawByte('}')
}
func easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(in *jlexer.Lexer, out *Column) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "title":
			out.Name = string(in.String())
		case "position":
			out.Pos = float64(in.Float64())
		case "tasks":
			if in.IsNull() {
				in.Skip()
				out.Tasks = nil
			} else {
				in.Delim('[')
				if out.Tasks == nil {
					if !in.IsDelim(']') {
						out.Tasks = make([]Task, 0, 1)
					} else {
						out.Tasks = []Task{}
					}
				} else {
					out.Tasks = (out.Tasks)[:0]
				}
				for !in.IsDelim(']') {
					var v13 Task
					easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels3(in, &v13)
					out.Tasks = append(out.Tasks, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels1(out *jwriter.Writer, in Column) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"position\":"
		out.RawString(prefix)
		out.Float64(float64(in.Pos))
	}
	if len(in.Tasks) != 0 {
		const prefix string = ",\"tasks\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v14, v15 := range in.Tasks {
				if v14 > 0 {
					out.RawByte(',')
				}
				easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels3(out, v15)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels3(in *jlexer.Lexer, out *Task) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "title":
			out.Name = string(in.String())
		case "description":
			out.About = string(in.String())
		case "level":
			out.Level = uint(in.Uint())
		case "deadline":
			out.Deadline = string(in.String())
		case "position":
			out.Pos = float64(in.Float64())
		case "cid":
			out.Cid = uint(in.Uint())
		case "members":
			if in.IsNull() {
				in.Skip()
				out.Members = nil
			} else {
				in.Delim('[')
				if out.Members == nil {
					if !in.IsDelim(']') {
						out.Members = make([]User, 0, 1)
					} else {
						out.Members = []User{}
					}
				} else {
					out.Members = (out.Members)[:0]
				}
				for !in.IsDelim(']') {
					var v16 User
					easyjson202377feDecodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(in, &v16)
					out.Members = append(out.Members, v16)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels3(out *jwriter.Writer, in Task) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	if in.Level != 0 {
		const prefix string = ",\"level\":"
		out.RawString(prefix)
		out.Uint(uint(in.Level))
	}
	if in.Deadline != "" {
		const prefix string = ",\"deadline\":"
		out.RawString(prefix)
		out.String(string(in.Deadline))
	}
	{
		const prefix string = ",\"position\":"
		out.RawString(prefix)
		out.Float64(float64(in.Pos))
	}
	{
		const prefix string = ",\"cid\":"
		out.RawString(prefix)
		out.Uint(uint(in.Cid))
	}
	if len(in.Members) != 0 {
		const prefix string = ",\"members\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v17, v18 := range in.Members {
				if v17 > 0 {
					out.RawByte(',')
				}
				easyjson202377feEncodeGithubComGoParkMailRu20201SIBIRSKAYAKORONAAppModels2(out, v18)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}