package fauxgl

import (
	"image"
	"math"
)
// Texture 纹理接口
type Texture interface {
	// Sample 对纹理进行采样
	// u, v: 纹理坐标
	// 返回值: 采样得到的颜色
	Sample(u, v float64) Color

	// BilinearSample 使用双线性插值对纹理进行采样
	// u, v: 纹理坐标
	// 返回值: 采样得到的颜色
	BilinearSample(u, v float64) Color
}


// LoadTexture 加载纹理
// path: 纹理文件路径
// 返回值: 纹理对象和错误信息
func LoadTexture(path string) (Texture, error) {
	im, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return NewImageTexture(im), nil
}

// ImageTexture 纹理对象
type ImageTexture struct {
	Width  int
	Height int
	Image  image.Image
}

// NewImageTexture 创建纹理对象
// im: 图像对象
// 返回值: 纹理对象
func NewImageTexture(im image.Image) Texture {
	size := im.Bounds().Max
	return &ImageTexture{size.X, size.Y, im}
}

// Sample 对纹理进行采样
// u, v: 纹理坐标
// 返回值: 采样得到的颜色
func (t *ImageTexture) Sample(u, v float64) Color {
	v = 1 - v
	u -= math.Floor(u)
	v -= math.Floor(v)
	x := int(u * float64(t.Width))
	y := int(v * float64(t.Height))
	return MakeColor(t.Image.At(x, y))
}

// BilinearSample 使用双线性插值对纹理进行采样
// u, v: 纹理坐标
// 返回值: 采样得到的颜色
func (t *ImageTexture) BilinearSample(u, v float64) Color {
	v = 1 - v
	u -= math.Floor(u)
	v -= math.Floor(v)
	x := u * float64(t.Width-1)
	y := v * float64(t.Height-1)
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	x -= float64(x0)
	y -= float64(y0)
	c00 := MakeColor(t.Image.At(x0, y0))
	c01 := MakeColor(t.Image.At(x0, y1))
	c10 := MakeColor(t.Image.At(x1, y0))
	c11 := MakeColor(t.Image.At(x1, y1))
	c := Color{}
	c = c.Add(c00.MulScalar((1 - x) * (1 - y)))
	c = c.Add(c10.MulScalar(x * (1 - y)))
	c = c.Add(c01.MulScalar((1 - x) * y))
	c = c.Add(c11.MulScalar(x * y))
	return c
}
