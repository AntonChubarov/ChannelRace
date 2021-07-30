package app

import (
	"RacersRace/domain"
	"fmt"
	"sort"
	"sync"
	"time"
)

type JudgeOfRace struct {
	RacersInfo      []domain.RacerInfo
	StepCommander   *StepCommander
	DisplayCommander *DisplayCommander
	InfoChannels    []chan domain.RacerInfo
	InfoCollector *InfoCollector
	StopChannel     chan bool
	IsInactiveRacer []bool
	InactiveCount   int
	MutexRacersInfo *sync.RWMutex
}

func NewRaceJudge(
	stepChannel []chan time.Time,
	infoChannels []chan domain.RacerInfo,
	displayChannel chan []domain.RacerInfo,
	stopChannel chan bool,
	) *JudgeOfRace {
	ri := make([]domain.RacerInfo, len(infoChannels))
	isInactive := make([]bool, len(infoChannels))
	mu := sync.RWMutex{}
	return &JudgeOfRace{
		RacersInfo:       ri,
		StepCommander:    NewStepCommander(stepChannel),
		DisplayCommander: NewDisplayCommander(&ri, displayChannel, &mu),
		InfoChannels:     infoChannels,
		StopChannel:      stopChannel,
		IsInactiveRacer:  isInactive,
		MutexRacersInfo:  &mu,
		InfoCollector:    NewInfoCollector(&ri, &infoChannels, &isInactive, &mu),
	}
}

func (j *JudgeOfRace) StartRace() {
	go j.StepCommander.Run()
	go j.DisplayCommander.Run()
	go j.InfoCollector.Run()
	go j.startToJudge()
}

func (j *JudgeOfRace) startToJudge() {
	sortedInfo := make([]domain.RacerInfo, len(j.RacersInfo))
	var nameOfRacerToStop string
	for {
		time.Sleep(domain.LoopSleepTime)
		j.MutexRacersInfo.RLock()

		copy(sortedInfo, j.RacersInfo)

		for i := range j.RacersInfo {
			if !j.IsInactiveRacer[i] && j.RacersInfo[i].StepInLap >= 20 {
				j.IsInactiveRacer[i] = true
				j.InactiveCount++
				fmt.Println(j.RacersInfo[i].Name, "was too slow!")
				if j.InactiveCount == len(j.InfoChannels)-1 {
					break
				}
			}
		}

		j.MutexRacersInfo.RUnlock()
		sort.SliceStable(sortedInfo, func(i, j int) bool {
			return sortedInfo[i].Score > sortedInfo[j].Score
		})

		if j.InactiveCount < len(j.InfoChannels)-1 {
			if sortedInfo[len(sortedInfo)-1-j.InactiveCount].Lap < sortedInfo[len(sortedInfo)-2-j.InactiveCount].Lap {
				nameOfRacerToStop = sortedInfo[len(sortedInfo)-1-j.InactiveCount].Name
				racerIndex := j.findRacerIndexByName(nameOfRacerToStop)
				j.IsInactiveRacer[racerIndex] = true
				j.InactiveCount++
			}
		}

		if j.InactiveCount == len(j.InfoChannels)-1 {
			j.MutexRacersInfo.RLock()
			fmt.Println("\nThe winner is " + sortedInfo[0].Name)
			for i := range sortedInfo{
				fmt.Println(sortedInfo[i].Name, "Score:", sortedInfo[i].Score)
			}
			j.MutexRacersInfo.RUnlock()
			j.StopChannel <- true
		}
	}
}

func (j *JudgeOfRace) findRacerIndexByName (name string) int {
	j.MutexRacersInfo.RLock()
	for i := range j.RacersInfo {
		if j.RacersInfo[i].Name == name {
			j.MutexRacersInfo.RUnlock()
			return i
		}
	}
	panic(fmt.Errorf("racer not found"))
}
