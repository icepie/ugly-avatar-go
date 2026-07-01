package uglyavatar

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestGenerateWithSeedIsDeterministic(t *testing.T) {
	a := GenerateWithSeed(42)
	b := GenerateWithSeed(42)
	if a.SVG() != b.SVG() {
		t.Fatal("same seed produced different SVG")
	}
}

func TestGeneratedAvatarHasExpectedParts(t *testing.T) {
	a := GenerateWithSeed(7)
	if len(a.FacePoints) != 402 {
		t.Fatalf("face points = %d, want 402", len(a.FacePoints))
	}
	if len(a.LeftEye.Upper) != 100 || len(a.RightEye.Upper) != 100 {
		t.Fatalf("unexpected eye point counts: left=%d right=%d", len(a.LeftEye.Upper), len(a.RightEye.Upper))
	}
	if len(a.MouthPoints) == 0 {
		t.Fatal("mouth points are empty")
	}
	svg := a.SVG()
	for _, want := range []string{`<svg`, `id="faceContour"`, `id="hairs"`, `id="mouth"`} {
		if !strings.Contains(svg, want) {
			t.Fatalf("svg missing %q", want)
		}
	}
}

func TestWritePNG(t *testing.T) {
	a := GenerateWithSeed(11)
	var b bytes.Buffer
	if err := a.WritePNG(&b, 128, 128); err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if got := img.Bounds().Dx(); got != 128 {
		t.Fatalf("png width = %d, want 128", got)
	}
	if got := img.Bounds().Dy(); got != 128 {
		t.Fatalf("png height = %d, want 128", got)
	}

	bg, err := parseRGB(a.BackgroundColor)
	if err != nil {
		t.Fatal(err)
	}
	minX, minY := 128, 128
	maxX, maxY := -1, -1
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			if pixelDifferent(img.At(x, y), bg) {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	for _, p := range []struct{ x, y int }{{0, 0}, {127, 0}, {0, 127}, {127, 127}} {
		if pixelDifferent(img.At(p.x, p.y), bg) {
			t.Fatalf("png corner (%d,%d) is not background: got %v, want %v", p.x, p.y, img.At(p.x, p.y), bg)
		}
	}
	if maxX-minX < 50 || maxY-minY < 50 {
		t.Fatalf("rendered content bbox too small, bbox=(%d,%d)-(%d,%d)", minX, minY, maxX, maxY)
	}
	if minX >= 64 || maxX <= 64 || minY >= 64 || maxY <= 64 {
		t.Fatalf("rendered content does not cover png center, bbox=(%d,%d)-(%d,%d)", minX, minY, maxX, maxY)
	}
}

func parseRGB(s string) (color.RGBA, error) {
	var c color.RGBA
	_, err := fmt.Sscanf(s, "rgb(%d, %d, %d)", &c.R, &c.G, &c.B)
	c.A = 255
	return c, err
}

func pixelDifferent(c color.Color, bg color.RGBA) bool {
	r, g, b, a := c.RGBA()
	return uint8(r>>8) != bg.R || uint8(g>>8) != bg.G || uint8(b>>8) != bg.B || uint8(a>>8) != bg.A
}
