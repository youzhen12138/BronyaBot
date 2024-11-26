package blockPuzzle

import (
	"BronyaBot/global"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// SliderCaptcha 用于封装滑块验证码的处理逻辑
type SliderCaptcha struct {
	Slider     image.Image // 滑块图片
	Background image.Image // 背景图片
}

// NewSliderCaptcha 初始化 SliderCaptcha
func NewSliderCaptcha(sliderBase64, backgroundBase64 string) (*SliderCaptcha, error) {
	slider, err := decodeBase64Image(sliderBase64)
	if err != nil {
		return nil, err
	}
	background, err := decodeBase64Image(backgroundBase64)
	if err != nil {
		return nil, err
	}
	return &SliderCaptcha{Slider: slider, Background: background}, nil
}

// decodeBase64Image 将 Base64 字符串解码为图像
func decodeBase64Image(data string) (image.Image, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(decoded))
	return img, err
}

// FindBestMatch 寻找滑块在背景中的最佳匹配位置
// FindBestMatch 寻找滑块在背景中的最佳匹配位置
//func (sc *SliderCaptcha) FindBestMatch() (float64, float64) {
//	bestX, bestY := 0.0, 0.0
//	bestScore := 0.0
//	stepSize := 0.01 // 设定更小的步长，精度提高
//
//	for y := 0.0; y <= float64(sc.Background.Bounds().Dy())-float64(sc.Slider.Bounds().Dy()); y += stepSize {
//		for x := 0.0; x <= float64(sc.Background.Bounds().Dx())-float64(sc.Slider.Bounds().Dx()); x += stepSize {
//			score := sc.computeMatch(int(x), int(y))
//			if score > bestScore {
//				bestScore = score
//				bestX, bestY = x, y
//			}
//		}
//	}
//
//	// 为了引入随机性，你可以加一个小的随机扰动
//	rand.Seed(time.Now().UnixNano())
//	randomFactor := 0.0001 * (rand.Float64() - 0.5) // 随机微调
//	bestX += randomFactor
//	bestY += randomFactor
//
//	global.Log.Info(fmt.Sprintf("Best match found at: (%.10f, %.10f) with score: %.4f", bestX, bestY, bestScore))
//	return bestX, bestY
//}

func (sc *SliderCaptcha) FindBestMatch() (int, int) {
	bestX, bestY := 0, 0
	bestScore := 0.0
	for y := 0; y <= sc.Background.Bounds().Dy()-sc.Slider.Bounds().Dy(); y++ {
		for x := 0; x <= sc.Background.Bounds().Dx()-sc.Slider.Bounds().Dx(); x++ {
			score := sc.computeMatch(x, y)
			if score > bestScore {
				bestScore = score
				bestX, bestY = x, y
			}
		}
	}
	global.Log.Info(fmt.Sprintf("Best match found at: (%d, %d) with score: %.4f", bestX, bestY, bestScore))
	return bestX, bestY
}

// computeMatch 计算滑块在背景指定位置的匹配度
func (sc *SliderCaptcha) computeMatch(offsetX, offsetY int) float64 {
	var total, match float64
	for y := 0; y < sc.Slider.Bounds().Dy(); y++ {
		for x := 0; x < sc.Slider.Bounds().Dx(); x++ {
			bgX, bgY := x+offsetX, y+offsetY
			if bgX >= sc.Background.Bounds().Dx() || bgY >= sc.Background.Bounds().Dy() {
				continue
			}

			sliderPixel := color.RGBAModel.Convert(sc.Slider.At(x, y)).(color.RGBA)
			bgPixel := color.RGBAModel.Convert(sc.Background.At(bgX, bgY)).(color.RGBA)

			// 跳过透明像素
			if sliderPixel.A > 0 {
				total++
				if math.Abs(float64(sliderPixel.R)-float64(bgPixel.R)) < 30 &&
					math.Abs(float64(sliderPixel.G)-float64(bgPixel.G)) < 30 &&
					math.Abs(float64(sliderPixel.B)-float64(bgPixel.B)) < 30 {
					match++
				}
			}
		}
	}
	if total == 0 {
		return 0
	}
	return match / total
}

// SaveImage 保存图像到文件，便于调试
func SaveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}

