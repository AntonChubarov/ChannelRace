package app

import (
	"math/rand"
	"time"
)

type RandomStepRacer struct {
	Name string
	judgeChan chan bool
	infoChan
	Step int
	Score int
	Lap int
	IsActive bool
}

func NewRandomStepRacer(name string, judgeChan chan bool) *RandomStepRacer {
	return &RandomStepRacer{
		Name: name,
		judgeChan: judgeChan,
		Step: 0,
		Score: 0,
		Lap:1,
		IsActive: true,
	}
}

func (r RandomStepRacer) StartRace() {
	for r.IsActive {
		<- r.judgeChan
		r.makeStep()
	}
}

func (r RandomStepRacer) GiveInfo() {
	panic("implement me")
}

func (r RandomStepRacer) StopRace() {
	panic("implement me")
}

func (r RandomStepRacer) makeStep() {
	rand.Seed(time.Now().UnixNano())
	points := 1 + rand.Intn(4)
	r.Step++
	r.Score += points
	r.Lap = 1 + r.Score / 50
}