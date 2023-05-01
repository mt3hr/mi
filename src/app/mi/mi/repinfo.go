package miapp

import (
	"fmt"
	"os"

	"github.com/mattn/go-zglob"
	mi "github.com/mt3hr/mi/src/app"
)

// MiRepFactories . 詳しくはLoadMiRepsを参照してください。
var MiRepFactories = map[string]MiRepFactory{}

// MiRepInfo . 詳しくはLoadMiRepsを参照してください。
type MiRepInfo struct {
	Type string
	File string
}

// MiRepFactory . 詳しくはLoadMiRepsを参照してください。
type MiRepFactory func(contentFile string) ([]mi.MiRep, error)

// LoadMiReps .
// 使い方を示します。
//
// まず使いたいRepInfoを想定します。
// 例えば次のようにします。
//
//	repInfo := RepInfo {
//		type: "db",
//		file: "hoge.db",
//	}
//
// 次に、そのrepInfoが読み込めるように、RepFactoryを作成し、RepFactoriesにを登録します。
// contentFileにはrepInfo.fileが渡されます。
// repFactory := func(contentFile string) ([]Rep, error) { //hogefuga }
// RepFactories[repInfo.type] = repFactory
//
// 最後に、LoadRepsにrepInfoをわたしてRepオブジェクトを読み込みます。
// reps, err := LoadReps(repInfo)
func LoadMiReps(repInfo *MiRepInfo) ([]mi.MiRep, error) {
	reps := []mi.MiRep{}
	factory, exist := MiRepFactories[repInfo.Type]
	if !exist {
		err := fmt.Errorf("unknown rep type %s", repInfo.Type)
		return nil, err
	}
	rep, err := factory(repInfo.File)
	if err != nil {
		err = fmt.Errorf("failed to load rep %s: %w", repInfo.File, err)
		return nil, err
	}
	reps = append(reps, rep...)
	return reps, nil
}

func init() {
	setEnv()
	registMiDirectoryToFactories()
}

func registMiDirectoryToFactories() {
	MiRepFactories["mi_db"] = func(contentFile string) ([]mi.MiRep, error) {
		reps := []mi.MiRep{}

		contentFile = os.ExpandEnv(contentFile)
		matches, _ := zglob.Glob(contentFile)
		for _, match := range matches {
			rep, err := mi.NewMiRepSQLite(match)
			if err != nil {
				err = fmt.Errorf("failed to load mi rep dir %s: %w", match, err)
				return nil, err
			}
			reps = append(reps, rep)
		}
		return reps, nil
	}
}
