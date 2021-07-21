package main

import (
	"RacersRace/app"
	"RacersRace/domain"
	"RacersRace/infrastructure"
	"strconv"
	"time"
)

func main() {
	racers := make([]*app.RandomStepRacer, domain.Racers,)
	stepChannel := make([]chan time.Time, domain.Racers)
	infoChannels := make([]chan domain.RacerInfo, domain.Racers)
	displayChannel := make(chan []domain.RacerInfo)
	stopChannel := make (chan bool)
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
	judge := app.NewRaceJudge(stepChannel, infoChannels, displayChannel, stopChannel)
	display := infrastructure.NewConsole(displayChannel)
	go display.StartShowRaceSatus()
	go judge.StartRace()
	for {
		<- stopChannel
		return
	}
}
