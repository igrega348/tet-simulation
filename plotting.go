package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Plot data
func plotData(times *mat.VecDense, history *mat.Dense) error {
	filename := "simulation_plot_1.png"
	p := plot.New()

	pts := make(plotter.XYs, times.Len())
	for i := 0; i < times.Len(); i++ {
		pts[i].X = times.AtVec(i)
		pts[i].Y = history.At(i, 0) // 1st dof
	}
	line1, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line1.LineStyle.Color = plotutil.Color(0)
	p.Add(line1)
	p.Legend.Add("Position DOF 1", line1)

	for i := 0; i < times.Len(); i++ {
		pts[i].X = times.AtVec(i)
		pts[i].Y = history.At(i, 2) // 3rd dof
	}
	line2, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line2.LineStyle.Color = plotutil.Color(1)
	p.Add(line2)
	p.Legend.Add("Velocity DOF 1", line2)

	p.Title.Text = "Simulation Results"
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "Position/velocity (units)"

	if err := p.Save(4*vg.Inch, 3*vg.Inch, filename); err != nil {
		return err
	}

	fmt.Printf("Plot saved to %s\n", filename)

	// Plot Velocity vs. Time
	p2 := plot.New()
	pts2 := make(plotter.XYs, times.Len())
	for i := 0; i < times.Len(); i++ {
		pts2[i].X = times.AtVec(i)
		pts2[i].Y = history.At(i, 1) // 2nd dof
	}
	line3, err := plotter.NewLine(pts2)
	if err != nil {
		return err
	}
	line3.LineStyle.Color = plotutil.Color(2)
	p2.Add(line3)
	p2.Legend.Add("Position DOF 2", line3)

	for i := 0; i < times.Len(); i++ {
		pts2[i].X = times.AtVec(i)
		pts2[i].Y = history.At(i, 3) // 4th dof
	}
	line4, err := plotter.NewLine(pts2)
	if err != nil {
		return err
	}
	line4.LineStyle.Color = plotutil.Color(3)
	p2.Add(line4)
	p2.Legend.Add("Velocity DOF 2", line4)

	p2.Title.Text = "Simulation Results"
	p2.X.Label.Text = "Time (s)"
	p2.Y.Label.Text = "Position/velocity (units)"

	if err := p2.Save(4*vg.Inch, 3*vg.Inch, "simulation_plot_2.png"); err != nil {
		return err
	}
	fmt.Println("Plot saved to simulation_plot_2.png")

	return nil
}
