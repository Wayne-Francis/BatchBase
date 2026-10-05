package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BatchStatus struct {
	IPBatch     string
	FPBatch     string
	BUResult    string
	AssayResult string
	EDResult    string
}

func handlerBatchStatus(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("batchstatus takes one argument\n")
	}

	batchStatus, err := getBatchStatus(s, cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("IP Batch: %s\n", batchStatus.IPBatch)
	fmt.Printf("FP Batch: %s\n", batchStatus.FPBatch)
	fmt.Printf("Blend: %s\n", batchStatus.BUResult)
	fmt.Printf("Assay: %s\n", batchStatus.AssayResult)
	fmt.Printf("EmittedDose: %s\n", batchStatus.EDResult)

	return nil
}

func getBatchStatus(s *state, batchLot string) (BatchStatus, error) {
	var batchStatus BatchStatus

	if strings.HasPrefix(batchLot, "IP-") {
		batchStatus.IPBatch = batchLot

		exists, err := s.db.CheckIPBatchExists(context.Background(), batchLot)
		if err != nil {
			return batchStatus, fmt.Errorf("error checking IP batch existence: %v", err)
		}
		if !exists {
			return batchStatus, fmt.Errorf("IP batch does not exist: %v", batchLot)
		}

		ipcResults, err := s.db.GetIPCQCResultsByIPBatch(
			context.Background(),
			batchLot,
		)
		if err != nil {
			return batchStatus, err
		}

		if len(ipcResults) == 0 {
			batchStatus.BUResult = "INCOMPLETE"
		} else if len(ipcResults) < 10 {
			batchStatus.BUResult = "INCOMPLETE"
		} else if len(ipcResults) == 10 {

			blendSpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"BlendUniformity",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving Blend spec: %v", err)
			}

			blendSpecMin, err := strconv.ParseFloat(blendSpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend minimum specification: %v", err)
			}

			blendSpecMax, err := strconv.ParseFloat(blendSpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend maximum specification: %v", err)
			}

			blendTest := true
			var blendResults []float64
			var blendFailures []string

			for _, result := range ipcResults {
				if result.TestName == "Blend" {
					blendResult, err := strconv.ParseFloat(result.Result, 64)
					if err != nil {
						return batchStatus, fmt.Errorf("invalid Blend result: %v", err)
					}

					blendResults = append(blendResults, blendResult)

					if blendResult < blendSpecMin || blendResult > blendSpecMax {
						blendTest = false
						blendFailures = append(
							blendFailures,
							fmt.Sprintf(
								"Replicate %s outside specification: %.2f",
								result.Replicate,
								blendResult,
							),
						)
					}
				}
			}

			sum := 0.0
			for _, result := range blendResults {
				sum += result
			}

			average := sum / float64(len(blendResults))

			var sumSquaredDifferences float64

			for _, result := range blendResults {
				difference := result - average
				sumSquaredDifferences += difference * difference
			}

			sd := math.Sqrt(
				sumSquaredDifferences / float64(len(blendResults)-1),
			)

			blendSpecMeanMin, err := strconv.ParseFloat(blendSpec.MeanMin.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend minimum specification: %v", err)
			}

			blendSpecMeanMax, err := strconv.ParseFloat(blendSpec.MeanMax.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend maximum specification: %v", err)
			}

			if average < blendSpecMeanMin || average > blendSpecMeanMax {
				blendTest = false
				blendFailures = append(
					blendFailures,
					fmt.Sprintf("Average outside specification: %.2f", average),
				)
			}

			rsd := (sd / average) * 100

			blendSpecRsdLimit, err := strconv.ParseFloat(blendSpec.RsdLimit.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend RSD specification: %v", err)
			}

			if rsd > blendSpecRsdLimit {
				blendTest = false
				blendFailures = append(
					blendFailures,
					fmt.Sprintf("RSD outside specification: %.2f", rsd),
				)
			}

			if blendTest {
				batchStatus.BUResult = "PASS"
			} else {
				batchStatus.BUResult = "FAIL"
				for _, failure := range blendFailures {
					fmt.Printf("  - %s\n", failure)
				}
			}

		} else {
			return batchStatus, fmt.Errorf(
				"error: too many IPC replicates found: %d",
				len(ipcResults),
			)
		}

		finishedProduct, err := s.db.GetFinishedProductByIPBatch(
			context.Background(),
			batchLot,
		)
		if err != nil && err != sql.ErrNoRows {
			return batchStatus, fmt.Errorf("error checking FP batch existence: %v", err)
		}

		if err == sql.ErrNoRows {
			batchStatus.FPBatch = "No finished product associated"
			batchStatus.AssayResult = "INCOMPLETE"
			batchStatus.EDResult = "INCOMPLETE"

			return batchStatus, nil
		}

		batchStatus.FPBatch = finishedProduct.FinishedProductBatch

		qcResults, err := s.db.GetQCReleaseResultsByFPBatch(
			context.Background(),
			finishedProduct.FinishedProductBatch,
		)
		if err != nil {
			return batchStatus, err
		}

		assayCount := 0
		assayResult := 0.0
		emittedDoseCount := 0

		for _, result := range qcResults {
			if result.TestName == "Assay" {
				assayCount++

				assayResult, err = strconv.ParseFloat(result.Result, 64)
				if err != nil {
					return batchStatus, fmt.Errorf("invalid Assay result: %v", err)
				}
			}

			if result.TestName == "EmittedDose" {
				emittedDoseCount++
			}
		}

		if assayCount == 0 {
			batchStatus.AssayResult = "INCOMPLETE"
		} else if assayCount > 1 {
			return batchStatus, fmt.Errorf(
				"error: too many Assay results found: %d",
				assayCount,
			)
		} else {
			assaySpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"Assay",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving Assay spec: %v", err)
			}

			assayMin, err := strconv.ParseFloat(assaySpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Assay minimum specification: %v", err)
			}

			assayMax, err := strconv.ParseFloat(assaySpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Assay maximum specification: %v", err)
			}

			if assayResult < assayMin || assayResult > assayMax {
				batchStatus.AssayResult = "FAIL"
			} else {
				batchStatus.AssayResult = "PASS"
			}
		}

		if emittedDoseCount == 0 {
			batchStatus.EDResult = "INCOMPLETE"
		} else if emittedDoseCount < 20 {
			batchStatus.EDResult = "INCOMPLETE"
		} else if emittedDoseCount > 20 {
			return batchStatus, fmt.Errorf(
				"error: too many EmittedDose results found: %d",
				emittedDoseCount,
			)
		} else {
			EDSpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"EmittedDose",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving EmittedDose spec: %v", err)
			}

			EDSpecMin, err := strconv.ParseFloat(EDSpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose minimum specification: %v", err)
			}

			EDSpecMax, err := strconv.ParseFloat(EDSpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose maximum specification: %v", err)
			}

			EDTest := true
			var EDResults []float64
			var EDFailures []string

			for _, result := range qcResults {
				if result.TestName == "EmittedDose" {
					EDResult, err := strconv.ParseFloat(result.Result, 64)
					if err != nil {
						return batchStatus, fmt.Errorf("invalid EmittedDose result: %v", err)
					}

					EDResults = append(EDResults, EDResult)

					if EDResult < EDSpecMin || EDResult > EDSpecMax {
						EDTest = false
						EDFailures = append(
							EDFailures,
							fmt.Sprintf(
								"Replicate %s outside specification: %.2f",
								result.Replicate,
								EDResult,
							),
						)
					}
				}
			}

			EDsum := 0.0
			for _, result := range EDResults {
				EDsum += result
			}

			EDaverage := EDsum / float64(len(EDResults))

			var EDsumSquaredDifferences float64

			for _, result := range EDResults {
				difference := result - EDaverage
				EDsumSquaredDifferences += difference * difference
			}

			EDsd := math.Sqrt(
				EDsumSquaredDifferences / float64(len(EDResults)-1),
			)

			EDSpecMeanMin, err := strconv.ParseFloat(EDSpec.MeanMin.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose minimum specification: %v", err)
			}

			EDSpecMeanMax, err := strconv.ParseFloat(EDSpec.MeanMax.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose maximum specification: %v", err)
			}

			if EDaverage < EDSpecMeanMin || EDaverage > EDSpecMeanMax {
				EDTest = false
				EDFailures = append(
					EDFailures,
					fmt.Sprintf(
						"Average outside specification: %.2f",
						EDaverage,
					),
				)
			}

			EDrsd := (EDsd / EDaverage) * 100

			EDRsdLimit, err := strconv.ParseFloat(EDSpec.RsdLimit.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose RSD specification: %v", err)
			}

			if EDrsd > EDRsdLimit {
				EDTest = false
				EDFailures = append(
					EDFailures,
					fmt.Sprintf(
						"RSD outside specification: %.2f",
						EDrsd,
					),
				)
			}

			if EDTest {
				batchStatus.EDResult = "PASS"
			} else {
				batchStatus.EDResult = "FAIL"
				for _, failure := range EDFailures {
					fmt.Printf("  - %s\n", failure)
				}
			}
		}

	} else if strings.HasPrefix(batchLot, "FP-") {
		batchStatus.FPBatch = batchLot

		exists, err := s.db.CheckFPBatchExists(context.Background(), batchLot)
		if err != nil {
			return batchStatus, fmt.Errorf("error checking FP batch existence: %v", err)
		}
		if !exists {
			return batchStatus, fmt.Errorf("FP batch does not exist: %v\n", batchLot)
		}

		qcResults, err := s.db.GetQCReleaseResultsByFPBatch(
			context.Background(),
			batchLot,
		)
		if err != nil {
			return batchStatus, err
		}

		assayCount := 0
		assayResult := 0.0
		emittedDoseCount := 0

		for _, result := range qcResults {
			if result.TestName == "Assay" {
				assayCount++

				assayResult, err = strconv.ParseFloat(result.Result, 64)
				if err != nil {
					return batchStatus, fmt.Errorf("invalid Assay result: %v", err)
				}
			}

			if result.TestName == "EmittedDose" {
				emittedDoseCount++
			}
		}

		if assayCount == 0 {
			batchStatus.AssayResult = "INCOMPLETE"
		} else if assayCount > 1 {
			return batchStatus, fmt.Errorf(
				"error: too many Assay results found: %d",
				assayCount,
			)
		} else {
			assaySpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"Assay",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving Assay spec: %v", err)
			}

			assayMin, err := strconv.ParseFloat(assaySpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Assay minimum specification: %v", err)
			}

			assayMax, err := strconv.ParseFloat(assaySpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Assay maximum specification: %v", err)
			}

			if assayResult < assayMin || assayResult > assayMax {
				batchStatus.AssayResult = "FAIL"
			} else {
				batchStatus.AssayResult = "PASS"
			}
		}

		if emittedDoseCount == 0 {
			batchStatus.EDResult = "INCOMPLETE"
		} else if emittedDoseCount < 20 {
			batchStatus.EDResult = "INCOMPLETE"
		} else if emittedDoseCount > 20 {
			return batchStatus, fmt.Errorf(
				"error: too many EmittedDose results found: %d",
				emittedDoseCount,
			)
		} else {
			EDSpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"EmittedDose",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving EmittedDose spec: %v", err)
			}

			EDSpecMin, err := strconv.ParseFloat(EDSpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose minimum specification: %v", err)
			}

			EDSpecMax, err := strconv.ParseFloat(EDSpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose maximum specification: %v", err)
			}

			EDTest := true
			var EDResults []float64
			var EDFailures []string

			for _, result := range qcResults {
				if result.TestName == "EmittedDose" {
					EDResult, err := strconv.ParseFloat(result.Result, 64)
					if err != nil {
						return batchStatus, fmt.Errorf("invalid EmittedDose result: %v", err)
					}

					EDResults = append(EDResults, EDResult)

					if EDResult < EDSpecMin || EDResult > EDSpecMax {
						EDTest = false
						EDFailures = append(
							EDFailures,
							fmt.Sprintf(
								"Replicate %s outside specification: %.2f",
								result.Replicate,
								EDResult,
							),
						)
					}
				}
			}

			EDsum := 0.0
			for _, result := range EDResults {
				EDsum += result
			}

			EDaverage := EDsum / float64(len(EDResults))

			var EDsumSquaredDifferences float64

			for _, result := range EDResults {
				difference := result - EDaverage
				EDsumSquaredDifferences += difference * difference
			}

			EDsd := math.Sqrt(
				EDsumSquaredDifferences / float64(len(EDResults)-1),
			)

			EDSpecMeanMin, err := strconv.ParseFloat(EDSpec.MeanMin.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose minimum specification: %v", err)
			}

			EDSpecMeanMax, err := strconv.ParseFloat(EDSpec.MeanMax.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose maximum specification: %v", err)
			}

			if EDaverage < EDSpecMeanMin || EDaverage > EDSpecMeanMax {
				EDTest = false
				EDFailures = append(
					EDFailures,
					fmt.Sprintf(
						"Average outside specification: %.2f",
						EDaverage,
					),
				)
			}

			EDrsd := (EDsd / EDaverage) * 100

			EDRsdLimit, err := strconv.ParseFloat(EDSpec.RsdLimit.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid EmittedDose RSD specification: %v", err)
			}

			if EDrsd > EDRsdLimit {
				EDTest = false
				EDFailures = append(
					EDFailures,
					fmt.Sprintf(
						"RSD outside specification: %.2f",
						EDrsd,
					),
				)
			}

			if EDTest {
				batchStatus.EDResult = "PASS"
			} else {
				batchStatus.EDResult = "FAIL"
				for _, failure := range EDFailures {
					fmt.Printf("  - %s\n", failure)
				}
			}
		}

		finishedProduct, err := s.db.GetFinishedProductByFinishedProductBatch(
			context.Background(),
			batchLot,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				batchStatus.IPBatch = "No IP batch associated"
				batchStatus.BUResult = "INCOMPLETE"

				fmt.Printf("FP batch has no IP batch associated: %v\n", batchLot)

				return batchStatus, nil
			}

			return batchStatus, fmt.Errorf(
				"error checking IP batch existence: %v",
				err,
			)
		}

		batchStatus.IPBatch = finishedProduct.InProcessBatchLot

		ipcResults, err := s.db.GetIPCQCResultsByIPBatch(
			context.Background(),
			finishedProduct.InProcessBatchLot,
		)
		if err != nil {
			return batchStatus, err
		}

		if len(ipcResults) == 0 {
			batchStatus.BUResult = "INCOMPLETE"
		} else if len(ipcResults) < 10 {
			batchStatus.BUResult = "INCOMPLETE"
		} else if len(ipcResults) == 10 {

			blendSpec, err := s.db.GetSpecsByTestName(
				context.Background(),
				"BlendUniformity",
			)
			if err != nil {
				return batchStatus, fmt.Errorf("error retrieving Blend spec: %v", err)
			}

			blendSpecMin, err := strconv.ParseFloat(blendSpec.MinResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend minimum specification: %v", err)
			}

			blendSpecMax, err := strconv.ParseFloat(blendSpec.MaxResult.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend maximum specification: %v", err)
			}

			blendTest := true
			var blendResults []float64
			var blendFailures []string

			for _, result := range ipcResults {
				if result.TestName == "Blend" {
					blendResult, err := strconv.ParseFloat(result.Result, 64)
					if err != nil {
						return batchStatus, fmt.Errorf("invalid Blend result: %v", err)
					}

					blendResults = append(blendResults, blendResult)

					if blendResult < blendSpecMin || blendResult > blendSpecMax {
						blendTest = false
						blendFailures = append(
							blendFailures,
							fmt.Sprintf(
								"Replicate %s outside specification: %.2f",
								result.Replicate,
								blendResult,
							),
						)
					}
				}
			}

			sum := 0.0
			for _, result := range blendResults {
				sum += result
			}

			average := sum / float64(len(blendResults))

			var sumSquaredDifferences float64

			for _, result := range blendResults {
				difference := result - average
				sumSquaredDifferences += difference * difference
			}

			sd := math.Sqrt(
				sumSquaredDifferences / float64(len(blendResults)-1),
			)

			blendSpecMeanMin, err := strconv.ParseFloat(blendSpec.MeanMin.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend minimum specification: %v", err)
			}

			blendSpecMeanMax, err := strconv.ParseFloat(blendSpec.MeanMax.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend maximum specification: %v", err)
			}

			if average < blendSpecMeanMin || average > blendSpecMeanMax {
				blendTest = false
				blendFailures = append(
					blendFailures,
					fmt.Sprintf("Average outside specification: %.2f", average),
				)
			}

			rsd := (sd / average) * 100

			blendSpecRsdLimit, err := strconv.ParseFloat(blendSpec.RsdLimit.String, 64)
			if err != nil {
				return batchStatus, fmt.Errorf("invalid Blend RSD specification: %v", err)
			}

			if rsd > blendSpecRsdLimit {
				blendTest = false
				blendFailures = append(
					blendFailures,
					fmt.Sprintf("RSD outside specification: %.2f", rsd),
				)
			}

			if blendTest {
				batchStatus.BUResult = "PASS"
			} else {
				batchStatus.BUResult = "FAIL"
				for _, failure := range blendFailures {
					fmt.Printf("  - %s\n", failure)
				}
			}

		} else {
			return batchStatus, fmt.Errorf(
				"error: too many IPC replicates found: %d",
				len(ipcResults),
			)
		}

	} else {
		return batchStatus, fmt.Errorf(
			"batchstatus takes a batch lot as an argument\n",
		)
	}

	return batchStatus, nil
}
