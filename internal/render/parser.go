package render

import (
	"strconv"
	"strings"
)

type Node interface {
	node()
}

type TextNode struct{ Content string }
type LinkNode struct{ Text, URL string }
type StreamNode struct{ Speed int }
type CloseStreamNode struct{}
type PauseNode struct{ Duration int }

func (TextNode) node()        {}
func (LinkNode) node()        {}
func (StreamNode) node()      {}
func (CloseStreamNode) node() {}
func (PauseNode) node()       {}

func Parse(input string) []Node {
	runes := []rune(input)
	var nodes []Node
	var buf []rune

	flush := func() {
		if len(buf) > 0 {
			nodes = append(nodes, TextNode{Content: string(buf)})
			buf = buf[:0]
		}
	}

	for i := 0; i < len(runes); {
		if runes[i] != '[' {
			buf = append(buf, runes[i])
			i++
			continue
		}

		node, end := parseTag(runes, i)
		if node == nil {
			node, end = tryLink(runes, i)
		}
		if node != nil {
			flush()
			nodes = append(nodes, node)
			i = end
			continue
		}

		buf = append(buf, runes[i])
		i++
	}

	flush()
	return nodes
}

// tryLink parses a CommonMark-style inline link at runes[pos].
func tryLink(runes []rune, pos int) (Node, int) {
	if pos >= len(runes) || runes[pos] != '[' {
		return nil, pos
	}

	// Link text: balanced brackets, backslash escapes.
	depth := 1
	j := pos + 1
	for j < len(runes) && depth > 0 {
		switch runes[j] {
		case '\\':
			j++
		case '[':
			depth++
		case ']':
			depth--
		}
		j++
	}
	if depth != 0 {
		return nil, pos
	}
	text := string(runes[pos+1 : j-1])

	// ']' must be immediately followed by '('.
	if j >= len(runes) || runes[j] != '(' {
		return nil, pos
	}
	j++

	// Destination: no spaces/control chars, balanced parens.
	destStart := j
	pDepth := 0
	for j < len(runes) {
		ch := runes[j]
		if ch == '(' {
			pDepth++
		} else if ch == ')' {
			if pDepth == 0 {
				break
			}
			pDepth--
		} else if ch <= ' ' {
			break
		}
		j++
	}
	url := string(runes[destStart:j])

	if j >= len(runes) || runes[j] != ')' {
		return nil, pos
	}
	j++

	return LinkNode{Text: text, URL: url}, j
}

const maxTagLen = 20

func parseTag(runes []rune, pos int) (Node, int) {
	if pos >= len(runes) || runes[pos] != '[' {
		return nil, pos
	}

	j := pos + 1
	for j < len(runes) && runes[j] != ']' {
		if j-pos > maxTagLen {
			return nil, pos
		}
		j++
	}
	if j >= len(runes) {
		return nil, pos
	}

	tag := string(runes[pos+1 : j])
	end := j + 1

	if tag == "/stream" {
		return CloseStreamNode{}, end
	}
	if after, ok := strings.CutPrefix(tag, "stream:"); ok {
		if n, err := strconv.Atoi(after); err == nil {
			return StreamNode{Speed: n}, end
		}
	}
	if after, ok := strings.CutPrefix(tag, "pause:"); ok {
		if n, err := strconv.Atoi(after); err == nil {
			return PauseNode{Duration: n}, end
		}
	}

	return nil, pos
}
