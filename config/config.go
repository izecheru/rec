package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Config struct {
	KiproIps []string `json:"kiproip"`
	// TODO ask razvan if he wants more things added to the config file
}

var kiprolen int

const FILEPATH = "./config/config.json"

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(FILEPATH)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(config Config) error {
	var filteredKiproIps []string
	for _, kiproIp := range config.KiproIps {
		// if we have an empty string we delete that entry
		if kiproIp != "" {
			filteredKiproIps = append(filteredKiproIps, kiproIp)
		}
	}
	config.KiproIps = filteredKiproIps
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling config: %v", err)
	}

	if err := os.WriteFile(FILEPATH, data, 0644); err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}

	return err
}

func ShowConfigWindow(owner walk.Form) (bool, *Config) {
	var configWin *walk.MainWindow
	var saveButton *walk.PushButton
	var addEntry *walk.PushButton

	conf, err := ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	kiproip := conf.KiproIps
	kiprolen := len(kiproip)
	originalValues := make([]string, kiprolen)
	copy(originalValues, kiproip)
	ipEdits := make([]*walk.LineEdit, kiprolen)
	changesMade := false

	children := make([]Widget, kiprolen+2) // 21 LineEdits + 1 Save Button + 1 add entry button

	for i := 0; i < kiprolen; i++ {
		index := i // capture loop variable
		children[i] = LineEdit{
			AssignTo: &ipEdits[index],
			Text:     originalValues[index],
			OnTextChanged: func() {
				for j := 0; j < kiprolen; j++ {
					if ipEdits[j].Text() != originalValues[j] {
						changesMade = true
						break
					}
				}
				saveButton.SetEnabled(changesMade)
			},
		}
	}

	children[kiprolen] = PushButton{
		AssignTo: &saveButton,
		Text:     "Save",
		Enabled:  false,
		OnClicked: func() {
			for i := 0; i < kiprolen; i++ {
				conf.KiproIps[i] = ipEdits[i].Text()
			}

			err := WriteConfig(*conf)
			if err != nil {
				log.Fatalf("Error writing config: %v", err)
				return
			}
			configWin.Close()
		},
	}

	children[kiprolen+1] = PushButton{
		AssignTo: &addEntry,
		Text:     "Add entry",
		Enabled:  true,
		OnClicked: func() {
			kiprolen += 1
			ipEdits := make([]walk.LineEdit, kiprolen)
			newKiproips := make([]string, kiprolen)
			copy(newKiproips, kiproip)
			kiproip = newKiproips
			fmt.Println("added 1 more ip to the list %v", len(ipEdits))

			err := WriteConfig(*conf)
			if err != nil {
				log.Fatalf("Error writing config: %v", err)
				return
			}
		},
	}

	MainWindow{
		AssignTo: &configWin,
		Title:    "Configuration",
		Size:     Size{Width: 400, Height: 800},
		Layout:   VBox{},
		Children: children,
	}.Create()

	configWin.Run()
	return changesMade, conf
}
