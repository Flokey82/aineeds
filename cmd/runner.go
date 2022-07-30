package main

import (
	"github.com/Flokey82/aineeds"
)

func main() {
	var entities []aineeds.Entity

	entities = append(entities, aineeds.NewAI(aineeds.NewBeing("AI 1")))
	entities = append(entities, aineeds.NewAI(aineeds.NewBeing("AI 2")))

	for {
		for _, e := range entities {
			e.Act()
		}
	}
}
