package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/beats/libbeat/paths"
	"github.com/stretchr/testify/assert"
)

func TestNewPHCreation(t *testing.T) {

	ph, err := NewPostionHandler("")
	assert.NotNil(t, ph)
	assert.Nil(t, err)

}

func TestReadPosition(t *testing.T) {
	ph, err := NewPostionHandler("")
	assert.NotNil(t, ph)
	assert.Nil(t, err)

	var position = "sometext"
	_, err1 := ph.WritePostiontoFile(position)
	assert.Nil(t, err1)

	var datafromFile = ""
	err2 := ph.ReadPositionfromFile(&datafromFile)
	assert.Nil(t, err2)
	assert.Equal(t, position, datafromFile)
}

func TestGetPostionFile(t *testing.T) {

	path, err := GetPostionFile()

	if err != nil {
		t.Errorf("Unable to get postion file path : %v", err)
	}
	ospath := paths.Paths.Home

	posFolder := filepath.Join(ospath, "logs")
	posfilePath := filepath.Join(posFolder, "pos.json")

	if path != posfilePath {
		t.Errorf("The expected %s path and actual %s path doesn't match", posfilePath, posfilePath)
	}

	created, err := exists(posFolder)

	if !created || err != nil {
		t.Errorf("Positon folder not created : %v", err)
	}

}

func TestWritePosition(t *testing.T) {
	ph, err := NewPostionHandler("")
	assert.NotNil(t, ph)
	assert.Nil(t, err)

	cases := []struct {
		name   string
		input  interface{}
		output bool
		err    error
	}{
		{"TestWritePosition", "sometext", true, nil},
		{"TestWritePosition_Exception", make(chan int), false, &json.UnsupportedTypeError{}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Running test case %s", test.name)
			output, err := ph.WritePostiontoFile(test.input)
			assert.IsType(t, test.err, err)
			assert.Equal(t, test.output, output)
		})
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}
