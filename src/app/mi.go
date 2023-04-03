package mi

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mt3hr/rykv/kyou"
)

const TimeLayout = kyou.TimeLayout

var (
	//go:embed mi/mi/embed
	embedDir embed.FS

	sqlCreateTables      string
	sqlGetTask           string
	sqlGetCheckStateInfo string
	sqlGetTaskTitleInfo  string
	sqlGetLimitInfo      string
	sqlGetBoardInfo      string

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
	sqlDelete           string
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
	sqlGetCheckStateInfoB, err := embedDir.ReadFile("/sql/GetCheckStateInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetCheckStateInfo = string(sqlGetCheckStateInfoB)
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
	sqlDeleteB, err := embedDir.ReadFile("/sql/Delete.sql")
	if err != nil {
		panic(err)
	}
	sqlDelete = string(sqlDeleteB)
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

func NewMiRepSQLite(dbFileName string) (MiRep, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		err = fmt.Errorf("error at open database %s: %w", dbFileName, err)
		return nil, err
	}
	_, err = db.Exec(sqlCreateTables)
	if err != nil {
		err = fmt.Errorf("error at create table to database at %s: %w", dbFileName, err)
		return nil, err
	}
	return &miRepSQLiteImpl{
		filename: dbFileName,
		db:       db,
		m:        &sync.Mutex{},
	}, nil
}

type miRepSQLiteImpl struct {
	filename string
	db       *sql.DB
	m        *sync.Mutex
}

func (m *miRepSQLiteImpl) GetTask(ctx context.Context, taskID string) (*Task, error) {
	statement := fmt.Sprintf(sqlGetTask, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	task := &Task{}
	createdTimeStr := ""
	err := row.Scan(&task.TaskID, &createdTimeStr)
	if err != nil {
		err = fmt.Errorf("error at get task %s: %w", taskID, err)
		return nil, err
	}

	task.CreatedTime, err = time.Parse(TimeLayout, createdTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}
	return task, nil
}

func (m *miRepSQLiteImpl) GetCheckStateInfo(ctx context.Context, checkStateID string) (*CheckStateInfo, error) {
	statement := fmt.Sprintf(sqlGetCheckStateInfo, escapeSQLite(checkStateID))
	row := m.db.QueryRowContext(ctx, statement)

	checkStateInfo := &CheckStateInfo{}

	updatedTimeStr := ""
	err := row.Scan(&checkStateInfo.CheckStateID,
		&checkStateInfo.TaskID,
		&updatedTimeStr,
		&checkStateInfo.IsChecked)
	if err != nil {
		err = fmt.Errorf("error at get check state info %s: %w", checkStateID, err)
		return nil, err
	}

	checkStateInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", checkStateID, err)
		return nil, err
	}
	return checkStateInfo, nil
}

func (m *miRepSQLiteImpl) GetTaskTitleInfo(ctx context.Context, taskTitleID string) (*TaskTitleInfo, error) {
	statement := fmt.Sprintf(sqlGetTaskTitleInfo, escapeSQLite(taskTitleID))
	row := m.db.QueryRowContext(ctx, statement)

	taskTitleInfo := &TaskTitleInfo{}
	updatedTimeStr := ""
	err := row.Scan(&taskTitleInfo.TaskTitleID,
		&taskTitleInfo.TaskID,
		&updatedTimeStr,
		&taskTitleInfo.Title)
	if err != nil {
		err = fmt.Errorf("error at get task title info %s: %w", taskTitleID, err)
		return nil, err
	}

	taskTitleInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskTitleID, err)
		return nil, err
	}
	return taskTitleInfo, nil
}

func (m *miRepSQLiteImpl) GetLimitInfo(ctx context.Context, limitInfoID string) (*LimitInfo, error) {
	statement := fmt.Sprintf(sqlGetLimitInfo, escapeSQLite(limitInfoID))
	row := m.db.QueryRowContext(ctx, statement)

	limitInfo := &LimitInfo{}
	updatedTimeStr := ""
	limitTimeStr := ""
	err := row.Scan(&limitInfo.LimitID,
		&limitInfo.TaskID,
		&updatedTimeStr,
		&limitTimeStr)
	if err != nil {
		err = fmt.Errorf("error at get limit info %s: %w", limitInfoID, err)
		return nil, err
	}

	limitInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", limitInfoID, err)
		return nil, err
	}

	limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse limit time %s: %w", limitInfoID, err)
		return nil, err
	}
	return limitInfo, nil
}

