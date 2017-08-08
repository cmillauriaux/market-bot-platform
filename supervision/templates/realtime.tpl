{{define "realtime"}}
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
            <h1>Preview</h1>
            <h2>Years statistics</h2>
            <div id="chartMonths" style="width: 100%; height: 500px;"></div>
            <script>
                var chart = AmCharts.makeChart("chartMonths", {
                    "type": "serial",
                    "theme": "light",
                    "marginRight": 80,
                    "valueAxes": [{
                        "position": "left",
                        "title": "Unique visitors"
                    }],
                    "graphs": [{
                        "id": "g1",
                        "fillAlphas": 0.4,
                        "valueField": "value",
                        "balloonText": "<div style='margin:5px; font-size:19px;'>Value:<b>[[value]]</b></div>"
                    }],
                    "chartScrollbar": {
                        "graph": "g1",
                        "scrollbarHeight": 80,
                        "backgroundAlpha": 0,
                        "selectedBackgroundAlpha": 0.1,
                        "selectedBackgroundColor": "#888888",
                        "graphFillAlpha": 0,
                        "graphLineAlpha": 0.5,
                        "selectedGraphFillAlpha": 0,
                        "selectedGraphLineAlpha": 1,
                        "autoGridCount": true,
                        "color": "#AAAAAA"
                    },
                    "chartCursor": {
                        "categoryBalloonDateFormat": "JJ:NN, DD MMMM",
                        "cursorPosition": "mouse"
                    },
                    "categoryField": "date",
                    "categoryAxis": {
                        "minPeriod": "mm",
                        "parseDates": true
                    },
                    "export": {
                        "enabled": true,
                        "dateFormat": "YYYY-MM-DD HH:NN:SS"
                    },
                    "dataProvider": 
                    [
                        {{range .Realtime.Values}}
                        {
                            "date": "{{.DisplayDate}}",
                            "value": {{.Value}} / 100
                        }, 
                        {{end}}
                    ]
                });

                chart.addListener("rendered", zoomChart);

                zoomChart();

                function zoomChart() {
                    chart.zoomToIndexes(chart.dataProvider.length - 40, chart.dataProvider.length - 1);
                }
            </script>
        </div>
        {{template "footer"}}
    </body>
</html>
{{end}}