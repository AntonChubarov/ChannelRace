package app

import (
	"RacersRace/domain"
	"time"
)

type JudgeOfRace struct {
	RacersInfo *[]domain.RacerInfo
	StepTicker *time.Ticker
	StepChannel chan time.Time
	DisplayTicker *time.Ticker
	DisplayChannel chan []domain.RacerInfo
}

func NewRaceJudge(stepChannel chan time.Time, infoChannels []chan domain.RacerInfo) *JudgeOfRace {
	return &JudgeOfRace{
		RacersInfo: runRacersInfoCollect(infoChannels),
		StepChannel: stepChannel,
	}
}

func (j *JudgeOfRace) StartRace() {
	go runStepTicker(j.StepChannel)
	go runDisplayTicker(j.DisplayChannel, j.RacersInfo)
}

func runRacersInfoCollect(infoChannels []chan domain.RacerInfo) *[]domain.RacerInfo {

	info := make([]domain.RacerInfo, len(infoChannels), len(infoChannels))
	go func(*[]domain.RacerInfo) {
		var ok bool
		var in domain.RacerInfo
		for {
			for i := range infoChannels {
				if in, ok = <-infoChannels[i]; ok {
					info[i] = in
				}
			}
		}
	}(&info)

	return &info
}

func runStepTicker(stepChannel chan time.Time) {
	stepTicker := time.NewTicker(time.Second)
	var s time.Time
	Loop:
	for {
		s = <- stepTicker.C
		select{
		case stepChannel <- s:
			continue
		default:
			break Loop
		}
	}
	return
}

func runDisplayTicker(displayChannel chan []domain.RacerInfo, racersInfo *[]domain.RacerInfo) {
	displayTicker := time.NewTicker(3 * time.Second)
Loop:
	for {
		select{
		case <- displayTicker.C:
			displayChannel <- *racersInfo
		default:
			break Loop
		}
	}
	return
}

func ToJudge() {

}