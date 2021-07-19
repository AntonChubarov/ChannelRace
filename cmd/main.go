package main

import (
	"RacersRace/app"
	"RacersRace/domain"
	"strconv"
	"time"
)

func main() {
	numOfRacers := 5;

	racers := make([]*app.RandomStepRacer, numOfRacers, numOfRacers)

	stepChannel := make(chan time.Time)

	infoChannels := make([]chan domain.RacerInfo, numOfRacers, numOfRacers)

	for i := range racers {
		racers[i] = app.NewRandomStepRacer("Racer " + strconv.Itoa(i), stepChannel, infoChannels[i])
		go racers[i].StartRace()
	}

	judge := app.NewRaceJudge(stepChannel, infoChannels)

	judge.StartRace()

	Loop:
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-stepChannel:
			continue
		default:
			break Loop
		}
	}

	return
}
