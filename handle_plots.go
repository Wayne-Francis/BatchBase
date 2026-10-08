package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type plotBatch struct {
	name    string
	results []float64
}

func handlerPlotAssay(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("plotassay takes no arguments")
	}

	results, err := s.db.GetAllQCReleaseResults(context.Background())
	if err != nil {
		return err
	}

	batches := []plotBatch{}

	for _, result := range results {
		if result.TestName != "Assay" {
			continue
		}

		value, err := strconv.ParseFloat(result.Result, 64)
		if err != nil {
			return fmt.Errorf("invalid Assay result %q: %v", result.Result, err)
		}

		found := false

		for i := range batches {
			if batches[i].name == result.FinishedProductBatch {
				batches[i].results = append(batches[i].results, value)
				found = true
				break
			}
		}

		if !found {
			batches = append(batches, plotBatch{
				name:    result.FinishedProductBatch,
				results: []float64{value},
			})
		}
	}

	if len(batches) == 0 {
		return fmt.Errorf("no Assay results found")
	}

	spec, err := s.db.GetSpecsByTestName(
		context.Background(),
		"Assay",
	)
	if err != nil {
		return fmt.Errorf("could not get Assay specification: %v", err)
	}

	return createQCPlot(
		"Assay",
		"Finished Product Batch",
		batches,
		float64Value(spec.MinResult),
		float64Value(spec.MaxResult),
	)
}

func handlerPlotBlend(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("plotblend takes no arguments")
	}

	results, err := s.db.GetIPCQCResults(context.Background())
	if err != nil {
		return err
	}

	batches := []plotBatch{}

	for _, result := range results {
		if result.TestName != "Blend Uniformity" {
			continue
		}

		value, err := strconv.ParseFloat(result.Result, 64)
		if err != nil {
			return fmt.Errorf("invalid Blend Uniformity result %q: %v", result.Result, err)
		}

		found := false

		for i := range batches {
			if batches[i].name == result.InProcessBatchLot {
				batches[i].results = append(batches[i].results, value)
				found = true
				break
			}
		}

		if !found {
			batches = append(batches, plotBatch{
				name:    result.InProcessBatchLot,
				results: []float64{value},
			})
		}
	}

	if len(batches) == 0 {
		return fmt.Errorf("no Blend Uniformity results found")
	}

	spec, err := s.db.GetSpecsByTestName(
		context.Background(),
		"Blend Uniformity",
	)
	if err != nil {
		return fmt.Errorf("could not get Blend Uniformity specification: %v", err)
	}

	return createQCPlot(
		"Blend Uniformity",
		"In-Process Batch",
		batches,
		float64Value(spec.MinResult),
		float64Value(spec.MaxResult),
	)
}

func handlerPlotEmittedDose(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("plotemitteddose takes no arguments")
	}

	results, err := s.db.GetAllQCReleaseResults(context.Background())
	if err != nil {
		return err
	}

	batches := []plotBatch{}

	for _, result := range results {
		if result.TestName != "EmittedDose" {
			continue
		}

		value, err := strconv.ParseFloat(result.Result, 64)
		if err != nil {
			return fmt.Errorf("invalid Emitted Dose result %q: %v", result.Result, err)
		}

		found := false

		for i := range batches {
			if batches[i].name == result.FinishedProductBatch {
				batches[i].results = append(batches[i].results, value)
				found = true
				break
			}
		}

		if !found {
			batches = append(batches, plotBatch{
				name:    result.FinishedProductBatch,
				results: []float64{value},
			})
		}
	}

	if len(batches) == 0 {
		return fmt.Errorf("no Emitted Dose results found")
	}

	spec, err := s.db.GetSpecsByTestName(
		context.Background(),
		"EmittedDose",
	)
	if err != nil {
		return fmt.Errorf("could not get Emitted Dose specification: %v", err)
	}

	return createQCPlot(
		"Emitted Dose",
		"Finished Product Batch",
		batches,
		float64Value(spec.MinResult),
		float64Value(spec.MaxResult),
	)
}

func createQCPlot(
	testName string,
	xAxisLabel string,
	batches []plotBatch,
	minSpec float64,
	maxSpec float64,
) error {

	p := plot.New()

	p.Title.Text = testName
	p.X.Label.Text = xAxisLabel
	p.Y.Label.Text = "Result (%)"

	points := make(plotter.XYs, 0)

	labels := make([]string, len(batches))

	for i, batch := range batches {
		labels[i] = batch.name

		for _, result := range batch.results {
			points = append(points, plotter.XY{
				X: float64(i),
				Y: result,
			})
		}
	}

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		return err
	}

	scatter.GlyphStyle.Radius = vg.Points(4)

	p.Add(scatter)

	minLine := plotter.NewFunction(func(x float64) float64 {
		return minSpec
	})
	minLine.Width = vg.Points(1)

	maxLine := plotter.NewFunction(func(x float64) float64 {
		return maxSpec
	})
	maxLine.Width = vg.Points(1)

	p.Add(minLine, maxLine)

	p.NominalX(labels...)

	if err := os.MkdirAll("plots", 0755); err != nil {
		return fmt.Errorf("could not create plots directory: %v", err)
	}

	filename := fmt.Sprintf(
		"plots/%s.png",
		safePlotFilename(testName),
	)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		return fmt.Errorf("could not save plot: %v", err)
	}

	fmt.Printf("Plot created: %s\n", filename)

	return nil
}

func safePlotFilename(testName string) string {
	switch testName {
	case "Blend Uniformity":
		return "blend_uniformity"
	case "Emitted Dose":
		return "emitted_dose"
	case "Assay":
		return "assay"
	default:
		return "qc_plot"
	}
}

func float64Value(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}
