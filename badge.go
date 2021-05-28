package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

type Badge struct {
	Color, URL, Message string
}

var DefaultBadge = Badge{
	Color:   "#7B42BC",
	URL:     "https://www.terraform.io/cloud",
	Message: "unknown",
}

//go:embed badge.tmpl
var badgeTemplate string
var tmpl = template.Must(template.New("badge").Parse(badgeTemplate))
var colorMap = map[string]string{
	"run:created":         "",
	"run:planning":        "",
	"run:needs_attention": "",
	"run:applying":        "",
	"run:completed":       "#123456",
	"run:errored":         "",
}

func (b *Badge) FromRun(r *Run) {
	n := r.Notifcations[0]
	color, ok := colorMap[n.Trigger]
	if !ok {
		color = DefaultBadge.Color
	}
	b.Color = color
	b.Message = n.Message
	b.URL = r.RunURL
}

// Render renders a Badge into SVG format.
func (b *Badge) Render(out io.Writer) error {
	if err := tmpl.Execute(out, b); err != nil {
		return fmt.Errorf("error rendering badge: %v", err)
	}
	return nil
}
