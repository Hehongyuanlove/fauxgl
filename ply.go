package fauxgl

import (
	"bufio"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

// ply格式
type plyFormat int

const (
	_                     plyFormat = iota
	plyAscii                        // ascii格式
	plyBinaryLittleEndian           // 小端字节序
	plyBinaryBigEndian              // 大端字节序
)

// ply格式映射
var plyFormatMapping = map[string]plyFormat{
	"ascii":                plyAscii,              // ascii格式
	"binary_little_endian": plyBinaryLittleEndian, // 小端字节序
	"binary_big_endian":    plyBinaryBigEndian,    // 大端字节序
}

// ply数据类型
type plyDataType int

const (
	plyNone    plyDataType = iota // 无类型
	plyInt8                       // 8位整型
	plyUint8                      // 8位无符号整型
	plyInt16                      // 16位整型
	plyUint16                     // 16位无符号整型
	plyInt32                      // 32位整型
	plyUint32                     // 32位无符号整型
	plyFloat32                    // 32位单精度浮点型
	plyFloat64                    // 64位双精度浮点型
)

// ply数据类型映射
var plyDataTypeMapping = map[string]plyDataType{
	"char":    plyInt8,    // 字符型
	"uchar":   plyUint8,   // 无符号字符型
	"short":   plyInt16,   // 短整型
	"ushort":  plyUint16,  // 无符号短整型
	"int":     plyInt32,   // 整型
	"uint":    plyUint32,  // 无符号整型
	"float":   plyFloat32, // 单精度浮点型
	"double":  plyFloat64, // 双精度浮点型
	"int8":    plyInt8,    // 8位整型
	"uint8":   plyUint8,   // 8位无符号整型
	"int16":   plyInt16,   // 16位整型
	"uint16":  plyUint16,  // 16位无符号整型
	"int32":   plyInt32,   // 32位整型
	"uint32":  plyUint32,  // 32位无符号整型
	"float32": plyFloat32, // 32位单精度浮点型
	"float64": plyFloat64, // 64位双精度浮点型
}

// ply属性
type plyProperty struct {
	name      string      // 属性名
	countType plyDataType // 计数类型
	dataType  plyDataType // 数据类型
}

// ply元素
type plyElement struct {
	name       string        // 元素名
	count      int           // 元素数量
	properties []plyProperty // 属性列表
}

// LoadPLY 从PLY文件中加载网格
func LoadPLY(path string) (*Mesh, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取头部信息
	reader := bufio.NewReader(file)
	var element plyElement
	var elements []plyElement
	format := plyAscii
	bytes := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		bytes += len(line)
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		if f[0] == "format" {
			format = plyFormatMapping[f[1]]
		}
		if f[0] == "element" {
			if element.count > 0 {
				elements = append(elements, element)
			}
			name := f[1]
			count, _ := strconv.ParseInt(f[2], 0, 0)
			element = plyElement{name, int(count), nil}
		}
		if f[0] == "property" {
			if f[1] == "list" {
				countType := plyDataTypeMapping[f[2]]
				dataType := plyDataTypeMapping[f[3]]
				name := f[4]
				property := plyProperty{name, countType, dataType}
				element.properties = append(element.properties, property)
			} else {
				countType := plyNone
				dataType := plyDataTypeMapping[f[1]]
				name := f[2]
				property := plyProperty{name, countType, dataType}
				element.properties = append(element.properties, property)
			}
		}
		if f[0] == "end_header" {
			if element.count > 0 {
				elements = append(elements, element)
			}
			break
		}
	}

	file.Seek(int64(bytes), 0)

	switch format {
	case plyBinaryBigEndian:
		return loadPlyBinary(file, elements, binary.BigEndian)
	case plyBinaryLittleEndian:
		return loadPlyBinary(file, elements, binary.LittleEndian)
	default:
		return loadPlyAscii(file, elements)
	}
}

