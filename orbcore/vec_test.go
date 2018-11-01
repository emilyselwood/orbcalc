package orbcore

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestRotation(t *testing.T) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop
	result1 := Rotate(mat.NewVecDense(3, []float64{1, 1, 1}), angle, AxisZ)
	result1 = Rotate(result1, angle, AxisX)
	result1 = Rotate(result1, angle, AxisZ)

	rot2 := RotationMatrix(angle, AxisZ)
	rot2.Mul(rot2, RotationMatrix(angle, AxisX))
	rot2.Mul(rot2, RotationMatrix(angle, AxisZ))

	result2 := mat.NewVecDense(3, []float64{1, 1, 1})
	result2.MulVec(rot2, result2)
	if !mat.EqualApprox(result1, result2, 0.0000000000000001) {
		t.Log("\nresult1:", mat.Formatted(result1, mat.Prefix("result1: "), mat.Squeeze()))
		t.Log("\nresult2:", mat.Formatted(result2, mat.Prefix("result2: "), mat.Squeeze()))
		t.Fail()
	}
}

func TestQuickerRotationMatrixForOrbit(t *testing.T) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop

	result1 := Rotate(mat.NewVecDense(3, []float64{1, 1, 1}), angle, AxisZ)
	result1 = Rotate(result1, angle, AxisX)
	result1 = Rotate(result1, angle, AxisZ)

	rot2 := QuickerRotationMatrixForOrbit(angle, angle, angle)

	result2 := mat.NewVecDense(3, []float64{1, 1, 1})
	result2.MulVec(rot2, result2)

	if !mat.EqualApprox(result1, result2, 0.0000000000000001) {
		t.Log("\nresult1:", mat.Formatted(result1, mat.Prefix("result1: "), mat.Squeeze()))
		t.Log("\nresult2:", mat.Formatted(result2, mat.Prefix("result2: "), mat.Squeeze()))
		t.Fail()
	}

}
