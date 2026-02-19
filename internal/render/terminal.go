package render

import (
	"fmt"
	"net/http"
)

const termWidth = 72

type termWrapper struct {
	w   http.ResponseWriter
	col int
}

func (tw *termWrapper) writeText(s string) {
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		if runes[i] == '\n' {
			fmt.Fprint(tw.w, "\n")
			tw.col = 0
			i++
			continue
		}

		// Count spaces
		spaces := 0
		for i < len(runes) && runes[i] == ' ' {
			spaces++
			i++
		}
		if i >= len(runes) || runes[i] == '\n' {
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
			fmt.Fprint(tw.w, word)
			tw.col += wordLen
		} else if tw.col+spaces+wordLen <= termWidth {
			for j := 0; j < spaces; j++ {
				fmt.Fprint(tw.w, " ")
			}
			fmt.Fprint(tw.w, word)
			tw.col += spaces + wordLen
		} else {
			fmt.Fprint(tw.w, "\n")
			fmt.Fprint(tw.w, word)
			tw.col = wordLen
		}
	}
}

func RenderTerminal(w http.ResponseWriter, nodes []Node) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	tw := &termWrapper{w: w}
	fmt.Fprint(w, "\n")

	for _, n := range nodes {
		switch n := n.(type) {
		case TextNode:
			tw.writeText(n.Content)
		case LinkNode:
			fmt.Fprintf(w, "\x1b]8;;%s\x07\x1b[4m", n.URL)
			tw.writeText(n.Text)
			fmt.Fprint(w, "\x1b[24m\x1b]8;;\x07")
		}
	}
}
