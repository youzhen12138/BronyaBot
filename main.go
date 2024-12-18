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
	global.Mail = core.InitMail()
	appService := service.NewAppService()
	appService.Init()

}
