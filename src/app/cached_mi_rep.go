package mi

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mt3hr/rykv/kyou"
)

func NewCachedMiRep(miRep MiRep) MiRep {
	result := &cachedMiRep{
		miRep: miRep,
		m:     &sync.Mutex{},
	}
	// result.UpdateCache(context.Background())
	return result
}

type cachedMiRep struct {
	cachedTaskInfo map[string]*TaskInfo
	miRep          MiRep
	m              *sync.Mutex
}

func (c *cachedMiRep) GetAllTasks(ctx context.Context) ([]*Task, error) {
	return c.miRep.GetAllTasks(ctx)
}

func (c *cachedMiRep) GetTask(ctx context.Context, taskID string) (*Task, error) {
	return c.miRep.GetTask(ctx, taskID)
}

func (c *cachedMiRep) GetAllCheckStateInfos(ctx context.Context) ([]*CheckStateInfo, error) {
	return c.miRep.GetAllCheckStateInfos(ctx)
}

func (c *cachedMiRep) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.CheckStateInfo, nil
}

func (c *cachedMiRep) GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.TaskTitleInfo, nil
}

func (c *cachedMiRep) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.LimitInfo, nil
}

func (c *cachedMiRep) GetLatestStartInfoFromTaskID(ctx context.Context, taskID string) (*StartInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.StartInfo, nil
}

func (c *cachedMiRep) GetLatestEndInfoFromTaskID(ctx context.Context, taskID string) (*EndInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.EndInfo, nil
}

func (c *cachedMiRep) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo.BoardInfo, nil
}

func (c *cachedMiRep) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	return c.miRep.GetCheckStateInfo(ctx, checkStateID)
}

func (c *cachedMiRep) GetTaskTitleInfo(ctx context.Context, taskTitleID string) (*TaskTitleInfo, error) {
	return c.miRep.GetTaskTitleInfo(ctx, taskTitleID)
}

func (c *cachedMiRep) GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error) {
	return c.miRep.GetLimitInfo(ctx, limitInfoID)
}

func (c *cachedMiRep) GetStartInfo(ctx context.Context, startInfoID string) (*StartInfo, error) {
	return c.miRep.GetStartInfo(ctx, startInfoID)
}

func (c *cachedMiRep) GetEndInfo(ctx context.Context, endInfoID string) (*EndInfo, error) {
	return c.miRep.GetEndInfo(ctx, endInfoID)
}

func (c *cachedMiRep) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	return c.miRep.GetBoardInfo(ctx, boardInfoID)
}

func (c *cachedMiRep) AddTask(task *Task) error {
	err := c.miRep.AddTask(task)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddCheckStateInfo(checkStateInfo *CheckStateInfo) error {
	err := c.miRep.AddCheckStateInfo(checkStateInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error {
	err := c.miRep.AddTaskTitleInfo(taskTitleInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddLimitInfo(limitInfo *LimitInfo) error {
	err := c.miRep.AddLimitInfo(limitInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddStartInfo(startInfo *StartInfo) error {
	err := c.miRep.AddStartInfo(startInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddEndInfo(endInfo *EndInfo) error {
	err := c.miRep.AddEndInfo(endInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) AddBoardInfo(boardInfo *BoardInfo) error {
	err := c.miRep.AddBoardInfo(boardInfo)
	if err != nil {
		return err
	}
	return c.UpdateCache(context.Background())
}

func (c *cachedMiRep) GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error) {
	matchTasks := []*Task{}
	for _, taskInfo := range c.cachedTaskInfo {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if (taskInfo.BoardInfo.BoardName == query.Board || query.Board == AllBoardName) &&
				strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), strings.ToLower(query.Word)) {
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
					matchTasks = append(matchTasks, taskInfo.Task)
				}
			}
		}
	}
	return matchTasks, nil
}

func (c *cachedMiRep) GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error) {
	taskInfo, exist := c.cachedTaskInfo[taskID]
	if !exist {
		err := fmt.Errorf("not exist task %s.", taskID)
		return nil, err
	}
	return taskInfo, nil
}

func (c *cachedMiRep) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}

	tasks, err := c.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		taskInfo, err := c.GetTaskInfo(ctx, task.TaskID)
		if err != nil {
			return nil, err
		}
		kyous = append(kyous, &kyou.Kyou{
			ID:          taskInfo.Task.TaskID,
			Time:        taskInfo.Task.CreatedTime,
			RepName:     c.RepName(),
			ImageSource: "",
		})
	}

	checkStateInfos, err := c.GetAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}
	for _, checkStateInfo := range checkStateInfos {
		if checkStateInfo.IsChecked {
			kyous = append(kyous, &kyou.Kyou{
				ID:          checkStateInfo.CheckStateID,
				Time:        checkStateInfo.UpdatedTime,
				RepName:     c.RepName(),
				ImageSource: "",
			})
		}
	}
	return kyous, nil
}

