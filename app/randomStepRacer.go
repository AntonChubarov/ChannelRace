package app

import (
	"RacersRace/domain"
	"math/rand"
	"time"
)

type RandomStepRacer struct {
	Name     string
	stepChan chan time.Time
	infoChan chan domain.RacerInfo
	Step     int
	Score    int
	Lap      int
}

func NewRandomStepRacer(name string, sChan chan time.Time, iChan chan domain.RacerInfo) *RandomStepRacer {
	return &RandomStepRacer{
		Name:     name,
		stepChan: sChan,
		infoChan: iChan,
		Step:     0,
		Score:    0,
		Lap:      1,
	}
}

func (r *RandomStepRacer) StartRace() {
	for {
		<- r.stepChan
		r.makeStep()
		r.infoChan <- domain.RacerInfo{
			Name: r.Name,
			Step: r.Step,
			Score: r.Score,
			Lap: r.Lap,
		}
		time.Sleep(domain.LoopSleepTime)
	}
}

func (r *RandomStepRacer) makeStep() {
	rand.Seed(time.Now().UnixNano())
	points := 1 + rand.Intn(6)
	r.Step++
	r.Score += points
	r.Lap = 1 + r.Score / 50
}