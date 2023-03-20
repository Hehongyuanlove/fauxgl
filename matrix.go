package fauxgl

import "math"

type Matrix struct {
	X00, X01, X02, X03 float64 // 第一行
	X10, X11, X12, X13 float64 // 第二行
	X20, X21, X22, X23 float64 // 第三行
	X30, X31, X32, X33 float64 // 第四行
}

// Identity 返回一个单位矩阵
func Identity() Matrix {
	return Matrix{
		1, 0, 0, 0, // 第一行
		0, 1, 0, 0, // 第二行
		0, 0, 1, 0, // 第三行
		0, 0, 0, 1} // 第四行
}

// Translate 返回一个平移矩阵
func Translate(v Vector) Matrix {
	return Matrix{
		1, 0, 0, v.X, // 第一行
		0, 1, 0, v.Y, // 第二行
		0, 0, 1, v.Z, // 第三行
		0, 0, 0, 1} // 第四行
}

// Scale 返回一个缩放矩阵
func Scale(v Vector) Matrix {
	return Matrix{
		v.X, 0, 0, 0, // 第一行
		0, v.Y, 0, 0, // 第二行
		0, 0, v.Z, 0, // 第三行
		0, 0, 0, 1} // 第四行
}

// Rotate 返回一个绕向量 v 旋转 a 角度的矩阵
func Rotate(v Vector, a float64) Matrix {
	v = v.Normalize()
	s := math.Sin(a)
	c := math.Cos(a)
	m := 1 - c
	return Matrix{
		m*v.X*v.X + c, m*v.X*v.Y + v.Z*s, m*v.Z*v.X - v.Y*s, 0, // 第一行
		m*v.X*v.Y - v.Z*s, m*v.Y*v.Y + c, m*v.Y*v.Z + v.X*s, 0, // 第二行
		m*v.Z*v.X + v.Y*s, m*v.Y*v.Z - v.X*s, m*v.Z*v.Z + c, 0, // 第三行
		0, 0, 0, 1} // 第四行
}

// RotateTo 返回一个将 a 向量旋转到 b 向量的矩阵
func RotateTo(a, b Vector) Matrix {
	dot := b.Dot(a)
	if dot == 1 {
		return Identity()
	} else if dot == -1 {
		return Rotate(a.Perpendicular(), math.Pi)
	} else {
		angle := math.Acos(dot)
		v := b.Cross(a).Normalize()
		return Rotate(v, angle)
	}
}

// Orient 返回一个定位矩阵
func Orient(position, size, up Vector, rotation float64) Matrix {
	m := Rotate(Vector{0, 0, 1}, rotation)
	m = m.Scale(size)
	m = m.RotateTo(Vector{0, 0, 1}, up)
	m = m.Translate(position)
	return m
}

// Frustum 返回一个透视投影矩阵
func Frustum(l, r, b, t, n, f float64) Matrix {
	t1 := 2 * n
	t2 := r - l
	t3 := t - b
	t4 := f - n
	return Matrix{
		t1 / t2, 0, (r + l) / t2, 0, // 第一行
		0, t1 / t3, (t + b) / t3, 0, // 第二行
		0, 0, (-f - n) / t4, (-t1 * f) / t4, // 第三行
		0, 0, -1, 0} // 第四行
}

// Orthographic 返回一个正交投影矩阵
func Orthographic(l, r, b, t, n, f float64) Matrix {
	return Matrix{
		2 / (r - l), 0, 0, -(r + l) / (r - l), // 第一行
		0, 2 / (t - b), 0, -(t + b) / (t - b), // 第二行
		0, 0, -2 / (f - n), -(f + n) / (f - n), // 第三行
		0, 0, 0, 1} // 第四行
}

// Perspective 返回一个透视投影矩阵
func Perspective(fovy, aspect, near, far float64) Matrix {
	ymax := near * math.Tan(fovy*math.Pi/360)
	xmax := ymax * aspect
	return Frustum(-xmax, xmax, -ymax, ymax, near, far)
}

// LookAt 返回一个观察矩阵
func LookAt(eye, center, up Vector) Matrix {
	z := eye.Sub(center).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)
	return Matrix{
		x.X, x.Y, x.Z, -x.Dot(eye), // 第一行
		y.X, y.Y, y.Z, -y.Dot(eye), // 第二行
		z.X, z.Y, z.Z, -z.Dot(eye), // 第三行
		0, 0, 0, 1, // 第四行
	}
}

