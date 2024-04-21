package geometry

import "fmt"

var (
	Identity3D = Matrix3D{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
)

type Matrix3D struct {
	A1 float64
	A2 float64
	A3 float64

	B1 float64
	B2 float64
	B3 float64

	C1 float64
	C2 float64
	C3 float64
}

func (m Matrix3D) String() string {
	return fmt.Sprintf("M[\n\t(%.3f, %.3f, %.3f),\n\t(%.3f, %.3f, %.3f),\n\t(%.3f, %.3f, %.3f)\n]",
		m.A1, m.A2, m.A3, m.B1, m.B2, m.B3, m.C1, m.C2, m.C3)
}

func (m Matrix3D) Add(n Matrix3D) Matrix3D {
	return Matrix3D{
		m.A1 + n.A1,
		m.A2 + n.A2,
		m.A3 + n.A3,

		m.B1 + n.B1,
		m.B2 + n.B2,
		m.B3 + n.B3,

		m.C1 + n.C1,
		m.C2 + n.C2,
		m.C3 + n.C3,
	}
}

func (m Matrix3D) ScalarMult(r float64) Matrix3D {
	return Matrix3D{
		m.A1 * r,
		m.A2 * r,
		m.A3 * r,

		m.B1 * r,
		m.B2 * r,
		m.B3 * r,

		m.C1 * r,
		m.C2 * r,
		m.C3 * r,
	}
}

func (m Matrix3D) MatrixMult(n Matrix3D) Matrix3D {
	return Matrix3D{
		m.A1*n.A1 + m.A2*n.B1 + m.A3*n.C1,
		m.A1*n.A2 + m.A2*n.B2 + m.A3*n.C2,
		m.A1*n.A3 + m.A2*n.B3 + m.A3*n.C3,

		m.B1*n.A1 + m.B2*n.B1 + m.B3*n.C1,
		m.B1*n.A2 + m.B2*n.B2 + m.B3*n.C2,
		m.B1*n.A3 + m.B2*n.B3 + m.B3*n.C3,

		m.C1*n.A1 + m.C2*n.B1 + m.C3*n.C1,
		m.C1*n.A2 + m.C2*n.B2 + m.C3*n.C2,
		m.C1*n.A3 + m.C2*n.B3 + m.C3*n.C3,
	}
}

func (m Matrix3D) Transpose() Matrix3D {
	return Matrix3D{
		m.A1,
		m.B1,
		m.C1,

		m.A2,
		m.B2,
		m.C2,

		m.A3,
		m.B3,
		m.C3,
	}
}

func (m Matrix3D) Determinant() float64 {
	return m.A1*m.B2*m.C3 +
		m.A2*m.B3*m.C1 +
		m.A3*m.B1*m.C2 -
		m.A3*m.B2*m.C1 -
		m.A2*m.B1*m.C3 -
		m.A1*m.B3*m.C2
}

// return Matrix and boolean indicating whether an inverse is even possible
func (m Matrix3D) Inverse() (Matrix3D, bool) {
	d := m.Determinant()
	if d == 0 {
		return Matrix3D{}, false
	}
	fmt.Printf("got %.3f, %.3f = %.3f \n", m.B1, m.C2, m.B1*m.C2)
	k := Matrix3D{
		m.B2*m.C3 - m.B3*m.C2,
		-m.B1*m.C3 + m.B3*m.C1,
		m.B1*m.C2 - m.B2*m.C1,

		-m.A2*m.C3 + m.C2*m.A3,
		m.A1*m.C3 - m.A3*m.C1,
		-m.A1*m.C2 + m.C1*m.A2,

		m.A2*m.B3 - m.A3*m.B2, // should be 0...
		-m.A1*m.B3 + m.B1*m.A3,
		m.A1*m.B2 - m.A2*m.B1,
	}
	return k.Transpose().ScalarMult(1 / d), true
}

func (m Matrix3D) MultVect(v Vector3D) Vector3D {
	return Vector3D{
		m.A1*v.X + m.A2*v.Y + m.A3*v.Z,
		m.B1*v.X + m.B2*v.Y + m.B3*v.Z,
		m.C1*v.X + m.C2*v.Y + m.C3*v.Z,
	}
}

func (m Matrix3D) toHomogenous() HomogeneusMatrix {
	return HomogeneusMatrix{
		m.A1,
		m.A2,
		m.A3,
		0,

		m.B1,
		m.B2,
		m.B3,
		0,

		m.C1,
		m.C2,
		m.C3,
		0,

		0,
		0,
		0,
		1,
	}
}
