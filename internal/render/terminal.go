package render

import (
	"net/http"
	"strings"
)

const termWidth = 72

var htmlStripper = strings.NewReplacer("<h1>", "", "</h1>", "")

type termWrapper struct {
	buf           strings.Builder
	col           int
	pendingSpaces int
}

func (tw *termWrapper) writeText(s string) {
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		if runes[i] == '\n' {
			tw.buf.WriteByte('\n')
			tw.col = 0
			tw.pendingSpaces = 0
			i++
			continue
		}

		spaces := tw.pendingSpaces
		tw.pendingSpaces = 0
		for i < len(runes) && runes[i] == ' ' {
			spaces++
			i++
		}
		if i >= len(runes) || runes[i] == '\n' {
			tw.pendingSpaces = spaces
			continue
		}

		// Collect the next word
		start := i
		for i < len(runes) && runes[i] != ' ' && runes[i] != '\n' {
			i++
		}
		word := string(runes[start:i])
		wordLen := i - start

		if tw.col == 0 {
			tw.buf.WriteString(word)
			tw.col += wordLen
		} else if tw.col+spaces+wordLen <= termWidth {
			tw.buf.WriteString(strings.Repeat(" ", spaces))
			tw.buf.WriteString(word)
			tw.col += spaces + wordLen
		} else {
			tw.buf.WriteByte('\n')
			tw.buf.WriteString(word)
			tw.col = wordLen
		}
	}
}

func RenderTerminal(w http.ResponseWriter, nodes []Node) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	tw := &termWrapper{}
	tw.buf.WriteByte('\n')

	for _, n := range nodes {
		switch n := n.(type) {
		case TextNode:
			tw.writeText(htmlStripper.Replace(n.Content))
		case LinkNode:
			if tw.col > 0 && tw.pendingSpaces > 0 {
				tw.buf.WriteString(strings.Repeat(" ", tw.pendingSpaces))
				tw.col += tw.pendingSpaces
				tw.pendingSpaces = 0
			}
			tw.buf.WriteString("\x1b]8;;" + n.URL + "\x07\x1b[4m")
			tw.writeText(n.Text)
			tw.buf.WriteString("\x1b[24m\x1b]8;;\x07")
		}
	}

	_, _ = w.Write([]byte(tw.buf.String()))
}
