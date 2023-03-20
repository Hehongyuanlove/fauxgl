package fauxgl

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

var (
	Discard     = Color{}           // 丢弃
	Transparent = Color{}           // 透明
	Black       = Color{0, 0, 0, 1} // 黑色
	White       = Color{1, 1, 1, 1} // 白色
)

type Color struct {
	R, G, B, A float64 // 颜色的 RGBA 值
}

// Gray 生成灰度颜色
func Gray(x float64) Color {
	return Color{x, x, x, 1}
}

// MakeColor 生成颜色
func MakeColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	const d = 0xffff
	return Color{float64(r) / d, float64(g) / d, float64(b) / d, float64(a) / d}
}

// HexColor 生成十六进制颜色
func HexColor(x string) Color {
	x = strings.Trim(x, "#")
	var r, g, b, a int
	a = 255
	switch len(x) {
	case 3:
		fmt.Sscanf(x, "%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 4:
		fmt.Sscanf(x, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 6:
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	const d = 0xff
	return Color{float64(r) / d, float64(g) / d, float64(b) / d, float64(a) / d}
}

// NRGBA 返回颜色的 NRGBA 值
func (c Color) NRGBA() color.NRGBA {
	const d = 0xff
	r := Clamp(c.R, 0, 1)
	g := Clamp(c.G, 0, 1)
	b := Clamp(c.B, 0, 1)
	a := Clamp(c.A, 0, 1)
	return color.NRGBA{uint8(r * d), uint8(g * d), uint8(b * d), uint8(a * d)}
}

// Opaque 返回不透明的颜色
func (a Color) Opaque() Color {
	return Color{a.R, a.G, a.B, 1}
}

// Alpha 返回指定透明度的颜色
func (a Color) Alpha(alpha float64) Color {
	return Color{a.R, a.G, a.B, alpha}
}

// Lerp 返回两个颜色之间的插值
func (a Color) Lerp(b Color, t float64) Color {
	return a.Add(b.Sub(a).MulScalar(t))
}

// Add 返回两个颜色的和
func (a Color) Add(b Color) Color {
	return Color{a.R + b.R, a.G + b.G, a.B + b.B, a.A + b.A}
}

// Sub 返回两个颜色的差
func (a Color) Sub(b Color) Color {
	return Color{a.R - b.R, a.G - b.G, a.B - b.B, a.A - b.A}
}

// Mul 返回两个颜色的积
func (a Color) Mul(b Color) Color {
	return Color{a.R * b.R, a.G * b.G, a.B * b.B, a.A * b.A}
}

// Div 返回两个颜色的商
func (a Color) Div(b Color) Color {
	return Color{a.R / b.R, a.G / b.G, a.B / b.B, a.A / b.A}
}

// AddScalar 返回颜色与标量的和
func (a Color) AddScalar(b float64) Color {
	return Color{a.R + b, a.G + b, a.B + b, a.A + b}
}

// SubScalar 返回颜色与标量的差
func (a Color) SubScalar(b float64) Color {
	return Color{a.R - b, a.G - b, a.B - b, a.A - b}
}

// MulScalar 返回颜色与标量的积
func (a Color) MulScalar(b float64) Color {
	return Color{a.R * b, a.G * b, a.B * b, a.A * b}
}

// DivScalar 返回颜色与标量的商
func (a Color) DivScalar(b float64) Color {
	return Color{a.R / b, a.G / b, a.B / b, a.A / b}
}

// Pow 返回颜色的幂
func (a Color) Pow(b float64) Color {
	return Color{math.Pow(a.R, b), math.Pow(a.G, b), math.Pow(a.B, b), math.Pow(a.A, b)}
}

// Min 返回两个颜色的最小值
func (a Color) Min(b Color) Color {
	return Color{math.Min(a.R, b.R), math.Min(a.G, b.G), math.Min(a.B, b.B), math.Min(a.A, b.A)}
}

// Max 返回两个颜色的最大值
func (a Color) Max(b Color) Color {
	return Color{math.Max(a.R, b.R), math.Max(a.G, b.G), math.Max(a.B, b.B), math.Max(a.A, b.A)}
}
