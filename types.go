package uglyavatar

import (
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	defaultSize = 500
	viewBox     = "-100 -100 200 200"
)

type Point struct {
	X float64
	Y float64
}

type Circle struct {
	CX          float64
	CY          float64
	R           float64
	StrokeWidth float64
}

type HairLine struct {
	Points      []Point
	StrokeWidth float64
}

type Eye struct {
	Upper   []Point
	Lower   []Point
	Contour []Point
	Pupil   []Circle
}

type Avatar struct {
	FaceScale           float64
	FacePoints          []Point
	FaceWidth           float64
	FaceHeight          float64
	Center              Point
	DistanceBetweenEyes float64
	EyeHeightOffset     float64
	LeftEyeOffsetX      float64
	LeftEyeOffsetY      float64
	RightEyeOffsetX     float64
	RightEyeOffsetY     float64
	LeftPupilShift      Point
	RightPupilShift     Point
	LeftEye             Eye
	RightEye            Eye
	Hairs               []HairLine
	HaventSleptForDays  bool
	HairColor           string
	BackgroundColor     string
	RainbowGradient     []string
	DyeColorOffset      string
	PointNose           bool
	RightNoseCenter     Point
	LeftNoseCenter      Point
	RightNose           []Circle
	LeftNose            []Circle
	NoseStrokeWidth     float64
	MouthPoints         []Point
	MouthStrokeWidth    float64
}

type Generator struct {
	r *rand.Rand
}

func New() *Generator {
	return NewWithRand(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func NewWithSeed(seed int64) *Generator {
	return NewWithRand(rand.New(rand.NewSource(seed)))
}

func NewWithRand(r *rand.Rand) *Generator {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &Generator{r: r}
}

func Generate() Avatar {
	return New().Generate()
}

func GenerateWithSeed(seed int64) Avatar {
	return NewWithSeed(seed).Generate()
}

func (a Avatar) WriteSVG(w io.Writer) error {
	_, err := io.WriteString(w, a.SVG())
	return err
}

func (a Avatar) SVG() string {
	return a.SVGWithSize(defaultSize, defaultSize)
}

func (a Avatar) SavePNG(path string, width, height int) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return a.WritePNG(out, width, height)
}
