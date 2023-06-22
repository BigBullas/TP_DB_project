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

func easyjsonDdc53814DecodeGithubComBigBullasTPDBProjectInternalModels(in *jlexer.Lexer, out *Info) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.Users = int64(in.Int64())
		case "forum":
			out.Forums = int64(in.Int64())
		case "thread":
			out.Threads = int64(in.Int64())
		case "post":
			out.Posts = int64(in.Int64())
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
func easyjsonDdc53814EncodeGithubComBigBullasTPDBProjectInternalModels(out *jwriter.Writer, in Info) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Users))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.Int64(int64(in.Forums))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int64(int64(in.Threads))
	}
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix)
		out.Int64(int64(in.Posts))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Info) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDdc53814EncodeGithubComBigBullasTPDBProjectInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Info) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDdc53814EncodeGithubComBigBullasTPDBProjectInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Info) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDdc53814DecodeGithubComBigBullasTPDBProjectInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Info) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDdc53814DecodeGithubComBigBullasTPDBProjectInternalModels(l, v)
}