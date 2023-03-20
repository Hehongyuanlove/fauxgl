package fauxgl

import "math"

// NewPlane函数返回一个平面网格
func NewPlane() *Mesh {
	v1 := Vector{-0.5, -0.5, 0}
	v2 := Vector{0.5, -0.5, 0}
	v3 := Vector{0.5, 0.5, 0}
	v4 := Vector{-0.5, 0.5, 0}
	return NewTriangleMesh([]*Triangle{
		NewTriangleForPoints(v1, v2, v3),
		NewTriangleForPoints(v1, v3, v4),
	})
}

// NewCube函数返回一个立方体网格
func NewCube() *Mesh {
	v := []Vector{
		{-1, -1, -1}, {-1, -1, 1}, {-1, 1, -1}, {-1, 1, 1},
		{1, -1, -1}, {1, -1, 1}, {1, 1, -1}, {1, 1, 1},
	}
	// 创建一个新的三角形网格
	mesh := NewTriangleMesh([]*Triangle{
		NewTriangleForPoints(v[3], v[5], v[7]), // 创建一个新的三角形
		NewTriangleForPoints(v[5], v[3], v[1]),
		NewTriangleForPoints(v[0], v[6], v[4]),
		NewTriangleForPoints(v[6], v[0], v[2]),
		NewTriangleForPoints(v[0], v[5], v[1]),
		NewTriangleForPoints(v[5], v[0], v[4]),
		NewTriangleForPoints(v[5], v[6], v[7]),
		NewTriangleForPoints(v[6], v[5], v[4]),
		NewTriangleForPoints(v[6], v[3], v[7]),
		NewTriangleForPoints(v[3], v[6], v[2]),
		NewTriangleForPoints(v[0], v[3], v[2]),
		NewTriangleForPoints(v[3], v[0], v[1]),
	})
	// 将网格缩放为0.5
	mesh.Transform(Scale(Vector{0.5, 0.5, 0.5}))
	return mesh
}

// NewCubeForBox函数返回一个盒子网格
func NewCubeForBox(box Box) *Mesh {
	m := Translate(Vector{0.5, 0.5, 0.5})
	m = m.Scale(box.Size())
	m = m.Translate(box.Min)
	cube := NewCube()
	cube.Transform(m)
	return cube
}

// NewCubeOutlineForBox函数返回一个盒子的轮廓线网格
func NewCubeOutlineForBox(box Box) *Mesh {
	x0 := box.Min.X
	y0 := box.Min.Y
	z0 := box.Min.Z
	x1 := box.Max.X
	y1 := box.Max.Y
	z1 := box.Max.Z
	return NewLineMesh([]*Line{
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y0, z1}), // 左下角前后两个点
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x0, y1, z1}), // 左上角前后两个点
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y0, z1}), // 右下角前后两个点
		NewLineForPoints(Vector{x1, y1, z0}, Vector{x1, y1, z1}), // 右上角前后两个点
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x0, y1, z0}), // 左下角上下两个点
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x0, y1, z1}), // 左上角上下两个点
		NewLineForPoints(Vector{x1, y0, z0}, Vector{x1, y1, z0}), // 右下角上下两个点
		NewLineForPoints(Vector{x1, y0, z1}, Vector{x1, y1, z1}), // 右上角上下两个点
		NewLineForPoints(Vector{x0, y0, z0}, Vector{x1, y0, z0}), // 左下角左右两个点
		NewLineForPoints(Vector{x0, y1, z0}, Vector{x1, y1, z0}), // 左上角左右两个点
		NewLineForPoints(Vector{x0, y0, z1}, Vector{x1, y0, z1}), // 右下角左右两个点
		NewLineForPoints(Vector{x0, y1, z1}, Vector{x1, y1, z1}), // 右上角左右两个点
	})
}

