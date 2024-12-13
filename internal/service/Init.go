package service

import (
	"BronyaBot/global"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/cx_service"
	"BronyaBot/internal/service/gongxueyun_service"
	"github.com/robfig/cron/v3"
	"time"
)

type AppService struct {
	users []entity.SignEntity
	cron  *cron.Cron
}

func NewAppService() *AppService {
	return &AppService{
		cron: cron.New(),
	}
}

func (svc *AppService) Init() {
	// 启动各模块服务
	svc.scheduleTasks()
	svc.cron.Start()
	//svc.StartTestCX()
}

func (svc *AppService) loadUsers() {
	if global.Config.Account.Gongxueyun.Off {
		global.Log.Info("已开启yaml配置  启用本地加载单用户模式")
		svc.users = append(svc.users, entity.SignEntity{
			ID:        -1,
			Username:  global.Config.Account.Gongxueyun.Phone,
			Password:  global.Config.Account.Gongxueyun.Password,
			Country:   global.Config.Account.Gongxueyun.Country,
			Province:  global.Config.Account.Gongxueyun.Province,
			City:      global.Config.Account.Gongxueyun.City,
			Area:      global.Config.Account.Gongxueyun.Area,
			Address:   global.Config.Account.Gongxueyun.Address,
			Latitude:  global.Config.Account.Gongxueyun.Latitude,
			Longitude: global.Config.Account.Gongxueyun.Longitude,
			Email:     global.Config.Account.Gongxueyun.Email,
		})
	} else {
		global.DB.Find(&svc.users)
		if len(svc.users) == 0 {
			global.Log.Warn("No users found in the database.")
		} else {
			global.Log.Info("Users loaded successfully.")
		}
	}
}
func (svc *AppService) scheduleTasks() {
	global.Log.Info("Scheduling tasks...")
	// 每天早上8点执行
	svc.cron.AddFunc("0 8 * * *", func() {
		global.Log.Info("Running task: 每天早上8点签到")
		svc.StartGongxueYun("sign")
		global.Log.Info("Task finished!")
	})
	// 每天晚上6点执行
	svc.cron.AddFunc("0 18 * * *", func() {
		global.Log.Info("Running task: 每天晚上6点签到")
		svc.StartGongxueYun("sign")
		global.Log.Info("Task finished!")
	})

	// 每周周五早上10点执行
	svc.cron.AddFunc("0 10 * * 5", func() {
		global.Log.Info("Running task: 每周周五早上10点签到")
		svc.StartGongxueYun("week")
		global.Log.Info("Task finished!")
	})

	// 每月最后一周的周一早上10点执行
	svc.cron.AddFunc("0 10 ? * 1L", func() {
		if isLastWeek(time.Now()) {
			global.Log.Info("Running task: 每月最后一周的周一早上10点签到")
			svc.StartGongxueYun("month")
			global.Log.Info("Task finished!")
		}
	})
}

func (svc *AppService) StartGongxueYun(_type string) {
	svc.loadUsers()
	global.Log.Info("Starting Gongxueyun module...")
	for i := range svc.users {
		ding := svc.createMoguDing(svc.users[i])
		ding.Run(_type)
	}
}

func (svc *AppService) StartTestCX() {
	cxLogic := cx_service.CxLogic{}
	cxLogic.Run()
}

func (svc *AppService) createMoguDing(user entity.SignEntity) *gongxueyun_service.MoguDing {
	return &gongxueyun_service.MoguDing{
		ID:          user.ID,
		PhoneNumber: user.Username,
		Password:    user.Password,
		Email:       user.Email,
		Sign: gongxueyun_service.SignInfo{
			City:      user.City,
			Area:      user.Area,
			Address:   user.Address,
			Country:   user.Country,
			Province:  user.Province,
			Latitude:  user.Latitude,
			Longitude: user.Longitude,
		},
	}
}
func isLastWeek(t time.Time) bool {
	_, week := t.ISOWeek()
	nextMonday := t.AddDate(0, 0, 7-int(t.Weekday())) // 下一个星期一
	nextMonthWeek, _ := nextMonday.ISOWeek()
	return week != nextMonthWeek
}
