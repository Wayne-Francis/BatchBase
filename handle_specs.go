package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
)

func handlerAddSpecs(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 6 {
		return fmt.Errorf("addspecs requires 6 arguments")
	}

	testName := cmd.Args[0]

	minResult, err := strconv.ParseFloat(cmd.Args[1], 64)
	if err != nil {
		return fmt.Errorf("invalid min result: %v", err)
	}

	maxResult, err := strconv.ParseFloat(cmd.Args[2], 64)
	if err != nil {
		return fmt.Errorf("invalid max result: %v", err)
	}

	meanMin, err := strconv.ParseFloat(cmd.Args[3], 64)
	if err != nil {
		return fmt.Errorf("invalid mean min: %v", err)
	}

	meanMax, err := strconv.ParseFloat(cmd.Args[4], 64)
	if err != nil {
		return fmt.Errorf("invalid mean max: %v", err)
	}

	rsdLimit, err := strconv.ParseFloat(cmd.Args[5], 64)
	if err != nil {
		return fmt.Errorf("invalid rsd limit: %v", err)
	}

	now := time.Now()

	_, err = s.db.AddSpecs(context.Background(), database.AddSpecsParams{
		TestName:  testName,
		MinResult: toNullString(minResult),
		MaxResult: toNullString(maxResult),
		MeanMin:   toNullString(meanMin),
		MeanMax:   toNullString(meanMax),
		RsdLimit:  toNullString(rsdLimit),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf(
		"New specs added:\nTest Name: %v\nMin Result: %v\nMax Result: %v\nMean Min: %v\nMean Max: %v\nRSD Limit: %v\nCreated At: %v\nUpdated At: %v\nCreated By: %v\n",
		testName,
		minResult,
		maxResult,
		meanMin,
		meanMax,
		rsdLimit,
		now,
		now,
		user.ID,
	)

	return nil
}

func handlerListSpecs(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("listspecs takes no arguments")
	}

	specs, err := s.db.GetSpecs(context.Background())
	if err != nil {
		return err
	}

	for _, spec := range specs {
		fmt.Printf(
			"Test Name: %v\nMin Result: %v\nMax Result: %v\nMean Min: %v\nMean Max: %v\nRSD Limit: %v\nCreated At: %v\nUpdated At: %v\nCreated By: %v\n\n",
			spec.TestName,
			spec.MinResult.String,
			spec.MaxResult.String,
			spec.MeanMin.String,
			spec.MeanMax.String,
			spec.RsdLimit.String,
			spec.CreatedAt,
			spec.UpdatedAt,
			spec.CreatedBy,
		)
	}

	return nil
}

func handlerSearchSpecsByTestName(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("searchspecs takes 1 argument")
	}

	testName := cmd.Args[0]

	spec, err := s.db.GetSpecsByTestName(context.Background(), testName)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Test Name: %v\nMin Result: %v\nMax Result: %v\nMean Min: %v\nMean Max: %v\nRSD Limit: %v\nCreated At: %v\nUpdated At: %v\nCreated By: %v\n",
		spec.TestName,
		spec.MinResult.String,
		spec.MaxResult.String,
		spec.MeanMin.String,
		spec.MeanMax.String,
		spec.RsdLimit.String,
		spec.CreatedAt,
		spec.UpdatedAt,
		spec.CreatedBy,
	)

	return nil
}

func handlerResetSpecs(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("resetspecs takes no arguments")
	}

	err := s.db.DeleteAllSpecs(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Specs have been deleted")
	return nil
}

func toNullString(value float64) sql.NullString {
	return sql.NullString{
		String: strconv.FormatFloat(value, 'f', 2, 64),
		Valid:  true,
	}
}
