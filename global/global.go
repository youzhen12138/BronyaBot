package global

import (
	"BronyaBot/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Log    *logrus.Logger
	Mail   *gomail.Dialer
)
