package v10

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
	graph := NewMyRegex(input).DebugFSM()
	renderTemplateToBrowser(fsmTemplate, graph)
}

// RenderRunner will render every step of the runner until it fails or succeeds. The template will then take care
// of hiding all but one of the steps to give the illusion of stepping through the input characters.
func RenderRunner(regex, input string) {
	newMyRegex := NewMyRegex(regex)
	debugSteps := newMyRegex.DebugMatch(input)

	var steps []Step
	for _, step := range debugSteps {
		steps = append(steps, Step{
			Graph:      step.runnerDrawing,
			InputSplit: threeSplitString(input, step.currentCharacterIndex),
		})
	}

	renderTemplateToBrowser(runnerTemplate, TemplateData{
		Steps: steps,
		Regex: regex,
	})
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

// threeSplitString divides a string into three pieces on a given index
func threeSplitString(s string, i int) []string {
	var left, middle, right string

	left = s[:i]
	if i < len(s) {
		middle = string(s[i])
		right = s[i+1:]
	}

	return []string{left, middle, right}
}
