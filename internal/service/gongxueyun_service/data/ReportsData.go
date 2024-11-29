package data

type ReportsInfo struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
	Flag int           `json:"flag"`
}
