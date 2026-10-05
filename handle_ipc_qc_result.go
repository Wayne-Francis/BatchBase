package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddIPCQCResult(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 4 {
		return fmt.Errorf("addipcqcresult requires a batch, test, date and at least one result")
	}

	inProcessBatchLot := cmd.Args[0]

	exists, err := s.db.CheckIPCBatchExists(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IPC batch existence: %v", err)
	}

	if !exists {
		return fmt.Errorf("in-process batch lot does not exist or has not been blended: %v", inProcessBatchLot)
	}

	testName := cmd.Args[1]

	testDate, err := time.Parse("02/01/2006", cmd.Args[2])
	if err != nil {
		return fmt.Errorf("invalid date format for test date: %v", err)
	}

	results := cmd.Args[3:]
	floatResults := []float64{}

	for _, result := range results {
		value, err := strconv.ParseFloat(result, 64)
		if err != nil {
			return fmt.Errorf("invalid QC result %q: %v", result, err)
		}

		floatResults = append(floatResults, value)
	}

	existingResults, err := s.db.CountIPCQCResults(
		context.Background(),
		database.CountIPCQCResultsParams{
			InProcessBatchLot: inProcessBatchLot,
			TestName:          testName,
		},
	)
	if err != nil {
		return fmt.Errorf("error checking existing IPC QC results: %v", err)
	}

	qcResults := addMultipleQCResults(floatResults, int(existingResults)+1)

	now := time.Now()

	for _, qcResult := range qcResults {
		_, err = s.db.AddIPCQCResults(context.Background(), database.AddIPCQCResultsParams{
			InProcessBatchLot: inProcessBatchLot,
			CreatedAt:         now,
			UpdatedAt:         now,
			TestName:          testName,
			Replicate:         int32(qcResult.Replicate),
			TestDate:          testDate,
			Result:            strconv.FormatFloat(qcResult.Result, 'f', 2, 64),
			CreatedBy:         user.ID,
		})

		if err != nil {
			return err
		}
	}

	fmt.Printf(
		"IPC QC results added:\nIn-process batch lot: %v\nTest name: %v\nReplicates added: %v\nTest date: %v\nCreated by: %v\n",
		inProcessBatchLot,
		testName,
		len(qcResults),
		testDate,
		user.ID,
	)

	return nil
}

func handlerListIPCQCResults(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listipcqcresults takes no arguments\n")
	}

	results, err := s.db.GetIPCQCResults(context.Background())
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Printf(
			"In-process Batch Lot: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.InProcessBatchLot,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func GetIPCQCResultsByIPCBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchipcqcresultsbyipcbatch takes 1 argument\n")
	}

	inProcessBatchLot := cmd.Args[0]

	ipcResults, err := s.db.GetIPCQCResultsByIPBatch(
		context.Background(),
		inProcessBatchLot,
	)
	if err != nil {
		return err
	}

	for _, result := range ipcResults {
		fmt.Printf(
			"In-process Batch Lot: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.InProcessBatchLot,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func GetIPCQCResultsByFPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchipcqcresultsbyfpbatch takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]

	ipcResults, err := s.db.GetIPCQCResultsByFPBatch(
		context.Background(),
		finishedProductBatch,
	)
	if err != nil {
		return err
	}

	for _, result := range ipcResults {
		fmt.Printf(
			"In-process Batch Lot: %v\nFinished Product Batch: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.InProcessBatchLot,
			result.FinishedProductBatch,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func handlerResetIPCQCResult(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetipcqcresult takes no arguments\n")
	}

	err := s.db.DeleteAllIPCQCResults(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("IPC QC results have been reset\n")
	return nil
}

func handlerDeleteIPCQCResult(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("deleteipcqcresult takes 1 argument\n")
	}

	inProcessBatchLot := cmd.Args[0]

	err := s.db.DeleteIPCQCResultsForIPBatch(context.Background(), inProcessBatchLot)
	if err != nil {
		return err
	}

	fmt.Printf("IPC QC results have been deleted\n")
	return nil
}
