package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Function to calculate the derivative of the system (example)
func f(t float64, y *mat.VecDense) *mat.VecDense {
	// Example: Simple harmonic oscillator
	dydt := mat.NewVecDense(y.Len(), nil)
	dydt.SetVec(0, y.AtVec(1))  // dy/dt = v
	dydt.SetVec(1, -y.AtVec(0)) // dv/dt = -y  (assuming mass=1, spring constant=1)
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

// Export data to CSV file
func exportDataToCSV(filename string, times []float64, positions []float64, velocities []float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err // Return the error instead of panicking
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Time", "Position", "Velocity"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i := 0; i < len(times); i++ {
		row := []string{
			strconv.FormatFloat(times[i], 'G', 6, 64),
			strconv.FormatFloat(positions[i], 'G', 6, 64),
			strconv.FormatFloat(velocities[i], 'G', 6, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	fmt.Printf("Simulation results written to %s\n", filename)

	return nil // Indicate success
}

// Plot data
func plotData(filename string, times []float64, positions []float64, velocities []float64) error {
	p := plot.New()

	pts := make(plotter.XYs, len(times))
	for i := range times {
		pts[i].X = times[i]
		pts[i].Y = positions[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	p.Add(line)

	p.Title.Text = "Simulation Results: Position vs. Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Position"

	if err := p.Save(4*vg.Inch, 3*vg.Inch, filename); err != nil {
		return err
	}

	fmt.Printf("Plot saved to %s\n", filename)

	// Plot Velocity vs. Time
	p2 := plot.New()
	pts2 := make(plotter.XYs, len(times))
	for i := range times {
		pts2[i].X = times[i]
		pts2[i].Y = velocities[i]
	}
	line2, err := plotter.NewLine(pts2)
	if err != nil {
		return err
	}
	p2.Add(line2)
	p2.Title.Text = "Simulation Results: Velocity vs. Time"
	p2.X.Label.Text = "Time"
	p2.Y.Label.Text = "Velocity"
	if err := p2.Save(4*vg.Inch, 3*vg.Inch, "simulation_plot_velocity.png"); err != nil {
		return err
	}
	fmt.Println("Velocity plot saved to simulation_plot_velocity.png")

	return nil
}

func main() {
	t := 0.0
	dt := 0.01
	y := mat.NewVecDense(2, []float64{1.0, 0.0})

	times := []float64{}
	positions := []float64{}
	velocities := []float64{}

	for i := 0; i < 100; i++ {
		y = rk4(t, y, dt)
		t += dt

		times = append(times, t)
		positions = append(positions, y.AtVec(0))
		velocities = append(velocities, y.AtVec(1))
	}

	csvFilename := "simulation_results.csv"
	err := exportDataToCSV(csvFilename, times, positions, velocities)
	if err != nil {
		fmt.Println("Error exporting data:", err)
		return
	}

	plotFilename := "simulation_plot.png" // Name for the position plot
	err = plotData(plotFilename, times, positions, velocities)
	if err != nil {
		fmt.Println("Error plotting data:", err)
		return
	}
}
