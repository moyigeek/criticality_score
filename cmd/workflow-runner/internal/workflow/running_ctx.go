package workflow

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/HUSTSecLab/criticality_score/pkg/logger"
)

type RunningHandler interface {
	Wait() error
	Stop() error
	Kill() error
}

type runningHandler struct {
	finish chan error
	isBusy bool
	stop   chan struct{}
	kill   chan struct{}
}

func newRunningHandler() *runningHandler {
	return &runningHandler{
		finish: make(chan error),
		stop:   make(chan struct{}),
		kill:   make(chan struct{}),
	}
}

func (h *runningHandler) Stop() error {
	if h.isBusy {
		h.stop <- struct{}{}
		return nil
	}
	return fmt.Errorf("no running process")
}

func (h *runningHandler) Kill() error {
	if h.isBusy {
		h.kill <- struct{}{}
		return nil
	}
	return fmt.Errorf("no running process")
}

func (h *runningHandler) Wait() error {
	return <-h.finish
}

type runningCtx struct {
	loggerFile     *os.File
	runningHandler *runningHandler
	flow           *WorkflowNode
}

func (ctx *runningCtx) runCmd() error {
	if ctx.flow.Cmd == nil {
		return nil
	}

	logger.Infof("Running command %v", ctx.flow.Cmd)

	n := ctx.flow
	logPrefix := n.LogPrefix
	if logPrefix == "" {
		logPrefix = n.Name
	}

	cmd := exec.Command(n.Cmd[0], n.Cmd[1:]...)

	finish := make(chan error, 1)

	// if stop is received, or stop, kill the process
	go func() {
		ctx.runningHandler.isBusy = true
		defer func() { ctx.runningHandler.isBusy = false }()

		if ctx.runningHandler == nil {
			return
		}

		select {
		case <-ctx.runningHandler.stop:
			logger.Infof("stopping process %s", n.Name)
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				logger.Errorf("failed to stop process %s: %v", n.Name, err)
			}
		case <-ctx.runningHandler.kill:
			logger.Infof("killing process %s", n.Name)
			err := cmd.Process.Kill()
			if err != nil {
				logger.Errorf("failed to kill process %s: %v", n.Name, err)
			}
		case <-finish:
		}
	}()

	cmd.Stdout = ctx.loggerFile
	cmd.Stderr = ctx.loggerFile

	cmdRetStatus := cmd.Run()

	finish <- cmdRetStatus
	return cmdRetStatus
}

func (ctx *runningCtx) Run() error {
	n := ctx.flow
	needUpdate := true

	if n.NeedUpdate != nil && !n.NeedUpdate() {
		needUpdate = false
	}

	if needUpdate {
		if n.RunBeforeCmd != nil {
			if err := n.RunBeforeCmd(); err != nil {
				return err
			}
		}

		if err := ctx.runCmd(); err != nil {
			return err
		}

		if n.RunAfterCmd != nil {
			if err := n.RunAfterCmd(); err != nil {
				return err
			}
		}
	}
	return nil
}
