package gongxueyun_service

import (
	"BronyaBot/global"
	"BronyaBot/internal/api"
	"BronyaBot/internal/service/gongxueyun_service/data"
	"BronyaBot/utils"
	"BronyaBot/utils/blockPuzzle"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type MoguDing struct {
	UserId        string     `json:"userId"`
	RoleKey       string     `json:"roleKey"`
	Authorization string     `json:"authorization"`
	PlanID        string     `json:"planId"`
	PlanName      string     `json:"planName"`
	PhoneNumber   string     `json:"phoneNumber"`
	Password      string     `json:"password"`
	Sign          SignStruct `json:"sign"`
	Email         string     `json:"email"`
}

func (m *MoguDing) Run() {
	global.Log.Infof("Starting sign-in process for user: %s", m.PhoneNumber)
	//m.GetBlock()
	if err := m.GetBlock(); err != nil {
		global.Log.Error(err.Error())
		return
	}

	if err := m.Login(); err != nil {
		global.Log.Error(err.Error())
		return
	}
	m.GetPlanId()
	m.SignIn()
}

type SignStruct struct {
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
}

var headers = map[string][]string{
	"User-Agent":   {"Mozilla/5.0 (Linux; U; Android 9; zh-cn; Redmi Note 5 Build/PKQ1.180904.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.141 Mobile Safari/537.36 XiaoMi/MiuiBrowser/11.10.8"},
	"Content-Type": {"application/json; charset=UTF-8"},
}
var clientUid = strings.ReplaceAll(uuid.New().String(), "-", "")
var comm = &commonParameters{}

func addHeader(key, value string) {
	// 检查 key 是否已经存在，若存在则追加到对应的值
	if _, exists := headers[key]; exists {
		//headers[key] = append(headers[key], value)
		headers[key] = []string{value}
	} else {
		// 若不存在，则新建一个字段
		headers[key] = []string{value}
	}
}
func GenerateRandomFloat(baseIntegerPart int) float64 {
	rand.Seed(time.Now().UnixNano())

	// Randomly adjust the integer part by ±1
	adjustment := rand.Intn(4) - 1 // Generates -1, 0, or 1
	integerPart := baseIntegerPart + adjustment

	// Calculate the maximum number of decimal places based on the integer part's length
	intPartLength := len(fmt.Sprintf("%d", integerPart))
	totalLength := rand.Intn(10) + 10 // Total length between 10 and 19
	decimalPlaces := totalLength - intPartLength

	if decimalPlaces <= 0 {
		decimalPlaces = 1 // Ensure at least one decimal place
	}

	// Generate a random decimal value with the specified number of decimal places
	decimalPart := rand.Float64() * math.Pow(10, float64(decimalPlaces))
	decimalPart = math.Trunc(decimalPart) / math.Pow(10, float64(decimalPlaces)) // Truncate to avoid floating-point imprecision

	return float64(integerPart) + decimalPart
}

func (mo *MoguDing) GetBlock() error {
	var maxRetries = 15
	for attempts := 1; attempts <= maxRetries; attempts++ {
		err := mo.processBlock()
		if err == nil {
			return nil
		}
		global.Log.Warning(fmt.Sprintf("Retrying captcha (%d/%d)", attempts, maxRetries))
		time.Sleep(6 * time.Second)
	}
	global.Log.Error("Failed to process captcha after maximum retries")
	return fmt.Errorf("failed to process captcha after maximum retries")
}

func (mo *MoguDing) processBlock() error {
	// 获取验证码数据
	requestData := map[string]string{
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
	comm.xY = string(marshal)
	comm.token = blockData.Data.Token
	comm.secretKey = blockData.Data.SecretKey
	cipher, _ := utils.NewAESECBPKCS5Padding(comm.secretKey, "base64")
	encrypt, _ := cipher.Encrypt(comm.xY)
	requestData = map[string]string{
		"pointJson":   encrypt,
		"token":       blockData.Data.Token,
		"captchaType": "blockPuzzle",
	}
	body, err = utils.SendRequest("POST", api.BaseApi+api.CHECK, requestData, headers)
	if err != nil {
		return fmt.Errorf("failed to verify captcha: %v", err)
	}

	// 解析验证结果
	var jsonContent data.CheckData
	if err := json.Unmarshal(body, &jsonContent); err != nil {
		return fmt.Errorf("failed to parse check response: %v", err)
	}
	if jsonContent.Code == 6111 {
		return fmt.Errorf("captcha verification failed, retry needed")
	}
	global.Log.Info("Captcha verification successful")
	padding, _ := utils.NewAESECBPKCS5Padding(blockData.Data.SecretKey, "base64")
	encrypt, err = padding.Encrypt(jsonContent.Data.Token + "---" + comm.xY)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to encrypt captcha: %v", err))
	}
	comm.captcha = encrypt
	return nil
}

func (mogu *MoguDing) Login() error {
	padding, _ := utils.NewAESECBPKCS5Padding(utils.MoGuKEY, "hex")
	encryptPhone, _ := padding.Encrypt(mogu.PhoneNumber)
	encryptPassword, _ := padding.Encrypt(mogu.Password)
	encryptTime, _ := padding.Encrypt(strconv.FormatInt(time.Now().UnixMilli(), 10))
	global.Log.Info("Login")
	requestData := map[string]string{
		"phone":     encryptPhone,
		"password":  encryptPassword,
		"captcha":   comm.captcha,
		"loginType": "android",
		"uuid":      clientUid,
		"device":    "android",
		"version":   "5.15.0",
		"t":         encryptTime,
	}
	var login = &data.Login{}
	var loginData = &data.LoginData{}
	body, err := utils.SendRequest("POST", api.BaseApi+api.LoginAPI, requestData, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(body, &login)
	if login.Code != 200 {
		global.Log.Error(login.Msg)
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
func (mogu *MoguDing) GetPlanId() {
	var planData = &data.PlanByStuData{}

	padding, _ := utils.NewAESECBPKCS5Padding(utils.MoGuKEY, "hex")
	encryptTime, _ := padding.Encrypt(strconv.FormatInt(time.Now().UnixMilli(), 10))

	sign := utils.CreateSign(mogu.UserId, mogu.RoleKey)

	addHeader("rolekey", mogu.RoleKey)
	addHeader("sign", sign)
	addHeader("authorization", mogu.Authorization)
	body := map[string]string{
		"pageSize": strconv.Itoa(999999),
		"t":        encryptTime,
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
	global.Log.Info("================")
	global.Log.Info(mogu.PlanID)
	global.Log.Info(mogu.PlanName)
	global.Log.Info("================")
	global.Log.Info(fmt.Sprintf("Get plan id successful"))
}
func (mogu *MoguDing) SignIn() {
	var resdata = &data.SaveData{}
	filling := dataStructureFilling(mogu)
	sign := utils.CreateSign(filling["device"], filling["type"], mogu.PlanID, mogu.UserId, filling["address"])
	addHeader("rolekey", mogu.UserId)
	addHeader("sign", sign)
	addHeader("authorization", mogu.Authorization)
	request, err := utils.SendRequest("POST", api.BaseApi+api.SignAPI, filling, headers)
	if err != nil {
		global.Log.Info(fmt.Sprintf("Failed to send request: %v", err))
	}
	json.Unmarshal(request, &resdata)
	global.Log.Info("================")
	global.Log.Info(resdata.Msg)
	utils.SendMail(mogu.Email, "测试邮件，请勿回复", resdata.Msg)
	global.Log.Info("================")
}

func dataStructureFilling(mogu *MoguDing) map[string]string {
	// 加载中国时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		global.Log.Error("Failed to load location: ", err)
		return nil
	}

	// 获取当前时间并格式化
	now := time.Now().In(loc)
	formattedTime := now.Format("2006-01-02 15:04:05")

	// 确定打卡类型
	typeStr := "START"
	if now.Hour() >= 12 {
		typeStr = "END"
	}

	// 加密当前时间戳
	padding, err := utils.NewAESECBPKCS5Padding(utils.MoGuKEY, "hex")
	if err != nil {
		global.Log.Error("Failed to initialize padding: ", err)
		return nil
	}

	encryptTime, err := padding.Encrypt(strconv.FormatInt(now.UnixMilli(), 10))
	if err != nil {
		global.Log.Error("Failed to encrypt timestamp: ", err)
		return nil
	}

	// 构造数据结构
	structuredData := data.SaveStructuredData{
		Address:    mogu.Sign.Address,
		City:       mogu.Sign.City,
		Area:       mogu.Sign.Area,
		Country:    mogu.Sign.Country,
		CreateTime: formattedTime,
		Device:     "{brand: Redmi Note 5, systemVersion: 14, Platform: Android}",
		//维度
		Latitude:  mogu.Sign.Latitude,
		Longitude: mogu.Sign.Longitude,
		Province:  mogu.Sign.Province,
		State:     "NORMAL",
		Type:      typeStr,
		UserId:    mogu.UserId,
		T:         encryptTime,
		PlanId:    mogu.PlanID,
	}

	// 转换为 JSON
	jsonData, err := json.Marshal(structuredData)
	if err != nil {
		global.Log.Error("Failed to marshal structuredData: ", err)
		return nil
	}

	// 转换为 map
	newMap := make(map[string]string)
	if err := json.Unmarshal(jsonData, &newMap); err != nil {
		global.Log.Error("Failed to unmarshal jsonData to map: ", err)
		return nil
	}

	return newMap
}
