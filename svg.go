package uglyavatar

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

type svgViewBox struct {
	X float64
	Y float64
	W float64
	H float64
}

func (v svgViewBox) String() string {
	return formatFloat(v.X) + " " + formatFloat(v.Y) + " " + formatFloat(v.W) + " " + formatFloat(v.H)
}

func (a Avatar) SVGWithSize(width, height int) string {
	return a.svgWithViewBoxAndTransform(width, height, svgViewBox{X: -100, Y: -100, W: 200, H: 200}, "")
}

func (a Avatar) svgWithViewBox(width, height int, vb svgViewBox) string {
	return a.svgWithViewBoxAndTransform(width, height, vb, "")
}

func (a Avatar) svgWithViewBoxAndTransform(width, height int, vb svgViewBox, contentTransform string) string {
	var b strings.Builder
	b.WriteString(`<svg viewBox="`)
	b.WriteString(vb.String())
	b.WriteString(`" xmlns="http://www.w3.org/2000/svg" width="`)
	b.WriteString(strconv.Itoa(width))
	b.WriteString(`" height="`)
	b.WriteString(strconv.Itoa(height))
	b.WriteString(`">`)
	b.WriteString(`<defs>`)
	b.WriteString(`<clipPath id="leftEyeClipPath"><polyline points="`)
	b.WriteString(points(a.LeftEye.Contour))
	b.WriteString(`"/></clipPath>`)
	b.WriteString(`<clipPath id="rightEyeClipPath"><polyline points="`)
	b.WriteString(points(a.RightEye.Contour))
	b.WriteString(`"/></clipPath>`)
	b.WriteString(`<filter id="fuzzy"><feTurbulence id="turbulence" baseFrequency="0.05" numOctaves="3" type="noise" result="noise"/><feDisplacementMap in="SourceGraphic" in2="noise" scale="2"/></filter>`)
	b.WriteString(`<linearGradient id="rainbowGradient" x1="0%" y1="0%" x2="100%" y2="0%">`)
	b.WriteString(`<stop offset="0%" style="stop-color: ` + html.EscapeString(svgColor(a.RainbowGradient[0])) + `; stop-opacity: 1"/>`)
	b.WriteString(`<stop offset="` + html.EscapeString(a.DyeColorOffset) + `" style="stop-color: ` + html.EscapeString(svgColor(a.RainbowGradient[1])) + `; stop-opacity: 1"/>`)
	b.WriteString(`<stop offset="100%" style="stop-color: ` + html.EscapeString(svgColor(a.RainbowGradient[2])) + `; stop-opacity: 1"/>`)
	b.WriteString(`</linearGradient></defs>`)
	b.WriteString(`<title>That's an ugly face</title>`)
	b.WriteString(`<desc>CREATED BY XUAN TANG, GO PORT AT github.com/txstc55/uglyavatar</desc>`)
	b.WriteString(`<rect x="` + formatFloat(vb.X) + `" y="` + formatFloat(vb.Y) + `" width="` + formatFloat(vb.W) + `" height="` + formatFloat(vb.H) + `" fill="` + html.EscapeString(svgColor(a.BackgroundColor)) + `"/>`)
	if contentTransform != "" {
		b.WriteString(`<g transform="` + contentTransform + `">`)
	}
	b.WriteString(`<polyline id="faceContour" points="` + points(a.FacePoints) + `" fill="#ffc9a9" stroke="black" stroke-width="` + formatFloat(3/a.FaceScale) + `" stroke-linejoin="round" filter="url(#fuzzy)"/>`)

	rightTransform := `translate(` + formatFloat(a.Center.X+a.DistanceBetweenEyes+a.RightEyeOffsetX) + ` ` + formatFloat(-(-a.Center.Y + a.EyeHeightOffset + a.RightEyeOffsetY)) + `)`
	leftTransform := `translate(` + formatFloat(-(a.Center.X + a.DistanceBetweenEyes + a.LeftEyeOffsetX)) + ` ` + formatFloat(-(-a.Center.Y + a.EyeHeightOffset + a.LeftEyeOffsetY)) + `)`

	b.WriteString(`<g transform="` + rightTransform + `"><polyline id="rightCountour" points="` + points(a.RightEye.Contour) + `" fill="white" stroke="white" stroke-width="0" stroke-linejoin="round" filter="url(#fuzzy)"/></g>`)
	b.WriteString(`<g transform="` + leftTransform + `"><polyline id="leftCountour" points="` + points(a.LeftEye.Contour) + `" fill="white" stroke="white" stroke-width="0" stroke-linejoin="round" filter="url(#fuzzy)"/></g>`)

	eyeStroke := 3.0
	if a.HaventSleptForDays {
		eyeStroke = 5.0
	}
	eyeStroke /= a.FaceScale
	b.WriteString(`<g transform="` + rightTransform + `">`)
	b.WriteString(`<polyline id="rightUpper" points="` + points(a.RightEye.Upper) + `" fill="none" stroke="black" stroke-width="` + formatFloat(eyeStroke) + `" stroke-linejoin="round" stroke-linecap="round" filter="url(#fuzzy)"/>`)
	b.WriteString(`<polyline id="rightLower" points="` + points(a.RightEye.Lower) + `" fill="none" stroke="black" stroke-width="` + formatFloat(eyeStroke) + `" stroke-linejoin="round" stroke-linecap="round" filter="url(#fuzzy)"/>`)
	writeCircles(&b, a.RightEye.Pupil, "rightEyeClipPath")
	b.WriteString(`</g>`)

	b.WriteString(`<g transform="` + leftTransform + `">`)
	b.WriteString(`<polyline id="leftUpper" points="` + points(a.LeftEye.Upper) + `" fill="none" stroke="black" stroke-width="` + formatFloat(eyeStroke) + `" stroke-linejoin="round" filter="url(#fuzzy)"/>`)
	b.WriteString(`<polyline id="leftLower" points="` + points(a.LeftEye.Lower) + `" fill="none" stroke="black" stroke-width="` + formatFloat(eyeStroke) + `" stroke-linejoin="round" filter="url(#fuzzy)"/>`)
	writeCircles(&b, a.LeftEye.Pupil, "leftEyeClipPath")
	b.WriteString(`</g>`)

	b.WriteString(`<g id="hairs">`)
	for _, hair := range a.Hairs {
		b.WriteString(`<polyline points="` + points(hair.Points) + `" fill="none" stroke="` + html.EscapeString(svgColor(a.HairColor)) + `" stroke-width="` + formatFloat(hair.StrokeWidth) + `" stroke-linejoin="round" filter="url(#fuzzy)"/>`)
	}
	b.WriteString(`</g>`)

	if a.PointNose {
		b.WriteString(`<g id="pointNose"><g id="rightNose">`)
		writeCircles(&b, a.RightNose, "")
		b.WriteString(`</g><g id="leftNose">`)
		writeCircles(&b, a.LeftNose, "")
		b.WriteString(`</g></g>`)
	} else {
		d := `M ` + formatFloat(a.LeftNoseCenter.X) + ` ` + formatFloat(a.LeftNoseCenter.Y) +
			`, Q` + formatFloat(a.RightNoseCenter.X) + ` ` + formatFloat(a.RightNoseCenter.Y*1.5) +
			`,` + formatFloat((a.LeftNoseCenter.X+a.RightNoseCenter.X)/2) + ` ` + formatFloat(-a.EyeHeightOffset*0.2)
		b.WriteString(`<g id="lineNose"><path d="` + d + `" fill="none" stroke="black" stroke-width="` + formatFloat(a.NoseStrokeWidth) + `" stroke-linejoin="round" filter="url(#fuzzy)"/></g>`)
	}
	b.WriteString(`<g id="mouth"><polyline points="` + points(a.MouthPoints) + `" fill="#d77f8c" stroke="black" stroke-width="` + formatFloat(a.MouthStrokeWidth) + `" stroke-linejoin="round" filter="url(#fuzzy)"/></g>`)
	if contentTransform != "" {
		b.WriteString(`</g>`)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

func writeCircles(b *strings.Builder, circles []Circle, clipPath string) {
	for _, c := range circles {
		b.WriteString(`<circle r="` + formatFloat(c.R) + `" cx="` + formatFloat(c.CX) + `" cy="` + formatFloat(c.CY) + `" stroke="black" fill="none" stroke-width="` + formatFloat(c.StrokeWidth) + `" filter="url(#fuzzy)"`)
		if clipPath != "" {
			b.WriteString(` clip-path="url(#` + clipPath + `)"`)
		}
		b.WriteString(`/>`)
	}
}

func points(ps []Point) string {
	var b strings.Builder
	for i, p := range ps {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(formatFloat(p.X))
		b.WriteByte(',')
		b.WriteString(formatFloat(p.Y))
	}
	return b.String()
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func svgColor(color string) string {
	var r, g, b int
	if _, err := fmt.Sscanf(color, "rgb(%d, %d, %d)", &r, &g, &b); err == nil {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	if _, err := fmt.Sscanf(color, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return color
}
