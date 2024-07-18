package main

import (
	"flag"
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/home2git"
)

var (
	flagPackage  = flag.String("package", "", "The package name of the homepage")
	flagHomepage = flag.String("homepage", "", "The homepage of the package")
	flagVerbose  = flag.Int("verbose", 0, "Verbose level") // TODO: log level
)

func main() {
	flag.Parse()
	pkg := *flagPackage
	homepage := *flagHomepage

	if pkg == "" || homepage == "" {
		flag.Usage()
		return
	}

	fmt.Print(home2git.HomepageToGit(homepage, pkg))
}
