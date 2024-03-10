package geometry

import "math"

func TranslationMatrix(v Vector3D) HomogeneusMatrix {
	return HomogeneusMatrix{
		1.0,
		0.0,
		0.0,
		v.X,

		0.0,
		1.0,
		0.0,
		v.Y,

		0.0,
		0.0,
		1.0,
		v.Z,

		0.0,
		0.0,
		0.0,
		1.0,
	}
}

func ScaleMatrix(t float64) HomogeneusMatrix {
	if t == 0.0 {
		panic("Cannot scale by 0.0")
	}
	return HomogeneusMatrix{
		1.0,
		0.0,
		0.0,
		0.0,

		0.0,
		1.0,
		0.0,
		0.0,

		0.0,
		0.0,
		1.0,
		0.0,

		0.0,
		0.0,
		0.0,
		1.0 / t,
	}
}

func RotateMatrixX(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		1,
		0,
		0,
		0,

		0,
		math.Cos(t),
		-math.Sin(t),
		0,

		0,
		math.Sin(t),
		math.Cos(t),
		0,

		0,
		0,
		0,
		1,
	}
}

func RotateMatrixY(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		math.Cos(t),
		0,
		math.Sin(t),
		0,

		0,
		1,
		0,
		0,

		-math.Sin(t),
		0,
		math.Cos(t),
		0,

		0,
		0,
		0,
		1,
	}
}

func RotateMatrixZ(t float64) HomogeneusMatrix {
	return HomogeneusMatrix{
		math.Cos(t),
		-math.Sin(t),
		0,
		0,

		math.Sin(t),
		math.Cos(t),
		0,
		0,

		0,
		0,
		1,
		0,

		0,
		0,
		0,
		1,
	}
}

// func RollPitchYawMatrix(r, p, y float64) HomogeneusMatrix {
// 	// `\cos\alpha\cos\beta &
// 	//       \cos\alpha\sin\beta\sin\gamma - \sin\alpha\cos\gamma &
// 	//       \cos\alpha\sin\beta\cos\gamma + \sin\alpha\sin\gamma \\

// 	//     \sin\alpha\cos\beta &
// 	//       \sin\alpha\sin\beta\sin\gamma + \cos\alpha\cos\gamma &
// 	//       \sin\alpha\sin\beta\cos\gamma - \cos\alpha\sin\gamma \\

// 	//    -\sin\beta & \cos\beta\sin\gamma & \cos\beta\cos\gamma \\
// 	//    `
// 	return HomogeneusMatrix{
// 		math.Cos(r) * math.Cos(p),
// 		math.Cos(r)*math.Sin(p)*math.Sin(y) - math.Sin(r)*math.Cos(y),
// 		math.Sin(r) * math.Cos(p),
// 		0.0,

// 		math.Cos(r) * math.Cos(p),
// 		math.Sin(r)*math.Sin(p)*math.Sin(y) + math.Cos(r)*math.Cos(y),
// 		math.Sin(r)*math.Sin(p)*math.Cos(y) - math.Cos(r)*math.Sin(y),
// 		0.0,

// 		-math.Sin(p),
// 		math.Cos(p) * math.Sin(y),
// 		math.Cos(p) * math.Cos(y),
// 		0.0,

// 		0.0,
// 		0.0,
// 		0.0,
// 		1.0,
// 	}
// }