//package blockPuzzle
//
//import (
//	"BronyaBot/global"
//	"bytes"
//	"encoding/base64"
//	"fmt"
//	"image"
//	"image/color"
//	"image/png"
//	"log"
//	"math"
//	"os"
//)
//
//// SliderCaptcha 用于封装滑块验证码的处理逻辑
//type SliderCaptcha struct {
//	Slider     image.Image // 滑块图片
//	Background image.Image // 背景图片
//}
//
//// NewSliderCaptcha 初始化 SliderCaptcha
//func NewSliderCaptcha(sliderBase64, backgroundBase64 string) (*SliderCaptcha, error) {
//	slider, err := decodeBase64Image(sliderBase64)
//	if err != nil {
//		return nil, err
//	}
//	background, err := decodeBase64Image(backgroundBase64)
//	if err != nil {
//		return nil, err
//	}
//	return &SliderCaptcha{Slider: slider, Background: background}, nil
//}
//
//// decodeBase64Image 将 Base64 字符串解码为图像
//func decodeBase64Image(data string) (image.Image, error) {
//	decoded, err := base64.StdEncoding.DecodeString(data)
//	if err != nil {
//		return nil, err
//	}
//	img, _, err := image.Decode(bytes.NewReader(decoded))
//	return img, err
//}
//
//// FindBestMatch 通过模板匹配法寻找滑块在背景中的最佳匹配位置
//func (sc *SliderCaptcha) FindBestMatch() (int, int) {
//	bestX, bestY := 0, 0
//	bestScore := math.MaxFloat64 // 初始化最优得分为最大值
//
//	// 遍历背景图像的每一个位置，计算与滑块的差异
//	for y := 0; y <= sc.Background.Bounds().Dy()-sc.Slider.Bounds().Dy(); y++ {
//		for x := 0; x <= sc.Background.Bounds().Dx()-sc.Slider.Bounds().Dx(); x++ {
//			score := sc.computeTemplateMatch(x, y)
//			if score < bestScore {
//				bestScore = score
//				bestX, bestY = x, y
//			}
//		}
//	}
//
//	// 输出最佳匹配的坐标及得分
//	global.Log.Info(fmt.Sprintf("Best match found at: (%d, %d) with score: %.4f", bestX, bestY, bestScore))
//	return bestX, bestY
//}
//
//// computeTemplateMatch 计算滑块与背景在指定位置的匹配度，基于模板匹配算法
//func (sc *SliderCaptcha) computeTemplateMatch(offsetX, offsetY int) float64 {
//	var total, diff float64
//
//	// 遍历滑块的每个像素，与背景图对应位置的像素比较
//	for y := 0; y < sc.Slider.Bounds().Dy(); y++ {
//		for x := 0; x < sc.Slider.Bounds().Dx(); x++ {
//			bgX, bgY := x+offsetX, y+offsetY
//			if bgX >= sc.Background.Bounds().Dx() || bgY >= sc.Background.Bounds().Dy() {
//				continue
//			}
//
//			// 获取滑块和背景的颜色
//			sliderPixel := color.RGBAModel.Convert(sc.Slider.At(x, y)).(color.RGBA)
//			bgPixel := color.RGBAModel.Convert(sc.Background.At(bgX, bgY)).(color.RGBA)
//
//			// 如果滑块像素是透明的，使用背景像素
//			if sliderPixel.A > 0 { // 仅当滑块像素不透明时参与计算
//				// 计算像素间的差异
//				colorDiff := math.Abs(float64(sliderPixel.R)-float64(bgPixel.R)) +
//					math.Abs(float64(sliderPixel.G)-float64(bgPixel.G)) +
//					math.Abs(float64(sliderPixel.B)-float64(bgPixel.B))
//
//				// 累加差异
//				total++
//				diff += colorDiff
//			}
//		}
//	}
//
//	// 返回平均差异，差异越小，匹配度越高
//	if total == 0 {
//		return math.MaxFloat64
//	}
//	return diff / total
//}
//
//// SaveImage 保存图像到文件，便于调试
//func SaveImage(img image.Image, filename string) {
//	file, err := os.Create(filename)
//	if err != nil {
//		log.Fatalf("Failed to create file: %v", err)
//	}
//	defer file.Close()
//	if err := png.Encode(file, img); err != nil {
//		log.Fatalf("Failed to encode image: %v", err)
//	}
//}

