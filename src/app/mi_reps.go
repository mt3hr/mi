package mi

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mt3hr/rykv/kyou"
)

type MiReps []MiRep

func (m MiReps) AddStartInfo(startInfo *StartInfo) error {
	return m.AddStartInfo(startInfo)
}

func (m MiReps) AddEndInfo(endInfo *EndInfo) error {
	return m.AddEndInfo(endInfo)
}

func (m MiReps) UpdateCache(ctx context.Context) error {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	errch := make(chan error, len(m))
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		go func(miRep MiRep) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil {
					errch <- err
				}
			default:
				err := miRep.UpdateCache(ctx)
				if err != nil {
					errch <- err
				}
			}
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			existErr = true
			break errloop
		case e := <-errch:
			err = fmt.Errorf("error at getAllCheckStateInfos: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return err
	}
	return nil
}

func (m MiReps) GetAllCheckStateInfos(ctx context.Context) ([]*CheckStateInfo, error) {
	checkStateInfos := map[string]*CheckStateInfo{}
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan []*CheckStateInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			matchTaskscheckStateInfos, err := miRep.GetAllCheckStateInfos(ctx)
			if err != nil {
				// errch <- err
				return
			}
			ch <- matchTaskscheckStateInfos
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at getAllCheckStateInfos: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

loop:
	for {
		select {
		case t := <-ch:
			if t == nil {
				continue loop
			}
			for _, checkStateInfo := range t {
				checkStateInfos[checkStateInfo.CheckStateID] = checkStateInfo
			}
		default:
			break loop
		}
	}

	allCheckStateInfos := []*CheckStateInfo{}
	for _, checkStateInfo := range checkStateInfos {
		if checkStateInfo == nil {
			continue
		}
		allCheckStateInfos = append(allCheckStateInfos, checkStateInfo)
	}
	return allCheckStateInfos, nil
}

func (m MiReps) SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error) {
	taskMap := map[string]*Task{}

	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan []*Task, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			tasks, err := miRep.SearchTasks(ctx, word, query)
			if err != nil {
				// errch <- err
				return
			}
			ch <- tasks
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at SearchTask: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

loop:
	for {
		select {
		case t := <-ch:
			if t == nil {
				continue loop
			}
			for _, task := range t {
				taskMap[task.TaskID] = task
			}
		default:
			break loop
		}
	}

	tasks := []*Task{}
	for _, task := range taskMap {
		if task == nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m MiReps) GetAllTasks(ctx context.Context) ([]*Task, error) {
	taskMap := map[string]*Task{}
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan []*Task, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			tasks, err := miRep.GetAllTasks(ctx)
			if err != nil {
				// errch <- err
				return
			}
			ch <- tasks
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}
loop:
	for {
		select {
		case t := <-ch:
			if t == nil {
				continue loop
			}
			for _, task := range t {
				taskMap[task.TaskID] = task
			}
		default:
			break loop
		}
	}

	tasks := []*Task{}
	for _, task := range taskMap {
		if task == nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m MiReps) GetTask(ctx context.Context, taskID string) (*Task, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *Task, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			tasks, err := miRep.GetTask(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- tasks
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}
	tasks := []*Task{}
loop:
	for {
		select {
		case task := <-ch:
			if task == nil {
				continue loop
			}
			tasks = append(tasks, task)
		default:
			break loop
		}
	}

	sort.Slice(tasks, func(i int, j int) bool {
		return tasks[i].CreatedTime.After(tasks[j].CreatedTime)
	})
	if len(tasks) == 0 {
		return nil, nil
	}
	return tasks[0], nil
}

func (m MiReps) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *CheckStateInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			checkStateInfo, err := miRep.GetLatestCheckStateInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- checkStateInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	checkStateInfos := []*CheckStateInfo{}
loop:
	for {
		select {
		case checkStateInfo := <-ch:
			if checkStateInfo == nil {
				continue loop
			}
			checkStateInfos = append(checkStateInfos, checkStateInfo)
		default:
			break loop
		}
	}

	sort.Slice(checkStateInfos, func(i int, j int) bool {
		return checkStateInfos[i].UpdatedTime.After(checkStateInfos[j].UpdatedTime)
	})
	if len(checkStateInfos) == 0 {
		return nil, nil
	}
	return checkStateInfos[0], nil
}

