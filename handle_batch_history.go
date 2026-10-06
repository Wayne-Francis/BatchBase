package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type BatchHistorySummaryReport struct {
	IPBatch       string
	FPBatch       string
	Materials     []string
	BlendStart    string
	BlendEnd      string
	FillStart     string
	FillEnd       string
	AssemblyStart string
	AssemblyEnd   string
	BUResult      string
	AssayResult   string
	EDResult      string
}

type BatchHistoryFullReport struct {
	IPBatch       string
	FPBatch       string
	Materials     []string
	BlendStart    string
	BlendEnd      string
	FillStart     string
	FillEnd       string
	AssemblyStart string
	AssemblyEnd   string
	BUResults     []float64
	AssayResults  []float64
	EDResults     []float64
}

func getBatchType(s *state, batchLot string) (string, error) {
	if strings.HasPrefix(batchLot, "IP-") {

		exists, err := s.db.CheckIPBatchExists(context.Background(), batchLot)
		if err != nil {
			return "", fmt.Errorf("error checking IP batch existence: %v", err)
		}
		if !exists {
			return "", fmt.Errorf("IP batch does not exist: %v", batchLot)
		}

		batchType := "IPBatch"
		return batchType, nil

	} else if strings.HasPrefix(batchLot, "FP-") {

		exists, err := s.db.CheckFPBatchExists(context.Background(), batchLot)
		if err != nil {
			return "", fmt.Errorf("error checking FP batch existence: %v", err)
		}
		if !exists {
			return "", fmt.Errorf("FP batch does not exist: %v", batchLot)
		}

		batchType := "FPBatch"
		return batchType, nil
	}

	return "", fmt.Errorf("invalid batch lot: %v", batchLot)
}

func handlerBatchHistorySummary(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("batchhistorysummary takes one argument\n")
	}

	var BatchHistorySummary BatchHistorySummaryReport
	batchLot := cmd.Args[0]

	batchType, err := getBatchType(s, batchLot)
	if err != nil {
		return err
	}

	// Determine IP and FP batch numbers
	if batchType == "IPBatch" {

		BatchHistorySummary.IPBatch = batchLot

		finishedProduct, err := s.db.GetFinishedProductByIPBatch(
			context.Background(),
			batchLot,
		)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error getting finished product: %v", err)
		}

		if err == sql.ErrNoRows {
			BatchHistorySummary.FPBatch = "No finished product associated"
		} else {
			BatchHistorySummary.FPBatch = finishedProduct.FinishedProductBatch
		}

	} else if batchType == "FPBatch" {

		BatchHistorySummary.FPBatch = batchLot

		finishedProduct, err := s.db.GetFinishedProductByFinishedProductBatch(
			context.Background(),
			batchLot,
		)
		if err != nil {
			return fmt.Errorf("error getting finished product: %v", err)
		}

		BatchHistorySummary.IPBatch = finishedProduct.InProcessBatchLot
	}

	// Get Material History for IP Batch
	materialUsages, err := s.db.GetMaterialUsageByIPBatch(
		context.Background(),
		BatchHistorySummary.IPBatch,
	)
	if err != nil {
		return fmt.Errorf("error getting material history: %v", err)
	}

	BatchHistorySummary.Materials = make([]string, len(materialUsages))

	for i, materialUsage := range materialUsages {
		BatchHistorySummary.Materials[i] = materialUsage.MaterialLot
	}

	// Get Blend History for IP Batch
	blend, err := s.db.GetBlendByIPBatch(
		context.Background(),
		BatchHistorySummary.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking blend existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistorySummary.BlendStart = "No blend associated"
		BatchHistorySummary.BlendEnd = "No blend associated"
	} else {
		BatchHistorySummary.BlendStart = blend.BlendStartDate.Format("02/01/2006")
		BatchHistorySummary.BlendEnd = blend.BlendEndDate.Format("02/01/2006")
	}

	// Get Fill History for IP Batch
	fill, err := s.db.GetFillByIPBatch(
		context.Background(),
		BatchHistorySummary.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking fill existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistorySummary.FillStart = "No fill associated"
		BatchHistorySummary.FillEnd = "No fill associated"
	} else {
		BatchHistorySummary.FillStart = fill.FillStartDate.Format("02/01/2006")
		BatchHistorySummary.FillEnd = fill.FillEndDate.Format("02/01/2006")
	}

	// Get Assembly History for IP Batch
	assembly, err := s.db.GetAssemblyByIPBatch(
		context.Background(),
		BatchHistorySummary.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking assembly existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistorySummary.AssemblyStart = "No assembly associated"
		BatchHistorySummary.AssemblyEnd = "No assembly associated"
	} else {
		BatchHistorySummary.AssemblyStart = assembly.AssemblyStartDate.Format("02/01/2006")
		BatchHistorySummary.AssemblyEnd = assembly.AssemblyEndDate.Format("02/01/2006")
	}

	// Get Batch QC Status
	batchStatus, err := getBatchStatus(
		s,
		BatchHistorySummary.IPBatch,
	)
	if err != nil {
		return err
	}

	BatchHistorySummary.BUResult = batchStatus.BUResult
	BatchHistorySummary.AssayResult = batchStatus.AssayResult
	BatchHistorySummary.EDResult = batchStatus.EDResult

	// Output
	fmt.Println("========================================")
	fmt.Println("          BATCH HISTORY SUMMARY")
	fmt.Println("========================================")

	fmt.Printf("\nIP Batch: %s\n", BatchHistorySummary.IPBatch)
	fmt.Printf("FP Batch: %s\n", BatchHistorySummary.FPBatch)

	fmt.Println("\nMaterials:")
	for _, material := range BatchHistorySummary.Materials {
		fmt.Printf("  %s\n", material)
	}

	fmt.Println("\nMANUFACTURING HISTORY")

	fmt.Println("\nBlend:")
	fmt.Printf("  Start: %s\n", BatchHistorySummary.BlendStart)
	fmt.Printf("  End:   %s\n", BatchHistorySummary.BlendEnd)

	fmt.Println("\nFill:")
	fmt.Printf("  Start: %s\n", BatchHistorySummary.FillStart)
	fmt.Printf("  End:   %s\n", BatchHistorySummary.FillEnd)

	fmt.Println("\nAssembly:")
	fmt.Printf("  Start: %s\n", BatchHistorySummary.AssemblyStart)
	fmt.Printf("  End:   %s\n", BatchHistorySummary.AssemblyEnd)

	fmt.Println("\nQC SUMMARY")

	fmt.Printf("Blend Uniformity: %s\n", BatchHistorySummary.BUResult)
	fmt.Printf("Assay:            %s\n", BatchHistorySummary.AssayResult)
	fmt.Printf("Emitted Dose:     %s\n", BatchHistorySummary.EDResult)

	return nil
}

