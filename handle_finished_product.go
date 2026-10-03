package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddfinishedproduct(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 4 {
		return fmt.Errorf("addfinishedproduct requires 4 arguments")
	}
	inProcessBatchLot := cmd.Args[0]

	ipExists, err := s.db.CheckIPBatchExists(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch existence: %v", err)
	}
	if !ipExists {
		return fmt.Errorf("in-process batch lot does not exist in batch material usage: %v", inProcessBatchLot)
	}

	blendExists, err := s.db.CheckIPBatchExistsInBlend(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch existence in blend: %v", err)
	}
	if !blendExists {
		return fmt.Errorf("IP batch lot has not been blended: %v", inProcessBatchLot)
	}

	fillExists, err := s.db.CheckIPBatchExistsInFill(context.Background(), inProcessBatchLot)
	if err != nil {
		return fmt.Errorf("error checking IP batch existence in fill: %v", err)
	}
	if !fillExists {
		return fmt.Errorf("IP batch lot has not been filled: %v", inProcessBatchLot)
	}
	_, err = s.db.GetFinishedProductByIPBatch(context.Background(), inProcessBatchLot)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking FP batch existence: %v", err)
	}
	if err == nil {
		return fmt.Errorf("finished product for IP batch lot already exists: %v", inProcessBatchLot)
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

func DeleteLotFromFinishedProduct(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("delete lot from finished product takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]
	assemblyExists, err := s.db.CheckFPBatchExistsInAssembly(context.Background(), finishedProductBatch)
	if err != nil {
		return fmt.Errorf("error checking assembly results for FP batch: %v", err)
	}

	if assemblyExists {
		return fmt.Errorf("cannot delete finished product batch: assembly results exist for finished product batch: %v", finishedProductBatch)
	}
	qcreleaseExists, err := s.db.CheckFPBatchExistsInQCRelease(context.Background(), finishedProductBatch)
	if err != nil {
		return fmt.Errorf("error checking QC release results for FP batch: %v", err)
	}

	if qcreleaseExists {
		return fmt.Errorf("cannot delete finished product batch: QC release results exist for finished product batch: %v", finishedProductBatch)
	}
	err = s.db.DeleteLotFromFinishedProducts(context.Background(), finishedProductBatch)
	if err != nil {
		return err
	}

	fmt.Printf("Finished product batch has been deleted from finished products\n")
	return nil
}
