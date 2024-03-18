package geometry

import "fmt"

type Line struct {
	A Point
	B Point
}

func (l Line) String() string {
	return fmt.Sprintf("Line %s %s", l.A, l.B)
}

// func (l Line) Project() ProjectedLine {
// 	return
// }

// type Projectedline struct {
// 	A ScenePixel
// 	B ScenePixel
// }

// func (l Projectedline) String() string {
// 	retur fmt.Sprintf("ProjectedLine %s %s", l.A, l.B)
// }
