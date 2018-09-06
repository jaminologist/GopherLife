package food

type Food struct {
	Name   string
	Energy int
}

func NewCarrot() Food {
	return Food{Name: "Carrot", Energy: 5}
}

func NewPotato() Food {
	return Food{Name: "Patato", Energy: 10}
}

func New() Food {
	return Food{Name: "Bean", Energy: 5}
}
