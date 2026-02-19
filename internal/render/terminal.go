package render

import (
	"fmt"
	"net/http"
)

func RenderTerminal(w http.ResponseWriter, nodes []Node) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	_, _ = fmt.Fprint(w, "\n")

	for _, n := range nodes {
		switch n := n.(type) {
		case TextNode:
			_, _ = fmt.Fprint(w, n.Content)
		case LinkNode:
			_, _ = fmt.Fprintf(w, "\x1b]8;;%s\x07\x1b[4m%s\x1b[24m\x1b]8;;\x07", n.URL, n.Text)
		}
	}
}