func (m MiReps) GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error) {
	taskTitleInfoMap := map[string]*TaskTitleInfo{}
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *TaskTitleInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			taskTitleInfo, err := miRep.GetLatestTaskTitleInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- taskTitleInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	taskTitleInfos := []*TaskTitleInfo{}
loop:
	for {
		select {
		case taskTitleInfo := <-ch:
			if taskTitleInfo == nil {
				continue loop
			}
			taskTitleInfoMap[taskTitleInfo.TaskTitleID] = taskTitleInfo
			taskTitleInfos = append(taskTitleInfos, taskTitleInfo)
		default:
			break loop
		}
	}

	sort.Slice(taskTitleInfos, func(i int, j int) bool {
		return taskTitleInfos[i].UpdatedTime.After(taskTitleInfos[j].UpdatedTime)
	})
	if len(taskTitleInfos) == 0 {
		return nil, nil
	}
	return taskTitleInfos[0], nil
}

func (m MiReps) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *LimitInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			limitInfo, err := miRep.GetLatestLimitInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- limitInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	limitInfos := []*LimitInfo{}
loop:
	for {
		select {
		case limitInfo := <-ch:
			if limitInfo == nil {
				continue loop
			}
			limitInfos = append(limitInfos, limitInfo)
		default:
			break loop
		}
	}

	sort.Slice(limitInfos, func(i int, j int) bool {
		return limitInfos[i].UpdatedTime.After(limitInfos[j].UpdatedTime)
	})
	if len(limitInfos) == 0 {
		return nil, nil
	}
	return limitInfos[0], nil
}

func (m MiReps) GetLatestStartInfoFromTaskID(ctx context.Context, taskID string) (*StartInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *StartInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			startInfo, err := miRep.GetLatestStartInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- startInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	startInfos := []*StartInfo{}
loop:
	for {
		select {
		case startInfo := <-ch:
			if startInfo == nil {
				continue loop
			}
			startInfos = append(startInfos, startInfo)
		default:
			break loop
		}
	}

	sort.Slice(startInfos, func(i int, j int) bool {
		return startInfos[i].UpdatedTime.After(startInfos[j].UpdatedTime)
	})
	if len(startInfos) == 0 {
		return nil, nil
	}
	return startInfos[0], nil
}

func (m MiReps) GetLatestEndInfoFromTaskID(ctx context.Context, taskID string) (*EndInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *EndInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			endInfo, err := miRep.GetLatestEndInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- endInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	endInfos := []*EndInfo{}
loop:
	for {
		select {
		case endInfo := <-ch:
			if endInfo == nil {
				continue loop
			}
			endInfos = append(endInfos, endInfo)
		default:
			break loop
		}
	}

	sort.Slice(endInfos, func(i int, j int) bool {
		return endInfos[i].UpdatedTime.After(endInfos[j].UpdatedTime)
	})
	if len(endInfos) == 0 {
		return nil, nil
	}
	return endInfos[0], nil
}

func (m MiReps) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *BoardInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			boardInfo, err := miRep.GetLatestBoardInfoFromTaskID(ctx, taskID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- boardInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	boardInfos := []*BoardInfo{}
loop:
	for {
		select {
		case boardInfo := <-ch:
			if boardInfo == nil {
				continue loop
			}
			boardInfos = append(boardInfos, boardInfo)
		default:
			break loop
		}
	}

	sort.Slice(boardInfos, func(i int, j int) bool {
		return boardInfos[i].UpdatedTime.After(boardInfos[j].UpdatedTime)
	})
	if len(boardInfos) == 0 {
		return nil, nil
	}
	return boardInfos[0], nil
}

