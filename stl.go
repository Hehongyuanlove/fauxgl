package fauxgl

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
)

// STLHeader 是 STL 文件头部
type STLHeader struct {
	_     [80]uint8 // 保留80字节
	Count uint32    // 三角形数量
}

// STLTriangle 是 STL 文件中的三角形
type STLTriangle struct {
	N, V1, V2, V3 [3]float32 // 法向量和三角形三个顶点
	_             uint16     // 保留2字节
}

// LoadSTL 从 STL 文件中加载网格
func LoadSTL(path string) (*Mesh, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()

	// 读取头部，获取预期的二进制大小
	header := STLHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	expectedSize := int64(header.Count)*50 + 84

	// 倒回文件开头
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// 解析 ASCII 或二进制 STL
	if size == expectedSize {
		return loadSTLB(file)
	} else {
		return loadSTLA(file)
	}
}

// 从 ASCII STL 文件中加载网格
func loadSTLA(file *os.File) (*Mesh, error) {
	var vertexes []Vector // 顶点数组
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			f := ParseFloats(fields[1:])
			vertexes = append(vertexes, Vector{f[0], f[1], f[2]})
		}
	}
	var triangles []*Triangle // 三角形数组
	for i := 0; i < len(vertexes); i += 3 {
		t := Triangle{}
		t.V1.Position = vertexes[i+0]
		t.V2.Position = vertexes[i+1]
		t.V3.Position = vertexes[i+2]
		t.FixNormals()
		triangles = append(triangles, &t)
	}
	return NewTriangleMesh(triangles), scanner.Err()
}

func makeFloat(b []byte) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(b)))
}

// 从二进制 STL 文件中加载网格
func loadSTLB(file *os.File) (*Mesh, error) {
	r := bufio.NewReader(file)
	header := STLHeader{}
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	count := int(header.Count)
	triangles := make([]*Triangle, count)
	_triangles := make([]Triangle, count)
	b := make([]byte, count*50)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	wn := runtime.NumCPU()
	ch := make(chan Box, wn)
	for wi := 0; wi < wn; wi++ {
		go func(wi, wn int) {
			var min, max Vector
			for i := wi; i < count; i += wn {
				j := i * 50
				v1 := Vector{makeFloat(b[j+12 : j+16]), makeFloat(b[j+16 : j+20]), makeFloat(b[j+20 : j+24])}
				v2 := Vector{makeFloat(b[j+24 : j+28]), makeFloat(b[j+28 : j+32]), makeFloat(b[j+32 : j+36])}
				v3 := Vector{makeFloat(b[j+36 : j+40]), makeFloat(b[j+40 : j+44]), makeFloat(b[j+44 : j+48])}
				t := &_triangles[i]
				t.V1.Position = v1
				t.V2.Position = v2
				t.V3.Position = v3
				n := t.Normal()
				t.V1.Normal = n
				t.V2.Normal = n
				t.V3.Normal = n
				if i == wi {
					min = v1
					max = v1
				}
				for _, v := range [3]Vector{v1, v2, v3} {
					if v.X < min.X {
						min.X = v.X
					}
					if v.Y < min.Y {
						min.Y = v.Y
					}
					if v.Z < min.Z {
						min.Z = v.Z
					}
					if v.X > max.X {
						max.X = v.X
					}
					if v.Y > max.Y {
						max.Y = v.Y
					}
					if v.Z > max.Z {
						max.Z = v.Z
					}
				}
				triangles[i] = t
			}
			ch <- Box{min, max}
		}(wi, wn)
	}
	box := EmptyBox
	for wi := 0; wi < wn; wi++ {
		box = box.Extend(<-ch)
	}
	mesh := NewTriangleMesh(triangles)
	mesh.box = &box
	return mesh, nil
}

// SaveSTL 将网格保存到 STL 文件中
func SaveSTL(path string, mesh *Mesh) error {
	// 创建文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建写入器
	w := bufio.NewWriter(file)

	// 写入头部
	header := STLHeader{}
	header.Count = uint32(len(mesh.Triangles))
	if err := binary.Write(w, binary.LittleEndian, &header); err != nil {
		return err
	}

	// 写入三角形
	for _, triangle := range mesh.Triangles {
		// 计算法向量
		n := triangle.Normal()

		// 创建 STL 三角形
		d := STLTriangle{}
		d.N[0] = float32(n.X)
		d.N[1] = float32(n.Y)
		d.N[2] = float32(n.Z)
		d.V1[0] = float32(triangle.V1.Position.X)
		d.V1[1] = float32(triangle.V1.Position.Y)
		d.V1[2] = float32(triangle.V1.Position.Z)
		d.V2[0] = float32(triangle.V2.Position.X)
		d.V2[1] = float32(triangle.V2.Position.Y)
		d.V2[2] = float32(triangle.V2.Position.Z)
		d.V3[0] = float32(triangle.V3.Position.X)
		d.V3[1] = float32(triangle.V3.Position.Y)
		d.V3[2] = float32(triangle.V3.Position.Z)

		// 写入 STL 三角形
		if err := binary.Write(w, binary.LittleEndian, &d); err != nil {
			return err
		}
	}

	// 刷新缓冲区
	w.Flush()

	return nil
}
