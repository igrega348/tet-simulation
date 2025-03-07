package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/mat"
)

var Ainv *mat.Dense
var P *mat.Dense
var params map[string]float64
var omega float64

// Driving to the system
func theta0(t float64) (float64, float64) {
	// return math.Sin(omega * t), omega * math.Cos(omega*t)
	return 0.0, 0.0
}

// Function to calculate the derivative of the system (example)
func f(t float64, y *mat.VecDense) *mat.VecDense {
	// Example: Simple harmonic oscillator
	dydt := mat.NewVecDense(y.Len(), nil)
	theta, thetadot := theta0(t)
	fvec := mat.NewVecDense(4, nil)
	fvec.SetVec(2, params["lamA"]*thetadot+params["kA"]*theta)

	dydt.MulVec(Ainv, fvec)
	px := mat.NewVecDense(4, nil)
	px.MulVec(P, y)
	px.ScaleVec(-1.0, px)
	dydt.AddVec(dydt, px)
	return dydt
}

// 4th order Runge-Kutta integration using gonum
func rk4(t float64, y *mat.VecDense, dt float64) *mat.VecDense {
	k1 := f(t, y)

	tempVec := mat.NewVecDense(y.Len(), nil)
	tempVec.ScaleVec(dt/2.0, k1)
	tempVec.AddVec(y, tempVec)
	k2 := f(t+dt/2, tempVec)

	tempVec.ScaleVec(dt/2.0, k2)
	tempVec.AddVec(y, tempVec)
	k3 := f(t+dt/2, tempVec)

	tempVec.ScaleVec(dt, k3)
	tempVec.AddVec(y, tempVec)
	k4 := f(t+dt, tempVec)

	result := mat.NewVecDense(y.Len(), nil)
	tempVec1 := mat.NewVecDense(y.Len(), nil)
	tempVec1.ScaleVec(1.0, k1)
	tempVec2 := mat.NewVecDense(y.Len(), nil)
	tempVec2.ScaleVec(2.0, k2)
	tempVec3 := mat.NewVecDense(y.Len(), nil)
	tempVec3.ScaleVec(2.0, k3)
	tempVec4 := mat.NewVecDense(y.Len(), nil)
	tempVec4.ScaleVec(1.0, k4)

	tempVec1.AddVec(tempVec1, tempVec2)
	tempVec1.AddVec(tempVec1, tempVec3)
	tempVec1.AddVec(tempVec1, tempVec4)
	result.ScaleVec(dt/6.0, tempVec1)
	result.AddVec(y, result)

	return result
}

func explicitEuler(t float64, y *mat.VecDense, dt float64) *mat.VecDense {
	fvec := f(t, y)
	result := mat.NewVecDense(y.Len(), nil)
	result.ScaleVec(dt, fvec)
	result.AddVec(y, result)
	return result
}

func assembleMatrices(params map[string]float64) (*mat.Dense, *mat.Dense) {
	// Assemble the matrices A and b
	A := mat.NewDense(4, 4, []float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, params["I1"], 0, 0, 0, 0, params["I2"]})
	kA := params["kA"]
	kB := params["kB"]
	lamA := params["lamA"]
	lamB := params["lamB"]
	B := mat.NewDense(4, 4, []float64{0, 0, -1, 0, 0, 0, 0, -1, kA + kB, -kB, lamA + lamB, -lamB, -kB, kB, -lamB, lamB})
	return A, B
}

func assembleKMatrix(params map[string]float64) *mat.Dense {
	K := mat.NewDense(2, 2, []float64{params["kA"] + params["kB"], -params["kB"], -params["kB"], params["kB"]})
	return K
}

func assembleMMatrix(params map[string]float64) *mat.Dense {
	M := mat.NewDense(2, 2, []float64{params["I1"], 0, 0, params["I2"]})
	return M
}

func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Excerpt(0))
	fmt.Printf("%v\n", fa)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	// import from file params.yaml
	var err error
	params, err = importParamsFromYAML("params.yaml")
	if err != nil {
		fmt.Println("Error importing parameters:", err)
		return
	}
	fmt.Println("Imported parameters:", params)
	A, B := assembleMatrices(params)
	fmt.Println("Matrix A:")
	matPrint(A)
	fmt.Println("Matrix B:")
	matPrint(B)
	Ainv = mat.NewDense(4, 4, nil)
	err = Ainv.Inverse(A)
	if err != nil {
		fmt.Println("Error calculating inverse of A:", err)
		return
	}
	fmt.Println("Inverse of A:")
	matPrint(Ainv)
	P = mat.NewDense(4, 4, nil)
	P.Mul(Ainv, B)
	omega = params["omega"]

	t := 0.0
	dt := 0.01
	tmax := params["tmax"]
	Nsteps := int(tmax / dt)
	th1 := params["th1"]
	th2 := params["th2"]
	th1dot := params["th1dot"]
	th2dot := params["th2dot"]
	// calculate initial conditions to match steady-state response
	// K := assembleKMatrix(params)
	// M := assembleMMatrix(params)
	// calculate steady-state response
	// thetass := mat.NewVecDense(2, nil)
	// KM := mat.NewDense(2, 2, nil)
	// KM.Scale(-omega*omega, M)
	// KM.Add(KM, K)
	// KM.Inverse(KM)
	// fvec := mat.NewVecDense(2, nil)
	// fvec.SetVec(0, 1.0)
	// thetass.MulVec(KM, fvec)
	// th1 = thetass.AtVec(0)
	// th2 = thetass.AtVec(1)
	// th1dot = 0.0
	// th2dot = 0.0
	// // rescale
	// // th2 = th2 / th1
	// // th1 = 1.0
	// fmt.Println("Initial conditions:", th1, th2, th1dot, th2dot)
	y := mat.NewVecDense(4, []float64{th1, th2, th1dot, th2dot})

	times := mat.NewVecDense(Nsteps, nil)
	history := mat.NewDense(Nsteps, 4, nil)

	for i := 0; i < Nsteps; i++ {
		y = rk4(t, y, dt)
		// y = explicitEuler(t, y, dt)
		t += dt

		times.SetVec(i, t)
		history.SetRow(i, y.RawVector().Data)
	}

	csvFilename := "simulation_results.csv"
	err = exportDataToCSV(csvFilename, times, history)
	if err != nil {
		fmt.Println("Error exporting data:", err)
		return
	}

	err = plotData(times, history)
	if err != nil {
		fmt.Println("Error plotting data:", err)
		return
	}
}
