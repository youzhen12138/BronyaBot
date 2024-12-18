package cx_service

import (
	"BronyaBot/global"
	"BronyaBot/internal/api"
	"BronyaBot/internal/service/cx_service/data"
	"BronyaBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (cx *CxLogic) loadAccount() {
	cx.Phone = global.Config.Account.Cx.Phone
	cx.Password = global.Config.Account.Cx.Password
}
func (cx *CxLogic) Run() {
	cx.loadAccount()
	global.Log.Info("Starting CX test module...")
	if err := cx.Login(); err != nil {
		global.Log.Error("CX login failed:", err)
		return
	}
	if err := cx.accInfo(); err != nil {
		global.Log.Error(err)
		return
	}

	if err := cx.PullCourse(); err != nil {
		global.Log.Error(err)
	}
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
func (c *CxLogic) accInfo() error {
	accinfoData := &data.AccinfoData{}
	resp, _ := http.NewRequest("GET", api.API_SSO_LOGIN, nil)
	resp.Header.Add("Cookie", c.cookie)
	resp.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	res, _ := c.client.Do(resp)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		global.Log.Error("CX login failed:", err)
	}
	err = json.Unmarshal(body, &accinfoData)
	if err != nil {
		global.Log.Error(err)
	}
	msg := accinfoData.Msg
	c.OssAccInfo.Puid = strconv.Itoa(msg.Puid)
	c.OssAccInfo.Phone = msg.Phone
	c.OssAccInfo.Uname = msg.Uname
	c.OssAccInfo.Name = msg.Name
	c.OssAccInfo.Sex = strconv.Itoa(msg.Sex)
	global.Log.Infof("账号登录成功: %s %s %s %s %s", c.OssAccInfo.Puid, c.OssAccInfo.Name, c.OssAccInfo.Sex, c.OssAccInfo.Phone, c.OssAccInfo.Uname)
	return nil
}
func (c *CxLogic) PullCourse() error {
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
	if pull.Result != 1 {
		return fmt.Errorf("拉取课程失败！")
	}
	global.Log.Infof("拉取成功！共%d个课程", len(pull.ChannelList))

	for i, s := range pull.ChannelList {
		global.Log.Infof("序号: %d 课程名: %s 老师名: %s 课程ID: %d 课程状态: %d - %s", i+1, s.Content.Course.Data[0].Name, s.Content.Course.Data[0].Teacherfactor, s.Content.Course.Data[0].Id, s.Content.Course.Data[0].Coursestate,
			func() string {
				switch s.Content.Course.Data[0].Coursestate {
				case 0:
					return "进行中"
				case 1:
					return "已结课"
				default:
					return "未知状态"
				}
			}())
	}
	global.Log.Info("==========\n")
	return nil
}
