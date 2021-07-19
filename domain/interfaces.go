package domain

type Racer interface {
	StartRace()
	StopRace()
}

type Visualiser interface {
	ShowRaceSatus([]Racer)
}