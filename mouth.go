package uglyavatar

import "math"

func (g *Generator) generateMouthShape0(_ []Point, faceHeight, faceWidth float64) []Point {
	mouthRightY := g.random(faceHeight/7, faceHeight/3.5)
	mouthLeftY := g.random(faceHeight/7, faceHeight/3.5)
	mouthRightX := g.random(faceWidth/10, faceWidth/2)
	mouthLeftX := -mouthRightX + g.random(-faceWidth/20, faceWidth/20)
	mouthRight := Point{X: mouthRightX, Y: mouthRightY}
	mouthLeft := Point{X: mouthLeftX, Y: mouthLeftY}
	control0 := Point{X: g.random(0, mouthRightX), Y: g.random(mouthLeftY+5, faceHeight/1.5)}
	control1 := Point{X: g.random(mouthLeftX, 0), Y: g.random(mouthLeftY+5, faceHeight/1.5)}

	points := make([]Point, 0, 200)
	for i := 0; i < 100; i++ {
		points = append(points, cubicBezier(mouthLeft, control1, control0, mouthRight, float64(i)/100))
	}
	if g.r.Float64() > 0.5 {
		for i := 0; i < 100; i++ {
			points = append(points, cubicBezier(mouthRight, control0, control1, mouthLeft, float64(i)/100))
		}
	} else {
		yOffsetPortion := g.random(0, 0.8)
		for i := 0; i < 100; i++ {
			t := float64(i) / 100
			points = append(points, Point{
				X: points[99].X*(1-t) + points[0].X*t,
				Y: (points[99].Y*(1-t)+points[0].Y*t)*(1-yOffsetPortion) + points[99-i].Y*yOffsetPortion,
			})
		}
	}
	return points
}

func (g *Generator) generateMouthShape1(_ []Point, faceHeight, faceWidth float64) []Point {
	mouthRightY := g.random(faceHeight/7, faceHeight/4)
	mouthLeftY := g.random(faceHeight/7, faceHeight/4)
	mouthRightX := g.random(faceWidth/10, faceWidth/2)
	mouthLeftX := -mouthRightX + g.random(-faceWidth/20, faceWidth/20)
	mouthRight := Point{X: mouthRightX, Y: mouthRightY}
	mouthLeft := Point{X: mouthLeftX, Y: mouthLeftY}
	control0 := Point{X: g.random(0, mouthRightX), Y: g.random(mouthLeftY+5, faceHeight/1.5)}
	control1 := Point{X: g.random(mouthLeftX, 0), Y: g.random(mouthLeftY+5, faceHeight/1.5)}

	points := make([]Point, 0, 200)
	for i := 0; i < 100; i++ {
		points = append(points, cubicBezier(mouthLeft, control1, control0, mouthRight, float64(i)/100))
	}
	center := Point{X: (mouthRight.X + mouthLeft.X) / 2, Y: points[25].Y/2 + points[75].Y/2}
	if g.r.Float64() > 0.5 {
		for i := 0; i < 100; i++ {
			points = append(points, cubicBezier(mouthRight, control0, control1, mouthLeft, float64(i)/100))
		}
	} else {
		yOffsetPortion := g.random(0, 0.8)
		for i := 0; i < 100; i++ {
			t := float64(i) / 100
			points = append(points, Point{
				X: points[99].X*(1-t) + points[0].X*t,
				Y: (points[99].Y*(1-t)+points[0].Y*t)*(1-yOffsetPortion) + points[99-i].Y*yOffsetPortion,
			})
		}
	}
	for i := range points {
		points[i].X -= center.X
		points[i].Y -= center.Y
		points[i].Y = -points[i].Y
		points[i].X *= 0.6
		points[i].Y *= 0.6
		points[i].X += center.X
		points[i].Y += center.Y * 0.8
	}
	return points
}

func (g *Generator) generateMouthShape2(_ []Point, faceHeight, faceWidth float64) []Point {
	center := Point{X: g.random(-faceWidth/8, faceWidth/8), Y: g.random(faceHeight/4, faceHeight/2.5)}
	points := g.eggShapePoints(g.random(faceWidth/4, faceWidth/10), g.random(faceHeight/10, faceHeight/20), 0.001, 50)
	rotation := g.random(-math.Pi/9.5, math.Pi/9.5)
	cosR := math.Cos(rotation)
	sinR := math.Sin(rotation)
	for i := range points {
		x := points[i].X
		y := points[i].Y
		points[i].X = x*cosR - y*sinR + center.X
		points[i].Y = x*sinR + y*cosR + center.Y
	}
	return points
}
