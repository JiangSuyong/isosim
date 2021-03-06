package server

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"isosim/web/data"
)

type dataSetManager struct{}

var instance *dataSetManager
var dataDir string

const defFileSuffix = ".srvdef.json"

var initDS = sync.Once{}

// Init verifies and initializes the dataDir passed in during in
// initialization
func Init(dirname string) error {
	dir, err := os.Open(dirname)
	if err != nil {
		return err
	}
	_ = dir.Close()
	dataDir = dirname
	return nil
}

// DataSetManager returns the singleton instance of the DataSetManager
func DataSetManager() *dataSetManager {

	initDS.Do(func() {
		instance = new(dataSetManager)

	})
	return instance
}

// ErrDataSetExists is an error that indicates that the dataset by the provided
// name already exists
var ErrDataSetExists = errors.New("isosim: data set exists")

// GetAll returns a list of all data sets (names only) for the given specId
// and msgId
func (dsm *dataSetManager) GetAll(specId string, msgId string) ([]string, error) {

	dir, err := os.Open(filepath.Join(dataDir, specId, msgId))
	if err != nil {
		return nil, err
	}

	if dirContents, err := dir.Readdir(-1); err != nil {
		return nil, err
	} else {
		var dataSets = make([]string, 0, 10)
		for _, ds := range dirContents {
			if !ds.IsDir() {
				dataSets = append(dataSets, ds.Name())
			}
		}

		return dataSets, nil

	}
}

// Get returns the content of a specific data set
func (dsm *dataSetManager) Get(specId string, msgId string, dsName string) ([]byte, error) {

	dsData, err := ioutil.ReadFile(filepath.Join(dataDir, specId, msgId, dsName))
	if err != nil {
		return nil, err

	}
	return dsData, nil

}

// Add add a new data-set for the given spec and msg
func (dsm *dataSetManager) Add(specId string, msgId string, name string, data string) error {

	log.Traceln("Adding data set - " + name + " data = " + data)
	exists, err := checkIfExists(specId, msgId, name)
	if err != nil {
		return err
	}
	if exists {
		return ErrDataSetExists
	}

	err = ioutil.WriteFile(filepath.Join(dataDir, specId, msgId, name), []byte(data), 0755)
	if err != nil {
		return err
	}
	return nil
}

// AddServerDef adds a new server definition
func (dsm *dataSetManager) AddServerDef(defString string) (string, error) {

	log.Traceln("Adding server definition - .. JSON = " + defString)

	serverDef := &data.ServerDef{}
	err := json.NewDecoder(bytes.NewBufferString(defString)).Decode(serverDef)
	if err != nil {
		return "", err
	}

	strSpecId := strconv.Itoa(serverDef.SpecId)
	dir, err := os.Open(filepath.Join(dataDir, strSpecId))
	if err != nil && os.IsNotExist(err) {
		//create dir if one doesn't exist
		os.Mkdir(filepath.Join(dataDir, strSpecId), 0755)
	} else {
		if err != nil {
			return "", err
		}
	}
	dir.Close()

	fileName := serverDef.ServerName
	fileName = strings.Replace(fileName, " ", "", -1)
	fileName = strings.Replace(fileName, ",", "", -1)
	fileName = fileName + defFileSuffix

	log.Debugln("Writing spec def to file = " + fileName)

	defFile, err := os.Open(filepath.Join(dataDir, strSpecId, fileName))
	if err == nil {
		return "", errors.New("isosim: server-def file exists")
	}
	defFile.Close()

	err = ioutil.WriteFile(filepath.Join(dataDir, strSpecId, fileName), []byte(defString), 0755)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// ServerDefinitions returns all existing server definitions
func (dsm *dataSetManager) ServerDefinitions(specId string) ([]string, error) {

	dir, err := os.Open(filepath.Join(dataDir, specId))
	if err != nil {
		return nil, err
	}
	dirContents, err := dir.Readdir(-1)
	res := make([]string, 0)

	for _, fileInfo := range dirContents {
		if strings.HasSuffix(fileInfo.Name(), defFileSuffix) {
			res = append(res, fileInfo.Name())
		}
	}

	return res, nil

}

// ServerDef returns a server definition by its name
func (dsm *dataSetManager) ServerDef(specId string, name string) ([]byte, error) {

	file, err := os.Open(filepath.Join(dataDir, specId, name))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)

}

// Update updates a data set
func (dsm *dataSetManager) Update(specId string, msgId string, name string, data string) error {

	log.Debugln("Updating data set - " + name + " data = " + data)

	err := ioutil.WriteFile(filepath.Join(dataDir, specId, msgId, name), []byte(data), 0755)
	if err != nil {
		return err
	}
	return nil
}

func checkIfExists(specId string, msgId string, name string) (bool, error) {

	//check if the dir exists for this spec and msg
	//and if not create one first

	dir, err := os.Open(filepath.Join(dataDir, specId, msgId))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(dataDir, specId, msgId), 0755)
		if err != nil {
			return false, err
		}
		dir, err = os.Open(filepath.Join(dataDir, specId, msgId))
		if err != nil {
			return false, err
		}

	}

	dirContents, err := dir.Readdir(-1)
	if err != nil {
		return false, err
	}
	for _, fileInfo := range dirContents {
		if fileInfo.Name() == name {
			return true, nil
		}
	}

	return false, nil

}
