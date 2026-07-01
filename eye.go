package uglyavatar

type eyeParams struct {
	heightUpper           float64
	heightLower           float64
	p0UpperRandX          float64
	p3UpperRandX          float64
	p0UpperRandY          float64
	p3UpperRandY          float64
	offsetUpperLeftRandY  float64
	offsetUpperRightRandY float64
	eyeTrueWidth          float64
	offsetUpperLeftX      float64
	offsetUpperRightX     float64
	offsetUpperLeftY      float64
	offsetUpperRightY     float64
	offsetLowerLeftX      float64
	offsetLowerRightX     float64
	offsetLowerLeftY      float64
	offsetLowerRightY     float64
	leftConverge0         float64
	rightConverge0        float64
	leftConverge1         float64
	rightConverge1        float64
}

type eyeShape struct {
	upper  []Point
	lower  []Point
	center Point
}

func (g *Generator) generateEyeParameters(width float64) eyeParams {
	heightUpper := g.r.Float64() * width / 1.2
	heightLower := g.r.Float64() * width / 1.2
	p0UpperRandX := g.r.Float64()*0.4 - 0.2
	p3UpperRandX := g.r.Float64()*0.4 - 0.2
	p0UpperRandY := g.r.Float64()*0.4 - 0.2
	p3UpperRandY := g.r.Float64()*0.4 - 0.2
	offsetUpperLeftRandY := g.r.Float64()
	offsetUpperRightRandY := g.r.Float64()

	p0Upper := Point{X: -width/2 + p0UpperRandX*width/16, Y: p0UpperRandY * heightUpper / 16}
	p3Upper := Point{X: width/2 + p3UpperRandX*width/16, Y: p3UpperRandY * heightUpper / 16}
	eyeTrueWidth := p3Upper.X - p0Upper.X
	offsetUpperLeftX := g.random(-eyeTrueWidth/10, eyeTrueWidth/2.3)
	offsetUpperRightX := g.random(-eyeTrueWidth/10, eyeTrueWidth/2.3)
	offsetUpperLeftY := offsetUpperLeftRandY * heightUpper
	offsetUpperRightY := offsetUpperRightRandY * heightUpper

	return eyeParams{
		heightUpper: heightUpper, heightLower: heightLower,
		p0UpperRandX: p0UpperRandX, p3UpperRandX: p3UpperRandX,
		p0UpperRandY: p0UpperRandY, p3UpperRandY: p3UpperRandY,
		offsetUpperLeftRandY: offsetUpperLeftRandY, offsetUpperRightRandY: offsetUpperRightRandY,
		eyeTrueWidth:      eyeTrueWidth,
		offsetUpperLeftX:  offsetUpperLeftX,
		offsetUpperRightX: offsetUpperRightX,
		offsetUpperLeftY:  offsetUpperLeftY,
		offsetUpperRightY: offsetUpperRightY,
		offsetLowerLeftX:  g.random(offsetUpperLeftX, eyeTrueWidth/2.1),
		offsetLowerRightX: g.random(offsetUpperRightX, eyeTrueWidth/2.1),
		offsetLowerLeftY:  g.random(-offsetUpperLeftY+5, heightLower),
		offsetLowerRightY: g.random(-offsetUpperRightY+5, heightLower),
		leftConverge0:     g.r.Float64(),
		rightConverge0:    g.r.Float64(),
		leftConverge1:     g.r.Float64(),
		rightConverge1:    g.r.Float64(),
	}
}

func (g *Generator) perturbEyeParams(p eyeParams) eyeParams {
	p.heightUpper += g.random(-p.heightUpper/2, p.heightUpper/2)
	p.heightLower += g.random(-p.heightLower/2, p.heightLower/2)
	p.p0UpperRandX += g.random(-p.p0UpperRandX/2, p.p0UpperRandX/2)
	p.p3UpperRandX += g.random(-p.p3UpperRandX/2, p.p3UpperRandX/2)
	p.p0UpperRandY += g.random(-p.p0UpperRandY/2, p.p0UpperRandY/2)
	p.p3UpperRandY += g.random(-p.p3UpperRandY/2, p.p3UpperRandY/2)
	p.offsetUpperLeftRandY += g.random(-p.offsetUpperLeftRandY/2, p.offsetUpperLeftRandY/2)
	p.offsetUpperRightRandY += g.random(-p.offsetUpperRightRandY/2, p.offsetUpperRightRandY/2)
	p.eyeTrueWidth += g.random(-p.eyeTrueWidth/2, p.eyeTrueWidth/2)
	p.offsetUpperLeftX += g.random(-p.offsetUpperLeftX/2, p.offsetUpperLeftX/2)
	p.offsetUpperRightX += g.random(-p.offsetUpperRightX/2, p.offsetUpperRightX/2)
	p.offsetUpperLeftY += g.random(-p.offsetUpperLeftY/2, p.offsetUpperLeftY/2)
	p.offsetUpperRightY += g.random(-p.offsetUpperRightY/2, p.offsetUpperRightY/2)
	p.offsetLowerLeftX += g.random(-p.offsetLowerLeftX/2, p.offsetLowerLeftX/2)
	p.offsetLowerRightX += g.random(-p.offsetLowerRightX/2, p.offsetLowerRightX/2)
	p.offsetLowerLeftY += g.random(-p.offsetLowerLeftY/2, p.offsetLowerLeftY/2)
	p.offsetLowerRightY += g.random(-p.offsetLowerRightY/2, p.offsetLowerRightY/2)
	p.leftConverge0 += g.random(-p.leftConverge0/2, p.leftConverge0/2)
	p.rightConverge0 += g.random(-p.rightConverge0/2, p.rightConverge0/2)
	p.leftConverge1 += g.random(-p.leftConverge1/2, p.leftConverge1/2)
	p.rightConverge1 += g.random(-p.rightConverge1/2, p.rightConverge1/2)
	return p
}

