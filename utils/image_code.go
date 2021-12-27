package utils

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"reflect"
	"time"
	//"io"
)

type ImageCode struct {
}

func GetRandStr(n int) (w []interface{}) {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	rand.Seed(time.Now().UnixNano())
	var randStr string
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	fmt.Println(randStr)

	for s := 0; s < n; s++ {
		w = append(w, string(randStr[s]))
	}
	return w
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "./SIMYOU.TTF", "filename of the ttf font")
	size     = flag.Float64("size", 20, "font size in points")
)

// 随机坐标
func getRandPos(width, height int) (x float64, y float64) {
	x = rand.Float64() * float64(width)
	y = rand.Float64() * float64(height)
	return x, y
}

// 随机颜色
func getRandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(rand.Intn(maxColor)))
	g = int(uint8(rand.Intn(maxColor)))
	b = int(uint8(rand.Intn(maxColor)))
	a = int(uint8(rand.Intn(255)))
	return r, g, b, a
}
func drawImageBygg(width, height int) {
	dc := gg.NewContext(width, height) // 56 => w*sin(45) + h*sin(45)  45度时，字体达到最大高度
	dc.SetRGB255(255, 255, 255)        // 设置背景色：末尾为透明度 1-0(1-不透明 0-透明)
	dc.Clear()
	//dc.SetRGBA(0, 9, 7, 1) // 设置字体色

	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: *size,
		DPI:  *dpi,
	})
	// 干扰线
	for i := 0; i < 6; i++ {
		x1, y1 := getRandPos(width, height)
		x2, y2 := getRandPos(width, height)
		r, g, b, a := getRandColor(255)
		w := float64(rand.Intn(3) + 1)
		dc.SetRGBA255(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}
	dc.SetFontFace(face)
	dc.SetRGBA(90, 0, 0, 1) // 设置字体色

	// 初始化用于计算坐标的变量
	fm := face.Metrics()
	ascent := float64(fm.Ascent.Round())  // 字体的基线到顶部距离
	decent := float64(fm.Descent.Round()) // 字体的基线到底部的距离
	w := float64(fm.Height.Round())       // 方块字，大多数应为等宽字，即和高度一样
	h := float64(fm.Height.Round())
	totalWidth := 0.0 // 目前已累积的图片宽度（需要用来计算字体位置）

	// 随机取汉字，定位倒立的字
	words := getRandomMembersFromMemberLibary(GetRandStr(6), 6)                               // 取8个字
	reverseWordsIndex := getRandomMembersFromMemberLibary([]interface{}{0, 1, 2, 3, 4, 5}, 1) // 随机2个倒立字

	for i, word := range words {
		degree := If(Contain(i, reverseWordsIndex), float64(RandInt64(150, 210)), float64(RandInt64(-30, 30))) // 随机角度，正向角度 -30~30，倒立角度 150~210
		x, y, leftCutSize, rightCS := getCoordByQuadrantAndDegree(w, h, ascent, decent, degree, totalWidth)
		dc.RotateAbout(gg.Radians(degree), 0, 0)
		dc.DrawStringAnchored(word.(string), x, y, 0, 0)
		dc.RotateAbout(-1*gg.Radians(degree), 0, 0)
		totalWidth = totalWidth + leftCutSize + rightCS
		fmt.Println("x:", x, "y:", y, "total:", totalWidth, "degree:", degree)
	}

	dc.Stroke()
	//writer := io.Writer()
	//dc.EncodePNG(writer)
	//dc.SavePNG("out.png")
	//dc.SaveJPG("out.jpeg", 200)
}

func getCoordByQuadrantAndDegree(w, h, ascent, descent, degree, beforTotalWidth float64) (x, y, leftCutSize, rightCutSize float64) {
	var totalWidth float64
	switch {
	case degree <= 0 && degree >= -40: // 第一象限：逆时针 -30度 ~ 0  <=>  330 ~ 360 （目前参数要传入负数）
		rd := -1 * degree // 转为正整数，便于计算
		leftCutSize = w * getDegreeSin(90-rd)
		rightCutSize = h * getDegreeSin(rd)

		offset := (leftCutSize + rightCutSize - w) / 2 // 横向偏移量（角度倾斜越厉害，占宽越多，通过偏移量分摊给它的左右边距来收窄）
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforTotalWidth + leftCutSize
		x = getDegreeSin(90-rd)*totalWidth - w
		y = ascent + getDegreeSin(rd)*totalWidth
	case degree >= 0 && degree <= 40: // 第四象限：顺时针 0 ~ 30度
		leftCutSize = h * getDegreeSin(degree)
		rightCutSize = w * getDegreeSin(90-degree)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforTotalWidth + leftCutSize // 现在totalwidth = 前面的宽 + 自己的左切边
		x = getDegreeSin(90-degree) * totalWidth
		y = ascent - getDegreeSin(degree)*totalWidth
	case degree >= 180 && degree <= 220: // 第二象限：顺时针 180 ~ 210度
		rd := degree - 180
		leftCutSize = h * getDegreeSin(rd)
		rightCutSize = w * getDegreeSin(90-rd)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforTotalWidth + leftCutSize
		x = -1 * (getDegreeSin(90-rd)*totalWidth + w)
		y = getDegreeSin(rd)*totalWidth - descent
	case degree >= 140 && degree <= 180: // 第三象限：顺时针 150 ~ 180度
		rd := 180 - degree
		leftCutSize = w * getDegreeSin(90-rd)
		rightCutSize = h * getDegreeSin(rd)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforTotalWidth + leftCutSize
		x = -1 * (getDegreeSin(90-rd) * totalWidth)
		y = -1 * (getDegreeSin(rd)*totalWidth + descent)
	default:
		panic(fmt.Sprintf("非法的参数：%f", degree))
	}
	return
}

func getDegreeSin(degree float64) float64 {
	return math.Sin(degree * math.Pi / 180)
}

//RandInt64 ...
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func getRandomMembersFromMemberLibary(lib []interface{}, size int) []interface{} {
	source, result := make([]interface{}, 0), make([]interface{}, 0)
	if size <= 0 || len(lib) == 0 {
		return result
	}
	for _, v := range lib {
		source = append(source, v)
	}
	if size >= len(lib) {
		return source
	}
	for i := 0; i < size; i++ {
		rand.Seed(time.Now().UnixNano())
		pos := rand.Intn(len(source))
		result = append(result, source[pos])
		source = append(source[:pos], source[pos+1:]...)
	}
	return result
}

//Contain ...
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

//If ...
func If(expr bool, trueVal float64, falseVal float64) float64 {
	if expr {
		return trueVal
	}
	return falseVal
}

func (ImageCode) GetCode() {
	drawImageBygg(120, 38)
}

