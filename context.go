package fauxgl

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

// Face 表示三角形的正反面
type Face int

const (
	_ Face = iota
	// FaceCW 表示三角形的正面为顺时针方向
	FaceCW
	// FaceCCW 表示三角形的正面为逆时针方向
	FaceCCW
)

// Cull 表示剔除模式
type Cull int

const (
	_ Cull = iota
	// CullNone 表示不剔除
	CullNone
	// CullFront 表示剔除正面
	CullFront
	// CullBack 表示剔除背面
	CullBack
)

// RasterizeInfo 表示光栅化信息
type RasterizeInfo struct {
	// TotalPixels 表示总像素数
	TotalPixels uint64
	// UpdatedPixels 表示更新的像素数
	UpdatedPixels uint64
}

// Add 将两个 RasterizeInfo 相加
func (info RasterizeInfo) Add(other RasterizeInfo) RasterizeInfo {
	return RasterizeInfo{
		info.TotalPixels + other.TotalPixels,
		info.UpdatedPixels + other.UpdatedPixels,
	}
}

// Context 是一个渲染上下文，包含颜色缓冲区、深度缓冲区、清除颜色、着色器、深度测试、颜色混合、线框模式、剔除模式、线宽、深度偏移、屏幕矩阵和锁等属性
type Context struct {
	Width        int          // 宽度
	Height       int          // 高度
	ColorBuffer  *image.NRGBA // 颜色缓冲区
	DepthBuffer  []float64    // 深度缓冲区
	ClearColor   Color        // 清除颜色
	Shader       Shader       // 着色器
	ReadDepth    bool         // 深度测试
	WriteDepth   bool         // 深度测试
	WriteColor   bool         // 颜色混合
	AlphaBlend   bool         // 颜色混合
	Wireframe    bool         // 线框模式
	FrontFace    Face         // 剔除模式
	Cull         Cull         // 剔除模式
	LineWidth    float64      // 线宽
	DepthBias    float64      // 深度偏移
	screenMatrix Matrix       // 屏幕矩阵
	locks        []sync.Mutex // 锁
}

// NewContext 创建一个新的渲染上下文
func NewContext(width, height int) *Context {
	dc := &Context{}
	dc.Width = width
	dc.Height = height
	dc.ColorBuffer = image.NewNRGBA(image.Rect(0, 0, width, height))
	dc.DepthBuffer = make([]float64, width*height)
	dc.ClearColor = Transparent
	dc.Shader = NewSolidColorShader(Identity(), Color{1, 0, 1, 1})
	dc.ReadDepth = true
	dc.WriteDepth = true
	dc.WriteColor = true
	dc.AlphaBlend = true
	dc.Wireframe = false
	dc.FrontFace = FaceCCW
	dc.Cull = CullBack
	dc.LineWidth = 2
	dc.DepthBias = 0
	dc.screenMatrix = Screen(width, height)
	dc.locks = make([]sync.Mutex, 256)
	dc.ClearDepthBuffer()
	return dc
}

// Image 返回颜色缓冲区的图像
func (dc *Context) Image() image.Image {
	return dc.ColorBuffer
}

// DepthImage 返回深度缓冲区的图像
func (dc *Context) DepthImage() image.Image {
	lo := math.MaxFloat64
	hi := -math.MaxFloat64
	for _, d := range dc.DepthBuffer {
		if d == math.MaxFloat64 {
			continue
		}
		if d < lo {
			lo = d
		}
		if d > hi {
			hi = d
		}
	}

	im := image.NewGray16(image.Rect(0, 0, dc.Width, dc.Height))
	var i int
	for y := 0; y < dc.Height; y++ {
		for x := 0; x < dc.Width; x++ {
			d := dc.DepthBuffer[i]
			t := (d - lo) / (hi - lo)
			if d == math.MaxFloat64 {
				t = 1
			}
			c := color.Gray16{uint16(t * 0xffff)}
			im.SetGray16(x, y, c)
			i++
		}
	}
	return im
}

