package orbcore

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

/*
Axis represents one of the cartisen axis
*/
type Axis int

/*
Constants for three major axis
*/
const (
	AxisX Axis = 0
	AxisY      = 1
	AxisZ      = 2
)

/*
Rotate rotates the vector around the given axis
*/
func Rotate(vec *mat.VecDense, angle float64, axis Axis) *mat.VecDense {
	rot := RotationMatrix(angle, axis)
	vec.MulVec(rot, vec)
	return vec
}

/*
RotationMatrixForOrbit creates a rotation matrix that is the three rotations we need for sorting out an orbit vector
Using this means we have to 4 matrix multiplications rather than 6

If this is in the hot path of your code consider using QuickerRotationMatrixForOrbit instead.
*/
func RotationMatrixForOrbit(angleZ1, angleX, angleZ2 float64) *mat.Dense {
	rot := RotationMatrix(angleZ1, AxisZ)
	rot.Mul(rot, RotationMatrix(angleX, AxisX))
	rot.Mul(rot, RotationMatrix(angleZ2, AxisZ))
	return rot
}

/*
QuickerRotationMatrixForOrbit builds a rotation matrix for an orbit by manually doing the matrix multiplication
This is a lot faster than the RotationMatrixForOrbit.

*/
func QuickerRotationMatrixForOrbit(angleZ1, angleX, angleZ2 float64) *mat.Dense {
	z1c := math.Cos(angleZ1)
	z1s := math.Sin(angleZ1)
	xc := math.Cos(angleX)
	xs := math.Sin(angleX)
	z2c := math.Cos(angleZ2)
	z2s := math.Sin(angleZ2)

	// TODO: simplify out the zeros. Left in for now to make the fomula easier to validate.
	a := ((z1c * 1 * z2c) + ((-z1s) * 0 * z2c) + (0 * xs * z2c)) +
		((z1c * 0 * z2s) + ((-z1s) * xc * z2s) + (0 * 0 * z2s)) +
		((z1c * 0 * 0) + ((-z1s) * (-xs) * 0) + (0 * xc * 0))

	b := ((z1c * 1 * (-z2s)) + ((-z1s) * 0 * (-z2s)) + (0 * xs * (-z2s))) +
		((z1c * 0 * z2c) + ((-z1s) * xc * z2c) + (0 * 0 * z2c)) +
		((z1c * 0 * 0) + ((-z1s) * (-xs) * 0) + (0 * xc * 0))

	c := ((z1c * 1 * 0) + ((-z1s) * 0 * 0) + (0 * xs * 0)) +
		((z1c * 0 * 0) + ((-z1s) * xc * 0) + (0 * 0 * 0)) +
		((z1c * 0 * 1) + ((-z1s) * (-xs) * 1) + (0 * xc * 1))

	d := ((z1s * 1 * z2c) + (z1c * 0 * z2c) + (0 * 0 * z2c)) +
		((z1s * 0 * z2s) + (z1c * xc * z2s) + (0 * 0 * z2s)) +
		((z1s * 0 * 0) + (z1c * (-xs) * 0) + (0 * xc * 0))

	e := ((z1s * 1 * (-z2s)) + (z1c * 0 * (-z2s)) + (0 * 0 * (-z2s))) +
		((z1s * 0 * z2c) + (z1c * xc * z2c) + (0 * 0 * z2c)) +
		((z1s * 0 * 0) + (z1c * (-xs) * 0) + (0 * xc * 0))

	f := ((z1s * 1 * 0) + (z1c * 0 * 0) + (0 * 0 * 0)) +
		((z1s * 0 * 0) + (z1c * 1 * 0) + (0 * 0 * 0)) +
		((z1s * 0 * 1) + (z1c * (-xs) * 1) + (0 * xc * 1))

	g := ((0 * 1 * z2c) + (0 * 0 * z2c) + (1 * 0 * z2c)) +
		((0 * 0 * z2s) + (0 * xc * z2s) + (1 * xs * z2s)) +
		((0 * 0 * 0) + (0 * -xs * 0) + (1 * xc * 0))

	h := ((0 * 1 * (-z2s)) + (0 * 0 * (-z2s)) + (1 * 0 * (-z2s))) +
		((0 * 0 * z2c) + (0 * xc * z2c) + (1 * xs * z2c)) +
		((0 * 0 * 0) + (0 * -xs * 0) + (1 * xc * 0))

	i := ((0 * 1 * 0) + (0 * 0 * 0) + (1 * 0 * 0)) +
		((0 * 0 * 0) + (0 * xc * 0) + (1 * xs * 0)) +
		((0 * 0 * 1) + (0 * (-xs) * 1) + (1 * xc * 1))

	rot := mat.NewDense(3, 3, []float64{a, b, c, d, e, f, g, h, i})

	return rot
}

/*
RotationMatrix returns a rotation matrix for the given angle and axis
Will panic if the axis is not one of AxisX, AxisY, AxisZ
*/
func RotationMatrix(angle float64, axis Axis) *mat.Dense {
	c := math.Cos(angle)
	s := math.Sin(angle)

	switch axis {
	case AxisX:
		return mat.NewDense(3, 3, []float64{
			1, 0, 0,
			0, c, -s,
			0, s, c,
		})
	case AxisY:
		return mat.NewDense(3, 3, []float64{
			c, 0, s,
			0, 1, 0,
			s, 0, c,
		})
	case AxisZ:
		return mat.NewDense(3, 3, []float64{
			c, -s, 0,
			s, c, 0,
			0, 0, 1,
		})
	default:
		panic(fmt.Sprintf("Unknown axis %v", axis))
	}

}
