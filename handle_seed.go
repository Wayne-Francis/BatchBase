package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Wayne_Francis/BatchBase/internal/database"
	"github.com/google/uuid"
)

func handlerSeed(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("seed takes no arguments")
	}

	ctx := context.Background()

	// ============================================================
	// WIPE DEVELOPMENT DATA
	// Delete in dependency order.
	// ============================================================

	if err := s.db.DeleteAllQCReleaseResults(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllAssembly(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllFinishedProducts(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllFill(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllIPCQCResults(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllBlend(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllMaterialUsage(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteAllSpecs(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteMaterials(ctx); err != nil {
		return err
	}

	if err := s.db.DeleteUsers(ctx); err != nil {
		return err
	}

	// ============================================================
	// SEED USER
	// ============================================================

	seedUserID := uuid.New()
	now := time.Now()

	_, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        seedUserID,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "seed_user",
	})
	if err != nil {
		return err
	}

	// The config stores the username, while CreatedBy stores the UUID.
	if err := s.cfg.SetUser("seed_user"); err != nil {
		return err
	}

	fmt.Println("Development database cleared.")
	fmt.Println("Creating seed data...")
	fmt.Println()

	// ============================================================
	// SPECS
	// ============================================================

	if err := seedSpec(
		ctx,
		s,
		seedUserID,
		"Blend Uniformity",
		85,
		115,
		90,
		110,
		5,
		now,
	); err != nil {
		return err
	}

	if err := seedSpec(
		ctx,
		s,
		seedUserID,
		"EmittedDose",
		75,
		125,
		90,
		110,
		5,
		now,
	); err != nil {
		return err
	}

	// Assay deliberately has no mean/RSD limits.
	if err := seedSpec(
		ctx,
		s,
		seedUserID,
		"Assay",
		90,
		110,
		0,
		0,
		0,
		now,
	); err != nil {
		return err
	}

	// ============================================================
	// MATERIALS
	// ============================================================

	materials := []struct {
		lot          string
		materialType string
	}{
		{"API-SX-001", "Salmeterol xinafoate"},
		{"API-SX-002", "Salmeterol xinafoate"},
		{"API-SX-003", "Salmeterol xinafoate"},
		{"API-SX-004", "Salmeterol xinafoate"},
		{"API-SX-005", "Salmeterol xinafoate"},
		{"API-SX-006", "Salmeterol xinafoate"},
		{"API-SX-007", "Salmeterol xinafoate"},

		{"LAC-MONO-001", "Lactose monohydrate"},
		{"LAC-MONO-002", "Lactose monohydrate"},
		{"LAC-MONO-003", "Lactose monohydrate"},
		{"LAC-MONO-004", "Lactose monohydrate"},
		{"LAC-MONO-005", "Lactose monohydrate"},
		{"LAC-MONO-006", "Lactose monohydrate"},
		{"LAC-MONO-007", "Lactose monohydrate"},

		{"LAC-FINE-001", "Fine lactose"},
		{"LAC-FINE-002", "Fine lactose"},
		{"LAC-FINE-003", "Fine lactose"},
		{"LAC-FINE-004", "Fine lactose"},
		{"LAC-FINE-005", "Fine lactose"},
		{"LAC-FINE-006", "Fine lactose"},
		{"LAC-FINE-007", "Fine lactose"},
	}

	for i, material := range materials {
		mfgDate := date(2026, time.August, 1+i)
		expDate := date(2028, time.August, 1+i)

		if err := seedMaterial(
			ctx,
			s,
			seedUserID,
			material.lot,
			material.materialType,
			mfgDate,
			expDate,
			now,
		); err != nil {
			return err
		}
	}

	// ============================================================
	// IP-260901
	// Full manufacturing chain + full QC
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260901", 1, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260901", date(2026, time.September, 1), date(2026, time.September, 2), now); err != nil {
		return err
	}

	if err := seedFill(ctx, s, seedUserID, "IP-260901", date(2026, time.September, 3), date(2026, time.September, 4), now); err != nil {
		return err
	}

	if err := seedFinishedProduct(ctx, s, seedUserID, "IP-260901", "FP-260901", "INHALER-001", "MOUTHPIECE-001", now); err != nil {
		return err
	}

	if err := seedAssembly(ctx, s, seedUserID, "FP-260901", date(2026, time.September, 5), date(2026, time.September, 6), now); err != nil {
		return err
	}

	if err := seedBlendQC(ctx, s, seedUserID, "IP-260901", []float64{
		101.2, 99.8, 100.5, 101.0, 98.9,
		100.3, 101.5, 99.6, 100.8, 99.7,
	}, now); err != nil {
		return err
	}

	if err := seedFinishedProductQC(
		ctx,
		s,
		seedUserID,
		"FP-260901",
		[]float64{
			100.8,
		},
		[]float64{
			101.2, 99.8, 100.5, 102.0, 98.9,
			100.7, 101.5, 99.4, 100.9, 101.8,
			99.7, 100.3, 101.1, 99.9, 102.2,
			100.6, 101.4, 98.8, 100.1, 101.0,
		},
		now,
	); err != nil {
		return err
	}

	// ============================================================
	// IP-260902
	// Full manufacturing chain + full QC
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260902", 2, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260902", date(2026, time.September, 3), date(2026, time.September, 4), now); err != nil {
		return err
	}

	if err := seedFill(ctx, s, seedUserID, "IP-260902", date(2026, time.September, 5), date(2026, time.September, 6), now); err != nil {
		return err
	}

	if err := seedFinishedProduct(ctx, s, seedUserID, "IP-260902", "FP-260902", "INHALER-002", "MOUTHPIECE-002", now); err != nil {
		return err
	}

	if err := seedAssembly(ctx, s, seedUserID, "FP-260902", date(2026, time.September, 7), date(2026, time.September, 8), now); err != nil {
		return err
	}

	if err := seedBlendQC(ctx, s, seedUserID, "IP-260902", []float64{
		97.4, 98.6, 99.2, 96.9, 100.1,
		98.0, 99.5, 97.8, 98.9, 99.7,
	}, now); err != nil {
		return err
	}

	if err := seedFinishedProductQC(
		ctx,
		s,
		seedUserID,
		"FP-260902",
		[]float64{
			97.8,
		},
		[]float64{
			96.2, 98.5, 97.1, 99.4, 95.8,
			98.0, 97.5, 99.1, 96.8, 98.7,
			97.9, 95.9, 98.3, 99.0, 96.6,
			97.4, 98.8, 96.1, 97.7, 98.2,
		},
		now,
	); err != nil {
		return err
	}

	// ============================================================
	// IP-260903
	// Full manufacturing chain + full QC
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260903", 3, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260903", date(2026, time.September, 5), date(2026, time.September, 6), now); err != nil {
		return err
	}

	if err := seedFill(ctx, s, seedUserID, "IP-260903", date(2026, time.September, 7), date(2026, time.September, 8), now); err != nil {
		return err
	}

	if err := seedFinishedProduct(ctx, s, seedUserID, "IP-260903", "FP-260903", "INHALER-003", "MOUTHPIECE-003", now); err != nil {
		return err
	}

	if err := seedAssembly(ctx, s, seedUserID, "FP-260903", date(2026, time.September, 9), date(2026, time.September, 10), now); err != nil {
		return err
	}

	if err := seedBlendQC(ctx, s, seedUserID, "IP-260903", []float64{
		100.6, 101.4, 99.2, 98.8, 100.1,
		99.5, 101.0, 99.7, 100.3, 98.9,
	}, now); err != nil {
		return err
	}

	if err := seedFinishedProductQC(
		ctx,
		s,
		seedUserID,
		"FP-260903",
		[]float64{
			99.4,
		},
		[]float64{
			96.8, 97.5, 98.1, 99.0, 97.2,
			98.7, 97.9, 96.9, 98.4, 97.6,
			99.1, 97.3, 98.0, 96.7, 98.8,
			97.7, 98.3, 97.1, 99.2, 98.5,
		},
		now,
	); err != nil {
		return err
	}

	// ============================================================
	// IP-260904
	// Manufacturing complete + Blend QC
	// NO QC RELEASE DATA
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260904", 4, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260904", date(2026, time.September, 7), date(2026, time.September, 8), now); err != nil {
		return err
	}

	if err := seedFill(ctx, s, seedUserID, "IP-260904", date(2026, time.September, 9), date(2026, time.September, 10), now); err != nil {
		return err
	}

	if err := seedFinishedProduct(ctx, s, seedUserID, "IP-260904", "FP-260904", "INHALER-004", "MOUTHPIECE-004", now); err != nil {
		return err
	}

	if err := seedAssembly(ctx, s, seedUserID, "FP-260904", date(2026, time.September, 11), date(2026, time.September, 12), now); err != nil {
		return err
	}

	if err := seedBlendQC(ctx, s, seedUserID, "IP-260904", []float64{
		99.4, 100.1, 98.7, 100.6, 99.8,
		100.3, 98.9, 99.7, 100.5, 99.2,
	}, now); err != nil {
		return err
	}

	// ============================================================
	// IP-260905
	// Blend + Fill + Blend QC
	// NO FINISHED PRODUCT
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260905", 5, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260905", date(2026, time.September, 9), date(2026, time.September, 10), now); err != nil {
		return err
	}

	if err := seedFill(ctx, s, seedUserID, "IP-260905", date(2026, time.September, 11), date(2026, time.September, 12), now); err != nil {
		return err
	}

	if err := seedBlendQC(ctx, s, seedUserID, "IP-260905", []float64{
		100.4, 99.6, 101.1, 100.7, 99.3,
		100.0, 101.4, 98.9, 100.8, 99.9,
	}, now); err != nil {
		return err
	}

	// ============================================================
	// IP-260906
	// Material usage + Blend
	// NO QC
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260906", 6, now); err != nil {
		return err
	}

	if err := seedBlend(ctx, s, seedUserID, "IP-260906", date(2026, time.September, 13), date(2026, time.September, 14), now); err != nil {
		return err
	}

	// ============================================================
	// IP-260907
	// Material usage only
	// ============================================================

	if err := seedMaterialUsageSet(ctx, s, seedUserID, "IP-260907", 7, now); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("           SEED COMPLETE")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("User: seed_user")
	fmt.Println()
	fmt.Println("IP batches:")
	fmt.Println("  IP-260901  -> Full chain + QC")
	fmt.Println("  IP-260902  -> Full chain + QC")
	fmt.Println("  IP-260903  -> Full chain + QC")
	fmt.Println("  IP-260904  -> Full chain, QC release pending")
	fmt.Println("  IP-260905  -> Fill complete, no finished product")
	fmt.Println("  IP-260906  -> Blend complete, no QC")
	fmt.Println("  IP-260907  -> Material usage only")
	fmt.Println()

	return nil
}

