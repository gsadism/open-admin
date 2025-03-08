package image

import (
	"bytes"
	"encoding/base64"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"
	"time"
)

const (
	DefaultDPI    = 75
	DefaultSize   = 20
	DefaultWith   = 120
	DefaultHeight = 40
)

const (
	Number int = iota
	Upper
	Lower
	Mixed
)

type GG struct {
	fontFile string
}

func NewGG(fontFile string) *GG {
	g := &GG{fontFile: fontFile}
	return g
}

func (t *GG) randStr(n int) (w []interface{}) {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz123456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	var randStr string
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}

	for s := 0; s < n; s++ {
		w = append(w, string(randStr[s]))
	}
	return w
}

func (t *GG) randNumber(n int) (w []interface{}) {
	chars := "0123456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	var randStr string
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}

	for s := 0; s < n; s++ {
		w = append(w, string(randStr[s]))
	}
	return w
}

func (t *GG) randUpperStr(n int) (w []interface{}) {
	chars := "ABCDEFGHIJKMNOPQRSTUVWXYZ"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	var randStr string
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}

	for s := 0; s < n; s++ {
		w = append(w, string(randStr[s]))
	}
	return w
}

func (t *GG) randLowerStr(n int) (w []interface{}) {
	chars := "abcdefghijklmnopqrstuvwxyz"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	var randStr string
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}

	for s := 0; s < n; s++ {
		w = append(w, string(randStr[s]))
	}
	return w
}

func (t *GG) randPos(width, height int) (float64, float64) {
	x := rand.Float64() * float64(width)
	y := rand.Float64() * float64(height)
	return x, y
}

func (t *GG) randColor(maxColor int) (int, int, int, int) {
	r := int(uint8(rand.Intn(maxColor)))
	g := int(uint8(rand.Intn(maxColor)))
	b := int(uint8(rand.Intn(maxColor)))
	a := int(uint8(rand.Intn(maxColor)))
	return r, g, b, a
}

func (t *GG) degreeSin(degree float64) float64 {
	return math.Sin(degree * math.Pi / 180)
}

func (t *GG) randInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func (t *GG) ifn(expr bool, trueVal float64, falseVal float64) float64 {
	if expr {
		return trueVal
	}
	return falseVal
}

func (t *GG) contain(obj interface{}, target interface{}) bool {
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
	default:
		panic("非法参数")
	}

	return false
}

