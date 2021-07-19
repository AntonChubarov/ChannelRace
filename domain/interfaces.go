package domain

type Racer interface {
	StartRace()
	GiveInfo()
	StopRace()
}

type Visualiser interface {
	ShowRaceSatus([]Racer)
}