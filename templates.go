package v5

type graphData struct {
	Graph string
	Input inputSearchString
}

type inputSearchString struct {
	Left, Middle, Right string
}

type templateData struct {
	Graphs []graphData
	Regex  string
}

const fsmTemplate = `
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
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
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<body onload="prev()">

<h1>Regex: ({{ .Regex }})</h1>

<div class="nav-buttons">
	<button id="prev" onClick="prev()">
		prev
	</button>
	<button id="next" onClick="next()">
		next
	</button>
	<p>Or use the arrow keys to step through the FSM</p>
</div>

{{ range .Graphs }}
<div class="graph">
	<p style='font-size:64px'>
		<span style='color:red'>{{ .Input.Left }}</span><span style='text-decoration-color:red;text-decoration-line:underline;'>{{ .Input.Middle }}</span><span>{{ .Input.Right }}</span>
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
		} else {
		  c[j].style.display = 'block'
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
		} else {
			c[j].style.display = 'block'
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
