package fauxgl

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)
// VOXHeader VOX文件头部
type VOXHeader struct {
	Magic   [4]byte // 魔数
	Version int32   // 版本号
}

// VOXChunk VOX文件块
type VOXChunk struct {
	ID            [4]byte // 块ID
	ContentBytes  int32   // 块内容大小
	ChildrenBytes int32   // 子块大小
}

// VOXVoxel VOX体素
type VOXVoxel struct {
	X, Y, Z, I uint8 // 体素坐标和颜色索引
}

// LoadVOX 从文件中读取voxel数据
func LoadVOX(path string) ([]Voxel, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取并检查头部
	header := VOXHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	if string(header.Magic[:]) != "VOX " {
		return nil, errors.New("invalid vox header")
	}
	if header.Version != 150 {
		return nil, errors.New("unsupported vox version")
	}

	var voxVoxels []VOXVoxel
	var palette [256]Color

	for i := range palette {
		x := voxDefaultPalette[i]
		r := float64((x>>0)&255) / 255
		g := float64((x>>8)&255) / 255
		b := float64((x>>16)&255) / 255
		a := float64((x>>24)&255) / 255
		palette[i] = Color{r, g, b, a}
	}

	for {
		// 读取chunk头部
		chunk := VOXChunk{}
		if err := binary.Read(file, binary.LittleEndian, &chunk); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}

		id := string(chunk.ID[:])
		switch id {
		case "XYZI":
			var numVoxels uint32
			if err := binary.Read(file, binary.LittleEndian, &numVoxels); err != nil {
				return nil, err
			}
			for i := 0; i < int(numVoxels); i++ {
				voxel := VOXVoxel{}
				if err := binary.Read(file, binary.LittleEndian, &voxel); err != nil {
					return nil, err
				}
				voxVoxels = append(voxVoxels, voxel)
			}
		case "RGBA":
			for i := 0; i <= 254; i++ {
				var color [4]uint8
				if err := binary.Read(file, binary.LittleEndian, &color); err != nil {
					return nil, err
				}
				r := float64(color[0]) / 255
				g := float64(color[1]) / 255
				b := float64(color[2]) / 255
				a := float64(color[3]) / 255
				palette[i+1] = Color{r, g, b, a}
			}
		default:
			file.Seek(int64(chunk.ContentBytes), 1)
		}
	}

	voxels := make([]Voxel, len(voxVoxels))
	for i, v := range voxVoxels {
		voxels[i] = Voxel{int(v.X), int(v.Y), int(v.Z), palette[v.I]}
	}

	return voxels, nil
}

