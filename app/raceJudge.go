package app

import "RacersRace/domain"

type JudgeOfRace struct {
	Racers []domain.Racer
}

func NewRaceJudge(racers []domain.Racer) *JudgeOfRace {
	return &JudgeOfRace{Racers: racers}
}

func (j *JudgeOfRace) StartRace() {

}