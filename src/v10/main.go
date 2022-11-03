package v10

import (
	"bytes"
	"fmt"
	"github.com/pkg/browser"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type CmdFlag string

const reduceEpsilon CmdFlag = "reduce-epsilon"

// Main just used for linking up the main functions
func Main(args []string) {
	args, flags := parseArgumentsAndFlags(args)

	switch args[0] {
	case "draw":
		if len(args) == 2 {
			RenderFSM(args[1], flags)
		} else if len(args) == 3 {
			RenderRunner(args[1], args[2], flags)
		}
	case "out":
		if len(args) == 4 {
			OutputRunnerToFile(args[1], args[2], args[3], flags)
		}
	default:
		panic("command not recognized")
	}
}

// RenderFSM will render just the finite state machine, and output the result to the browser
func RenderFSM(input string, flags Set[CmdFlag]) {
	reducers := getReducersFromFlags(flags)
	graph := NewMyRegex(input, reducers...).DebugFSM()
	html := buildFsmHtml(graph)
	outputToBrowser(html)
}

// RenderRunner will render every step of the runner until it fails or succeeds. The template will then take care
// of hiding all but one of the steps to give the illusion of stepping through the input characters. It will
// then output the result to the browser.
func RenderRunner(regex, input string, flags Set[CmdFlag]) {
	data := buildRunnerTemplateData(regex, input, getReducersFromFlags(flags))
	htmlRunner := buildRunnerHTML(data)
	outputToBrowser(htmlRunner)
}

// OutputRunnerToFile will render every step of the runner, the same as RenderRunner, and then write the html to
// a file.
func OutputRunnerToFile(regex, input, filePath string, flags Set[CmdFlag]) {
	data := buildRunnerTemplateData(regex, input, getReducersFromFlags(flags))
	htmlRunner := buildRunnerHTML(data)
	outputToFile(htmlRunner, filePath)
}

func buildFsmHtml(graph string) string {
	return renderWithTemplate(fsmTemplate, graph)
}

func buildRunnerHTML(data TemplateData) string {
	return renderWithTemplate(runnerTemplate, data)
}

func buildRunnerTemplateData(regex string, input string, reducers []Reducer) TemplateData {
	newMyRegex := NewMyRegex(regex, reducers...)
	debugSteps := newMyRegex.DebugMatch(input)

	var steps []Step
	for _, step := range debugSteps {
		steps = append(steps, Step{
			Graph:      step.runnerDrawing,
			InputSplit: threeSplitString(input, step.currentCharacterIndex),
		})
	}

	data := TemplateData{
		Steps: steps,
		Regex: regex,
	}
	return data
}

func getReducersFromFlags(flags Set[CmdFlag]) []Reducer {
	var reducers []Reducer

	for flag := range flags {
		switch flag {
		case reduceEpsilon:
			reducers = append(reducers, &epsilonReducer{})
		}
	}
	return reducers
}

func outputToFile(html, path string) {
	containingDir := filepath.Dir(path)
	err := os.MkdirAll(containingDir, 0750)
	if err != nil {
		panic(err)
	}

	if filepath.Ext(path) == "" {
		path += ".html"
	}

	if filepath.Ext(path) != ".html" {
		panic("only .html extension permitted")
	}

	err = os.WriteFile(path, []byte(html), 0750)
	if err != nil {
		panic(err)
	}
}

func parseArgumentsAndFlags(args []string) ([]string, Set[CmdFlag]) {
	flagSet := NewSet[CmdFlag]()
	var arguments = []string{}
	for _, arg := range args {
		switch arg {
		case "--reduce-epsilons":
			flagSet.add(reduceEpsilon)
		default:
			if strings.HasPrefix(arg, "--") {
				panic(fmt.Sprintf("flag '%s' not recognized", arg))
			}
			arguments = append(arguments, arg)
		}
	}
	return arguments, flagSet
}

func renderWithTemplate(tmplt string, data any) string {
	t, err := template.New("graph").Parse(tmplt)
	if err != nil {
		panic(err)
	}
	w := bytes.Buffer{}
	err = t.Execute(&w, data)
	if err != nil {
		panic(err)
	}
	return w.String()
}

func outputToBrowser(html string) {
	reader := strings.NewReader(html)
	err := browser.OpenReader(reader)
	if err != nil {
		panic(err)
	}
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
