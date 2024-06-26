package miapp

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	mi "github.com/mt3hr/mi/src/app"
	"github.com/mt3hr/rykv/kyou"
	"github.com/mt3hr/rykv/registrep"
	"github.com/mt3hr/rykv/tag"
	"github.com/mt3hr/rykv/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const get_board_struct_address = "/api/get_board_struct"
const get_tag_struct_address = "/api/get_tag_struct"
const add_task_address = "/api/add_task"
const update_task_address = "/api/update_task"
const delete_task_address = "/api/delete_task"
const get_task_address = "/api/get_task"
const get_tasks_from_board_address = "/api/get_board_task"
const add_tag_address = "/api/add_tag"
const add_text_address = "/api/add_text"
const get_tags_related_task_address = "/api/get_task_tag"
const get_texts_related_task_address = "/api/get_task_text"
const get_tag_address = "/api/get_tag"
const get_text_address = "/api/get_text"
const delete_tag_address = "/api/delete_tag"
const delete_text_address = "/api/delete_text"
const get_tag_names_address = "/api/get_tag_names"
const get_board_names_address = "/api/get_board_names"
const get_application_config_address = "/api/get_application_config"
const get_board_struct_method = "post"
const get_tag_struct_method = "post"
const add_task_method = "post"
const update_task_method = "post"
const delete_task_method = "post"
const get_task_method = "post"
const get_tasks_from_board_method = "post"
const add_tag_method = "post"
const add_text_method = "post"
const get_tags_related_task_method = "post"
const get_texts_related_task_method = "post"
const delete_tag_method = "post"
const delete_text_method = "post"
const get_tag_method = "post"
const get_text_method = "post"
const get_tag_names_method = "post"
const get_board_names_method = "post"
const get_application_config_method = "post"

var (
	htmlFS embed.FS = mi.EmbedDir

	ConfigFileName     = ""        // 使用するコンフィグファイルのpath
	LoadedConfig       = &Config{} // loadConfigで読み込まれたコンフィグ
	LoadedRepositories = &Repositories{}
)

func PersistentPreRun(_ *cobra.Command, _ []string) {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = loadBoardStruct()
	if err != nil {
		log.Fatal(err)
	}
	err = loadTagStructFromFile()
	if err != nil {
		log.Fatal(err)
	}
	LoadedConfig.ApplicationConfig.HiddenTags = append(LoadedConfig.ApplicationConfig.HiddenTags, kyou.DeletedTagName)
}

// Config .
// コンフィグファイルのデータモデル
type Config struct {
	Repositories struct {
		MiReps []*MiRepInfo
		MiRep  *MiRepInfo

		TagRep  *tag.TagRepInfo
		TextRep *text.TextRepInfo

		TagReps  []*tag.TagRepInfo
		TextReps []*text.TextRepInfo
	}

	ServerConfig struct {
		LocalOnly bool
		Address   string
		TLS       struct {
			Enable   bool
			CertFile string
			KeyFile  string
		}
		EnableDeleteAction bool
	}

	ApplicationConfig ApplicationConfig `yaml:"ApplicationConfig" json:"application_config"`
}

type ApplicationConfig struct {
	HiddenTags  []string `json:"hidden_tags"`
	UnCheckTags []string `json:"un_check_tags"`

	BoardStruct      interface{} `yaml:"BoardStruct" json:"board_struct"`
	TagStruct        interface{} `yaml:"TagStruct" json:"tag_struct"`
	DefaultBoardName string      `yaml:"DefaultBoardName" json:"default_board_name"`
	EnableHotReload  bool        `yaml:"EnableHotReload" json:"enable_hot_reload"`
}

func getConfigFile() string {
	return ConfigFileName
}
func getConfig() *Config {
	return LoadedConfig
}
func getConfigName() string {
	return "mi_config"
}
func getConfigExt() string {
	return ".yaml"
}

// CreateDefaultConfigYAML .
// デフォルトのYAMLコンフィグを返します。
// エクスポートされているのはidfから呼び出されるためです
func CreateDefaultConfigYAML() string {
	return `ServerConfig:
  # trueにするとlocalhost以外からのリクエストをブロックします。
  LocalOnly: true

  # mi server でサーバーをたてるときに使うアドレス
  Address: ":2734"

  # 設定するとhttpsで接続するようになります
  TLS:
    Enable: false
    CertFile: ""
    KeyFile: ""

  # 削除機能を有効化するかどうか
  # 有効無効に関わらず削除済みの情報は非表示となります
  # 削除機能はHiddenTags機能で実現されています。
  # 削除済みの情報を閲覧したい場合はTagStructに"deleted"を追加してください
  EnableDeleteAction: false

ApplicationConfig:
  # タスクに更新があったときに板をリロードするかどうか
  EnableHotReload: false

  # 読み込み時にチェックを入れないTag
  UnCheckTags: []

  # ここに記されたタグのついた情報は、チェックが入っていない限り検索結果に反映されません。
  # UncheckTagsと組み合わせて使います。
  # 削除機能もこの機能で実現されています。
  HiddenTags: []

  # Device, Type, Rep, Tagの階層構造の設定。
  BoardStruct:
    All: board
    Inbox: board
  TagStruct: 
    no tag: tag

  # デフォルトの板名
  DefaultBoardName: "Inbox"

Repositories:
  # タスク情報の保存先データベースファイル
  MiRep:
    type: mi_db
    file: $HOME/Mi.db

  # タスク情報データベースファイル郡
  MiReps:
  - type: mi_db
    file: $HOME/Mi.db
  
  # タグ記録時の保存先データベースファイル
  TagRep:
    type: db
    file: $HOME/Tag.db

  # タグ情報源データベースファイル郡
  TagReps:
  - type: db
    file: $HOME/Tag.db

  # テキスト記録時の保存先データベースファイル
  TextRep:
    type: db
    file: $HOME/Text.db

  # テキスト情報源データベースファイル郡
  TextReps:
  - type: db
    file: $HOME/Text.db
`
}

// BoardStructの順番を保証するために
func loadBoardStruct() error {
	configOpt := getConfigFile()
	configName := getConfigName()
	configExt := getConfigExt()
	configPaths := []string{}
	configFileName := ""
	var b []byte

	if configOpt != "" {
		// コンフィグファイルが明示的に指定された場合はそれを
		configPaths = append(configPaths, configOpt)
	} else {
		// 実行ファイルの親ディレクトリ、カレントディレクトリ、ホームディレクトリの順に
		exe, err := os.Executable()
		if err != nil {
			err = fmt.Errorf("error at get executable file path: %w", err)
			log.Printf(err.Error())
		} else {
			configPaths = append(configPaths, filepath.Join(filepath.Dir(exe), configName+configExt))
		}

		configPaths = append(configPaths, filepath.Join(".", configName+configExt))

		home, err := homedir.Dir()
		if err != nil {
			err = fmt.Errorf("error at get user home directory: %w", err)
			log.Printf(err.Error())
		} else {
			configPaths = append(configPaths, filepath.Join(home, configName+configExt))
		}
	}

	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); err == nil {
			configFileName = configPath
			break
		}
	}

	b, err := os.ReadFile(configFileName)
	if err != nil {
		err = fmt.Errorf("error at read file %s: %w", configFileName, err)
		return err
	}

	m := yaml.MapSlice{}
	boardStructMap := yaml.MapSlice{}
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		err = fmt.Errorf("error at yaml unmarshall: %w", err)
		return err
	}
	for _, v := range m {
		if v.Key == "ApplicationConfig" {
			i, ok := v.Value.(yaml.MapSlice)
			if !ok {
				err = fmt.Errorf("configファイルが変です。多分ApplicationConfigの項目がありません")
				return err
			}
			for _, v := range i {
				if v.Key == "BoardStruct" {
					boardStructMap, ok = v.Value.(yaml.MapSlice)
					if !ok {
						err = fmt.Errorf("configファイルが変です。多分ApplicationConfigの項目、BoardStructがありません")
						return err
					}
				}
			}
		}
	}
	LoadedConfig.ApplicationConfig.BoardStruct = MapSlice(boardStructMap)
	return nil
}

// MapSlice . yaml.MapSliceをJSONにするために用意されたものです
type MapSlice yaml.MapSlice

// MapItem . yaml.MapItemをJSONにするために用意されたものです
type MapItem yaml.MapItem