// NewLatLngSphere函数返回一个经纬度球网格
func NewLatLngSphere(latStep, lngStep int) *Mesh {
	var triangles []*Triangle
	for lat0 := -90; lat0 < 90; lat0 += latStep {
		lat1 := lat0 + latStep
		v0 := float64(lat0+90) / 180
		v1 := float64(lat1+90) / 180
		for lng0 := -180; lng0 < 180; lng0 += lngStep {
			lng1 := lng0 + lngStep
			u0 := float64(lng0+180) / 360
			u1 := float64(lng1+180) / 360
			if lng1 >= 180 {
				lng1 -= 360
			}
			p00 := LatLngToXYZ(float64(lat0), float64(lng0)) // 经纬度转换为三维坐标
			p01 := LatLngToXYZ(float64(lat0), float64(lng1))
			p10 := LatLngToXYZ(float64(lat1), float64(lng0))
			p11 := LatLngToXYZ(float64(lat1), float64(lng1))
			t1 := NewTriangleForPoints(p00, p01, p11)
			t2 := NewTriangleForPoints(p00, p11, p10)
			if lat0 != -90 {
				t1.V1.Texture = Vector{u0, v0, 0} // 设置纹理坐标
				t1.V2.Texture = Vector{u1, v0, 0}
				t1.V3.Texture = Vector{u1, v1, 0}
				triangles = append(triangles, t1)
			}
			if lat1 != 90 {
				t2.V1.Texture = Vector{u0, v0, 0}
				t2.V2.Texture = Vector{u1, v1, 0}
				t2.V3.Texture = Vector{u0, v1, 0}
				triangles = append(triangles, t2)
			}
		}
	}
	return NewTriangleMesh(triangles)
}

// NewSphere函数返回一个球体网格
func NewSphere(detail int) *Mesh {
	var triangles []*Triangle
	ico := NewIcosahedron() // 创建一个新的正二十面体
	for _, t := range ico.Triangles {
		v1 := t.V1.Position // 获取三角形的三个顶点
		v2 := t.V2.Position
		v3 := t.V3.Position
		triangles = append(triangles, newSphereHelper(detail, v1, v2, v3)...) // 将三角形分解为更小的三角形
	}
	return NewTriangleMesh(triangles) // 返回新的三角形网格
}

// newSphereHelper函数将三角形分解为更小的三角形
func newSphereHelper(detail int, v1, v2, v3 Vector) []*Triangle {
	if detail == 0 {
		t := NewTriangleForPoints(v1, v2, v3)
		return []*Triangle{t}
	}
	var triangles []*Triangle
	v12 := v1.Add(v2).DivScalar(2).Normalize() // 获取三角形的中点
	v13 := v1.Add(v3).DivScalar(2).Normalize()
	v23 := v2.Add(v3).DivScalar(2).Normalize()
	triangles = append(triangles, newSphereHelper(detail-1, v1, v12, v13)...) // 递归分解三角形
	triangles = append(triangles, newSphereHelper(detail-1, v2, v23, v12)...)
	triangles = append(triangles, newSphereHelper(detail-1, v3, v13, v23)...)
	triangles = append(triangles, newSphereHelper(detail-1, v12, v23, v13)...)
	return triangles
}

// NewCylinder函数返回一个圆柱网格
func NewCylinder(step int, capped bool) *Mesh {
	var triangles []*Triangle
	for a0 := 0; a0 < 360; a0 += step {
		a1 := (a0 + step) % 360
		r0 := Radians(float64(a0))
		r1 := Radians(float64(a1))
		x0 := math.Cos(r0)
		y0 := math.Sin(r0)
		x1 := math.Cos(r1)
		y1 := math.Sin(r1)
		p00 := Vector{x0, y0, -0.5}               // 底部左侧点
		p10 := Vector{x1, y1, -0.5}               // 底部右侧点
		p11 := Vector{x1, y1, 0.5}                // 顶部右侧点
		p01 := Vector{x0, y0, 0.5}                // 顶部左侧点
		t1 := NewTriangleForPoints(p00, p10, p11) // 底部三角形
		t2 := NewTriangleForPoints(p00, p11, p01) // 顶部三角形
		triangles = append(triangles, t1)
		triangles = append(triangles, t2)
		if capped {
			p0 := Vector{0, 0, -0.5}                 // 底部中心点
			p1 := Vector{0, 0, 0.5}                  // 顶部中心点
			t3 := NewTriangleForPoints(p0, p10, p00) // 底部三角形
			t4 := NewTriangleForPoints(p1, p01, p11) // 顶部三角形
			triangles = append(triangles, t3)
			triangles = append(triangles, t4)
		}
	}
	return NewTriangleMesh(triangles)
}

