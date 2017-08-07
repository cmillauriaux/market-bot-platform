{{define "realtime"}}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>

    <body>
        {{template "header"}}
        <a href="/">Home</a>
        <div>
        <h1>Preview</h1>
        <h2>Realtime statistics</h2>
        <table>
            {{range .Realtime.Values}}
                <tr>
                    <td>{{.OrderID}}</td>
                    <td>{{.Date}}</td>
                    <td>{{.Value}}</td>
                    <td>{{.Quantity}}</td>
                </tr>
            {{end}}
        </table>
        </div>
        {{template "footer"}}
    </body>
</html>
{{end}}