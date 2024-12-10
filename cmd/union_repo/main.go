package main
import (
	"flag"
	"github.com/HUSTSecLab/criticality_score/pkg/union_repo"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.yaml", "Path to the config file")

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)

	union_repo.Run()
}