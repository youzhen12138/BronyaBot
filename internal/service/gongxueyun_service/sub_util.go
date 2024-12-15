package gongxueyun_service

import (
	"BronyaBot/global"
	"BronyaBot/internal/api"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/gongxueyun_service/data"
	"BronyaBot/utils"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 获取格式化的当前时间
func GetFormattedTime() (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		global.Log.Error("Failed to load location: ", err)
		return "", err
	}
	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05"), nil
}

// 获取打卡类型（START 或 END）
func GetClockType() (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		global.Log.Error("Failed to load location: ", err)
		return "", err
	}
	now := time.Now().In(loc)
	if now.Hour() >= 12 {
		return "END", nil
	}
	return "START", nil
}
func DataStructureFilling(mogu *MoguDing) map[string]interface{} {
	// 确定打卡类型
	typeStr, err := GetClockType()
	if err != nil {
		return nil
	}
	// 加密当前时间戳
	encryptTime, err := EncryptTimestamp(time.Now().UnixMilli())
	if err != nil {
		global.Log.Error("Failed to encrypt timestamp: ", err)
		return nil
	}
	formattedTime, err := GetFormattedTime()
	if err != nil {
		global.Log.Error("Failed to get formatted time: ", err)
	}
	// 直接构造 map，而不是先构造结构体再转换为 map
	return map[string]interface{}{
		"address":    mogu.Sign.Address,
		"city":       mogu.Sign.City,
		"area":       mogu.Sign.Area,
		"country":    mogu.Sign.Country,
		"createTime": formattedTime,
		"device":     "{brand: Redmi Note 5, systemVersion: 14, Platform: Android}",
		"latitude":   mogu.Sign.Latitude,
		"longitude":  mogu.Sign.Longitude,
		"province":   mogu.Sign.Province,
		"state":      "NORMAL",
		"type":       typeStr,
		"userId":     mogu.UserId,
		"t":          encryptTime,
		"planId":     mogu.PlanID,
	}
}
func SubmitStructureFilling(mogu *MoguDing, content string, title string, Retype string) map[string]interface{} {
	//formattedTime, err := GetFormattedTime()
	//if err != nil {
	//	global.Log.Error("Failed to get formatted time: ", err)
	//}
	timestamp, _ := EncryptTimestamp(time.Now().UnixMilli())
	submitData := data.SubmitData{
		Weeks:      mogu.WeekTime.Week,
		Content:    content,
		PlanId:     mogu.PlanID,
		ReportType: Retype,
		//ReportTime: formattedTime,
		Title:     title,
		JobId:     mogu.CommParameters.JobId,
		T:         timestamp,
		StartTime: mogu.WeekTime.StartTime,
		EndTime:   mogu.WeekTime.EndTime,
	}
	return data.SubmitDataFunc(submitData)
}

// 加密时间戳的通用方法
func EncryptTimestamp(timestamp int64) (string, error) {
	padding, err := utils.NewAESECBPKCS5Padding(utils.MoGuKEY, "hex")
	if err != nil {
		return "", fmt.Errorf("failed to initialize padding: %v", err)
	}

	encryptTime, err := padding.Encrypt(strconv.FormatInt(timestamp, 10))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt timestamp: %v", err)
	}
	return encryptTime, nil
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
func GenerateReportAI(userInput string, wordLimit int) string {
	// Define request data
	resdata := &data.AIData{}

	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + global.Config.AI.Token},
	}
	data := map[string]interface{}{
		"max_tokens":  4096,
		"top_k":       4,
		"temperature": 0.5,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("According to the information provided by the user, write an article strictly according to the template. Do not use Markdown syntax, HTML tags, or any special formatting. The output must be plain text and match the job description. The content should be fluent, conform to Chinese grammatical conventions, Repeat: DO NOT USE MARKDOWN ,and have more than %d characters. ", wordLimit),
			},
			{
				"role":    "user",
				"content": "模板：实习地点：xxxx\n\n工作内容：\n\nxzzzx\n\n工作总结：\n\nxxxxxx\n\n遇到问题：\n\nxzzzx\n\n自我评价：\n\nxxxxxx",
			},
			{
				"role":    "user",
				"content": "不能使用markdown格式输出注意！！！这是重要的",
			},
			{
				"role":    "user",
				"content": userInput,
			},
		},
		"model":  "4.0Ultra",
		"stream": false,
	}
	request, err := utils.SendRequest("POST", api.XUNFEIAPI, data, headers)
	if err != nil {
		global.Log.Error("Failed to send request: ", err)
	}
	json.Unmarshal(request, resdata)
	global.Log.Info("generate Successful!")
	global.Log.Infof("输入信息消耗token: %d\t大模型输出信息消耗token: %d\t总token: %d", resdata.Usage.PromptTokens, resdata.Usage.CompletionTokens, resdata.Usage.TotalTokens)
	return resdata.Choices[0].Message.Content
}
func LoadUsers() []entity.SignEntity {
	if global.Config.Account.Gongxueyun.Off {
		global.Log.Info("已开启yaml配置 启用本地加载单用户模式")
		return []entity.SignEntity{
			{
				ID:        -1,
				Username:  global.Config.Account.Gongxueyun.Phone,
				Password:  global.Config.Account.Gongxueyun.Password,
				Country:   global.Config.Account.Gongxueyun.Country,
				Province:  global.Config.Account.Gongxueyun.Province,
				City:      global.Config.Account.Gongxueyun.City,
				Area:      global.Config.Account.Gongxueyun.Area,
				Address:   global.Config.Account.Gongxueyun.Address,
				Latitude:  global.Config.Account.Gongxueyun.Latitude,
				Longitude: global.Config.Account.Gongxueyun.Longitude,
				Email:     global.Config.Account.Gongxueyun.Email,
			},
		}
	} else {
		var users []entity.SignEntity
		global.DB.Find(&users)
		if len(users) == 0 {
			global.Log.Warn("No users found in the database.")
		} else {
			global.Log.Info("Users loaded successfully.")
		}
		return users
	}
}
