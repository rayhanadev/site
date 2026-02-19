package render

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

type Template struct {
	Head, Tail string
}

// iOS Safari won't start rendering until it receives ~1KB of renderable data.
// HTML comments don't count toward this threshold, so we use zero-width spaces
// (U+200B) instead. Each is 3 bytes in UTF-8, so 342 characters = 1026 bytes,
// just over the 1KB threshold.
var safariPad = strings.Repeat("\u200B", 342)

func RenderHTML(w http.ResponseWriter, f http.Flusher, r *http.Request, tmpl Template, nodes []Node) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, tmpl.Head)
	_, _ = fmt.Fprint(w, safariPad)
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
			_, _ = fmt.Fprintf(w, `<a href="%s"%s>`, html.EscapeString(n.URL), linkTarget(n.URL))
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
		case ClipNode:
			_, _ = fmt.Fprint(w, `<span style="display:block;overflow:hidden;white-space:nowrap">`)
			f.Flush()
		case CloseClipNode:
			_, _ = fmt.Fprint(w, "</span>")
			f.Flush()
		}
	}

	_, _ = fmt.Fprint(w, tmpl.Tail)
	f.Flush()
}

func RenderHTMLInstant(w http.ResponseWriter, tmpl Template, nodes []Node) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, tmpl.Head)

	for _, n := range nodes {
		switch n := n.(type) {
		case TextNode:
			_, _ = fmt.Fprint(w, n.Content)
		case LinkNode:
			_, _ = fmt.Fprintf(w, `<a href="%s"%s>%s</a>`, html.EscapeString(n.URL), linkTarget(n.URL), n.Text)
		case ClipNode:
			_, _ = fmt.Fprint(w, `<span style="display:block;overflow:hidden;white-space:nowrap">`)
		case CloseClipNode:
			_, _ = fmt.Fprint(w, "</span>")
		}
	}

	_, _ = fmt.Fprint(w, tmpl.Tail)
}

func linkTarget(url string) string {
	if strings.Contains(url, "ring.purduehackers.com") {
		return ""
	}
	return ` target="_blank" rel="noopener noreferrer"`
}

func streamTextHTML(w http.ResponseWriter, f http.Flusher, r *http.Request, stream *streamer, text string) {
	runes := []rune(text)
	for i := 0; i < len(runes); {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		if runes[i] == '<' {
			j := i + 1
			for j < len(runes) && runes[j] != '>' {
				j++
			}
			if j < len(runes) {
				j++
			}
			_, _ = fmt.Fprint(w, string(runes[i:j]))
			f.Flush()
			i = j
			continue
		}

		_, _ = fmt.Fprint(w, string(runes[i]))
		f.Flush()
		if runes[i] != '\n' && runes[i] != ' ' {
			stream.delay()
		}
		i++
	}
}
