//  draw.go -- canvas drawing operations

//#%#% a crude first hack.
//#%#% will need a good rewrite with error checking, clipping, etc.

package graphics

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var _ = fmt.Printf // enable debugging

//  VPainter.Reset() establishes or reestablishes initial conditions:
//		origin = center of canvas
//		current location = origin
//		orientation = towards top
//		drawing size = 1 pt
//		color = black
func (v *VPainter) Reset() *VPainter {
	v.Dx = float64(v.Width) / (2 * v.PixPerPt) // offset to origin
	v.Dy = float64(v.Height) / (2 * v.PixPerPt)
	v.Xloc = 0 // current location
	v.Yloc = 0
	v.Aim = -90                          // orientation = towards top
	v.Size = 1                           // drawing width = 1 pt
	v.VColor = NewColor(1, 1, 1, 1)      // color = white
	v.Rect(-v.Dx, -v.Dy, 2*v.Dx, 2*v.Dy) // clear the canvas
	v.VColor = NewColor(0, 0, 0, 1)      // color = black
	return v
}

//  VPainter.Goto(x,y,o) sets the current location to (x, y).
//  If o is a number (not nil) it sets the current orientation to that.
func (v *VPainter) Goto(x, y float64, o interface{}) *VPainter {
	v.Xloc = x
	v.Yloc = y
	if o != nil {
		v.Aim = o.(float64)
	}
	return v
}

//  VPainter.Forward(d) draws a line by moving the pen forward d units.
func (v *VPainter) Forward(d float64) *VPainter {
	s, c := math.Sincos(v.Aim * (math.Pi / 180))
	x := v.Xloc + d*c
	y := v.Yloc + d*s
	v.Line(v.Xloc, v.Yloc, x, y)
	v.Xloc = x
	v.Yloc = y
	return v
}

//  VPainter.Line(x1, y1, x2, y2) draws a line.
//  #%#% in a really dumb way. should stroke, not draw a zillion points.
func (v *VPainter) Line(x1, y1, x2, y2 float64) *VPainter {
	dx := x2 - x1
	dy := y2 - y1
	dmax := math.Max(math.Abs(dx), math.Abs(dy))
	n := int(math.Ceil(float64(v.PixPerPt) * dmax))
	dx /= float64(n)
	dy /= float64(n)
	for i := 0; i <= n; i++ {
		v.Point(x1, y1)
		x1 += dx
		y1 += dy
	}
	return v
}

//  VPainter.Point(x, y) draws a point (a disc based on pen size).
func (v *VPainter) Point(x, y float64) *VPainter {
	v.Disc(x, y, v.Size)
	return v
}

//  VPainter.Rect(x, y, w, h) draws a rectangle.
func (v *VPainter) Rect(x, y, w, h float64) *VPainter {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	x = x + v.Dx
	y = y + v.Dy
	r := image.Rect(v.ToPx(x), v.ToPx(y), v.ToPx(x+w), v.ToPx(y+h))
	draw.Draw(v.Canvas.Image, r,
		image.NewUniform(v.VColor), image.Point{}, draw.Over)
	return v
}

//  VPainter.Overlay(x, y, c) copies an image from canvas c.
//  The origin of canvas c is aligned with (x,y) of the destination.
//  #%#% Should allow subimage and scaling specification somehow.
func (v *VPainter) Overlay(x, y float64, c *VPainter) *VPainter {
	f := c.PixPerPt / v.PixPerPt        // scaling (sampling) factor
	w := float64(c.Width) / c.PixPerPt  // width in points
	h := float64(c.Height) / c.PixPerPt // height in points
	x0 := v.ToPx(x + v.Dx - c.Dx)
	y0 := v.ToPx(y + v.Dy - c.Dy)
	x1 := x0 + v.ToPx(w)
	y1 := y0 + v.ToPx(h)
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 > v.Width {
		x1 = v.Width
	}
	if y1 > v.Height {
		y1 = v.Height
	}
	//#%#% could add a fast path calling draw.Draw if densities agree
	for j := y0; j < y1; j++ {
		jj := int(f * float64(j-y0))
		for i := x0; i < x1; i++ {
			ii := int(f * float64(i-x0))
			v.Image.Set(i, j, c.Image.At(ii, jj))
		}
	}
	return v
}

//  VPainter.Text(x, y, s) draws a string of text characters.
func (v *VPainter) Text(x, y float64, s string) *VPainter {
	v.VFont.Typeset(v, v.ToPx(x+v.Dx), v.ToPx(y+v.Dy), s)
	return v
}

//  VPainter.Disc(x, y, d) draws a circle of diameter d at (x,y),
func (v *VPainter) Disc(x, y, d float64) *VPainter {
	xx := v.ToPx(x + v.Dx) // center quantized to raster coordinates
	yy := v.ToPx(y + v.Dy)
	r := v.PixPerPt * d / 2 // radius in pixels as floating value
	mask := &zircle{r}
	b := mask.Bounds()
	dstr := image.Rect(xx+b.Min.X, yy+b.Min.Y, xx+b.Max.X, yy+b.Max.Y)
	draw.DrawMask(v.Canvas.Image, dstr, image.NewUniform(v.VColor), image.ZP,
		mask, mask.Bounds().Min, draw.Over)
	return v
}

//  Adapted from "Drawing Through a Mask"
//  http://blog.golang.org/go-imagedraw-package

type zircle struct { // a circle of radius r at (0,0)
	r float64
}

func (z *zircle) ColorModel() color.Model {
	return color.AlphaModel
}

func (z *zircle) Bounds() image.Rectangle {
	i := int(math.Ceil(z.r))
	return image.Rect(-i, -i, i, i)
}

func (z *zircle) At(x, y int) color.Color {
	if float64(x*x+y*y) < z.r*z.r {
		return color.Alpha{255}
	} else {
		return color.Alpha{0}
	}
}
