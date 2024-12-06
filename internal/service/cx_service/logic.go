package cx_service

import (
	"BronyaBot/global"
	"BronyaBot/internal/api"
	"BronyaBot/internal/service/cx_service/data"
	"BronyaBot/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CxLogic struct {
	cookie   string
	Phone    string
	Password string
}

func (cx *CxLogic) Run() {
	global.Log.Info("Starting CX test module...")
	if err := cx.Login(); err != nil {
		global.Log.Error("CX login failed:", err)
		return
	}
	cx.PullCourse()
}

func (c *CxLogic) Login() error {

	global.Log.Info("=====开始登录=====")
	phone, _ := utils.AESCBCEncrypt([]byte(c.Phone))
	password, _ := utils.AESCBCEncrypt([]byte(c.Password))
	postres, err := http.PostForm(api.API_LOGIN_WEB, url.Values{
		"fid":               {"-1"},
		"uname":             {phone},
		"password":          {password},
		"t":                 {"true"},
		"forbidotherlogin":  {"0"},
		"validate":          {""},
		"doubleFactorLogin": {"0"},
		"independentId":     {"0"},
		"independentNameId": {"0"},
	})
	if err != nil {
		global.Log.Error(err.Error())
	}
	defer postres.Body.Close()
	var jsonContent map[string]interface{}
	body, _ := io.ReadAll(postres.Body)
	json.Unmarshal(body, &jsonContent)
	if jsonContent["status"] == false {
		global.Log.Errorf("登录失败 %t", jsonContent["status"])
		return err
	}
	global.Log.Infof("%s login, status: %t ", c.Phone, jsonContent["status"])
	values := postres.Header.Values("Set-Cookie")
	for _, v := range values {
		c.cookie += strings.ReplaceAll(strings.ReplaceAll(v, "HttpOnly", ""), "Path=/", "")
	}
	return nil
}
func (c *CxLogic) PullCourse() {
	global.Log.Info("=====拉取课程=====")
	client := &http.Client{}
	request, err := http.NewRequest("GET", api.API_CLASS_LST, nil)
	if err != nil {
		global.Log.Error(err.Error())
	}
	request.Header.Add("Cookie", c.cookie)
	request.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	res, err := client.Do(request)
	if err != nil {
		global.Log.Error(err.Error())
	}
	pull := data.Course{}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		global.Log.Error(err.Error())
	}
	json.Unmarshal(body, &pull)
	for i := range pull.ChannelList {
		global.Log.Info(pull.ChannelList[i].Content)
	}

}