func (m *miRepSQLiteImpl) GetBoardInfo(ctx context.Context, boardInfoID string) (*BoardInfo, error) {
	statement := fmt.Sprintf(sqlGetBoardInfo, escapeSQLite(boardInfoID))
	row := m.db.QueryRowContext(ctx, statement)

	boardInfo := &BoardInfo{}
	updatedTimeStr := ""
	err := row.Scan(&boardInfo.BoardInfoID,
		&boardInfo.TaskID,
		&updatedTimeStr,
		&boardInfo.BoardName)
	if err != nil {
		err = fmt.Errorf("error at get board info %s: %w", boardInfoID, err)
		return nil, err
	}

	boardInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", boardInfoID, err)
		return nil, err
	}

	return boardInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestCheckStateInfoFromTaskID(ctx context.Context, taskID string) (*CheckStateInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestCheckStateFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	checkStateInfo := &CheckStateInfo{}

	updatedTimeStr := ""
	err := row.Scan(&checkStateInfo.CheckStateID,
		&checkStateInfo.TaskID,
		&updatedTimeStr,
		&checkStateInfo.IsChecked)
	if err != nil {
		err = fmt.Errorf("error at get check state info %s: %w", taskID, err)
		return nil, err
	}

	checkStateInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}
	return checkStateInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestTaskTitleInfoFromTaskID(ctx context.Context, taskID string) (*TaskTitleInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestTaskTitleInfoFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	taskTitleInfo := &TaskTitleInfo{}
	updatedTimeStr := ""
	err := row.Scan(&taskTitleInfo.TaskTitleID,
		&taskTitleInfo.TaskID,
		&updatedTimeStr,
		&taskTitleInfo.Title)
	if err != nil {
		err = fmt.Errorf("error at get task title info %s: %w", taskID, err)
		return nil, err
	}

	taskTitleInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}
	return taskTitleInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestLimitInfoFromTaskID(ctx context.Context, taskID string) (*LimitInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestLimitInfoFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	limitInfo := &LimitInfo{}
	updatedTimeStr := ""
	limitTimeStr := ""
	err := row.Scan(&limitInfo.LimitID,
		&limitInfo.TaskID,
		&updatedTimeStr,
		&limitTimeStr)
	if err != nil {
		err = fmt.Errorf("error at get limit info %s: %w", taskID, err)
		return nil, err
	}

	limitInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}

	limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse limit time %s: %w", taskID, err)
		return nil, err
	}
	return limitInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestBoardInfoFromTaskID(ctx context.Context, taskID string) (*BoardInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestBoardInfoFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	boardInfo := &BoardInfo{}
	updatedTimeStr := ""
	err := row.Scan(&boardInfo.BoardInfoID,
		&boardInfo.TaskID,
		&updatedTimeStr,
		&boardInfo.BoardName)
	if err != nil {
		err = fmt.Errorf("error at get board info %s: %w", taskID, err)
		return nil, err
	}

	boardInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}

	return boardInfo, nil
}

func (m *miRepSQLiteImpl) AddTask(task *Task) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlAddTask,
		escapeSQLite(task.TaskID),
		escapeSQLite(task.CreatedTime.Format(TimeLayout)))
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at add task to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}

func (m *miRepSQLiteImpl) AddCheckStateInfo(checkStateInfo *CheckStateInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlAddCheckState,
		escapeSQLite(checkStateInfo.CheckStateID),
		escapeSQLite(checkStateInfo.TaskID),
		escapeSQLite(checkStateInfo.UpdatedTime.Format(TimeLayout)),
		checkStateInfo.IsChecked)
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at add check state to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}

func (m *miRepSQLiteImpl) AddTaskTitleInfo(taskTitleInfo *TaskTitleInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlAddTaskTitleInfo,
		escapeSQLite(taskTitleInfo.TaskTitleID),
		escapeSQLite(taskTitleInfo.TaskID),
		escapeSQLite(taskTitleInfo.UpdatedTime.Format(TimeLayout)),
		escapeSQLite(taskTitleInfo.Title))
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at add task title info to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}

func (m *miRepSQLiteImpl) AddLimitInfo(limitInfo *LimitInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlAddLimitInfo,
		escapeSQLite(limitInfo.LimitID),
		escapeSQLite(limitInfo.TaskID),
		escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
		escapeSQLite(limitInfo.Limit.Format(TimeLayout)))
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at add task limit info to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}

func (m *miRepSQLiteImpl) AddBoardInfo(boardInfo *BoardInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlAddBoardInfo,
		escapeSQLite(boardInfo.BoardInfoID),
		escapeSQLite(boardInfo.TaskID),
		escapeSQLite(boardInfo.UpdatedTime.Format(TimeLayout)),
		escapeSQLite(boardInfo.BoardName))
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at add board info to to database %s: %w", m.filename, err)
		return err
	}
	return nil
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
	m.m.Lock()
	defer m.m.Unlock()
	statement := fmt.Sprintf(sqlDelete,
		escapeSQLite(id),
		escapeSQLite(id),
		escapeSQLite(id),
		escapeSQLite(id),
		escapeSQLite(id))
	_, err := m.db.Exec(statement)
	if err != nil {
		err = fmt.Errorf("error at delete to database %s: %w", m.filename, err)
		return err
	}
	return nil
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

func escapeSQLite(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}
