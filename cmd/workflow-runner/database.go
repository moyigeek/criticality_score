package main

import (
	"math"
	"time"

	"github.com/HUSTSecLab/criticality_score/cmd/workflow-runner/internal/workflow"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
)

func getPreCmdFn(w *workflow.WorkflowNode) func() error {
	return func() error {
		jobid := "job03292390"
		name := w.Name
		event := "start"
		time := time.Now()

		logger.Infof("Write to db not implemented yet: %s, %s, %s, %s", jobid, time, name, event)
		return nil
	}
}

func getPostCmdFn(w *workflow.WorkflowNode) func(e error) error {

	return func(e error) error {

		jobid := "job03292390"
		name := w.Name
		event := "end"
		time := time.Now()

		if e != nil {
			event = "error"
		}

		logger.Infof("Write to db not implemented yet: %s, %s, %s, %s", jobid, time, name, event)
		return nil
	}
}

var updateInterval = map[*workflow.WorkflowNode]time.Duration{
	&srcAllGitMetricsNeedUpdate: 24 * time.Hour,
	&srcDistributionNeedUpdate:  24 * time.Hour,
	&srcDepsDevNeedUpdate:       24 * time.Hour,
	&srcGitPlatformNeedUpdate:   24 * time.Hour,
	&srcGitlinkNeedUpdate:       24 * time.Hour,
}

func getNeedUpdateFn(w *workflow.WorkflowNode) func() bool {
	return func() bool {
		logger.Infof("Read from db not implemeted", w.Name)
		// TODO: read from db
		lastUpdateTime := time.Now().Add(-time.Hour * 12)

		if time.Since(lastUpdateTime) > updateInterval[w] {
			return true
		}
		return true
	}
}

func getNextUpdateDuration() time.Duration {
	duration := time.Duration(math.MaxInt64)

	for _, d := range updateInterval {
		// TODO: read from db
		lastUpdateTime := time.Now().Add(-time.Hour * 12)

		waitTime := lastUpdateTime.Add(d).Sub(time.Now())

		if waitTime < duration {
			duration = waitTime
		}
	}

	if duration < 0 {
		return 0
	}
	return duration
}
