package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/HUSTSecLab/criticality_score/cmd/workflow-runner/internal/workflow"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
)

var handler workflow.RunningHandler

func StopCurrentWorkflow() {
	if handler != nil {
		handler.Stop()
	}
}

func main() {
	initCmds()
	initSources()
	initTasks()

	// catch interrupt signal
	c := make(chan os.Signal, 1)
	needStop := make(chan struct{}, 1)
	signal.Notify(c, os.Interrupt)

	var err error

	go func() {
		logger.Info("start rpc server...")
		StartRpcServer()
	}()

	go func() {
		for {
			<-c
			StopCurrentWorkflow()
			needStop <- struct{}{}
		}
	}()

	for {
		handler, err = taskCalcScore.StartWorkflow(nil)
		if err != nil {
			logger.Error("failed to start workflow", err)
		}

		err = handler.Wait()
		if err != nil {
			logger.Error("workflow running failed", err)
		}

		select {
		case <-needStop:
			return
		default:
		}

		waitTime := getNextUpdateDuration()
		logger.Infof("wait for %s to start next workflow", waitTime)

		select {
		case <-needStop:
			return
		case <-time.After(waitTime):
		}
	}
}