// LookAtDirection 返回一个观察矩阵
func LookAtDirection(forward, up Vector) Matrix {
	z := forward.Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)
	return Matrix{
		x.X, x.Y, x.Z, 0, // 第一行
		y.X, y.Y, y.Z, 0, // 第二行
		z.X, z.Y, z.Z, 0, // 第三行
		0, 0, 0, 1, // 第四行
	}
}

// Screen 返回一个屏幕矩阵
func Screen(w, h int) Matrix {
	w2 := float64(w) / 2
	h2 := float64(h) / 2
	return Matrix{
		w2, 0, 0, w2, // 第一行
		0, -h2, 0, h2, // 第二行
		0, 0, 0.5, 0.5, // 第三行
		0, 0, 0, 1, // 第四行
	}
}

// Viewport 返回一个视口矩阵
func Viewport(x, y, w, h float64) Matrix {
	l := x
	b := y
	r := x + w
	t := y + h
	return Matrix{
		(r - l) / 2, 0, 0, (r + l) / 2, // 第一行
		0, (t - b) / 2, 0, (t + b) / 2, // 第二行
		0, 0, 0.5, 0.5, // 第三行
		0, 0, 0, 1, // 第四行
	}
}

// Translate 返回一个平移矩阵
func (m Matrix) Translate(v Vector) Matrix {
	return Translate(v).Mul(m)
}

// Scale 返回一个缩放矩阵
func (m Matrix) Scale(v Vector) Matrix {
	return Scale(v).Mul(m)
}

// Rotate 返回一个绕向量 v 旋转 a 角度的矩阵
func (m Matrix) Rotate(v Vector, a float64) Matrix {
	return Rotate(v, a).Mul(m)
}

// RotateTo 返回一个将 a 向量旋转到 b 向量的矩阵
func (m Matrix) RotateTo(a, b Vector) Matrix {
	return RotateTo(a, b).Mul(m)
}

// Frustum 返回一个透视投影矩阵
func (m Matrix) Frustum(l, r, b, t, n, f float64) Matrix {
	return Frustum(l, r, b, t, n, f).Mul(m)
}

// Orthographic 返回一个正交投影矩阵
func (m Matrix) Orthographic(l, r, b, t, n, f float64) Matrix {
	return Orthographic(l, r, b, t, n, f).Mul(m)
}

// Perspective 返回一个透视投影矩阵
func (m Matrix) Perspective(fovy, aspect, near, far float64) Matrix {
	return Perspective(fovy, aspect, near, far).Mul(m)
}

// LookAt 返回一个观察矩阵
func (m Matrix) LookAt(eye, center, up Vector) Matrix {
	return LookAt(eye, center, up).Mul(m)
}

// LookAtDirection 返回一个观察矩阵
func (m Matrix) LookAtDirection(forward, up Vector) Matrix {
	return LookAtDirection(forward, up).Mul(m)
}

// Screen 返回一个屏幕矩阵
func (m Matrix) Screen(w, h int) Matrix {
	return Screen(w, h).Mul(m)
}

// Viewport 返回一个视口矩阵
func (m Matrix) Viewport(x, y, w, h float64) Matrix {
	return Viewport(x, y, w, h).Mul(m)
}

// MulScalar 返回矩阵的标量积
func (a Matrix) MulScalar(b float64) Matrix {
	return Matrix{
		a.X00 * b, a.X01 * b, a.X02 * b, a.X03 * b,
		a.X10 * b, a.X11 * b, a.X12 * b, a.X13 * b,
		a.X20 * b, a.X21 * b, a.X22 * b, a.X23 * b,
		a.X30 * b, a.X31 * b, a.X32 * b, a.X33 * b,
	}
}

