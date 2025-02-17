package fauxgl

import (
	"math"
	"math/rand"
)

// Vector结构体表示一个三维向量
type Vector struct {
	X, Y, Z float64 // X、Y、Z分别表示三维向量的三个分量
}

// V函数返回一个新的Vector
func V(x, y, z float64) Vector {
	return Vector{x, y, z}
}

// RandomUnitVector函数返回一个随机的单位向量
func RandomUnitVector() Vector {
	for {
		x := rand.Float64()*2 - 1
		y := rand.Float64()*2 - 1
		z := rand.Float64()*2 - 1
		if x*x+y*y+z*z > 1 {
			continue
		}
		return Vector{x, y, z}.Normalize()
	}
}

// VectorW函数返回一个新的VectorW
func (a Vector) VectorW() VectorW {
	return VectorW{a.X, a.Y, a.Z, 1}
}

// IsDegenerate函数判断向量是否退化
func (a Vector) IsDegenerate() bool {
	nan := math.IsNaN(a.X) || math.IsNaN(a.Y) || math.IsNaN(a.Z)
	inf := math.IsInf(a.X, 0) || math.IsInf(a.Y, 0) || math.IsInf(a.Z, 0)
	return nan || inf
}

// Length函数返回向量的长度
func (a Vector) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Less函数比较两个向量的大小
func (a Vector) Less(b Vector) bool {
	if a.X != b.X {
		return a.X < b.X
	}
	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.Z < b.Z
}

// Distance函数返回两个向量之间的距离
func (a Vector) Distance(b Vector) float64 {
	return a.Sub(b).Length()
}

