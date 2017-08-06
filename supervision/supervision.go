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
	app.Start()
}
