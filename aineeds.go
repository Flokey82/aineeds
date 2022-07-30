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

const (
	DefaultHP                = 5 // Default starting HP
	HungerDeathThreshold     = 5 // Being dies if hunger exceeds this value
	HungerSatiationThreshold = 2 // The satiation need is triggered if hunger exceeds this threshold
)

// Being represents an entity capable of performing actions.
type Being struct {
	name   string
	MaxHP  int
	HP     int
	Hunger int
	Dead   bool
}

// NewBeing returns a new being with the given name.
func NewBeing(name string, maxHP int) *Being {
	return &Being{
		name:  name,
		MaxHP: maxHP,
		HP:    maxHP,
	}
}

// Name returns the name of the being.
func (b *Being) Name() string {
	return b.name
}

// Act is called on each cycle.
func (b *Being) Act() {
	b.Log("acts")
	b.Hunger++
	if b.Hunger > HungerDeathThreshold {
		b.Die("starved")
	}
}

// Eat resets the hunger value back to 0.
func (b *Being) Eat() {
	b.Hunger = 0
	b.Log("eats")
}

// Damage the being from an attacker.
func (b *Being) Damage(attacker *Being) {
	b.HP--
	attacker.Log("attacks")
	b.Log("takes damage")
	if b.HP <= 0 {
		b.Die("damage")
	}
}

// Heal heals the being up to max health.
func (b *Being) Heal() {
	b.HP = b.MaxHP
	b.Log("heals")
}

// Die indicates that the being has died because of 'reason'.
func (b *Being) Die(reason string) {
	b.Dead = true
	b.Log(fmt.Sprintf("dies (%s)", reason))
}

// Log logs an event encountered by the being.
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

// AI implements a very basic needs based AI.
type AI struct {
	Being     *Being        // underlying being
	Enemy     *Being        // current enemy
	Needs     [NeedMax]bool // determines which needs need to be met
	Prioities []Need        // determines the order of needs
}

// NewAI returns a new AI controlled being.
func NewAI(being *Being) *AI {
	return &AI{
		Being: being,
		Prioities: []Need{
			NeedSurvival,
			NeedHealth,
			NeedSatiation,
			NeedRevenge,
		},
	}
}

// Name returns the name of the underlying being.
func (a *AI) Name() string {
	return a.Being.Name()
}

// Act is called on every cycle.
func (a *AI) Act() {
	if a.Being.Dead {
		return // Dead beings don't act.
	}
	defer a.Being.Act()

	// Observe values changed during the being.act.
	// TODO: Register the observation of values -> needs
	// instead of hardcoding.
	if a.Being.HP < a.Being.MaxHP {
		a.Needs[NeedHealth] = true
	}
	if a.Being.HP < 3 {
		a.Needs[NeedSurvival] = true
	}
	if a.Being.Hunger > HungerSatiationThreshold {
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

// ActOnNeed causes the AI to act on the given need.
func (a *AI) ActOnNeed(need Need) {
	// TODO: Register needs somehow instead of hardcoding
	// needs -> actions.
	switch need {
	case NeedSurvival, NeedHealth:
		a.Heal()
	case NeedSatiation:
		a.Eat()
	case NeedRevenge:
		a.Damage(a.Enemy)
	}
}

// Eat is called when the AI acts on 'NeedSatiation'.
func (a *AI) Eat() {
	a.Needs[NeedSatiation] = false
	a.Being.Eat()
}

// Heal is called when the AI acts on 'NeedHealth', 'NeedSurvival'.
func (a *AI) Heal() {
	a.Needs[NeedHealth] = false
	a.Needs[NeedSurvival] = false
	a.Being.Heal()
}

// Damage is called when the AI is damaged by an attacker.
func (a *AI) Damage(attacker *Being) {
	a.Enemy = attacker
	a.Needs[NeedRevenge] = true // TODO: Find a better way to set this.
	a.Being.Damage(attacker)
}
