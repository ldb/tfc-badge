package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

type Badge struct {
	Color, URL, Message string
	Width               int
}

const MessageWidthMultiplier = 15

var DefaultBadge = Badge{
	Color:   "#7B42BC",
	URL:     "https://www.terraform.io/cloud",
	Message: "unknown",
	Width:   150,
}

//go:embed badge.tmpl
var badgeTemplate string
var tmpl = template.Must(template.New("badge").Parse(badgeTemplate))
var colorMap = map[string]string{
	"run:created":         "#7b42bc",
	"run:planning":        "#1563ff",
	"run:needs_attention": "#fa8f37",
	"run:applying":        "#1563ff",
	"run:completed":       "#2eb039",
	"run:errored":         "#c73445",
}

func (b *Badge) FromRun(r *Run) {
	n := r.Notifications[0]
	color, ok := colorMap[n.Trigger]
	if !ok {
		color = DefaultBadge.Color
	}
	b.Color = color
	b.Message = n.Message
	b.Width = len(n.Message) * MessageWidthMultiplier
		b.URL = r.RunURL
}

// Render renders a Badge into SVG format.
func (b *Badge) Render(out io.Writer) error {
	if err := tmpl.Execute(out, b); err != nil {
		return fmt.Errorf("error rendering badge: %v", err)
	}
	return nil
}
