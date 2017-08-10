{{define "history"}}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
        {{template "includes"}}
    </head>

    <body class="fix-header">
        <div id="wrapper">
            {{template "header"}}
            <div id="page-wrapper">
                <div class="row bg-title">
                    <div class="col-md-12">
                        <h4 class="page-title">Last statistics</h4>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-12">
                        <div class="white-box">
                            <div class="row">
                                <div class="col-md-2">
                                    <h2>30 days statistics</h2>
                                    {{with .Last30DaysStatistics}}
                                        <table class="table">
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
                                            <tr>
                                                <td>Variation</td>
                                                {{if .Summary.UpwardVariation}}
                                                    <td>Up</td>
                                                {{else}}
                                                    <td>Down</td>
                                                {{end}}
                                            </tr>
                                        </table>
                                    {{end}}
                                </div>
                                <div class="col-md-10">
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
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        {{template "footer"}}
    </body>
</html>
{{end}}