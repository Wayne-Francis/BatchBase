package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func AddQCReleaseResults(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 4 {
		return fmt.Errorf("addqcreleaseresults requires a batch, test, date and at least one result")
	}

	finishedProductBatch := cmd.Args[0]

	exists, err := s.db.CheckFPBatchExists(context.Background(), finishedProductBatch)
	if err != nil {
		return fmt.Errorf("error checking FP batch existence: %v", err)
	}

	if !exists {
		return fmt.Errorf("finished product batch does not exist or has not been assembled: %v", finishedProductBatch)
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

	existingResults, err := s.db.CountQCReleaseResults(
		context.Background(),
		database.CountQCReleaseResultsParams{
			FinishedProductBatch: finishedProductBatch,
			TestName:             testName,
		},
	)
	if err != nil {
		return fmt.Errorf("error checking existing QC release results: %v", err)
	}

	qcResults := addMultipleQCResults(floatResults, int(existingResults)+1)

	now := time.Now()

	for _, qcResult := range qcResults {
		_, err = s.db.AddQCReleaseResults(context.Background(), database.AddQCReleaseResultsParams{
			FinishedProductBatch: finishedProductBatch,
			CreatedAt:            now,
			UpdatedAt:            now,
			TestName:             testName,
			Replicate:            int32(qcResult.Replicate),
			TestDate:             testDate,
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
		"QC release results added:\nFinished product batch: %v\nTest name: %v\nReplicates added: %v\nTest date: %v\nCreated by: %v\n",
		finishedProductBatch,
		testName,
		len(qcResults),
		testDate,
		user.ID,
	)

	return nil

}

func handlerListQCReleaseResults(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listqcreleaseresults takes no arguments\n")
	}

	results, err := s.db.GetAllQCReleaseResults(context.Background())
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Printf(
			"Finished Product Batch: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.FinishedProductBatch,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func GetQCReleaseResultsByFPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchqcreleaseresultsbyfpbatch takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]

	qcResults, err := s.db.GetQCReleaseResultsByFPBatch(
		context.Background(),
		finishedProductBatch,
	)
	if err != nil {
		return err
	}

	for _, result := range qcResults {
		fmt.Printf(
			"Finished Product Batch: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.FinishedProductBatch,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func GetQCReleaseResultsByIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchqcreleaseresultsbyipbatch takes 1 argument\n")
	}

	inProcessBatchLot := cmd.Args[0]

	qcResults, err := s.db.GetQCReleaseResultsByIPBatch(
		context.Background(),
		inProcessBatchLot,
	)
	if err != nil {
		return err
	}

	for _, result := range qcResults {
		fmt.Printf(
			"Finished Product Batch: %v\nTest Name: %v\nReplicate: %v\nTest Date: %v\nResult: %v\n",
			result.FinishedProductBatch,
			result.TestName,
			result.Replicate,
			result.TestDate,
			result.Result,
		)
	}

	return nil
}

func handlerResetQCReleaseResult(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetqcreleaseresult takes no arguments\n")
	}

	err := s.db.DeleteAllQCReleaseResults(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("QC release results have been reset\n")
	return nil
}

func handlerDeleteQCReleaseResult(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("deleteqcreleaseresult takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]

	err := s.db.DeleteQCReleaseResultsForFPBatch(context.Background(), finishedProductBatch)
	if err != nil {
		return err
	}

	fmt.Printf("QC release results have been deleted\n")
	return nil
}