// ============================================================
// HELPERS
// ============================================================

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func seedMaterial(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	lot string,
	materialType string,
	mfgDate time.Time,
	expDate time.Time,
	now time.Time,
) error {
	_, err := s.db.AddRawMaterial(ctx, database.AddRawMaterialParams{
		MaterialLot:  lot,
		CreatedAt:    now,
		UpdatedAt:    now,
		MaterialType: materialType,
		MfgDate:      mfgDate,
		ExpDate:      expDate,
		CreatedBy:    userID,
	})

	return err
}

func seedMaterialUsageSet(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	ipBatch string,
	index int,
	now time.Time,
) error {
	materials := []string{
		fmt.Sprintf("API-SX-%03d", index),
		fmt.Sprintf("LAC-MONO-%03d", index),
		fmt.Sprintf("LAC-FINE-%03d", index),
	}

	for _, materialLot := range materials {
		_, err := s.db.AddMaterialUsage(ctx, database.AddMaterialUsageParams{
			InProcessBatchLot: ipBatch,
			MaterialLot:       materialLot,
			CreatedAt:         now,
			UpdatedAt:         now,
			CreatedBy:         userID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func seedBlend(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	ipBatch string,
	startDate time.Time,
	endDate time.Time,
	now time.Time,
) error {
	_, err := s.db.AddBlend(ctx, database.AddBlendParams{
		InProcessBatchLot: ipBatch,
		BlendStartDate:    startDate,
		BlendEndDate:      endDate,
		CreatedAt:         now,
		UpdatedAt:         now,
		CreatedBy:         userID,
	})

	return err
}

func seedFill(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	ipBatch string,
	startDate time.Time,
	endDate time.Time,
	now time.Time,
) error {
	_, err := s.db.AddFill(ctx, database.AddFillParams{
		InProcessBatchLot: ipBatch,
		FillStartDate:     startDate,
		FillEndDate:       endDate,
		CreatedAt:         now,
		UpdatedAt:         now,
		CreatedBy:         userID,
	})

	return err
}

func seedFinishedProduct(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	ipBatch string,
	fpBatch string,
	component1 string,
	component2 string,
	now time.Time,
) error {
	_, err := s.db.AddFinishedProduct(ctx, database.AddFinishedProductParams{
		InProcessBatchLot:    ipBatch,
		FinishedProductBatch: fpBatch,
		Component1Batch:      component1,
		Component2Batch:      component2,
		CreatedAt:            now,
		UpdatedAt:            now,
		CreatedBy:            userID,
	})

	return err
}

func seedAssembly(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	fpBatch string,
	startDate time.Time,
	endDate time.Time,
	now time.Time,
) error {
	_, err := s.db.AddAssembly(ctx, database.AddAssemblyParams{
		FinishedProductBatch: fpBatch,
		CreatedAt:            now,
		UpdatedAt:            now,
		AssemblyStartDate:    startDate,
		AssemblyEndDate:      endDate,
		CreatedBy:            userID,
	})

	return err
}

func seedSpec(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	testName string,
	minResult float64,
	maxResult float64,
	meanMin float64,
	meanMax float64,
	rsdLimit float64,
	now time.Time,
) error {
	var meanMinValue sql.NullString
	var meanMaxValue sql.NullString
	var rsdLimitValue sql.NullString

	// Assay deliberately has no mean/RSD specifications.
	if meanMin != 0 || meanMax != 0 || rsdLimit != 0 {
		meanMinValue = toNullString(meanMin)
		meanMaxValue = toNullString(meanMax)
		rsdLimitValue = toNullString(rsdLimit)
	}

	_, err := s.db.AddSpecs(ctx, database.AddSpecsParams{
		TestName:  testName,
		MinResult: toNullString(minResult),
		MaxResult: toNullString(maxResult),
		MeanMin:   meanMinValue,
		MeanMax:   meanMaxValue,
		RsdLimit:  rsdLimitValue,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: userID,
	})

	return err
}

func seedBlendQC(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	ipBatch string,
	results []float64,
	now time.Time,
) error {
	for i, result := range results {
		_, err := s.db.AddIPCQCResults(ctx, database.AddIPCQCResultsParams{
			InProcessBatchLot: ipBatch,
			CreatedAt:         now,
			UpdatedAt:         now,
			TestName:          "Blend Uniformity",
			Replicate:         int32(i + 1),
			TestDate:          now,
			Result:            strconv.FormatFloat(result, 'f', 2, 64),
			CreatedBy:         userID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func seedFinishedProductQC(
	ctx context.Context,
	s *state,
	userID uuid.UUID,
	fpBatch string,
	assayResults []float64,
	emittedDoseResults []float64,
	now time.Time,
) error {
	for i, result := range assayResults {
		_, err := s.db.AddQCReleaseResults(ctx, database.AddQCReleaseResultsParams{
			FinishedProductBatch: fpBatch,
			CreatedAt:            now,
			UpdatedAt:            now,
			TestName:             "Assay",
			Replicate:            int32(i + 1),
			TestDate:             now,
			Result:               strconv.FormatFloat(result, 'f', 2, 64),
			CreatedBy:            userID,
		})

		if err != nil {
			return err
		}
	}

	for i, result := range emittedDoseResults {
		_, err := s.db.AddQCReleaseResults(ctx, database.AddQCReleaseResultsParams{
			FinishedProductBatch: fpBatch,
			CreatedAt:            now,
			UpdatedAt:            now,
			TestName:             "EmittedDose",
			Replicate:            int32(i + 1),
			TestDate:             now,
			Result:               strconv.FormatFloat(result, 'f', 2, 64),
			CreatedBy:            userID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
