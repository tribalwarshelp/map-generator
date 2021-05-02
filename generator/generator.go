package generator

import (
	"fmt"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"

	"github.com/pkg/errors"
)

const (
	defaultBackgroundColor      = "#000"
	defaultGridLineColor        = "#fff"
	defaultContinentNumberColor = "#fff"
	defaultMapSize              = 1000
	defaultQuality              = 100
)

type Marker struct {
	Larger   bool
	Villages []*twmodel.Village `json:"villages" gqlgen:"villages" xml:"villages"`
	// Color in HEX format
	Color string `json:"color" gqlgen:"color" xml:"color"`
}

type Config struct {
	Markers     []*Marker
	Destination io.Writer
	// Default 1000x1000
	MapSize          int
	ContinentGrid    bool
	ContinentNumbers bool
	// BackgroundColor in HEX format
	BackgroundColor string
	// GridLineColor in HEX format
	GridLineColor string
	// ContinentNumberColor in HEX format
	ContinentNumberColor string
	Scale                float32
	CenterX              int
	CenterY              int
	PNG                  bool
	// Quality of image (1-100), only for JPEG format
	Quality int
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
	if cfg.Quality <= 0 || cfg.Quality > 100 {
		cfg.Quality = defaultQuality
	}
	if cfg.Scale < 1 {
		cfg.Scale = 1
	}
	if cfg.CenterX <= 0 {
		cfg.CenterX = cfg.MapSize / 2
	}
	if cfg.CenterY <= 0 {
		cfg.CenterY = cfg.MapSize / 2
	}

	cfg.CenterX = int(float32(cfg.CenterX) * cfg.Scale)
	cfg.CenterY = int(float32(cfg.CenterY) * cfg.Scale)
}

func Generate(cfg Config) error {
	cfg.init()
	upLeft := image.Point{0, 0}
	lowRight := image.Point{cfg.MapSize, cfg.MapSize}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	mapSizeDividedBy10 := cfg.MapSize / 10
	imgHalfWidth := cfg.MapSize / 2
	imgHalfHeight := imgHalfWidth

	//background
	backgroundColor, err := parseHexColorFast(cfg.BackgroundColor)
	if err != nil {
		return errors.Wrap(err, "map-generator")
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	//markers
	for _, marker := range cfg.Markers {
		m := marker
		parsedColor, err := parseHexColorFast(m.Color)
		if err != nil {
			return errors.Wrap(err, "map-generator")
		}
		for _, village := range m.Villages {
			limit := 0
			if m.Larger {
				limit = 2
			}
			rect := image.Rect(village.X-limit, village.Y-limit, village.X+limit+1, village.Y+limit+1)
			draw.Draw(img, rect, &image.Uniform{parsedColor}, image.Point{}, draw.Over)
		}
	}

	//Continents
	if cfg.ContinentGrid {
		gridLineColor, err := parseHexColorFast(cfg.GridLineColor)
		if err != nil {
			return errors.Wrap(err, "map-generator")
		}
		for y := mapSizeDividedBy10; y < cfg.MapSize; y += mapSizeDividedBy10 {
			for x := 0; x < cfg.MapSize; x++ {
				img.Set(x, y, gridLineColor)
				img.Set(y, x, gridLineColor)
			}
		}
	}

	if cfg.ContinentNumbers {
		continent := 0
		continentNumberColor, err := parseHexColorFast(cfg.ContinentNumberColor)
		if err != nil {
			return errors.Wrap(err, "map-generator")
		}
		for y := mapSizeDividedBy10; y <= cfg.MapSize; y += mapSizeDividedBy10 {
			for x := mapSizeDividedBy10; x <= cfg.MapSize; x += mapSizeDividedBy10 {
				continentStr := fmt.Sprintf("%d", continent)
				if continent < 10 {
					continentStr = fmt.Sprintf("0%d", continent)
				}
				drawText(img, x-16, y-3, continentNumberColor, continentStr)
				continent++
			}
		}
	}

	var resizedImg image.Image = img
	if cfg.Scale != 1 {
		width := int(float32(cfg.MapSize) * cfg.Scale)
		resizedImg = imaging.Resize(img, width, width, imaging.NearestNeighbor)
	}

	b := resizedImg.Bounds()
	centered := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	draw.Draw(centered, b, resizedImg, image.Point{
		X: cfg.CenterX - imgHalfWidth,
		Y: cfg.CenterY - imgHalfHeight,
	}, draw.Src)

	if cfg.PNG {
		if err := png.Encode(cfg.Destination, centered); err != nil {
			return errors.Wrap(err, "map-generator")
		}
	} else {
		if err := jpeg.Encode(cfg.Destination, centered, &jpeg.Options{
			Quality: cfg.Quality,
		}); err != nil {
			return errors.Wrap(err, "map-generator")
		}
	}
	return nil
}
