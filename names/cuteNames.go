package names

import (
	"math/rand"
)

//CuteName Returns a random cute name
func CuteName() string {
	return collection[rand.Intn(len(collection))] + "-" + collection[rand.Intn(len(collection))]
}

var collection = []string{
	"Abby",
	"Allie",
	"Annie",
	"Angel",
	"Amos",
	"Apple",
	"Archie",
	"Bailey",
	"BammBamm",
	"Beau",
	"Bella",
	"Biscuit",
	"Bluebell",
	"Bolt",
	"Bonnie",
	"Boots",
	"Buddy",
	"Callie",
	"Caramel",
	"Charlie",
	"Cherry",
	"Chloe",
	"Clover",
	"Coco",
	"Colby",
	"Cooper",
	"Cricket",
	"Cupcake",
	"Dash",
	"Dixie",
	"Dolce",
	"Donut",
	"Eddie",
	"Ellie",
	"Finn",
	"Flower",
	"Frankie",
	"Frodo",
	"Gigi",
	"Goldie",
	"Goose",
	"Grace",
	"Gulliver",
	"Gus",
	"Harper",
	"Henry",
	"Herbie",
	"Izzy",
	"Jack",
	"Jojo",
	"Josie",
	"Kai",
	"Katie",
	"Kiki",
	"Lennon",
	"Leo",
	"Lillie",
	"Lily",
	"Lottie",
	"Lucy",
	"Lulu",
	"Maggie",
	"Marshmallow",
	"Marley",
	"Mojo",
	"Molly",
	"Nala",
	"Nessie",
	"Oliver",
	"Otis",
	"Otto",
	"Oscar",
	"Pancakes",
	"Pansy",
	"Peaches",
	"Pip",
	"Pluto",
	"Poppy",
	"Pumpkin",
	"Queen",
	"Quinn",
	"Rhubarb",
	"Riley",
	"River",
	"Rosie",
	"Sadie",
	"Sophie",
	"Stella",
	"Sugar",
	"Tillie",
	"Toby",
	"Tucker",
	"Violet",
	"Waffles",
	"Winnie",
	"Winston",
	"Zacky",
	"Zoe",
	"Zuzu",
}
