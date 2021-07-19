package infrastructure

import "RacersRace/domain"

type Console struct {

}

func (c Console) ShowRaceSatus(racers []domain.Racer) {
	for r := range racers {

	}
}

func NewConsole() *Console {
	return &Console{}
}