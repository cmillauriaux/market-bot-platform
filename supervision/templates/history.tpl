{{define "history"}}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
        <script src="https://www.amcharts.com/lib/3/amcharts.js"></script>
        <script src="https://www.amcharts.com/lib/3/serial.js"></script>
        <script src="https://www.amcharts.com/lib/3/plugins/export/export.min.js"></script>
        <link rel="stylesheet" href="https://www.amcharts.com/lib/3/plugins/export/export.css" type="text/css" media="all" />
        <script src="https://www.amcharts.com/lib/3/themes/light.js"></script>
    </head>

    <body>
        {{template "header"}}
        <div>
            <h1>History</h1>
            
        </div>
        {{template "footer"}}
    </body>
</html>
{{end}}