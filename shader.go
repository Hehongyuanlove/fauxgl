package fauxgl

import "math"

// Shader 接口
type Shader interface {
	Vertex(Vertex) Vertex // 顶点着色器
	Fragment(Vertex) Color // 片元着色器
}

// SolidColorShader 渲染单一颜色
type SolidColorShader struct {
	Matrix Matrix // 变换矩阵
	Color  Color // 颜色
}

// NewSolidColorShader 创建一个渲染单一颜色的着色器
func NewSolidColorShader(matrix Matrix, color Color) *SolidColorShader {
	return &SolidColorShader{matrix, color}
}

// Vertex 顶点着色器
func (shader *SolidColorShader) Vertex(v Vertex) Vertex {
	v.Output = shader.Matrix.MulPositionW(v.Position)
	return v
}

// Fragment 片元着色器
func (shader *SolidColorShader) Fragment(v Vertex) Color {
	return shader.Color
}

// TextureShader 渲染纹理
type TextureShader struct {
	Matrix  Matrix
	Texture Texture
}

// NewTextureShader 创建一个渲染纹理的着色器
func NewTextureShader(matrix Matrix, texture Texture) *TextureShader {
	return &TextureShader{matrix, texture}
}

// Vertex 顶点着色器
func (shader *TextureShader) Vertex(v Vertex) Vertex {
	v.Output = shader.Matrix.MulPositionW(v.Position)
	return v
}

// Fragment 片元着色器
func (shader *TextureShader) Fragment(v Vertex) Color {
	return shader.Texture.BilinearSample(v.Texture.X, v.Texture.Y)
}

// PhongShader 实现冯氏着色法
type PhongShader struct {
	Matrix         Matrix
	LightDirection Vector
	CameraPosition Vector
	ObjectColor    Color
	AmbientColor   Color
	DiffuseColor   Color
	SpecularColor  Color
	Texture        Texture
	SpecularPower  float64
}
// NewPhongShader 创建一个实现冯氏着色法的着色器
func NewPhongShader(matrix Matrix, lightDirection, cameraPosition Vector) *PhongShader {
	ambient := Color{0.2, 0.2, 0.2, 1}
	diffuse := Color{0.8, 0.8, 0.8, 1}
	specular := Color{1, 1, 1, 1}
	return &PhongShader{
		matrix, lightDirection, cameraPosition,
		Discard, ambient, diffuse, specular, nil, 32}
}

// Vertex 顶点着色器
func (shader *PhongShader) Vertex(v Vertex) Vertex {
	v.Output = shader.Matrix.MulPositionW(v.Position)
	return v
}

// Fragment 片元着色器
func (shader *PhongShader) Fragment(v Vertex) Color {
	light := shader.AmbientColor
	color := v.Color
	if shader.ObjectColor != Discard {
		color = shader.ObjectColor
	}
	if shader.Texture != nil {
		color = shader.Texture.BilinearSample(v.Texture.X, v.Texture.Y)
	}
	diffuse := math.Max(v.Normal.Dot(shader.LightDirection), 0)
	light = light.Add(shader.DiffuseColor.MulScalar(diffuse))
	if diffuse > 0 && shader.SpecularPower > 0 {
		camera := shader.CameraPosition.Sub(v.Position).Normalize()
		reflected := shader.LightDirection.Negate().Reflect(v.Normal)
		specular := math.Max(camera.Dot(reflected), 0)
		if specular > 0 {
			specular = math.Pow(specular, shader.SpecularPower)
			light = light.Add(shader.SpecularColor.MulScalar(specular))
		}
	}
	return color.Mul(light).Min(White).Alpha(color.A)
}
