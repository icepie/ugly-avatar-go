# uglyavatar

Go library port of the Vue ugly avatar generator in `../`.

## Install

```bash
go get github.com/txstc55/uglyavatar
```

## Usage

```go
package main

import (
	"os"

	"github.com/txstc55/uglyavatar"
)

func main() {
	avatar := uglyavatar.GenerateWithSeed(42)
	_ = os.WriteFile("avatar.svg", []byte(avatar.SVG()), 0o644)
	_ = avatar.SavePNG("avatar.png", 500, 500)
}
```

## API

```go
avatar := uglyavatar.Generate()
avatar = uglyavatar.GenerateWithSeed(42)

svg := avatar.SVG()
svg = avatar.SVGWithSize(500, 500)

pngBytes, err := avatar.PNG(500, 500)
err = avatar.WritePNG(writer, 500, 500)
err = avatar.SavePNG("avatar.png", 500, 500)
img, err := avatar.RenderImage(500, 500)
```

The package exposes the generated geometry through `Avatar`, so callers can use the built-in SVG/PNG renderers or draw the face with their own renderer. Use `NewWithSeed` or `GenerateWithSeed` when repeatable output matters.

This is a port of the original non-commercial project by Xuan Tang. Check the parent repository license before redistributing or using it commercially.
