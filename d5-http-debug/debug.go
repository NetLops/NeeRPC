package NeeRPC

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	debugText = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>NeeRPC Services</title>

</head>
<body>
{{range .}}
<hr>
Service {{.Name}}
<hr>
<table>
    <th align="center">Method</th><th align="center">Calls</th>
    {{range $name, $mtype := .Method}}
    <tr>
        <td align="left font=fixed">{{$name}}({{$mtype.ArgType}}, {{$mtype.ReplyType}}) error</td>
        <td align="center">{{$mtype.NumCalls}}</td>
    </tr>
    {{end}}
</table>
{{end}}
</body>
</html>`
)

var debug = template.Must(template.New("Rpc debug").Parse(debugText))

type debugHTTP struct {
	*Server
}

type debugService struct {
	Name   string
	Method map[string]*methodType
}

// Runs at /debug/neerpc
func (server debugHTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Build a sorted version of the data
	var services []debugService
	server.serviceMap.Range(func(key, value any) bool {
		svc := value.(*service)
		services = append(services, debugService{
			Name:   key.(string),
			Method: svc.method,
		})
		return true
	})
	err := debug.Execute(w, services)
	if err != nil {
		_, _ = fmt.Fprintln(w, "rpc: server executing template:", err.Error())
	}
}
