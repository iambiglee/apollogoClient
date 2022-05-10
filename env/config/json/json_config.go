package json

import (
	"encoding/json"
	"errors"
	"github.com/apollogoClient/v1/utils"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
}

func (c ConfigFile) Load(fileName string, unmarshal func([]byte) (interface{}, error)) (interface{}, error) {
	fs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Fail to read files:" + err.Error())
	}
	config, loadErr := unmarshal(fs)

	if utils.IsNotNil(loadErr) {
		return nil, errors.New("Load Json Config fail:" + loadErr.Error())
	}

	return config, nil

}

func (c ConfigFile) Write(content interface{}, configPath string) error {
	if content == nil {
		return errors.New("content is null can not write backup file")
	}
	file, e := os.Create(configPath)
	if e != nil {
		return e
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(content)
}
