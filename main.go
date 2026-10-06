package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Wayne_Francis/BatchBase/internal/config"
	"github.com/Wayne_Francis/BatchBase/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Dburl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}

	// Users
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("resetusers", handlerResetUsers)
	cmds.register("users", handlerUsers)

	// Materials
	cmds.register("addmaterial", middlewareLoggedIn(handlerAddMaterial))
	cmds.register("listmaterials", handlerListMaterials)
	cmds.register("searchrawmaterialbylot", handlerSearchRawMaterialByLot)
	cmds.register("deletematerialfrommaterials", handlerDeleteMaterialFromMaterials)
	cmds.register("resetmaterials", handlerResetMaterials)

	// Material Usage
	cmds.register("addmaterialusage", middlewareLoggedIn(handlerAddMaterialUsage))
	cmds.register("listmaterialusage", handlerListMaterialUsage)
	cmds.register("searchmaterialusagebyipbatch", handlerSearchMaterialUsageByIPBatch)
	cmds.register("searchmaterialusagebymateriallot", handlerSearchMaterialUsageByMaterialLot)
	cmds.register("deletematerialfromusage", DeleteMaterialFromUsage)
	cmds.register("resetmaterialusage", handlerResetMaterialUsage)

	// Blend
	cmds.register("addblend", middlewareLoggedIn(handlerAddBlend))
	cmds.register("listblends", handlerListBlends)
	cmds.register("searchblendbyipbatch", handlerSearchBlendbyIPBatch)
	cmds.register("deletelotfromblend", DeleteLotFromBlend)
	cmds.register("resetblends", handlerResetBlend)

	// Fill
	cmds.register("addfill", middlewareLoggedIn(handlerAddFill))
	cmds.register("listfills", handlerListFills)
	cmds.register("searchfillbyipbatch", handlerSearchFillbyIPBatch)
	cmds.register("deletelotfromfill", DeleteLotFromFill)
	cmds.register("resetfills", handlerResetFill)

	// Finished Product
	cmds.register("addfinishedproduct", middlewareLoggedIn(handlerAddfinishedproduct))
	cmds.register("listfinishedproducts", handlerListFinishedProducts)
	cmds.register("searchfinishedproductbyipbatch", handlerSearchFinishedProductByIPBatch)
	cmds.register("deletelotfromfinishedproduct", DeleteLotFromFinishedProduct)
	cmds.register("resetfinishedproducts", handlerResetFinishedProducts)

	// Assembly
	cmds.register("addassembly", middlewareLoggedIn(handlerAddAssembly))
	cmds.register("listassemblies", handlerListAssembly)
	cmds.register("searchassemblybyfpbatch", handlerSearchAssemblyByFPBatch)
	cmds.register("searchassemblybyipbatch", handlerSearchAssemblyByIPBatch)
	cmds.register("deletefpbatchfromassembly", handlerDeleteFPBatchFromAssembly)
	cmds.register("resetassemblies", handlerResetAssembly)

	// IPC QC
	cmds.register("addipcqcresult", middlewareLoggedIn(handlerAddIPCQCResult))
	cmds.register("listipcqcresults", handlerListIPCQCResults)
	cmds.register("searchipcqcresultsbyipcbatch", GetIPCQCResultsByIPCBatch)
	cmds.register("searchipcqcresultsbyfpbatch", GetIPCQCResultsByFPBatch)
	cmds.register("deleteipcqcresult", handlerDeleteIPCQCResult)
	cmds.register("resetipcqcresults", handlerResetIPCQCResult)

	// QC Release
	cmds.register("addqcreleaseresults", middlewareLoggedIn(AddQCReleaseResults))
	cmds.register("listqcreleaseresults", handlerListQCReleaseResults)
	cmds.register("searchqcreleaseresultsbyipbatch", GetQCReleaseResultsByIPBatch)
	cmds.register("searchqcreleaseresultsbyfpbatch", GetQCReleaseResultsByFPBatch)
	cmds.register("deleteqcreleaseresult", handlerDeleteQCReleaseResult)
	cmds.register("resetqcreleaseresults", handlerResetQCReleaseResult)

	// Specs
	cmds.register("addspecs", middlewareLoggedIn(handlerAddSpecs))
	cmds.register("listspecs", handlerListSpecs)
	cmds.register("searchspecsbytestname", handlerSearchSpecsByTestName)
	cmds.register("resetspecs", handlerResetSpecs)
	cmds.register("deletespecs", handlerDeleteSpecs)

	// Batch History
	cmds.register("materialhistory", handlerMaterialHistory)
	cmds.register("batchstatus", handlerBatchStatus)
	cmds.register("batchhistorysummary", handlerBatchHistorySummary)
	cmds.register("batchhistoryfull", handlerBatchHistoryFull)

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("please type commands")
	}

	c := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.run(s, c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