// Mul 返回矩阵的乘积
func (a Matrix) Mul(b Matrix) Matrix {
	m := Matrix{}
	m.X00 = a.X00*b.X00 + a.X01*b.X10 + a.X02*b.X20 + a.X03*b.X30
	m.X10 = a.X10*b.X00 + a.X11*b.X10 + a.X12*b.X20 + a.X13*b.X30
	m.X20 = a.X20*b.X00 + a.X21*b.X10 + a.X22*b.X20 + a.X23*b.X30
	m.X30 = a.X30*b.X00 + a.X31*b.X10 + a.X32*b.X20 + a.X33*b.X30
	m.X01 = a.X00*b.X01 + a.X01*b.X11 + a.X02*b.X21 + a.X03*b.X31
	m.X11 = a.X10*b.X01 + a.X11*b.X11 + a.X12*b.X21 + a.X13*b.X31
	m.X21 = a.X20*b.X01 + a.X21*b.X11 + a.X22*b.X21 + a.X23*b.X31
	m.X31 = a.X30*b.X01 + a.X31*b.X11 + a.X32*b.X21 + a.X33*b.X31
	m.X02 = a.X00*b.X02 + a.X01*b.X12 + a.X02*b.X22 + a.X03*b.X32
	m.X12 = a.X10*b.X02 + a.X11*b.X12 + a.X12*b.X22 + a.X13*b.X32
	m.X22 = a.X20*b.X02 + a.X21*b.X12 + a.X22*b.X22 + a.X23*b.X32
	m.X32 = a.X30*b.X02 + a.X31*b.X12 + a.X32*b.X22 + a.X33*b.X32
	m.X03 = a.X00*b.X03 + a.X01*b.X13 + a.X02*b.X23 + a.X03*b.X33
	m.X13 = a.X10*b.X03 + a.X11*b.X13 + a.X12*b.X23 + a.X13*b.X33
	m.X23 = a.X20*b.X03 + a.X21*b.X13 + a.X22*b.X23 + a.X23*b.X33
	m.X33 = a.X30*b.X03 + a.X31*b.X13 + a.X32*b.X23 + a.X33*b.X33
	return m
}

// MulPosition 返回矩阵与位置向量的乘积
func (a Matrix) MulPosition(b Vector) Vector {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z + a.X03
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z + a.X13
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z + a.X23
	return Vector{x, y, z}
}

// MulPositionW 返回矩阵与位置向量的乘积
func (a Matrix) MulPositionW(b Vector) VectorW {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z + a.X03
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z + a.X13
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z + a.X23
	w := a.X30*b.X + a.X31*b.Y + a.X32*b.Z + a.X33
	return VectorW{x, y, z, w}
}

// MulDirection 返回矩阵与方向向量的乘积
func (a Matrix) MulDirection(b Vector) Vector {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z
	return Vector{x, y, z}.Normalize()
}

// MulBox 返回矩阵与包围盒的乘积
func (a Matrix) MulBox(box Box) Box {
	// http://dev.theomader.com/transform-bounding-boxes/
	r := Vector{a.X00, a.X10, a.X20}
	u := Vector{a.X01, a.X11, a.X21}
	b := Vector{a.X02, a.X12, a.X22}
	t := Vector{a.X03, a.X13, a.X23}
	xa := r.MulScalar(box.Min.X)
	xb := r.MulScalar(box.Max.X)
	ya := u.MulScalar(box.Min.Y)
	yb := u.MulScalar(box.Max.Y)
	za := b.MulScalar(box.Min.Z)
	zb := b.MulScalar(box.Max.Z)
	xa, xb = xa.Min(xb), xa.Max(xb)
	ya, yb = ya.Min(yb), ya.Max(yb)
	za, zb = za.Min(zb), za.Max(zb)
	min := xa.Add(ya).Add(za).Add(t)
	max := xb.Add(yb).Add(zb).Add(t)
	return Box{min, max}
}

// Transpose 返回矩阵的转置矩阵
func (a Matrix) Transpose() Matrix {
	return Matrix{
		a.X00, a.X10, a.X20, a.X30,
		a.X01, a.X11, a.X21, a.X31,
		a.X02, a.X12, a.X22, a.X32,
		a.X03, a.X13, a.X23, a.X33}
}

// Determinant 返回矩阵的行列式
func (a Matrix) Determinant() float64 {
	return (a.X00*a.X11*a.X22*a.X33 - a.X00*a.X11*a.X23*a.X32 +
		a.X00*a.X12*a.X23*a.X31 - a.X00*a.X12*a.X21*a.X33 +
		a.X00*a.X13*a.X21*a.X32 - a.X00*a.X13*a.X22*a.X31 -
		a.X01*a.X12*a.X23*a.X30 + a.X01*a.X12*a.X20*a.X33 -
		a.X01*a.X13*a.X20*a.X32 + a.X01*a.X13*a.X22*a.X30 -
		a.X01*a.X10*a.X22*a.X33 + a.X01*a.X10*a.X23*a.X32 +
		a.X02*a.X13*a.X20*a.X31 - a.X02*a.X13*a.X21*a.X30 +
		a.X02*a.X10*a.X21*a.X33 - a.X02*a.X10*a.X23*a.X31 +
		a.X02*a.X11*a.X23*a.X30 - a.X02*a.X11*a.X20*a.X33 -
		a.X03*a.X10*a.X21*a.X32 + a.X03*a.X10*a.X22*a.X31 -
		a.X03*a.X11*a.X22*a.X30 + a.X03*a.X11*a.X20*a.X32 -
		a.X03*a.X12*a.X20*a.X31 + a.X03*a.X12*a.X21*a.X30)
}