// ClearColorBufferWith 使用指定颜色清除颜色缓冲区
func (dc *Context) ClearColorBufferWith(color Color) {
	c := color.NRGBA()
	for y := 0; y < dc.Height; y++ {
		i := dc.ColorBuffer.PixOffset(0, y)
		for x := 0; x < dc.Width; x++ {
			dc.ColorBuffer.Pix[i+0] = c.R
			dc.ColorBuffer.Pix[i+1] = c.G
			dc.ColorBuffer.Pix[i+2] = c.B
			dc.ColorBuffer.Pix[i+3] = c.A
			i += 4
		}
	}
}

// ClearColorBuffer 使用清除颜色清除颜色缓冲区
func (dc *Context) ClearColorBuffer() {
	dc.ClearColorBufferWith(dc.ClearColor)
}

// ClearDepthBufferWith 使用指定值清除深度缓冲区
func (dc *Context) ClearDepthBufferWith(value float64) {
	for i := range dc.DepthBuffer {
		dc.DepthBuffer[i] = value
	}
}

// ClearDepthBuffer 使用最大值清除深度缓冲区
func (dc *Context) ClearDepthBuffer() {
	dc.ClearDepthBufferWith(math.MaxFloat64)
}

// edge 计算三角形的边
func edge(a, b, c Vector) float64 {
	return (b.X-c.X)*(a.Y-c.Y) - (b.Y-c.Y)*(a.X-c.X)
}

// rasterize 光栅化三角形
func (dc *Context) rasterize(v0, v1, v2 Vertex, s0, s1, s2 Vector) RasterizeInfo {
	var info RasterizeInfo

	// 整数边界框
	min := s0.Min(s1.Min(s2)).Floor()
	max := s0.Max(s1.Max(s2)).Ceil()
	x0 := int(min.X)
	x1 := int(max.X)
	y0 := int(min.Y)
	y1 := int(max.Y)

	// 前向差分变量
	p := Vector{float64(x0) + 0.5, float64(y0) + 0.5, 0}
	w00 := edge(s1, s2, p)
	w01 := edge(s2, s0, p)
	w02 := edge(s0, s1, p)
	a01 := s1.Y - s0.Y
	b01 := s0.X - s1.X
	a12 := s2.Y - s1.Y
	b12 := s1.X - s2.X
	a20 := s0.Y - s2.Y
	b20 := s2.X - s0.X

	// 倒数
	ra := 1 / edge(s0, s1, s2)
	r0 := 1 / v0.Output.W
	r1 := 1 / v1.Output.W
	r2 := 1 / v2.Output.W
	ra12 := 1 / a12
	ra20 := 1 / a20
	ra01 := 1 / a01

	// 遍历边界框中的所有像素
	for y := y0; y <= y1; y++ {
		var d float64
		d0 := -w00 * ra12
		d1 := -w01 * ra20
		d2 := -w02 * ra01
		if w00 < 0 && d0 > d {
			d = d0
		}
		if w01 < 0 && d1 > d {
			d = d1
		}
		if w02 < 0 && d2 > d {
			d = d2
		}
		d = float64(int(d))
		if d < 0 {
			// 在病态情况下发生
			d = 0
		}
		w0 := w00 + a12*d
		w1 := w01 + a20*d
		w2 := w02 + a01*d
		wasInside := false
		for x := x0 + int(d); x <= x1; x++ {
			b0 := w0 * ra
			b1 := w1 * ra
			b2 := w2 * ra
			w0 += a12
			w1 += a20
			w2 += a01
			// 检查是否在三角形内部
			if b0 < 0 || b1 < 0 || b2 < 0 {
				if wasInside {
					break
				}
				continue
			}
			wasInside = true
			// 检查深度缓冲区以进行早期中止
			i := y*dc.Width + x
			if i < 0 || i >= len(dc.DepthBuffer) {
				// TODO: 裁剪舍入误差；修复
				// TODO: 也可能是由于粗线超出屏幕
				continue
			}
			info.TotalPixels++
			z := b0*s0.Z + b1*s1.Z + b2*s2.Z
			bz := z + dc.DepthBias
			if dc.ReadDepth && bz > dc.DepthBuffer[i] { // safe w/out lock?
				continue
			}
			// 透视校正插值顶点数据
			b := VectorW{b0 * r0, b1 * r1, b2 * r2, 0}
			b.W = 1 / (b.X + b.Y + b.Z)
			v := InterpolateVertexes(v0, v1, v2, b)
			// 调用片段着色器
			color := dc.Shader.Fragment(v)
			if color == Discard {
				continue
			}
			// 原子更新缓冲区
			lock := &dc.locks[(x+y)&255]
			lock.Lock()
			// 再次检查深度缓冲区
			if bz <= dc.DepthBuffer[i] || !dc.ReadDepth {
				info.UpdatedPixels++
				if dc.WriteDepth {
					// 更新深度缓冲区
					dc.DepthBuffer[i] = z
				}
				if dc.WriteColor {
					// 更新颜色缓冲区
					if dc.AlphaBlend && color.A < 1 {
						sr, sg, sb, sa := color.NRGBA().RGBA()
						a := (0xffff - sa) * 0x101
						j := dc.ColorBuffer.PixOffset(x, y)
						dr := &dc.ColorBuffer.Pix[j+0]
						dg := &dc.ColorBuffer.Pix[j+1]
						db := &dc.ColorBuffer.Pix[j+2]
						da := &dc.ColorBuffer.Pix[j+3]
						*dr = uint8((uint32(*dr)*a/0xffff + sr) >> 8)
						*dg = uint8((uint32(*dg)*a/0xffff + sg) >> 8)
						*db = uint8((uint32(*db)*a/0xffff + sb) >> 8)
						*da = uint8((uint32(*da)*a/0xffff + sa) >> 8)
					} else {
						dc.ColorBuffer.SetNRGBA(x, y, color.NRGBA())
					}
				}
			}
			lock.Unlock()
		}
		w00 += b12
		w01 += b20
		w02 += b01
	}

	return info
}

