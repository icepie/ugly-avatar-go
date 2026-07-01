package uglyavatar

import (
	"math"
	"sort"
)

func (g *Generator) random(min, max float64) float64 {
	return g.r.Float64()*(max-min) + min
}

func (g *Generator) chance(p float64) bool {
	return g.r.Float64() < p
}

func cubicBezier(p0, p1, p2, p3 Point, t float64) Point {
	mt := 1 - t
	return Point{
		X: mt*mt*mt*p0.X + 3*mt*mt*t*p1.X + 3*mt*t*t*p2.X + t*t*t*p3.X,
		Y: mt*mt*mt*p0.Y + 3*mt*mt*t*p1.Y + 3*mt*t*t*p2.Y + t*t*t*p3.Y,
	}
}

func factorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	return float64(n) * factorial(n-1)
}

func binomialCoefficient(n, k int) float64 {
	return factorial(n) / (factorial(k) * factorial(n-k))
}

func calculateBezierPoint(t float64, controlPoints []Point) Point {
	var p Point
	n := len(controlPoints) - 1
	for i := 0; i <= n; i++ {
		coef := binomialCoefficient(n, i)
		a := math.Pow(1-t, float64(n-i))
		b := math.Pow(t, float64(i))
		p.X += coef * a * b * controlPoints[i].X
		p.Y += coef * a * b * controlPoints[i].Y
	}
	return p
}

func computeBezierCurve(controlPoints []Point, numberOfPoints int) []Point {
	curve := make([]Point, 0, numberOfPoints+1)
	for i := 0; i <= numberOfPoints; i++ {
		curve = append(curve, calculateBezierPoint(float64(i)/float64(numberOfPoints), controlPoints))
	}
	return curve
}

func reversePoints(points []Point) []Point {
	out := make([]Point, len(points))
	for i := range points {
		out[i] = points[len(points)-1-i]
	}
	return out
}

func copyPoints(points []Point) []Point {
	out := make([]Point, len(points))
	copy(out, points)
	return out
}

func (g *Generator) eggShapePoints(a, b, k float64, segmentPoints int) []Point {
	result := make([]Point, 0, segmentPoints*4)
	for i := 0; i < segmentPoints; i++ {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/1.1/float64(segmentPoints), math.Pi/1.1/float64(segmentPoints))
		y := math.Sin(degree) * b
		x := math.Sqrt(((1-(y*y)/(b*b))/(1+k*y))*a*a) + g.random(-a/200, a/200)
		result = append(result, Point{X: x, Y: y})
	}
	for i := segmentPoints; i > 0; i-- {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/1.1/float64(segmentPoints), math.Pi/1.1/float64(segmentPoints))
		y := math.Sin(degree) * b
		x := -math.Sqrt(((1-(y*y)/(b*b))/(1+k*y))*a*a) + g.random(-a/200, a/200)
		result = append(result, Point{X: x, Y: y})
	}
	for i := 0; i < segmentPoints; i++ {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/1.1/float64(segmentPoints), math.Pi/1.1/float64(segmentPoints))
		y := -math.Sin(degree) * b
		x := -math.Sqrt(((1-(y*y)/(b*b))/(1+k*y))*a*a) + g.random(-a/200, a/200)
		result = append(result, Point{X: x, Y: y})
	}
	for i := segmentPoints; i > 0; i-- {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/1.1/float64(segmentPoints), math.Pi/1.1/float64(segmentPoints))
		y := -math.Sin(degree) * b
		x := math.Sqrt(((1-(y*y)/(b*b))/(1+k*y))*a*a) + g.random(-a/200, a/200)
		result = append(result, Point{X: x, Y: y})
	}
	return result
}

func rectangularIntersection(radian, a, b float64) Point {
	if radian < 0 {
		radian = 0
	}
	if radian > math.Pi/2 {
		radian = math.Pi / 2
	}
	if math.Abs(radian-math.Pi/2) < 0.0001 {
		return Point{X: 0, Y: b}
	}
	m := math.Tan(radian)
	y := m * a
	if y < b {
		return Point{X: a, Y: y}
	}
	return Point{X: b / m, Y: b}
}

func (g *Generator) rectangularFaceContourPoints(a, b float64, segmentPoints int) []Point {
	result := make([]Point, 0, segmentPoints*4)
	for i := 0; i < segmentPoints; i++ {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/11/float64(segmentPoints), math.Pi/11/float64(segmentPoints))
		result = append(result, rectangularIntersection(degree, a, b))
	}
	for i := segmentPoints; i > 0; i-- {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/11/float64(segmentPoints), math.Pi/11/float64(segmentPoints))
		p := rectangularIntersection(degree, a, b)
		result = append(result, Point{X: -p.X, Y: p.Y})
	}
	for i := 0; i < segmentPoints; i++ {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/11/float64(segmentPoints), math.Pi/11/float64(segmentPoints))
		p := rectangularIntersection(degree, a, b)
		result = append(result, Point{X: -p.X, Y: -p.Y})
	}
	for i := segmentPoints; i > 0; i-- {
		degree := (math.Pi/2/float64(segmentPoints))*float64(i) +
			g.random(-math.Pi/11/float64(segmentPoints), math.Pi/11/float64(segmentPoints))
		p := rectangularIntersection(degree, a, b)
		result = append(result, Point{X: p.X, Y: -p.Y})
	}
	return result
}

func sortInts(values []int) {
	sort.Ints(values)
}
