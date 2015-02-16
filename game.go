package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type Game struct {
	deck       []Card
	topCard    int
	clocks     int
	fuses      int
	piles      map[string]int
	turn       int
	tableTalk  []chan *Information
	turnAction chan *Action
	players    []*Player
}

type Card int

func (c Card) Color() string {
	return "rgybw"[c/10 : c/10+1]
}

func (c Card) Number() int {
	return []int{1, 1, 1, 2, 2, 3, 3, 4, 4, 5}[c%10]
}

func (c Card) String() string {
	return fmt.Sprintf("%d%s", c.Number(), c.Color())
}

type ActionType int

const (
	Play ActionType = iota
	Discard
	Inform
)

type Action struct {
	t ActionType
	c Card
	p int
	i *Information
}

func (a *Action) String() string {
	r := ""
	switch a.t {
	case Play:
		r += "Play " + a.c.String()
		r += fmt.Sprintf(" %d", a.p)
	case Discard:
		r += "Discard " + a.c.String()
		r += fmt.Sprintf(" %d", a.p)
	case Inform:
		r += "Tell " + a.i.String()
	}
	return r
}

type Information struct {
	from           int
	to             int
	characteristic interface{}
	positions      []int
}

func (info *Information) String() string {
	r := fmt.Sprintf("Player %d, %v at", info.to,
		info.characteristic)
	for _, p := range info.positions {
		r += fmt.Sprintf(" %d,", p)
	}
	return r
}

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

	// initialize piles
	g.piles = make(map[string]int)

	// initialize players
	g.players = make([]*Player, 5)
	g.tableTalk = make([]chan *Information, 5)
	g.turnAction = make(chan *Action)
	for i := range g.players {
		g.tableTalk[i] = make(chan *Information)
		g.players[i] = &Player{
			g,
			i,
			make([]Card, 0, 4),
			make([][]Card, 0, 4),
			make([]Action, 0, 4),
			g.tableTalk[i],
			g.turnAction,
		}
	}

	// deal cards
	for i := 0; i < 4; i++ {
		for _, p := range g.players {
			p.Draw()
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

func (g *Game) Play(c Card) bool {
	// check if it's playable
	if g.piles[c.Color()] == c.Number()-1 {
		g.piles[c.Color()] = c.Number()
		// bonus clock for playing a 5
		if c.Number() == 5 {
			g.clocks++
		}
		all5 := true
		for _, n := range g.piles {
			all5 = all5 && (n == 5)
		}
		// you win!
		if all5 {
			return true
		}
	} else {
		g.fuses--
		// game over!
		if g.fuses == 0 {
			return true
		}
	}
	return false
}

func (g *Game) Discard(c Card) {
	g.clocks++
}

func (g *Game) String() string {
	r := ""
	score := 0
	for c, n := range g.piles {
		r += fmt.Sprintf("%d%s ", n, c)
		score += n
	}
	r += fmt.Sprintf(" score: %d  clocks: %d  fuses: %d\n", score, g.clocks,
		g.fuses)
	for i := range g.deck {
		if i == g.topCard {
			r += "|"
		} else {
			r += " "
		}
		r += g.deck[i].String()
		if i%10 == 9 {
			r += "\n"
		}
	}

	for i, p := range g.players {
		if i == g.turn {
			r += "> "
		} else {
			r += "  "
		}
		r += fmt.Sprintf("Player %d, %s\n", i, p.String())

	}
	return r
}

func (g *Game) Start() {
	for _, p := range g.players {
		go p.Play()
	}
	for g.topCard < 50 {
		g.tableTalk[g.turn] <- &Information{0, 0, true, nil}
		a := <-g.turnAction
		fmt.Println(g.turn, a)
		fmt.Println(g)
		switch a.t {
		case Play:
			over := g.Play(a.c)
			if over {
				fmt.Println("Game over")
				fmt.Println(g)
				return
			}
		case Discard:
			g.Discard(a.c)
		case Inform:
			g.clocks--
			for _, c := range g.tableTalk {
				c <- a.i
			}
		}
		g.turn++
		g.turn %= 5
	}
}
