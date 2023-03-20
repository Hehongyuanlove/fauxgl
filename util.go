package fauxgl

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 角度转弧度
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// 弧度转角度
func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// 经纬度转空间坐标
func LatLngToXYZ(lat, lng float64) Vector {
	lat, lng = Radians(lat), Radians(lng)
	x := math.Cos(lat) * math.Cos(lng)
	y := math.Cos(lat) * math.Sin(lng)
	z := math.Sin(lat)
	return Vector{x, y, z}
}

// 加载模型
func LoadMesh(path string) (*Mesh, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".stl":
		return LoadSTL(path)
	case ".obj":
		return LoadOBJ(path)
	case ".ply":
		return LoadPLY(path)
	case ".3ds":
		return Load3DS(path)
	}
	return nil, fmt.Errorf("unrecognized mesh extension: %s", ext)
}

// 加载图片
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

// 保存图片
func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

// 解析浮点数
func ParseFloats(items []string) []float64 {
	result := make([]float64, len(items))
	for i, item := range items {
		f, _ := strconv.ParseFloat(item, 64)
		result[i] = f
	}
	return result
}

// 限制范围
func Clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// 限制范围
func ClampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// 取绝对值
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 四舍五入
func Round(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	} else {
		return int(math.Floor(a + 0.5))
	}
}

// 四舍五入到指定小数位
func RoundPlaces(a float64, places int) float64 {
	shift := powersOfTen[places]
	return float64(Round(a*shift)) / shift
}

var powersOfTen = []float64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16}
