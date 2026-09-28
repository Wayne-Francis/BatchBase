package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddBlend(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 3 {
		return fmt.Errorf("addblend requires 3 arguments")
	}
	inProcessBatchLot := cmd.Args[0]
	exists, err := s.db.CheckIPBatchExists(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("IP batch lot does not exist: %v", inProcessBatchLot)
	}
	blendStartDate, err := time.Parse("02/01/2006", cmd.Args[1])
	if err != nil {
		return fmt.Errorf("invalid date format for blend start date: %v", err)
	}
	blendEndDate, err := time.Parse("02/01/2006", cmd.Args[2])
	if err != nil {
		return fmt.Errorf("invalid date format for blend end date: %v", err)
	}
	_, err = s.db.AddBlend(context.Background(), database.AddBlendParams{
		InProcessBatchLot: inProcessBatchLot,
		BlendStartDate:    blendStartDate,
		BlendEndDate:      blendEndDate,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         user.ID,
	})

	if err != nil {
		return err
	}
	fmt.Printf("new blend added:\nIn-process batch lot: %v,\nBlend start date: %v,\nBlend end date: %v,\ncreated at: %v,\nupdated at: %v\ncreated by: %v\n", inProcessBatchLot, blendStartDate, blendEndDate, time.Now(), time.Now(), user.ID)
	return nil
}

func handlerListBlends(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listblends takes no arguments\n")
	}
	blends, err := s.db.GetBlendList(context.Background())
	if err != nil {
		return err
	}
	for _, blend := range blends {
		fmt.Printf("In-process Batch Lot: %v\nBlend Start Date: %v\nBlend End Date: %v\nIn-process Blend Hold End Date: %v\n", blend.InProcessBatchLot, blend.BlendStartDate, blend.BlendEndDate, blend.InProcessBlendHoldEndDate)
	}
	return nil
}

func handlerSearchBlendbyIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchblend takes 1 arguments\n")
	}
	inProcessBatchLot := cmd.Args[0]
	blend, err := s.db.GetBlendByIPBatch(context.Background(), inProcessBatchLot)
	if err != nil {
		return err
	}
	fmt.Printf("In-process Batch Lot: %v\nBlend Start Date: %v\nBlend End Date: %v\nIn-process Blend Hold End Date: %v\n", blend.InProcessBatchLot, blend.BlendStartDate, blend.BlendEndDate, blend.InProcessBlendHoldEndDate)
	return nil
}

func handlerResetBlend(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetblends takes no arguments\n")
	}
	err := s.db.DeleteAllBlend(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Blends have been deleted\n")
	return nil
}