func (c *cachedMiRep) GetKyousByTime(ctx context.Context, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error) {
	allKyous, err := c.GetAllKyous(ctx)
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

func (c *cachedMiRep) GetContentHTML(ctx context.Context, id string) (string, error) {
	return c.miRep.GetContentHTML(ctx, id)
}

func (c *cachedMiRep) GetPath(ctx context.Context, id string) (string, error) {
	return c.miRep.GetPath(ctx, id)
}

func (c *cachedMiRep) Delete(id string) error {
	return c.miRep.Delete(id)
}

func (c *cachedMiRep) Close() error {
	return c.miRep.Close()
}

func (c *cachedMiRep) Path() string {
	return c.miRep.Path()
}

func (c *cachedMiRep) RepName() string {
	return c.miRep.RepName()
}

func (c *cachedMiRep) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	word = strings.ToLower(word)
	kyous := []*kyou.Kyou{}

	checkStateInfos, err := c.GetAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	for _, taskInfo := range c.cachedTaskInfo {
		tasks = append(tasks, taskInfo.Task)
	}

	if tasks != nil {
		for _, task := range tasks {
			taskInfo, err := c.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), word) || word == "" || strings.Contains(word, "タスク作成") {
				kyous = append(kyous, &kyou.Kyou{
					ID:          taskInfo.Task.TaskID,
					Time:        taskInfo.Task.CreatedTime,
					RepName:     c.RepName(),
					ImageSource: "",
				})
				for _, checkStateInfo := range checkStateInfos {
					if checkStateInfo.IsChecked && (word == "" || strings.Contains(word, "タスクチェック")) {
						if checkStateInfo.TaskID == taskInfo.Task.TaskID {
							kyous = append(kyous, &kyou.Kyou{
								ID:          checkStateInfo.CheckStateID,
								Time:        checkStateInfo.UpdatedTime,
								RepName:     c.RepName(),
								ImageSource: "",
							})
						}
					}
				}
			}
		}
	}
	return kyous, nil
}

func (c *cachedMiRep) SearchByTime(ctx context.Context, word string, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error) {
	word = strings.ToLower(word)
	kyous := []*kyou.Kyou{}

	checkStateInfos, err := c.GetAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	for _, taskInfo := range c.cachedTaskInfo {
		tasks = append(tasks, taskInfo.Task)
	}

	if tasks != nil {
		for _, task := range tasks {
			taskInfo, err := c.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), word) || word == "" || strings.Contains(word, "タスク作成") && (taskInfo.Task.CreatedTime.After(startTime) && taskInfo.Task.CreatedTime.Before(endTime)) {
				kyous = append(kyous, &kyou.Kyou{
					ID:          taskInfo.Task.TaskID,
					Time:        taskInfo.Task.CreatedTime,
					RepName:     c.RepName(),
					ImageSource: "",
				})
				for _, checkStateInfo := range checkStateInfos {
					if checkStateInfo.IsChecked && (word == "" || strings.Contains(word, "タスクチェック")) && (checkStateInfo.UpdatedTime.After(startTime) && checkStateInfo.UpdatedTime.Before(endTime)) {
						if checkStateInfo.TaskID == taskInfo.Task.TaskID {
							kyous = append(kyous, &kyou.Kyou{
								ID:          checkStateInfo.CheckStateID,
								Time:        checkStateInfo.UpdatedTime,
								RepName:     c.RepName(),
								ImageSource: "",
							})
						}
					}
				}
			}
		}
	}
	return kyous, nil
}

func (c *cachedMiRep) SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error) {
	matchTasks := []*Task{}
	taskInfos := map[string]*TaskInfo{}
	tasks := []*Task{}
	for _, taskInfo := range c.cachedTaskInfo {
		tasks = append(tasks, taskInfo.Task)
	}

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			taskInfo, err := c.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if (taskInfo.BoardInfo.BoardName == query.Board || query.Board == AllBoardName) &&
				strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), strings.ToLower(query.Word)) {
				isMatch := false
				switch query.CheckState {
				case NoCheckOnly:
					isMatch = !taskInfo.CheckStateInfo.IsChecked
				case CheckOnly:
					isMatch = taskInfo.CheckStateInfo.IsChecked
				case All:
					isMatch = true
				}

				if strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), word) {
					if isMatch {
						matchTasks = append(matchTasks, task)
						taskInfos[task.TaskID] = taskInfo
					}
				}
			}
		}
	}

	/* 外でやるのでけします
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
			return limitI.After(limitJ)
		})
	}
	*/

	return matchTasks, nil
}

func (c *cachedMiRep) UpdateCache(ctx context.Context) error {
	// c.m.Lock()
	// defer c.m.Unlock()
	allTasks, err := c.GetAllTasks(ctx)
	if err != nil {
		return err
	}

	c.cachedTaskInfo = map[string]*TaskInfo{}
	allTaskInfos := map[string]*TaskInfo{}
	for _, task := range allTasks {
		taskInfo, err := c.miRep.GetTaskInfo(ctx, task.TaskID)
		if err != nil {
			return err
		}
		allTaskInfos[task.TaskID] = taskInfo
	}
	c.cachedTaskInfo = allTaskInfos

	return ctx.Err()
}
