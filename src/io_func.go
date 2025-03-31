package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/mat"
	"gopkg.in/yaml.v3"
)

// Export data to CSV file
func exportDataToCSV(filename string, times *mat.VecDense, history *mat.Dense, variableNames ...string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err // Return the error instead of panicking
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	rows, cols := history.Dims()

	// Use provided variable names or default to y0, y1, ...
	header := []string{"Time"}
	if len(variableNames) == cols {
		header = append(header, variableNames...)
	} else {
		for i := 0; i < cols; i++ {
			header = append(header, fmt.Sprintf("y%d", i))
		}
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i := 0; i < rows; i++ {
		row := make([]string, cols+1)
		row[0] = strconv.FormatFloat(times.AtVec(i), 'G', 6, 64)
		for j := 0; j < cols; j++ {
			row[j+1] = strconv.FormatFloat(history.At(i, j), 'G', 6, 64)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	log.Info().Msgf("Simulation results written to %s", filename)

	return nil // Indicate success
}

// import parameters from yaml into dictionary
func importParamsFromYAML(filename string) (map[string]float64, error) {
	// parameters to be imported are
	// I1, I2, kA, kB, lamA, lamB
	// Open the YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal().Msgf("Error reading file: %v", err)
	}
	// unmarshal
	out := map[string]interface{}{}
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		log.Fatal().Msgf("Error unmarshalling YAML: %v", err)
	}
	// convert to float64
	params := map[string]float64{}
	for k, v := range out {
		params[k], err = toFloat64(v)
		if err != nil {
			return params, fmt.Errorf("error converting %s to float64: %v", k, err)
		}
	}
	return params, err
}

func toFloat64(data interface{}) (float64, error) {
	switch t := data.(type) {
	case int:
		return float64(t), nil
	case float64:
		return t, nil
	default:
		return 0.0, fmt.Errorf("data is not an integer or float64")
	}
}
