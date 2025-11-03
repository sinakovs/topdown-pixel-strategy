package world

type PlayerID string

type player struct {
	ID       PlayerID
	UnitsIDs map[UnitID]*unit
}
