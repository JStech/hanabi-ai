package main

import (
	"errors"
	"math/rand"
)

type Game struct {
	deck    []Card
	topCard int
	clocks  int
	fuses   int
	piles   map[string]int
	turn    []chan bool
	players []*Player
	result  chan bool
}

type Card int

func (c Card) Color() string {
	return "rgybw"[c/10 : c/10+1]
}

func (c Card) Number() int {
	return []int{1, 1, 1, 2, 2, 3, 3, 4, 4, 5}[c%10]
}

type Action int

const (
	Play Action = iota
	Discard
)

func NewGame(seed int64) *Game {
	// create game
	g := new(Game)
	// initialize and shuffle the deck
	rand.Seed(seed)
	g.deck = make([]Card, 50)
	for i := 0; i < 50; i++ {
		g.deck[i] = Card(i)
	}
	for i := 49; i > 0; i-- {
		r := rand.Int31n(int32(i + 1))
		g.deck[i], g.deck[r] = g.deck[r], g.deck[i]
	}
	g.topCard = 0

	// initalize counters
	g.clocks = 8
	g.fuses = 3

	// initialize players
	g.players = make([]*Player, 5)
	g.turn = make([]chan bool, 5)
	for i := range g.players {
		g.players[i] = &Player{
			g,
			make([]Card, 0, 4),
			make([][]Card, 0, 4),
			make([]struct {
				a   Action
				pos int
			}, 0, 4),
			g.turn[i],
			g.turn[(i+1)%5],
		}
	}

	// result channel
	g.result = make(chan bool)

	// deal cards
	for i := 0; i < 4; i++ {
		for _, p := range g.players {
			p.Draw(g.deck[g.topCard])
			g.topCard++
		}
	}
	return g
}

func (g *Game) Draw() (c Card, err error) {
	if g.topCard >= 50 {
		err = errors.New("Out of cards")
		return
	}
	c = g.deck[g.topCard]
	g.topCard++
	return
}

func (g *Game) Play(c Card) {
	if g.piles[c.Color()] == c.Number()-1 {
		g.piles[c.Color()] = c.Number()
		if c.Number() == 5 {
			g.clocks++
		}
	} else {
		g.fuses--
		if g.fuses == 0 {
			g.result <- false
		}
	}
}
