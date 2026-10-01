package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddIPCQCResult(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 4 {
		return fmt.Errorf("addipcqcresult requires a batch, test, date and at least one result")
	}

	InProcessBatchLot := cmd.Args[0]

	exists, err := s.db.CheckIPCBatchExists(context.Background(), InProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IPC batch existence: %v", err)
	}

	if !exists {
		return fmt.Errorf("in-process batch lot has not been blended: %v", InProcessBatchLot)
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

	qcResults := addMultipleQCResults(floatResults)

	now := time.Now()

	for _, qcResult := range qcResults {
		_, err = s.db.AddIPCQCResults(context.Background(), database.AddIPCQCResultsParams{
			InProcessBatchLot: InProcessBatchLot,
			CreatedAt:         now,
			UpdatedAt:         now,
			TestName:          testName,
			Replicate:         int32(qcResult.Replicate),
			TestDate:          testDate,
			Result: sql.NullString{
				String: strconv.FormatFloat(qcResult.Result, 'f', 2, 64),
				Valid:  true,
			},
			CreatedBy: user.ID,
		})

		if err != nil {
			return err
		}
	}

	fmt.Printf(
		"IPC QC results added:\nIn-process batch lot: %v\nTest name: %v\nReplicates added: %v\nTest date: %v\nCreated by: %v\n",
		InProcessBatchLot,
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
