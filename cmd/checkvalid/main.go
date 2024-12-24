package main
import (
	"flag"
	"log"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/checkvalid"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var outputFile = flag.String("output", "output.csv", "path to the output file")
func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	invalidLinks := checkvalid.CheckVaild(db)
	checkvalid.WriteCsv(invalidLinks, *outputFile)
	log.Println("checkvalid finished")
}