// line 绘制线段
func (dc *Context) line(v0, v1 Vertex, s0, s1 Vector) RasterizeInfo {
	n := s1.Sub(s0).Perpendicular().MulScalar(dc.LineWidth / 2)
	s0 = s0.Add(s0.Sub(s1).Normalize().MulScalar(dc.LineWidth / 2))
	s1 = s1.Add(s1.Sub(s0).Normalize().MulScalar(dc.LineWidth / 2))
	s00 := s0.Add(n)
	s01 := s0.Sub(n)
	s10 := s1.Add(n)
	s11 := s1.Sub(n)
	info1 := dc.rasterize(v1, v0, v0, s11, s01, s00)
	info2 := dc.rasterize(v1, v1, v0, s10, s11, s00)
	return info1.Add(info2)
}

// wireframe 绘制线框
func (dc *Context) wireframe(v0, v1, v2 Vertex, s0, s1, s2 Vector) RasterizeInfo {
	info1 := dc.line(v0, v1, s0, s1)
	info2 := dc.line(v1, v2, s1, s2)
	info3 := dc.line(v2, v0, s2, s0)
	return info1.Add(info2).Add(info3)
}

// drawClippedLine 绘制裁剪后的线段
func (dc *Context) drawClippedLine(v0, v1 Vertex) RasterizeInfo {
	// 规范化设备坐标
	ndc0 := v0.Output.DivScalar(v0.Output.W).Vector()
	ndc1 := v1.Output.DivScalar(v1.Output.W).Vector()

	// 屏幕坐标
	s0 := dc.screenMatrix.MulPosition(ndc0)
	s1 := dc.screenMatrix.MulPosition(ndc1)

	// 光栅化
	return dc.line(v0, v1, s0, s1)
}

