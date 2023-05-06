package mi

import (
	"context"

	"github.com/mt3hr/rykv/kyou"
	"github.com/mt3hr/rykv/tag"
)

type miRepT struct {
	rep           MiRep
	deleteTagReps tag.DeleteTagReps
}

func (m *miRepT) SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error) {
	return m.rep.SearchTasks(ctx, word, query)
}

func (m *miRepT) GetAllTasks(ctx context.Context) ([]*Task, error) {
	return m.rep.GetAllTasks(ctx)
}

func (m *miRepT) GetTask(ctx context.Context, taskID string) (*Task, error) {
	return m.rep.GetTask(ctx, taskID)
}

func (m *miRepT) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	return m.rep.GetLatestCheckStateInfoFromTaskID(ctx, taskID)
}

func (m *miRepT) GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error) {
	return m.rep.GetLatestTaskTitleInfoFromTaskID(ctx, taskID)
}

func (m *miRepT) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	return m.rep.GetLatestLimitInfoFromTaskID(ctx, taskID)
}

func (m *miRepT) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	return m.rep.GetLatestBoardInfoFromTaskID(ctx, taskID)
}

func (m *miRepT) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	return m.rep.GetCheckStateInfo(ctx, checkStateID)
}

func (m *miRepT) GetTaskTitleInfo(ctx context.Context, taskTitleID string) (*TaskTitleInfo, error) {
	return m.rep.GetTaskTitleInfo(ctx, taskTitleID)
}

func (m *miRepT) GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error) {
	return m.rep.GetLimitInfo(ctx, limitInfoID)
}

func (m *miRepT) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	return m.rep.GetBoardInfo(ctx, boardInfoID)
}

func (m *miRepT) AddTask(task *Task) error {
	return m.rep.AddTask(task)
}

func (m *miRepT) AddCheckStateInfo(checkStateInfo *CheckStateInfo) error {
	return m.rep.AddCheckStateInfo(checkStateInfo)
}

func (m *miRepT) AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error {
	return m.rep.AddTaskTitleInfo(taskTitleInfo)
}

func (m *miRepT) AddLimitInfo(limitInfo *LimitInfo) error {
	return m.rep.AddLimitInfo(limitInfo)
}

func (m *miRepT) AddBoardInfo(boardInfo *BoardInfo) error {
	return m.rep.AddBoardInfo(boardInfo)
}

func (m *miRepT) GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error) {
	return m.rep.GetTasksAtBoard(ctx, query)
}

func (m *miRepT) GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error) {
	return m.rep.GetTaskInfo(ctx, taskID)
}

// このRepからすべてのKyouを取得する
func (m *miRepT) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	return m.rep.GetAllKyous(ctx)
}

// このRepのもつKyouの内容をHTMLで取得する
func (m *miRepT) GetContentHTML(ctx context.Context, id string) (string, error) {
	return m.rep.GetContentHTML(ctx, id)
}

// このRepのもつKyouのPathを取得する
func (m *miRepT) GetPath(ctx context.Context, id string) (string, error) {
	return m.rep.GetPath(ctx, id)
}

// このRepからKyouを削除する
func (m *miRepT) Delete(id string) error {
	return m.deleteTagReps.Delete(id)
}

// このRepを閉じる
func (m *miRepT) Close() error {
	return m.rep.Close()
}

// このRepのPathを取得する
func (m *miRepT) Path() string {
	return m.rep.Path()
}

// このRepの名前を取得する
func (m *miRepT) RepName() string {
	return m.rep.RepName()
}

// このRepから単語が含まれるKyouを取得する
func (m *miRepT) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	return m.rep.Search(ctx, word)
}

// なにか情報をキャッシュしているならば、最新の状態に更新する。
func (m *miRepT) UpdateCache(ctx context.Context) error {
	return m.rep.UpdateCache(ctx)
}

func WrapMiRepT(rep MiRep, deleteTagReps tag.DeleteTagReps) MiRep {
	return &miRepT{
		rep:           rep,
		deleteTagReps: deleteTagReps,
	}
}
