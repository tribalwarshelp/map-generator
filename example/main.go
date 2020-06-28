package main

import (
	"log"
	"os"
	"time"

	"github.com/tribalwarshelp/map-generator/generator"
)

func main() {
	t1 := time.Now()
	f, _ := os.Create("image.jpg")
	defer f.Close()

	generator.Generate(generator.Config{
		Destination:      f,
		Scale:            10,
		ContinentGrid:    true,
		ContinentNumbers: true,
	})
	log.Print(time.Now().Sub(t1).String())
}
