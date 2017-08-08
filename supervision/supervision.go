package supervision

import (
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/ivpusic/neo"
)

func RunServer(history *history.History) {
	app := neo.App()
	app.Templates(
		"./supervision/templates/*",
	)
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
