package v8EpsilonConcatenation

type TemplateData struct {
	Steps []Step
	Regex string
	Input string
}

type Step struct {
	Graph      string
	InputSplit []string
}

const fsmTemplate = `
<script src="https://cdn.jsdelivr.net/npm/mermaid@9.1.7/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});
</script>
<div class="mermaid">
    {{ . }}
</div>
<div>
<span style="white-space: pre-wrap">{{ . }}</span>
</div>
`

const runnerTemplate = `
<script src="https://cdn.jsdelivr.net/npm/mermaid@9.1.7/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<body onload="prev()" style="text-align:center;margin:3vw;overflow: hidden">

<h1>Regex: ({{ .Regex }})</h1>

<div class="nav-buttons">
	<button id="prev" onClick="prev()">
		prev
	</button>
	<button id="next" onClick="next()">
		next
	</button>
	<p>Or use the arrow keys to Step through the FSM</p>
</div>

{{ range $i, $s := .Steps }}
<div class="graph" {{ if ne $i 0 }} style="visibility:hidden;" {{ else }} style="visibility:visible" {{ end }}>
	<p style='font-size:64px'>
		<span style='color:red'>{{ index .InputSplit 0 }}</span><span style='text-decoration-color:red;text-decoration-line:underline;'>{{ index .InputSplit 1 }}</span><span>{{ index .InputSplit 2 }}</span>
	</p>
	<div class="mermaid">
		{{ .Graph }}
	</div>
</div>
{{ end }}

<script type="text/javascript">
let i = 1

function next() {
  const c = document.getElementsByClassName('graph') 
  if (i >= c.length - 1) return 
	i++
	for (let j = 0; j < c.length; j++) {
		if (i != j)	{
		  c[j].style.display = 'none' 
		  c[j].style.visibility = 'hidden' 
		} else {
		  c[j].style.display = 'block'
		  c[j].style.visibility = 'visible' 
		}	
	}
}

function prev() {
	if (i <= 0) return
	i--
	const c = document.getElementsByClassName('graph') 
	for (let j = 0; j < c.length; j++) {
		if (i != j)	{
		  c[j].style.display = 'none' 
		  c[j].style.visibility = 'hidden' 
		} else {
		  c[j].style.display = 'block'
		  c[j].style.visibility = 'visible' 
		}	
	}
}

function checkKey(e) {
  	if (e.which === 37 || e.which === 40) {
		prev()
	} else if (e.which === 39 || e.which === 38) {
		next()
	}	
}

</script>
<script>document.onkeydown = checkKey;</script>
<div>
</div>
`
