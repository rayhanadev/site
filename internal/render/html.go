package render

import (
	"fmt"
	"html"
	"net/http"
	"time"
)

// Template holds the HTML template split around the body placeholder.
type Template struct {
	Head, Tail string
}

func RenderHTML(w http.ResponseWriter, f http.Flusher, r *http.Request, tmpl Template, nodes []Node) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, tmpl.Head)
	f.Flush()

	stream := newStreamer(14)

	for _, n := range nodes {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		switch n := n.(type) {
		case TextNode:
			streamTextHTML(w, f, r, stream, n.Content)
		case LinkNode:
			_, _ = fmt.Fprintf(w, `<a href="%s">`, html.EscapeString(n.URL))
			f.Flush()
			streamTextHTML(w, f, r, stream, n.Text)
			_, _ = fmt.Fprint(w, "</a>")
			f.Flush()
		case StreamNode:
			stream.push(n.Speed)
		case CloseStreamNode:
			stream.pop()
		case PauseNode:
			time.Sleep(time.Duration(n.Duration) * time.Millisecond)
		}
	}

	_, _ = fmt.Fprint(w, tmpl.Tail)
	f.Flush()
}

func streamTextHTML(w http.ResponseWriter, f http.Flusher, r *http.Request, stream *streamer, text string) {
	for _, ch := range text {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		_, _ = fmt.Fprint(w, html.EscapeString(string(ch)))
		f.Flush()
		if ch != '\n' && ch != ' ' {
			stream.delay()
		}
	}
}
