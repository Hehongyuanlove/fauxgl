package fauxgl

import "math"

var EmptyBox = Box{}

type Box struct {
	Min, Max Vector // 最小点和最大点
}


// BoxForBoxes 返回一个包含所有盒子的最小盒子
func BoxForBoxes(boxes []Box) Box {
	if len(boxes) == 0 {
		return EmptyBox
	}
	x0 := boxes[0].Min.X
	y0 := boxes[0].Min.Y
	z0 := boxes[0].Min.Z
	x1 := boxes[0].Max.X
	y1 := boxes[0].Max.Y
	z1 := boxes[0].Max.Z
	for _, box := range boxes {
		x0 = math.Min(x0, box.Min.X)
		y0 = math.Min(y0, box.Min.Y)
		z0 = math.Min(z0, box.Min.Z)
		x1 = math.Max(x1, box.Max.X)
		y1 = math.Max(y1, box.Max.Y)
		z1 = math.Max(z1, box.Max.Z)
	}
	return Box{Vector{x0, y0, z0}, Vector{x1, y1, z1}}
}

// Volume 返回盒子的体积
func (a Box) Volume() float64 {
	s := a.Size()
	return s.X * s.Y * s.Z
}

// Anchor 返回盒子的锚点
func (a Box) Anchor(anchor Vector) Vector {
	return a.Min.Add(a.Size().Mul(anchor))
}

// Center 返回盒子的中心点
func (a Box) Center() Vector {
	return a.Anchor(Vector{0.5, 0.5, 0.5})
}

// Size 返回盒子的大小
func (a Box) Size() Vector {
	return a.Max.Sub(a.Min)
}

// Extend 扩展盒子
func (a Box) Extend(b Box) Box {
	if a == EmptyBox {
		return b
	}
	return Box{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// Offset 偏移盒子
func (a Box) Offset(x float64) Box {
	return Box{a.Min.SubScalar(x), a.Max.AddScalar(x)}
}

// Translate 移动盒子
func (a Box) Translate(v Vector) Box {
	return Box{a.Min.Add(v), a.Max.Add(v)}
}

// Contains 判断点是否在盒子内
func (a Box) Contains(b Vector) bool {
	return a.Min.X <= b.X && a.Max.X >= b.X &&
		a.Min.Y <= b.Y && a.Max.Y >= b.Y &&
		a.Min.Z <= b.Z && a.Max.Z >= b.Z
}

// ContainsBox 判断盒子是否在盒子内
func (a Box) ContainsBox(b Box) bool {
	return a.Min.X <= b.Min.X && a.Max.X >= b.Max.X &&
		a.Min.Y <= b.Min.Y && a.Max.Y >= b.Max.Y &&
		a.Min.Z <= b.Min.Z && a.Max.Z >= b.Max.Z
}

// Intersects 判断盒子是否相交
func (a Box) Intersects(b Box) bool {
	return !(a.Min.X > b.Max.X || a.Max.X < b.Min.X || a.Min.Y > b.Max.Y ||
		a.Max.Y < b.Min.Y || a.Min.Z > b.Max.Z || a.Max.Z < b.Min.Z)
}

// Intersection 返回两个盒子的交集
func (a Box) Intersection(b Box) Box {
	if !a.Intersects(b) {
		return EmptyBox
	}
	min := a.Min.Max(b.Min)
	max := a.Max.Min(b.Max)
	min, max = min.Min(max), min.Max(max)
	return Box{min, max}
}

// Transform 变换盒子
func (a Box) Transform(m Matrix) Box {
	return m.MulBox(a)
}
