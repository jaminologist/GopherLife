package controllers

import "strconv"

type FormData struct {
	DisplayName        string
	Name               string
	Value              string
	Type               string
	BootStrapFormWidth int
}

func FormDataWidth(width int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Width",
		Type:               "Number",
		Name:               "width",
		Value:              strconv.Itoa(width),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

func FormDataHeight(height int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Height",
		Type:               "Number",
		Name:               "height",
		Value:              strconv.Itoa(height),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

func FormDataInitialPopulation(initialPopulation int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Initial Population",
		Type:               "Number",
		Name:               "initialPopulation",
		Value:              strconv.Itoa(initialPopulation),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

func FormDataMaxPopulation(maxPopulation int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Max Population",
		Type:               "Number",
		Name:               "maxPopulation",
		Value:              strconv.Itoa(maxPopulation),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

func FormDataSnakeSlowDown(slowdown int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Snake SlowDown",
		Type:               "Number",
		Name:               "SnakeSlowDown",
		Value:              strconv.Itoa(slowdown),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}