func (m MiReps) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *CheckStateInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			checkStateInfo, err := miRep.GetCheckStateInfo(ctx, checkStateID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- checkStateInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	checkStateInfos := []*CheckStateInfo{}
loop:
	for {
		select {
		case checkStateInfo := <-ch:
			if checkStateInfo == nil {
				continue loop
			}
			checkStateInfos = append(checkStateInfos, checkStateInfo)
		default:
			break loop
		}
	}

	sort.Slice(checkStateInfos, func(i int, j int) bool {
		return checkStateInfos[i].UpdatedTime.After(checkStateInfos[j].UpdatedTime)
	})
	for _, checkStateInfo := range checkStateInfos {
		return checkStateInfo, nil
	}

	return nil, fmt.Errorf("check state not found. checkStateID=%s", checkStateID)
}

func (m MiReps) GetTaskTitleInfo(ctx context.Context, taskTitleID string) (*TaskTitleInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *TaskTitleInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			taskTitleInfo, err := miRep.GetTaskTitleInfo(ctx, taskTitleID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- taskTitleInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	taskTitleInfos := []*TaskTitleInfo{}
loop:
	for {
		select {
		case taskTitleInfo := <-ch:
			if taskTitleInfo == nil {
				continue loop
			}
			taskTitleInfos = append(taskTitleInfos, taskTitleInfo)
		default:
			break loop
		}
	}

	sort.Slice(taskTitleInfos, func(i int, j int) bool {
		return taskTitleInfos[i].UpdatedTime.After(taskTitleInfos[j].UpdatedTime)
	})
	for _, taskTitleInfo := range taskTitleInfos {
		return taskTitleInfo, nil
	}

	return nil, fmt.Errorf("task title not found. taskTitleID=%s", taskTitleID)
}

func (m MiReps) GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *LimitInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			limitInfo, err := miRep.GetLimitInfo(ctx, limitInfoID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- limitInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	limitInfos := []*LimitInfo{}
loop:
	for {
		select {
		case limitInfo := <-ch:
			if limitInfo == nil {
				continue loop
			}
			limitInfos = append(limitInfos, limitInfo)
		default:
			break loop
		}
	}

	sort.Slice(limitInfos, func(i int, j int) bool {
		return limitInfos[i].UpdatedTime.After(limitInfos[j].UpdatedTime)
	})
	for _, limitInfo := range limitInfos {
		return limitInfo, nil
	}

	return nil, fmt.Errorf("limit not found. limitInfoID=%s", limitInfoID)
}

func (m MiReps) GetStartInfo(ctx context.Context, startInfoID string) (*StartInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *StartInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			startInfo, err := miRep.GetStartInfo(ctx, startInfoID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- startInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	startInfos := []*StartInfo{}
loop:
	for {
		select {
		case startInfo := <-ch:
			if startInfo == nil {
				continue loop
			}
			startInfos = append(startInfos, startInfo)
		default:
			break loop
		}
	}

	sort.Slice(startInfos, func(i int, j int) bool {
		return startInfos[i].UpdatedTime.After(startInfos[j].UpdatedTime)
	})
	for _, limitInfo := range startInfos {
		return limitInfo, nil
	}

	return nil, fmt.Errorf("limit not found. limitInfoID=%s", startInfoID)
}

func (m MiReps) GetEndInfo(ctx context.Context, endInfoID string) (*EndInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *EndInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			endInfo, err := miRep.GetEndInfo(ctx, endInfoID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- endInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	endInfos := []*EndInfo{}
loop:
	for {
		select {
		case endInfo := <-ch:
			if endInfo == nil {
				continue loop
			}
			endInfos = append(endInfos, endInfo)
		default:
			break loop
		}
	}

	sort.Slice(endInfos, func(i int, j int) bool {
		return endInfos[i].UpdatedTime.After(endInfos[j].UpdatedTime)
	})
	for _, limitInfo := range endInfos {
		return limitInfo, nil
	}

	return nil, fmt.Errorf("limit not found. limitInfoID=%s", endInfoID)
}

func (m MiReps) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	existErr := false
	var err error
	wg := &sync.WaitGroup{}
	ch := make(chan *BoardInfo, len(m))
	errch := make(chan error, len(m))
	defer close(ch)
	defer close(errch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			boardInfo, err := miRep.GetBoardInfo(ctx, boardInfoID)
			if err != nil {
				// errch <- err
				return
			}
			ch <- boardInfo
		}(miRep)
	}
	wg.Wait()
errloop:
	for {
		select {
		case e := <-errch:
			err = fmt.Errorf("error at GetAllTasks: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, err
	}

	boardInfos := []*BoardInfo{}
loop:
	for {
		select {
		case boardInfo := <-ch:
			if boardInfo == nil {
				continue loop
			}
			boardInfos = append(boardInfos, boardInfo)
		default:
			break loop
		}
	}

	sort.Slice(boardInfos, func(i int, j int) bool {
		return boardInfos[i].UpdatedTime.After(boardInfos[j].UpdatedTime)
	})
	for _, boardInfo := range boardInfos {
		return boardInfo, nil
	}

	return nil, fmt.Errorf("boardInfo not found. boardInfoID=%s", boardInfoID)
}

func (m MiReps) AddTask(task *Task) error {
	return fmt.Errorf("not implemented") // 追加する対象を特定できないので実装しない
}

func (m MiReps) AddCheckStateInfo(checkStateInfo *CheckStateInfo) error {
	return fmt.Errorf("not implemented") // 追加する対象を特定できないので実装しない
}

func (m MiReps) AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error {
	return fmt.Errorf("not implemented") // 追加する対象を特定できないので実装しない
}

func (m MiReps) AddLimitInfo(limitInfo *LimitInfo) error {
	return fmt.Errorf("not implemented") // 追加する対象を特定できないので実装しない
}

func (m MiReps) AddBoardInfo(boardInfo *BoardInfo) error {
	return fmt.Errorf("not implemented") // 追加する対象を特定できないので実装しない
}

func (m MiReps) GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error) {
	matchTasks := []*Task{}
	taskInfos := map[string]*TaskInfo{}
	tasks, err := m.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			taskInfo, err := m.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if taskInfo == nil || taskInfo.BoardInfo == nil { // taskInfo.BoardInfo==nil????
				continue
			}
			if (taskInfo.BoardInfo.BoardName == query.Board || query.Board == AllBoardName) &&
				(query.Word == "" || strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), strings.ToLower(query.Word))) {
				isMatch := false
				switch query.CheckState {
				case NoCheckOnly:
					isMatch = !taskInfo.CheckStateInfo.IsChecked
				case CheckOnly:
					isMatch = taskInfo.CheckStateInfo.IsChecked
				case All:
					isMatch = true
				}
				if isMatch {
					matchTasks = append(matchTasks, task)
					taskInfos[task.TaskID] = taskInfo
				}
			}
		}
	}

	switch query.SortType {
	case CreatedTimeDesc:
		sort.Slice(matchTasks, func(i int, j int) bool {
			return matchTasks[i].CreatedTime.After(matchTasks[j].CreatedTime)
		})
	case LimitTimeAsc:
		sort.Slice(matchTasks, func(i int, j int) bool {
			if taskInfos[matchTasks[i].TaskID].LimitInfo.Limit == nil && taskInfos[matchTasks[j].TaskID].LimitInfo.Limit == nil {
				return false
			}
			if taskInfos[matchTasks[i].TaskID].LimitInfo.Limit != nil && taskInfos[matchTasks[j].TaskID].LimitInfo.Limit == nil {
				return true
			}
			if taskInfos[matchTasks[i].TaskID].LimitInfo.Limit == nil && taskInfos[matchTasks[j].TaskID].LimitInfo.Limit != nil {
				return false
			}
			limitI := *taskInfos[matchTasks[i].TaskID].LimitInfo.Limit
			limitJ := *taskInfos[matchTasks[j].TaskID].LimitInfo.Limit
			return limitI.Before(limitJ)
		})
	}

	return matchTasks, nil
}

