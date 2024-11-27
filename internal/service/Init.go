package service

import (
	"BronyaBot/global"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/cx_service"
	"BronyaBot/internal/service/gongxueyun_service"
	"context"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type AppService struct {
	users []entity.SignEntity
}

func NewAppService() *AppService {
	return &AppService{}
}

func (svc *AppService) Init() {
	// 加载用户数据
	svc.loadUsers()

	// 启动各模块服务
	svc.StartGongxueYun()
	//svc.StartTestCX()
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
	global.Log.Info("Starting Gongxueyun module...")
	concurrentLimit := 5                        // 最大并发数
	sem := make(chan struct{}, concurrentLimit) // 信号量
	var wg sync.WaitGroup

	// 创建限流器，每秒最多允许 1 个请求
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)

	for _, user := range svc.users {
		wg.Add(1)
		go func(user entity.SignEntity) {
			defer wg.Done()
			// 获取信号量，限制并发数
			sem <- struct{}{}
			// 限制请求频率
			if err := limiter.Wait(context.Background()); err != nil {
				global.Log.Error("Error waiting for rate limiter:", err)
				return
			}
			// 执行用户的操作
			ding := svc.createMoguDing(user)
			ding.Run()
			// 释放信号量
			<-sem
		}(user)
	}

	wg.Wait() // 等待所有 goroutine 完成
}

//func (svc *AppService) StartGongxueYun() {
//	global.Log.Info("Starting Gongxueyun module...")
//	for _, user := range svc.users {
//		ding := svc.createMoguDing(user)
//		ding.Run()
//	}
//}

func (svc *AppService) StartTestCX() {
	global.Log.Info("Starting CX test module...")
	logic := cx_service.CxLogic{
		Phone:    "1111111",
		Password: "1111111",
	}
	if err := logic.Login(); err != nil {
		global.Log.Error("CX login failed:", err)
		return
	}
	logic.PullCourse()
}

func (svc *AppService) createMoguDing(user entity.SignEntity) *gongxueyun_service.MoguDing {
	return &gongxueyun_service.MoguDing{
		PhoneNumber: user.Username,
		Password:    user.Password,
		Email:       user.Email,
		Sign: gongxueyun_service.SignStruct{
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