// Inverse 返回矩阵的逆矩阵
func (a Matrix) Inverse() Matrix {
	m := Matrix{}
	d := a.Determinant()
	m.X00 = (a.X12*a.X23*a.X31 - a.X13*a.X22*a.X31 + a.X13*a.X21*a.X32 - a.X11*a.X23*a.X32 - a.X12*a.X21*a.X33 + a.X11*a.X22*a.X33) / d
	m.X01 = (a.X03*a.X22*a.X31 - a.X02*a.X23*a.X31 - a.X03*a.X21*a.X32 + a.X01*a.X23*a.X32 + a.X02*a.X21*a.X33 - a.X01*a.X22*a.X33) / d
	m.X02 = (a.X02*a.X13*a.X31 - a.X03*a.X12*a.X31 + a.X03*a.X11*a.X32 - a.X01*a.X13*a.X32 - a.X02*a.X11*a.X33 + a.X01*a.X12*a.X33) / d
	m.X03 = (a.X03*a.X12*a.X21 - a.X02*a.X13*a.X21 - a.X03*a.X11*a.X22 + a.X01*a.X13*a.X22 + a.X02*a.X11*a.X23 - a.X01*a.X12*a.X23) / d
	m.X10 = (a.X13*a.X22*a.X30 - a.X12*a.X23*a.X30 - a.X13*a.X20*a.X32 + a.X10*a.X23*a.X32 + a.X12*a.X20*a.X33 - a.X10*a.X22*a.X33) / d
	m.X11 = (a.X02*a.X23*a.X30 - a.X03*a.X22*a.X30 + a.X03*a.X20*a.X32 - a.X00*a.X23*a.X32 - a.X02*a.X20*a.X33 + a.X00*a.X22*a.X33) / d
	m.X12 = (a.X03*a.X12*a.X30 - a.X02*a.X13*a.X30 - a.X03*a.X10*a.X32 + a.X00*a.X13*a.X32 + a.X02*a.X10*a.X33 - a.X00*a.X12*a.X33) / d
	m.X13 = (a.X02*a.X13*a.X20 - a.X03*a.X12*a.X20 + a.X03*a.X10*a.X22 - a.X00*a.X13*a.X22 - a.X02*a.X10*a.X23 + a.X00*a.X12*a.X23) / d
	m.X20 = (a.X11*a.X23*a.X30 - a.X13*a.X21*a.X30 + a.X13*a.X20*a.X31 - a.X10*a.X23*a.X31 - a.X11*a.X20*a.X33 + a.X10*a.X21*a.X33) / d
	m.X21 = (a.X03*a.X21*a.X30 - a.X01*a.X23*a.X30 - a.X03*a.X20*a.X31 + a.X00*a.X23*a.X31 + a.X01*a.X20*a.X33 - a.X00*a.X21*a.X33) / d
	m.X22 = (a.X01*a.X13*a.X30 - a.X03*a.X11*a.X30 + a.X03*a.X10*a.X31 - a.X00*a.X13*a.X31 - a.X01*a.X10*a.X33 + a.X00*a.X11*a.X33) / d
	m.X23 = (a.X03*a.X11*a.X20 - a.X01*a.X13*a.X20 - a.X03*a.X10*a.X21 + a.X00*a.X13*a.X21 + a.X01*a.X10*a.X23 - a.X00*a.X11*a.X23) / d
	m.X30 = (a.X12*a.X21*a.X30 - a.X11*a.X22*a.X30 - a.X12*a.X20*a.X31 + a.X10*a.X22*a.X31 + a.X11*a.X20*a.X32 - a.X10*a.X21*a.X32) / d
	m.X31 = (a.X01*a.X22*a.X30 - a.X02*a.X21*a.X30 + a.X02*a.X20*a.X31 - a.X00*a.X22*a.X31 - a.X01*a.X20*a.X32 + a.X00*a.X21*a.X32) / d
	m.X32 = (a.X02*a.X11*a.X30 - a.X01*a.X12*a.X30 - a.X02*a.X10*a.X31 + a.X00*a.X12*a.X31 + a.X01*a.X10*a.X32 - a.X00*a.X11*a.X32) / d
	m.X33 = (a.X01*a.X12*a.X20 - a.X02*a.X11*a.X20 + a.X02*a.X10*a.X21 - a.X00*a.X12*a.X21 - a.X01*a.X10*a.X22 + a.X00*a.X11*a.X22) / d
	return m
}