func handlerBatchHistoryFull(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("batchhistoryfull takes one argument\n")
	}

	var BatchHistoryFull BatchHistoryFullReport
	batchLot := cmd.Args[0]

	batchType, err := getBatchType(s, batchLot)
	if err != nil {
		return err
	}

	// Determine IP and FP batch numbers
	if batchType == "IPBatch" {

		BatchHistoryFull.IPBatch = batchLot

		finishedProduct, err := s.db.GetFinishedProductByIPBatch(
			context.Background(),
			batchLot,
		)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error getting finished product: %v", err)
		}

		if err == sql.ErrNoRows {
			BatchHistoryFull.FPBatch = "No finished product associated"
		} else {
			BatchHistoryFull.FPBatch = finishedProduct.FinishedProductBatch
		}

	} else if batchType == "FPBatch" {

		BatchHistoryFull.FPBatch = batchLot

		finishedProduct, err := s.db.GetFinishedProductByFinishedProductBatch(
			context.Background(),
			batchLot,
		)
		if err != nil {
			return fmt.Errorf("error getting finished product: %v", err)
		}

		BatchHistoryFull.IPBatch = finishedProduct.InProcessBatchLot
	}

	// Get Material History for IP Batch
	materialUsages, err := s.db.GetMaterialUsageByIPBatch(
		context.Background(),
		BatchHistoryFull.IPBatch,
	)
	if err != nil {
		return fmt.Errorf("error getting material history: %v", err)
	}

	BatchHistoryFull.Materials = make([]string, len(materialUsages))

	for i, materialUsage := range materialUsages {
		BatchHistoryFull.Materials[i] = materialUsage.MaterialLot
	}

	// Get Blend History for IP Batch
	blend, err := s.db.GetBlendByIPBatch(
		context.Background(),
		BatchHistoryFull.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking blend existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistoryFull.BlendStart = "No blend associated"
		BatchHistoryFull.BlendEnd = "No blend associated"
	} else {
		BatchHistoryFull.BlendStart = blend.BlendStartDate.Format("02/01/2006")
		BatchHistoryFull.BlendEnd = blend.BlendEndDate.Format("02/01/2006")
	}

	// Get Fill History for IP Batch
	fill, err := s.db.GetFillByIPBatch(
		context.Background(),
		BatchHistoryFull.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking fill existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistoryFull.FillStart = "No fill associated"
		BatchHistoryFull.FillEnd = "No fill associated"
	} else {
		BatchHistoryFull.FillStart = fill.FillStartDate.Format("02/01/2006")
		BatchHistoryFull.FillEnd = fill.FillEndDate.Format("02/01/2006")
	}

	// Get Assembly History for IP Batch
	assembly, err := s.db.GetAssemblyByIPBatch(
		context.Background(),
		BatchHistoryFull.IPBatch,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking assembly existence: %v", err)
	}

	if err == sql.ErrNoRows {
		BatchHistoryFull.AssemblyStart = "No assembly associated"
		BatchHistoryFull.AssemblyEnd = "No assembly associated"
	} else {
		BatchHistoryFull.AssemblyStart = assembly.AssemblyStartDate.Format("02/01/2006")
		BatchHistoryFull.AssemblyEnd = assembly.AssemblyEndDate.Format("02/01/2006")
	}

	// Get IPC QC Results
	IPCResults, err := s.db.GetIPCQCResultsByIPBatch(
		context.Background(),
		BatchHistoryFull.IPBatch,
	)
	if err != nil {
		return fmt.Errorf("error getting IPC QC results: %v", err)
	}

	for _, result := range IPCResults {
		if result.TestName == "Blend Uniformity" {

			value, err := strconv.ParseFloat(result.Result, 64)
			if err != nil {
				return fmt.Errorf(
					"error converting Blend Uniformity result: %v",
					err,
				)
			}

			BatchHistoryFull.BUResults = append(
				BatchHistoryFull.BUResults,
				value,
			)
		}
	}

	// Get Finished Product QC Results
	if BatchHistoryFull.FPBatch != "No finished product associated" {

		ReleaseResults, err := s.db.GetQCReleaseResultsByFPBatch(
			context.Background(),
			BatchHistoryFull.FPBatch,
		)
		if err != nil {
			return fmt.Errorf("error getting QC release results: %v", err)
		}

		for _, result := range ReleaseResults {

			value, err := strconv.ParseFloat(result.Result, 64)
			if err != nil {
				return fmt.Errorf("error converting QC result: %v", err)
			}

			if result.TestName == "Assay" {
				BatchHistoryFull.AssayResults = append(
					BatchHistoryFull.AssayResults,
					value,
				)
			}

			if result.TestName == "EmittedDose" {
				BatchHistoryFull.EDResults = append(
					BatchHistoryFull.EDResults,
					value,
				)
			}
		}
	}

	// Output
	fmt.Println("========================================")
	fmt.Println("           FULL BATCH HISTORY")
	fmt.Println("========================================")

	fmt.Printf("\nIP Batch: %s\n", BatchHistoryFull.IPBatch)
	fmt.Printf("FP Batch: %s\n", BatchHistoryFull.FPBatch)

	fmt.Println("\nMaterials:")
	if len(BatchHistoryFull.Materials) == 0 {
		fmt.Println("  No materials associated")
	} else {
		for _, material := range BatchHistoryFull.Materials {
			fmt.Printf("  %s\n", material)
		}
	}

	fmt.Println("\nMANUFACTURING HISTORY")

	fmt.Println("\nBlend:")
	fmt.Printf("  Start: %s\n", BatchHistoryFull.BlendStart)
	fmt.Printf("  End:   %s\n", BatchHistoryFull.BlendEnd)

	fmt.Println("\nFill:")
	fmt.Printf("  Start: %s\n", BatchHistoryFull.FillStart)
	fmt.Printf("  End:   %s\n", BatchHistoryFull.FillEnd)

	fmt.Println("\nAssembly:")
	fmt.Printf("  Start: %s\n", BatchHistoryFull.AssemblyStart)
	fmt.Printf("  End:   %s\n", BatchHistoryFull.AssemblyEnd)

	fmt.Println("\nIPC QC - BLEND UNIFORMITY")

	if len(BatchHistoryFull.BUResults) == 0 {
		fmt.Println("  No Blend Uniformity results")
	} else {
		for i, result := range BatchHistoryFull.BUResults {
			fmt.Printf("  Replicate %d: %.2f\n", i+1, result)
		}
	}

	fmt.Println("\nFINISHED PRODUCT QC")

	fmt.Println("\nAssay:")
	if len(BatchHistoryFull.AssayResults) == 0 {
		fmt.Println("  No Assay results")
	} else {
		for i, result := range BatchHistoryFull.AssayResults {
			fmt.Printf("  Replicate %d: %.2f\n", i+1, result)
		}
	}

	fmt.Println("\nEmitted Dose:")
	if len(BatchHistoryFull.EDResults) == 0 {
		fmt.Println("  No Emitted Dose results")
	} else {
		for i, result := range BatchHistoryFull.EDResults {
			fmt.Printf("  Replicate %d: %.2f\n", i+1, result)
		}
	}

	return nil
}
