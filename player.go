package main

import (
	"strings"
)

type Player struct {
	game       *Game
	position   int
	hand       []*Card
	knowledge  [][]Characteristic
	plan       []*Action
	toPlay     []*Card
	tableTalk  chan *Information
	turnAction chan *Action
}

func (p *Player) Draw() {
	c, err := p.game.Draw()
	if err != nil {
		// DO SOMETHING because the game's ending
	} else {
		p.hand = append(p.hand, c)
		p.knowledge = append(p.knowledge, make([]Characteristic, 0, 2))
	}
}

func (p *Player) GetDiagonal(from int) []*Card {
	d := make([]*Card, 0, 4)
	g := p.game
	for i := 1; i < 5; i++ {
		if (i+from)%5 == p.position {
			d = append(d, nil)
		}
		d = append(d, g.players[(i+from)%5].hand[i-1])
	}
	return d
}

func (p *Player) GetPlayDiscard(diag []*Card) (play, discard []int) {
	// list of playable/discardable positions
	play = make([]int, 0, 4)
	discard = make([]int, 0, 4)
	// imagine what the state will be when all called cards have been played
	pseudoPiles := make(map[Color]Number)
	for c, n := range p.game.piles {
		pseudoPiles[c] = n
	}
	// "play" all the cards in the play queue
	var toPlay []*Card
	copy(toPlay, p.toPlay)
	for len(toPlay) > 0 {
		for i := 0; i < len(toPlay); i++ {
			if pseudoPiles[toPlay[i].color] == toPlay[i].number-1 {
				pseudoPiles[toPlay[i].color]++
				toPlay = append(toPlay[:i], toPlay[i+1:]...)
				i--
			}
		}
	}

	// count distinct plays and all discards
	distinctPlays := make(map[Card]bool)
	for i, c := range diag {
		if c == nil {
			continue
		}
		if pseudoPiles[c.color] >= c.number {
			discard = append(discard, i)
		} else if pseudoPiles[c.color] == c.number-1 && !distinctPlays[*c] {
			play = append(play, i)
			distinctPlays[*c] = true
		}
	}
	return
}

func (p *Player) Play() {
	g := p.game
	for info := range p.tableTalk {
		switch info.characteristic.(type) {
		case int:
			if info.from == p.position {
				continue
			}
			// message: playable cards
			if info.to == p.position {
				// TODO: implement knowledge
			}
			relPos := (info.to - info.from + 5) % 5
			d := p.GetDiagonal(info.from)
			playable, _ := p.GetPlayDiscard(d)
			// message received: playable card
			if len(playable) == relPos-1 {
				cardPos := (p.position-info.from+5)%5 - 1
				p.plan = append(p.plan, &Action{Play, nil, cardPos, nil})
			}
			// queue other players' playable cards
			for _, i := range playable {
				p.toPlay = append(p.toPlay, d[i])
			}
		case string:
			if info.from == p.position {
				continue
			}
			// message: discardable cards
			if info.to == p.position {
				// TODO: implement knowledge
			}
			relPos := (info.to - info.from + 5) % 5
			d := p.GetDiagonal(info.from)
			_, discardable := p.GetPlayDiscard(d)
			// message received: discardable card
			if len(discardable) == relPos-2 {
				cardPos := (p.position-info.from+5)%5 - 1
				p.plan = append(p.plan, &Action{Discard, nil, cardPos, nil})
			}
		case bool:
			a := new(Action)
			if len(p.plan) > 0 {
				// execute the plan
				action := p.plan[0]
				p.plan = p.plan[1:]
				a.t = action.t
				a.c = p.hand[action.p]
				p.hand = append(p.hand[:action.p], p.hand[action.p+1:]...)
				p.knowledge = append(p.knowledge[:action.p],
					p.knowledge[action.p+1:]...)
				for i := range p.plan {
					if p.plan[i].p > action.p {
						p.plan[i].p--
					}
				}
				p.Draw()
			} else if g.clocks == 0 {
				// nothing else to do
				a.t = Discard
				a.c = p.hand[0]
				p.hand = p.hand[1:]
				p.knowledge = p.knowledge[1:]
				p.Draw()
			} else {
				a.t = Inform
				a.i = &Information{p.position, 0, "", make([]int, 0)}
				diag := p.GetDiagonal(p.position)
				playable, discardable := p.GetPlayDiscard(diag)
				// tell players about discardable cards
				if g.clocks < 3 && len(discardable) > 0 ||
					len(discardable) > len(playable)+2 || len(playable) == 0 {
					if len(discardable) == 4 {
						discardable = discardable[:len(discardable)-1]
					}
					a.i.to = (p.position + len(discardable) + 1) % 5
					p := g.players[a.i.to]
					a.i.characteristic = p.hand[0].color
					for c := range p.hand {
						if p.hand[c].color == a.i.characteristic {
							a.i.positions = append(a.i.positions, c)
						}
					}
				} else {
					// tell players about playable cards
					a.i.to = (p.position + len(playable)) % 5
					p := g.players[a.i.to]
					a.i.characteristic = p.hand[0].number
					for c := range p.hand {
						if p.hand[c].number == a.i.characteristic {
							a.i.positions = append(a.i.positions, c)
						}
					}
				}
			}
			p.turnAction <- a
		}
	}
}

func (p *Player) String() string {
	r := make([]string, 0, len(p.hand))
	for _, c := range p.hand {
		r = append(r, c.String())
	}
	for _, p := range p.plan {
		r = append(r, p.String())
	}
	return strings.Join(r, " ")
}
