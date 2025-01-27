package task

import (
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
	"github.com/lib/pq"
)

func Collect(gitLink string) {
	gmr := repository.NewGitMetricsRepository(storage.GetDefaultAppDatabaseContext())

	recordFail := func(e error) {
		logger.Errorf("Collecting %s Failed: %s", gitLink, e)
		gmr.InsertOrUpdateFailed(&repository.FailedGitMetric{
			GitLink:    sqlutil.ToData(gitLink),
			Message:    sqlutil.ToNullable(e.Error()),
			UpdateTime: sqlutil.ToNullable(time.Now()),
		})
	}

	recordSuccess := func(repo *git.Repo) {
		logger.Infof("%s Collected", gitLink)

		err := gmr.InsertOrUpdate(&repository.GitMetric{
			GitLink:          sqlutil.ToData(gitLink),
			CreatedSince:     sqlutil.ToNullable(repo.CreatedSince),
			UpdatedSince:     sqlutil.ToNullable(repo.UpdatedSince),
			ContributorCount: sqlutil.ToNullable(repo.ContributorCount),
			CommitFrequency:  sqlutil.ToNullable(repo.CommitFrequency),
			OrgCount:         sqlutil.ToNullable(repo.OrgCount),
			License:          sqlutil.ToNullable(pq.StringArray(repo.Licenses)),
			Language:         sqlutil.ToNullable(pq.StringArray(repo.Languages)),
		})

		if err != nil {
			logger.Fatalf("Inserting %s Failed", gitLink)
		}
	}

	u := url.ParseURL(gitLink)
	r, err := collector.Collect(&u, config.GetGitStoragePath())
	if err != nil {
		recordFail(err)
		return
	}
	repo, err := git.ParseRepo(r)
	if err != nil {
		recordFail(err)
		return
	}

	recordSuccess(repo)
}
