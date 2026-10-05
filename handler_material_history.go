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

	for _, batch := range iPBatches {
		fmt.Printf("\nIP Batch Lot: %v\n", batch.InProcessBatchLot)

		materials, err := s.db.GetMaterialUsageByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)
		if err != nil {
			return err
		}

		for _, material := range materials {
			fmt.Printf("Material Lot: %v\n", material.MaterialLot)
		}

		ipcResults, err := s.db.GetIPCQCResultsByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)
		if err != nil {
			return err
		}

		fmt.Println("\nIPC QC Results:")
		for _, result := range ipcResults {
			fmt.Printf(
				"In-process Batch Lot: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n\n",
				result.InProcessBatchLot,
				result.TestName,
				result.Replicate,
				result.TestDate,
				result.Result,
			)
		}

		finishedProduct, err := s.db.GetFinishedProductByIPBatch(
			context.Background(),
			batch.InProcessBatchLot,
		)

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error checking FP batch existence: %v", err)
		}

		if err == sql.ErrNoRows {
			fmt.Println("FP Batch: No finished product associated")
		} else {
			fmt.Printf("FP Batch: %v\n", finishedProduct.FinishedProductBatch)

			qcResults, err := s.db.GetQCReleaseResultsByFPBatch(
				context.Background(),
				finishedProduct.FinishedProductBatch,
			)
			if err != nil {
				return err
			}

			fmt.Println("\nFinished Product QC Results:")
			for _, result := range qcResults {
				fmt.Printf(
					"Finished Product Batch: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n\n",
					result.FinishedProductBatch,
					result.TestName,
					result.Replicate,
					result.TestDate,
					result.Result,
				)
			}
		}

		batchStatus, err := getBatchStatus(s, batch.InProcessBatchLot)
		if err != nil {
			return err
		}

		fmt.Println("Batch Status:")
		fmt.Printf("Blend: %s\n", batchStatus.BUResult)
		fmt.Printf("Assay: %s\n", batchStatus.AssayResult)
		fmt.Printf("EmittedDose: %s\n", batchStatus.EDResult)
	}

	return nil
}
