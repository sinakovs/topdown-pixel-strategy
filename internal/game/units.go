package game

func (g *Game) moveUnitTowardsTarget() {
	if g.unit.x < g.target.x {
		g.unit.x++
	} else if g.unit.x > g.target.x {
		g.unit.x--
	}

	if g.unit.y < g.target.y {
		g.unit.y++
	} else if g.unit.y > g.target.y {
		g.unit.y--
	}
}
