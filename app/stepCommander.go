package app

import (
	"RacersRace/domain"
	"time"
)

type StepCommander struct {
	stepChannel []chan time.Time
	stepTicker  *time.Ticker
}

func NewStepCommander(stepChannel []chan time.Time) *StepCommander {
	return &StepCommander{stepChannel: stepChannel}
}

func (s *StepCommander) Run() {
	s.stepTicker = time.NewTicker(domain.StepTime)
	var signal time.Time

	for {
		select {
		case signal = <-s.stepTicker.C:
			for i := range s.stepChannel {
				s.stepChannel[i] <- signal
			}
		default:
			continue
		}
		time.Sleep(domain.LoopSleepTime)
	}
}