package color

type Texture interface {
	// a,b range from (0,1), when used for triangles, only the lower triangluar values will be called
	// TODO: Texture is surprisingly similar to Frame, maybe there's some generalizations?
	GetTextureColor(b, c float64) Color
}
