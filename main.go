package main

import (
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"rec/config"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const FILEPATH = "./config/reconfig.json"

var Recs [rows][cols]*walk.CheckBox
var Stats [rows]*walk.Label
var Rem [rows]*walk.Label
var SetAll [cols]*walk.CheckBox

type _Name struct {
	Name, Number, Rem             string
	NameSync, NumberSync, RemSync bool
}

type Data struct {
	Data string
}

// TODO if needed those can be loaded in the config as well
const rows = 24
const cols = 4

var kiproip []string
var NRKIPRO int
var Name []_Name
var UMDsync = &sync.Mutex{}

func main() {
	var mainWin *walk.MainWindow

	conf, _ := config.ReadConfig()
	kiproip = conf.KiproIps

	mwin := MainWindow{
		AssignTo: &mainWin,
		Title:    "REC Bucatari",
		MinSize:  Size{650, 400},
		MaxSize:  Size{650, 400},
		Size:     Size{650, 400},
		Layout:   HBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					Label{Text: "   ---:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   ---:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   PGM 1:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 01:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 02:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 03:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 04:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 05:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 06:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 07:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 08:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 09:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 10:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 11:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Cam 12:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 01:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 02:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 03:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 04:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 05:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 07:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 06:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 07:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 08:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Dome 09:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "   Martor:", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			VSpacer{MinSize: Size{10, 10}, MaxSize: Size{10, 10}},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},

				Children: []Widget{
					PushButton{Text: "START", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					PushButton{OnClicked: func() { BtnClicked(true, 1, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 2, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 3, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 4, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 5, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 6, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 7, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 8, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 9, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 10, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 11, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 12, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 13, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 14, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 15, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 16, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 17, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 18, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 19, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 20, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 21, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 22, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(true, 23, 0) }, Text: "Rec", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					PushButton{Text: "STOP", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					PushButton{OnClicked: func() { BtnClicked(false, 1, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 2, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 3, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 4, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 5, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 6, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 7, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 8, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 9, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 10, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 11, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 12, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 13, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 14, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 15, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 16, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 17, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 18, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 19, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 20, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 21, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 22, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{OnClicked: func() { BtnClicked(false, 23, 0) }, Text: "Stop", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					CheckBox{OnCheckStateChanged: func() { Setall(0) }, AssignTo: &SetAll[0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					CheckBox{AssignTo: &Recs[0][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[1][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[2][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[3][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[4][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[5][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[6][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[7][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[8][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[9][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[10][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[11][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[12][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[13][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[14][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[15][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[16][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[17][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[18][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[19][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[20][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[21][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[22][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[23][0], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Rec", OnClicked: func() { BtnClicked(true, 0, 1) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Stop", OnClicked: func() { BtnClicked(false, 0, 1) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					CheckBox{OnCheckStateChanged: func() { Setall(1) }, AssignTo: &SetAll[1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					CheckBox{AssignTo: &Recs[0][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[1][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[2][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[3][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[4][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[5][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[6][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[7][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[8][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[9][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[10][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[11][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[12][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[13][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[14][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[15][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[16][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[17][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[18][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[19][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[20][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[21][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[22][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[23][1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Rec", OnClicked: func() { BtnClicked(true, 0, 2) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Stop", OnClicked: func() { BtnClicked(false, 0, 2) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					CheckBox{OnCheckStateChanged: func() { Setall(2) }, AssignTo: &SetAll[2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					CheckBox{AssignTo: &Recs[0][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[1][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[2][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[3][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[4][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[5][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[6][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[7][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[8][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[9][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[10][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[11][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[12][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[13][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[14][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[15][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[16][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[17][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[18][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[19][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[20][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[21][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[22][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[23][2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Rec", OnClicked: func() { BtnClicked(true, 0, 3) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Stop", OnClicked: func() { BtnClicked(false, 0, 3) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					CheckBox{OnCheckStateChanged: func() { Setall(3) }, AssignTo: &SetAll[3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					CheckBox{AssignTo: &Recs[0][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[1][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[2][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[3][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[4][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[5][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[6][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[7][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[8][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[9][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[10][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[11][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[12][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[13][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[14][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[15][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[16][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[17][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[18][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[19][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[20][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[21][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[22][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					CheckBox{AssignTo: &Recs[23][3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Rec", OnClicked: func() { BtnClicked(true, 0, 4) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{Text: "Stop", OnClicked: func() { BtnClicked(false, 0, 4) }, MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					Label{Text: "", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					Label{AssignTo: &Stats[0], MinSize: Size{100, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[4], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[5], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[6], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[7], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[8], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[9], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[10], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[11], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[12], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[13], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[14], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[15], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[16], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[17], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[18], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[19], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[20], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[21], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[22], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Stats[23], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					PushButton{
						Text: "Open Config",
						OnClicked: func() {
							changesMade, conf := config.ShowConfigWindow(mainWin)
							if changesMade {
								kiproip = conf.KiproIps
							}
						},
					},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					Label{Text: "", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 10}, MaxSize: Size{0, 10}},
					Label{AssignTo: &Rem[0], MinSize: Size{50, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[1], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[2], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[3], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[4], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[5], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[6], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[7], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[8], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[9], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[10], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[11], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[12], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[13], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[14], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[15], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[16], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[17], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[18], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[19], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[20], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[21], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[22], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{AssignTo: &Rem[23], MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
					Label{Text: "", MinSize: Size{0, 25}, MaxSize: Size{0, 25}},
				},
			},
			VSpacer{MinSize: Size{5, 0}},
		},
	}
	go GetStateLoop()
	mwin.Run()
}

func BtnClicked(rec bool, row, col int) {
	if col == 0 {
		Rec(rec, row-1)
	}
	if row == 0 {
		for i := 0; i < rows; i++ {
			if Recs[i][col-1].CheckState() == 1 {
				Rec(rec, i)
			}
			//Clip List
			if i == 0 {
				go GetList(kiproip[i])
			}
		}
	}
}

func Rec(rec bool, recnr int) {
	if recnr == rows-1 { //este martorul
		conn, err := net.Dial("tcp", "172.22.0.113:6789")
		if err != nil {
			println("Eroare conectare la martor")
		}
		encoder := gob.NewEncoder(conn)
		data := *new(Data)
		if rec == true {
			data.Data = "rec"
		} else {
			data.Data = "stop"
		}
		encoder.Encode(data)
		_ = conn.Close()
	} else if recnr == rows-2 { //este vtr
		println("VTR")
	} else {
		if rec == true {
			go http.Get("http://" + kiproip[recnr] + "/config?action=set&paramid=eParamID_TransportCommand&value=3")
			go PrintUMD(recnr)
		} else {
			go http.Get("http://" + kiproip[recnr] + "/config?action=set&paramid=eParamID_TransportCommand&value=4")
		}
		go GetState(recnr)
	}
}

func Setall(col int) {
	if SetAll[col].CheckState() == 1 {
		for i := 0; i < rows-3; i++ {
			Recs[i][col].SetChecked(true)

		}
	} else {
		for i := 0; i < rows; i++ {
			Recs[i][col].SetChecked(false)
		}
	}
}

// ********** State **********
func GetState(recnr int) {
	resp, err := http.Get("http://" + kiproip[recnr] + "/options?eParamID_TransportState")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_body, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return
	}
	// separ optiunile
	body := string(_body)
	var opts []string
	_opts1 := strings.Split(body, "{")
	for i := range _opts1 {
		_opts2 := strings.Split(_opts1[i], "}")
		opts = append(opts, _opts2...)
	}
	status := ""
	//caut optiunea selectata
	for _, opt := range opts {
		if strings.Contains(opt, "selected:\"true\"") {
			elem := "short_text:"
			index := strings.Index(opt, elem)
			if index >= 0 {
				opt = opt[index:]
				opt = opt[strings.Index(opt, "\"")+1:]
				status = opt[:strings.Index(opt, "\"")]
			}
		}
	}
	Stats[recnr].SetText(status)
}

func PrintUMD(i int) {
	go func() {
		Name[i].NameSync = false
		val1, err1 := GetText(kiproip[i], "eParamID_CustomClipName")
		if err1 != nil {
			Name[i].Name = "???"
		}
		Name[i].Name = val1
		Name[i].NameSync = true
	}()

	go func() {
		Name[i].NumberSync = false
		val2old := ""
		for lm := 0; lm < 50; lm++ {
			time.Sleep(100 * time.Millisecond)
			val2, err2 := GetText(kiproip[i], "eParamID_CustomTake")
			if err2 != nil {
				Name[i].Number = "???"
			} else {
				val, err3 := strconv.Atoi(val2)
				if err3 != nil {
					val2 = "???"
				} else {
					val2 = strconv.Itoa(val - 1)
				}
				Name[i].Number = val2
				Name[i].NumberSync = true

				if val2old == "" || val2old == "???" {
					val2old = val2
				} else {
					if val2 != val2old && val2 != "???" && val2old != "???" {
						//go SendUMD(i)
						break
					}
				}
			}
		}
	}()

	go SendUMD(i)
	//go GetStateLoop()
}

func GetText(ip, param string) (string, error) {
	resp, err := http.Get("http://" + ip + "/options?" + param)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}
	if j := strings.Index(string(body), "text:\""); j >= 0 {
		if k := strings.Index(string(body[j+6:]), "\""); k >= 0 {
			return string(body[j+6 : j+6+k]), nil
		}
	}
	return "", errors.New("Unknown error")
}

func SendUMD(nr int) error {
	for slp := 0; Name[nr].NameSync == false || Name[nr].NumberSync == false || slp < 20; slp++ {
		time.Sleep(100 * time.Millisecond)
	}
	if Name[nr].NameSync == false || Name[nr].NumberSync == false {
		return errors.New("Eroare send umd")
	}
	UMDsync.Lock()
	defer UMDsync.Unlock()
	conn, err := net.Dial("udp", "172.22.30.82:9801")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := Name[nr].Name + "_" + Name[nr].Number
	if len(text) > 16 {
		text = text[:13-len(Name[nr].Number)] + "..." + text[len(text)-len(Name[nr].Number):]
	}
	for len(text) < 16 {
		text += " "
	}

	cmd := make([]byte, 2)
	cmd[0] = 128
	for i := 0; i < nr; i++ {
		cmd[0]++
	}
	cmd[1] = 0
	cmd = append(cmd, []byte(text)...)
	err1 := binary.Write(conn, binary.BigEndian, cmd)
	if err1 != nil {
		return err1
	}
	return nil
}

func SendUMDRem(nr int, val string) error {
	// for slp := 0; Name[nr].RemSync == false || slp < 20; slp++ {
	// 	time.Sleep(100 * time.Millisecond)
	// }
	// if Name[nr].RemSync == false || Name[nr].RemSync == false {
	// 	return errors.New("Eroare SendUMD")
	// }
	if Name[nr].Rem == val {
		return nil
	}
	Name[nr].Rem = val
	UMDsync.Lock()
	defer UMDsync.Unlock()
	conn, err := net.Dial("udp", "172.22.30.82:9801")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := val + "%"
	// if len(text) > 16 {
	// 	text = text[:13-len(Name[nr].Number)] + "%" + text[len(text)-len(Name[nr].Number):]
	// }
	for len(text) < 16 {
		text += " "
	}

	cmd := make([]byte, 2)
	cmd[0] = 192
	for i := 0; i < nr; i++ {
		cmd[0]++
	}
	cmd[1] = 0
	cmd = append(cmd, []byte(text)...)
	err1 := binary.Write(conn, binary.BigEndian, cmd)
	if err1 != nil {
		return err1
	}
	return nil
}

var faster = false
var faster_lock = false

func Faster() {
	ttf := 30 * time.Second
	for faster_lock {
		time.Sleep(50 * time.Millisecond)
		ttf -= 50 * time.Millisecond
	}
	faster = true
	faster_lock = true
	if ttf > 0 {
		time.Sleep(ttf * time.Second)
	}
	faster = false
	faster_lock = false
}

func GetStateLoop() {
	for {
		if faster {
			time.Sleep(200 * time.Millisecond)
		} else {
			time.Sleep(2 * time.Second)
		}
		for i := 0; i < NRKIPRO; i++ {
			time.Sleep(10 * time.Millisecond)
			go GetState(i)
			go GetRemaining(i)
		}
	}

}

func GetRemaining(recnr int) {
	resp, err := http.Get("http://" + kiproip[recnr] + "/options?eParamID_CurrentMediaAvailable")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_body, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return
	}
	// separ optiunile
	body := string(_body)
	status := ""
	//caut optiunea selectata
	elem := "value:"
	index := strings.Index(body, elem)
	if index >= 0 {
		body = body[index:]
		body = body[strings.Index(body, "\"")+1:]
		status = body[:strings.Index(body, "\"")]
	}
	Rem[recnr].SetText(status + "%")
	go SendUMDRem(recnr, status)
}

// ************ Clip List ***************
type attributes struct {
	Starting, Format string
}

type clip struct {
	Clipname, Duration string
	Attributes         attributes

	Name, StartTC, StopTC string
}

type clips struct {
	Clips []clip
}

var Clips clips

func GetList(ip string) clips {
	resp, err := http.Get("http://" + ip + "/clips?action=get_clips")
	if err != nil {
		return *new(clips)
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return *new(clips)
	}
	_body := string(body)
	_body = strings.Replace(_body, "Starting TC", "Starting", 100)
	body = []byte(_body)
	err2 := json.Unmarshal(body, &Clips)
	if err2 != nil {
		println(err2.Error())
	}
	for i := range Clips.Clips {
		Clips.Clips[i].Name = Clips.Clips[i].Clipname[:strings.Index(Clips.Clips[i].Clipname, ".mov")]
		if len(Clips.Clips[i].Attributes.Starting) > 8 {
			Clips.Clips[i].StartTC = Clips.Clips[i].Attributes.Starting[:8]
		}
		secsStart := StringsToSecs(Clips.Clips[i].Attributes.Starting)
		secsDur := StringsToSecs(Clips.Clips[i].Duration)
		secsStop := secsStart + secsDur
		sec := secsStop % 60
		min := ((secsStop - sec) % 3600) / 60
		hrs := (secsStop - (min * 60) - sec) / 3600
		Clips.Clips[i].StopTC = strconv.Itoa(hrs) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)

		println(Clips.Clips[i].Name)
		println(Clips.Clips[i].StartTC)
		println(Clips.Clips[i].StopTC)
	}
	println(len(Clips.Clips))
	return *new(clips)
}

func StringsToSecs(s string) int {
	times := strings.Split(s, ":")
	secs := 0
	if len(times) > 3 {
		secs = (Atoi(times[0]) * 3600) + (Atoi(times[1]) * 60) + Atoi(times[2])
	}
	return secs
}

func Atoi(s string) int {
	_s, _ := strconv.Atoi(s)
	return _s
}