func generateEyePoints(p eyeParams, width float64) eyeShape {
	p0Upper := Point{X: -width/2 + p.p0UpperRandX*width/16, Y: p.p0UpperRandY * p.heightUpper / 16}
	p3Upper := Point{X: width/2 + p.p3UpperRandX*width/16, Y: p.p3UpperRandY * p.heightUpper / 16}
	p0Lower := p0Upper
	p3Lower := p3Upper

	p1Upper := Point{X: p0Upper.X + p.offsetUpperLeftX, Y: p0Upper.Y + p.offsetUpperLeftY}
	p2Upper := Point{X: p3Upper.X - p.offsetUpperRightX, Y: p3Upper.Y + p.offsetUpperRightY}
	p1Lower := Point{X: p0Lower.X + p.offsetLowerLeftX, Y: p0Lower.Y - p.offsetLowerLeftY}
	p2Lower := Point{X: p3Lower.X - p.offsetLowerRightX, Y: p3Lower.Y - p.offsetLowerRightY}

	upper := make([]Point, 100)
	upperLeftControl := make([]Point, 100)
	upperRightControl := make([]Point, 100)
	upperLeftControlPoint := Point{
		X: p0Upper.X*(1-p.leftConverge0) + p1Lower.X*p.leftConverge0,
		Y: p0Upper.Y*(1-p.leftConverge0) + p1Lower.Y*p.leftConverge0,
	}
	upperRightControlPoint := Point{
		X: p3Upper.X*(1-p.rightConverge0) + p2Lower.X*p.rightConverge0,
		Y: p3Upper.Y*(1-p.rightConverge0) + p2Lower.Y*p.rightConverge0,
	}
	for t := 0; t < 100; t++ {
		tf := float64(t) / 100
		upper[t] = cubicBezier(p0Upper, p1Upper, p2Upper, p3Upper, tf)
		upperLeftControl[t] = cubicBezier(upperLeftControlPoint, p0Upper, p1Upper, p2Upper, tf)
		upperRightControl[t] = cubicBezier(p1Upper, p2Upper, p3Upper, upperRightControlPoint, tf)
	}
	for i := 0; i < 75; i++ {
		weight := pow2((75 - float64(i)) / 75)
		upper[i] = Point{
			X: upper[i].X*(1-weight) + upperLeftControl[i+25].X*weight,
			Y: upper[i].Y*(1-weight) + upperLeftControl[i+25].Y*weight,
		}
		upper[i+25] = Point{
			X: upper[i+25].X*weight + upperRightControl[i].X*(1-weight),
			Y: upper[i+25].Y*weight + upperRightControl[i].Y*(1-weight),
		}
	}

	lower := make([]Point, 100)
	lowerLeftControl := make([]Point, 100)
	lowerRightControl := make([]Point, 100)
	lowerLeftControlPoint := Point{
		X: p0Lower.X*(1-p.leftConverge0) + p1Upper.X*p.leftConverge0,
		Y: p0Lower.Y*(1-p.leftConverge0) + p1Upper.Y*p.leftConverge0,
	}
	lowerRightControlPoint := Point{
		X: p3Lower.X*(1-p.rightConverge1) + p2Upper.X*p.rightConverge1,
		Y: p3Lower.Y*(1-p.rightConverge1) + p2Upper.Y*p.rightConverge1,
	}
	for t := 0; t < 100; t++ {
		tf := float64(t) / 100
		lower[t] = cubicBezier(p0Lower, p1Lower, p2Lower, p3Lower, tf)
		lowerLeftControl[t] = cubicBezier(lowerLeftControlPoint, p0Lower, p1Lower, p2Lower, tf)
		lowerRightControl[t] = cubicBezier(p1Lower, p2Lower, p3Lower, lowerRightControlPoint, tf)
	}
	for i := 0; i < 75; i++ {
		weight := pow2((75 - float64(i)) / 75)
		lower[i] = Point{
			X: lower[i].X*(1-weight) + lowerLeftControl[i+25].X*weight,
			Y: lower[i].Y*(1-weight) + lowerLeftControl[i+25].Y*weight,
		}
		lower[i+25] = Point{
			X: lower[i+25].X*weight + lowerRightControl[i].X*(1-weight),
			Y: lower[i+25].Y*weight + lowerRightControl[i].Y*(1-weight),
		}
	}
	for i := 0; i < 100; i++ {
		lower[i].Y = -lower[i].Y
		upper[i].Y = -upper[i].Y
	}

	center := Point{X: upper[50].X/2 + lower[50].X/2, Y: upper[50].Y/2 + lower[50].Y/2}
	for i := 0; i < 100; i++ {
		lower[i].X -= center.X
		lower[i].Y -= center.Y
		upper[i].X -= center.X
		upper[i].Y -= center.Y
	}
	return eyeShape{upper: upper, lower: lower}
}

func pow2(v float64) float64 {
	return v * v
}

func (g *Generator) generateBothEyes(width float64) (eyeShape, eyeShape) {
	leftParams := g.generateEyeParameters(width)
	rightParams := g.perturbEyeParams(leftParams)
	left := generateEyePoints(leftParams, width)
	right := generateEyePoints(rightParams, width)
	for i := range left.upper {
		left.upper[i].X = -left.upper[i].X
	}
	for i := range left.lower {
		left.lower[i].X = -left.lower[i].X
	}
	return left, right
}
