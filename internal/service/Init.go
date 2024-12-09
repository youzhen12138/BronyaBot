package service

import (
	"BronyaBot/global"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/cx_service"
	"BronyaBot/internal/service/gongxueyun_service"
	"sync"
)

type AppService struct {
	users []entity.SignEntity
}

func NewAppService() *AppService {
	return &AppService{}
}

func (svc *AppService) Init() {
	// 启动各模块服务
	//ai := gongxueyun_service.GenerateReportAI("周报 工作地点: 浙江省宁波高新区光信路69号菁华路58号B座4001室 公司名:宁波晨希网络科技有限公司  岗位职责: 电话催收客服", 1500)
	//global.Log.Info("\n" + ai)
	svc.StartGongxueYun()
	//svc.StartTestCX()
	//utils.CreateSign("")
}

func (svc *AppService) loadUsers() {
	global.DB.Find(&svc.users)
	if len(svc.users) == 0 {
		global.Log.Warn("No users found in the database.")
	} else {
		global.Log.Info("Users loaded successfully.")
	}
}
func (svc *AppService) StartGongxueYun() {
	svc.loadUsers()
	global.Log.Info("Starting Gongxueyun module...")
	// 创建一个 Mutex 来保证每次只有一个 goroutine 执行相关操作
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, user := range svc.users {
		wg.Add(1)
		go func(user entity.SignEntity) {
			defer wg.Done()
			mu.Lock()
			ding := svc.createMoguDing(user)
			ding.Run()
			defer mu.Unlock()
		}(user)
	}
	wg.Wait() // 等待所有 goroutine 完成
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
