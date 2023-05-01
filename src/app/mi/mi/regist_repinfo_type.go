package miapp

import (
	"fmt"
	"os"
	"sort"

	"github.com/mattn/go-zglob"
	mi "github.com/mt3hr/mi/src/app"
	"github.com/mt3hr/rykv/tag"
	"github.com/mt3hr/rykv/text"
)

func init() {
	setEnv()

	MiRepFactories["mi_db"] = func(contentFile string) ([]mi.MiRep, error) {
		reps := []mi.MiRep{}

		contentFile = os.ExpandEnv(contentFile)
		matches, _ := zglob.Glob(contentFile)
		for _, match := range matches {
			rep, err := mi.NewMiRepSQLite(match)
			if err != nil {
				err = fmt.Errorf("failed to NewMiRepDir %s: %w", match, err)
				return nil, err
			}
			reps = append(reps, rep)
		}
		return reps, nil
	}

	tag.TagRepFactories["db"] = func(contentFile string) ([]tag.TagRep, error) {
		reps := []tag.TagRep{}

		contentFile = os.ExpandEnv(contentFile)
		matches, _ := zglob.Glob(contentFile)
		sort.Strings(matches)
		for _, match := range matches {
			rep, err := tag.NewTagRepSQLite(match)
			if err != nil {
				err = fmt.Errorf("error at new tag rep sqlite %s: %w", match, err)
				return nil, err
			}
			reps = append(reps, rep)
		}
		return reps, nil
	}

	text.TextRepFactories["db"] = func(contentFile string) ([]text.TextRep, error) {
		reps := []text.TextRep{}

		contentFile = os.ExpandEnv(contentFile)
		matches, _ := zglob.Glob(contentFile)
		sort.Strings(matches)
		for _, match := range matches {
			rep, err := text.NewTextRepSQLite(match)
			if err != nil {
				err = fmt.Errorf("error at new text rep sqlite %s: %w", match, err)
				return nil, err
			}
			reps = append(reps, rep)
		}
		return reps, nil
	}
}