// 从PLY文件中加载网格（ASCII格式）
func loadPlyAscii(file *os.File, elements []plyElement) (*Mesh, error) {
	scanner := bufio.NewScanner(file)
	var vertexes []Vector
	var triangles []*Triangle
	for _, element := range elements {
		for i := 0; i < element.count; i++ {
			scanner.Scan()
			line := scanner.Text()
			f := strings.Fields(line)
			fi := 0
			vertex := Vector{}
			for _, property := range element.properties {
				if property.name == "x" {
					vertex.X, _ = strconv.ParseFloat(f[fi], 64)
					vertex.X, _ = strconv.ParseFloat(f[fi], 64) // 解析x坐标
				}
				if property.name == "y" {
					vertex.Y, _ = strconv.ParseFloat(f[fi], 64)
					vertex.Y, _ = strconv.ParseFloat(f[fi], 64) // 解析y坐标
				}
				if property.name == "z" {
					vertex.Z, _ = strconv.ParseFloat(f[fi], 64)
					vertex.Z, _ = strconv.ParseFloat(f[fi], 64) // 解析z坐标
				}
				if property.name == "vertex_indices" {
					i1, _ := strconv.ParseInt(f[fi+1], 0, 0)
					i2, _ := strconv.ParseInt(f[fi+2], 0, 0)
					i3, _ := strconv.ParseInt(f[fi+3], 0, 0)
					t := Triangle{}
					t.V1.Position = vertexes[i1]
					t.V2.Position = vertexes[i2]
					t.V3.Position = vertexes[i3]
					t.FixNormals()
					triangles = append(triangles, &t)
					fi += 3
				}
				fi++
			}
			if element.name == "vertex" {
				vertexes = append(vertexes, vertex)
				vertexes = append(vertexes, vertex) // 添加顶点
			}
		}
	}
	return NewTriangleMesh(triangles), nil
}

// 从PLY文件中加载网格（二进制格式）
func loadPlyBinary(file *os.File, elements []plyElement, order binary.ByteOrder) (*Mesh, error) {
	var vertexes []Vector     // 顶点列表
	var triangles []*Triangle // 三角形列表
	for _, element := range elements {
		for i := 0; i < element.count; i++ {
			var vertex Vector   // 顶点
			var points []Vector // 点列表
			for _, property := range element.properties {
				if property.countType == plyNone { // 非列表类型
					value, err := readPlyFloat(file, order, property.dataType) // 读取浮点数
					if err != nil {
						return nil, err
					}
					if property.name == "x" { // x坐标
						vertex.X = value // 设置x坐标
					}
					if property.name == "y" { // y坐标
						vertex.Y = value // 设置y坐标
					}
					if property.name == "z" { // z坐标
						vertex.Z = value // 设置z坐标
					}
				} else { // 列表类型
					count, err := readPlyInt(file, order, property.countType) // 读取计数
					if err != nil {
						return nil, err
					}
					for j := 0; j < count; j++ {
						value, err := readPlyInt(file, order, property.dataType) // 读取整数
						if err != nil {
							return nil, err
						}
						if property.name == "vertex_indices" { // 顶点索引
							points = append(points, vertexes[value]) // 添加点
						}
					}
				}
			}
			if element.name == "vertex" { // 顶点
				vertexes = append(vertexes, vertex) // 添加顶点
			}
			if element.name == "face" { // 面
				t := Triangle{}
				t.V1.Position = points[0]         // 设置第一个顶点
				t.V2.Position = points[1]         // 设置第二个顶点
				t.V3.Position = points[2]         // 设置第三个顶点
				t.FixNormals()                    // 修正法向量
				triangles = append(triangles, &t) // 添加三角形
			}
		}
	}
	return NewTriangleMesh(triangles), nil // 返回网格
}

// 从PLY文件中读取整数
func readPlyInt(file *os.File, order binary.ByteOrder, dataType plyDataType) (int, error) {
	value, err := readPlyFloat(file, order, dataType) // 读取浮点数
	return int(value), err
}

// 从PLY文件中读取浮点数
func readPlyFloat(file *os.File, order binary.ByteOrder, dataType plyDataType) (float64, error) {
	switch dataType {
	case plyInt8:
		var value int8
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyUint8:
		var value uint8
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyInt16:
		var value int16
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyUint16:
		var value uint16
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyInt32:
		var value int32
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyUint32:
		var value uint32
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyFloat32:
		var value float32
		err := binary.Read(file, order, &value)
		return float64(value), err
	case plyFloat64:
		var value float64
		err := binary.Read(file, order, &value)
		return float64(value), err
	default:
		return 0, nil
	}
}
