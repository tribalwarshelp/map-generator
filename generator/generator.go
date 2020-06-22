package generator

import (
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

const (
	defaultBackgroundColor      = "#000"
	defaultGridLineColor        = "#fff"
	defaultContinentNumberColor = "#fff"
	defaultMapSize              = 1000
)

type Marker struct {
	Villages []*models.Village `json:"villages" gqlgen:"villages" xml:"villages"`
	Color    string            `json:"color" gqlgen:"color" xml:"color"`
}

type Config struct {
	Markers              []*Marker
	Destination          io.Writer
	MapSize              int
	ContinentGrid        bool
	ContinentNumbers     bool
	BackgroundColor      string
	GridLineColor        string
	ContinentNumberColor string
}

func (cfg *Config) init() {
	if cfg.BackgroundColor == "" {
		cfg.BackgroundColor = defaultBackgroundColor
	}
	if cfg.GridLineColor == "" {
		cfg.GridLineColor = defaultGridLineColor
	}
	if cfg.ContinentNumberColor == "" {
		cfg.ContinentNumberColor = defaultContinentNumberColor
	}
	if cfg.MapSize <= 0 {
		cfg.MapSize = defaultMapSize
	}
}

func Generate(cfg Config) error {
	cfg.init()
	upLeft := image.Point{0, 0}
	lowRight := image.Point{cfg.MapSize, cfg.MapSize}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	backgroundColor, err := parseHexColorFast(cfg.BackgroundColor)
	if err != nil {
		return errors.Wrap(err, "map-generator")
	}

	// Background.
	for y := 0; y < cfg.MapSize; y++ {
		for x := 0; x < cfg.MapSize; x++ {
			img.Set(x, y, backgroundColor)
		}
	}

	// Markers
	for _, marker := range cfg.Markers {
		parsedColor, err := parseHexColorFast(marker.Color)
		if err != nil {
			return err
		}
		for _, village := range marker.Villages {
			img.Set(village.X, village.Y, parsedColor)
		}
	}

	//Continents
	if cfg.ContinentGrid {
		gridLineColor, err := parseHexColorFast(cfg.GridLineColor)
		if err != nil {
			return errors.Wrap(err, "map-generator")
		}
		for y := cfg.MapSize / 10; y < cfg.MapSize; y += cfg.MapSize / 10 {
			for x := 0; x < cfg.MapSize; x++ {
				img.Set(x, y, gridLineColor)
			}
		}
		for x := cfg.MapSize / 10; x < cfg.MapSize; x += cfg.MapSize / 10 {
			for y := 0; y < cfg.MapSize; y++ {
				img.Set(x, y, gridLineColor)
			}
		}
	}

	if cfg.ContinentNumbers {
		continent := 0
		continentNumberColor, err := parseHexColorFast(cfg.ContinentNumberColor)
		if err != nil {
			return errors.Wrap(err, "map-generator")
		}
		for y := cfg.MapSize / 10; y <= cfg.MapSize; y += cfg.MapSize / 10 {
			for x := cfg.MapSize / 10; x <= cfg.MapSize; x += cfg.MapSize / 10 {
				continentStr := fmt.Sprintf("%d", continent)
				if continent < 10 {
					continentStr = fmt.Sprintf("0%d", continent)
				}
				drawText(img, x-16, y-3, continentNumberColor, continentStr)
				continent++
			}
		}
	}

	if err := png.Encode(cfg.Destination, img); err != nil {
		return errors.Wrap(err, "map-generator")
	}
	return nil
}