func (m MiReps) GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error) {
	taskInfo := &TaskInfo{}
	var err error
	taskInfo.Task, err = m.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.CheckStateInfo, err = m.GetLatestCheckStateInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.TaskTitleInfo, err = m.GetLatestTaskTitleInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.LimitInfo, err = m.GetLatestLimitInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.StartInfo, err = m.GetLatestStartInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.EndInfo, err = m.GetLatestEndInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	taskInfo.BoardInfo, err = m.GetLatestBoardInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return taskInfo, nil
}

func (m MiReps) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}
	wg := &sync.WaitGroup{}
	ch := make(chan []*kyou.Kyou, len(m))
	defer close(ch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			kyous, err := miRep.GetAllKyous(ctx)
			if err != nil {
				// errch <- err
				return
			}
			ch <- kyous
		}(miRep)
	}
	wg.Wait()
loop:
	for {
		select {
		case collectedKyous := <-ch:
			kyous = append(kyous, collectedKyous...)
		default:
			break loop
		}
	}
	return kyous, nil
}

// 期間を指定してKyouを取得する
func (m MiReps) GetKyousByTime(ctx context.Context, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error) {
	allKyous, err := m.GetAllKyous(ctx)
	if err != nil {
		return nil, err
	}

	matchKyous := []*kyou.Kyou{}
	for _, kyou := range allKyous {
		if kyou.Time.After(startTime) && kyou.Time.Before(endTime) {
			matchKyous = append(matchKyous, kyou)
		}
	}
	return matchKyous, nil
}

