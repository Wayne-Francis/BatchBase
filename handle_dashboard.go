package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func startDashboard(s *state, cmds *commands) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("========================================")
		fmt.Println("              BATCHBASE")
		fmt.Println("   Pharmaceutical Batch Management")
		fmt.Println("========================================")
		fmt.Println()

		// Live system overview
		userCount, err := s.db.CountUsers(context.Background())
		if err != nil {
			fmt.Println("Error loading user count:", err)
			return
		}

		materialCount, err := s.db.CountMaterials(context.Background())
		if err != nil {
			fmt.Println("Error loading material count:", err)
			return
		}

		ipCount, err := s.db.CountInProcessBatches(context.Background())
		if err != nil {
			fmt.Println("Error loading in-process batch count:", err)
			return
		}

		fpCount, err := s.db.CountFinishedProducts(context.Background())
		if err != nil {
			fmt.Println("Error loading finished product count:", err)
			return
		}

		fmt.Println("SYSTEM OVERVIEW")
		fmt.Println("----------------------------------------")
		fmt.Printf("Users:                  %v\n", userCount)
		fmt.Printf("Materials:              %v\n", materialCount)
		fmt.Printf("In-Process Batches:     %v\n", ipCount)
		fmt.Printf("Finished Products:      %v\n", fpCount)
		fmt.Println()

		fmt.Println("MAIN MENU")
		fmt.Println("----------------------------------------")
		fmt.Println("1. Users")
		fmt.Println("2. Materials")
		fmt.Println("3. Manufacturing")
		fmt.Println("4. Quality Control")
		fmt.Println("5. Batch History")
		fmt.Println("6. Plots")
		fmt.Println("0. Exit")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			usersMenu(s, cmds, reader)
		case "2":
			materialsMenu(s, cmds, reader)
		case "3":
			manufacturingMenu(s, cmds, reader)
		case "4":
			qualityControlMenu(s, cmds, reader)
		case "5":
			batchHistoryMenu(s, cmds, reader)
		case "6":
			plotsMenu(s, cmds, reader)
		case "0":
			fmt.Println("Leaving BatchBase. Goodbye!")
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func usersMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("USERS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. List Users")
		fmt.Println("2. Register User")
		fmt.Println("3. Switch Current User")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			err := cmds.run(s, command{Name: "users"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			name, err := readDashboardInput(reader, "User name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "register",
				Args: []string{name},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			name, err := readDashboardInput(reader, "User name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "login",
				Args: []string{name},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func materialsMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("MATERIALS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Material")
		fmt.Println("2. List Materials")
		fmt.Println("3. Search Material by Lot")
		fmt.Println("4. Material History")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			materialLot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			materialType, err := readDashboardInput(reader, "Material type: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			mfgDate, err := readDashboardInput(reader, "Manufacturing date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			expDate, err := readDashboardInput(reader, "Expiry date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addmaterial",
				Args: []string{
					materialLot,
					materialType,
					mfgDate,
					expDate,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listmaterials"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			lot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchrawmaterialbylot",
				Args: []string{lot},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			lot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "materialhistory",
				Args: []string{lot},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func manufacturingMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("MANUFACTURING")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Material Usage")
		fmt.Println("2. Blends")
		fmt.Println("3. Fills")
		fmt.Println("4. Finished Products")
		fmt.Println("5. Assembly")
		fmt.Println("6. Full Batch History")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			materialUsageMenu(s, cmds, reader)
		case "2":
			blendMenu(s, cmds, reader)
		case "3":
			fillMenu(s, cmds, reader)
		case "4":
			finishedProductMenu(s, cmds, reader)
		case "5":
			assemblyMenu(s, cmds, reader)
		case "6":
			batch, err := readDashboardInput(reader, "IP or FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "batchhistoryfull",
				Args: []string{batch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func materialUsageMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("MATERIAL USAGE")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Material Usage")
		fmt.Println("2. List Material Usage")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("4. Search by Material Lot")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			materialLot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addmaterialusage",
				Args: []string{
					ipBatch,
					materialLot,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listmaterialusage"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchmaterialusagebyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			materialLot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchmaterialusagebymateriallot",
				Args: []string{materialLot},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func blendMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("BLENDS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Blend")
		fmt.Println("2. List Blends")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			startDate, err := readDashboardInput(reader, "Blend start date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			endDate, err := readDashboardInput(reader, "Blend end date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addblend",
				Args: []string{
					ipBatch,
					startDate,
					endDate,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listblends"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchblendbyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func fillMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("FILLS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Fill")
		fmt.Println("2. List Fills")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			startDate, err := readDashboardInput(reader, "Fill start date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			endDate, err := readDashboardInput(reader, "Fill end date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addfill",
				Args: []string{
					ipBatch,
					startDate,
					endDate,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listfills"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchfillbyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func finishedProductMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("FINISHED PRODUCTS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Finished Product")
		fmt.Println("2. List Finished Products")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			component1, err := readDashboardInput(reader, "Component 1 batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			component2, err := readDashboardInput(reader, "Component 2 batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addfinishedproduct",
				Args: []string{
					ipBatch,
					fpBatch,
					component1,
					component2,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listfinishedproducts"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchfinishedproductbyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func assemblyMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("ASSEMBLY")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Assembly")
		fmt.Println("2. List Assemblies")
		fmt.Println("3. Search by FP Batch")
		fmt.Println("4. Search by IP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			startDate, err := readDashboardInput(reader, "Assembly start date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			endDate, err := readDashboardInput(reader, "Assembly end date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "addassembly",
				Args: []string{
					fpBatch,
					startDate,
					endDate,
				},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listassemblies"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchassemblybyfpbatch",
				Args: []string{fpBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchassemblybyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func qualityControlMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("QUALITY CONTROL")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. IPC QC")
		fmt.Println("2. QC Release")
		fmt.Println("3. Specifications")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipcQCMenu(s, cmds, reader)
		case "2":
			qcReleaseMenu(s, cmds, reader)
		case "3":
			specsMenu(s, cmds, reader)
		case "0":
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func ipcQCMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("IPC QC")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add IPC QC Results")
		fmt.Println("2. List IPC QC Results")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("4. Search by FP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			testName, err := readDashboardInput(reader, "Test name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			testDate, err := readDashboardInput(reader, "Test date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			results, err := readDashboardResults(reader)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			args := []string{
				ipBatch,
				testName,
				testDate,
			}
			args = append(args, results...)

			err = cmds.run(s, command{
				Name: "addipcqcresult",
				Args: args,
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listipcqcresults"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchipcqcresultsbyipcbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchipcqcresultsbyfpbatch",
				Args: []string{fpBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func qcReleaseMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("QC RELEASE")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add QC Release Results")
		fmt.Println("2. List QC Release Results")
		fmt.Println("3. Search by IP Batch")
		fmt.Println("4. Search by FP Batch")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			testName, err := readDashboardInput(reader, "Test name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			testDate, err := readDashboardInput(reader, "Test date (DD/MM/YYYY): ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			results, err := readDashboardResults(reader)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			args := []string{
				fpBatch,
				testName,
				testDate,
			}
			args = append(args, results...)

			err = cmds.run(s, command{
				Name: "addqcreleaseresults",
				Args: args,
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listqcreleaseresults"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			ipBatch, err := readDashboardInput(reader, "IP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchqcreleaseresultsbyipbatch",
				Args: []string{ipBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			fpBatch, err := readDashboardInput(reader, "FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchqcreleaseresultsbyfpbatch",
				Args: []string{fpBatch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func specsMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("SPECIFICATIONS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Add Specification")
		fmt.Println("2. List Specifications")
		fmt.Println("3. Search by Test Name")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			testName, err := readDashboardInput(reader, "Test name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			minResult, err := readDashboardInput(reader, "Minimum result: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			maxResult, err := readDashboardInput(reader, "Maximum result: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Println()
			fmt.Println("Does this specification include mean and RSD limits?")
			fmt.Println("1. Yes")
			fmt.Println("2. No")

			includeMeanRSD, err := readDashboardInput(reader, "Select an option: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			var args []string

			switch includeMeanRSD {
			case "1":
				meanMin, err := readDashboardInput(reader, "Mean minimum: ")
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				meanMax, err := readDashboardInput(reader, "Mean maximum: ")
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				rsdLimit, err := readDashboardInput(reader, "RSD limit: ")
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				args = []string{
					testName,
					minResult,
					maxResult,
					meanMin,
					meanMax,
					rsdLimit,
				}

			case "2":
				args = []string{
					testName,
					minResult,
					maxResult,
				}

			default:
				fmt.Println("Invalid option.")
				continue
			}

			err = cmds.run(s, command{
				Name: "addspecs",
				Args: args,
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "listspecs"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			testName, err := readDashboardInput(reader, "Test name: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "searchspecsbytestname",
				Args: []string{testName},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func batchHistoryMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("BATCH HISTORY")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Batch Status")
		fmt.Println("2. Batch History Summary")
		fmt.Println("3. Full Batch History")
		fmt.Println("4. Material History")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			batch, err := readDashboardInput(reader, "IP or FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "batchstatus",
				Args: []string{batch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			batch, err := readDashboardInput(reader, "IP or FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "batchhistorysummary",
				Args: []string{batch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			batch, err := readDashboardInput(reader, "IP or FP batch: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "batchhistoryfull",
				Args: []string{batch},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "4":
			lot, err := readDashboardInput(reader, "Material lot: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = cmds.run(s, command{
				Name: "materialhistory",
				Args: []string{lot},
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func plotsMenu(s *state, cmds *commands, reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("PLOTS")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("1. Assay")
		fmt.Println("2. Blend Uniformity")
		fmt.Println("3. Emitted Dose")
		fmt.Println("0. Back")

		choice, err := readDashboardInput(reader, "Select an option: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch choice {
		case "1":
			err := cmds.run(s, command{Name: "plotassay"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "2":
			err := cmds.run(s, command{Name: "plotblend"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "3":
			err := cmds.run(s, command{Name: "plotemitteddose"})
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "0":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func readDashboardInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func readDashboardResults(reader *bufio.Reader) ([]string, error) {
	fmt.Print("Enter results separated by spaces: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("at least one result is required")
	}

	return strings.Fields(input), nil
}
