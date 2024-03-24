package geometry

import "fmt"

type Matrix2D struct {
	A1 float64
	A2 float64

	B1 float64
	B2 float64
}

func (m Matrix2D) String() string {
	return fmt.Sprintf("M[\n\t(%.3f, %.3f),\n\t(%.3f, %.3f),\n]",
		m.A1, m.A2, m.B1, m.B2)
}

func (m Matrix2D) Add(n Matrix2D) Matrix2D {
	return Matrix2D{
		m.A1 + n.A1,
		m.A2 + n.A2,

		m.B1 + n.B1,
		m.B2 + n.B2,
	}
}

func (m Matrix2D) ScalarMult(r float64) Matrix2D {
	return Matrix2D{
		m.A1 * r,
		m.A2 * r,

		m.B1 * r,
		m.B2 * r,
	}
}

func (m Matrix2D) MatrixMult(n Matrix2D) Matrix2D {
	return Matrix2D{
		m.A1*n.A1 + m.A2*n.B1,
		m.A1*n.A2 + m.A2*n.B2,

		m.B1*n.A1 + m.B2*n.B1,
		m.B1*n.A2 + m.B2*n.B2,
	}
}

func (m Matrix2D) Transpose() Matrix2D {
	return Matrix2D{
		m.A1,
		m.B1,

		m.A2,
		m.B2,
	}
}

func (m Matrix2D) Determinant() float64 {
	return m.A1*m.B2 -
		m.A2*m.B1
}

// return Matrix and boolean indicating whether an inverse is even possible
func (m Matrix2D) Inverse() (Matrix2D, bool) {
	d := m.Determinant()
	if d == 0 {
		return Matrix2D{}, false
	}
	return Matrix2D{
		m.B2, -m.A2,
		-m.B1, m.A1,
	}.ScalarMult(1 / d), true
}

func (m Matrix2D) MultVect(v Vector2D) Vector2D {
	return Vector2D{
		m.A1*v.X + m.A2*v.Y,
		m.B1*v.X + m.B2*v.Y,
	}
}
