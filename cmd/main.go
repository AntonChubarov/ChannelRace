package main

import (
	"RacersRace/app"
	"RacersRace/domain"
	"RacersRace/infrastructure"
	"strconv"
	"time"
)

func main() {
	numOfRacers := 5;

	racers := make([]*app.RandomStepRacer, numOfRacers, numOfRacers)

	stepChannel := make([]chan time.Time, numOfRacers, numOfRacers)

	infoChannels := make([]chan domain.RacerInfo, numOfRacers, numOfRacers)

	displayChannel := make(chan []domain.RacerInfo)

	for i := range racers {
		racerNumber := strconv.Itoa(i)
		if i < 10 {
			racerNumber = "0" + racerNumber
		}
		stepChannel[i] = make(chan time.Time)
		infoChannels[i] = make(chan domain.RacerInfo)
		racers[i] = app.NewRandomStepRacer("Racer" + racerNumber, stepChannel[i], infoChannels[i])
		go racers[i].StartRace()
	}

	judge := app.NewRaceJudge(stepChannel, infoChannels, displayChannel)

	display := infrastructure.NewConsole(displayChannel)

	go display.StartShowRaceSatus()

	go judge.StartRace()

	for {
		time.Sleep(time.Second)
	}
}
