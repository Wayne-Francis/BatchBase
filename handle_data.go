package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
	"github.com/google/uuid"
)

func handlerAddMaterial(s *state, cmd command) error {
	if len(cmd.Args) != 4 {
		return fmt.Errorf("addmaterial requires 4 arguments")
	}
	materiallot := cmd.Args[0]
	materiallot_uuid, err := uuid.Parse(materiallot)
	if err != nil {
		return fmt.Errorf("invalid UUID for material lot: %v", err)
	}
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
		MaterialLot:  materiallot_uuid,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		MaterialType: materialtype,
		MfgDate:      mfgdate_time,
		ExpDate:      expdate_time,
	})

	if err != nil {
		return err
	}
	fmt.Printf("new material added:\nlot: %v,\ntype: %v,\nmfg date: %v,\nexp date: %v\n", materiallot_uuid, materialtype, mfgdate_time, expdate_time)
	return nil
}

func handlerMaterials(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("materials takes no arguments\n")
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
