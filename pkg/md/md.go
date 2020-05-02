package md

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type Elem interface {
	String() string
}

type JoinType struct {
	elems []Elem
	with  string
}

func (p JoinType) String() string {
	s := ""
	for _, e := range p.elems {
		s += p.with + e.String()
	}
	return strings.TrimPrefix(s, p.with)
}

func Paragraph(elems ...Elem) JoinType {
	return JoinType{elems: elems, with: " "}
}

func Doc(elems ...Elem) JoinType {
	return JoinType{elems: elems, with: "\n\n"}
}

type TextType struct {
	content string
}

func (t TextType) String() string {
	return t.content
}

func Text(text string) TextType {
	return TextType{content: text}
}

type HeadlineType struct {
	level   int
	content string
}

func (h HeadlineType) String() string {
	return strings.Repeat("#", h.level) + " " + h.content
}

func Headline(level int, content string) HeadlineType {
	return HeadlineType{
		level:   level,
		content: content,
	}
}

type SurroundType struct {
	body     Elem
	surround string
}

func (s SurroundType) String() string {
	return s.surround + s.body.String() + s.surround
}

func Bold(e Elem) SurroundType {
	return SurroundType{body: e, surround: "**"}
}

func Italic(e Elem) SurroundType {
	return SurroundType{body: e, surround: "*"}
}

func Code(e Elem) SurroundType {
	return SurroundType{body: e, surround: "`"}
}

type CodeBlockType struct {
	lang    string
	snippet string
}

func (c CodeBlockType) String() string {
	return fmt.Sprintf("```%s\n%s\n```", c.lang, c.snippet)
}

func CodeBlock(lang, snippet string) CodeBlockType {
	return CodeBlockType{lang: lang, snippet: snippet}
}

type ListType struct {
	elems []Elem
}

func (l ListType) String() string {
	s := ""
	for _, e := range l.elems {
		switch t := e.(type) {
		case ListType:
			s += "\n  " + strings.Join(strings.Split(t.String(), "\n"), "\n  ")
		default:
			s += "\n* " + t.String()
		}
	}
	return strings.TrimPrefix(s, "\n")
}

func List(elems ...Elem) ListType {
	return ListType{elems: elems}
}

type LinkType struct {
	desc Elem
	href string
}

func (l LinkType) String() string {
	return fmt.Sprintf("[%s](%s)", l.desc.String(), l.href)
}

func Link(desc Elem, href string) LinkType {
	return LinkType{
		desc: desc,
		href: href,
	}
}

type FrontmatterType struct {
	yaml string
}

func (f FrontmatterType) String() string {
	return "---\n" + f.yaml + "---"
}

func Frontmatter(data map[string]interface{}) FrontmatterType {
	d, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}

	return FrontmatterType{yaml: string(d)}
}
