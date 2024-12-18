package gongxueyun_service

import (
	"BronyaBot/global"
	"BronyaBot/internal/api"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/gongxueyun_service/data"
	"BronyaBot/utils"
	"BronyaBot/utils/blockPuzzle"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func (m *MoguDing) Run(runType string) {
	if err := m.GetBlock(); err != nil {
		utils.SendMail(m.Email, "Block-Error", err.Error())
		global.Log.Error(err.Error())
		return
	}
	if err := m.Login(); err != nil {
		utils.SendMail(m.Email, "Login-Error-测试邮件请勿回复", err.Error())
		global.Log.Error(err.Error())
		return
	}

	if err := m.GetPlanId(); err != nil {
		global.Log.Error(err.Error())
		return
	}
	if err := m.GetJobInfo(); err != nil {
		global.Log.Error("Failed to get job info: %v", err)
		return
	}
	m.getWeeksTime()
	if runType == "sign" {
		m.SignIn()
	} else if runType == "week" {
		m.getSubmittedReportsInfo("week")
		m.SubmitReport("week", 1500)
	} else if runType == "month" {
		m.getSubmittedReportsInfo("month")
		m.SubmitReport("month", 1600)
	}

}

var headers = map[string][]string{
	"User-Agent":   {"Mozilla/5.0 (Linux; U; Android 9; zh-cn; Redmi Note 5 Build/PKQ1.180904.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.141 Mobile Safari/537.36 XiaoMi/MiuiBrowser/11.10.8"},
	"Content-Type": {"application/json; charset=UTF-8"},
	"host":         {"api.moguding.net:9000"},
}
var clientUid = strings.ReplaceAll(uuid.New().String(), "-", "")

func addHeader(key, value string) {
	if _, exists := headers[key]; exists {
		headers[key] = []string{value}
	} else {
		headers[key] = []string{value}
	}
}
func addStandardHeaders(roleKey, userId, authorization string) {
	addHeader("rolekey", roleKey)
	addHeader("userid", userId)
	addHeader("authorization", authorization)

}

func (mo *MoguDing) GetBlock() error {
	global.Log.Infof("Starting sign-in process for user: %s", mo.Email)
	var maxRetries = 15
	for attempts := 1; attempts <= maxRetries; attempts++ {
		err := mo.processBlock()
		if err == nil {
			return nil
		}
		global.Log.Warning(fmt.Sprintf("Retrying captcha (%d/%d)", attempts, maxRetries))
		time.Sleep(10 * time.Second)
	}
	global.Log.Error("Failed to process captcha after maximum retries")
	return fmt.Errorf("failed to process captcha after maximum retries")
}
func (mo *MoguDing) processBlock() error {
	// 获取验证码数据
	requestData := map[string]interface{}{
		"clientUid":   clientUid,
		"captchaType": "blockPuzzle",
	}
	body, err := utils.SendRequest("POST", api.BaseApi+api.BlockPuzzle, requestData, headers)
	if err != nil {
		return fmt.Errorf("failed to fetch block puzzle: %v", err)
	}
	// 解析响应数据
	blockData := &data.BlockRes{}
	if err := json.Unmarshal(body, &blockData); err != nil {
		return fmt.Errorf("failed to parse block puzzle response: %v", err)
	}
	// 初始化滑块验证码
	captcha, err := blockPuzzle.NewSliderCaptcha(blockData.Data.JigsawImageBase64, blockData.Data.OriginalImageBase64)
	if err != nil {
		return fmt.Errorf("failed to initialize captcha: %v", err)
	}
	x, _ := captcha.FindBestMatch()

	// 加密并验证
	xY := map[string]string{"x": strconv.FormatFloat(GenerateRandomFloat(x), 'f', -1, 64), "y": strconv.Itoa(5)}
	global.Log.Info(fmt.Sprintf("Captcha matched at: xY=%s", xY))

	marshal, err := json.Marshal(xY)

	mo.CommParameters.xY = string(marshal)
	mo.CommParameters.token = blockData.Data.Token
	mo.CommParameters.secretKey = blockData.Data.SecretKey
	cipher, _ := utils.NewAESECBPKCS5Padding(mo.CommParameters.secretKey, "base64")
	encrypt, _ := cipher.Encrypt(mo.CommParameters.xY)
	requestData = map[string]interface{}{
		"pointJson":   encrypt,
		"token":       blockData.Data.Token,
		"captchaType": "blockPuzzle",
	}
	body, err = utils.SendRequest("POST", api.BaseApi+api.CHECK, requestData, headers)
	if err != nil {
		return fmt.Errorf("failed to verify captcha: %v", err)
	}

	// 解析验证结果
	jsonContent := &data.CheckData{}
	if err := json.Unmarshal(body, &jsonContent); err != nil {
		return fmt.Errorf("failed to parse check response: %v", err)
	}
	if jsonContent.Code == 6111 {
		return fmt.Errorf("captcha verification failed, retry needed")
	}
	global.Log.Info("Captcha verification successful")
	padding, _ := utils.NewAESECBPKCS5Padding(blockData.Data.SecretKey, "base64")
	encrypt, err = padding.Encrypt(jsonContent.Data.Token + "---" + mo.CommParameters.xY)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to encrypt captcha: %v", err))
	}
	mo.CommParameters.captcha = encrypt
	return nil
}
func (mogu *MoguDing) Login() error {
	padding, _ := utils.NewAESECBPKCS5Padding(utils.MoGuKEY, "hex")
	encryptPhone, _ := padding.Encrypt(mogu.PhoneNumber)
	encryptPassword, _ := padding.Encrypt(mogu.Password)
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	requestData := map[string]interface{}{
		"phone":     encryptPhone,
		"password":  encryptPassword,
		"captcha":   mogu.CommParameters.captcha,
		"loginType": "android",
		"uuid":      clientUid,
		"device":    "android",
		"version":   "5.15.0",
		"t":         timestamp,
	}
	var login = &data.Login{}
	var loginData = &data.LoginData{}
	body, err := utils.SendRequest("POST", api.BaseApi+api.LoginAPI, requestData, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(body, &login)
	if login.Code != 200 {
		return fmt.Errorf(login.Msg)

	}
	decrypt, err := padding.Decrypt(login.Data)
	json.Unmarshal([]byte(decrypt), &loginData)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to decrypt data: %v", err))
	}
	mogu.RoleKey = loginData.RoleKey
	mogu.UserId = loginData.UserId
	mogu.Authorization = loginData.Token
	global.Log.Info("================")
	global.Log.Info(loginData.NikeName)
	global.Log.Info(loginData.Phone)
	global.Log.Info("================")
	global.Log.Info("Login successful")
	return nil
}
func (mogu *MoguDing) GetPlanId() error {
	defaultId := "6686304d065db846edab7d4565065abc"
	planData := &data.PlanByStuData{}
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	sign := utils.CreateSign(mogu.UserId, mogu.RoleKey)
	addHeader("sign", sign)
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	body := map[string]interface{}{
		"pageSize": strconv.Itoa(999999),
		"t":        timestamp,
	}
	request, err := utils.SendRequest("POST", api.BaseApi+api.GetPlanIDAPI, body, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(request, &planData)
	for i := range planData.Data {
		mogu.PlanID = planData.Data[i].PlanId
		mogu.PlanName = planData.Data[i].PlanName
	}
	if strings.EqualFold(mogu.PlanID, defaultId) {
		return fmt.Errorf(mogu.PlanName)
	}
	global.Log.Info("================")
	global.Log.Info(mogu.PlanID)
	global.Log.Info(mogu.PlanName)
	global.Log.Info("================")
	return nil
}
func (mogu *MoguDing) GetJobInfo() error {
	job := &data.JobInfoData{}
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	body := map[string]interface{}{
		"planId": mogu.PlanID,
		"t":      timestamp,
	}
	request, err := utils.SendRequest("POST", api.BaseApi+api.GetJobInfoAPI, body, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(request, &job)
	if job.Data.JobId == "" {
		return fmt.Errorf("job info not found")
	} else {
		mogu.JobInfo.JobName = job.Data.JobName
		mogu.JobInfo.Address = job.Data.Address
		mogu.JobInfo.CompanyName = job.Data.CompanyName
	}
	return nil
}
func (mogu *MoguDing) SignIn() {
	resdata := &data.SaveData{}
	filling := DataStructureFilling(mogu)
	sign := utils.CreateSign(filling["device"].(string), filling["type"].(string), mogu.PlanID, mogu.UserId, filling["address"].(string))
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	addHeader("sign", sign)
	request, err := utils.SendRequest("POST", api.BaseApi+api.SignAPI, filling, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}

	json.Unmarshal(request, &resdata)
	global.Log.Info("================")
	global.Log.Info(resdata.Msg)
	global.Log.Info("================")
	if resdata.Msg == "success" {
		mogu.updateSignState(1)
	} else {
		mogu.updateSignState(0)
	}
	utils.SendMail(mogu.Email, "检查是否打卡完成", resdata.Msg+"\n如果未成功请联系管理员")

}
func (mogu *MoguDing) updateSignState(state int) {
	// 更新数据库表中的 state 字段
	if mogu.ID != -1 {
		err := global.DB.Model(&entity.SignEntity{}).Where("username = ?", mogu.PhoneNumber).Update("state", state).Error
		if err != nil {
			global.Log.Error(fmt.Sprintf("Failed to update state for user %s: %v", mogu.PhoneNumber, err))
		} else {
			global.Log.Info(fmt.Sprintf("Successfully updated state for user %s to %d", mogu.PhoneNumber, state))
		}
	}
}

// 获取已经提交的日报、周报或月报的数量。
func (mogu *MoguDing) getSubmittedReportsInfo(reportType string) {
	report := &data.ReportsInfo{}
	sign := utils.CreateSign(mogu.UserId, mogu.RoleKey, reportType)
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	addHeader("sign", sign)
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	body := map[string]interface{}{
		"currPage":   1,
		"pageSize":   10,
		"reportType": reportType,
		"planId":     mogu.PlanID,
		"t":          timestamp,
	}
	request, err := utils.SendRequest("POST", api.BaseApi+api.GetWeekCountAPI, body, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(request, &report)
	if report.Flag == 0 {
		global.Log.Warning("未发现之前存在报告，初始化报告为0")
		mogu.ReportStruct.CreateTime = ""
		mogu.ReportStruct.ReportId = ""
		mogu.ReportStruct.ReportType = ""
		mogu.ReportStruct.Flag = 0
		return
	} else {
		mogu.ReportStruct.CreateTime = report.Data[0].CreateTime
		mogu.ReportStruct.ReportId = report.Data[0].ReportId
		mogu.ReportStruct.ReportType = report.Data[0].ReportType
		mogu.ReportStruct.Flag = report.Flag
	}
}

// 获取提交周时间
func (mogu *MoguDing) getWeeksTime() {
	week := &data.WeeksData{}
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	body := map[string]interface{}{
		"t": timestamp,
	}
	request, err := utils.SendRequest("POST", api.BaseApi+api.GetWeeks, body, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(request, &week)
	if len(week.Data) > 0 {
		mogu.WeekTime.Week = week.Data[0].Weeks
		mogu.WeekTime.StartTime = week.Data[0].StartTime
		mogu.WeekTime.EndTime = week.Data[0].EndTime
		mogu.WeekTime.IsDefault = week.Data[0].IsDefault
		mogu.WeekTime.Flag = week.Flag
	}
}

// SubmitReport
// 提交定时报告
func (mogu *MoguDing) SubmitReport(reportType string, limit int) {
	res := &data.RepResData{}
	var _t string
	switch reportType {
	case "week":
		_t = "周报"
	case "month":
		_t = "月报"
	case "day":
		_t = "日报"
	}
	input := fmt.Sprintf("报告类型: %s 工作地点: %s 公司名: %s 岗位职责: %s", _t, mogu.JobInfo.Address, mogu.JobInfo.CompanyName, mogu.JobInfo.JobName)
	ai := GenerateReportAI(input, limit)
	addStandardHeaders(mogu.RoleKey, mogu.UserId, mogu.Authorization)
	filling := SubmitStructureFilling(mogu, ai, "报告", reportType)
	sign := utils.CreateSign(mogu.UserId, reportType, mogu.PlanID, "报告")
	addHeader("sign", sign)
	request, _ := utils.SendRequest("POST", api.BaseApi+api.SubmitAReport, filling, headers)
	json.Unmarshal(request, &res)
	global.Log.Info(fmt.Sprintf("Submit report: %v", res))
	utils.SendMail(mogu.Email, strconv.Itoa(res.Code), res.Msg+"\n如果未成功请联系管理员")
}
