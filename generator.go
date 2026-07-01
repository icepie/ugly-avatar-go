package uglyavatar

import "math"

func (g *Generator) Generate() Avatar {
	var a Avatar
	a.FaceScale = 1.5 + g.r.Float64()*0.6
	a.HaventSleptForDays = g.r.Float64() > 0.8
	a.FacePoints, a.FaceWidth, a.FaceHeight, a.Center = g.generateFaceContourPoints(100)

	left, right := g.generateBothEyes(a.FaceWidth / 2)
	a.LeftEye.Upper = left.upper
	a.LeftEye.Lower = left.lower
	a.RightEye.Upper = right.upper
	a.RightEye.Lower = right.lower
	a.LeftEye.Contour = append(copyPoints(left.upper[10:90]), reversePoints(left.lower[10:90])...)
	a.RightEye.Contour = append(copyPoints(right.upper[10:90]), reversePoints(right.lower[10:90])...)

	a.DistanceBetweenEyes = g.random(a.FaceWidth/4.5, a.FaceWidth/4)
	a.EyeHeightOffset = g.random(a.FaceHeight/8, a.FaceHeight/6)
	a.LeftEyeOffsetX = g.random(-a.FaceWidth/20, a.FaceWidth/10)
	a.LeftEyeOffsetY = g.random(-a.FaceHeight/50, a.FaceHeight/50)
	a.RightEyeOffsetX = g.random(-a.FaceWidth/20, a.FaceWidth/10)
	a.RightEyeOffsetY = g.random(-a.FaceHeight/50, a.FaceHeight/50)

	leftInd0 := int(math.Floor(g.random(10, float64(len(left.upper)-10))))
	rightInd0 := int(math.Floor(g.random(10, float64(len(right.upper)-10))))
	leftInd1 := int(math.Floor(g.random(10, float64(len(left.upper)-10))))
	rightInd1 := int(math.Floor(g.random(10, float64(len(right.upper)-10))))
	leftLerp := g.random(0.2, 0.8)
	rightLerp := g.random(0.2, 0.8)
	a.LeftPupilShift = Point{
		X: left.upper[leftInd0].X*leftLerp + left.lower[leftInd1].X*(1-leftLerp),
		Y: left.upper[leftInd0].Y*leftLerp + left.lower[leftInd1].Y*(1-leftLerp),
	}
	a.RightPupilShift = Point{
		X: right.upper[rightInd0].X*rightLerp + right.lower[rightInd1].X*(1-rightLerp),
		Y: right.upper[rightInd0].Y*rightLerp + right.lower[rightInd1].Y*(1-rightLerp),
	}
	a.LeftEye.Pupil = g.pupilCircles(a.LeftPupilShift)
	a.RightEye.Pupil = g.pupilCircles(a.RightPupilShift)

	numHairLines := make([]int, 4)
	for i := range numHairLines {
		numHairLines[i] = int(math.Floor(g.random(0, 50)))
	}
	var hairs [][]Point
	if g.r.Float64() > 0.3 {
		hairs = g.generateHairLines0(a.FacePoints, numHairLines[0]+10)
	}
	if g.r.Float64() > 0.3 {
		hairs = append(hairs, g.generateHairLines1(a.FacePoints, int(float64(numHairLines[1])/1.5+10))...)
	}
	if g.r.Float64() > 0.5 {
		hairs = append(hairs, g.generateHairLines2(a.FacePoints, numHairLines[2]*3+10)...)
	}
	if g.r.Float64() > 0.5 {
		hairs = append(hairs, g.generateHairLines3(a.FacePoints, numHairLines[3]*3+10)...)
	}
	a.Hairs = make([]HairLine, 0, len(hairs))
	for _, h := range hairs {
		a.Hairs = append(a.Hairs, HairLine{Points: h, StrokeWidth: 0.5 + g.r.Float64()*2.5})
	}

	a.RightNoseCenter = Point{X: g.random(a.FaceWidth/18, a.FaceWidth/12), Y: g.random(0, a.FaceHeight/5)}
	a.LeftNoseCenter = Point{
		X: g.random(-a.FaceWidth/18, -a.FaceWidth/12),
		Y: a.RightNoseCenter.Y + g.random(-a.FaceHeight/30, a.FaceHeight/20),
	}
	a.PointNose = g.r.Float64() > 0.5
	if a.PointNose {
		a.RightNose = g.noseCircles(a.RightNoseCenter)
		a.LeftNose = g.noseCircles(a.LeftNoseCenter)
	} else {
		a.NoseStrokeWidth = 2.5 + g.r.Float64()
	}

	if g.r.Float64() > 0.1 {
		a.HairColor = hairColors[int(math.Floor(g.r.Float64()*10))]
	} else {
		a.HairColor = "url(#rainbowGradient)"
		a.DyeColorOffset = formatFloat(g.random(0, 100)) + "%"
		a.RainbowGradient = []string{
			hairColors[int(math.Floor(g.r.Float64()*10))],
			hairColors[int(math.Floor(g.r.Float64()*float64(len(hairColors))))],
			hairColors[int(math.Floor(g.r.Float64()*float64(len(hairColors))))],
		}
	}
	if len(a.RainbowGradient) == 0 {
		a.RainbowGradient = []string{hairColors[0], hairColors[1], hairColors[2]}
		a.DyeColorOffset = "50%"
	}
	a.BackgroundColor = backgroundColors[int(math.Floor(g.r.Float64()*float64(len(backgroundColors))))]

	switch int(math.Floor(g.r.Float64() * 3)) {
	case 0:
		a.MouthPoints = g.generateMouthShape0(a.FacePoints, a.FaceHeight, a.FaceWidth)
	case 1:
		a.MouthPoints = g.generateMouthShape1(a.FacePoints, a.FaceHeight, a.FaceWidth)
	default:
		a.MouthPoints = g.generateMouthShape2(a.FacePoints, a.FaceHeight, a.FaceWidth)
	}
	a.MouthStrokeWidth = 2.7 + g.r.Float64()*0.5
	return a
}

func (g *Generator) pupilCircles(center Point) []Circle {
	circles := make([]Circle, 0, 10)
	for i := 0; i < 10; i++ {
		circles = append(circles, Circle{
			R:           g.r.Float64()*2 + 3,
			CX:          center.X + g.r.Float64()*5 - 2.5,
			CY:          center.Y + g.r.Float64()*5 - 2.5,
			StrokeWidth: 1 + g.r.Float64()*0.5,
		})
	}
	return circles
}

func (g *Generator) noseCircles(center Point) []Circle {
	circles := make([]Circle, 0, 10)
	for i := 0; i < 10; i++ {
		circles = append(circles, Circle{
			R:           g.r.Float64()*2 + 1,
			CX:          center.X + g.r.Float64()*4 - 2,
			CY:          center.Y + g.r.Float64()*4 - 2,
			StrokeWidth: 1 + g.r.Float64()*0.5,
		})
	}
	return circles
}
