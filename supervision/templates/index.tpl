{{define "index"}}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>

    <body>
        {{template "header"}}
        <div>
        <h1>Preview</h1>
        <h2>Realtime statistics</h2>
        {{with .InstantStatistics}}
        <table>
                <tr>
                    <td>Realtime from</td>
                    <td>{{.Date}}</td>
                </tr>
                <tr>
                    <td>Realtime to</td>
                    <td>{{.DateFin}}</td>
                </tr>
                <tr>
                    <td>Min value</td>
                    <td>{{.Min}}</td>
                </tr>
                <tr>
                    <td>Max value</td>
                    <td>{{.Max}}</td>
                </tr>
                <tr>
                    <td>Avg value</td>
                    <td>{{.Value}}</td>
                </tr>
                <tr>
                    <td>Delta (min/max)</td>
                    <td>{{.Delta}}</td>
                </tr>
                <tr>
                    <td>Quantity</td>
                    <td>{{.Quantity}}</td>
                </tr>
        </table>
        {{end}}
        </div>
        <a href="/realtime">Realtime details</a>
        {{template "footer"}}
    </body>
</html>
{{end}}