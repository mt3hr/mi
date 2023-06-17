package miapp

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	TagStructFile string
)

func getTagStructFile() string {
	return TagStructFile
}

func getDefaultTagStructFileName() string {
	return "$HOME/tag_struct.yaml"
}

// TagStructの順番を保証するために
func loadTagStructFromFile() error {
	tagStructFileName := getTagStructFile()
	if tagStructFileName == "" {
		tagStructFileName = getDefaultTagStructFileName()
	}
	tagStructFileName = os.ExpandEnv(tagStructFileName)

	_, err := os.Stat(tagStructFileName)
	if err != nil {
		content := `TagStruct: 
    notag: tag
`
		_, err = os.Create(tagStructFileName)
		if err != nil {
			return err
		}
		err = os.WriteFile(tagStructFileName, []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	b, err := os.ReadFile(tagStructFileName)
	if err != nil {
		err = fmt.Errorf("error at read file %s: %w", tagStructFileName, err)
		return err
	}

	m := yaml.MapSlice{}
	tagStructMap := yaml.MapSlice{}
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		err = fmt.Errorf("error at yaml unmarshall: %w", err)
		return err
	}
	for _, v := range m {
		if v.Key == "TagStruct" {
			var ok bool
			tagStructMap, ok = v.Value.(yaml.MapSlice)
			if !ok {
				err = fmt.Errorf("TagStruct.yamlファイルが変です")
				return err
			}
		}
	}
	LoadedConfig.ApplicationConfig.TagStruct = MapSlice(tagStructMap)
	return nil
}
