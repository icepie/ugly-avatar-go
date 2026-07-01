package uglyavatar

import "math"

func (g *Generator) generateHairLines0(faceContour []Point, numHairLines int) [][]Point {
	face := faceContour[:len(faceContour)-2]
	results := make([][]Point, 0, numHairLines)
	for i := 0; i < numHairLines; i++ {
		numHairPoints := 20 + int(math.Floor(g.random(-5, 5)))
		hairLine := make([]Point, 0, numHairPoints)
		indexOffset := int(math.Floor(g.random(30, 140)))
		for j := 0; j < numHairPoints; j++ {
			hairLine = append(hairLine, face[(len(face)-(j+indexOffset))%len(face)])
		}
		d0 := computeBezierCurve(hairLine, numHairPoints)
		hairLine = hairLine[:0]
		indexOffset = int(math.Floor(g.random(30, 140)))
		for j := 0; j < numHairPoints; j++ {
			hairLine = append(hairLine, face[(len(face)-(-j+indexOffset))%len(face)])
		}
		d1 := computeBezierCurve(hairLine, numHairPoints)
		d := make([]Point, 0, numHairPoints)
		for j := 0; j < numHairPoints; j++ {
			weight := pow2(float64(j) / float64(numHairPoints))
			d = append(d, Point{
				X: d0[j].X*weight + d1[j].X*(1-weight),
				Y: d0[j].Y*weight + d1[j].Y*(1-weight),
			})
		}
		results = append(results, d)
	}
	return results
}

func (g *Generator) generateHairLines1(faceContour []Point, numHairLines int) [][]Point {
	face := faceContour[:len(faceContour)-2]
	results := make([][]Point, 0, numHairLines)
	for i := 0; i < numHairLines; i++ {
		numHairPoints := 20 + int(math.Floor(g.random(-5, 5)))
		hairLine := make([]Point, 0, numHairPoints+1)
		indexStart := int(math.Floor(g.random(20, 160)))
		hairLine = append(hairLine, face[(len(face)-indexStart)%len(face)])
		for j := 1; j < numHairPoints+1; j++ {
			indexStart = int(math.Floor(g.random(20, 160)))
			hairLine = append(hairLine, face[(len(face)-indexStart)%len(face)])
		}
		results = append(results, computeBezierCurve(hairLine, numHairPoints))
	}
	return results
}

func (g *Generator) generateHairLines2(faceContour []Point, numHairLines int) [][]Point {
	face := faceContour[:len(faceContour)-2]
	results := make([][]Point, 0, numHairLines)
	pickedIndices := make([]int, 0, numHairLines)
	for i := 0; i < numHairLines; i++ {
		pickedIndices = append(pickedIndices, int(math.Floor(g.random(10, 180))))
	}
	sortInts(pickedIndices)
	for i := 0; i < numHairLines; i++ {
		numHairPoints := 20 + int(math.Floor(g.random(-5, 5)))
		hairLine := make([]Point, 0, numHairPoints)
		indexOffset := pickedIndices[i]
		lower := g.random(0.8, 1.4)
		reverse := 1
		if g.r.Float64() <= 0.5 {
			reverse = -1
		}
		for j := 0; j < numHairPoints; j++ {
			powerScale := g.random(0.1, 3)
			portion := (1-math.Pow(float64(j)/float64(numHairPoints), powerScale))*(1-lower) + lower
			p := face[(len(face)-(reverse*j+indexOffset))%len(face)]
			hairLine = append(hairLine, Point{X: p.X * portion, Y: p.Y * portion})
		}
		d := computeBezierCurve(hairLine, numHairPoints)
		if g.r.Float64() > 0.7 {
			d = reversePoints(d)
		}
		if len(results) == 0 {
			results = append(results, d)
			continue
		}
		last := results[len(results)-1][len(results[len(results)-1])-1]
		dist := math.Hypot(d[0].X-last.X, d[0].Y-last.Y)
		if g.r.Float64() > 0.5 && dist < 100 {
			results[len(results)-1] = append(results[len(results)-1], d...)
		} else {
			results = append(results, d)
		}
	}
	return results
}

func (g *Generator) generateHairLines3(faceContour []Point, numHairLines int) [][]Point {
	face := faceContour[:len(faceContour)-2]
	results := make([][]Point, 0, numHairLines)
	pickedIndices := make([]int, 0, numHairLines)
	for i := 0; i < numHairLines; i++ {
		pickedIndices = append(pickedIndices, int(math.Floor(g.random(10, 180))))
	}
	sortInts(pickedIndices)
	splitPoint := int(math.Floor(g.random(0, 200)))
	for i := 0; i < numHairLines; i++ {
		numHairPoints := 30 + int(math.Floor(g.random(-8, 8)))
		hairLine := make([]Point, 0, numHairPoints)
		indexOffset := pickedIndices[i]
		lower := g.random(1, 2.3)
		if g.r.Float64() > 0.9 {
			lower = g.random(0, 1)
		}
		reverse := -1
		if indexOffset > splitPoint {
			reverse = 1
		}
		for j := 0; j < numHairPoints; j++ {
			powerScale := g.random(0.1, 3)
			portion := (1-math.Pow(float64(j)/float64(numHairPoints), powerScale))*(1-lower) + lower
			p := face[(len(face)-(reverse*j*2+indexOffset))%len(face)]
			hairLine = append(hairLine, Point{X: p.X * portion, Y: p.Y})
		}
		results = append(results, computeBezierCurve(hairLine, numHairPoints))
	}
	return results
}
