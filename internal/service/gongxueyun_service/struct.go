package gongxueyun_service

type MoguDing struct {
	ID             int              `json:"ID"`
	UserId         string           `json:"userId"`
	RoleKey        string           `json:"roleKey"`
	Authorization  string           `json:"authorization"`
	PlanID         string           `json:"planId"`
	PlanName       string           `json:"planName"`
	PhoneNumber    string           `json:"phoneNumber"`
	Password       string           `json:"password"`
	Sign           SignInfo         `json:"sign"`
	CommParameters commonParameters `json:"commParameters"`
	Email          string           `json:"email"`
}
type SignInfo struct {
	//	构造签到信息
	Address   string `json:"address"`
	City      string `json:"city"`
	Area      string `json:"area"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Province  string `json:"province"`
}
type commonParameters struct {
	token     string
	secretKey string
	xY        string
	captcha   string
	JobId     string
}