// MarshalJSON . JSONにMarshalします。
func (m MapSlice) MarshalJSON() ([]byte, error) {
	jsonStr := "{"
	for i, item := range m {
		if i != 0 {
			jsonStr += ","
		}
		switch value := interface{}(item.Value).(type) {
		case yaml.MapSlice:
			itemJSON, err := json.Marshal(MapSlice(value))
			if err != nil {
				err = fmt.Errorf("error at marshal json: %w", err)
				return nil, err
			}
			jsonStr += fmt.Sprintf(`"%s": %s`, item.Key, string(itemJSON))
		case yaml.MapItem:
			ValueJSON, err := json.Marshal(MapItem(value))
			if err != nil {
				err = fmt.Errorf("error at marshal json: %w", err)
				return nil, err
			}
			jsonStr += fmt.Sprintf(`"%s": "%s"`, item.Key, string(ValueJSON))
		case string:
			jsonStr += fmt.Sprintf(`"%s": "%s"`, item.Key, value)
		default:
			err := fmt.Errorf("変な型が渡されました %s", reflect.TypeOf(item.Value))
			return nil, err
		}
	}
	jsonStr += "}"
	return []byte(jsonStr), nil
}

// MarshalJSON . JSONにMarshalします。
func (m MapItem) MarshalJSON() ([]byte, error) {
	jsonStr := "{"
	switch value := interface{}(m.Value).(type) {
	case yaml.MapSlice:
		itemJSON, err := json.Marshal(MapSlice(value))
		if err != nil {
			err = fmt.Errorf("error at marshal json: %w", err)
			return nil, err
		}
		jsonStr += fmt.Sprintf(`"%s": %s`, m.Key, string(itemJSON))
	case string:
		jsonStr += fmt.Sprintf(`"%s": "%s" `, m.Key, value)
	default:
		err := fmt.Errorf("変な型が渡されました %s", reflect.TypeOf(m.Value))
		return nil, err
	}
	jsonStr += "}"
	return []byte(jsonStr), nil
}

func loadConfig() error {
	configOpt := getConfigFile()
	config := getConfig()
	configName := getConfigName()
	configExt := getConfigExt()

	v := viper.New()
	configPaths := []string{}
	if configOpt != "" {
		// コンフィグファイルが明示的に指定された場合はそれを
		v.SetConfigFile(configOpt)
		configPaths = append(configPaths, configOpt)
	} else {
		// 実行ファイルの親ディレクトリ、カレントディレクトリ、ホームディレクトリの順に
		v.SetConfigName(configName)
		exe, err := os.Executable()
		if err != nil {
			err = fmt.Errorf("error at get executable file path: %w", err)
			log.Printf(err.Error())
		} else {
			v.AddConfigPath(filepath.Dir(exe))
			configPaths = append(configPaths, filepath.Join(filepath.Dir(exe), configName+configExt))
		}

		v.AddConfigPath(".")
		configPaths = append(configPaths, filepath.Join(".", configName+configExt))

		home, err := homedir.Dir()
		if err != nil {
			err = fmt.Errorf("error at get user home directory: %w", err)
			log.Printf(err.Error())
		} else {
			v.AddConfigPath(home)
			configPaths = append(configPaths, filepath.Join(home, configName+configExt))
		}
	}

	// 読み込んでcfgを作成する
	existConfigPath := false
	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); err == nil {
			existConfigPath = true
			break
		}
	}
	if !existConfigPath {
		// コンフィグファイルが指定されていなくてコンフィグファイルが見つからなかった場合、
		// ホームディレクトリにデフォルトコンフィグファイルを作成する。
		// できなければカレントディレクトリにコンフィグファイルを作成する。
		if configOpt == "" {
			configDir := ""
			home, err := homedir.Dir()
			if err != nil {
				err = fmt.Errorf("error at get user home directory: %w", err)
				log.Printf(err.Error())
				configDir = "."
			} else {
				configDir = home
			}

			configFileName := filepath.Join(configDir, configName+configExt)
			err = os.WriteFile(configFileName, []byte(CreateDefaultConfigYAML()), os.ModePerm)
			if err != nil {
				err = fmt.Errorf("error at write file to %s: %w", configFileName, err)
				return err
			}
			v.SetConfigFile(configFileName)
		} else {
			err := fmt.Errorf("コンフィグファイルが見つかりませんでした。")
			return err
		}
	}

	err := v.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("error at read in config: %w", err)
		return err
	}

	err = v.Unmarshal(config)
	if err != nil {
		err = fmt.Errorf("error at unmarshal config file: %w", err)
		return err
	}

	// 各DBファイルの作成
	if config.Repositories.MiRep == nil {
		err := fmt.Errorf("configファイルのRepositories.MiRepの項目が設定されていないかあるいは不正です")
		return err
	}
	if config.Repositories.TagRep == nil {
		err := fmt.Errorf("configファイルのRepositories.TagRepの項目が設定されていないかあるいは不正です")
		return err
	}
	if config.Repositories.TextRep == nil {
		err := fmt.Errorf("configファイルのRepositories.TextRepの項目が設定されていないかあるいは不正です")
		return err
	}
	files := []string{
		os.ExpandEnv(config.Repositories.MiRep.File),
		os.ExpandEnv(config.Repositories.TagRep.File),
		os.ExpandEnv(config.Repositories.TextRep.File),
	}

	for _, filename := range files {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("error at create file %s: %w", filename, err)
			return err
		}
		defer f.Close()
	}

	return nil
}

// configの情報をもとにrepositoriesを読み込む
func LoadRepositories() error {
	r := &Repositories{}

	if LoadedConfig.Repositories.MiRep == nil {
		err := fmt.Errorf("configファイルのRepositories.MiRepの項目が設定されていないかあるいは不正です")
		return err
	}
	reps, err := LoadMiReps(LoadedConfig.Repositories.MiRep)
	if err != nil {
		err = fmt.Errorf("error at load rep: %w", err)
		return err
	}
	r.MiRep = reps[0]

	addedMiReps := []mi.MiRep{}
	addedMiReps = append(addedMiReps, r.MiRep)
	if LoadedConfig.Repositories.MiReps == nil {
		err := fmt.Errorf("configファイルのRepositories.MiRepsの項目が設定されていないかあるいは不正です")
		return err
	}
	for _, repInfo := range LoadedConfig.Repositories.MiReps {
		reps, err := LoadMiReps(repInfo)
		if err != nil {
			err = fmt.Errorf("error at load reps: %w", err)
			return err
		}
		for _, miRep := range reps {
			for _, existMiRep := range addedMiReps {
				if filepath.Base(miRep.Path()) == filepath.Base(existMiRep.Path()) {
					continue
				}
			}
			r.MiReps = append(r.MiReps, miRep)
			addedMiReps = append(addedMiReps, miRep)
		}
	}
	r.MiReps = []mi.MiRep{
		mi.NewCachedMiRep(r.MiReps),
	}
	r.MiReps = append(r.MiReps, r.MiRep)

	if LoadedConfig.Repositories.TagReps == nil {
		err := fmt.Errorf("configファイルのRepositories.TagRepsの項目が設定されていないかあるいは不正です")
		return err
	}
	for _, tagRepInfo := range LoadedConfig.Repositories.TagReps {
		tagReps, err := tag.LoadTagReps(tagRepInfo)
		if err != nil {
			err = fmt.Errorf("error at load tag reps type=%s file=%s: %w", tagRepInfo.Type, tagRepInfo.File, err)
			return err
		}
		for _, tagRep := range tagReps {
			r.TagReps = append(r.TagReps, tag.NewCachedTagRep(tagRep, time.Hour*24))
		}
	}

	if LoadedConfig.Repositories.TextReps == nil {
		err := fmt.Errorf("configファイルのRepositories.TextRepsの項目が設定されていないかあるいは不正です")
		return err
	}
	for _, textRepInfo := range LoadedConfig.Repositories.TextReps {
		textReps, err := text.LoadTextReps(textRepInfo)
		if err != nil {
			err = fmt.Errorf("error at load text reps type=%s file=%s: %w", textRepInfo.Type, textRepInfo.File, err)
			return err
		}
		r.TextReps = append(r.TextReps, textReps...)
	}

	if LoadedConfig.Repositories.TagRep == nil {
		err := fmt.Errorf("configファイルのRepositories.TagRepの項目が設定されていないかあるいは不正です")
		return err
	}
	writetoTagRep, err := tag.LoadTagReps(LoadedConfig.Repositories.TagRep)
	if err != nil {
		err = fmt.Errorf("error at load write to tag rep: %w", err)
		return err
	}
	if len(writetoTagRep) != 1 {
		err = fmt.Errorf("見つかったtag repの数が1つではありませんでした。")
		return err
	}
	r.TagRep = tag.NewCachedTagRep(writetoTagRep[0], time.Hour*24)

	if LoadedConfig.Repositories.TextRep == nil {
		err := fmt.Errorf("configファイルのRepositories.TextRepの項目が設定されていないかあるいは不正です")
		return err
	}
	writetoTextRep, err := text.LoadTextReps(LoadedConfig.Repositories.TextRep)
	if err != nil {
		err = fmt.Errorf("error at to load write to text rep: %w", err)
		return err
	}
	if len(writetoTextRep) != 1 {
		err = fmt.Errorf("見つかったtext repの数が1つではありませんでした。")
		return err
	}
	r.TextRep = writetoTextRep[0]

	r.DeleteTagReps = tag.NewDeleteTagReps(r.TagRep, r.TagReps, time.Hour*24*365)

	LoadedRepositories = r
	return nil
}

