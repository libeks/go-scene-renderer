package geometry

import "fmt"

type Line struct {
	A Point
	B Point
}

func (l Line) String() string {
	return fmt.Sprintf("Line %s %s", l.A, l.B)
}
