package mi

import (
	"context"
	"database/sql"
	"embed"
	"path/filepath"
	"sync"
	"time"

	"github.com/mt3hr/rykv/kyou"
)

var (
	//go:embed mi/mi/embed
	embedDir embed.FS

	sqlCreateTables     string
	sqlGetTask          string
	sqlGetCheckState    string
	sqlGetTaskTitleInfo string
	sqlGetLimitInfo     string
	sqlGetBoardInfo     string

	sqlGetLatestTaskFromTaskID          string
	sqlGetLatestCheckStateFromTaskID    string
	sqlGetLatestTaskTitleInfoFromTaskID string
	sqlGetLatestLimitInfoFromTaskID     string
	sqlGetLatestBoardInfoFromTaskID     string

	sqlAddTask          string
	sqlAddCheckState    string
	sqlAddTaskTitleInfo string
	sqlAddLimitInfo     string
	sqlAddBoardInfo     string
)

func init() {
	sqlCreateTablesB, err := embedDir.ReadFile("/sql/CreateTables.sql")
	if err != nil {
		panic(err)
	}
	sqlCreateTables = string(sqlCreateTablesB)
	sqlGetTaskB, err := embedDir.ReadFile("/sql/GetTask.sql")
	if err != nil {
		panic(err)
	}
	sqlGetTask = string(sqlGetTaskB)
	sqlGetCheckStateB, err := embedDir.ReadFile("/sql/GetCheckState.sql")
	if err != nil {
		panic(err)
	}
	sqlGetCheckState = string(sqlGetCheckStateB)
	sqlGetTaskTitleInfoB, err := embedDir.ReadFile("/sql/GetTaskTitleInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetTaskTitleInfo = string(sqlGetTaskTitleInfoB)
	sqlGetLimitInfoB, err := embedDir.ReadFile("/sql/GetLimitInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLimitInfo = string(sqlGetLimitInfoB)
	sqlGetBoardInfoB, err := embedDir.ReadFile("/sql/GetBoardInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetBoardInfo = string(sqlGetBoardInfoB)

	sqlGetLatestTaskFromTaskIDB, err := embedDir.ReadFile("/sql/GetLatestTaskFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestTaskFromTaskID = string(sqlGetLatestTaskFromTaskIDB)
	sqlGetLatestCheckStateFromTaskIDB, err := embedDir.ReadFile("/sql/GetLatestCheckStateFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestCheckStateFromTaskID = string(sqlGetLatestCheckStateFromTaskIDB)
	sqlGetLatestTaskTitleInfoFromTaskIDB, err := embedDir.ReadFile("/sql/GetLatestTaskTitleInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestTaskTitleInfoFromTaskID = string(sqlGetLatestTaskTitleInfoFromTaskIDB)
	sqlGetLatestLimitInfoFromTaskIDB, err := embedDir.ReadFile("/sql/GetLatestLimitInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestLimitInfoFromTaskID = string(sqlGetLatestLimitInfoFromTaskIDB)
	sqlGetLatestBoardInfoFromTaskIDB, err := embedDir.ReadFile("/sql/GetLatestBoardInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestBoardInfoFromTaskID = string(sqlGetLatestBoardInfoFromTaskIDB)

	sqlAddTaskB, err := embedDir.ReadFile("/sql/AddTask.sql")
	if err != nil {
		panic(err)
	}
	sqlAddTask = string(sqlAddTaskB)
	sqlAddCheckStateB, err := embedDir.ReadFile("/sql/AddCheckState.sql")
	if err != nil {
		panic(err)
	}
	sqlAddCheckState = string(sqlAddCheckStateB)
	sqlAddTaskTitleInfoB, err := embedDir.ReadFile("/sql/AddTaskTitleInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddTaskTitleInfo = string(sqlAddTaskTitleInfoB)
	sqlAddLimitInfoB, err := embedDir.ReadFile("/sql/AddLimitInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddLimitInfo = string(sqlAddLimitInfoB)
	sqlAddBoardInfoB, err := embedDir.ReadFile("/sql/AddBoardInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddBoardInfo = string(sqlAddBoardInfoB)
}

type MiRep interface {
	GetTask(ctx context.Context, taskID string) (*Task, error)
	GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error)
	GetTaskTitleInfo(cxt context.Context, taskTitleID string) (*TaskTitleInfo, error)
	GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error)
	GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error)

	GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error)
	GetLatestTaskTitleInfoFromTaskID(cxt context.Context, taskID string) (*TaskTitleInfo, error)
	GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error)
	GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error)

	AddTask(task *Task) error
	AddCheckStateInfo(checkStateInfo *CheckStateInfo) error
	AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error
	AddLimitInfo(limitInfo *LimitInfo) error
	AddBoardInfo(boardInfo *BoardInfo) error

	GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error)
	GetContentHTML(ctx context.Context, id string) (string, error)
	GetPath(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string) error
	Close() error
	Path() string
	RepName() string
	Search(ctx context.Context, word string) ([]*kyou.Kyou, error)
	UpdateCache() error
}

type Task struct {
	TaskID      string    `json:"task_id"`
	CreatedTime time.Time `json:"created_time"`
}

type CheckStateInfo struct {
	CheckStateID string    `json:"check_state_id"`
	TaskID       string    `json:"task_id"`
	UpdatedTime  time.Time `json:"updated_time"`
	IsChecked    bool      `json:"is_checked"`
}

type TaskTitleInfo struct {
	TaskTitleID string    `json:"task_title_id"`
	TaskID      string    `json:"task_id"`
	UpdatedTime time.Time `json:"updated_time"`
	Title       string    `json:"title"`
}

type LimitInfo struct {
	LimitID     string    `json:"limit_id"`
	TaskID      string    `json:"task_id"`
	UpdatedTime time.Time `json:"updated_time"`
	Limit       time.Time `json:"limit"`
}

type BoardInfo struct {
	BoardInfoID string    `json:"board_info_id"`
	TaskID      string    `json:"task_id"`
	UpdatedTime time.Time `json:"updated_time"`
	BoardName   string    `json:"board_name"`
}

type miRepSQLiteImpl struct {
	filename string
	db       *sql.DB
	m        *sync.Mutex
}

func (m *miRepSQLiteImpl) GetTask(ctx context.Context, taskID string) (*Task, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetTaskTitleInfo(cxt context.Context, taskTitleID string) (*TaskTitleInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetLatestTaskTitleInfoFromTaskID(cxt context.Context, taskID string) (*TaskTitleInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) AddTask(task *Task) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) AddCheckStateInfo(checkStateInfo *CheckStateInfo) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) AddLimitInfo(limitInfo *LimitInfo) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) AddBoardInfo(boardInfo *BoardInfo) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetContentHTML(ctx context.Context, id string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) GetPath(ctx context.Context, id string) (string, error) {
	return m.filename, nil
}

func (m *miRepSQLiteImpl) Delete(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) Close() error {
	return m.db.Close()
}

func (m *miRepSQLiteImpl) Path() string {
	return m.filename
}

func (m *miRepSQLiteImpl) RepName() string {
	base := filepath.Base(m.Path())
	ext := filepath.Ext(base)
	withoutExt := base[:len(base)-len(ext)]
	return withoutExt
}

func (m *miRepSQLiteImpl) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	panic("not implemented") // TODO: Implement
}

func (m *miRepSQLiteImpl) UpdateCache() error {
	return nil
}
