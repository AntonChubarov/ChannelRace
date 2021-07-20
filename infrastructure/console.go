package infrastructure

import (
	"RacersRace/domain"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Console struct {
	DisplayChannel chan []domain.RacerInfo
}

func NewConsole(displayChannel chan []domain.RacerInfo) *Console {
	return &Console{
		DisplayChannel: displayChannel,
	}
}

func (c Console) StartShowRaceSatus() {
	var info []domain.RacerInfo
	for {
		info = <- c.DisplayChannel
		CallClear()
		showRaceInfo(info)
	}
}

func showRaceInfo(info []domain.RacerInfo)  {
	for i := range info {
		printRacerString(info[i])
	}
}

func printRacerString(racerInfo domain.RacerInfo) {
	racerString := ""
	racerString += racerInfo.Name
	racerString += " |"
	for i := 0; i < 50; i++ {
		if i < racerInfo.Score % 50 {
			racerString += "-"
		}
		if i == racerInfo.Score % 50 {
			racerString += ">"
		}
		if i > racerInfo.Score % 50 {
			racerString += " "
		}
	}
	racerString += " | "
	racerString += "Lap: "
	racerString += strconv.Itoa(racerInfo.Lap)
	fmt.Println(racerString)
}


// Console clear

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

