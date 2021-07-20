package app

import (
	"RacersRace/domain"
	"time"
)

type JudgeOfRace struct {
	RacersInfo     []domain.RacerInfo
	StepTicker     *time.Ticker
	StepChannel    []chan time.Time
	DisplayTicker  *time.Ticker
	DisplayChannel chan []domain.RacerInfo
	InfoChannels   []chan domain.RacerInfo
}

func NewRaceJudge(stepChannel []chan time.Time, infoChannels []chan domain.RacerInfo, displayChannel chan []domain.RacerInfo) *JudgeOfRace {
	return &JudgeOfRace{
		RacersInfo: make([]domain.RacerInfo, len(infoChannels)),
		StepChannel:    stepChannel,
		DisplayChannel: displayChannel,
		InfoChannels:   infoChannels,
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
			if in, ok = <-j.InfoChannels[i]; ok {
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
		time.Sleep(200*time.Millisecond)
	}
}
