package api

var BaseApi = "https://api.moguding.net:9000"

const (
	BlockPuzzle            = "/session/captcha/v1/get"
	CHECK                  = "/session/captcha/v1/check"
	LoginAPI               = "/session/user/v6/login"
	GetPlanIDAPI           = "/practice/plan/v3/getPlanByStu"
	SignAPI                = "/attendence/clock/v4/save"
	GetWeekCountAPI        = "/practice/paper/v2/listByStu"
	GetWeeks               = "/practice/paper/v1/getWeeks1"
	SubmitWeeklyOrDailyAPI = "/practice/paper/v2/save"
)
