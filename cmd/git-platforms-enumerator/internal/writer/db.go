package writer

import (
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repositories"
)

type DatabaseWriter struct {
	dbCtx       storage.AppDatabaseContext
	repo        repositories.PlatformLinkRepository
	tablePrefix string

	buffer     []string
	bufferSize int
}

func NewDatabaseWriter(ctx storage.AppDatabaseContext, tablePrefix string) *DatabaseWriter {
	return &DatabaseWriter{
		dbCtx:       ctx,
		tablePrefix: tablePrefix,
		buffer:      make([]string, 0),
		bufferSize:  1000,
	}
}

func (w *DatabaseWriter) Open() error {
	repo := repositories.NewPlatformLinkRepository(w.dbCtx, repositories.PlatformLinkTablePrefix(w.tablePrefix))
	w.repo = repo
	return repo.ClearLinks()
}

func (w *DatabaseWriter) Close() error {
	w.flush()
	return nil
}

func (w *DatabaseWriter) flush() error {
	err := w.repo.BatchInsertLinks(w.buffer)
	if err != nil {
		logger.Error("Failed to insert links: %v", err)
		return err
	}
	w.buffer = make([]string, 0)
	return nil
}

func (w *DatabaseWriter) Write(url string) error {
	w.buffer = append(w.buffer, url)

	if len(w.buffer) >= w.bufferSize {
		err := w.flush()
		if err != nil {
			logger.Error("Failed to flush buffer: %v", err)
			return err
		}
	}

	return nil
}
