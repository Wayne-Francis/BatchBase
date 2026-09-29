package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddfinishedproduct(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 4 {
		return fmt.Errorf("addfinishedproduct requires 4 arguments")
	}
	inProcessBatchLot := cmd.Args[0]
	exists, err := s.db.CheckIPBatchExistsInFill(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch fill: %v", err)
	}
	if !exists {
		if !exists {
			return fmt.Errorf("IP batch lot has not been filled: %v", inProcessBatchLot)
		}
	}
	finishedProductBatch := cmd.Args[1]
	component1Batch := cmd.Args[2]
	component2Batch := cmd.Args[3]
	_, err = s.db.AddFinishedProduct(context.Background(), database.AddFinishedProductParams{
		InProcessBatchLot:    inProcessBatchLot,
		FinishedProductBatch: finishedProductBatch,
		Component1Batch:      component1Batch,
		Component2Batch:      component2Batch,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("new finished product added:\nIn-process batch lot: %v,\nFinished product batch: %v,\nComponent 1 batch: %v,\nComponent 2 batch: %v,\ncreated at: %v,\nupdated at: %v\ncreated by: %v\n", inProcessBatchLot, finishedProductBatch, component1Batch, component2Batch, time.Now(), time.Now(), user.ID)
	return nil
}

func handlerListFinishedProducts(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listfinishedproducts takes no arguments\n")
	}
	finishedProducts, err := s.db.GetFinishedProduct(context.Background())
	if err != nil {
		return err
	}
	for _, finishedProduct := range finishedProducts {
		fmt.Printf("In-process Batch Lot: %v\nFinished Product Batch: %v\nComponent 1 Batch: %v\nComponent 2 Batch: %v\n", finishedProduct.InProcessBatchLot, finishedProduct.FinishedProductBatch, finishedProduct.Component1Batch, finishedProduct.Component2Batch)
	}
	return nil
}

func handlerSearchFinishedProductByIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchfinishedproducts takes 1 arguments\n")
	}
	inProcessBatchLot := cmd.Args[0]
	finishedProduct, err := s.db.GetFinishedProductByIPBatch(context.Background(), inProcessBatchLot)
	if err != nil {
		return err
	}
	fmt.Printf("In-process Batch Lot: %v\nFinished Product Batch: %v\nComponent 1 Batch: %v\nComponent 2 Batch: %v\n", finishedProduct.InProcessBatchLot, finishedProduct.FinishedProductBatch, finishedProduct.Component1Batch, finishedProduct.Component2Batch)
	return nil
}

func handlerResetFinishedProducts(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetfinishedproducts takes no arguments\n")
	}
	err := s.db.DeleteAllFinishedProducts(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Finished products have been deleted\n")
	return nil
}
