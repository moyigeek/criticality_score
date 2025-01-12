package workflow

import (
	"testing"

	"github.com/samber/lo"
)

func TestWorkflow1(t *testing.T) {
	a := WorkflowNode{
		Name: "a",
	}

	b := WorkflowNode{
		Name: "b",
	}

	c := WorkflowNode{
		Name:         "c",
		Dependencies: []*WorkflowNode{&a, &b},
	}

	result, err := caculateBuildSequence(&c, true)

	if err != nil {
		t.Errorf("caculateBuildSequence failed %v", err)
	}

	if !(len(result) == 2 && lo.IndexOf(result[0], &a) != -1 && lo.IndexOf(result[0], &b) != -1 && result[1][0] == &c) {
		t.Errorf("caculateBuildSequence failed")
	}
}

func TestCirculear(t *testing.T) {
	a := WorkflowNode{
		Name: "a",
	}

	b := WorkflowNode{
		Name:         "b",
		Dependencies: []*WorkflowNode{&a},
	}

	c := WorkflowNode{
		Name:         "c",
		Dependencies: []*WorkflowNode{&b},
	}

	a.Dependencies = []*WorkflowNode{&c}

	_, err := caculateBuildSequence(&c, false)

	if err == nil {
		t.Errorf("caculateBuildSequence failed")
	}
}
