package app

import (
	"RacersRace/domain"
	"fmt"
	"sort"
	"time"
)

type JudgeOfRace struct {
	RacersInfo     []domain.RacerInfo
	StepTicker     *time.Ticker
	StepChannel    []chan time.Time
	DisplayTicker  *time.Ticker
	DisplayChannel chan []domain.RacerInfo
	InfoChannels   []chan domain.RacerInfo
	// trial
	InactiveRacers []bool
	InactiveCount int
}

func NewRaceJudge(stepChannel []chan time.Time, infoChannels []chan domain.RacerInfo, displayChannel chan []domain.RacerInfo) *JudgeOfRace {
	return &JudgeOfRace{
		RacersInfo: make([]domain.RacerInfo, len(infoChannels)),
		StepChannel:    stepChannel,
		DisplayChannel: displayChannel,
		InfoChannels:   infoChannels,
		InactiveRacers: make([]bool, len(infoChannels)),
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
			if in, ok = <-j.InfoChannels[i]; ok && !j.InactiveRacers[i] {
				j.RacersInfo[i] = in
			}
		}
		time.Sleep(200*time.Millisecond)
	}
}

func (j *JudgeOfRace) runStepTicker() {
	j.StepTicker = time.NewTicker(time.Second)
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
		time.Sleep(200*time.Millisecond)
	}
}

func (j *JudgeOfRace) runDisplayTicker() {
	j.DisplayTicker = time.NewTicker(3 * time.Second)

	for {
		select {
		case <-j.DisplayTicker.C:
			j.DisplayChannel <- j.RacersInfo
		default:
			continue
		}
		time.Sleep(200*time.Millisecond)
	}
}

func (j *JudgeOfRace) startToJudge() {
	for {
		time.Sleep(500*time.Millisecond)
		sortedInfo := j.RacersInfo
		sort.SliceStable(sortedInfo, func(i, j int) bool {
			return sortedInfo[i].Lap > sortedInfo[j].Lap
		})
		var nameOfRacerToStop string
		if sortedInfo[len(sortedInfo) - 1 - j.InactiveCount].Lap < sortedInfo[len(sortedInfo) - 2 - j.InactiveCount].Lap {
			nameOfRacerToStop = sortedInfo[len(sortedInfo) - 1].Name
			racerIndex := j.findRacerIndexByName(nameOfRacerToStop)
			if _, ok := <- j.InfoChannels[racerIndex]; ok {
				// trial
				j.InactiveRacers[racerIndex] = true
				j.InactiveCount++
			}
		}


		//fmt.Println(sortedInfo)
	}
}

func (j *JudgeOfRace) findRacerIndexByName (name string) int {
	for i := range j.RacersInfo {
		if j.RacersInfo[i].Name == name {
			return i
		}
	}
	panic(fmt.Errorf("racer not found"))
}
