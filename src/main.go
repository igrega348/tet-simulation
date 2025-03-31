// Author: Ivan Grega
// License: MIT
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/gonum/mat"
)

var Ainv *mat.Dense
var P *mat.Dense
var params map[string]float64
var omega float64

// Driving to the system
func theta0(t float64) (float64, float64) {
	return math.Sin(omega * t), omega * math.Cos(omega*t)
	// return 0.0, 0.0
}

// Function to calculate the derivative of the system
func f(t float64, y *mat.VecDense) *mat.VecDense {
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

// 4th order Runge-Kutta integration
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

func DenseToSymDense(A *mat.Dense) *mat.SymDense {
	// Convert a Dense matrix to a SymDense matrix
	rows, cols := A.Dims()
	if rows != cols {
		log.Error().Msg("Matrix is not square")
		return nil
	}
	symA := mat.NewSymDense(rows, nil)
	for i := 0; i < rows; i++ {
		for j := i; j < cols; j++ {
			symA.SetSym(i, j, A.At(i, j))
		}
	}
	return symA
}

func get_dt(params map[string]float64) float64 {
	// calculate analytical estimates of natural frequencies
	K := assembleKMatrix(params)
	M := assembleMMatrix(params)
	KM := mat.NewDense(2, 2, nil)
	KM.Inverse(M)
	KM.Mul(KM, K)
	KMsym := DenseToSymDense(KM)
	// calculate eigenvalues = omega^2
	var eigen mat.EigenSym
	eigen.Factorize(KMsym, false)
	// calculate natural frequencies. f=omega_n/(2*pi)
	omega_n := []float64{0, 0}
	eigen.Values(omega_n)
	// square root of eigenvalues
	fmax := 0.0
	for i := 0; i < 2; i++ {
		omega_n[i] = math.Sqrt(math.Abs(omega_n[i]))
		if omega_n[i] > fmax {
			fmax = omega_n[i]
		}
	}
	log.Info().Msgf("Natural frequencies omega_n: %v", omega_n)
	// set dt to 1/10 of the period of the highest frequency or 1/10 of the period of driving
	fmax = max(fmax, params["omega"]/(2.0*math.Pi))
	dt := 1.0 / (10.0 * fmax)
	return dt
}

func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Excerpt(0))
	log.Info().Msgf("\n%v\n", fa)
}

func simulate(
	params_file string,
) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// import from file params.yaml
	var err error
	params, err = importParamsFromYAML(params_file)
	if err != nil {
		fmt.Println("Error importing parameters:", err)
		return
	}
	log.Info().Msgf("Imported parameters: %v", params)
	A, B := assembleMatrices(params)
	log.Info().Msg("Matrix A:")
	matPrint(A)
	log.Info().Msg("Matrix B:")
	matPrint(B)
	Ainv = mat.NewDense(4, 4, nil)
	err = Ainv.Inverse(A)
	if err != nil {
		fmt.Println("Error calculating inverse of A:", err)
		return
	}
	log.Info().Msg("Inverse of A:")
	matPrint(Ainv)
	P = mat.NewDense(4, 4, nil)
	P.Mul(Ainv, B)
	omega = float64(params["omega"])

	t := 0.0
	dt := get_dt(params)
	tmax := float64(params["tmax"])
	Nsteps := int(tmax / dt)
	th1 := float64(params["th1"])
	th2 := float64(params["th2"])
	th1dot := float64(params["th1dot"])
	th2dot := float64(params["th2dot"])
	log.Info().Msgf("Time step dt=%.2g", dt)
	log.Info().Msgf("Initial conditions: %.2g %.2f %.2g %.2g", th1, th2, th1dot, th2dot)
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
	varNames := []string{"theta1", "theta2", "theta1dot", "theta2dot"}
	err = exportDataToCSV(csvFilename, times, history, varNames...)
	if err != nil {
		log.Error().Msgf("Error exporting data to CSV: %v", err)
		return
	}

	err = plotData(times, history)
	if err != nil {
		log.Error().Msgf("Error plotting data: %v", err)
		return
	}
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "v",
				Usage: "Enable verbose logging",
			},
			&cli.StringFlag{
				Name:  "params",
				Usage: "Path to the parameters file",
			},
		},
		Action: func(c *cli.Context) error {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			if c.Bool("v") {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
				log.Info().Msg("Verbose logging enabled")
			} else {
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			}
			simulate(c.String("params"))
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}
