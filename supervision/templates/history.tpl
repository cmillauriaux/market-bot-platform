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
            <h1>Last statistics</h1>
            <h2>30 days statistics</h2>
            {{with .Last30DaysStatistics}}
                <table>
                    <tr>
                        <td>Open</td>
                        <td>{{.Summary.Open}}</td>
                    </tr>
                    <tr>
                        <td>Close</td>
                        <td>{{.Summary.Close}}</td>
                    </tr>
                    <tr>
                        <td>Min</td>
                        <td>{{.Summary.Min}}</td>
                    </tr>
                    <tr>
                        <td>Max</td>
                        <td>{{.Summary.Max}}</td>
                    </tr>
                    <tr>
                        <td>Average</td>
                        <td>{{.Summary.Value}}</td>
                    </tr>
                </table>
            {{end}}
            <div id="chart30Days" style="width: 100%; height: 700px;"></div>
            <script>
                var chart = AmCharts.makeChart("chart30Days", {
                    "type": "serial",
                    "theme": "light",
                    "marginRight": 40,
                    "marginLeft": 40,
                    "autoMarginOffset": 20,
                    "mouseWheelZoomEnabled":true,
                    "dataDateFormat": "YYYY-MM-DD",
                    "valueAxes": [{
                        "id": "v1",
                        "axisAlpha": 0,
                        "position": "left",
                        "ignoreAxisWidth":true
                    }],
                    "balloon": {
                        "borderThickness": 1,
                        "shadowAlpha": 0
                    },
                    "graphs": [{
                        "id": "g1",
                        "balloon":{
                        "drop":true,
                        "adjustBorderColor":false,
                        "color":"#ffffff"
                        },
                        "bullet": "round",
                        "bulletBorderAlpha": 1,
                        "bulletColor": "#FFFFFF",
                        "bulletSize": 5,
                        "hideBulletsCount": 50,
                        "lineThickness": 2,
                        "title": "red line",
                        "useLineColorForBulletBorder": true,
                        "valueField": "value",
                        "balloonText": "<span style='font-size:18px;'>[[value]]</span>"
                    }],
                    "chartScrollbar": {
                        "graph": "g1",
                        "oppositeAxis":false,
                        "offset":30,
                        "scrollbarHeight": 80,
                        "backgroundAlpha": 0,
                        "selectedBackgroundAlpha": 0.1,
                        "selectedBackgroundColor": "#888888",
                        "graphFillAlpha": 0,
                        "graphLineAlpha": 0.5,
                        "selectedGraphFillAlpha": 0,
                        "selectedGraphLineAlpha": 1,
                        "autoGridCount":true,
                        "color":"#AAAAAA"
                    },
                    "chartCursor": {
                        "pan": true,
                        "valueLineEnabled": true,
                        "valueLineBalloonEnabled": true,
                        "cursorAlpha":1,
                        "cursorColor":"#258cbb",
                        "limitToGraph":"g1",
                        "valueLineAlpha":0.2,
                        "valueZoomable":true
                    },
                    "valueScrollbar":{
                    "oppositeAxis":false,
                    "offset":50,
                    "scrollbarHeight":10
                    },
                    "categoryField": "date",
                    "categoryAxis": {
                        "parseDates": true,
                        "dashLength": 1,
                        "minorGridEnabled": true
                    },
                    "export": {
                        "enabled": true
                    },
                    "dataProvider": 
                    [
                        {{range .Last30DaysStatistics.Details}}
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