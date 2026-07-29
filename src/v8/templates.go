package v8

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
<link href="https://fonts.googleapis.com/css?family=Poppins:300,400" rel="stylesheet">
<script src="https://cdn.jsdelivr.net/npm/mermaid@9.1.7/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<style>
:root { color: #666; font-family: "Poppins", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; font-size: 16px; font-weight: 300; }
* { box-sizing: border-box; }
body { margin: 0; overflow: hidden; padding: 18px 12px 8px; text-align: center; }
h1 { color: #555; font-size: 1.35rem; font-weight: 400; line-height: 1.3; margin: 0 0 14px; overflow-wrap: anywhere; }
.nav-buttons { align-items: center; display: flex; flex-wrap: wrap; gap: 8px; justify-content: center; }
button { background: #9cc8fa; border: 2px solid #9cc8fa; border-radius: 3px; color: #33485f; cursor: pointer; font: inherit; font-weight: 400; line-height: 1.2; padding: 7px 16px; }
button:hover, button:focus-visible { background: #248aff; border-color: #248aff; color: #fff; outline: none; }
button:focus-visible { box-shadow: 0 0 0 3px rgba(36, 138, 255, .25); }
button:disabled { background: #f5f7f9; border-color: #d8e9fc; color: #aaa; cursor: default; }
.hint { color: #888; flex-basis: 100%; font-size: .82rem; margin: 0; }
.graph { margin-top: 12px; }
.input { color: #555; font-size: clamp(2rem, 9vw, 3.6rem); line-height: 1.15; margin: 0 0 8px; overflow-wrap: anywhere; }
.input-processed { color: #d95050; }
.input-current { text-decoration-color: #d95050; text-decoration-line: underline; text-decoration-thickness: 3px; text-underline-offset: 4px; }
.mermaid { display: flex; justify-content: center; width: 100%; }
.mermaid svg { height: auto; max-width: 100% !important; }
@media (max-width: 480px) { body { padding-inline: 4px; } h1 { font-size: 1.15rem; } .hint { font-size: .75rem; } }
</style>
<body onload="prev()">

<h1>Regex: ({{ .Regex }})</h1>

<div class="nav-buttons">
	<button id="prev" onClick="prev()">
		Previous
	</button>
	<button id="next" onClick="next()">
		Next
	</button>
	<p class="hint">You can also use the arrow keys to step through the FSM.</p>
</div>

{{ range $i, $s := .Steps }}
<div class="graph" {{ if ne $i 0 }} style="display:none;visibility:hidden;" {{ else }} style="visibility:visible" {{ end }}>
	<p class="input">
		<span class="input-processed">{{ index .InputSplit 0 }}</span><span class="input-current">{{ index .InputSplit 1 }}</span><span>{{ index .InputSplit 2 }}</span>
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
	updateButtons(c.length)
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
	updateButtons(c.length)
}

function updateButtons(length) {
	document.getElementById('prev').disabled = i <= 0
	document.getElementById('next').disabled = i >= length - 1
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