var voxDefaultPalette = []uint{
	0x00000000, 0xffffffff, 0xffccffff, 0xff99ffff,
	0xff66ffff, 0xff33ffff, 0xff00ffff, 0xffffccff,
	0xffccccff, 0xff99ccff, 0xff66ccff, 0xff33ccff,
	0xff00ccff, 0xffff99ff, 0xffcc99ff, 0xff9999ff,
	0xff6699ff, 0xff3399ff, 0xff0099ff, 0xffff66ff,
	0xffcc66ff, 0xff9966ff, 0xff6666ff, 0xff3366ff,
	0xff0066ff, 0xffff33ff, 0xffcc33ff, 0xff9933ff,
	0xff6633ff, 0xff3333ff, 0xff0033ff, 0xffff00ff,
	0xffcc00ff, 0xff9900ff, 0xff6600ff, 0xff3300ff,
	0xff0000ff, 0xffffffcc, 0xffccffcc, 0xff99ffcc,
	0xff66ffcc, 0xff33ffcc, 0xff00ffcc, 0xffffcccc,
	0xffcccccc, 0xff99cccc, 0xff66cccc, 0xff33cccc,
	0xff00cccc, 0xffff99cc, 0xffcc99cc, 0xff9999cc,
	0xff6699cc, 0xff3399cc, 0xff0099cc, 0xffff66cc,
	0xffcc66cc, 0xff9966cc, 0xff6666cc, 0xff3366cc,
	0xff0066cc, 0xffff33cc, 0xffcc33cc, 0xff9933cc,
	0xff6633cc, 0xff3333cc, 0xff0033cc, 0xffff00cc,
	0xffcc00cc, 0xff9900cc, 0xff6600cc, 0xff3300cc,
	0xff0000cc, 0xffffff99, 0xffccff99, 0xff99ff99,
	0xff66ff99, 0xff33ff99, 0xff00ff99, 0xffffcc99,
	0xffcccc99, 0xff99cc99, 0xff66cc99, 0xff33cc99,
	0xff00cc99, 0xffff9999, 0xffcc9999, 0xff999999,
	0xff669999, 0xff339999, 0xff009999, 0xffff6699,
	0xffcc6699, 0xff996699, 0xff666699, 0xff336699,
	0xff006699, 0xffff3399, 0xffcc3399, 0xff993399,
	0xff663399, 0xff333399, 0xff003399, 0xffff0099,
	0xffcc0099, 0xff990099, 0xff660099, 0xff330099,
	0xff000099, 0xffffff66, 0xffccff66, 0xff99ff66,
	0xff66ff66, 0xff33ff66, 0xff00ff66, 0xffffcc66,
	0xffcccc66, 0xff99cc66, 0xff66cc66, 0xff33cc66,
	0xff00cc66, 0xffff9966, 0xffcc9966, 0xff999966,
	0xff669966, 0xff339966, 0xff009966, 0xffff6666,
	0xffcc6666, 0xff996666, 0xff666666, 0xff336666,
	0xff006666, 0xffff3366, 0xffcc3366, 0xff993366,
	0xff663366, 0xff333366, 0xff003366, 0xffff0066,
	0xffcc0066, 0xff990066, 0xff660066, 0xff330066,
	0xff000066, 0xffffff33, 0xffccff33, 0xff99ff33,
	0xff66ff33, 0xff33ff33, 0xff00ff33, 0xffffcc33,
	0xffcccc33, 0xff99cc33, 0xff66cc33, 0xff33cc33,
	0xff00cc33, 0xffff9933, 0xffcc9933, 0xff999933,
	0xff669933, 0xff339933, 0xff009933, 0xffff6633,
	0xffcc6633, 0xff996633, 0xff666633, 0xff336633,
	0xff006633, 0xffff3333, 0xffcc3333, 0xff993333,
	0xff663333, 0xff333333, 0xff003333, 0xffff0033,
	0xffcc0033, 0xff990033, 0xff660033, 0xff330033,
	0xff000033, 0xffffff00, 0xffccff00, 0xff99ff00,
	0xff66ff00, 0xff33ff00, 0xff00ff00, 0xffffcc00,
	0xffcccc00, 0xff99cc00, 0xff66cc00, 0xff33cc00,
	0xff00cc00, 0xffff9900, 0xffcc9900, 0xff999900,
	0xff669900, 0xff339900, 0xff009900, 0xffff6600,
	0xffcc6600, 0xff996600, 0xff666600, 0xff336600,
	0xff006600, 0xffff3300, 0xffcc3300, 0xff993300,
	0xff663300, 0xff333300, 0xff003300, 0xffff0000,
	0xffcc0000, 0xff990000, 0xff660000, 0xff330000,
	0xff0000ee, 0xff0000dd, 0xff0000bb, 0xff0000aa,
	0xff000088, 0xff000077, 0xff000055, 0xff000044,
	0xff000022, 0xff000011, 0xff00ee00, 0xff00dd00,
	0xff00bb00, 0xff00aa00, 0xff008800, 0xff007700,
	0xff005500, 0xff004400, 0xff002200, 0xff001100,
	0xffee0000, 0xffdd0000, 0xffbb0000, 0xffaa0000,
	0xff880000, 0xff770000, 0xff550000, 0xff440000,
	0xff220000, 0xff110000, 0xffeeeeee, 0xffdddddd,
	0xffbbbbbb, 0xffaaaaaa, 0xff888888, 0xff777777,
	0xff555555, 0xff444444, 0xff222222, 0xff111111,
}
