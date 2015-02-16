package main

type Player struct {
	game      *Game
	hand      []Card
	knowledge [][]Card
	plan      []struct {
		a   Action
		pos int
	}
	turnIn  chan bool
	turnOut chan bool
}

func (p *Player) Draw(c Card) {
	p.hand = append(p.hand, c)
	p.knowledge = append(p.knowledge, make([]Card, 0))
}
