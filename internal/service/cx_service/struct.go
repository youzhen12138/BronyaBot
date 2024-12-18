package cx_service

import "net/http"

type CxLogic struct {
	cookie     string
	Phone      string
	Password   string
	client     http.Client
	OssAccInfo oss_login
}
type oss_login struct {
	Puid  string `json:"puid"`
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Phone string `json:"phone"`
	Uname string `json:"uname"`
}
