package app

import (
	"RacersRace/domain"
	"sync"
	"time"
)

type DisplayCommander struct {
	racersInfo      *[]domain.RacerInfo
	displayTicker   *time.Ticker
	displayChannel  chan []domain.RacerInfo
	mutexRacersInfo *sync.RWMutex
}

func NewDisplayCommander(racersInfo *[]domain.RacerInfo,
	displayChannel chan []domain.RacerInfo,
	mutexRacersInfo *sync.RWMutex,
) *DisplayCommander {
	return &DisplayCommander{racersInfo: racersInfo,
		displayChannel: displayChannel,
		mutexRacersInfo: mutexRacersInfo,
	}
}

func (d *DisplayCommander) Run() {
	d.displayTicker = time.NewTicker(domain.DisplayTime)

	for {
		select {
		case <-d.displayTicker.C:
			d.mutexRacersInfo.RLock()
			d.displayChannel <- *d.racersInfo
			d.mutexRacersInfo.RUnlock()
		default:
			continue
		}
		time.Sleep(domain.LoopSleepTime)
	}
}
