package mi

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
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

var (
	//go:embed mi/mi/embed
	embedDir embed.FS

	sqlCreateTables          string
	sqlGetAllTasks           string
	sqlGetTask               string
	sqlGetCheckStateInfo     string
	sqlGetAllCheckStateInfos string
	sqlGetTaskTitleInfo      string
	sqlGetLimitInfo          string
	sqlGetBoardInfo          string

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
	sqlGetAllTasksB, err := embedDir.ReadFile("/sql/GetAllTasks.sql")
	if err != nil {
		panic(err)
	}
	sqlGetAllTasks = string(sqlGetAllTasksB)

	sqlGetCheckStateInfoB, err := embedDir.ReadFile("/sql/GetCheckStateInfo.sql")
	if err != nil {
		panic(err)
	}
	sqlGetCheckStateInfo = string(sqlGetCheckStateInfoB)
	sqlGetTaskTitleInfoB, err := embedDir.ReadFile("/sql/GetTaskTitleInfo.sql")
	sqlGetAllCheckStateInfosB, err := embedDir.ReadFile("/sql/GetAllCheckStateInfos.sql")
	if err != nil {
		panic(err)
	}
	sqlGetAllCheckStateInfos = string(sqlGetAllCheckStateInfosB)
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
		*limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr.String)
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
		*limitInfo.Limit, err = time.Parse(TimeLayout, limitTimeStr.String)
		if err != nil {
			err = fmt.Errorf("error at parse limit time %s: %w", taskID, err)
			return nil, err
		}
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

func (m *miRepSQLiteImpl) GetTasksAtBoard(ctx context.Context, query *SearchTaskQuery) ([]*Task, error) {
	//TODO タグによる絞り込みはここからではできないのでhandlerあたりで絞って
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

	return nil, nil
}

func (m *miRepSQLiteImpl) GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error) {
	taskInfo := &TaskInfo{}
	var err error
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

	checkStateInfos, err := m.getAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}
	for _, checkStateInfo := range checkStateInfos {
		kyous = append(kyous, &kyou.Kyou{
			ID:          checkStateInfo.CheckStateID,
			Time:        checkStateInfo.UpdatedTime,
			RepName:     m.RepName(),
			ImageSource: "",
		})
	}
	return kyous, nil
}

func (m *miRepSQLiteImpl) GetContentHTML(ctx context.Context, id string) (string, error) {
	tasks, _ := m.GetAllTasks(ctx)
	if tasks != nil {
		for _, task := range tasks {
			taskInfo, err := m.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return "", err
			}
			return `<p>タスク作成: ` + taskInfo.TaskTitleInfo.Title + `</p>`, nil
		}
	}

	checkStateInfos, _ := m.getAllCheckStateInfos(ctx)
	if checkStateInfos != nil {
		for _, checkStateInfo := range checkStateInfos {
			taskInfo, err := m.GetTaskInfo(ctx, checkStateInfo.TaskID)
			if err != nil {
				return "", err
			}
			if taskInfo.CheckStateInfo.IsChecked {
				return `<p>☑` + taskInfo.TaskTitleInfo.Title + `</p>`, nil
			} else {
				return `<p>□` + taskInfo.TaskTitleInfo.Title + `</p>`, nil
			}
		}
	}
	return "", fmt.Errorf("not found kyou %s", id)
}

func (m *miRepSQLiteImpl) GetPath(ctx context.Context, id string) (string, error) {
	return m.filename, nil
}

func (m *miRepSQLiteImpl) Delete(ctx context.Context, id string) error {
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
	base := filepath.Base(m.Path())
	ext := filepath.Ext(base)
	withoutExt := base[:len(base)-len(ext)]
	return withoutExt
}

func (m *miRepSQLiteImpl) Search(ctx context.Context, word string) ([]*kyou.Kyou, error) {
	word = strings.ToLower(word)
	kyous := []*kyou.Kyou{}

	checkStateInfos, err := m.getAllCheckStateInfos(ctx)
	if err != nil {
		return nil, err
	}

	tasks, _ := m.GetAllTasks(ctx)
	if tasks != nil {
		for _, task := range tasks {
			taskInfo, err := m.GetTaskInfo(ctx, task.TaskID)
			if err != nil {
				return nil, err
			}
			if strings.Contains(strings.ToLower(taskInfo.TaskTitleInfo.Title), word) {
				kyous = append(kyous, &kyou.Kyou{
					ID:          taskInfo.Task.TaskID,
					Time:        taskInfo.Task.CreatedTime,
					RepName:     m.RepName(),
					ImageSource: "",
				})
			}
			for _, checkStateInfo := range checkStateInfos {
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
	return kyous, nil
}

func (m *miRepSQLiteImpl) UpdateCache() error {
	return nil
}

func (m *miRepSQLiteImpl) getAllCheckStateInfos(ctx context.Context) ([]*CheckStateInfo, error) {
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