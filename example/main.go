package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/tribalwarshelp/map-generator/generator"
	"github.com/tribalwarshelp/shared/models"
)

func main() {
	villages := generateVillages(10)
	// villages2 := generateVillages(1)
	t1 := time.Now()
	f, _ := os.Create("image.jpeg")
	defer f.Close()

	err := generator.Generate(generator.Config{
		Destination:      f,
		Scale:            1,
		ContinentGrid:    true,
		ContinentNumbers: true,
		Markers: []*generator.Marker{
			&generator.Marker{
				Color:    "#fff",
				Villages: villages,
				Larger:   false,
			},
			&generator.Marker{
				Color: "#ff0",
				Villages: []*models.Village{
					&models.Village{
						X: 500,
						Y: 500,
					},
				},
				Larger: true,
			},
		},
	})
	log.Println(time.Now().Sub(t1).String(), err)
}

func generateVillages(ch int) []*models.Village {
	villages := []*models.Village{}
	for y := 0; y <= 1000; y++ {
		for x := 0; x <= 1000; x++ {
			if rand.Intn(100) <= ch {
				villages = append(villages, &models.Village{
					X: x,
					Y: y,
				})
			}
		}
	}
	return villages
}
