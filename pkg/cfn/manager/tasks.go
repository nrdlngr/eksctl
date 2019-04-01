package manager

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

type TaskSet struct {
	tasks []Task
}

// Task is a common interface for the stack manager tasks
type Task interface {
	Do(chan error) error
	Describe() string
}

func (t *TaskSet) Append(task Task) {
	t.tasks = append(t.tasks, task)
}

func (t *TaskSet) RunAllAsSequence() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	for _, task := range t.tasks {
		if Run(appendErr, task); len(errs) > 0 {
			return errs
		}
	}
	return nil
}

func (t *TaskSet) RunAllInParallel() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	if Run(appendErr, t.tasks...); len(errs) > 0 {
		return errs
	}
	return nil
}

type taskWithoutParams struct {
	info string
	call func(chan error) error
}

func (t *taskWithoutParams) Describe() string { return t.info }
func (t *taskWithoutParams) Do(errs chan error) error {
	return t.call(errs)
}

type taskWithNameParam struct {
	info string
	name string
	call func(chan error, string) error
}

func (t *taskWithNameParam) Describe() string { return t.info }
func (t *taskWithNameParam) Do(errs chan error) error {
	return t.call(errs, t.name)
}

type taskWithNodeGroupSpec struct {
	info      string
	nodeGroup *api.NodeGroup
	call      func(chan error, *api.NodeGroup) error
}

func (t *taskWithNodeGroupSpec) Describe() string { return t.info }
func (t *taskWithNodeGroupSpec) Do(errs chan error) error {
	return t.call(errs, t.nodeGroup)
}

type taskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(*Stack, chan error) error
}

func (t *taskWithStackSpec) Describe() string { return t.info }
func (t *taskWithStackSpec) Do(errs chan error) error {
	return t.call(t.stack, errs)
}

// Run a series of tasks in parallel and wait for each task them to complete;
// passError should take any errors and do what it needs to in
// a given context, e.g. during serial CLI-driven execution one
// can keep errors in a slice, while in a daemon channel maybe
// more suitable
func Run(passError func(error), tasks ...Task) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	for t := range tasks {
		go func(t int) {
			defer wg.Done()
			logger.Debug("task %d started - %s", t, tasks[t].Describe())
			errs := make(chan error)
			if err := tasks[t].Do(errs); err != nil {
				passError(err)
				return
			}
			if err := <-errs; err != nil {
				passError(err)
				return
			}
			logger.Debug("task %d returned without errors", t)
		}(t)
	}
	logger.Debug("waiting for %d tasks to complete", len(tasks))
	wg.Wait()
}

// RunSingleTask runs a task with a proper error handling
func (c *StackCollection) RunSingleTask(t Task) []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}
	if Run(appendErr, t); len(errs) > 0 {
		return errs
	}
	return nil
}

// CreateClusterWithNodeGroups runs all tasks required to create
// the stacks (a cluster and one or more nodegroups); any errors
// will be returned as a slice as soon as one of the tasks or group
// of tasks is completed
func (c *StackCollection) CreateClusterWithNodeGroups(onlySubset sets.String) []error {
	doCreateClusterStack := &taskWithoutParams{
		call: c.createClusterTask,
	}
	if errs := c.RunSingleTask(doCreateClusterStack); len(errs) > 0 {
		return errs
	}

	return c.CreateAllNodeGroups(onlySubset)
}

// CreateAllNodeGroups runs all tasks required to create the node groups;
// any errors will be returned as a slice as soon as one of the tasks
// or group of tasks is completed
func (c *StackCollection) CreateAllNodeGroups(onlySubset sets.String) []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	createAllNodeGroups := []Task{}
	for i := range c.spec.NodeGroups {
		ng := c.spec.NodeGroups[i]
		if onlySubset != nil && !onlySubset.Has(ng.Name) {
			continue
		}
		t := &taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create nodegroup %q", ng.Name),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		}
		createAllNodeGroups = append(createAllNodeGroups, t)
	}
	if Run(appendErr, createAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}

// CreateOneNodeGroup runs a task to create a single node groups;
// any errors will be returned as a slice as soon as the tasks is
// completed
func (c *StackCollection) CreateOneNodeGroup(ng *api.NodeGroup) []error {
	return c.RunSingleTask(&taskWithNodeGroupSpec{
		info:      fmt.Sprintf("create nodegroup %q", ng.Name),
		nodeGroup: ng,
		call:      c.createNodeGroupTask,
	})
}

// // DeleteAllNodeGroups deletes all nodegroups without waiting
// func (c *StackCollection) DeleteAllNodeGroups(call taskFunc) []error {
// 	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
// 	if err != nil {
// 		return []error{err}
// 	}

// 	errs := []error{}
// 	for _, s := range nodeGroupStacks {
// 		if err := c.DeleteNodeGroup(getNodeGroupName(s)); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}

// 	return errs
// }

// WaitDeleteAllNodeGroups runs all tasks required to delete all the nodegroup
// stacks and wait for all nodegroups to be deleted; any errors will be returned
// as a slice as soon as the group of tasks is completed
func (c *StackCollection) WaitDeleteAllNodeGroups(force bool) []error {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return []error{err}
	}

	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	deleteAllNodeGroups := []Task{}
	for i := range nodeGroupStacks {
		t := &taskWithNameParam{
			info: fmt.Sprintf("delete nodegroup stack %q", nodeGroupStacks[i]),
			call: c.waitDeleteNodeGroupTask,
			name: getNodeGroupName(nodeGroupStacks[i]),
		}
		if force {
			t.call = c.waitForceDeleteNodeGroupTask
		}
		deleteAllNodeGroups = append(deleteAllNodeGroups, t)
	}
	if Run(appendErr, deleteAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}

// // SquentialWaitDeleteStacks will delete each of the given stacks one after another
// func (c *StackCollection) SquentialWaitDeleteStacks(stacks []*Stack) []error {
// 	errs := []error{}
// 	appendErr := func(err error) {
// 		errs = append(errs, err)
// 	}

// 	for _, s := range stacks {
// 		t := &taskWithStackSpec{
// 			info:  fmt.Sprintf("delete stack %q", *s.StackName),
// 			stack: s,
// 			call:  c.WaitDeleteStackBySpec,
// 		}
// 		if Run(appendErr, t); len(errs) > 0 {
// 			return errs
// 		}
// 	}

// 	return nil
// }