//package blockPuzzle
//
//import (
//	"BronyaBot/global"
//	"bytes"
//	"encoding/base64"
//	"fmt"
//	"image"
//	"image/color"
//	"image/png"
//	"log"
//	"math"
//	"os"
//)
//
//// SliderCaptcha 用于封装滑块验证码的处理逻辑
//type SliderCaptcha struct {
//	Slider     image.Image // 滑块图片
//	Background image.Image // 背景图片
//}
//
//// NewSliderCaptcha 初始化 SliderCaptcha
//func NewSliderCaptcha(sliderBase64, backgroundBase64 string) (*SliderCaptcha, error) {
//	slider, err := decodeBase64Image(sliderBase64)
//	if err != nil {
//		return nil, err
//	}
//	background, err := decodeBase64Image(backgroundBase64)
//	if err != nil {
//		return nil, err
//	}
//	return &SliderCaptcha{Slider: slider, Background: background}, nil
//}
//
//// decodeBase64Image 将 Base64 字符串解码为图像
//func decodeBase64Image(data string) (image.Image, error) {
//	decoded, err := base64.StdEncoding.DecodeString(data)
//	if err != nil {
//		return nil, err
//	}
//	img, _, err := image.Decode(bytes.NewReader(decoded))
//	return img, err
//}
//
//// FindBestMatch 寻找滑块在背景中的最佳匹配位置
//func (sc *SliderCaptcha) FindBestMatch() (int, int) {
//	bestX, bestY := 0, 0
//	bestScore := 0.0
//
//	// 找到滑块在背景中的最佳匹配位置
//	for y := 0; y <= sc.Background.Bounds().Dy()-sc.Slider.Bounds().Dy(); y++ {
//		for x := 0; x <= sc.Background.Bounds().Dx()-sc.Slider.Bounds().Dx(); x++ {
//			score := sc.computeMatch(x, y)
//			if score > bestScore {
//				bestScore = score
//				bestX, bestY = x, y
//			}
//		}
//	}
//
//	global.Log.Info(fmt.Sprintf("Best match found at: (%d, %d) with score: %.4f", bestX, bestY, bestScore))
//	return bestX, bestY
//}
//
//// computeMatch 计算滑块在背景指定位置的匹配度
//func (sc *SliderCaptcha) computeMatch(offsetX, offsetY int) float64 {
//	var total, match float64
//
//	// 只匹配滑块的非透明区域，且忽略背景中的纯白色区域（#ffffff）
//	for y := 0; y < sc.Slider.Bounds().Dy(); y++ {
//		for x := 0; x < sc.Slider.Bounds().Dx(); x++ {
//			bgX, bgY := x+offsetX, y+offsetY
//			if bgX >= sc.Background.Bounds().Dx() || bgY >= sc.Background.Bounds().Dy() {
//				continue
//			}
//
//			sliderPixel := color.RGBAModel.Convert(sc.Slider.At(x, y)).(color.RGBA)
//			bgPixel := color.RGBAModel.Convert(sc.Background.At(bgX, bgY)).(color.RGBA)
//
//			// 忽略滑块中的透明像素（A=0）和背景中的白色像素（#ffffff）
//			if sliderPixel.A > 0 && (sliderPixel.R != 255 || sliderPixel.G != 255 || sliderPixel.B != 255) {
//				total++
//				// 计算颜色差异，容忍一定的差异
//				if math.Abs(float64(sliderPixel.R)-float64(bgPixel.R)) < 30 &&
//					math.Abs(float64(sliderPixel.G)-float64(bgPixel.G)) < 30 &&
//					math.Abs(float64(sliderPixel.B)-float64(bgPixel.B)) < 30 {
//					match++
//				}
//			}
//		}
//	}
//
//	// 返回匹配度
//	if total == 0 {
//		return 0
//	}
//	return match / total
//}
//
//// SaveImage 保存图像到文件，便于调试
//func SaveImage(img image.Image, filename string) {
//	file, err := os.Create(filename)
//	if err != nil {
//		log.Fatalf("Failed to create file: %v", err)
//	}
//	defer file.Close()
//	if err := png.Encode(file, img); err != nil {
//		log.Fatalf("Failed to encode image: %v", err)
//	}
//}
