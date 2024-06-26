package main

import (
	"flag"
	"fmt"
	"image/color"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// A type of function that passes in a pdf, that allows you to accomplish various functions included with the gofpdf library
// with easier syntax
type PDFOption func(*gofpdf.Fpdf)

// Create a fill color for a shape being drawn on the pdf
func FillColor(c color.RGBA) PDFOption {
	return func(pdf *gofpdf.Fpdf) {
		r, g, b := rgb(c)
		pdf.SetFillColor(r, g, b)
	}
}

// Return the red, green, and blue values for a provided color
func rgb(c color.RGBA) (int, int, int) {
	alpha := float64(c.A) / 255.0
	alphaWhite := int(255 * (1.0 - alpha))
	r := int(float64(c.R)*alpha) + alphaWhite
	g := int(float64(c.G)*alpha) + alphaWhite
	b := int(float64(c.B)*alpha) + alphaWhite
	return r, g, b
}

// Create a new pdf, with x,y used to track where we are in the pdf
type PDF struct {
	fpdf *gofpdf.Fpdf
	x, y float64
}

// For shifting where we are looking at in the pdf
func (p *PDF) Move(xDelta, yDelta float64) {
	p.x, p.y = p.x+xDelta, p.y+yDelta
	p.fpdf.MoveTo(p.x, p.y)
}

// For moving where we are in the pdf
func (p *PDF) MoveAbs(x, y float64) {
	p.x, p.y = x, y
	p.fpdf.MoveTo(p.x, p.y)
}

// For creating a new text at location p.x, p.y
func (p *PDF) Text(text string) {
	p.fpdf.Text(p.x, p.y, text)
}

// For creating a new Polygon
func (p *PDF) Polygon(pts []gofpdf.PointType, opts ...PDFOption) {
	for _, opt := range opts {
		opt(p.fpdf)
	}
	p.fpdf.Polygon(pts, "F")
}

func main() {
	//NOTE: Name of pdf is cert.pdf

	//When running the code, instead of calling go run . by itself, call go run . -name="YourName" to
	//dynamically insert your name into the certificate
	name := flag.String("name", "", "the name of the person who completed the course")
	flag.Parse()

	fpdf := gofpdf.New(gofpdf.OrientationLandscape, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	w, h := fpdf.GetPageSize()
	fpdf.AddPage()
	pdf := PDF{
		fpdf: fpdf,
	}
	primary := color.RGBA{103, 60, 79, 255}
	secondary := color.RGBA{103, 60, 79, 220}

	//Top and Bottom Banner
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{0, h / 9.0},
		{w - (w / 6.0), 0},
	}, FillColor(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{w / 6.0, 0},
		{w, 0},
		{w, h / 9.0},
	}, FillColor(primary))

	pdf.Polygon([]gofpdf.PointType{
		{w, h},
		{w, h - h/8.0},
		{w / 6, h},
	}, FillColor(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - h/8.0},
		{w - (w / 6), h},
	}, FillColor(primary))

	//Certificate of Completion Text
	fpdf.SetFont("times", "B", 50)
	fpdf.SetTextColor(50, 50, 50)
	pdf.MoveAbs(0, 100)
	_, lineHt := fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "Certificate of Completion", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*2.0)

	//Certificate Award Text
	fpdf.SetFont("arial", "", 28)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "This certificate is awarded to", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*2.0)

	//Name of person receiving certificate
	fpdf.SetFont("times", "B", 42)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, *name, gofpdf.AlignCenter)
	pdf.Move(0, lineHt*1.75)

	//Course Completion Text
	fpdf.SetFont("arial", "", 22)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt*1.5, "For successfully completing the Creating a PDF exercise in the course Gophercises programming course for budding Gophers (Go developers)", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*4.5)

	//Logo
	fpdf.ImageOptions("../Bisou.png", w/2.0-50.0, pdf.y, 100.0, 0, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")

	//Lines above "Date" and "Student Signature"
	pdf.Move(0, 65.0)
	fpdf.SetFillColor(100, 100, 100)
	fpdf.Rect(60, pdf.y, 240, 1.0, "F")
	fpdf.Rect(490.0, pdf.y, 240, 1.0, "F")

	//Date and Student Signature Text + Imported Student Signature and Dynamic Date
	fpdf.SetFont("arial", "", 12)
	pdf.Move(0, lineHt/1.5)
	fpdf.SetTextColor(100, 100, 100)
	pdf.MoveAbs(60+105.0, pdf.y)
	pdf.Text("Date")
	pdf.MoveAbs(490.0+60.0, pdf.y)
	pdf.Text("Student Signature")
	pdf.MoveAbs(60.0, pdf.y-lineHt/1.5)
	fpdf.SetFont("times", "", 22)
	_, lineHt = fpdf.GetFontSize()
	pdf.Move(0, -lineHt)
	fpdf.SetTextColor(50, 50, 50)
	year, mo, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d/%d/%d", mo, day, year)
	fpdf.CellFormat(240.0, lineHt, dateStr, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignCenter, false, 0, "")
	pdf.MoveAbs(490.0, pdf.y)

	sig, err := gofpdf.SVGBasicFileParse("../projPhase3certificate/sig.svg")
	if err != nil {
		panic(err)
	}
	pdf.Move(0, -(0.45*sig.Ht - lineHt))
	fpdf.SVGBasicWrite(&sig, 0.5)

	//drawGrid(fpdf)
	err = fpdf.OutputFileAndClose("cert.pdf")
	if err != nil {
		panic(err)
	}
}

// Draw a grid of equally spaced vertical and horizontal lines over the pdf
func drawGrid(pdf *gofpdf.Fpdf) {
	w, h := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)
	for x := 0.0; x < w; x = x + (w / 20.0) {
		pdf.SetTextColor(200, 200, 200)
		pdf.Line(x, 0, x, h)
		_, lineHt := pdf.GetFontSize()
		pdf.Text(x, lineHt, fmt.Sprintf("%d", int(x)))
	}
	for y := 0.0; y < h; y = y + (w / 20.0) {
		pdf.SetTextColor(80, 80, 80)
		pdf.Line(0, y, w, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}
