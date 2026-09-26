package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddMaterial(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 4 {
		return fmt.Errorf("addmaterial requires 4 arguments")
	}
	materiallot := cmd.Args[0]
	materialtype := cmd.Args[1]
	mfgdate := cmd.Args[2]
	mfgdate_time, err := time.Parse("02/01/2006", mfgdate)
	if err != nil {
		return fmt.Errorf("invalid date format for mfg date: %v", err)
	}
	expdate := cmd.Args[3]
	expdate_time, err := time.Parse("02/01/2006", expdate)
	if err != nil {
		return fmt.Errorf("invalid date format for exp date: %v", err)
	}
	_, err = s.db.AddRawMaterial(context.Background(), database.AddRawMaterialParams{
		MaterialLot:  materiallot,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		MaterialType: materialtype,
		MfgDate:      mfgdate_time,
		ExpDate:      expdate_time,
		CreatedBy:    user.ID,
	})

	if err != nil {
		return err
	}
	fmt.Printf("new material added:\nlot: %v,\ntype: %v,\nmfg date: %v,\nexp date: %v\ncreated by: %v\n", materiallot, materialtype, mfgdate_time, expdate_time, user.ID)
	return nil
}

func handlerListMaterials(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listmaterials takes no arguments\n")
	}
	materials, err := s.db.GetRawMaterials(context.Background())
	if err != nil {
		return err
	}
	for _, material := range materials {
		fmt.Printf("Lot: %v\nType: %v\nMfg Date: %v\nExp Date: %v\n", material.MaterialLot, material.MaterialType, material.MfgDate, material.ExpDate)
	}
	return nil
}

func handlerSearchRawMaterialByLot(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchrawmaterialbylot takes 1 argument\n")
	}
	materialLot := cmd.Args[0]
	material, err := s.db.GetRawMaterialByLot(context.Background(), materialLot)
	if err != nil {
		return err
	}
	fmt.Printf("Lot: %v\nType: %v\nMfg Date: %v\nExp Date: %v\n", material.MaterialLot, material.MaterialType, material.MfgDate, material.ExpDate)
	return nil
}

func handlerResetMaterials(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetmaterials takes no arguments\n")
	}
	err := s.db.DeleteMaterials(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Materials have been deleted\n")
	return nil
}
