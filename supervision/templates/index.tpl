{{define "index"}}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>

    <body>
        {{template "header"}}
        <div>
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