func (t *GG) randomMembersFromMemberLibrary(lib []interface{}, size int) []interface{} {
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

func (t *GG) cordByQuadrantAndDegree(w, h, ascent, descent, degree, beforeTotalWidth float64) (x, y, leftCutSize, rightCutSize float64) {
	var totalWidth float64
	switch {
	case degree <= 0 && degree >= -40: // 第一象限：逆时针 -30度 ~ 0  <=>  330 ~ 360 （目前参数要传入负数）
		rd := -1 * degree // 转为正整数，便于计算
		leftCutSize = w * t.degreeSin(90-rd)
		rightCutSize = h * t.degreeSin(rd)

		offset := (leftCutSize + rightCutSize - w) / 2 // 横向偏移量（角度倾斜越厉害，占宽越多，通过偏移量分摊给它的左右边距来收窄）
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforeTotalWidth + leftCutSize
		x = t.degreeSin(90-rd)*totalWidth - w
		y = ascent + t.degreeSin(rd)*totalWidth
	case degree >= 0 && degree <= 40: // 第四象限：顺时针 0 ~ 30度
		leftCutSize = h * t.degreeSin(degree)
		rightCutSize = w * t.degreeSin(90-degree)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforeTotalWidth + leftCutSize // 现在totalwidth = 前面的宽 + 自己的左切边
		x = t.degreeSin(90-degree) * totalWidth
		y = ascent - t.degreeSin(degree)*totalWidth
	case degree >= 180 && degree <= 220: // 第二象限：顺时针 180 ~ 210度
		rd := degree - 180
		leftCutSize = h * t.degreeSin(rd)
		rightCutSize = w * t.degreeSin(90-rd)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforeTotalWidth + leftCutSize
		x = -1 * (t.degreeSin(90-rd)*totalWidth + w)
		y = t.degreeSin(rd)*totalWidth - descent
	case degree >= 140 && degree <= 180: // 第三象限：顺时针 150 ~ 180度
		rd := 180 - degree
		leftCutSize = w * t.degreeSin(90-rd)
		rightCutSize = h * t.degreeSin(rd)

		offset := (leftCutSize + rightCutSize - w) / 2
		leftCutSize, rightCutSize = leftCutSize-offset, rightCutSize-offset

		totalWidth = beforeTotalWidth + leftCutSize
		x = -1 * (t.degreeSin(90-rd) * totalWidth)
		y = -1 * (t.degreeSin(rd)*totalWidth + descent)
	default:
		panic(any("非法参数"))
	}
	return
}

func (t *GG) drawImageByGgToBase64(species int, n, width, height int, dpi float64, size float64) (string, string, error) {
	if n <= 0 {
		n = 4
	}
	if n >= 6 {
		n = 6
	}
	dc := gg.NewContext(width, height)
	dc.SetRGB255(255, 255, 255) // 设置背景色：末尾为透明度 1-0(1-不透明 0-透明)
	dc.Clear()
	//dc.SetRGBA(0, 9, 7, 1) // 设置字体色
	fontBytes, err := ioutil.ReadFile(t.fontFile)
	if err != nil {
		return "", "", err
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return "", "", err
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: size,
		DPI:  dpi,
	})

	// 干扰线
	for i := 0; i < 6; i++ {
		x1, y1 := t.randPos(width, height)
		x2, y2 := t.randPos(width, height)
		r, g, b, a := t.randColor(255)
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
	code := t.randStr(n)
	switch species {
	case Number:
		code = t.randNumber(n)
	case Upper:
		code = t.randUpperStr(n)
	case Lower:
		code = t.randLowerStr(n)
	case Mixed:
		code = t.randStr(n)
	default:
		code = t.randStr(n)
	}
	words := t.randomMembersFromMemberLibrary(code, n)
	reverseWordsIndex := t.randomMembersFromMemberLibrary([]interface{}{0, 1, 2, 3, 4, 5}, 1) // 随机2个倒立字

	for i, word := range words {
		degree := t.ifn(t.contain(i, reverseWordsIndex), float64(t.randInt64(150, 210)), float64(t.randInt64(-30, 30))) // 随机角度，正向角度 -30~30，倒立角度 150~210
		x, y, leftCutSize, rightCS := t.cordByQuadrantAndDegree(w, h, ascent, decent, degree, totalWidth)
		dc.RotateAbout(gg.Radians(degree), 0, 0)
		dc.DrawStringAnchored(word.(string), x, y, 0, 0)
		dc.RotateAbout(-1*gg.Radians(degree), 0, 0)
		totalWidth = totalWidth + leftCutSize + rightCS
	}

	dc.Stroke()
	img := dc.Image()
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	var b bytes.Buffer
	err = encoder.Encode(&b, img)
	if err != nil {
		return "", "", err
	}
	res := base64.StdEncoding.EncodeToString(b.Bytes())

	// 验证码
	var codeStr string
	for _, v := range code {
		str, ok := v.(string)
		if ok {
			codeStr += str
		}
	}
	return codeStr, "data:image/jpeg;base64," + res, nil
}

func (t *GG) Image(spec, n, width, height int, dpi, size float64) (string, string, error) {
	if code, img, err := t.drawImageByGgToBase64(spec, n, width, height, dpi, size); err != nil {
		return "", "", err
	} else {
		return code, img, err
	}
}

func (t *GG) Default(spec, n int) (string, string, error) {
	if code, img, err := t.drawImageByGgToBase64(spec, n, DefaultWith, DefaultHeight, DefaultDPI, DefaultSize); err != nil {
		return "", "", err
	} else {
		return code, img, err
	}
}
