package main

import (
	"log"
	"os"

	"database/sql"

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
	s := &state{cfg: &cfg, db: dbQueries}
	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("resetusers", handlerResetUsers)
	cmds.register("users", handlerUsers)
	cmds.register("addmaterial", middlewareLoggedIn(handlerAddMaterial))
	cmds.register("listmaterials", handlerListMaterials)
	cmds.register("searchrawmaterialbylot", handlerSearchRawMaterialByLot)
	cmds.register("resetmaterials", handlerResetMaterials)
	cmds.register("addmaterialusage", middlewareLoggedIn(handlerAddMaterialUsage))
	cmds.register("listmaterialusage", handlerListMaterialUsage)
	cmds.register("searchmaterialusagebyipbatch", handlerSearchMaterialUsageByIPBatch)
	cmds.register("searchmaterialusagebymateriallot", handlerSearchMaterialUsageByMaterialLot)
	cmds.register("resetmaterialusage", handlerResetMaterialUsage)
	cmds.register("addblend", middlewareLoggedIn(handlerAddBlend))
	cmds.register("listblends", handlerListBlends)
	cmds.register("searchblendbyipbatch", handlerSearchBlendbyIPBatch)
	cmds.register("resetblends", handlerResetBlend)
	cmds.register("addfill", middlewareLoggedIn(handlerAddFill))
	cmds.register("listfills", handlerListFills)
	cmds.register("searchfillbyipbatch", handlerSearchFillbyIPBatch)
	cmds.register("resetfills", handlerResetFill)
	cmds.register("addfinishedproduct", middlewareLoggedIn(handlerAddfinishedproduct))
	cmds.register("listfinishedproducts", handlerListFinishedProducts)
	cmds.register("searchfinishedproductbyipbatch", handlerSearchFinishedProductByIPBatch)
	cmds.register("resetfinishedproducts", handlerResetFinishedProducts)
	cmds.register("addassembly", middlewareLoggedIn(handlerAddAssembly))
	cmds.register("listassemblies", handlerListAssembly)
	cmds.register("searchassemblybyfpbatch", handlerSearchAssemblyByFPBatch)
	cmds.register("searchassemblybyipbatch", handlerSearchAssemblyByIPBatch)
	cmds.register("resetassemblies", handlerResetAssembly)
	cmds.register("addipcqcresult", middlewareLoggedIn(handlerAddIPCQCResult))
	cmds.register("listipcqcresults", handlerListIPCQCResults)
	cmds.register("searchipcqcresultsbyipcbatch", GetIPCQCResultsByIPCBatch)
	cmds.register("searchipcqcresultsbyfpbatch", GetIPCQCResultsByFPBatch)
	cmds.register("resetipcqcresults", handlerResetIPCQCResult)
	cmds.register("addqcreleaseresults", middlewareLoggedIn(AddQCReleaseResults))
	cmds.register("listqcreleaseresults", handlerListQCReleaseResults)
	cmds.register("searchqcreleaseresultsbyipbatch", GetQCReleaseResultsByIPBatch)
	cmds.register("searchqcreleaseresultsbyfpbatch", GetQCReleaseResultsByFPBatch)
	cmds.register("resetqcreleaseresults", handlerResetQCReleaseResult)
	cmds.register("addspecs", middlewareLoggedIn(handlerAddSpecs))
	cmds.register("listspecs", handlerListSpecs)
	cmds.register("searchspecsbytestname", handlerSearchSpecsByTestName)
	cmds.register("resetspecs", handlerResetSpecs)
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
