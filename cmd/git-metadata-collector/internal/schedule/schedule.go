package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
)

const FetchSize = 200
const IdleInterval = 30 * time.Second

var tasks []string = make([]string, 0, FetchSize)
var manualTasks []string = make([]string, 0)
var muTasks sync.Mutex

func fetchTasksFromDatabase() error {
	ctx := storage.GetDefaultAppDatabaseContext()
	r := repository.NewRankedGitTaskRepository(ctx)

	for {
		ok := false

		result, err := r.Query(FetchSize)
		if err != nil {
			return err
		}

		for r := range result {
			ok = true
			tasks = append(tasks, *r.GitLink)
		}

		if !ok {
			time.Sleep(IdleInterval)
		} else {
			break
		}
	}
	return nil
}

func AddManualTask(task string) {
	muTasks.Lock()
	defer muTasks.Unlock()

	manualTasks = append(manualTasks, task)
}

func GetTask() (string, error) {
	// TODO: optimize the performance
	muTasks.Lock()
	defer muTasks.Unlock()

	if len(manualTasks) > 0 {
		task := manualTasks[0]
		manualTasks = manualTasks[1:]
		return task, nil
	}

	if len(tasks) == 0 {
		err := fetchTasksFromDatabase()

		if err != nil {
			return "", err
		}
	}

	if len(tasks) > 0 {
		task := tasks[0]
		tasks = tasks[1:]
		return task, nil
	}

	return "", fmt.Errorf("no task available")
}
