package render

import (
	"fmt"
	"net/http"
	"time"
)

func RenderTerminal(w http.ResponseWriter, f http.Flusher, r *http.Request, nodes []Node) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	_, _ = fmt.Fprint(w, "\n")
	f.Flush()

	stream := newStreamer(16)

	for _, n := range nodes {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		switch n := n.(type) {
		case TextNode:
			streamTextPlain(w, f, r, stream, n.Content)
		case LinkNode:
			_, _ = fmt.Fprintf(w, "\x1b]8;;%s\x07\x1b[4m", n.URL)
			streamTextPlain(w, f, r, stream, n.Text)
			_, _ = fmt.Fprint(w, "\x1b[24m\x1b]8;;\x07")
			f.Flush()
		case StreamNode:
			stream.push(n.Speed)
		case CloseStreamNode:
			stream.pop()
		case PauseNode:
			time.Sleep(time.Duration(n.Duration) * time.Millisecond)
		}
	}
}

func streamTextPlain(w http.ResponseWriter, f http.Flusher, r *http.Request, stream *streamer, text string) {
	for _, ch := range text {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		_, _ = fmt.Fprint(w, string(ch))
		f.Flush()
		if ch != '\n' && ch != ' ' {
			stream.delay()
		}
	}
}
