package names

import "math/rand"

func GetCuteName() string {

	return firstNames[rand.Intn(len(firstNames))]

}

var firstNames = []string{
	"Airy",
	"Foofy",
	"River",
	"Rob",
	"Zacky",
}
