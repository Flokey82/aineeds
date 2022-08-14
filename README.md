# aineeds: Needs based AI

This package is an implementation of http://www.roguebasin.com/index.php/Need-based_AI_in_JavaScript and a port of https://github.com/ondras/need-based-ai/ to Go.

Lifted straight from roguebasin:

The AI logic is based on the famous [Hierarchy of needs](https://en.wikipedia.org/wiki/Maslow%27s_hierarchy_of_needs) and prioritizes behavior accordingly.

Survival is the basic need to survive, i.e. to maintain at least a minimal amount of hitpoints. Health is a need to be healthy, to regenerate as many hitpoints as possible. Satiation represents the need to fight hunger. Finally, Revenge is a need to avenge any damage that was done to us.

These needs are initially set to FALSE, meaning "satisfied" / "need not felt". A more complex iteration might use float values for a more nuanced state of needs.