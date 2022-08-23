package v5

import (
	"bytes"
	"fmt"
	"github.com/pkg/browser"
	"html/template"
	"os"
	"strings"
)

// Main just used for linking up the main functions
func Main() {
	main()
}

func main() {
	switch os.Args[2] {
	case "draw":
		Draw()
	default:
		fmt.Println("command not recognized")
	}
}

func Draw() {
	tokens := lex(os.Args[3])
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
