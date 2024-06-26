package mi

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mt3hr/rykv/kyou"
)

const AllBoardName = "All"

type MiRep interface {
	GetAllTasks(ctx context.Context) ([]*Task, error)
	GetTask(ctx context.Context, taskID string) (*Task, error)
	GetAllCheckStateInfos(ctx context.Context) ([]*CheckStateInfo, error)

	GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error)
	GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error)
	GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error)
	GetLatestStartInfoFromTaskID(ctx context.Context, taskID string) (*StartInfo, error)
	GetLatestEndInfoFromTaskID(ctx context.Context, taskID string) (*EndInfo, error)
	GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error)

	GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error)
	GetTaskTitleInfo(ctx context.Context, taskTitleID string) (*TaskTitleInfo, error)
	GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error)
	GetStartInfo(ctx context.Context, startInfoID string) (*StartInfo, error)
	GetEndInfo(ctx context.Context, endInfoID string) (*EndInfo, error)
	GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error)

	AddTask(task *Task) error
	AddCheckStateInfo(checkStateInfo *CheckStateInfo) error
	AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error
	AddLimitInfo(limitInfo *LimitInfo) error
	AddStartInfo(startInfo *StartInfo) error
	AddEndInfo(endInfo *EndInfo) error
	AddBoardInfo(boardInfo *BoardInfo) error

	GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error)
	GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error)

	GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error)
	GetKyousByTime(ctx context.Context, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error)
	GetContentHTML(ctx context.Context, id string) (string, error)
	GetPath(ctx context.Context, id string) (string, error)
	Delete(id string) error
	Close() error
	Path() string
	RepName() string
	Search(ctx context.Context, word string) ([]*kyou.Kyou, error)
	SearchByTime(ctx context.Context, word string, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error)
	SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error)
	UpdateCache(ctx context.Context) error
}
