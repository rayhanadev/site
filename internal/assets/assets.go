package assets

import _ "embed"

//go:embed content/content.txt
var Content []byte

//go:embed templates/template.html
var TemplateHTML string
