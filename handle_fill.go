package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddFill(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 3 {
		return fmt.Errorf("addfill requires 3 arguments")
	}
	inProcessBatchLot := cmd.Args[0]
	exists, err := s.db.CheckIPBatchExistsInBlend(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("IP batch lot has not been blended: %v", inProcessBatchLot)
	}
	fillStartDate, err := time.Parse("02/01/2006", cmd.Args[1])
	if err != nil {
		return fmt.Errorf("invalid date format for fill start date: %v", err)
	}
	fillEndDate, err := time.Parse("02/01/2006", cmd.Args[2])
	if err != nil {
		return fmt.Errorf("invalid date format for fill end date: %v", err)
	}
	_, err = s.db.AddFill(context.Background(), database.AddFillParams{
		InProcessBatchLot: inProcessBatchLot,
		FillStartDate:     fillStartDate,
		FillEndDate:       fillEndDate,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         user.ID,
	})

	if err != nil {
		return err
	}
	fmt.Printf("new fill added:\nIn-process batch lot: %v,\nFill start date: %v,\nFill end date: %v,\ncreated at: %v,\nupdated at: %v\ncreated by: %v\n", inProcessBatchLot, fillStartDate, fillEndDate, time.Now(), time.Now(), user.ID)
	return nil
}

func handlerListFills(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listfills takes no arguments\n")
	}
	fills, err := s.db.GetFillList(context.Background())
	if err != nil {
		return err
	}
	for _, fill := range fills {
		fmt.Printf("In-process Batch Lot: %v\nFill Start Date: %v\nFill End Date: %v\nIn-process Fill Hold End Date: %v\n", fill.InProcessBatchLot, fill.FillStartDate, fill.FillEndDate, fill.InProcessDiscHoldEndDate)
	}
	return nil
}

func handlerSearchFillbyIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchfill takes 1 arguments\n")
	}
	inProcessBatchLot := cmd.Args[0]
	fill, err := s.db.GetFillByIPBatch(context.Background(), inProcessBatchLot)
	if err != nil {
		return err
	}
	fmt.Printf("In-process Batch Lot: %v\nFill Start Date: %v\nFill End Date: %v\nIn-process Fill Hold End Date: %v\n", fill.InProcessBatchLot, fill.FillStartDate, fill.FillEndDate, fill.InProcessDiscHoldEndDate)
	return nil
}

func handlerResetFill(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetfills takes no arguments\n")
	}
	err := s.db.DeleteAllFill(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Fills have been deleted\n")
	return nil
}
