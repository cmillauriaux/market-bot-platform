{{define "bot"}}
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
                        <h4 class="page-title">Bots</h4>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-12">
                        <div class="white-box">
                            <div class="row">
                                <div class="col-md-12">
                                    <table class="table">
                                        <h2>Parameters</h2>
                                        <tr>
                                            <td>Bot</td>
                                            <td>{{.Display}}</td>
                                        </tr>
                                        <tr>
                                            <td>Wallet</td>
                                            <td>{{.Wallet}}</td>
                                        </tr>
                                        <tr>
                                            <td>Transactions value</td>
                                            <td>{{.GetTotalTransactionValue}}</td>
                                        </tr>
                                        <tr>
                                            <td>Orders value</td>
                                            <td>{{.GetTotalOrdersValue}}</td>
                                        </tr>
                                        <tr>
                                            <td>Total value</td>
                                            <td>{{.GetTotalValue}}</td>
                                        </tr>
                                    </table>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-12">
                                    <h2>History</h2>
                                    <table class="table">
                                        <tr>
                                            <th>OrderID</th>
                                            <th>Sell/Buy</th>
                                            <th>Date</th>
                                            <th>Original price</th>
                                            <th>Order price</th>
                                            <th>Tx price</th>
                                            <th>Quantity</th>
                                            <th>Plus value</th>
                                        </tr>
                                        {{range .GetHistory}}
                                            <tr>
                                                <td>{{.OrderID}}</td>
                                                <td>
                                                    {{if .Sell}}
                                                        SELL
                                                    {{end}}
                                                    {{if .Buy}}
                                                        BUY
                                                    {{end}}
                                                </td>
                                                <td>{{.DisplayDate}}</td>
                                                <td>{{.OriginalValue}}</td>
                                                <td>{{.OrderValue}}</td>
                                                <td>{{.TransactionValue}}</td>
                                                <td>{{.Quantity}}</td>
                                                <td>{{.GetPlusValue}}</td>
                                            </tr>
                                        {{end}}
                                    </table>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-12">
                                    <h2>Transactions</h2>
                                    <table class="table">
                                        <tr>
                                            <th>OrderID</th>
                                            <th>Sell/Buy</th>
                                            <th>Date</th>
                                            <th>Original price</th>
                                            <th>Transaction price</th>
                                            <th>Quantity</th>
                                            <th>In progress</th>
                                        </tr>
                                        {{range .Transactions}}
                                            <tr>
                                                <td>{{.OrderID}}</td>
                                                <td>
                                                    {{if .Sell}}
                                                        SELL
                                                    {{end}}
                                                    {{if .Buy}}
                                                        BUY
                                                    {{end}}
                                                </td>
                                                <td>{{.DisplayDate}}</td>
                                                <td>{{.TransactionValue}}</td>
                                                <td>{{.OrderValue}}</td>
                                                <td>{{.Quantity}}</td>
                                                <td>{{.InProgress}}</td>
                                            </tr>
                                        {{end}}
                                    </table>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-12">
                                    <h2>Orders</h2>
                                    <table class="table">
                                        <tr>
                                            <th>OrderID</th>
                                            <th>Sell/Buy</th>
                                            <th>Date</th>
                                            <th>Original price</th>
                                            <th>Order price</th>
                                            <th>Quantity</th>
                                            <th>In progress</th>
                                        </tr>
                                        {{range .Orders}}
                                            <tr>
                                                <td>{{.OrderID}}</td>
                                                <td>
                                                    {{if .Sell}}
                                                        SELL
                                                    {{end}}
                                                    {{if .Buy}}
                                                        BUY
                                                    {{end}}
                                                </td>
                                                <td>{{.DisplayDate}}</td>
                                                <td>{{.OriginalValue}}</td>
                                                <td>{{.OrderValue}}</td>
                                                <td>{{.Quantity}}</td>
                                                <td>{{.InProgress}}</td>
                                            </tr>
                                        {{end}}
                                    </table>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-12">
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
                                            "balloonText": "<div style='margin:5px; font-size:19px;'>Value:<b>[[value]]</b></div>",
                                            "bulletField": "bullet"
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
                                            {{range .Statistics.ForeverStatistics.Details}}
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