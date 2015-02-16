package main

import (
	"strings"
)

type Player struct {
	game      *Game
	hand      []Card
	knowledge [][]Card
	plan      []struct {
		a   Action
		pos int
	}
	startTurn  chan bool
	turnAction chan *Action
}

func (p *Player) Draw(c Card) {
	p.hand = append(p.hand, c)
	p.knowledge = append(p.knowledge, make([]Card, 0))
}

func (p *Player) Play() {
	for _ = range p.startTurn {

		p.turnAction <- new(Action)
	}
}

func (p *Player) String() string {
	r := make([]string, 0, len(p.hand))
	for _, c := range p.hand {
		r = append(r, c.String())
	}
	return strings.Join(r, " ")
}
