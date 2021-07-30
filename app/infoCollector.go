package app

import (
	"RacersRace/domain"
	"sync"
	"time"
)

type InfoCollector struct {
	racersInfo      *[]domain.RacerInfo
	infoChannels    *[]chan domain.RacerInfo
	isInactiveRacer *[]bool
	mutexRacersInfo *sync.RWMutex
}

func NewInfoCollector(racersInfo *[]domain.RacerInfo,
	infoChannels *[]chan domain.RacerInfo,
	isInactiveRacer *[]bool,
	mutexRacersInfo *sync.RWMutex,
) *InfoCollector {
	return &InfoCollector{racersInfo: racersInfo,
		infoChannels: infoChannels,
		isInactiveRacer: isInactiveRacer,
		mutexRacersInfo: mutexRacersInfo,
	}
}


func (ic *InfoCollector) Run() {
	var ok bool
	var in domain.RacerInfo
	for {
		for i := range *ic.infoChannels {
			if in, ok = <- (*ic.infoChannels)[i]; ok && !(*ic.isInactiveRacer)[i] {
				ic.mutexRacersInfo.Lock()
				(*ic.racersInfo)[i] = in
				ic.mutexRacersInfo.Unlock()
			}
		}
		time.Sleep(domain.LoopSleepTime)
	}
}