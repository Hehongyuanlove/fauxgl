package fauxgl

import (
	"encoding/binary"
	"io"
	"os"
)

// Load3DS 从3DS文件中加载网格
func Load3DS(filename string) (*Mesh, error) {
	type ChunkHeader struct {
		ChunkID uint16
		Length  uint32
	}

	file, err := os.Open(filename) // 打开文件
	if err != nil {
		return nil, err
	}
	defer file.Close() // 关闭文件

	var vertices []Vector // 顶点列表
	var faces []*Triangle // 面列表
	var triangles []*Triangle // 三角形列表
	for {
		header := ChunkHeader{}
		if err := binary.Read(file, binary.LittleEndian, &header); err != nil { // 读取块头
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch header.ChunkID {
		case 0x4D4D: // 主块
		case 0x3D3D: // 3D编辑器块
		case 0x4000: // 对象块
			_, err := readNullTerminatedString(file) // 读取对象名称
			if err != nil {
				return nil, err
			}
		case 0x4100: // 三角形对象块
		case 0x4110: // 顶点列表
			v, err := readVertexList(file) // 读取顶点列表
			if err != nil {
				return nil, err
			}
			vertices = v
		case 0x4120: // 面列表
			f, err := readFaceList(file, vertices) // 读取面列表
			if err != nil {
				return nil, err
			}
			faces = f
			triangles = append(triangles, faces...)
		case 0x4150: // 平滑组列表
			err := readSmoothingGroups(file, faces) // 读取平滑组列表
			if err != nil {
				return nil, err
			}
		// case 0x4160:
		// 	matrix, err := readLocalAxis(file)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	for i, v := range vertices {
		// 		vertices[i] = matrix.MulPosition(v)
		// 	}
		default:
			file.Seek(int64(header.Length-6), 1) // 跳过未知块
		}
	}

	return NewTriangleMesh(triangles), nil // 返回三角形网格
}

// readSmoothingGroups 读取平滑组列表
func readSmoothingGroups(file *os.File, triangles []*Triangle) error {
	groups := make([]uint32, len(triangles)) // 平滑组列表
	if err := binary.Read(file, binary.LittleEndian, &groups); err != nil { // 读取平滑组列表
		return err
	}
	var tables [32]map[Vector][]Vector // 平滑组表
	for i := 0; i < 32; i++ {
		tables[i] = make(map[Vector][]Vector)
	}
	for i, g := range groups {
		t := triangles[i]
		n := t.Normal()
		for j := 0; j < 32; j++ {
			if g&1 == 1 {
				tables[j][t.V1.Position] = append(tables[j][t.V1.Position], n)
				tables[j][t.V2.Position] = append(tables[j][t.V2.Position], n)
				tables[j][t.V3.Position] = append(tables[j][t.V3.Position], n)
			}
			g >>= 1
		}
	}
	for i, g := range groups {
		t := triangles[i]
		var n1, n2, n3 Vector
		for j := 0; j < 32; j++ {
			if g&1 == 1 {
				for _, v := range tables[j][t.V1.Position] {
					n1 = n1.Add(v)
				}
				for _, v := range tables[j][t.V2.Position] {
					n2 = n2.Add(v)
				}
				for _, v := range tables[j][t.V3.Position] {
					n3 = n3.Add(v)
				}
			}
			g >>= 1
		}
		t.V1.Normal = n1.Normalize()
		t.V2.Normal = n2.Normalize()
		t.V3.Normal = n3.Normalize()
	}
	return nil
}

// readLocalAxis 读取本地坐标系
func readLocalAxis(file *os.File) (Matrix, error) {
	var m [4][3]float32
	if err := binary.Read(file, binary.LittleEndian, &m); err != nil {
		return Matrix{}, err
	}
	matrix := Matrix{
		float64(m[0][0]), float64(m[0][1]), float64(m[0][2]), float64(m[3][0]),
		float64(m[1][0]), float64(m[1][1]), float64(m[1][2]), float64(m[3][1]),
		float64(m[2][0]), float64(m[2][1]), float64(m[2][2]), float64(m[3][2]),
		0, 0, 0, 1,
	}
	return matrix, nil
}

// readVertexList 读取顶点列表
func readVertexList(file *os.File) ([]Vector, error) {
	var count uint16
	if err := binary.Read(file, binary.LittleEndian, &count); err != nil { // 读取顶点数
		return nil, err
	}
	result := make([]Vector, count) // 顶点列表
	for i := range result {
		var v [3]float32
		if err := binary.Read(file, binary.LittleEndian, &v); err != nil { // 读取顶点坐标
			return nil, err
		}
		result[i] = Vector{float64(v[0]), float64(v[1]), float64(v[2])}
	}
	return result, nil
}

// readFaceList 读取面列表
func readFaceList(file *os.File, vertices []Vector) ([]*Triangle, error) {
	var count uint16
	if err := binary.Read(file, binary.LittleEndian, &count); err != nil { // 读取面数
		return nil, err
	}
	result := make([]*Triangle, count) // 面列表
	for i := range result {
		var v [4]uint16
		if err := binary.Read(file, binary.LittleEndian, &v); err != nil { // 读取面顶点索引
			return nil, err
		}
		result[i] = NewTriangleForPoints(
			vertices[v[0]], vertices[v[1]], vertices[v[2]])
	}
	return result, nil
}

// readNullTerminatedString 读取以null结尾的字符串
func readNullTerminatedString(file *os.File) (string, error) {
	var bytes []byte
	buf := make([]byte, 1)
	for {
		n, err := file.Read(buf)
		if err != nil {
			return "", err
		} else if n == 1 {
			if buf[0] == 0 {
				break
			}
			bytes = append(bytes, buf[0])
		}
	}
	return string(bytes), nil
}
