package local

import "strings"

type Query struct {
	Text  string
	Limit int
}

func NewQuery(text string, limit int) Query {
	text = strings.TrimSpace(text)
	if limit <= 0 {
		limit = 20
	}

	return Query{
		Text:  text,
		Limit: limit,
	}
}

func (q Query) Empty() bool {
	return strings.TrimSpace(q.Text) == ""
}
