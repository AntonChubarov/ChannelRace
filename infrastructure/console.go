package infrastructure

import (
	"RacersRace/domain"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
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
		time.Sleep(domain.LoopSleepTime)
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
	for i := 0; i < domain.StepsInLap; i++ {
		if i < racerInfo.Score % domain.StepsInLap {
			racerString += "-"
		}
		if i == racerInfo.Score % domain.StepsInLap {
			racerString += ">"
		}
		if i > racerInfo.Score % domain.StepsInLap {
			racerString += " "
		}
	}
	racerString += " | "
	racerString += "Lap: "
	racerString += strconv.Itoa(racerInfo.Lap)
	fmt.Println(racerString)
}


// Console clear

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Println("unable to perform console command")
		}
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Println("unable to perform console command")
		}
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else { // unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

