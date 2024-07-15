package config

import (
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
)

type Config struct {
	KiproIps []string `json:"kiproip"`
	// TODO ask razvan if he wants more things added to the config file
}

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

	conf, err := ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if len(conf.KiproIps) != 21 {
		log.Fatalf("Expected 21 IP addresses, got %d", len(conf.KiproIps))
	}

	kiproip := conf.KiproIps
	originalValues := make([]string, len(kiproip))
	copy(originalValues, kiproip)

	var ipEdits [21]*walk.LineEdit
	changesMade := false

	children := make([]Widget, 22) // 21 LineEdits + 1 Save Button

	for i := 0; i < 21; i++ {
		index := i // capture loop variable
		children[i] = LineEdit{
			AssignTo: &ipEdits[index],
			Text:     originalValues[index],
			OnTextChanged: func() {
				for j := 0; j < 21; j++ {
					if ipEdits[j].Text() != originalValues[j] {
						changesMade = true
						break
					}
				}
				saveButton.SetEnabled(changesMade)
			},
		}
	}

	children[21] = PushButton{
		AssignTo: &saveButton,
		Text:     "Save",
		Enabled:  false,
		OnClicked: func() {
			for i := 0; i < 21; i++ {
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

func main() {
	ShowConfigWindow(nil)
}