func LaunchServer() error {
	router := registrep.Router

	/*
		//旧版DB対応
		allTasks, err := LoadedRepositories.MiReps.GetAllTasks(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, task := range allTasks {
			start, err := LoadedRepositories.MiReps.GetLatestStartInfoFromTaskID(context.TODO(), task.TaskID)
			if start == nil || err != nil {
				LoadedRepositories.MiRep.AddStartInfo(&mi.StartInfo{
					StartID:     uuid.New().String(),
					TaskID:      task.TaskID,
					UpdatedTime: time.Now(),
					Start:       nil,
				})
			}
			end, err := LoadedRepositories.MiReps.GetLatestEndInfoFromTaskID(context.TODO(), task.TaskID)
			if end == nil || err != nil {
				LoadedRepositories.MiRep.AddEndInfo(&mi.EndInfo{
					EndID:       uuid.New().String(),
					TaskID:      task.TaskID,
					UpdatedTime: time.Now(),
					End:         nil,
				})
			}
		}
	*/

	router.HandleFunc(get_board_struct_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetBoardStructRequest{}
		response := &GetBoardStructResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}
		response.BoardStruct = LoadedConfig.ApplicationConfig.BoardStruct
	}).Methods(get_board_struct_method)

	router.HandleFunc(get_tag_struct_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTagStructRequest{}
		response := &GetTagStructResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}
		response.TagStruct = LoadedConfig.ApplicationConfig.TagStruct
	}).Methods(get_tag_struct_method)

	router.HandleFunc(add_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &AddTaskRequest{}
		response := &AddTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		if request.TaskInfo.Task.TaskID != request.TaskInfo.TaskTitleInfo.TaskID ||
			request.TaskInfo.Task.TaskID != request.TaskInfo.CheckStateInfo.TaskID ||
			request.TaskInfo.Task.TaskID != request.TaskInfo.LimitInfo.TaskID ||
			request.TaskInfo.Task.TaskID != request.TaskInfo.StartInfo.TaskID ||
			request.TaskInfo.Task.TaskID != request.TaskInfo.EndInfo.TaskID ||
			request.TaskInfo.Task.TaskID != request.TaskInfo.BoardInfo.TaskID {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "TaskIDが一致しません")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = LoadedRepositories.MiRep.AddTask(request.TaskInfo.Task)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "タスクの追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddTaskTitleInfo(request.TaskInfo.TaskTitleInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "タイトル情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddCheckStateInfo(request.TaskInfo.CheckStateInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "チェック情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddLimitInfo(request.TaskInfo.LimitInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "期限情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddStartInfo(request.TaskInfo.StartInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "開始情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddEndInfo(request.TaskInfo.EndInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "終了情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.MiRep.AddBoardInfo(request.TaskInfo.BoardInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の追加に失敗しました")
			response.Errors = append(response.Errors, "板情報の追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(add_task_method)

	router.HandleFunc(update_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &UpdateTaskRequest{}
		response := &UpdateTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		currentTaskTitleInfo, err := LoadedRepositories.MiReps.GetLatestTaskTitleInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクのタイトル情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currentCheckStateInfo, err := LoadedRepositories.MiReps.GetLatestCheckStateInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクのチェック情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currentLimitInfo, err := LoadedRepositories.MiReps.GetLatestLimitInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクの期限情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currentStartInfo, err := LoadedRepositories.MiReps.GetLatestStartInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクの開始情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currentEndInfo, err := LoadedRepositories.MiReps.GetLatestEndInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクの終了情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currentBoardInfo, err := LoadedRepositories.MiReps.GetLatestBoardInfoFromTaskID(r.Context(), request.TaskInfo.Task.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスクの更新に失敗しました")
			response.Errors = append(response.Errors, "タスクの板情報取得時にエラーが発生しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if request.TaskInfo.TaskTitleInfo.Title != currentTaskTitleInfo.Title {
			err := LoadedRepositories.MiRep.AddTaskTitleInfo(request.TaskInfo.TaskTitleInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクのタイトル情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if request.TaskInfo.CheckStateInfo.IsChecked != currentCheckStateInfo.IsChecked {
			err := LoadedRepositories.MiRep.AddCheckStateInfo(request.TaskInfo.CheckStateInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクのチェック情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if (request.TaskInfo.LimitInfo.Limit != nil && currentLimitInfo.Limit == nil) ||
			(request.TaskInfo.LimitInfo.Limit == nil && currentLimitInfo.Limit != nil) ||
			(request.TaskInfo.LimitInfo.Limit != nil && currentLimitInfo.Limit != nil && !request.TaskInfo.LimitInfo.Limit.Equal(*currentLimitInfo.Limit)) {
			err := LoadedRepositories.MiRep.AddLimitInfo(request.TaskInfo.LimitInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクの期限情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if (request.TaskInfo.StartInfo.Start != nil && currentStartInfo.Start == nil) ||
			(request.TaskInfo.StartInfo.Start == nil && currentStartInfo.Start != nil) ||
			(request.TaskInfo.StartInfo.Start != nil && currentStartInfo.Start != nil && !request.TaskInfo.StartInfo.Start.Equal(*currentStartInfo.Start)) {
			err := LoadedRepositories.MiRep.AddStartInfo(request.TaskInfo.StartInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクの開始情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if (request.TaskInfo.EndInfo.End != nil && currentEndInfo.End == nil) ||
			(request.TaskInfo.EndInfo.End == nil && currentEndInfo.End != nil) ||
			(request.TaskInfo.EndInfo.End != nil && currentEndInfo.End != nil && !request.TaskInfo.EndInfo.End.Equal(*currentEndInfo.End)) {
			err := LoadedRepositories.MiRep.AddEndInfo(request.TaskInfo.EndInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクの終了情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if request.TaskInfo.BoardInfo.BoardName != currentBoardInfo.BoardName {
			err := LoadedRepositories.MiRep.AddBoardInfo(request.TaskInfo.BoardInfo)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの更新に失敗しました")
				response.Errors = append(response.Errors, "タスクの板情報更新時にエラーが発生しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}).Methods(update_task_method)

	router.HandleFunc(delete_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &DeleteTaskRequest{}
		response := &DeleteTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		deleted := false
		for _, taskRep := range LoadedRepositories.MiReps {
			err = taskRep.Delete(request.TaskID)
			if err != nil {
				response.Errors = append(response.Errors, "タスクの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deleted = true
			break
		}
		if !deleted {
			if err != nil {
				response.Errors = append(response.Errors, "タスクの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		LoadedRepositories.DeleteTagReps.UpdateCache(r.Context())

		err = LoadedRepositories.TagRep.UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = tag.TagReps(LoadedRepositories.TagReps).UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(delete_task_method)

	router.HandleFunc(get_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTaskRequest{}
		response := &GetTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}
		taskInfo, err := LoadedRepositories.MiReps.GetTaskInfo(r.Context(), request.TaskID)
		if err != nil {
			response.Errors = append(response.Errors, "タスク情報の取得に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.TaskInfo = taskInfo
	}).Methods(get_task_method)

	router.HandleFunc(get_tasks_from_board_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTaskFromBoardRequest{}
		response := &GetTaskFromBoardResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		if request.UpdateCache {
			err := LoadedRepositories.UpdateCache(r.Context())
			if err != nil {
				response.Errors = append(response.Errors, "板内タスク情報の取得に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		boardsTasksMap := map[string]*mi.Task{}
		if request.Query.Word != "" {
			words, notWords := parseWords(request.Query.Word)

			boardsTasksMap, err = filterWords(r.Context(), LoadedRepositories.MiReps, LoadedRepositories.TextReps, words, notWords, false, request.Query)
			if err != nil {
				response.Errors = append(response.Errors, "板内タスク情報の取得に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			boatdsTasks, err := LoadedRepositories.MiReps.GetTasksAtBoard(r.Context(), request.Query)
			if err != nil {
				response.Errors = append(response.Errors, "板内タスク情報の取得に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, task := range boatdsTasks {
				boardsTasksMap[task.TaskID] = task
			}
		}

		boardsTasksMap, err = filterTags(r.Context(), boardsTasksMap, LoadedRepositories.TagReps, request.Query.Tags, Or)
		if err != nil {
			response.Errors = append(response.Errors, "板内タスク情報の取得に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		boardsTasks := []*mi.Task{}
		for _, task := range boardsTasksMap {
			boardsTasks = append(boardsTasks, task)
		}

		boardsTaskInfos := []*mi.TaskInfo{}
		for _, task := range boardsTasks {
			taskInfo, err := LoadedRepositories.MiReps.GetTaskInfo(r.Context(), task.TaskID)
			if err != nil {
				response.Errors = append(response.Errors, "タスク情報の取得に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			boardsTaskInfos = append(boardsTaskInfos, taskInfo)
		}

		switch request.Query.SortType {
		case mi.CreatedTimeDesc:
			sort.Slice(boardsTaskInfos, func(i int, j int) bool {
				return boardsTaskInfos[i].Task.CreatedTime.After(boardsTaskInfos[j].Task.CreatedTime)
			})
		case mi.LimitTimeAsc:
			hasLimitTaskInfos := []*mi.TaskInfo{}
			noLimitTaskInfos := []*mi.TaskInfo{}

			for _, taskInfo := range boardsTaskInfos {
				if taskInfo.LimitInfo.Limit == nil {
					noLimitTaskInfos = append(noLimitTaskInfos, taskInfo)
				} else {
					hasLimitTaskInfos = append(hasLimitTaskInfos, taskInfo)
				}
			}

			sort.Slice(hasLimitTaskInfos, func(i int, j int) bool {
				limitI := *hasLimitTaskInfos[i].LimitInfo.Limit
				limitJ := *hasLimitTaskInfos[j].LimitInfo.Limit
				return limitI.Before(limitJ)
			})

			sort.Slice(noLimitTaskInfos, func(i int, j int) bool {
				return noLimitTaskInfos[i].Task.CreatedTime.After(noLimitTaskInfos[j].Task.CreatedTime)
			})

			boardsTaskInfos = append(hasLimitTaskInfos, noLimitTaskInfos...)
		case mi.StartTimeDesc:
			hasStartTaskInfos := []*mi.TaskInfo{}
			noStartTaskInfos := []*mi.TaskInfo{}

			for _, taskInfo := range boardsTaskInfos {
				if taskInfo.StartInfo.Start == nil {
					noStartTaskInfos = append(noStartTaskInfos, taskInfo)
				} else {
					hasStartTaskInfos = append(hasStartTaskInfos, taskInfo)
				}
			}

			sort.Slice(hasStartTaskInfos, func(i int, j int) bool {
				startI := *hasStartTaskInfos[i].StartInfo.Start
				startJ := *hasStartTaskInfos[j].StartInfo.Start
				return startI.Before(startJ)
			})

			sort.Slice(noStartTaskInfos, func(i int, j int) bool {
				return noStartTaskInfos[i].Task.CreatedTime.After(noStartTaskInfos[j].Task.CreatedTime)
			})

			boardsTaskInfos = append(hasStartTaskInfos, noStartTaskInfos...)
		case mi.EndTimeDesc:
			hasEndTaskInfos := []*mi.TaskInfo{}
			noEndTaskInfos := []*mi.TaskInfo{}

			for _, taskInfo := range boardsTaskInfos {
				if taskInfo.EndInfo.End == nil {
					noEndTaskInfos = append(noEndTaskInfos, taskInfo)
				} else {
					hasEndTaskInfos = append(hasEndTaskInfos, taskInfo)
				}
			}

			sort.Slice(hasEndTaskInfos, func(i int, j int) bool {
				endI := *hasEndTaskInfos[i].EndInfo.End
				endJ := *hasEndTaskInfos[j].EndInfo.End
				return endI.Before(endJ)
			})

			sort.Slice(noEndTaskInfos, func(i int, j int) bool {
				return noEndTaskInfos[i].Task.CreatedTime.After(noEndTaskInfos[j].Task.CreatedTime)
			})

			boardsTaskInfos = append(hasEndTaskInfos, noEndTaskInfos...)
		}

		response.BoardsTasks = boardsTaskInfos
	}).Methods(get_tasks_from_board_method)

	router.HandleFunc(add_tag_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &AddTagRequest{}
		response := &AddTagResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		tagInfo := &tag.Tag{
			ID:     uuid.New().String(),
			Target: request.TaskID,
			Tag:    request.Tag,
			Time:   time.Now(),
		}

		err = LoadedRepositories.TagRep.AddTag(tagInfo)
		if err != nil {
			response.Errors = append(response.Errors, "タグの追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = LoadedRepositories.TagRep.UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tag.TagReps(LoadedRepositories.TagReps).UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(add_tag_method)

	router.HandleFunc(add_text_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &AddTextRequest{}
		response := &AddTextResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		textInfo := &text.Text{
			ID:     uuid.New().String(),
			Target: request.TaskID,
			Text:   request.Text,
			Time:   time.Now(),
		}

		err = LoadedRepositories.TextRep.AddText(textInfo)
		if err != nil {
			response.Errors = append(response.Errors, "テキストの追加に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(add_text_method)

	router.HandleFunc(get_tags_related_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTagsRelatedTaskRequest{}
		response := &GetTagsRelatedTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		tags := map[string]*tag.Tag{}
		wg := &sync.WaitGroup{}
		ch := make(chan []*tag.Tag, len(LoadedRepositories.TagReps))
		for _, tagRep := range LoadedRepositories.TagReps {
			tagRep := tagRep
			wg.Add(1)
			go func(tagRep tag.TagRep) {
				defer wg.Done()
				matchTags, err := tagRep.GetTagsByTarget(r.Context(), request.TaskID)
				if err != nil {
					response.Errors = append(response.Errors, "タグの取得に失敗しました")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				ch <- matchTags
			}(tagRep)
		}
		wg.Wait()
	loop:
		for {
			select {
			case t := <-ch:
				if t == nil {
					continue loop
				}
				for _, tag := range t {
					tags[tag.ID] = tag
				}
			default:
				break loop
			}
		}
		tagList := []*tag.Tag{}
		for _, matchTag := range tags {
			tagList = append(tagList, matchTag)
		}
		response.Tags = tagList
	}).Methods(get_tags_related_task_method)

	router.HandleFunc(get_texts_related_task_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTextsRelatedTaskRequest{}
		response := &GetTextsRelatedTaskResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		texts := map[string]*text.Text{}
		wg := &sync.WaitGroup{}
		ch := make(chan []*text.Text, len(LoadedRepositories.TextReps))
		for _, textRep := range LoadedRepositories.TextReps {
			textRep := textRep
			wg.Add(1)
			go func(textRep text.TextRep) {
				defer wg.Done()
				matchTexts, err := textRep.GetTextsByTarget(r.Context(), request.TaskID)
				if err != nil {
					response.Errors = append(response.Errors, "テキストの取得に失敗しました")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				ch <- matchTexts
			}(textRep)
		}
		wg.Wait()
	loop:
		for {
			select {
			case t := <-ch:
				if t == nil {
					continue loop
				}
				for _, text := range t {
					texts[text.ID] = text
				}
			default:
				break loop
			}
		}
		textList := []*text.Text{}
		for _, matchText := range texts {
			textList = append(textList, matchText)
		}
		response.Texts = textList

	}).Methods(get_texts_related_task_method)

	router.HandleFunc(delete_tag_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &DeleteTagRequest{}
		response := &DeleteTagResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		deleted := false
		for _, tagRep := range LoadedRepositories.TagReps {
			err = tagRep.Delete(request.TagID)
			if err != nil {
				response.Errors = append(response.Errors, "タグの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deleted = true
			break
		}
		if !deleted {
			if err != nil {
				response.Errors = append(response.Errors, "タグの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		LoadedRepositories.DeleteTagReps.UpdateCache(r.Context())

		err = LoadedRepositories.TagRep.UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = tag.TagReps(LoadedRepositories.TagReps).UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(delete_tag_method)

	router.HandleFunc(delete_text_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &DeleteTextRequest{}
		response := &DeleteTextResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		deleted := false
		for _, textRep := range LoadedRepositories.TextReps {
			err = textRep.Delete(request.TextID)
			if err != nil {
				response.Errors = append(response.Errors, "テキストの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deleted = true
			break
		}
		if !deleted {
			if err != nil {
				response.Errors = append(response.Errors, "テキストの削除に失敗しました")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		LoadedRepositories.DeleteTagReps.UpdateCache(r.Context())

		err = LoadedRepositories.TagRep.UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = tag.TagReps(LoadedRepositories.TagReps).UpdateCache(r.Context())
		if err != nil {
			response.Errors = append(response.Errors, "タグのキャッシュ更新に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(delete_text_method)

	router.HandleFunc(get_tag_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTagRequest{}
		response := &GetTagResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		matchTag := &tag.Tag{}
		wg := &sync.WaitGroup{}
		ch := make(chan *tag.Tag, len(LoadedRepositories.TagReps))
		defer close(ch)
		for _, tagRep := range LoadedRepositories.TagReps {
			tagRep := tagRep
			wg.Add(1)
			go func(tagRep tag.TagRep) {
				defer wg.Done()
				tag, err := tagRep.GetTagByID(r.Context(), request.TagID)
				if err != nil {
					return
				}
				ch <- tag
			}(tagRep)
		}
		wg.Wait()
	loop:
		for {
			select {
			case t := <-ch:
				if t == nil {
					continue loop
				}
				matchTag = t
			default:
				break loop
			}
		}
		response.Tag = matchTag
	}).Methods(get_tag_method)

	router.HandleFunc(get_text_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTextRequest{}
		response := &GetTextResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		matchText := &text.Text{}
		wg := &sync.WaitGroup{}
		ch := make(chan *text.Text, len(LoadedRepositories.TextReps))
		defer close(ch)
		for _, textRep := range LoadedRepositories.TextReps {
			textRep := textRep
			wg.Add(1)
			go func(textRep text.TextRep) {
				defer wg.Done()
				text, err := textRep.GetTextByID(r.Context(), request.TextID)
				if err != nil {
					return
				}
				ch <- text
			}(textRep)
		}
		wg.Wait()
	loop:
		for {
			select {
			case t := <-ch:
				matchText = t
			default:
				break loop
			}
		}

		response.Text = matchText
	}).Methods(get_text_method)

	router.HandleFunc(get_tag_names_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetTagNamesRequest{}
		response := &GetTagNamesResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		tags := []string{}
		tagNames := map[string]interface{}{}
		wg := &sync.WaitGroup{}
		ch := make(chan []*tag.Tag, len(LoadedRepositories.TagReps))
		for _, tagRep := range LoadedRepositories.TagReps {
			tagRep := tagRep
			wg.Add(1)
			go func(tagRep tag.TagRep) {
				defer wg.Done()
				loadedTags, err := tagRep.GetAllTags(r.Context())
				if err != nil {
					response.Errors = append(response.Errors, "タグ一覧の取得に失敗しました")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				ch <- loadedTags
			}(tagRep)
		}
		wg.Wait()
	loop:
		for {
			select {
			case t := <-ch:
				if t == nil {
					continue loop
				}
				for _, tag := range t {
					tagNames[tag.Tag] = tag
				}
			default:
				break loop
			}
		}

		for tagName := range tagNames {
			if tagName != kyou.DeletedTagName {
				tags = append(tags, strings.TrimSpace(tagName))
			}
		}
		sort.Slice(tags, func(i, j int) bool {
			return tags[i] < tags[j]
		})
		response.TagNames = tags
	}).Methods(get_tag_names_method)

	router.HandleFunc(get_board_names_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetBoardNamesRequest{}
		response := &GetBoardNamesResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}
		ch := make(chan []string, len(LoadedRepositories.MiReps))
		for _, miRep := range LoadedRepositories.MiReps {
			miRep := miRep
			wg.Add(1)
			go func(miRep mi.MiRep) {
				defer wg.Done()
				repsBoardNames := []string{}
				tasks, err := miRep.GetAllTasks(r.Context())
				if err != nil {
					response.Errors = append(response.Errors, "板一覧の取得に失敗しました")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				for _, task := range tasks {
					boardInfo, err := miRep.GetLatestBoardInfoFromTaskID(r.Context(), task.TaskID)
					if err != nil {
						// response.Errors = append(response.Errors, "板一覧の取得に失敗しました")
						// w.WriteHeader(http.StatusInternalServerError)
						// errch <- err
						continue
					}
					if boardInfo != nil {
						repsBoardNames = append(repsBoardNames, boardInfo.BoardName)
					}
				}
				ch <- repsBoardNames
			}(miRep)
		}
		boardNamesList := []string{}
		boardStructMap := getConfig().ApplicationConfig.BoardStruct.(MapSlice)
		for _, boardNameObj := range boardStructMap {
			if boardNameObj.Value.(string) == "board" {
				boardName := boardNameObj.Key.(string)
				if boardName == mi.AllBoardName {
					continue
				}
				boardNamesList = append(boardNamesList, boardName)
			}
		}
		wg.Wait()
	loop:
		for {
			select {
			case repsBoardNames := <-ch:
				if repsBoardNames == nil {
					continue loop
				}
				for _, repsBoardName := range repsBoardNames {
					boardNamesList = append(boardNamesList, repsBoardName)
				}
			default:
				break loop
			}
		}

		boardNames := []string{}
		for _, boardName := range boardNamesList {
			existBoardName := false
			for _, boardNameInList := range boardNames {
				if boardName == boardNameInList {
					existBoardName = true
					break
				}
			}
			if existBoardName {
				continue
			}
			boardNames = append(boardNames, boardName)
		}
		response.BoardNames = boardNames
	}).Methods(get_board_names_method)

	router.HandleFunc(get_application_config_address, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		request := &GetApplicationConfigRequest{}
		response := &GetApplicationConfigResponse{}

		defer r.Body.Close()
		defer func() {
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				panic(err)
			}
		}()

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			panic(err)
		}

		response.ApplicationConfig = LoadedConfig.ApplicationConfig
	}).Methods(get_application_config_method)

	html, err := fs.Sub(htmlFS, "mi/mi/embed/html")
	if err != nil {
		return err
	}
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(html)).ServeHTTP(w, r)
	})

	var handler http.Handler = router
	if LoadedConfig.ServerConfig.LocalOnly {
		h := handler
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spl := strings.Split(r.RemoteAddr, ":")
			remoteHost := strings.Join(spl[:len(spl)-1], ":")
			switch remoteHost {
			case "localhost":
				fallthrough
			case "127.0.0.1":
				fallthrough
			case "[::1]":
				fallthrough
			case "::1":
				h.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusForbidden)
		})
	}

	if LoadedConfig.ServerConfig.TLS.Enable {
		err = http.ListenAndServeTLS(
			LoadedConfig.ServerConfig.Address,
			os.ExpandEnv(LoadedConfig.ServerConfig.TLS.CertFile),
			os.ExpandEnv(LoadedConfig.ServerConfig.TLS.KeyFile),
			handler)
		if err != nil {
			err = fmt.Errorf("failed to launch server: %w", err)
			return err
		}
	} else {
		err = http.ListenAndServe(LoadedConfig.ServerConfig.Address, handler)
		if err != nil {
			err = fmt.Errorf("failed to launch server: %w", err)
			return err
		}
	}
	return nil
}

func sortTagsByTime(tags []*tag.Tag) {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Time.Before(tags[j].Time)
	})
}

func sortTextsByTime(texts []*text.Text) {
	sort.Slice(texts, func(i, j int) bool {
		return texts[i].Time.Before(texts[j].Time)
	})
}

func parseWords(word string) (words, notWords []string) {
	nextIsNotWord := false
	for _, word := range strings.Split(word, " ") {
		for _, word := range strings.Split(word, "　") {
			if strings.HasPrefix(word, "-") {
				nextIsNotWord = true
				word = strings.TrimPrefix(word, "-")
			}
			switch word {
			case "":
				continue
			case "-":
				nextIsNotWord = true
			default:
				if nextIsNotWord {
					notWords = append(notWords, word)
				} else {
					words = append(words, word)
				}
				nextIsNotWord = false
			}
		}
	}
	return words, notWords
}

type Repositories struct {
	MiReps   mi.MiReps
	TagReps  []tag.TagRep
	TextReps []text.TextRep

	MiRep   mi.MiRep
	TagRep  tag.TagRep
	TextRep text.TextRep

	DeleteTagReps tag.DeleteTagReps
}

func (r *Repositories) UpdateCache(ctx context.Context) error {
	var err error
	err = r.MiRep.UpdateCache(ctx)
	if err != nil {
		err = fmt.Errorf("error at update cache at %s: %w", r.MiRep.RepName(), err)
		return err
	}

	err = r.MiReps.UpdateCache(ctx)
	if err != nil {
		err = fmt.Errorf("error at update cache at mi reps: %w", err)
		return err
	}

	for _, rep := range r.TagReps {
		err := rep.UpdateCache(ctx)
		if err != nil {
			err = fmt.Errorf("error at update cache at %s: %w", rep.Path(), err)
			return err
		}
	}
	for _, rep := range r.TextReps {
		err := rep.UpdateCache(ctx)
		if err != nil {
			err = fmt.Errorf("error at update cache at %s: %w", rep.Path(), err)
			return err
		}
	}

	err = r.TagRep.UpdateCache(ctx)
	if err != nil {
		err = fmt.Errorf("error at update cache at %s: %w", r.TagRep.Path(), err)
		return err
	}
	err = r.TextRep.UpdateCache(ctx)
	if err != nil {
		err = fmt.Errorf("error at update cache at %s: %w", r.TextRep.Path(), err)
		return err
	}

	err = r.DeleteTagReps.UpdateCache(ctx)
	if err != nil {
		err = fmt.Errorf("error at update cache at delete tag reps: %w", err)
		return err
	}

	return nil
}

// Close . Repositories内のすべてのrepを閉じる
func (r *Repositories) Close() error {
	var err error
	for _, rep := range r.MiReps {
		err := rep.Close()
		if err != nil {
			err = fmt.Errorf("error at close %s: %w", rep.Path(), err)
			return err
		}
	}
	for _, rep := range r.TagReps {
		err := rep.Close()
		if err != nil {
			err = fmt.Errorf("error at close %s: %w", rep.Path(), err)
			return err
		}
	}
	for _, rep := range r.TextReps {
		err := rep.Close()
		if err != nil {
			err = fmt.Errorf("error at close %s: %w", rep.Path(), err)
			return err
		}
	}
	err = r.TagRep.Close()
	if err != nil {
		err = fmt.Errorf("error at close %s: %w", r.TagRep.Path(), err)
		return err
	}
	err = r.TextRep.Close()
	if err != nil {
		err = fmt.Errorf("error at close %s: %w", r.TextRep.Path(), err)
		return err
	}
	return nil
}

func WrapT(repos *Repositories) (*Repositories, error) {
	repos.MiReps = wrapMiRepsT(repos.MiReps, repos.DeleteTagReps)
	repos.TagReps = wrapTagRepsT(repos.TagReps, repos.DeleteTagReps)
	repos.TextReps = wrapTextRepsT(repos.TextReps, repos.DeleteTagReps)
	return repos, nil
}

func wrapMiRepsT(reps []mi.MiRep, deleteTagReps tag.DeleteTagReps) []mi.MiRep {
	wrapedReps := []mi.MiRep{}
	for _, rep := range reps {
		wrapedReps = append(wrapedReps, mi.WrapMiRepT(rep, deleteTagReps))
	}
	return wrapedReps
}

func wrapTagRepsT(reps []tag.TagRep, deleteTagReps tag.DeleteTagReps) []tag.TagRep {
	wrapedReps := []tag.TagRep{}
	for _, rep := range reps {
		wrapedReps = append(wrapedReps, tag.WrapTagRepT(rep, deleteTagReps))
	}
	return wrapedReps
}

func wrapTextRepsT(reps []text.TextRep, deleteTagReps tag.DeleteTagReps) []text.TextRep {
	wrapedReps := []text.TextRep{}
	for _, rep := range reps {
		wrapedReps = append(wrapedReps, text.WrapTextRepT(rep, deleteTagReps))
	}
	return wrapedReps
}

// NoTag . tagが一つもついていないkyouに自動的につけられるタグ名
const NoTag = `no tag`

// TagFilterMode .
// タグの検索モード。And, Or, Onlyのいずれか
type TagFilterMode string

// TagFilterModeの一覧
const (
	And  TagFilterMode = "and"
	Or   TagFilterMode = "or"
	Only TagFilterMode = "only"
)

func filterTags(ctx context.Context, matchTasks map[string]*mi.Task, tagReps []tag.TagRep, tags []string, mode TagFilterMode) (map[string]*mi.Task, error) {
	// タグを持っていないidを取得する
	noHaveTagTasks := map[string]*mi.Task{}
	haveTagTasks := map[string]struct{}{}

	existErr := false
	var e error

	wg := &sync.WaitGroup{}
	tagsch := make(chan []*tag.Tag, len(tagReps))
	errch := make(chan error, len(tagReps))
	defer close(tagsch)
	defer close(errch)
	for _, tagrep := range tagReps {
		wg.Add(1)
		tagrep := tagrep
		go func(tagrep tag.TagRep) {
			defer wg.Done()
			tags, err := tagrep.GetAllTags(ctx)
			if err != nil {
				err = fmt.Errorf("error at get all tags from tagrep %s: %w", tagrep.Path(), err)
				errch <- err
			}
			tagsch <- tags
		}(tagrep)
	}
	wg.Wait()

errloop1:
	for {
		select {
		case e := <-errch:
			e = fmt.Errorf("error at filter tags: %w", e)
			existErr = true
		default:
			break errloop1
		}
	}
	if existErr {
		return nil, e
	}

loop1:
	for {
		select {
		case t := <-tagsch:
			if t == nil {
				continue loop1
			}
			for _, tag := range t {
				haveTagTasks[tag.Target] = struct{}{}
			}
		default:
			break loop1
		}
	}

	for _, task := range matchTasks {
		if _, exist := haveTagTasks[task.TaskID]; !exist {
			noHaveTagTasks[task.TaskID] = task
		}
	}

	if mode == Or {
		wg := &sync.WaitGroup{}
		tasksch := make(chan map[string]*mi.Task, len(tagReps))
		errch := make(chan error, len(tagReps))
		defer close(tasksch)
		defer close(errch)

		// tagがあり、or検索の場合は、タグにヒットしたやつすべて
		temp := map[string]*mi.Task{}
		for _, tagrep := range tagReps {
			wg.Add(1)
			tagrep := tagrep
			go func(tagrep tag.TagRep) {
				defer wg.Done()
				t := map[string]*mi.Task{}
				for _, tagname := range tags {
					tags, err := tagrep.GetTagsByName(ctx, tagname)
					if err != nil {
						err = fmt.Errorf("error at get tag by name %s from tagrep %s: %w", tagname, tagrep.Path(), err)
						errch <- err
					}
					for _, tag := range tags {
						if task, exist := matchTasks[tag.Target]; exist {
							t[task.TaskID] = task
						}
					}
				}
				tasksch <- t
			}(tagrep)
		}
		wg.Wait()
	errloop2:
		for {
			select {
			case e := <-errch:
				e = fmt.Errorf("error at filter tags: %w", e)
				existErr = true
			default:
				break errloop2
			}
		}
		if existErr {
			return nil, e
		}

	loop2:
		for {
			select {
			case tasks := <-tasksch:
				if tasks == nil {
					continue loop2
				}
				for _, task := range tasks {
					temp[task.TaskID] = task
				}
			default:
				break loop2
			}
		}

		// notagが含まれたらタグを持っていないkyouを追加する
		for _, tag := range tags {
			if tag == NoTag {
				for _, task := range noHaveTagTasks {
					temp[task.TaskID] = task
				}
			}
		}
		matchTasks = map[string]*mi.Task{}
		for _, task := range temp {
			_, exist := matchTasks[task.TaskID]
			if !exist {
				matchTasks[task.TaskID] = task
			}
		}
		return filterHiddenTags(ctx, matchTasks, tagReps, tags)
	}

	temp := []*mi.Task{}
	for _, tag := range tags {
		if tag == NoTag {
			for _, task := range noHaveTagTasks {
				temp = append(temp, task)
			}
		}
	}
	for i, tagname := range tags {
		switch i {
		case 0:
			wg := &sync.WaitGroup{}
			tasksch := make(chan map[string]*mi.Task, len(tagReps))
			errch := make(chan error, len(tagReps))
			defer close(tasksch)
			defer close(errch)

			for _, tagrep := range tagReps {
				wg.Add(1)
				tagrep := tagrep
				go func(tagrep tag.TagRep) {
					defer wg.Done()
					t := map[string]*mi.Task{}
					tags, err := tagrep.GetTagsByName(ctx, tagname)
					if err != nil {
						err = fmt.Errorf("error at get tags by name %s from tagrep %s: %w", tagname, tagrep.Path(), err)
						errch <- err
					}
					for _, tag := range tags {
						if task, exist := matchTasks[tag.Target]; exist {
							t[task.TaskID] = task
						}
					}
					tasksch <- t
				}(tagrep)
			}
			wg.Wait()
		errloop3:
			for {
				select {
				case e := <-errch:
					e = fmt.Errorf("error at filter tags: %w", e)
					existErr = true
				default:
					break errloop3
				}
			}
			if existErr {
				return nil, e
			}

		loop3:
			for {
				select {
				case tasks := <-tasksch:
					if tasks == nil {
						continue loop3
					}
					for _, task := range tasks {
						temp = append(temp, task)
					}
				default:
					break loop3
				}
			}
		default:
			wg := &sync.WaitGroup{}
			tasksch := make(chan []*mi.Task, len(tagReps))
			errch := make(chan error, len(tagReps))
			defer close(tasksch)
			defer close(errch)
			temppp := []*mi.Task{}
			for _, tagrep := range tagReps {
				wg.Add(1)
				tagrep := tagrep
				go func(tagrep tag.TagRep) {
					t := []*mi.Task{}
					defer wg.Done()
					tags, err := tagrep.GetTagsByName(ctx, tagname)
					if err != nil {
						err = fmt.Errorf("failed to get tag by name %s from tagrep %s: %w", tagname, tagrep.Path(), err)
						errch <- err
					}

					tasks := []*mi.Task{}
					for _, tag := range tags {
						if id, exist := matchTasks[tag.Target]; exist {
							tasks = append(tasks, id)
						}
					}

					for _, existTask := range temp {
						exist := false
						for _, task := range tasks {
							if existTask.TaskID == task.TaskID {
								exist = true
							}
						}
						if exist {
							t = append(t, existTask)
						}
					}
					tasksch <- t
				}(tagrep)
			}
			wg.Wait()
		errloop4:
			for {
				select {
				case e := <-errch:
					e = fmt.Errorf("error at filter tags: %w", e)
					existErr = true
				default:
					break errloop4
				}
			}
			if existErr {
				return nil, e
			}

		loop4:
			for {
				select {
				case tasks := <-tasksch:
					if tasks == nil {
						continue loop4
					}
					for _, task := range tasks {
						temppp = append(temp, task)
					}
				default:
					break loop4
				}
			}

			temp = temppp
		}
	}
	matchTasks = map[string]*mi.Task{}
	for _, task := range temp {
		_, exist := matchTasks[task.TaskID]
		if !exist {
			matchTasks[task.TaskID] = task
		}
	}

	// OnlyModeでNoTagが含まれたらAnd検索結果と同義なので
	if mode == And || (mode == Only && equal([]string{NoTag}, tags)) {
		return filterHiddenTags(ctx, matchTasks, tagReps, tags)
	} else if mode == Only {
		allTags := []*tag.Tag{}

		wg := &sync.WaitGroup{}
		tagsch := make(chan []*tag.Tag, len(tagReps))
		errch := make(chan error, len(tagReps))
		defer close(tagsch)
		defer close(errch)
		for _, tagrep := range tagReps {
			wg.Add(1)
			tagrep := tagrep
			go func(tagrep tag.TagRep) {
				defer wg.Done()
				tags, err := tagrep.GetAllTags(ctx)
				if err != nil {
					err = fmt.Errorf("error at get all tags from %s: %w", tagrep.Path(), err)
					errch <- err
				}
				tagsch <- tags
			}(tagrep)
		}
		wg.Wait()
	errloop5:
		for {
			select {
			case e := <-errch:
				e = fmt.Errorf("error at filter tags: %w", e)
				existErr = true
			default:
				break errloop5
			}
		}
		if existErr {
			return nil, e
		}

	loop5:
		for {
			select {
			case t := <-tagsch:
				if t == nil {
					continue loop5
				}
				allTags = append(allTags, t...)
			default:
				break loop5
			}
		}

		// requestされたtagじゃないものがあったら除去する
		sortedTags := sort.StringSlice(tags)
		unMatchTaskTasks := map[string]struct{}{}
		for target := range matchTasks {
			attachedTagsMap := map[string]struct{}{}
			for _, tag := range allTags {
				if tag.Target == target {
					attachedTagsMap[tag.Tag] = struct{}{}
				}
			}
			attachedTags := []string{}
			for attachedTag := range attachedTagsMap {
				attachedTags = append(attachedTags, attachedTag)
			}
			sort.Strings(attachedTags)
			if !equal(sortedTags, attachedTags) {
				unMatchTaskTasks[target] = struct{}{}
			}
		}
		for unMatchTaskID := range unMatchTaskTasks {
			delete(matchTasks, unMatchTaskID)
		}
		return filterHiddenTags(ctx, matchTasks, tagReps, tags)
	}
	err := fmt.Errorf("invalid 'mode' value: %s", mode)
	return nil, err
}

func filterHiddenTags(ctx context.Context, matchTasks map[string]*mi.Task, tagReps []tag.TagRep, tags []string) (map[string]*mi.Task, error) {
	var e error
	existErr := false
loop1:
	for _, hiddenTag := range LoadedConfig.ApplicationConfig.HiddenTags {
		for _, tag := range tags {
			if hiddenTag == tag {
				continue loop1
			}
		}
		wg := &sync.WaitGroup{}
		tagsch := make(chan []*tag.Tag, len(tagReps))
		errch := make(chan error, len(tagReps))
		defer close(tagsch)
		defer close(errch)
		for _, tagrep := range tagReps {
			wg.Add(1)
			tagrep := tagrep
			go func(tagrep tag.TagRep) {
				defer wg.Done()
				tags, err := tagrep.GetTagsByName(ctx, hiddenTag)
				if err != nil {
					err = fmt.Errorf("error at get tags by name from %s: %w", tagrep.Path(), err)
					errch <- err
				}
				tagsch <- tags
			}(tagrep)
		}
		wg.Wait()
	errloop:
		for {
			select {
			case e := <-errch:
				e = fmt.Errorf("error at filter deleted tags: %w", e)
				existErr = true
			default:
				break errloop
			}
		}
		if existErr {
			return nil, e
		}

	loop2:
		for {
			select {
			case t := <-tagsch:
				if t == nil {
					continue loop2
				}
				for _, tag := range t {
					if _, exist := matchTasks[tag.Target]; exist {
						delete(matchTasks, tag.Target)
					}
				}
			default:
				break loop2
			}
		}
	}

	return matchTasks, nil
}

func filterWords(ctx context.Context, reps mi.MiReps, textReps []text.TextRep, words []string, notWords []string, and bool, query *mi.SearchTaskQuery) (map[string]*mi.Task, error) {
	existErr := false
	var e error
	matchTasks := map[string]*mi.Task{}
	// wordsがないときにはRep内のすべてのID
	if len(words) == 0 {
		allTasks := []*mi.Task{}

		wg := &sync.WaitGroup{}
		tasksch := make(chan []*mi.Task, len(reps))
		errch := make(chan error, len(reps))
		defer close(tasksch)
		defer close(errch)

		for _, rep := range reps {
			wg.Add(1)
			rep := rep
			go func(rep mi.MiRep) {
				defer wg.Done()
				tasks, err := rep.GetAllTasks(ctx)
				if err != nil {
					err = fmt.Errorf("error at get all tasks from %s: %w", rep.Path(), err)
					errch <- err
				}
				tasksch <- tasks
			}(rep)
		}
		wg.Wait()

	errloop:
		for {
			select {
			case e := <-errch:
				e = fmt.Errorf("error at filter tags: %w", e)
				existErr = true
			default:
				break errloop
			}
		}
		if existErr {
			return nil, e
		}

	loop:
		for {
			select {
			case tasks := <-tasksch:
				if tasks == nil {
					continue loop
				}
				allTasks = append(allTasks, tasks...)
			default:
				break loop
			}
		}

		// 重複がないようにMapに詰める
		for _, task := range allTasks {
			if _, exist := matchTasks[task.TaskID]; !exist {
				matchTasks[task.TaskID] = task
			}
		}

		// notWordsにhitしたものを外す
		if len(notWords) != 0 {
			notMatchTasks, err := orSearch(ctx, reps, textReps, notWords, query)
			if err != nil {
				err := fmt.Errorf("error at orSearch: %w", err)
				return nil, err
			}
			for _, notMatchTask := range notMatchTasks {
				if _, exist := matchTasks[notMatchTask.TaskID]; exist {
					delete(matchTasks, notMatchTask.TaskID)
				}
			}
		}
		return matchTasks, nil
	}
	// wordsの長さが1のときはor検索を使う（速いので）
	if len(words) == 1 {
		and = false
	}

	tasks := []*mi.Task{}
	var err error
	if and {
		tasks, err = andSearch(ctx, reps, textReps, words, query)
		if err != nil {
			err = fmt.Errorf("failed to and search: %w", err)
			return nil, err
		}
	} else {
		tasks, err = orSearch(ctx, reps, textReps, words, query)
		if err != nil {
			err = fmt.Errorf("failed to or search: %w", err)
			return nil, err
		}
	}

	// 重複がないようにMapに詰める
	for _, task := range tasks {
		if _, exist := matchTasks[task.TaskID]; !exist {
			matchTasks[task.TaskID] = task
		}
	}

	// notWordsにhitしたものを外す
	notTasks, err := orSearch(ctx, reps, textReps, notWords, query)
	if err != nil {
		err := fmt.Errorf("error at orSearch: %w", err)
		return nil, err
	}
	for _, notID := range notTasks {
		if _, exist := matchTasks[notID.TaskID]; exist {
			delete(matchTasks, notID.TaskID)
		}
	}
	return matchTasks, nil
}

func orSearch(ctx context.Context, reps mi.MiReps, textReps []text.TextRep, words []string, query *mi.SearchTaskQuery) ([]*mi.Task, error) {
	existErr := false
	var e error

	matchTasks := []*mi.Task{}
	allTasks := []*mi.Task{}

	wg := &sync.WaitGroup{}
	tasksch := make(chan []*mi.Task, len(reps))
	errch := make(chan error, len(reps))
	taskscht := make(chan []*mi.Task, len(textReps))
	errcht := make(chan error, len(reps))
	defer close(tasksch)
	defer close(errch)
	defer close(taskscht)
	defer close(errcht)
	for _, rep := range reps {
		wg.Add(1)
		rep := rep
		go func(rep mi.MiRep) {
			defer wg.Done()
			tasks, err := rep.GetAllTasks(ctx)
			if err != nil {
				err = fmt.Errorf("error at get all tasks from %s: %w", rep.Path(), err)
				errch <- err
			}
			tasksch <- tasks
		}(rep)
	}
	wg.Wait()

errloop:
	for {
		select {
		case e := <-errch:
			e = fmt.Errorf("error at filter tags: %w", e)
			existErr = true
		default:
			break errloop
		}
	}
	if existErr {
		return nil, e
	}

loop:
	for {
		select {
		case tasks := <-tasksch:
			if tasks == nil {
				continue loop
			}
			allTasks = append(allTasks, tasks...)
		default:
			break loop
		}
	}

	// repにSearchしてヒットしたもの
	for _, rep := range reps {
		wg.Add(1)
		rep := rep
		go func(rep mi.MiRep) {
			defer wg.Done()
			for _, word := range words {
				matchTasksInRep, err := rep.SearchTasks(ctx, word, query)
				if err != nil {
					err = fmt.Errorf("error at search %s in %s: %w", word, rep.Path(), err)
					errch <- err
				}
				tasksch <- matchTasksInRep
			}
		}(rep)

		wg.Wait()

	errloop2:
		for {
			select {
			case e := <-errch:
				e = fmt.Errorf("error at filter tags: %w", e)
				existErr = true
			default:
				break errloop2
			}
		}
		if existErr {
			return nil, e
		}

	loop2:
		for {
			select {
			case tasks := <-tasksch:
				if tasks == nil {
					continue loop2
				}
				matchTasks = append(matchTasks, tasks...)
			default:
				break loop2
			}
		}

	}
	//textRepにSearchしてヒットしたもの
	for _, textRep := range textReps {
		wg.Add(1)
		textRep := textRep
		go func(textRep text.TextRep) {
			defer wg.Done()
			matchTasks := []*mi.Task{}

			for _, word := range words {
				matchTexts, err := textRep.Search(ctx, word)
				if err != nil {
					err = fmt.Errorf("error at search %s in %s: %w", word, textRep.Path(), err)
					errcht <- err
				}
				for _, text := range matchTexts {
					for _, task := range allTasks {
						if task.TaskID == text.Target {
							matchTasks = append(matchTasks, task)
						}
					}
				}
				taskscht <- matchTasks
			}
		}(textRep)
	}
	wg.Wait()

errloop3:
	for {
		select {
		case e := <-errcht:
			e = fmt.Errorf("error at filter tags: %w", e)
			existErr = true
		default:
			break errloop3
		}
	}
	if existErr {
		return nil, e
	}

loop3:
	for {
		select {
		case tasks := <-taskscht:
			if tasks == nil {
				continue loop3
			}
			allTasks = append(allTasks, tasks...)
		default:
			break loop3
		}
	}

	// idが完全に一致するものも
	for _, task := range allTasks {
		for _, word := range words {
			if task.TaskID == word {
				matchTasks = append(matchTasks, task)
			}
		}
	}
	return matchTasks, nil
}

func andSearch(ctx context.Context, reps mi.MiReps, textReps []text.TextRep, words []string, query *mi.SearchTaskQuery) ([]*mi.Task, error) {
	existErr := false
	var e error

	// searchで見つかったかどうか := map[id]map[word]
	m := map[string]map[string]bool{}
	hitTasks := map[string]*mi.Task{}
	allTasks := []*mi.Task{}

	allTasksMap := map[string]*mi.Task{}

	wg := &sync.WaitGroup{}
	tasksch := make(chan map[string]*mi.Task, len(reps))
	errch := make(chan error, len(reps))
	defer close(tasksch)
	defer close(errch)
	taskscht := make(chan map[string]*mi.Task, len(textReps))
	errcht := make(chan error, len(textReps))
	defer close(taskscht)
	defer close(errcht)
	for _, rep := range reps {
		wg.Add(1)
		rep := rep
		go func(rep mi.MiRep) {
			defer wg.Done()
			tasksMap := map[string]*mi.Task{}

			tasks, err := rep.GetAllTasks(ctx)
			if err != nil {
				err = fmt.Errorf("error at get all task from %s: %w", rep.Path(), err)
				errch <- err
			}
			for _, task := range tasks {
				if _, exist := allTasksMap[task.TaskID]; !exist {
					tasksMap[task.TaskID] = task
				}
			}
			tasksch <- tasksMap
		}(rep)
	}
	for _, task := range allTasksMap {
		allTasks = append(allTasks, task)
	}
	wg.Wait()

errloop1:
	for {
		select {
		case e := <-errch:
			e = fmt.Errorf("error at filter tags: %w", e)
			existErr = true
		default:
			break errloop1
		}
	}
	if existErr {
		return nil, e
	}

loop1:
	for {
		select {
		case tasks := <-tasksch:
			if tasks == nil {
				continue loop1
			}
			for taskid, task := range tasks {
				allTasksMap[taskid] = task
			}
		default:
			break loop1
		}
	}

	for _, word := range words {
		for _, rep := range reps {
			tasks, err := rep.SearchTasks(ctx, word, query)
			if err != nil {
				err = fmt.Errorf("error at search %s from %s: %w", word, rep.RepName(), err)
				return nil, err
			}
			for _, task := range tasks {
				if _, exist := m[task.TaskID]; !exist {
					m[task.TaskID] = map[string]bool{}
				}
				m[task.TaskID][word] = true
				hitTasks[task.TaskID] = task
			}
		}
		for _, textRep := range textReps {
			texts, err := textRep.Search(ctx, word)
			if err != nil {
				err = fmt.Errorf("error at search %s from %s: %w", word, textRep.Path(), err)
				return nil, err
			}
			for _, text := range texts {
				texts, err := textRep.GetTextsByTarget(ctx, text.ID)
				if err != nil {
					err = fmt.Errorf("error at get texts by target %s from %s: %w", text.ID, textRep.Path(), err)
					return nil, err
				}

				for _, text := range texts {
					found := false
					for _, task := range allTasks {
						if task.TaskID == text.Target {
							found = true
							break
						}
					}
					if !found {
						// repが分散しているとtargetの存在しないtextが出現し得るのでその場合はcontinue
						continue
					}

					if _, exist := m[text.Target]; !exist {
						m[text.Target] = map[string]bool{}
					}
					m[text.Target][word] = true
				}
			}
		}
	}

	for _, word := range words {
		for _, wordMap := range m {
			if _, exist := wordMap[word]; !exist {
				wordMap[word] = false
			}
		}
	}

	tasks := []*mi.Task{}
	ids := []string{}
	for id, wordMap := range m {
		allMatch := true
		for _, exist := range wordMap {
			if !exist && allMatch {
				allMatch = false
				break
			}
		}
		if allMatch {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		task, exist := hitTasks[id]
		if !exist {
			found := false
			for _, k := range allTasks {
				if k.TaskID == task.TaskID {
					found = true
					task = k
					break
				}
			}
			if !found {
				err := fmt.Errorf("not found %s from all reps", id)
				return nil, err
			}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func equal(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func setEnv() {
	// HOME
	home := os.Getenv("HOME")
	if home == "" {
		home, err := homedir.Dir()
		if err != nil {
			err = fmt.Errorf("error at get user home directory: %w", err)
			log.Printf(err.Error())
		} else {
			os.Setenv("HOME", home)
		}
	}

	// EXE
	exe := os.Getenv("EXE")
	if exe == "" {
		exe, err := os.Executable()
		if err != nil {
			err = fmt.Errorf("error at get executable file path: %w", err)
			log.Printf(err.Error())
		} else {
			os.Setenv("EXE", exe)
		}
	}
}