func (m MiReps) GetContentHTML(ctx context.Context, id string) (string, error) {
	task, _ := m.GetLatestTaskTitleInfoFromTaskID(ctx, id)
	if task != nil {
		if task.TaskID == id {
			titleInfo, err := m.GetLatestTaskTitleInfoFromTaskID(ctx, task.TaskID)
			if err == nil {
				return `<p>タスク作成:<br/>` + titleInfo.Title + `</p>`, nil
			}
		}
	}

	checkStateInfos, _ := m.GetAllCheckStateInfos(ctx)
	for _, checkStateInfo := range checkStateInfos {
		if checkStateInfo.CheckStateID == id {
			titleInfo, err := m.GetLatestTaskTitleInfoFromTaskID(ctx, checkStateInfo.TaskID)
			if err != nil {
				panic(err)
			}

			if titleInfo != nil {
				if checkStateInfo.IsChecked {
					return `<p>タスクチェック:<br/>` + titleInfo.Title + `</p>`, nil
				} else {
					return `<p>タスク未チェック:<br/>` + titleInfo.Title + `</p>`, nil
				}
			}
		}
	}
	return "", fmt.Errorf("not found kyou %s", id)
}

func (m MiReps) GetPath(ctx context.Context, id string) (string, error) {
	return "Mi", nil
}

func (m MiReps) Delete(id string) error {
	for _, miRep := range m {
		miRep.Delete(id)
	}
	return nil
}

func (m MiReps) Close() error {
	for _, miRep := range m {
		miRep.Close()
	}
	return nil
}

func (m MiReps) Path() string {
	return "Mi"
}

func (m MiReps) RepName() string {
	return "Mi"
}

func (m MiReps) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}
	wg := &sync.WaitGroup{}
	ch := make(chan []*kyou.Kyou, len(m))
	defer close(ch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			kyous, err := miRep.Search(ctx, word)
			if err != nil {
				// errch <- err
				return
			}
			ch <- kyous
		}(miRep)
	}
	wg.Wait()
loop:
	for {
		select {
		case collectedKyous := <-ch:
			kyous = append(kyous, collectedKyous...)
		default:
			break loop
		}
	}
	return kyous, nil
}

// このRepから単語が含まれるKyouを取得する。日時範囲指定
func (m MiReps) SearchByTime(ctx context.Context, word string, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}
	wg := &sync.WaitGroup{}
	ch := make(chan []*kyou.Kyou, len(m))
	defer close(ch)
	for _, miRep := range m {
		wg.Add(1)
		miRep := miRep
		go func(miRep MiRep) {
			defer wg.Done()
			kyous, err := miRep.SearchByTime(ctx, word, startTime, endTime)
			if err != nil {
				// errch <- err
				return
			}
			ch <- kyous
		}(miRep)
	}
	wg.Wait()
loop:
	for {
		select {
		case collectedKyous := <-ch:
			kyous = append(kyous, collectedKyous...)
		default:
			break loop
		}
	}
	return kyous, nil
}
