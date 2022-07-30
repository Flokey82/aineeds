// Package aineeds is a port of https://github.com/ondras/need-based-ai/.
// NOTE: this is just a very basic implementation and needs to be refined.
package aineeds

import (
	"fmt"
	"log"
)

type Entity interface {
	Name() string
	Act()
	Eat()
	Heal()
	Damage(*Being)
}

type Being struct {
	name   string
	HP     int
	Hunger int
	Dead   bool
}

func NewBeing(name string) *Being {
	return &Being{
		name: name,
		HP:   10,
	}
}

func (b *Being) Name() string {
	return b.name
}

func (b *Being) Act() {
	b.Log("acts")
	b.Hunger++
	if b.Hunger > 5 {
		b.Die("starved")
	}
}

func (b *Being) Eat() {
	b.Hunger = 0
	b.Log("eats")
}

func (b *Being) Damage(attacker *Being) {
	b.HP--
	attacker.Log("attacks")
	b.Log("takes damage")
	if b.HP <= 0 {
		b.Die("damage")
	}
}

func (b *Being) Heal() {
	b.HP = 5
	b.Log("heals")
}

func (b *Being) Die(reason string) {
	b.Dead = true
	b.Log(fmt.Sprintf("dies (%s)", reason))
}

func (b *Being) Log(text string) {
	log.Printf("%s: %s", b.Name(), text)
}

type Need int

const (
	NeedSurvival Need = iota
	NeedHealth
	NeedSatiation
	NeedRevenge
	NeedMax
)

type AI struct {
	Being     *Being
	Enemy     *Being
	Needs     [NeedMax]bool
	Prioities []Need
}

func NewAI(being *Being) *AI {
	return &AI{
		Being: being,
		Needs: [NeedMax]bool{
			NeedSurvival:  false,
			NeedHealth:    false,
			NeedSatiation: false,
			NeedRevenge:   false,
		},
		Prioities: []Need{
			NeedSurvival,
			NeedHealth,
			NeedSatiation,
			NeedRevenge,
		},
	}
}

func (a *AI) Name() string {
	return a.Being.Name()
}

func (a *AI) Act() {
	if a.Being.Dead {
		return // Dead beings don't act.
	}
	defer a.Being.Act()

	// Observe values changed during the being.act.
	if a.Being.Hunger > 2 {
		a.Needs[NeedSatiation] = true
	}
	if a.Enemy != nil && a.Enemy.Dead {
		a.Enemy = nil
		a.Needs[NeedRevenge] = false
	}

	// Act on the needs.
	for _, need := range a.Prioities {
		if a.Needs[need] {
			a.Being.Log(fmt.Sprintf("needs %d", need))
			a.ActOnNeed(need)
			return
		}
	}
	a.Being.Log("idle")
}

func (a *AI) ActOnNeed(need Need) {
	switch need {
	case NeedSurvival, NeedHealth:
		a.Heal()
	case NeedSatiation:
		a.Eat()
	case NeedRevenge:
		a.Damage(a.Enemy)
	}
}

func (a *AI) Eat() {
	a.Needs[NeedSatiation] = false
	a.Being.Eat()
}

func (a *AI) Heal() {
	a.Needs[NeedHealth] = false
	a.Needs[NeedSurvival] = false
	a.Being.Heal()
}

func (a *AI) Damage(attacker *Being) {
	a.Enemy = attacker
	a.Needs[NeedRevenge] = true
	a.Needs[NeedHealth] = true
	if a.Being.HP < 3 {
		a.Needs[NeedSurvival] = true
	}
	a.Being.Damage(attacker)
}
