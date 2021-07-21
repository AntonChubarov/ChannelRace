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
	StepTicker      *time.Ticker
	StepChannel     []chan time.Time
	DisplayTicker   *time.Ticker
	DisplayChannel  chan []domain.RacerInfo
	InfoChannels    []chan domain.RacerInfo
	StopChannel     chan bool
	IsInactiveRacer []bool
	InactiveCount   int
	MutexRacersInfo sync.RWMutex
}

func NewRaceJudge(
	stepChannel []chan time.Time,
	infoChannels []chan domain.RacerInfo,
	displayChannel chan []domain.RacerInfo,
	stopChannel chan bool,
	) *JudgeOfRace {
	return &JudgeOfRace{
		RacersInfo:      make([]domain.RacerInfo, len(infoChannels)),
		StepChannel:     stepChannel,
		DisplayChannel:  displayChannel,
		InfoChannels:    infoChannels,
		StopChannel:     stopChannel,
		IsInactiveRacer: make([]bool, len(infoChannels)),
		MutexRacersInfo: sync.RWMutex{},
	}
}

func (j *JudgeOfRace) StartRace() {
	go j.runStepTicker()
	go j.runDisplayTicker()
	go j.runRacersInfoCollect()
	go j.startToJudge()
}

func (j *JudgeOfRace) runRacersInfoCollect() {
	var ok bool
	var in domain.RacerInfo
	for {
		for i := range j.InfoChannels {
			if in, ok = <-j.InfoChannels[i]; ok && !j.IsInactiveRacer[i] {
				j.MutexRacersInfo.Lock()
				j.RacersInfo[i] = in
				j.MutexRacersInfo.Unlock()
			}
		}
		time.Sleep(domain.LoopSleepTime)
	}
}

func (j *JudgeOfRace) runStepTicker() {
	j.StepTicker = time.NewTicker(domain.StepTime)
	var s time.Time

	for {
		select {
		case s = <-j.StepTicker.C:
			for i := range j.StepChannel {
				j.StepChannel[i] <- s
			}
		default:
			continue
		}
		time.Sleep(domain.LoopSleepTime)
	}
}

func (j *JudgeOfRace) runDisplayTicker() {
	j.DisplayTicker = time.NewTicker(domain.DisplayTime)

	for {
		select {
		case <-j.DisplayTicker.C:
			j.MutexRacersInfo.RLock()
			j.DisplayChannel <- j.RacersInfo
			j.MutexRacersInfo.RUnlock()
		default:
			continue
		}
		time.Sleep(domain.LoopSleepTime)
	}
}

func (j *JudgeOfRace) startToJudge() {
	sortedInfo := make([]domain.RacerInfo, len(j.RacersInfo))
	for {
		time.Sleep(domain.LoopSleepTime)
		j.MutexRacersInfo.RLock()

		copy(sortedInfo, j.RacersInfo)

		j.MutexRacersInfo.RUnlock()
		sort.SliceStable(sortedInfo, func(i, j int) bool {
			return sortedInfo[i].Lap > sortedInfo[j].Lap
		})
		var nameOfRacerToStop string
		if sortedInfo[len(sortedInfo) - 1 - j.InactiveCount].Lap < sortedInfo[len(sortedInfo) - 2 - j.InactiveCount].Lap {
			nameOfRacerToStop = sortedInfo[len(sortedInfo) - 1 - j.InactiveCount].Name
			racerIndex := j.findRacerIndexByName(nameOfRacerToStop)
				j.IsInactiveRacer[racerIndex] = true
				j.InactiveCount++
		}
		if j.InactiveCount == len(j.InfoChannels)-1 {
			j.MutexRacersInfo.RLock()
			for i := range j.RacersInfo {
				if !j.IsInactiveRacer[i] {
					fmt.Println("\nThe winner is " + j.RacersInfo[i].Name)
				}
			}
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
