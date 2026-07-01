package uglyavatar

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func (a Avatar) PNG(width, height int) ([]byte, error) {
	var b bytes.Buffer
	if err := a.WritePNG(&b, width, height); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (a Avatar) WritePNG(w io.Writer, width, height int) error {
	img, err := a.RenderImage(width, height)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

func (a Avatar) RenderImage(width, height int) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be positive")
	}
	icon, err := oksvg.ReadIconStream(strings.NewReader(a.pngSVG(width, height)))
	if err != nil {
		return nil, err
	}
	icon.SetTarget(0, 0, float64(width), float64(height))

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	if bg, ok := parseRGBColor(a.BackgroundColor); ok {
		draw.Draw(rgba, rgba.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)
	}
	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	icon.Draw(rasterx.NewDasher(width, height, scanner), 1)
	return rgba, nil
}

func (a Avatar) pngSVG(width, height int) string {
	vb := a.pngViewBox(width, height)
	sx := float64(width) / vb.W
	sy := float64(height) / vb.H
	tx := -vb.X * sx
	ty := -vb.Y * sy
	transform := "matrix(" + formatFloat(sx) + " 0 0 " + formatFloat(sy) + " " + formatFloat(tx) + " " + formatFloat(ty) + ")"
	return a.svgWithViewBoxAndTransform(width, height, svgViewBox{
		X: 0,
		Y: 0,
		W: float64(width),
		H: float64(height),
	}, transform)
}

func (a Avatar) pngViewBox(width, height int) svgViewBox {
	b := newBounds()
	b.addPoints(a.FacePoints)
	b.addPoints(a.MouthPoints)
	for _, hair := range a.Hairs {
		b.addPoints(hair.Points)
	}

	rightTx := Point{X: a.Center.X + a.DistanceBetweenEyes + a.RightEyeOffsetX, Y: -(-a.Center.Y + a.EyeHeightOffset + a.RightEyeOffsetY)}
	leftTx := Point{X: -(a.Center.X + a.DistanceBetweenEyes + a.LeftEyeOffsetX), Y: -(-a.Center.Y + a.EyeHeightOffset + a.LeftEyeOffsetY)}
	b.addTranslatedPoints(a.RightEye.Contour, rightTx)
	b.addTranslatedPoints(a.RightEye.Upper, rightTx)
	b.addTranslatedPoints(a.RightEye.Lower, rightTx)
	b.addTranslatedCircles(a.RightEye.Pupil, rightTx)
	b.addTranslatedPoints(a.LeftEye.Contour, leftTx)
	b.addTranslatedPoints(a.LeftEye.Upper, leftTx)
	b.addTranslatedPoints(a.LeftEye.Lower, leftTx)
	b.addTranslatedCircles(a.LeftEye.Pupil, leftTx)

	if a.PointNose {
		b.addCircles(a.RightNose)
		b.addCircles(a.LeftNose)
	} else {
		b.addPoint(a.LeftNoseCenter)
		b.addPoint(a.RightNoseCenter)
		b.addPoint(Point{X: (a.LeftNoseCenter.X + a.RightNoseCenter.X) / 2, Y: -a.EyeHeightOffset * 0.2})
	}

	b.addPoint(Point{X: -100, Y: -100})
	b.addPoint(Point{X: 100, Y: 100})
	return b.viewBox(12, float64(width)/float64(height))
}

type pointBounds struct {
	minX float64
	minY float64
	maxX float64
	maxY float64
	set  bool
}

func newBounds() pointBounds {
	return pointBounds{}
}

func (b *pointBounds) addPoint(p Point) {
	if !b.set {
		b.minX, b.maxX = p.X, p.X
		b.minY, b.maxY = p.Y, p.Y
		b.set = true
		return
	}
	if p.X < b.minX {
		b.minX = p.X
	}
	if p.X > b.maxX {
		b.maxX = p.X
	}
	if p.Y < b.minY {
		b.minY = p.Y
	}
	if p.Y > b.maxY {
		b.maxY = p.Y
	}
}

func (b *pointBounds) addPoints(points []Point) {
	for _, p := range points {
		b.addPoint(p)
	}
}

func (b *pointBounds) addTranslatedPoints(points []Point, tx Point) {
	for _, p := range points {
		b.addPoint(Point{X: p.X + tx.X, Y: p.Y + tx.Y})
	}
}

func (b *pointBounds) addCircles(circles []Circle) {
	for _, c := range circles {
		b.addPoint(Point{X: c.CX - c.R - c.StrokeWidth, Y: c.CY - c.R - c.StrokeWidth})
		b.addPoint(Point{X: c.CX + c.R + c.StrokeWidth, Y: c.CY + c.R + c.StrokeWidth})
	}
}

func (b *pointBounds) addTranslatedCircles(circles []Circle, tx Point) {
	for _, c := range circles {
		b.addPoint(Point{X: c.CX + tx.X - c.R - c.StrokeWidth, Y: c.CY + tx.Y - c.R - c.StrokeWidth})
		b.addPoint(Point{X: c.CX + tx.X + c.R + c.StrokeWidth, Y: c.CY + tx.Y + c.R + c.StrokeWidth})
	}
}

func (b pointBounds) viewBox(padding, aspect float64) svgViewBox {
	minX := b.minX - padding
	minY := b.minY - padding
	maxX := b.maxX + padding
	maxY := b.maxY + padding
	w := maxX - minX
	h := maxY - minY
	if aspect <= 0 {
		aspect = 1
	}
	boxAspect := w / h
	if boxAspect > aspect {
		newH := w / aspect
		delta := (newH - h) / 2
		minY -= delta
		h = newH
	} else {
		newW := h * aspect
		delta := (newW - w) / 2
		minX -= delta
		w = newW
	}
	return svgViewBox{X: minX, Y: minY, W: w, H: h}
}

func parseRGBColor(s string) (color.RGBA, bool) {
	var r, g, b uint8
	if _, err := fmt.Sscanf(s, "rgb(%d, %d, %d)", &r, &g, &b); err == nil {
		return color.RGBA{R: r, G: g, B: b, A: 255}, true
	}
	if _, err := fmt.Sscanf(s, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
		return color.RGBA{R: r, G: g, B: b, A: 255}, true
	}
	return color.RGBA{}, false
}
