package v5

import (
	"bytes"
	"fmt"
	"github.com/pkg/browser"
	"html/template"
	"strings"
)

// Main just used for linking up the main functions
func Main(args []string) {
	switch args[0] {
	case "draw":
		Draw(args[1])
	default:
		fmt.Println("command not recognized")
	}
}

func Draw(input string) {
	tokens := lex(input)
	parser := NewParser()
	ast := parser.Parse(tokens)
	head, _ := ast.compile()

	stateDrawing := head.Draw()

	t, err := template.New("graph").Parse(`
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});
</script>
<div class="mermaid">
    {{ . }}
</div>
<div>
<span style="white-space: pre-wrap">{{ . }}</span>
</div>
`)
	if err != nil {
		panic(err)
	}
	w := bytes.Buffer{}
	err = t.Execute(&w, stateDrawing)
	if err != nil {
		return
	}

	reader := strings.NewReader(w.String())
	err = browser.OpenReader(reader)
	if err != nil {
		panic(err)
	}
	return
}
