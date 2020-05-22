package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/paths"
)

var lock sync.Mutex

//PositionHandler holds the position of read events from flat file
type PositionHandler struct {
	PostionFilePath string
	rwmutex         sync.RWMutex
}

//NewPostionHandler create new position handler object
func NewPostionHandler(posFileName string) (*PositionHandler, error) {

	posfilepath, err := GetPostionFile()

	if err != nil {
		logp.Err("Unable to create position handler%v", err)
		return nil, err
	}

	pos := new(PositionHandler)
	pos.PostionFilePath = posfilepath
	logp.Debug("SophosCentralBeat", "Position file path is  : %s", posfilepath)
	return pos, err
}

//ReadPositionfromFile read position of events from where next read will start
func (position *PositionHandler) ReadPositionfromFile(v interface{}) error {

	file := position.PostionFilePath

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)

	if err != nil {

		logp.Err("Unable to read postion from file %v", err)
	}

	return err

}

//WritePostiontoFile write position of current read events.
func (position *PositionHandler) WritePostiontoFile(v interface{}) (bool, error) {

	file := position.PostionFilePath

	jtkn, err := json.Marshal(v)

	if err != nil {
		logp.Err("Unable to save postion to file %v", err)
		return false, err
	}

	ioutil.WriteFile(file, jtkn, 0644)
	return true, err

}

//GetPostionFile is used to get position file path
func GetPostionFile() (string, error) {

	tokenCacheDir := paths.Paths.Home

	posFolder := filepath.Join(tokenCacheDir, "logs")

	err := os.MkdirAll(posFolder, 0700)
	if err != nil {
		return "", err
	}

	return filepath.Join(posFolder, url.QueryEscape("pos.json")), nil
}
