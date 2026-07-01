package uglyavatar

func (g *Generator) generateFaceContourPoints(numPoints int) ([]Point, float64, float64, Point) {
	faceSizeX0 := g.random(50, 100)
	faceSizeY0 := g.random(70, 100)
	faceSizeY1 := g.random(50, 80)
	faceSizeX1 := g.random(70, 100)
	faceK0 := g.random(0.001, 0.005)
	if !g.chance(0.5) {
		faceK0 *= -1
	}
	faceK1 := g.random(0.001, 0.005)
	if !g.chance(0.5) {
		faceK1 *= -1
	}
	face0TranslateX := g.random(-5, 5)
	face0TranslateY := g.random(-15, 15)
	face1TranslateY := g.random(-5, 5)
	face1TranslateX := g.random(-5, 25)

	var results0 []Point
	if g.r.Float64() > 0.1 {
		results0 = g.eggShapePoints(faceSizeX0, faceSizeY0, faceK0, numPoints)
	} else {
		results0 = g.rectangularFaceContourPoints(faceSizeX0, faceSizeY0, numPoints)
	}
	var results1 []Point
	if g.r.Float64() > 0.3 {
		results1 = g.eggShapePoints(faceSizeX1, faceSizeY1, faceK1, numPoints)
	} else {
		results1 = g.rectangularFaceContourPoints(faceSizeX1, faceSizeY1, numPoints)
	}

	for i := range results0 {
		results0[i].X += face0TranslateX
		results0[i].Y += face0TranslateY
		results1[i].X += face1TranslateX
		results1[i].Y += face1TranslateY
	}

	result := make([]Point, 0, len(results0)+2)
	var center Point
	quarter := len(results0) / 4
	for i := range results0 {
		other := results1[(i+quarter)%len(results0)]
		p := Point{
			X: results0[i].X*0.7 + other.Y*0.3,
			Y: results0[i].Y*0.7 - other.X*0.3,
		}
		result = append(result, p)
		center.X += p.X
		center.Y += p.Y
	}
	center.X /= float64(len(result))
	center.Y /= float64(len(result))
	for i := range result {
		result[i].X -= center.X
		result[i].Y -= center.Y
	}

	width := result[0].X - result[len(result)/2].X
	height := result[len(result)/4].Y - result[(len(result)*3)/4].Y
	result = append(result, result[0], result[1])
	return result, width, height, Point{}
}
