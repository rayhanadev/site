package render

import "strings"

var botPatterns = []string{
	"bot",
	"crawl",
	"spider",
	"slurp",
	"sogou",
	"ia_archiver",
	"chatgpt-user",
	"google-extended",
	"anthropic-ai",
	"claude-web",
	"cohere-ai",
	"facebookexternalhit",
	"whatsapp",
}

var terminalPatterns = []string{
	"curl",
	"wget",
	"httpie",
}

func IsBot(ua string) bool {
	lower := strings.ToLower(ua)
	for _, p := range botPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func IsTerminal(ua string) bool {
	lower := strings.ToLower(ua)
	for _, p := range terminalPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
