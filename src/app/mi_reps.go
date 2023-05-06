package mi

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mt3hr/rykv/kyou"
)

type MiReps []MiRep

func (m MiReps) UpdateCache(ctx context.Context) error {
	return
}

func (m MiReps) SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error) {
	taskMap := map[string]*Task{}
	for _, miRep := range m {
		tasks, err := miRep.SearchTasks(ctx, word, query)
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			taskMap[task.TaskID] = task
		}
	}
	tasks := []*Task{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m MiReps) GetAllTasks(ctx context.Context) ([]*Task, error) {
	taskMap := map[string]*Task{}
	for _, miRep := range m {
		tasks, err := miRep.GetAllTasks(ctx)
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			taskMap[task.TaskID] = task
		}
	}
	tasks := []*Task{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m MiReps) GetTask(ctx context.Context, taskID string) (*Task, error) {
	for _, miRep := range m {
		task, err := miRep.GetTask(ctx, taskID)
		if err != nil {
			continue
		}
		return task, nil
	}
	return nil, fmt.Errorf("task not found. taskID=%s", taskID)
}

func (m MiReps) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	checkStateInfos := []*CheckStateInfo{}
	for _, miRep := range m {
		checkState, err := miRep.GetLatestCheckStateInfoFromTaskID(ctx, taskID)
		if err != nil {
			continue
		}
		checkStateInfos = append(checkStateInfos, checkState)
	}

	sort.Slice(checkStateInfos, func(i int, j int) bool {
		return checkStateInfos[i].UpdatedTime.After(checkStateInfos[j].UpdatedTime)
	})
	for _, checkStateInfo := range checkStateInfos {
		return checkStateInfo, nil
	}

	return nil, fmt.Errorf("check state not found. taskID=%s", taskID)
}

func (m MiReps) GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error) {
	taskTitleInfos := []*TaskTitleInfo{}
	for _, miRep := range m {
		titleInfo, err := miRep.GetLatestTaskTitleInfoFromTaskID(ctx, taskID)
		if err != nil {
			continue
		}
		taskTitleInfos = append(taskTitleInfos, titleInfo)
	}

	sort.Slice(taskTitleInfos, func(i int, j int) bool {
		return taskTitleInfos[i].UpdatedTime.After(taskTitleInfos[j].UpdatedTime)
	})
	for _, taskTitleInfo := range taskTitleInfos {
		return taskTitleInfo, nil
	}

	return nil, fmt.Errorf("title not found. taskID=%s", taskID)
}

func (m MiReps) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	limitInfos := []*LimitInfo{}
	for _, miRep := range m {
		limitInfo, err := miRep.GetLatestLimitInfoFromTaskID(ctx, taskID)
		if err != nil {
			continue
		}
		limitInfos = append(limitInfos, limitInfo)
	}

	sort.Slice(limitInfos, func(i int, j int) bool {
		return limitInfos[i].UpdatedTime.After(limitInfos[j].UpdatedTime)
	})
	for _, limitInfo := range limitInfos {
		return limitInfo, nil
	}

	return nil, fmt.Errorf("limit not found. taskID=%s", taskID)
}

func (m MiReps) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	boardInfos := []*BoardInfo{}
	for _, miRep := range m {
		boardInfo, err := miRep.GetLatestBoardInfoFromTaskID(ctx, taskID)
		if err != nil {
			continue
		}
		boardInfos = append(boardInfos, boardInfo)
	}

	sort.Slice(boardInfos, func(i int, j int) bool {
		return boardInfos[i].UpdatedTime.After(boardInfos[j].UpdatedTime)
	})
	for _, boardInfo := range boardInfos {
		return boardInfo, nil
	}

	return nil, fmt.Errorf("board not found. taskID=%s", taskID)
}

func (m MiReps) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	checkStateInfos := []*CheckStateInfo{}
	for _, miRep := range m {
		checkStateInfo, err := miRep.GetCheckStateInfo(ctx, checkStateID)
		if err != nil {
			continue
		}
		checkStateInfos = append(checkStateInfos, checkStateInfo)
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
	taskTitleInfos := []*TaskTitleInfo{}
	for _, miRep := range m {
		taskTitleInfo, err := miRep.GetTaskTitleInfo(ctx, taskTitleID)
		if err != nil {
			continue
		}
		taskTitleInfos = append(taskTitleInfos, taskTitleInfo)
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
	limitInfos := []*LimitInfo{}
	for _, miRep := range m {
		limitInfo, err := miRep.GetLimitInfo(ctx, limitInfoID)
		if err != nil {
			continue
		}
		limitInfos = append(limitInfos, limitInfo)
	}

	sort.Slice(limitInfos, func(i int, j int) bool {
		return limitInfos[i].UpdatedTime.After(limitInfos[j].UpdatedTime)
	})
	for _, limitInfo := range limitInfos {
		return limitInfo, nil
	}

	return nil, fmt.Errorf("limit not found. limitInfoID=%s", limitInfoID)
}

func (m MiReps) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	boardInfos := []*BoardInfo{}
	for _, miRep := range m {
		boardInfo, err := miRep.GetBoardInfo(ctx, boardInfoID)
		if err != nil {
			continue
		}
		boardInfos = append(boardInfos, boardInfo)
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
			if taskInfo.BoardInfo.BoardName == query.Board &&
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
	taskInfo.BoardInfo, err = m.GetLatestBoardInfoFromTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return taskInfo, nil
}

func (m MiReps) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}
	for _, miRep := range m {
		kyou, err := miRep.GetAllKyous(ctx)
		if err != nil {
			continue
		}
		kyous = append(kyous, kyou...)
	}
	return kyous, nil
}

func (m MiReps) GetContentHTML(ctx context.Context, id string) (string, error) {
	for _, miRep := range m {
		contentHTML, err := miRep.GetContentHTML(ctx, id)
		if err != nil {
			continue
		}
		return contentHTML, nil
	}
	return "", fmt.Errorf("not found kyou id=%s", id)
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
	for _, miRep := range m {
		kyou, err := miRep.Search(ctx, word)
		if err != nil {
			continue
		}
		kyous = append(kyous, kyou...)
	}
	return kyous, nil
}

func (m MiReps) UpdateCache() error {
	return nil
}
