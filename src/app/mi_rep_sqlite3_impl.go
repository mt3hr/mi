package mi

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mt3hr/rykv/kyou"
)

const TimeLayout = kyou.TimeLayout

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

func (m *miRepSQLiteImpl) SearchTasks(ctx context.Context, word string, query *SearchTaskQuery) ([]*Task, error) {
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

var (
	//go:embed mi/mi/embed
	EmbedDir embed.FS

	sqlCreateTables          string
	sqlGetAllTasks           string
	sqlGetTask               string
	sqlGetCheckStateInfo     string
	sqlGetAllCheckStateInfos string
	sqlGetTaskTitleInfo      string
	sqlGetLimitInfo          string
	sqlGetStartInfo          string
	sqlGetEndInfo            string
	sqlGetBoardInfo          string

	sqlGetLatestCheckStateFromTaskID    string
	sqlGetLatestTaskTitleInfoFromTaskID string
	sqlGetLatestLimitInfoFromTaskID     string
	sqlGetLatestStartInfoFromTaskID     string
	sqlGetLatestEndInfoFromTaskID       string
	sqlGetLatestBoardInfoFromTaskID     string

	sqlAddTask           string
	sqlAddCheckStateInfo string
	sqlAddTaskTitleInfo  string
	sqlAddLimitInfo      string
	sqlAddStartInfo      string
	sqlAddEndInfo        string
	sqlAddBoardInfo      string
	sqlDelete            string
)

func init() {
	Prepare()
}

func Prepare() {
	sqlCreateTablesB, err := EmbedDir.ReadFile("mi/mi/embed/sql/CreateTables.sql")
	if err != nil {
		panic(err)
	}
	sqlCreateTables = string(sqlCreateTablesB)
	sqlGetTaskB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetTask.sql")
	if err != nil {
		panic(err)
	}
	sqlGetTask = string(sqlGetTaskB)
	sqlGetAllTasksB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetAllTasks.sql")
	if err != nil {
		panic(err)
	}
	sqlGetAllTasks = string(sqlGetAllTasksB)

	sqlGetCheckStateInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetCheckStateInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetCheckStateInfo = string(sqlGetCheckStateInfoB)
	sqlGetTaskTitleInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetTaskTitleInfo.sql")
	sqlGetAllCheckStateInfosB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetAllCheckStateInfos.sql")
	if err != nil {
		panic(err)
	}
	sqlGetAllCheckStateInfos = string(sqlGetAllCheckStateInfosB)
	if err != nil {
		panic(err)
	}
	sqlGetTaskTitleInfo = string(sqlGetTaskTitleInfoB)
	sqlGetLimitInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLimitInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLimitInfo = string(sqlGetLimitInfoB)
	sqlGetStartInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetStartInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetStartInfo = string(sqlGetStartInfoB)
	sqlGetEndInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetEndInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetEndInfo = string(sqlGetEndInfoB)
	sqlGetBoardInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetBoardInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetBoardInfo = string(sqlGetBoardInfoB)
	sqlGetLatestCheckStateFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestCheckStateInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestCheckStateFromTaskID = string(sqlGetLatestCheckStateFromTaskIDB)
	sqlGetLatestTaskTitleInfoFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestTaskTitleInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestTaskTitleInfoFromTaskID = string(sqlGetLatestTaskTitleInfoFromTaskIDB)
	sqlGetLatestLimitInfoFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestLimitInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestLimitInfoFromTaskID = string(sqlGetLatestLimitInfoFromTaskIDB)
	sqlGetLatestStartInfoFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestStartInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestStartInfoFromTaskID = string(sqlGetLatestStartInfoFromTaskIDB)
	sqlGetLatestEndInfoFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestEndInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestEndInfoFromTaskID = string(sqlGetLatestEndInfoFromTaskIDB)

	sqlGetLatestBoardInfoFromTaskIDB, err := EmbedDir.ReadFile("mi/mi/embed/sql/GetLatestBoardInfoFromTaskID.sql")
	if err != nil {
		panic(err)
	}
	sqlGetLatestBoardInfoFromTaskID = string(sqlGetLatestBoardInfoFromTaskIDB)

	sqlAddTaskB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddTask.sql")
	if err != nil {
		panic(err)
	}
	sqlAddTask = string(sqlAddTaskB)
	sqlAddCheckStateInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddCheckStateInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddCheckStateInfo = string(sqlAddCheckStateInfoB)
	sqlAddTaskTitleInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddTaskTitleInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddTaskTitleInfo = string(sqlAddTaskTitleInfoB)
	sqlAddLimitInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddLimitInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddLimitInfo = string(sqlAddLimitInfoB)
	sqlAddStartInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddStartInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddStartInfo = string(sqlAddStartInfoB)
	sqlAddEndInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddEndInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddEndInfo = string(sqlAddEndInfoB)
	sqlAddBoardInfoB, err := EmbedDir.ReadFile("mi/mi/embed/sql/AddBoardInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlAddBoardInfo = string(sqlAddBoardInfoB)
	sqlDeleteB, err := EmbedDir.ReadFile("mi/mi/embed/sql/Delete.sql")
	if err != nil {
		panic(err)
	}
	sqlDelete = string(sqlDeleteB)
}

func (m *miRepSQLiteImpl) GetAllTasks(ctx context.Context) ([]*Task, error) {
	tasks := []*Task{}
	statement := sqlGetAllTasks
	rows, err := m.db.QueryContext(ctx, statement)
	if err != nil {
		err = fmt.Errorf("error at get all tasks: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			task := &Task{}
			createdTimeStr := ""
			err := rows.Scan(&task.TaskID, &createdTimeStr)
			if err != nil {
				return nil, err
			}

			task.CreatedTime, err = time.Parse(TimeLayout, createdTimeStr)
			if err != nil {
				err = fmt.Errorf("error at parse time: %w", err)
				return nil, err
			}
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
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
	limitTimeStr := sql.NullString{}
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

	if limitTimeStr.Valid {
		limitInfo.Limit = &time.Time{}
		*limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse limit time %s: %w", limitInfoID, err)
			return nil, err
		}
	}
	return limitInfo, nil
}
func (m *miRepSQLiteImpl) GetStartInfo(ctx context.Context, limitInfoID string) (*StartInfo, error) {
	statement := fmt.Sprintf(sqlGetStartInfo, escapeSQLite(limitInfoID))
	row := m.db.QueryRowContext(ctx, statement)

	limitInfo := &StartInfo{}
	updatedTimeStr := ""
	limitTimeStr := sql.NullString{}
	err := row.Scan(&limitInfo.StartID,
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

	if limitTimeStr.Valid {
		limitInfo.Start = &time.Time{}
		*limitInfo.Start, err = time.Parse(TimeLayout, limitTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse limit time %s: %w", limitInfoID, err)
			return nil, err
		}
	}
	return limitInfo, nil
}

func (m *miRepSQLiteImpl) GetEndInfo(ctx context.Context, limitInfoID string) (*EndInfo, error) {
	statement := fmt.Sprintf(sqlGetEndInfo, escapeSQLite(limitInfoID))
	row := m.db.QueryRowContext(ctx, statement)

	limitInfo := &EndInfo{}
	updatedTimeStr := ""
	limitTimeStr := sql.NullString{}
	err := row.Scan(&limitInfo.EndID,
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

	if limitTimeStr.Valid {
		limitInfo.End = &time.Time{}
		*limitInfo.End, err = time.Parse(TimeLayout, limitTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse limit time %s: %w", limitInfoID, err)
			return nil, err
		}
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
	limitTimeStr := sql.NullString{}
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

	if limitTimeStr.Valid {
		limitInfo.Limit = &time.Time{}
		*limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse limit time %s: %w", taskID, err)
			return nil, err
		}
	}
	return limitInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestStartInfoFromTaskID(ctx context.Context, taskID string) (*StartInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestStartInfoFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	startInfo := &StartInfo{}
	updatedTimeStr := ""
	startTimeStr := sql.NullString{}
	err := row.Scan(&startInfo.StartID,
		&startInfo.TaskID,
		&updatedTimeStr,
		&startTimeStr)
	if err != nil {
		err = fmt.Errorf("error at get start info %s: %w", taskID, err)
		return nil, err
	}

	startInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}

	if startTimeStr.Valid {
		startInfo.Start = &time.Time{}
		*startInfo.Start, err = time.Parse(TimeLayout, startTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse start time %s: %w", taskID, err)
			return nil, err
		}
	}
	return startInfo, nil
}

func (m *miRepSQLiteImpl) GetLatestEndInfoFromTaskID(ctx context.Context, taskID string) (*EndInfo, error) {
	statement := fmt.Sprintf(sqlGetLatestEndInfoFromTaskID, escapeSQLite(taskID))
	row := m.db.QueryRowContext(ctx, statement)

	endInfo := &EndInfo{}
	updatedTimeStr := ""
	endTimeStr := sql.NullString{}
	err := row.Scan(&endInfo.EndID,
		&endInfo.TaskID,
		&updatedTimeStr,
		&endTimeStr)
	if err != nil {
		err = fmt.Errorf("error at get end info %s: %w", taskID, err)
		return nil, err
	}

	endInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
	if err != nil {
		err = fmt.Errorf("error at parse time %s: %w", taskID, err)
		return nil, err
	}

	if endTimeStr.Valid {
		endInfo.End = &time.Time{}
		*endInfo.End, err = time.Parse(TimeLayout, endTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse end time %s: %w", taskID, err)
			return nil, err
		}
	}
	return endInfo, nil
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
	statement := fmt.Sprintf(sqlAddCheckStateInfo,
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
	statement := sqlAddLimitInfo
	var err error
	if limitInfo.Limit != nil {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.LimitID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			escapeSQLite(limitInfo.Limit.Format(TimeLayout)))
	} else {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.LimitID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			sql.NullString{})
	}
	if err != nil {
		err = fmt.Errorf("error at add task limit info to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}

func (m *miRepSQLiteImpl) AddStartInfo(limitInfo *StartInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := sqlAddStartInfo
	var err error
	if limitInfo.Start != nil {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.StartID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			escapeSQLite(limitInfo.Start.Format(TimeLayout)))
	} else {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.StartID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			sql.NullString{})
	}
	if err != nil {
		err = fmt.Errorf("error at add task limit info to to database %s: %w", m.filename, err)
		return err
	}
	return nil
}
func (m *miRepSQLiteImpl) AddEndInfo(limitInfo *EndInfo) error {
	m.m.Lock()
	defer m.m.Unlock()
	statement := sqlAddEndInfo
	var err error
	if limitInfo.End != nil {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.EndID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			escapeSQLite(limitInfo.End.Format(TimeLayout)))
	} else {
		_, err = m.db.Exec(statement,
			escapeSQLite(limitInfo.EndID),
			escapeSQLite(limitInfo.TaskID),
			escapeSQLite(limitInfo.UpdatedTime.Format(TimeLayout)),
			sql.NullString{})
	}
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

func (m *miRepSQLiteImpl) GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error) {
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
					matchTasks = append(matchTasks, task)
					taskInfos[task.TaskID] = taskInfo
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

func (m *miRepSQLiteImpl) GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error) {
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

func (m *miRepSQLiteImpl) GetAllKyous(ctx context.Context) ([]*kyou.Kyou, error) {
	kyous := []*kyou.Kyou{}

	tasks, err := m.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		taskInfo, err := m.GetTaskInfo(ctx, task.TaskID)
		if err != nil {
			return nil, err
		}
		kyous = append(kyous, &kyou.Kyou{
			ID:          taskInfo.Task.TaskID,
			Time:        taskInfo.Task.CreatedTime,
			RepName:     m.RepName(),
			ImageSource: "",
		})
	}

	checkStateInfos, err := m.GetAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}
	for _, checkStateInfo := range checkStateInfos {
		if checkStateInfo.IsChecked {
			kyous = append(kyous, &kyou.Kyou{
				ID:          checkStateInfo.CheckStateID,
				Time:        checkStateInfo.UpdatedTime,
				RepName:     m.RepName(),
				ImageSource: "",
			})
		}
	}
	return kyous, nil
}

func (m *miRepSQLiteImpl) GetKyousByTime(ctx context.Context, startTime time.Time, endTime time.Time) ([]*kyou.Kyou, error) {
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

func (m *miRepSQLiteImpl) GetContentHTML(ctx context.Context, id string) (string, error) {
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
				continue
			}
			if checkStateInfo.IsChecked {
				return `<p>タスクチェック:<br/>` + titleInfo.Title + `</p>`, nil
			}
		}
	}
	return "", fmt.Errorf("not found kyou %s", id)
}

func (m *miRepSQLiteImpl) GetPath(ctx context.Context, id string) (string, error) {
	return m.filename, nil
}

func (m *miRepSQLiteImpl) Delete(id string) error {
	// タスクだった場合、関連情報もDBから消したほうがよくない？と思ったけど別DBにある可能性があるので消さなくていいか
	// いやそんなことないか、RepsでFor回されたときに消せるわ
	//TODO あでも接続されていないDBがあったときめんどいな、どうしよう
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
	return "Mi"
	// base := filepath.Base(m.Path())
	// ext := filepath.Ext(base)
	// withoutExt := base[:len(base)-len(ext)]
	// return withoutExt
}

func (m *miRepSQLiteImpl) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	word = strings.ToLower(word)
	kyous := []*kyou.Kyou{}

	checkStateInfos, err := m.GetAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := m.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	if tasks != nil {
		for _, task := range tasks {
			taskInfo, err := m.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), word) || word == "" || strings.Contains(word, "タスク作成") {
				kyous = append(kyous, &kyou.Kyou{
					ID:          taskInfo.Task.TaskID,
					Time:        taskInfo.Task.CreatedTime,
					RepName:     m.RepName(),
					ImageSource: "",
				})
				for _, checkStateInfo := range checkStateInfos {
					if checkStateInfo.IsChecked && (word == "" || strings.Contains(word, "タスクチェック")) {
						if checkStateInfo.TaskID == taskInfo.Task.TaskID {
							kyous = append(kyous, &kyou.Kyou{
								ID:          checkStateInfo.CheckStateID,
								Time:        checkStateInfo.UpdatedTime,
								RepName:     m.RepName(),
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

func (m *miRepSQLiteImpl) UpdateCache(ctx context.Context) error {
	return nil
}

func (m *miRepSQLiteImpl) GetAllCheckStateInfos(ctx context.Context) ([]*CheckStateInfo, error) {
	checkStateInfos := []*CheckStateInfo{}
	statement := sqlGetAllCheckStateInfos
	rows, err := m.db.QueryContext(ctx, statement)
	if err != nil {
		err = fmt.Errorf("error at get all tasks: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			checkStateInfo := &CheckStateInfo{}
			updatedTimeStr := ""
			err := rows.Scan(&checkStateInfo.CheckStateID,
				&checkStateInfo.TaskID,
				&updatedTimeStr,
				&checkStateInfo.IsChecked)
			if err != nil {
				err = fmt.Errorf("error at get check state info: %w", err)
				return nil, err
			}

			checkStateInfo.UpdatedTime, err = time.Parse(TimeLayout, updatedTimeStr)
			if err != nil {
				err = fmt.Errorf("error at parse time: %w", err)
				return nil, err
			}

			checkStateInfos = append(checkStateInfos, checkStateInfo)
		}
	}
	return checkStateInfos, nil
}

func escapeSQLite(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}
