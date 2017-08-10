package supervision

import (
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/ivpusic/neo"
)

func RunServer(history *history.History) {
	app := neo.App()
	// Compile templates
	app.Templates(
		"./supervision/templates/*",
	)

	// Serve static files
	app.Serve("/bootstrap", "./supervision/templates/bootstrap")
	app.Serve("/css", "./supervision/templates/css")
	app.Serve("/js", "./supervision/templates/js")
	app.Serve("/less", "./supervision/templates/less")
	app.Serve("/plugins", "./supervision/templates/plugins")

	// Serve pages
	app.Get("/", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Tpl("index", history)
	})
	app.Get("/history", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Tpl("history", history)
	})
	app.Get("/history-years", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Tpl("history-years", history)
	})
	app.Get("/realtime", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Tpl("realtime", history)
	})
	app.Start()
}
