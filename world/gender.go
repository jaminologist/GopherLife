package world

import "math/rand"

const (
	Male   Gender = 0
	Female Gender = 1
)

var genders = [2]Gender{Male, Female}

func (gender Gender) String() string {

	names := [...]string{
		"Male",
		"Female"}

	if gender < Male || gender > Female {
		return "Unknown"
	}

	return names[gender]
}

//Opposite Returns the opposite gender
func (gender Gender) Opposite() Gender {

	switch gender {
	case Male:
		return Female
	case Female:
		return Male
	default:
		return Male
	}
}

//GetRandomGender Returns a random gender.
func GetRandomGender() Gender {
	return genders[rand.Intn(len(genders))]
}
