package fauxgl

// Line 结构体表示一条线段
type Line struct {
	V1, V2 Vertex
}

// NewLine 用于创建一条线段
func NewLine(v1, v2 Vertex) *Line {
	return &Line{v1, v2}
}

// NewLineForPoints 用于创建一条线段
func NewLineForPoints(p1, p2 Vector) *Line {
	v1 := Vertex{Position: p1}
	v2 := Vertex{Position: p2}
	return NewLine(v1, v2)
}

// BoundingBox 用于获取线段的包围盒
func (l *Line) BoundingBox() Box {
	min := l.V1.Position.Min(l.V2.Position)
	max := l.V1.Position.Max(l.V2.Position)
	return Box{min, max}
}

// Transform 用于对线段进行变换
func (l *Line) Transform(matrix Matrix) {
	l.V1.Position = matrix.MulPosition(l.V1.Position)
	l.V2.Position = matrix.MulPosition(l.V2.Position)
	l.V1.Normal = matrix.MulDirection(l.V1.Normal)
	l.V2.Normal = matrix.MulDirection(l.V2.Normal)
}
