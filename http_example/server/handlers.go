package main

import (
	"html/template"
	"net/http"
)

// {{if .}} 到 {{end}} 的代码段仅在当前数据项（这里是点 .）的值非空时才会执行。 也就是说，当字符串为空时，此部分模板段会被忽略。
// 其中两段 {{.}} 表示要将数据显示在模板中 （即将查询字符串显示在 Web 页面上）。HTML 模板包将自动对文本进行转义， 因此文本的显示是安全的。
const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`

var templ = template.Must(template.New("qr").Parse(templateStr))

func HandlerQR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}