// LengthSquared函数返回向量的长度的平方
func (a Vector) LengthSquared() float64 {
	return a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

// DistanceSquared函数返回两个向量之间距离的平方
func (a Vector) DistanceSquared(b Vector) float64 {
	return a.Sub(b).LengthSquared()
}

// Lerp函数返回两个向量之间的插值
func (a Vector) Lerp(b Vector, t float64) Vector {
	return a.Add(b.Sub(a).MulScalar(t))
}

// LerpDistance函数返回两个向量之间距离的插值
func (a Vector) LerpDistance(b Vector, d float64) Vector {
	return a.Add(b.Sub(a).Normalize().MulScalar(d))
}

// Dot函数返回两个向量的点积
func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross函数返回两个向量的叉积
func (a Vector) Cross(b Vector) Vector {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return Vector{x, y, z}
}

// Normalize函数返回向量的单位向量
func (a Vector) Normalize() Vector {
	r := 1 / math.Sqrt(a.X*a.X+a.Y*a.Y+a.Z*a.Z)
	return Vector{a.X * r, a.Y * r, a.Z * r}
}

// Negate函数返回向量的相反向量
func (a Vector) Negate() Vector {
	return Vector{-a.X, -a.Y, -a.Z}
}

// Abs函数返回向量的绝对值
func (a Vector) Abs() Vector {
	return Vector{math.Abs(a.X), math.Abs(a.Y), math.Abs(a.Z)}
}

// Add函数返回两个向量的和
func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Sub函数返回两个向量的差
func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Mul函数返回两个向量的积
func (a Vector) Mul(b Vector) Vector {
	return Vector{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Div函数返回两个向量的商
func (a Vector) Div(b Vector) Vector {
	return Vector{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

// Mod函数返回两个向量的模
func (a Vector) Mod(b Vector) Vector {
	// as implemented in GLSL
	x := a.X - b.X*math.Floor(a.X/b.X)
	y := a.Y - b.Y*math.Floor(a.Y/b.Y)
	z := a.Z - b.Z*math.Floor(a.Z/b.Z)
	return Vector{x, y, z}
}

// AddScalar函数返回向量加上标量的结果
func (a Vector) AddScalar(b float64) Vector {
	return Vector{a.X + b, a.Y + b, a.Z + b}
}

// SubScalar函数返回向量减去标量的结果
func (a Vector) SubScalar(b float64) Vector {
	return Vector{a.X - b, a.Y - b, a.Z - b}
}

// MulScalar函数返回向量乘以标量的结果
func (a Vector) MulScalar(b float64) Vector {
	return Vector{a.X * b, a.Y * b, a.Z * b}
}

// DivScalar函数返回向量除以标量的结果
func (a Vector) DivScalar(b float64) Vector {
	return Vector{a.X / b, a.Y / b, a.Z / b}
}

// Min函数返回两个向量的最小值
func (a Vector) Min(b Vector) Vector {
	return Vector{math.Min(a.X, b.X), math.Min(a.Y, b.Y), math.Min(a.Z, b.Z)}
}

// Max函数返回两个向量的最大值
func (a Vector) Max(b Vector) Vector {
	return Vector{math.Max(a.X, b.X), math.Max(a.Y, b.Y), math.Max(a.Z, b.Z)}
}

// Floor函数返回向量的下取整
func (a Vector) Floor() Vector {
	return Vector{math.Floor(a.X), math.Floor(a.Y), math.Floor(a.Z)}
}

// Ceil函数返回向量的上取整
func (a Vector) Ceil() Vector {
	return Vector{math.Ceil(a.X), math.Ceil(a.Y), math.Ceil(a.Z)}
}

// Round函数返回向量的四舍五入
func (a Vector) Round() Vector {
	return a.RoundPlaces(0)
}

// RoundPlaces函数返回向量的四舍五入到指定小数位数
func (a Vector) RoundPlaces(n int) Vector {
	x := RoundPlaces(a.X, n)
	y := RoundPlaces(a.Y, n)
	z := RoundPlaces(a.Z, n)
	return Vector{x, y, z}
}

// MinComponent函数返回向量的最小分量
func (a Vector) MinComponent() float64 {
	return math.Min(math.Min(a.X, a.Y), a.Z)
}

// MaxComponent函数返回向量的最大分量
func (a Vector) MaxComponent() float64 {
	return math.Max(math.Max(a.X, a.Y), a.Z)
}

// Reflect函数返回向量在法向量上的反射向量
func (i Vector) Reflect(n Vector) Vector {
	return i.Sub(n.MulScalar(2 * n.Dot(i)))
}

// Perpendicular函数返回向量的垂直向量
func (a Vector) Perpendicular() Vector {
	if a.X == 0 && a.Y == 0 {
		if a.Z == 0 {
			return Vector{}
		}
		return Vector{0, 1, 0}
	}
	return Vector{-a.Y, a.X, 0}.Normalize()
}

// SegmentDistance函数返回点到线段的距离
func (p Vector) SegmentDistance(v Vector, w Vector) float64 {
	l2 := v.DistanceSquared(w)
	if l2 == 0 {
		return p.Distance(v)
	}
	t := p.Sub(v).Dot(w.Sub(v)) / l2
	if t < 0 {
		return p.Distance(v)
	}
	if t > 1 {
		return p.Distance(w)
	}
	return v.Add(w.Sub(v).MulScalar(t)).Distance(p)
}

// VectorW结构体表示一个四维向量
type VectorW struct {
	X, Y, Z, W float64 // X、Y、Z、W分别表示四维向量的四个分量
}

// Vector函数返回VectorW的Vector部分
func (a VectorW) Vector() Vector {
	return Vector{a.X, a.Y, a.Z}
}

// Outside函数判断VectorW是否在超出范围
func (a VectorW) Outside() bool {
	x, y, z, w := a.X, a.Y, a.Z, a.W
	return x < -w || x > w || y < -w || y > w || z < -w || z > w
}

// Dot函数返回两个VectorW的点积
func (a VectorW) Dot(b VectorW) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

// Add函数返回两个VectorW的和
func (a VectorW) Add(b VectorW) VectorW {
	return VectorW{a.X + b.X, a.Y + b.Y, a.Z + b.Z, a.W + b.W}
}

// Sub函数返回两个VectorW的差
func (a VectorW) Sub(b VectorW) VectorW {
	return VectorW{a.X - b.X, a.Y - b.Y, a.Z - b.Z, a.W - b.W}
}

// MulScalar函数返回VectorW乘以标量的结果
func (a VectorW) MulScalar(b float64) VectorW {
	return VectorW{a.X * b, a.Y * b, a.Z * b, a.W * b}
}

// DivScalar函数返回VectorW除以标量的结果
func (a VectorW) DivScalar(b float64) VectorW {
	return VectorW{a.X / b, a.Y / b, a.Z / b, a.W / b}
}
