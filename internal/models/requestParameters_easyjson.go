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

func easyjson410d6225DecodeGithubComBigBullasTPDBProjectInternalModels(in *jlexer.Lexer, out *RequestParameters) {
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
		case "desc":
			out.Desc = bool(in.Bool())
		case "limit":
			out.Limit = int(in.Int())
		case "since":
			out.Since = string(in.String())
		case "sort":
			out.Sort = string(in.String())
		case "sinceInt":
			out.SinceInt = int(in.Int())
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
func easyjson410d6225EncodeGithubComBigBullasTPDBProjectInternalModels(out *jwriter.Writer, in RequestParameters) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"desc\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Desc))
	}
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix)
		out.Int(int(in.Limit))
	}
	{
		const prefix string = ",\"since\":"
		out.RawString(prefix)
		out.String(string(in.Since))
	}
	{
		const prefix string = ",\"sort\":"
		out.RawString(prefix)
		out.String(string(in.Sort))
	}
	{
		const prefix string = ",\"sinceInt\":"
		out.RawString(prefix)
		out.Int(int(in.SinceInt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestParameters) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson410d6225EncodeGithubComBigBullasTPDBProjectInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestParameters) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson410d6225EncodeGithubComBigBullasTPDBProjectInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestParameters) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson410d6225DecodeGithubComBigBullasTPDBProjectInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestParameters) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson410d6225DecodeGithubComBigBullasTPDBProjectInternalModels(l, v)
}