// NewCone函数返回一个圆锥网格
func NewCone(step int, capped bool) *Mesh {
	var triangles []*Triangle
	for a0 := 0; a0 < 360; a0 += step {
		a1 := (a0 + step) % 360
		r0 := Radians(float64(a0))
		r1 := Radians(float64(a1))
		x0 := math.Cos(r0)
		y0 := math.Sin(r0)
		x1 := math.Cos(r1)
		y1 := math.Sin(r1)
		p00 := Vector{x0, y0, -0.5}              // 底部左侧点
		p10 := Vector{x1, y1, -0.5}              // 底部右侧点
		p1 := Vector{0, 0, 0.5}                  // 顶部中心点
		t1 := NewTriangleForPoints(p00, p10, p1) // 底部三角形
		triangles = append(triangles, t1)
		if capped {
			p0 := Vector{0, 0, -0.5}                 // 底部中心点
			t2 := NewTriangleForPoints(p0, p10, p00) // 顶部三角形
			triangles = append(triangles, t2)
		}
	}
	return NewTriangleMesh(triangles) // 返回新的三角形网格
}

// NewIcosahedron函数返回一个正二十面体网格
func NewIcosahedron() *Mesh {
	const a = 0.8506507174597755 // 五边形的内角余弦值
	const b = 0.5257312591858783 // 五边形的外角正弦值
	vertices := []Vector{
		{-a, -b, 0}, // 顶点0
		{-a, b, 0},  // 顶点1
		{-b, 0, -a}, // 顶点2
		{-b, 0, a},  // 顶点3
		{0, -a, -b}, // 顶点4
		{0, -a, b},  // 顶点5
		{0, a, -b},  // 顶点6
		{0, a, b},   // 顶点7
		{b, 0, -a},  // 顶点8
		{b, 0, a},   // 顶点9
		{a, -b, 0},  // 顶点10
		{a, b, 0},   // 顶点11
	}
	indices := [][3]int{
		{0, 3, 1},   // 三角形0
		{1, 3, 7},   // 三角形1
		{2, 0, 1},   // 三角形2
		{2, 1, 6},   // 三角形3
		{4, 0, 2},   // 三角形4
		{4, 5, 0},   // 三角形5
		{5, 3, 0},   // 三角形6
		{6, 1, 7},   // 三角形7
		{6, 7, 11},  // 三角形8
		{7, 3, 9},   // 三角形9
		{8, 2, 6},   // 三角形10
		{8, 4, 2},   // 三角形11
		{8, 6, 11},  // 三角形12
		{8, 10, 4},  // 三角形13
		{8, 11, 10}, // 三角形14
		{9, 3, 5},   // 三角形15
		{10, 5, 4},  // 三角形16
		{10, 9, 5},  // 三角形17
		{11, 7, 9},  // 三角形18
		{11, 9, 10}, // 三角形19
	}
	triangles := make([]*Triangle, len(indices))
	for i, idx := range indices {
		p1 := vertices[idx[0]]
		p2 := vertices[idx[1]]
		p3 := vertices[idx[2]]
		triangles[i] = NewTriangleForPoints(p1, p2, p3) // 创建新的三角形
	}
	return NewTriangleMesh(triangles) // 返回新的三角形网格
}
