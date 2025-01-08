package main

import (
	"flag"
	"strings"
	"time"

	"github.com/HUSTSecLab/criticality_score/cmd/git-platforms-enumerator/internal/enumerator"
	"github.com/HUSTSecLab/criticality_score/cmd/git-platforms-enumerator/internal/writer"
	log "github.com/sirupsen/logrus"
)

// dateFlag implements the flag.Value interface to simplify the input and validation of
// dates from the command line.
type dateFlag time.Time

const dateFormat = "2006-01-02"

func (d *dateFlag) Set(value string) error {
	t, err := time.Parse(dateFormat, value)
	if err != nil {
		return err
	}
	*d = dateFlag(t)
	return nil
}

func (d *dateFlag) String() string {
	return (*time.Time)(d).Format(dateFormat)
}

func (d *dateFlag) Time() time.Time {
	return time.Time(*d)
}

func main() {
	// flags
	var (
		flagConfig      = flag.String("config", "", "path to the configuration file, when output type is db")
		flagPlatforms   = flag.String("platforms", "", "comma separated list of platforms to enumerate")
		flagOutputType  = flag.String("output", "stdout", "output type: allow stdout, file, db")
		flagOutputFilev = flag.String("output-file", "", "output file")
		flagJobs        = flag.Int("jobs", 10, "number of concurrent jobs")
		flagTake        = flag.Int("take", 1000, "number of repositories to enumerate")
	)

	// github flags
	var (
		flagMinStars        = flag.Int("min-stars", 100, "minimum number of stars")
		flagStarOverlap     = flag.Int("star-overlap", 5, "minimum number of stars overlap")
		flagRequireMinStars = flag.Bool("require-min-stars", false, "require minimum number of stars")
		flagQuery           = flag.String("query", "is:public", "sets the base query")
		flagStartDate       = dateFlag(enumerator.GithubEpochDate)
		flagEndDate         = dateFlag(time.Now().UTC().Truncate(time.Hour * 24))
	)

	flag.Var(&flagStartDate, "start-date", "start date for the search")
	flag.Var(&flagEndDate, "end-date", "end date for the search")
	flag.Parse()

	configPath := *flagConfig
	platforms := strings.Split(*flagPlatforms, ",")

	for _, platform := range platforms {
		var w writer.Writer
		var tablePrefix string
		var en enumerator.Enumerator

		switch platform {
		case "github":
			tablePrefix = "github"
			githubConfig := enumerator.GithubEnumeratorConfig{
				MinStars:        *flagMinStars,
				StarOverlap:     *flagStarOverlap,
				RequireMinStars: *flagRequireMinStars,
				Query:           *flagQuery,
				StartDate:       flagStartDate.Time(),
				EndDate:         flagEndDate.Time(),
				Workers:         *flagJobs,
			}
			en = enumerator.NewGithubEnumerator(&githubConfig)
		case "gitlab":
			tablePrefix = "gitlab_links"
			en = enumerator.NewGitlabEnumerator(*flagTake, *flagJobs)
		case "bitbucket":
			tablePrefix = "bitbucket_links"
			en = enumerator.NewBitBucketEnumerator(*flagTake)
		default:
			panic("unknown platform")
		}

		switch *flagOutputType {
		case "stdout":
			w = writer.NewStdOutWriter()
		case "file":
			w = writer.NewTextFileWriter(*flagOutputFilev)
		case "db":
			w = writer.NewDatabaseWriter(*&configPath, tablePrefix)
		default:
			panic("unknown output type")
		}

		en.SetWriter(w)

		err := en.Enumerate()

		if err != nil {
			log.WithError(err).Errorf("failed to enumerate %s", platform)
		}

	}

}
