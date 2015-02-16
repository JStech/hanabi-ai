package main

type Card struct {
	color  Color
	number Number
}

type Characteristic interface {
	String() string
	Equals(other Characteristic) bool
	Orthogonal(other Characteristic) bool
	And(other Characteristic) *Card
}

func NewCard(n int) *Card {
	c := Card{
		Color("rgybw"[n/10 : n/10+1]),
		Number([]int{1, 1, 1, 2, 2, 3, 3, 4, 4, 5}[n%10]),
	}
	return &c
}

type Color string

func (c Color) String() string {
	return string(c)
}

func (c Color) Equals(other Characteristic) bool {
	switch other.(type) {
	case Color:
		return c == other
	}
	return false
}

func (c Color) Orthogonal(other Characteristic) bool {
	switch other.(type) {
	case Number:
		return true
	}
	return false
}

func (c Color) And(other Characteristic) *Card {
	return &Card{c, other.(Number)}
}

type Number int

func (n Number) String() string {
	return "12345"[n-1 : n]
}

func (n Number) Equals(other Characteristic) bool {
	switch other.(type) {
	case Number:
		return n == other
	}
	return false
}

func (n Number) Orthogonal(other Characteristic) bool {
	switch other.(type) {
	case Color:
		return true
	}
	return false
}

func (n Number) And(other Characteristic) *Card {
	return &Card{other.(Color), n}
}

func (c Card) String() string {
	return c.color.String() + c.number.String()
}
