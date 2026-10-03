package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddAssembly(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 3 {
		return fmt.Errorf("addassembly requires 3 arguments")
	}

	finishedProductBatch := cmd.Args[0]

	exists, err := s.db.CheckFPBatchExistsInFinishedProduct(context.Background(), finishedProductBatch)
	if err != nil {
		return fmt.Errorf("error checking FP batch existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("finished product batch does not exist: %v", finishedProductBatch)
	}

	assemblyStartDate, err := time.Parse("02/01/2006", cmd.Args[1])
	if err != nil {
		return fmt.Errorf("invalid date format for assembly start date: %v", err)
	}

	assemblyEndDate, err := time.Parse("02/01/2006", cmd.Args[2])
	if err != nil {
		return fmt.Errorf("invalid date format for assembly end date: %v", err)
	}

	now := time.Now()

	_, err = s.db.AddAssembly(context.Background(), database.AddAssemblyParams{
		FinishedProductBatch: finishedProductBatch,
		CreatedAt:            now,
		UpdatedAt:            now,
		AssemblyStartDate:    assemblyStartDate,
		AssemblyEndDate:      assemblyEndDate,
		CreatedBy:            user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf(
		"New assembly added:\nFinished product batch: %v\nAssembly start date: %v\nAssembly end date: %v\nCreated at: %v\nUpdated at: %v\nCreated by: %v\n",
		finishedProductBatch,
		assemblyStartDate,
		assemblyEndDate,
		now,
		now,
		user.ID,
	)

	return nil
}

func handlerListAssembly(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listassembly takes no arguments\n")
	}

	assemblies, err := s.db.GetAssembly(context.Background())
	if err != nil {
		return err
	}

	for _, assembly := range assemblies {
		fmt.Printf(
			"Finished Product Batch: %v\nAssembly Start Date: %v\nAssembly End Date: %v\nFinished Product Expiry: %v\n",
			assembly.FinishedProductBatch,
			assembly.AssemblyStartDate,
			assembly.AssemblyEndDate,
			assembly.FinishedProductExpiry,
		)
	}

	return nil
}

func handlerSearchAssemblyByFPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchassemblybyfpbatch takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]

	assembly, err := s.db.GetAssemblyByFinishedProductBatch(
		context.Background(),
		finishedProductBatch,
	)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Finished Product Batch: %v\nAssembly Start Date: %v\nAssembly End Date: %v\nFinished Product Expiry: %v\n",
		assembly.FinishedProductBatch,
		assembly.AssemblyStartDate,
		assembly.AssemblyEndDate,
		assembly.FinishedProductExpiry,
	)

	return nil
}

func handlerSearchAssemblyByIPBatch(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchassemblybyipbatch takes 1 argument\n")
	}

	inProcessBatchLot := cmd.Args[0]

	assembly, err := s.db.GetAssemblyByIPBatch(
		context.Background(),
		inProcessBatchLot,
	)
	if err != nil {
		return err
	}

	fmt.Printf(
		"In-process Batch Lot: %v\nFinished Product Batch: %v\nAssembly Start Date: %v\nAssembly End Date: %v\nFinished Product Expiry: %v\n",
		assembly.InProcessBatchLot,
		assembly.FinishedProductBatch,
		assembly.AssemblyStartDate,
		assembly.AssemblyEndDate,
		assembly.FinishedProductExpiry,
	)

	return nil
}

func handlerResetAssembly(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetassembly takes no arguments\n")
	}

	err := s.db.DeleteAllAssembly(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Assembly has been deleted\n")
	return nil
}

func handlerDeleteFPBatchFromAssembly(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("deletefpbatchfromassembly takes 1 argument\n")
	}

	finishedProductBatch := cmd.Args[0]
	exists, err := s.db.CheckFPBatchExistsInQCRelease(context.Background(), finishedProductBatch)
	if err != nil {
		return fmt.Errorf("error checking QC release results for FP batch: %v", err)
	}

	if exists {
		return fmt.Errorf("cannot delete assembly: QC release results exist for finished product batch: %v", finishedProductBatch)
	}

	err = s.db.DeleteFPBatchFromAssembly(context.Background(), finishedProductBatch)
	if err != nil {
		return err
	}

	fmt.Printf("Finished product batch has been deleted from assembly\n")
	return nil
}
