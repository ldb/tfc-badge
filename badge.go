package main

import (
	_ "embed"
	"fmt"
	"golang.org/x/text/width"
	"html/template"
	"io"
	"math"
)

type Badge struct {
	Color, URL, Message string
	Width               int
}

const (
	MessageWidthMultiplier = 6.5
	LabelOffset            = 82
	MessageDefaultWidth    = 180
	MessageMargin          = 0
)

var DefaultBadge = Badge{
	Color:   "#7B42BC",
	URL:     "https://www.terraform.io/cloud",
	Message: "unknown",
	Width:   MessageDefaultWidth,
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
	b.Width = LabelOffset + int(math.Ceil(float64(GetWidthUTF8String(n.Message))*MessageWidthMultiplier)) + MessageMargin
	if b.Width < MessageDefaultWidth {
		b.Width = MessageDefaultWidth
	}
	b.URL = r.RunURL
}

// Render renders a Badge into SVG format.
func (b *Badge) Render(out io.Writer) error {
	if err := tmpl.Execute(out, b); err != nil {
		return fmt.Errorf("error rendering badge: %v", err)
	}
	return nil
}

func GetWidthUTF8String(s string) int {
	size := 0
	for _, runeValue := range s {
		p := width.LookupRune(runeValue)
		if p.Kind() == width.EastAsianWide {
			size += 2
			continue
		}
		size += 1

	}
	return size
}
