package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"transifex"
	"transifex/format"
)

type configElement struct {
	Type      string             `json:"type"`
	Structure string             `json:"structure"`
	Resources []LocalizationFile `json:"resources"`
}

type LocalizationFile struct {
	transifex.BaseResource
	Fname        string   `json:"fname"`
	Categories   []string `json:"categories"`
	Dir          string   `json:"dir"`
	Translations map[string]string
	Format       format.Format
	FileLocator  format.FileLocator
}

func (f *LocalizationFile) init(rootDir string, elem configElement) error {
	f.I18nType = elem.Type
	f.Category = strings.Join(f.Categories, " ")
	f.Format = format.Formats[elem.Type]
	f.FileLocator = format.FileLocators[elem.Structure]

	var readErr error
	f.Translations, readErr = f.FileLocator.List(filepath.Join(rootDir, f.Dir), f.Fname, f.Format.Ext())

	if readErr != nil {
		return readErr
	}

	return nil
}

// type LocalizationFile struct {
// 	transifex.BaseResource
// 	Filename     string
// 	Structure    format.FileLocator
// 	Format       format.Format
// 	Translations map[string]string
// }

func ReadConfig(configFile, rootDir, sourceLang string) (files []LocalizationFile, err error) {
	if sourceLang == "" {
		return nil, fmt.Errorf("Source lang is empty.")
	}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Unable to read %s", configFile)
		return nil, err
	}

	var jsonData []configElement
	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		return nil, err
	}

	files = []LocalizationFile{}
	for _, elem := range jsonData {
		for _, f := range elem.Resources {
			if err = f.init(rootDir, elem); err != nil {
				return nil, err
			}
			files = append(files, f)
		}
	}

	return files, nil
}
