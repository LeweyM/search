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
		if len(args) == 2 {
			RenderFSM(args[1])
		} else if len(args) == 3 {
			RenderRunner(args[1], args[2])
		}
	default:
		fmt.Println("command not recognized")
	}
}

func RenderFSM(input string) {
	tokens := lex(input)
	parser := NewParser()
	ast := parser.Parse(tokens)
	head, _ := ast.compile()

	graph, _ := head.Draw()

	renderTemplateToBrowser(fsmTemplate, graph)
}

// RenderRunner will render every step of the runner until it fails or succeeds. The template will then take care
// of hiding all but one of the steps to give the illusion of stepping through the input characters.
func RenderRunner(regex, input string) {
	tokens := lex(regex)
	parser := NewParser()
	ast := parser.Parse(tokens)
	head, _ := ast.compile()

	allGraphSteps := NewRunner(head).drawAllGraphSteps(input)

	var graphs []graphData
	for letterIndex, graphStep := range allGraphSteps {
		graphs = append(graphs, graphData{
			Graph: graphStep,
			Input: splitSearchString(input, letterIndex),
		})
	}

	renderTemplateToBrowser(runnerTemplate, templateData{
		Graphs: graphs,
		Regex:  regex,
	})
}

// splitSearchString divides the search input string into three pieces so that we can render in the browser:
// 1. What has already been processed
// 2. The next character to process
// 3. What is yet to be processed
func splitSearchString(input string, currentLetterIndex int) inputSearchString {
	var left, middle, right string

	left = input[:currentLetterIndex]
	if currentLetterIndex < len(input) {
		middle = string(input[currentLetterIndex])
		right = input[currentLetterIndex+1:]
	}

	return inputSearchString{
		Left:   left,
		Middle: middle,
		Right:  right,
	}
}

func renderTemplateToBrowser(tmplt string, data any) {
	t, err := template.New("graph").Parse(tmplt)
	if err != nil {
		panic(err)
	}
	w := bytes.Buffer{}
	err = t.Execute(&w, data)
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(w.String())
	err = browser.OpenReader(reader)
	if err != nil {
		panic(err)
	}
	return
}
