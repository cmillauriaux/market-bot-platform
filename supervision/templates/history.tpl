{{define "history"}}
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
        <h2>Years statistics</h2>
        <table>
                {{range .YearsStatistics}}
                    <tr>
                        <td>{{.Date}} à {{.DateFin}}</td>
                        <td>{{.Min}}-{{.Max}} ({{.Delta}}%)</td>
                        <td>{{.Value}}</td>
                        <td>{{.Quantity}}</td>
                    </tr>
                {{end}}
        </table>
        <h2>Months statistics</h2>
        <table>
                {{range .MonthsStatistics}}
                    <tr>
                        <td>{{.Date}} à {{.DateFin}}</td>
                        <td>{{.Min}}-{{.Max}} ({{.Delta}}%)</td>
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