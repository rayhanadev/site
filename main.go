package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/rayhanadev/site/internal/assets"
	"github.com/rayhanadev/site/internal/render"
)

func main() {
	content := string(assets.Content)
	nodes := render.Parse(content)

	parts := strings.SplitN(assets.TemplateHTML, "{{BODY}}", 2)
	if len(parts) != 2 {
		log.Fatal("template.html missing {{BODY}} placeholder")
	}
	tmpl := render.Template{Head: parts[0], Tail: parts[1]}

	addr := envOr("LISTEN_ADDR", ":3000")
	certPath := envOr("TLS_CERT_PATH", "configs/certs/cert.pem")
	keyPath := envOr("TLS_KEY_PATH", "configs/certs/key.pem")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Alt-Svc", `h3="`+addr+`"; ma=86400`)
		accept := r.Header.Get("Accept")

		if strings.Contains(accept, "text/html") {
			render.RenderHTML(w, flusher, r, tmpl, nodes)
		} else {
			render.RenderTerminal(w, flusher, r, nodes)
		}
	})

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	tc := &tls.Config{Certificates: []tls.Certificate{cert}}

	h3 := &http3.Server{Addr: addr, Handler: mux, TLSConfig: tc}
	h2 := &http.Server{
		Addr:              addr,
		Handler:           mux,
		TLSConfig:         tc,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Cert/key paths empty because TLSConfig already contains the loaded certificate.
	go func() { log.Fatal(h2.ListenAndServeTLS("", "")) }()
	log.Println("listening on", addr)
	log.Fatal(h3.ListenAndServe())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