// drawClippedTriangle 绘制裁剪后的三角形
func (dc *Context) drawClippedTriangle(v0, v1, v2 Vertex) RasterizeInfo {
	// 规范化设备坐标
	ndc0 := v0.Output.DivScalar(v0.Output.W).Vector()
	ndc1 := v1.Output.DivScalar(v1.Output.W).Vector()
	ndc2 := v2.Output.DivScalar(v2.Output.W).Vector()

	// 背面剔除
	a := (ndc1.X-ndc0.X)*(ndc2.Y-ndc0.Y) - (ndc2.X-ndc0.X)*(ndc1.Y-ndc0.Y)
	if a < 0 {
		v0, v1, v2 = v2, v1, v0
		ndc0, ndc1, ndc2 = ndc2, ndc1, ndc0
	}
	if dc.Cull == CullFront {
		a = -a
	}
	if dc.FrontFace == FaceCW {
		a = -a
	}
	if dc.Cull != CullNone && a <= 0 {
		return RasterizeInfo{}
	}

	// 屏幕坐标
	s0 := dc.screenMatrix.MulPosition(ndc0)
	s1 := dc.screenMatrix.MulPosition(ndc1)
	s2 := dc.screenMatrix.MulPosition(ndc2)

	// 光栅化
	if dc.Wireframe {
		return dc.wireframe(v0, v1, v2, s0, s1, s2)
	} else {
		return dc.rasterize(v0, v1, v2, s0, s1, s2)
	}
}

// DrawLine 绘制线段
func (dc *Context) DrawLine(t *Line) RasterizeInfo {
	// 调用顶点着色器
	v1 := dc.Shader.Vertex(t.V1)
	v2 := dc.Shader.Vertex(t.V2)

	if v1.Outside() || v2.Outside() {
		// 裁剪到视图体积
		line := ClipLine(NewLine(v1, v2))
		if line != nil {
			return dc.drawClippedLine(line.V1, line.V2)
		} else {
			return RasterizeInfo{}
		}
	} else {
		// 无需裁剪
		return dc.drawClippedLine(v1, v2)
	}
}

// DrawTriangle 绘制三角形
func (dc *Context) DrawTriangle(t *Triangle) RasterizeInfo {
	// 调用顶点着色器
	v1 := dc.Shader.Vertex(t.V1)
	v2 := dc.Shader.Vertex(t.V2)
	v3 := dc.Shader.Vertex(t.V3)

	if v1.Outside() || v2.Outside() || v3.Outside() {
		// 裁剪到视图体积
		triangles := ClipTriangle(NewTriangle(v1, v2, v3))
		var result RasterizeInfo
		for _, t := range triangles {
			info := dc.drawClippedTriangle(t.V1, t.V2, t.V3)
			result = result.Add(info)
		}
		return result
	} else {
		// 无需裁剪
		return dc.drawClippedTriangle(v1, v2, v3)
	}
}

// DrawLines 绘制线段集合
func (dc *Context) DrawLines(lines []*Line) RasterizeInfo {
	wn := runtime.NumCPU()
	ch := make(chan RasterizeInfo, wn)
	for wi := 0; wi < wn; wi++ {
		go func(wi int) {
			var result RasterizeInfo
			for i, l := range lines {
				if i%wn == wi {
					info := dc.DrawLine(l)
					result = result.Add(info)
				}
			}
			ch <- result
		}(wi)
	}
	var result RasterizeInfo
	for wi := 0; wi < wn; wi++ {
		result = result.Add(<-ch)
	}
	return result
}

// DrawTriangles 绘制三角形集合
func (dc *Context) DrawTriangles(triangles []*Triangle) RasterizeInfo {
	wn := runtime.NumCPU()
	ch := make(chan RasterizeInfo, wn)
	for wi := 0; wi < wn; wi++ {
		go func(wi int) {
			var result RasterizeInfo
			for i, t := range triangles {
				if i%wn == wi {
					info := dc.DrawTriangle(t)
					result = result.Add(info)
				}
			}
			ch <- result
		}(wi)
	}
	var result RasterizeInfo
	for wi := 0; wi < wn; wi++ {
		result = result.Add(<-ch)
	}
	return result
}

// DrawMesh 绘制网格
func (dc *Context) DrawMesh(mesh *Mesh) RasterizeInfo {
	info1 := dc.DrawTriangles(mesh.Triangles)
	info2 := dc.DrawLines(mesh.Lines)
	return info1.Add(info2)
}
