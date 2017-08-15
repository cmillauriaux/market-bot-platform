package supervision

import (
	"github.com/cmillauriaux/market-bot-platform/bots"
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/ivpusic/neo"
)

func RunServer(history *history.History, bots []bots.Bot, basePath string) {
	app := neo.App()
	// Compile templates
	app.Templates(
		basePath + "/templates/*",
	)

	// Serve static files
	app.Serve("/bootstrap", basePath+"/templates/bootstrap")
	app.Serve("/css", basePath+"./templates/css")
	app.Serve("/js", basePath+"./templates/js")
	app.Serve("/less", basePath+"./templates/less")
	app.Serve("/plugins", basePath+"./templates/plugins")

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
	app.Get("/bots", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Tpl("bots", bots)
	})
	app.Get("/bot/:id", func(ctx *neo.Ctx) (int, error) {
		botID := ctx.Req.Params.Get("id")
		for _, bot := range bots {
			if bot.GetID() == botID {
				return 200, ctx.Res.Tpl("bot", bot)
			}
		}
		return 404, nil
	})
	app.Start()
}
