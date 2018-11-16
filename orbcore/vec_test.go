package orbcore

import (
	"fmt"
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
	if !mat.EqualApprox(result1, result2, 0.0000000001) {
		t.Log("\nresult1:", mat.Formatted(result1, mat.Prefix("result1: "), mat.Squeeze()))
		t.Log("\nresult2:", mat.Formatted(result2, mat.Prefix("result2: "), mat.Squeeze()))
		t.Fail()
	}
}

func TestQuickerRotationMatrixForOrbit(t *testing.T) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop

	result1 := Rotate(mat.NewVecDense(3, []float64{10000, 10000, 10000}), angle, AxisZ)
	result1 = Rotate(result1, angle, AxisX)
	result1 = Rotate(result1, angle, AxisZ)

	rot2 := QuickerRotationMatrixForOrbit(angle, angle, angle)

	result2 := mat.NewVecDense(3, []float64{10000, 10000, 10000})
	result2.MulVec(rot2, result2)

	if !mat.EqualApprox(result1, result2, 0.0000000001) {
		t.Log("\nresult1:", mat.Formatted(result1, mat.Prefix("result1: "), mat.Squeeze()))
		t.Log("\nresult2:", mat.Formatted(result2, mat.Prefix("result2: "), mat.Squeeze()))
		t.Fail()
	}

}

func TestTestQuickerRotationMatrixForOrbitLong(t *testing.T) {
	for i := 0.0; i < 180.0; i = i + 0.1 {
		t.Run(fmt.Sprintf("deg:%v", i), func(t2 *testing.T) {
			angle := float64(i) * math.Pi / 180.0

			result1 := Rotate(mat.NewVecDense(3, []float64{5e+09, 5e+9, 0}), angle, AxisZ)
			result1 = Rotate(result1, angle, AxisX)
			result1 = Rotate(result1, angle, AxisZ)

			rot2 := QuickerRotationMatrixForOrbit(angle, angle, angle)

			result2 := mat.NewVecDense(3, []float64{5e+9, 5e+9, 0})
			result2.MulVec(rot2, result2)

			if !mat.EqualApprox(result1, result2, 0.0000001) {
				t2.Log("\nresult1:", mat.Formatted(result1, mat.Prefix("result1: "), mat.Squeeze()))
				t2.Log("\nresult2:", mat.Formatted(result2, mat.Prefix("result2: "), mat.Squeeze()))
				t2.Fail()
			}
		})
	}
}

func TestRotationMatrixGeneration(t *testing.T) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop
	m1 := RotationMatrix(angle, AxisX)
	m2 := QuickerRotationMatrixForOrbit(0, angle, 0)

	if !mat.EqualApprox(m1, m2, 0.000000000000001) {
		t.Log("\nresult1:", mat.Formatted(m1, mat.Prefix("m1: "), mat.Squeeze()))
		t.Log("\nresult2:", mat.Formatted(m2, mat.Prefix("m2: "), mat.Squeeze()))
		t.Fail()
	}

}

func BenchmarkQuickerRotationMatrix(b *testing.B) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop
	for n := 0; n < b.N; n++ {
		QuickerRotationMatrixForOrbit(angle, angle, angle)
	}
}

func BenchmarkRotationMatrix(b *testing.B) {
	angle := 45 * math.Pi / 180.0 // Not using the function in orbconvert because of an import loop
	for n := 0; n < b.N; n++ {
		RotationMatrixForOrbit(angle, angle, angle)
	}
}
