package main

import (
	"context"
	"database/sql"
	"fmt"
)

func handlerMaterialHistory(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("materialhistory takes one argument\n")
	}

	materialLot := cmd.Args[0]

	iPBatches, err := s.db.GetIPBatchesByMaterialLot(context.Background(), materialLot)
	if err != nil {
		return err
	}

	fmt.Println("========================================")
	fmt.Println("           MATERIAL HISTORY")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("Material Lot: %v\n", materialLot)

	// Material history summary
	fmt.Println()
	fmt.Printf("Used In: %d IP Batches\n", len(iPBatches))

	if len(iPBatches) > 0 {
		fmt.Println()
		fmt.Printf("%-15s %-15s %s\n", "IP Batch", "FP Batch", "Status")
		fmt.Println("----------------------------------------")

		for _, batch := range iPBatches {
			finishedProduct, err := s.db.GetFinishedProductByIPBatch(
				context.Background(),
				batch.InProcessBatchLot,
			)

			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("error checking FP batch existence: %v", err)
			}

			fpBatch := "None"
			if err == nil {
				fpBatch = finishedProduct.FinishedProductBatch
			}

			batchStatus, err := getBatchStatus(s, batch.InProcessBatchLot)
			if err != nil {
				return err
			}

			status := "INCOMPLETE"

			if batchStatus.BUResult == "PASS" &&
				batchStatus.AssayResult == "PASS" &&
				batchStatus.EDResult == "PASS" {
				status = "PASS"
			}

			fmt.Printf(
				"%-15s %-15s %s\n",
				batch.InProcessBatchLot,
				fpBatch,
				status,
			)
		}
	} else {
		fmt.Println("No batches found")
	}

	for _, batch := range iPBatches {
		fmt.Println()
		fmt.Println("BATCH")
		fmt.Println("----------------------------------------")
		fmt.Printf("IP Batch: %v\n", batch.InProcessBatchLot)

		materials, err := s.db.GetMaterialUsageByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("MATERIALS USED")
		fmt.Println("----------------------------------------")

		for _, material := range materials {
			fmt.Printf("%v\n", material.MaterialLot)
		}

		ipcResults, err := s.db.GetIPCQCResultsByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("IPC QC - BLEND UNIFORMITY")
		fmt.Println("----------------------------------------")

		if len(ipcResults) == 0 {
			fmt.Println("No IPC QC results")
		} else {
			fmt.Printf("%-13s %s\n", "Replicate", "Result")

			for _, result := range ipcResults {
				fmt.Printf(
					"%-13v %v\n",
					result.Replicate,
					result.Result,
				)
			}
		}

		finishedProduct, err := s.db.GetFinishedProductByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)

		fmt.Println()
		fmt.Println("FINISHED PRODUCT")
		fmt.Println("----------------------------------------")

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error checking FP batch existence: %v", err)
		}

		if err == sql.ErrNoRows {
			fmt.Println("No finished product associated")
		} else {
			fmt.Printf("FP Batch: %v\n", finishedProduct.FinishedProductBatch)

			qcResults, err := s.db.GetQCReleaseResultsByFPBatch(
				context.Background(),
				finishedProduct.FinishedProductBatch,
			)
			if err != nil {
				return err
			}

			fmt.Println()
			fmt.Println("FINISHED PRODUCT QC")
			fmt.Println("----------------------------------------")

			if len(qcResults) == 0 {
				fmt.Println("No finished product QC results")
			} else {
				fmt.Printf("%-18s %-13s %s\n", "Test", "Replicate", "Result")

				for _, result := range qcResults {
					testName := result.TestName

					if testName == "EmittedDose" {
						testName = "Emitted Dose"
					}

					fmt.Printf(
						"%-18s %-13v %v\n",
						testName,
						result.Replicate,
						result.Result,
					)
				}
			}
		}

		batchStatus, err := getBatchStatus(s, batch.InProcessBatchLot)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("BATCH STATUS")
		fmt.Println("----------------------------------------")
		fmt.Printf("%-18s %s\n", "Blend Uniformity:", batchStatus.BUResult)
		fmt.Printf("%-18s %s\n", "Assay:", batchStatus.AssayResult)
		fmt.Printf("%-18s %s\n", "Emitted Dose:", batchStatus.EDResult)
	}

	return nil
}
