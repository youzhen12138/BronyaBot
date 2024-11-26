package main

import (
	"BronyaBot/core"
	"BronyaBot/global"
	"BronyaBot/internal/service"
)

func main() {
	core.InitConf()
	global.Log = core.InitLogger()
	global.DB = core.InitGorm()
	appService := service.NewAppService()
	appService.Init()
}
