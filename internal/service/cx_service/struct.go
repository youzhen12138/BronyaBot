package cx_service

import (
	"BronyaBot/internal/service/cx_service/data"
	"net/http"
)

type CxLogic struct {
	cookie     string
	Phone      string
	Password   string
	client     http.Client
	OssAccInfo oss_login
	ClassesLst data.Course
}
type oss_login struct {
	Puid       string `json:"puid"`
	Name       string `json:"name"`
	Sex        string `json:"sex"`
	Phone      string `json:"phone"`
	Uname      string `json:"uname"`
	SchoolName string `json:"schoolname"`
}
