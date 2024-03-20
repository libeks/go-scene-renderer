package geometry

import "fmt"

var (
	HomogeneusIdentity = HomogeneusMatrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
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
	return m.Transpose().ScalarMult(1 / d), true
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
		0.0,

		m.B1,
		m.B2,
		m.B3,
		0.0,

		m.C1,
		m.C2,
		m.C3,
		0.0,

		0.0,
		0.0,
		0.0,
		1.0,
	}
}

type HomogeneusMatrix struct {
	A1 float64
	A2 float64
	A3 float64
	A4 float64

	B1 float64
	B2 float64
	B3 float64
	B4 float64

	C1 float64
	C2 float64
	C3 float64
	C4 float64

	D1 float64
	D2 float64
	D3 float64
	D4 float64
}

func (m HomogeneusMatrix) String() string {
	return fmt.Sprintf("M[\n\t(%.3f, %.3f, %.3f, %.3f),\n\t(%.3f, %.3f, %.3f, %.3f),\n\t(%.3f, %.3f, %.3f, %.3f),\n\t(%.3f, %.3f, %.3f, %.3f)\n]",
		m.A1, m.A2, m.A3, m.A4, m.B1, m.B2, m.B3, m.B4, m.C1, m.C2, m.C3, m.C4, m.D1, m.D2, m.D3, m.D4)
}

func (m HomogeneusMatrix) to3D() Matrix3D {
	if m.D1 == 0 && m.D2 == 0 && m.D3 == 0 {
		return Matrix3D{
			m.A1,
			m.A2,
			m.A3,

			m.B1,
			m.B2,
			m.B3,

			m.C1,
			m.C2,
			m.C3,
		}
	}
	panic(fmt.Sprintf("Cannot extract 3D matrix from Homogenous matrix %s", m))
}

func (m HomogeneusMatrix) isHomogenous() bool {
	return m.D1 == 0 && m.D2 == 0 && m.D3 == 0
}

func (m HomogeneusMatrix) Add(n HomogeneusMatrix) HomogeneusMatrix {
	return HomogeneusMatrix{
		m.A1 + n.A1,
		m.A2 + n.A2,
		m.A3 + n.A3,
		m.A4 + n.A4,

		m.B1 + n.B1,
		m.B2 + n.B2,
		m.B3 + n.B3,
		m.B4 + n.B4,

		m.C1 + n.C1,
		m.C2 + n.C2,
		m.C3 + n.C3,
		m.C4 + n.C4,

		m.D1 + n.D1,
		m.D2 + n.D2,
		m.D3 + n.D3,
		m.D4 + n.D4,
	}
}

func (m HomogeneusMatrix) ScalarMult(r float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		m.A1 * r,
		m.A2 * r,
		m.A3 * r,
		m.A4 * r,

		m.B1 * r,
		m.B2 * r,
		m.B3 * r,
		m.B4 * r,

		m.C1 * r,
		m.C2 * r,
		m.C3 * r,
		m.C4 * r,

		m.D1 * r,
		m.D2 * r,
		m.D3 * r,
		m.D4 * r,
	}
}

func (m HomogeneusMatrix) MatrixMult(n HomogeneusMatrix) HomogeneusMatrix {
	return HomogeneusMatrix{
		m.A1*n.A1 + m.A2*n.B1 + m.A3*n.C1 + m.A4*n.D1,
		m.A1*n.A2 + m.A2*n.B2 + m.A3*n.C2 + m.A4*n.D2,
		m.A1*n.A3 + m.A2*n.B3 + m.A3*n.C3 + m.A4*n.D3,
		m.A1*n.A4 + m.A2*n.B4 + m.A3*n.C4 + m.A4*n.D4,

		m.B1*n.A1 + m.B2*n.B1 + m.B3*n.C1 + m.B4*n.D1,
		m.B1*n.A2 + m.B2*n.B2 + m.B3*n.C2 + m.B4*n.D2,
		m.B1*n.A3 + m.B2*n.B3 + m.B3*n.C3 + m.B4*n.D3,
		m.B1*n.A4 + m.B2*n.B4 + m.B3*n.C4 + m.B4*n.D4,

		m.C1*n.A1 + m.C2*n.B1 + m.C3*n.C1 + m.C4*n.D1,
		m.C1*n.A2 + m.C2*n.B2 + m.C3*n.C2 + m.C4*n.D2,
		m.C1*n.A3 + m.C2*n.B3 + m.C3*n.C3 + m.C4*n.D3,
		m.C1*n.A4 + m.C2*n.B4 + m.C3*n.C4 + m.C4*n.D4,

		m.D1*n.A1 + m.D2*n.B1 + m.D3*n.C1 + m.D4*n.D1,
		m.D1*n.A2 + m.D2*n.B2 + m.D3*n.C2 + m.D4*n.D2,
		m.D1*n.A3 + m.D2*n.B3 + m.D3*n.C3 + m.D4*n.D3,
		m.D1*n.A4 + m.D2*n.B4 + m.D3*n.C4 + m.D4*n.D4,
	}
}

func (m HomogeneusMatrix) Transpose() HomogeneusMatrix {
	return HomogeneusMatrix{
		m.A1,
		m.B1,
		m.C1,
		m.D1,

		m.A2,
		m.B2,
		m.C2,
		m.D2,

		m.A3,
		m.B3,
		m.C3,
		m.D3,

		m.A4,
		m.B4,
		m.C4,
		m.D4,
	}
}

func (m HomogeneusMatrix) Determinant() float64 {
	if m.isHomogenous() {
		return m.to3D().Determinant() * m.D4
	}
	panic(fmt.Sprintf("Cannot extract determinant from Homogenous matrix %s", m))
}

// return Matrix and boolean indicating whether an inverse is even possible
func (m HomogeneusMatrix) Inverse() (HomogeneusMatrix, bool) {
	if !m.isHomogenous() || m.Determinant() == 0.0 {
		return HomogeneusMatrix{}, false
	}
	// from https://mathematica.stackexchange.com/a/106260
	upperInverse, ok := m.to3D().Inverse()
	if !ok {
		return HomogeneusMatrix{}, false
	}
	// have to add upper right 1x3 vector
	upperThird := Vector3D{m.D1, m.D2, m.D3}
	afterUpperThird := upperInverse.MultVect(upperThird)
	answer := upperInverse.toHomogenous()
	answer.D1 = afterUpperThird.X
	answer.D2 = afterUpperThird.Y
	answer.D3 = afterUpperThird.Z

	return answer, true
}

func (m HomogeneusMatrix) MultVect(v HomogenousVector) HomogenousVector {
	return HomogenousVector{
		m.A1*v.X + m.A2*v.Y + m.A3*v.Z + m.A4*v.T,
		m.B1*v.X + m.B2*v.Y + m.B3*v.Z + m.B4*v.T,
		m.C1*v.X + m.C2*v.Y + m.C3*v.Z + m.C4*v.T,
		m.D1*v.X + m.D2*v.Y + m.D3*v.Z + m.D4*v.T,
	}
}

// matrices are multiplied, from right to left, the -1 is the first one, then -2, up to 0.
func MatrixProduct(in ...HomogeneusMatrix) HomogeneusMatrix {
	result := HomogeneusIdentity
	for i := len(in) - 1; i >= 0; i-- {
		result = in[i].MatrixMult(result)
	}
	return result

}
