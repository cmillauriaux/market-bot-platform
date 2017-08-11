{{define "index"}}
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
                        <h4 class="page-title">Dashboard</h4>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-6">
                        <div class="row">
                            <div class="col-md-2">
                                <h2>Last hour</h2>
                                {{with .LastHourStatistics}}
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
                                        <td>Delta</td>
                                        <td>
                                            {{.Summary.Delta}}
                                            {{if .Summary.UpwardVariation}}
                                                <i class="fa fa-long-arrow-up green"></i>
                                            {{else}}
                                                <i class="fa fa-long-arrow-down red"></i>
                                            {{end}}
                                        </td>
                                    </tr>
                                    </table>
                                {{end}}
                            </div>
                            <div class="col-md-10">
                                <div id="chartYears" style="width: 100%; height: 500px;"></div>
                                <script>
                                    var chart = AmCharts.makeChart("chartYears", {
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
                                        "autoGridCount": false,
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
                                        {{range .LastHourStatistics.Details}}
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
                    <div class="col-md-6">
                        <div class="row">
                            <div class="col-md-2">
                                <h2>Last 6 hours</h2>
                                {{with .LastSixHoursStatistics}}
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
                                            <td>Delta</td>
                                            <td>
                                                {{.Summary.Delta}}
                                                {{if .Summary.UpwardVariation}}
                                                    <i class="fa fa-long-arrow-up green"></i>
                                                {{else}}
                                                    <i class="fa fa-long-arrow-down red"></i>
                                                {{end}}
                                            </td>
                                        </tr>
                                    </table>
                                {{end}}
                            </div>
                            <div class="col-md-10">
                                <div id="chartLastSixHours" style="width: 100%; height: 500px;"></div>
                                <script>
                                    var chart = AmCharts.makeChart("chartLastSixHours", {
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
                                        "autoGridCount": false,
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
                                        {{range .LastSixHoursStatistics.Details}}
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
                <div class="row">
                    <div class="col-md-6">
                        <div class="row">
                            <div class="col-md-2">
                                <h2>Last 24 hours</h2>
                                {{with .LastDayStatistics}}
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
                                        <td>Delta</td>
                                        <td>
                                            {{.Summary.Delta}}
                                            {{if .Summary.UpwardVariation}}
                                                <i class="fa fa-long-arrow-up green"></i>
                                            {{else}}
                                                <i class="fa fa-long-arrow-down red"></i>
                                            {{end}}
                                        </td>
                                    </tr>
                                    </table>
                                {{end}}
                            </div>
                            <div class="col-md-10">
                                <div id="chartDay" style="width: 100%; height: 500px;"></div>
                                <script>
                                    var chart = AmCharts.makeChart("chartDay", {
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
                                        "autoGridCount": false,
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
                                        {{range .LastDayStatistics.Details}}
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
                    <div class="col-md-6">
                        <div class="row">
                            <div class="col-md-2">
                                <h2>Last 6 days</h2>
                                {{with .Last7DaysStatistics}}
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
                                        <td>Delta</td>
                                        <td>
                                            {{.Summary.Delta}}
                                            {{if .Summary.UpwardVariation}}
                                                <i class="fa fa-long-arrow-up green"></i>
                                            {{else}}
                                                <i class="fa fa-long-arrow-down red"></i>
                                            {{end}}
                                        </td>
                                    </tr>
                                    </table>
                                {{end}}
                            </div>
                            <div class="col-md-10">
                                <div id="chart7Days" style="width: 100%; height: 500px;"></div>
                                <script>
                                    var chart = AmCharts.makeChart("chart7Days", {
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
                                        "autoGridCount": false,
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
                                        {{range .Last7DaysStatistics.Details}}
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
        {{template "footer"}}
    </body>
</html>
{{end}}