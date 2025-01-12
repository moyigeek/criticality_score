package main

import (
	"github.com/HUSTSecLab/criticality_score/cmd/workflow-runner/internal/workflow"
)

var (
	taskCalcScore               workflow.WorkflowNode
	taskUpdateDistruibution     workflow.WorkflowNode
	taskUpdateGitMetrics        workflow.WorkflowNode
	taskUpdateGitMetricsPartial workflow.WorkflowNode
	taskUpdateDepsDev           workflow.WorkflowNode
	taskSyncGitMetrics          workflow.WorkflowNode
	taskEnumeratePlatforms      workflow.WorkflowNode

	srcDistributionNeedUpdate  workflow.WorkflowNode
	srcGitlinkNeedUpdate       workflow.WorkflowNode // triggered manually
	srcGitPlatformNeedUpdate   workflow.WorkflowNode
	srcAllGitMetricsNeedUpdate workflow.WorkflowNode
	srcDepsDevNeedUpdate       workflow.WorkflowNode
)

func initTasks() {
	/** calculate score **/
	taskCalcScore.Name = "calc-score"
	taskCalcScore.Description = "Calculate the total score"
	taskCalcScore.Cmd = []string{"bash", "-c", "sleep 10; echo 'taskCalcScore'"}
	taskCalcScore.Dependencies = []*workflow.WorkflowNode{
		&taskUpdateDistruibution,
		&taskUpdateGitMetrics,
		&taskUpdateGitMetricsPartial,
		&taskUpdateDepsDev,
	}

	/** update distribution **/
	taskUpdateDistruibution.Name = "update-distribution"
	taskUpdateDistruibution.Description = "Update the distribution"
	taskUpdateDistruibution.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskUpdateDistruibution'"}
	taskUpdateDistruibution.Dependencies = []*workflow.WorkflowNode{
		&srcDistributionNeedUpdate,
		&taskSyncGitMetrics,
	}

	/** update git metrics **/
	taskUpdateGitMetrics.Name = "update-git-metrics"
	taskUpdateGitMetrics.Description = "Update the git metrics all"
	taskUpdateGitMetrics.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskUpdateGitMetrics'"}
	taskUpdateGitMetrics.Dependencies = []*workflow.WorkflowNode{
		&srcAllGitMetricsNeedUpdate,
	}

	/** update git metrics partial **/
	taskUpdateGitMetricsPartial.Name = "update-git-metrics-partial"
	taskUpdateGitMetricsPartial.Description = "Update the git metrics partial"
	taskUpdateGitMetricsPartial.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskUpdateGitMetricsPartial'"}
	taskUpdateGitMetricsPartial.Dependencies = []*workflow.WorkflowNode{
		&taskSyncGitMetrics,
	}

	/** update deps dev **/
	taskUpdateDepsDev.Name = "update-deps-dev"
	taskUpdateDepsDev.Description = "Update the dev dependencies"
	taskUpdateDepsDev.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskUpdateDepsDev'"}
	taskUpdateDepsDev.Dependencies = []*workflow.WorkflowNode{
		&srcDepsDevNeedUpdate,
		&taskUpdateGitMetrics,
		&taskUpdateGitMetricsPartial,
	}

	/** sync git metrics **/
	taskSyncGitMetrics.Name = "sync-git-metrics"
	taskSyncGitMetrics.Description = "Sync the git metrics"
	taskSyncGitMetrics.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskSyncGitMetrics'"}
	taskSyncGitMetrics.Dependencies = []*workflow.WorkflowNode{
		&srcGitlinkNeedUpdate,
		&taskEnumeratePlatforms,
	}

	/** enumerate platforms **/
	taskEnumeratePlatforms.Name = "enumerate-platforms"
	taskEnumeratePlatforms.Description = "Enumerate the platforms"
	taskEnumeratePlatforms.Cmd = []string{"bash", "-c", "sleep 1; echo 'taskEnumeratePlatforms'"}
	taskEnumeratePlatforms.Dependencies = []*workflow.WorkflowNode{
		&srcGitPlatformNeedUpdate,
	}
}

func initSources() {
	srcDistributionNeedUpdate.Name = "src-distribution-need-update"
	srcDistributionNeedUpdate.Description = "Check if the distribution needs update"
	srcDistributionNeedUpdate.NeedUpdate = func() bool { return false }

	srcGitlinkNeedUpdate.Name = "src-gitlink-need-update"
	srcGitlinkNeedUpdate.Description = "Check if the gitlink needs update"
	srcGitlinkNeedUpdate.NeedUpdate = func() bool { return true }

	srcGitPlatformNeedUpdate.Name = "src-git-platform-need-update"
	srcGitPlatformNeedUpdate.Description = "Check if the git platform needs update"
	srcGitPlatformNeedUpdate.NeedUpdate = func() bool { return false }

	srcAllGitMetricsNeedUpdate.Name = "src-all-git-metrics-need-update"
	srcAllGitMetricsNeedUpdate.Description = "Check if all git metrics need update"
	srcAllGitMetricsNeedUpdate.NeedUpdate = func() bool { return false }

	srcDepsDevNeedUpdate.Name = "src-deps-dev-need-update"
	srcDepsDevNeedUpdate.Description = "Check if the dev dependencies need update"
	srcDepsDevNeedUpdate.NeedUpdate = func() bool { return false }
}
