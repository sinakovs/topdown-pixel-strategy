package world

type UnitID string

type unit struct {
	ID       UnitID
	OwnerID  PlayerID
	WorldCol uint32
	WorldRow uint32
}
