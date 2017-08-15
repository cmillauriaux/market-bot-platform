{{define "bots"}}
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
                            {{range .}}
                                <div class="row">
                                    <div class="col-md-12">
                                        <h2>Parameters</h2>
                                        <table class="table">
                                            <tr>
                                                <td><a href="/bot/{{.GetID}}">Details</a></td>
                                            </tr>
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
                            {{end}}
                        </div>
                    </div>
                </div>
            </div>
        </div>
        {{template "footer"}}
    </body>
</html>
{{end}}