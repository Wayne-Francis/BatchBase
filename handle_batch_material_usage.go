package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddMaterialUsage(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addmaterialusage requires 2 arguments")
	}
	inProcessBatchLot := cmd.Args[0]
	materialLot := cmd.Args[1]
	_, err := s.db.AddMaterialUsage(context.Background(), database.AddMaterialUsageParams{
		InProcessBatchLot: inProcessBatchLot,
		MaterialLot:       materialLot,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         user.ID,
	})

	if err != nil {
		return err
	}
	fmt.Printf("new material usage added:\nIn-process batch lot: %v,\nMaterial lot: %v,\ncreated at: %v,\nupdated at: %v\ncreated by: %v\n", inProcessBatchLot, materialLot, time.Now(), time.Now(), user.ID)
	return nil
}

func handlerListMaterialUsage(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listmaterialusage takes no arguments\n")
	}
	materialusages, err := s.db.GetMaterialUsageList(context.Background())
	if err != nil {
		return err
	}
	for _, materialusage := range materialusages {
		fmt.Printf("In-process Batch Lot: %v\nMaterial Lot: %v\nCreated At: %v\nUpdated At: %v\nCreated By: %v\n", materialusage.InProcessBatchLot, materialusage.MaterialLot, materialusage.CreatedAt, materialusage.UpdatedAt, materialusage.CreatedBy)
	}
	return nil
}

func handlerSearchMaterialUsageByIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchipbatch takes 1 arguments\n")
	}
	inProcessBatchLot := cmd.Args[0]
	materialUsages, err := s.db.GetMaterialUsageByIPBatch(context.Background(), inProcessBatchLot)
	if err != nil {
		return err
	}
	for _, materialUsage := range materialUsages {
		fmt.Printf("Material Lot: %v\n", materialUsage.MaterialLot)
	}

	return nil
}

func handlerSearchMaterialUsageByMaterialLot(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchmaterialusagebymateriallot takes 1 arguments\n")
	}
	materialLot := cmd.Args[0]
	IPBatches, err := s.db.GetIPBatchesByMaterialLot(context.Background(), materialLot)
	if err != nil {
		return err
	}
	for _, ipBatch := range IPBatches {
		fmt.Printf("IPBatch: %v\n", ipBatch.InProcessBatchLot)
	}
	return nil
}

func handlerResetMaterialUsage(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetmaterialsusage takes no arguments\n")
	}
	err := s.db.DeleteAllMaterialUsage(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Material usage have been deleted\n")
	return nil
}
