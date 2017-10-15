package gg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

// PlotXY draws line graphs for given X-Y values
//
//                |←——————————————————— W ————————————————————→|
//                       |←——————————— ww ——————————→|
//                            |←—————— w ——————→|
//              (0,0)
//   ———————————  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ——→ xScr
//    ↑           ┃      |             '                       ┃
//    |           ┃   (x0,y0)          ' TR                    ┃
//    | ————————  ┃      ┌───────────────────────────┐ ——————— ┃ ——→ x
//    |   ↑       ┃      │ (p0,q0)     ' DV          │    ↑    ┃
//    |   |  ———  ┃      │    o─────────────────┐    │    |    ┃ ——→ p
//    |   |   ↑   ┃      │    │                 │    │    |    ┃
//    |           ┃ —LR— │-DH-│                 │-DH-│-RR-|-RL-┃
//       hh   h   ┃      │    │                 │    │    |    ┃
//    H           ┃      │    ↑ yReal           │    │    |    ┃
//        |   ↓   ┃      │    │                 │    │    |    ┃
//    |   |  ———  ┃      │    ●─────────────────o    │    |    ┃ ——→ xReal
//    |   ↓       ┃      │             ' DV  (pf,qf) │    ↓    ┃
//    |  ———————  ┃ ———— └───────────────────────────┘ ——————— ┃
//    |           ┃                    ' BR       (xf,yf)      ┃
//    ↓           ┃                    ' BL                    ┃
//   ———————————  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
//                |      |    |                              (W,H)
//                ↓      ↓    ↓
//               yScr    y    q
//
//    W: figure width          ww: plot area width     w: x-y area width
//    H: figure height         hh: plot area height    h: x-y area height
//
//    LR : Left Ruler          TR : Top Ruler          xScr ,yScr : "screen" coordinates
//    BR : Bottom Ruler        RR : Right Ruler        x    ,y    : plot area coordinates
//    BL : Bottom Legend       RL : Right Legend       xReal,yReal: x-y (real data) coords
//    DH : Delta Horizontal    DV : Delta Vertical
//
type PlotXY struct {

	// title and labels
	Title  string // Title
	Xlabel string // XLabel
	Ylabel string // YLabel

	// options
	EqualScale  bool // equal scale factors
	DrawGrid    bool // draw grid
	DrawBorders bool // draw borders
	DeltaH      int  // DH increment
	DeltaV      int  // DV increment

	// legend
	LegendOn    bool // legend is on
	LegAtBottom bool // legend at bottom
	LegLineLen  int  // length of legend line indicator
	LegGap      int  // legend: gap between icons
	LegTxtGap   int  // legend: gap between line and text
	LegNrow     int  // legend: number of rows

	// ticks
	NumTicksX      int    // number of x-ticks
	NumTicksY      int    // number of y-ticks
	TicksFormat    string // format of ticks numbers
	TicksNumDigits int    // number of digits of ticks
	TicksLength    int    // length of tick lines

	// styles
	StyleFG *PlotArgs // style: foreground
	StyleFR *PlotArgs // style: frame
	StylePL *PlotArgs // style: plot area
	StyleGD *PlotArgs // style: grid
	StyleBR *PlotArgs // style: bottom ruler
	StyleLR *PlotArgs // style: left ruler
	StyleTR *PlotArgs // style: top ruler
	StyleRR *PlotArgs // style: right ruler

	// font typefaces
	FontTitle  *truetype.Font // font for title text
	FontTicks  *truetype.Font // font for ticks text
	FontLabel  *truetype.Font // font for x-y labels
	FontLegend *truetype.Font // font for legend text

	// font sizes
	FsizeTitle  int // font size of title text
	FsizeTicks  int // font size of ticks text
	FsizeLabels int // font size of x-y labels
	FsizeLegend int // font size of legend text

	// curves and ticks
	dataX  [][]float64 // all curves x-data
	dataY  [][]float64 // all curves y-data
	curves []*PlotArgs // all curves properties
	xticks []float64   // x-ticks
	yticks []float64   // y-ticks

	// scale and positions
	sfx float64 // x scale factor
	sfy float64 // y scale factor
	p0  int     // x-origin of plotting area
	q0  int     // y-origin of plotting area
	pf  int     // x-max of plotting area
	qf  int     // y-max of plotting area

	// limits
	xmin      float64 // minimum x value (real coordinates)
	ymin      float64 // minimum y value (real coordinates)
	xmax      float64 // maximum x value (real coordinates)
	ymax      float64 // maximum y value (real coordinates)
	xminFix   float64 // fixed minimum x value (real coordinates)
	yminFix   float64 // fixed minimum y value (real coordinates)
	xmaxFix   float64 // fixed maximum x value (real coordinates)
	ymaxFix   float64 // fixed maximum y value (real coordinates)
	xminFixOn bool    // use or not xminFix
	xmaxFixOn bool    // use or not xmaxFix
	yminFixOn bool    // use or not yminFix
	ymaxFixOn bool    // use or not ymaxFix

	// constants
	cteEps   float64 // constant machine eps
	cteSqEps float64 // constant sqrt(eps)
	cteMin   float64 // constant min float
	cteMax   float64 // constant max float
}

// NewPlotXY creates a new PlotXY object
func NewPlotXY() (o *PlotXY) {

	// title and labels
	o = new(PlotXY)
	o.Title = "Plotting with PlotXY"
	o.Xlabel = "x"
	o.Ylabel = "y"

	// options
	o.EqualScale = false
	o.DrawGrid = true
	o.DrawBorders = true
	o.DeltaH = 8
	o.DeltaV = 8

	// legend
	o.LegendOn = true
	o.LegAtBottom = true
	o.LegLineLen = 30
	o.LegGap = 10
	o.LegTxtGap = 4
	o.LegNrow = 1

	// ticks
	o.NumTicksX = 10
	o.NumTicksY = 10
	o.TicksFormat = "%g"
	o.TicksNumDigits = 7
	o.TicksLength = 6

	// styles
	o.StyleFG = &PlotArgs{C: "#000"}    // foreground
	o.StyleFR = &PlotArgs{C: "#fff"}    // frame
	o.StylePL = &PlotArgs{C: "#faf7ec"} // plot area
	o.StyleGD = &PlotArgs{C: "#b7b7b7"} // grid
	o.StyleBR = &PlotArgs{C: "#e0e0e0"} // bottom ruler
	o.StyleLR = &PlotArgs{C: "#e0e0e0"} // left ruler
	o.StyleTR = &PlotArgs{C: "#e0e0e0"} // top ruler
	o.StyleRR = &PlotArgs{C: "#e0e0e0"} // right ruler

	// fonts
	var err1, err2, err3, err4 error
	o.FontTitle, err1 = truetype.Parse(goregular.TTF)
	o.FontTicks, err2 = truetype.Parse(gomono.TTF)
	o.FontLabel, err3 = truetype.Parse(gomono.TTF)
	o.FontLegend, err4 = truetype.Parse(gomono.TTF)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("cannot load fonts")
	}

	// font sizes
	o.FsizeTitle = 16
	o.FsizeTicks = 12
	o.FsizeLabels = 14
	o.FsizeLegend = 12

	// constants
	o.cteEps = math.Nextafter(1.0, 2.0) - 1.0
	o.cteSqEps = math.Sqrt(o.cteEps)
	o.cteMin = math.SmallestNonzeroFloat64
	o.cteMax = math.MaxFloat64
	return
}

// AddCurve adds curve to graph and returns the new curve properties
func (o *PlotXY) AddCurve(name string, x, y []float64) (curve *PlotArgs) {

	// check
	if len(x) != len(y) {
		panic("lengths of x and y must be the same")
	}

	// add (real) coordinates and curve properties
	ncurves := len(o.curves)
	curve = &PlotArgs{L: name, C: GetColor(ncurves, 0)}
	o.dataX = append(o.dataX, x)
	o.dataY = append(o.dataY, y)
	o.curves = append(o.curves, curve)

	// find limits
	if ncurves == 0 { // first curve
		if len(x) > 1 {
			o.xmin = x[0]
			o.xmax = x[0]
			o.ymin = y[0]
			o.ymax = y[0]
		} else {
			o.xmin = 0.0
			o.xmax = 1.0
			o.ymin = 0.0
			o.ymax = 1.0
		}
	}
	for i := 1; i < len(x); i++ {
		o.xmin = min(o.xmin, x[i])
		o.xmax = max(o.xmax, x[i])
		o.ymin = min(o.ymin, y[i])
		o.ymax = max(o.ymax, y[i])
	}
	return
}

// SetMinX sets minimum x value. Use fixed=true to prevent automatic updates
func (o *PlotXY) SetMinX(value float64, fixed bool) {
	o.xminFix = value
	o.xminFixOn = fixed
}

// SetMinY sets minimum y value. Use fixed=true to prevent automatic updates
func (o *PlotXY) SetMinY(value float64, fixed bool) {
	o.yminFix = value
	o.yminFixOn = fixed
}

// SetMaxX sets maximum x value. Use fixed=true to prevent automatic updates
func (o *PlotXY) SetMaxX(value float64, fixed bool) {
	o.xmaxFix = value
	o.xmaxFixOn = fixed
}

// SetMaxY sets maximum y value. Use fixed=true to prevent automatic updates
func (o *PlotXY) SetMaxY(value float64, fixed bool) {
	o.ymaxFix = value
	o.ymaxFixOn = fixed
}

// Render draws PlotXY
func (o *PlotXY) Render(dc *Context) {

	// check number of curves
	ncurves := len(o.curves)
	if ncurves < 1 {
		return
	}

	// x-y limits
	if o.xminFixOn {
		o.xmin = o.xminFix
	}
	if o.xmaxFixOn {
		o.xmax = o.xmaxFix
	}
	if o.yminFixOn {
		o.ymin = o.yminFix
	}
	if o.ymaxFixOn {
		o.ymax = o.ymaxFix
	}
	if math.Abs(o.xmax-o.xmin) <= o.cteEps {
		o.xmin = o.xmin - 1
		o.xmax = o.xmax + 1
	}
	if math.Abs(o.ymax-o.ymin) <= o.cteEps {
		o.ymin = o.ymin - 1
		o.ymax = o.ymax + 1
	}

	// ticks values
	bnumtck := o.NumTicksX
	lnumtck := o.NumTicksY
	if math.Abs(o.xmax-o.xmin) <= o.cteEps {
		bnumtck = 3
	}
	if math.Abs(o.ymax-o.ymin) <= o.cteEps {
		lnumtck = 3
	}
	o.xticks = o.pretty(o.xmin, o.xmax, bnumtck)
	o.yticks = o.pretty(o.ymin, o.ymax, lnumtck)

	// bottom: tick text height
	txt := o.fmtNum(o.TicksFormat, o.TicksNumDigits, o.xticks[0])
	_, tickH := o.measureTxt(dc, o.FontTicks, o.FsizeTicks, txt)

	// bottom: x-label text height
	_, xlblH := o.measureTxt(dc, o.FontLabel, o.FsizeLabels, o.Xlabel)

	// left: tick text width
	tickW := 0
	for _, value := range o.yticks {
		txt = o.fmtNum(o.TicksFormat, o.TicksNumDigits, value)
		tw, _ := o.measureTxt(dc, o.FontTicks, o.FsizeTicks, txt)
		tickW = imax(tickW, tw)
	}

	// left: y-label text width
	ylblW, _ := o.measureTxt(dc, o.FontLabel, o.FsizeLabels, o.Ylabel)

	// legend dimensions
	legTxtW := 0 // legend txt width
	legTxtH := 0 // legend txt height
	legH := 0    // legend total height
	if o.LegendOn {
		o.setFont(dc, o.FontLegend, o.FsizeLegend, "")
		for _, curve := range o.curves {
			tw, th := dc.MeasureString(curve.L)
			legTxtW = imax(legTxtW, int(tw))
			legTxtH = imax(legTxtH, int(th))
			legTxtH = imax(legTxtH, curve.markerSize())
		}
		legH = (2 + legTxtH) * o.LegNrow
	}

	// height of scales
	xscaleH := o.TicksLength + tickH + xlblH + 2 // height of x-scale (ticks + label)
	yscaleW := o.TicksLength + tickW + ylblW + 2 // width of y-scale (ticks + label)

	// auxiliary variables
	LR := 0 // Left ruler thickness (screen coordinates)
	RR := 6 // Right ruler thickness (screen coordinates)
	BR := 0 // Bottom ruler thickness (screen coordinates)
	TR := 6 // Top ruler thickness (screen coordinates)
	RL := 0 // right legend thickness
	BL := 0 // bottom legend thickness
	BR = imax(BR, xscaleH)
	LR = imax(LR, yscaleW)
	if o.LegAtBottom {
		BR += legH
	} else {
		RR = o.LegLineLen + o.LegTxtGap + o.LegGap + legTxtW + 2 // width of legend "icon"
	}

	// height of title
	if o.Title != "" {
		_, th := o.measureTxt(dc, o.FontTitle, o.FsizeTitle, o.Title)
		TR = th + 12
	}

	// derived variables
	W := dc.Width()
	H := dc.Height()
	ww := imax(1, W-(LR+RR+RL))
	hh := imax(1, H-(TR+BR+BL))
	w := imax(1, ww-2*o.DeltaH)
	h := imax(1, hh-2*o.DeltaV)
	x0 := LR
	y0 := TR
	o.p0 = x0 + o.DeltaH
	o.q0 = y0 + o.DeltaV
	xf := x0 + ww
	yf := y0 + hh
	o.pf = o.p0 + w
	o.qf = o.q0 + h

	// scaling factors
	o.sfx = float64(w) / (o.xmax - o.xmin)
	o.sfy = float64(h) / (o.ymax - o.ymin)
	if o.sfx <= o.cteEps {
		o.sfx = 1.0
	}
	if o.sfy <= o.cteEps {
		o.sfy = 1.0
	}
	if o.EqualScale {
		sf := o.sfx
		if o.sfx > o.sfy {
			sf = o.sfy
		}
		o.sfx = sf
		o.sfy = sf
	}

	// draw background of plot-area
	o.StylePL.Rect(dc, false, false, x0, y0, ww, hh)

	// draw grid
	if o.DrawGrid {

		// vertical lines
		for i := 0; i < len(o.xticks); i++ {
			x := o.xScr(o.xticks[i])
			if x >= x0 && x <= xf {
				o.StyleGD.Line(dc, false, true, x, y0, x, yf)
			}
		}

		// horizontal lines
		for i := 0; i < len(o.yticks); i++ {
			y := o.yScr(o.yticks[i])
			if y >= y0 && y <= yf {
				o.StyleGD.Line(dc, false, true, x0, y, xf, y)
			}
		}
	}

	// draw curves
	for k, curve := range o.curves {

		// draw markers
		if curve.M != "" {
			idx := 0
			for i := 0; i < len(o.dataX[k]); i++ {
				if i >= idx {
					curve.DrawMarker(dc, o.xScr(o.dataX[k][i]), o.yScr(o.dataY[k][i]))
					idx += curve.Me
				}
			}
		}

		// draw lines
		if curve.Ls != "none" {
			if len(o.dataX[k]) > 1 {
				curve.Activate(dc, false, true)
				dc.MoveTo(float64(o.xScr(o.dataX[k][0])), float64(o.yScr(o.dataY[k][0])))
				for i := 0; i < len(o.dataX[k]); i++ {
					dc.LineTo(float64(o.xScr(o.dataX[k][i])), float64(o.yScr(o.dataY[k][i])))
				}
				dc.Stroke()
			}
		}
	}

	// compute legend data, where the legend "icon" dimensions are:
	//
	//        |← legLineLen →|← labelLen →|
	//   [gap][      line    |    txt     ]     example:     ——x——Curve1
	//
	hei := 2 + legTxtH          // icon height
	lll := o.LegLineLen         // length of legend line
	hll := lll / 2              // half length of legend line
	xl := o.LegGap              // initial x-coord on icon line
	yl := yf + xscaleH + hei/2  // initial y-coord on icon line
	col := 0                    // column number
	ncol := ncurves / o.LegNrow // number of columns
	if ncurves%o.LegNrow > 0 {
		ncol++
	}
	if !o.LegAtBottom {
		xl = xf + o.LegGap
		yl = TR + o.LegGap
		ncol = 1
	}

	// bottom ruler
	if BR > 1 {

		// clear background
		o.StyleBR.Rect(dc, false, false, 0, yf, W, imax(1, H-hh-TR))

		// draw ticks and text
		for _, x := range o.xticks {
			xi := o.xScr(x)
			if xi >= x0 && xi <= xf {
				o.StyleFG.Line(dc, false, true, xi, yf, xi, yf+o.TicksLength)
				txt = o.fmtNum(o.TicksFormat, o.TicksNumDigits, x)
				o.text(dc, o.FontTicks, o.FsizeTicks, "", txt, xi, yf+o.TicksLength, 0.5, 1.0)
			}
		}

		// x-label
		if o.Xlabel != "" {
			xmid := (o.xScr(o.xmin) + o.xScr(o.xmax)) / 2
			ymid := yf + o.TicksLength + tickH
			o.text(dc, o.FontLabel, o.FsizeLabels, "", o.Xlabel, xmid, ymid, 0.5, 1.0)
		}

		// legend @ bottom side
		//
		//        |← LegHlen →|
		//   [gap][   line    |txt][gap][line|txt] ...  ←  yl
		//        ↑                     ↑
		//        x                     x
		//
		if o.LegAtBottom && o.LegendOn {
			for _, curve := range o.curves {

				// icon={line,marker} and label
				if curve.M != "" {
					curve.DrawMarker(dc, xl+hll, yl)
				}
				if curve.Ls != "none" {
					curve.Line(dc, false, true, xl, yl, xl+lll, yl)
				}
				if curve.L != "" {
					xt := xl + lll + o.LegTxtGap
					o.text(dc, o.FontLegend, o.FsizeLegend, "", curve.L, xt, yl, 0.0, 0.3)
				}

				// update column position
				tw := legTxtW
				if o.LegNrow < 2 {
					txtw, _ := dc.MeasureString(curve.L)
					tw = int(txtw)
				}
				xl += lll + o.LegTxtGap + tw + o.LegGap

				// update row position
				if o.LegNrow > 1 {
					if col == ncol-1 {
						col = -1
						xl = o.LegGap
						yl += hei
					}
					col++
				}
			}
		}
	}

	// left ruler
	if LR > 1 {

		// clear background
		o.StyleLR.Rect(dc, false, false, 0, 0, LR, hh+TR)

		// draw ticks and text
		for _, y := range o.yticks {
			yi := o.yScr(y)
			if yi >= y0 && yi <= yf {
				o.StyleFG.Line(dc, false, true, x0-o.TicksLength, yi, x0, yi)
				txt = o.fmtNum(o.TicksFormat, o.TicksNumDigits, y)
				o.text(dc, o.FontTicks, o.FsizeTicks, "", txt, x0-o.TicksLength-1, yi, 1.0, 0.4)
			}
		}

		// y-label
		if o.Ylabel != "" {
			xmid := x0 - o.TicksLength - tickW
			ymid := (o.yScr(o.ymin) + o.yScr(o.ymax)) / 2
			o.text(dc, o.FontLabel, o.FsizeLabels, "", o.Ylabel, xmid, ymid, 1.0, 0.3)
		}
	}

	// top ruler
	if TR > 1 {

		// clear background
		o.StyleTR.Rect(dc, false, false, LR, 0, imax(1, W-LR), TR)

		// draw title
		if o.Title != "" {
			o.text(dc, o.FontTitle, o.FsizeTitle, "", o.Title, W/2, TR/2, 0.5, 0.5)
		}
	}

	// right ruler
	if RR > 1 {

		// clear background
		o.StyleRR.Rect(dc, false, false, xf, y0, imax(1, W-ww-LR), hh)

		// legend @ right side
		//
		//        |← LegHlen →|
		//   [gap][   line    |txt]  ← yl
		//   [gap][   line    |txt]
		//        ↑
		//        xl
		//
		if !o.LegAtBottom && o.LegendOn {
			for _, curve := range o.curves {

				// icon={line,marker} and label
				if curve.M != "" {
					curve.DrawMarker(dc, xl+hll, yl)
				}
				if curve.Ls != "none" {
					curve.Line(dc, false, true, xl, yl, xl+lll, yl)
				}
				if curve.L != "" {
					xt := xl + lll + o.LegTxtGap
					o.text(dc, o.FontLegend, o.FsizeLegend, "", curve.L, xt, yl, 0.0, 0.3)
				}

				// update row position
				yl += hei
			}
		}
	}

	// frame
	if o.DrawBorders {
		dc.SetRGB(0, 0, 0)
		dc.DrawRectangle(float64(x0), float64(y0), float64(ww), float64(hh))
		dc.DrawRectangle(0, 0, float64(W), float64(H))
		dc.Stroke()
	}
}

// auxiliary ////////////////////////////////////////////////////////////////////////////////////

// xScr converts real x-coords to to screen coordinates
func (o *PlotXY) xScr(x float64) int {
	return o.p0 + int(o.sfx*(x-o.xmin))
}

// yScr converts real y-coords to to screen coordinates
func (o *PlotXY) yScr(y float64) int {
	return o.qf - int(o.sfy*(y-o.ymin))
}

// text draws text
func (o *PlotXY) text(dc *Context, f *truetype.Font, size int, clr, txt string, x, y int, ax, ay float64) {
	o.setFont(dc, f, size, clr)
	dc.DrawStringAnchored(txt, float64(x), float64(y), ax, ay)
}

// setFont sets font
func (o *PlotXY) setFont(dc *Context, f *truetype.Font, size int, clr string) {
	if f == nil {
		var err error
		f, err = truetype.Parse(goregular.TTF)
		if err != nil {
			panic(err)
		}
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: float64(size),
	})
	dc.SetFontFace(face)
	dc.SetHexColor("")
}

// fmtNum formats number
func (o *PlotXY) fmtNum(format string, ndigits int, x float64) (l string) {
	val := o.truncate(ndigits, x)
	l = fmt.Sprintf(format, val)
	return
}

// truncate returns a truncated float
func (o *PlotXY) truncate(ndigits int, x float64) (val float64) {
	s := fmt.Sprintf("%."+fmt.Sprintf("%d", ndigits)+"f", x)
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return
}

// measureTxt returns string measures in screen units
func (o *PlotXY) measureTxt(dc *Context, font *truetype.Font, fsz int, txt string) (w, h int) {
	o.setFont(dc, font, fsz, "")
	tw, th := dc.MeasureString(txt)
	return int(tw), int(th)
}

// pretty format ////////////////////////////////////////////////////////////////////////////////

// compute pretty scale numbers
func (o *PlotXY) pretty(Lo, Hi float64, nDiv int) (vals []float64) {

	// constants
	roundingEps := o.cteSqEps
	epsCorrection := 0.0
	shrinkSml := 0.75
	h := 1.5
	h5 := 0.5 + 1.5*h

	// local variables
	minN := int(int(nDiv) / int(3))
	lo := Lo
	hi := Hi
	dx := hi - lo
	cell := 1.0  // cell := "scale" here
	ub := 0.0    // upper bound on cell/unit
	isml := true // is small ?

	// check range
	if !(dx == 0 && hi == 0) { // hi=lo=0

		cell := math.Abs(hi)
		ub := 1.0 + 1.5/(1.0+h5)
		ndiv := 1

		if math.Abs(lo) > math.Abs(hi) {
			cell = math.Abs(lo)
		}
		if h5 >= 1.5*h+0.5 {
			ub = 1.0 + 1.0/(1.0+h)
		}
		if nDiv > 1 {
			ndiv = nDiv
		}
		isml = dx < cell*ub*float64(ndiv)*o.cteEps*3 // added times 3, as several calculations here
	}

	// set cell
	if isml {
		if cell > 10 {
			cell = 9 + cell/10
		}
		cell *= shrinkSml
		if minN > 1 {
			cell /= float64(minN)
		}
	} else {
		cell = dx
		if nDiv > 1 {
			cell /= float64(nDiv)
		}
	}
	if cell < 20*o.cteMin {
		cell = 20 * o.cteMin // very small range.. corrected
	} else if cell*10 > o.cteMax {
		cell = 0.1 * o.cteMax // very large range.. corrected
	}

	// find base and unit
	bas := math.Pow(10.0, math.Floor(math.Log10(cell))) // base <= cell < 10*base
	unit := bas
	ub = 2 * bas
	if ub-cell < h*(cell-unit) {
		unit = ub
		ub = 5 * bas
		if ub-cell < h5*(cell-unit) {
			unit = ub
			ub = 10 * bas
			if ub-cell < h*(cell-unit) {
				unit = ub
			}
		}
	}

	// find number of
	ns := math.Floor(lo/unit + roundingEps)
	nu := math.Ceil(hi/unit - roundingEps)
	if epsCorrection > 0 && (epsCorrection > 1 || !isml) {
		if lo > 0 {
			lo *= (1 - o.cteEps)
		} else {
			lo = -o.cteMin
		}
		if hi > 0 {
			hi *= (1 + o.cteEps)
		} else {
			hi = +o.cteMin
		}
	}
	for ns*unit > lo+roundingEps*unit {
		ns -= 1.0
	}
	for nu*unit < hi-roundingEps*unit {
		nu += 1.0
	}

	// find number of divisions
	ndiv := int(0.5 + nu - ns)
	if ndiv < minN {
		k := minN - ndiv
		if ns >= 0.0 {
			nu += float64(k / 2)
			ns -= float64(k/2 + k%2)
		} else {
			ns -= float64(k / 2)
			nu += float64(k/2 + k%2)
		}
		ndiv = minN
	}
	ndiv++

	// ensure that result covers original range
	if ns*unit < lo {
		lo = ns * unit
	}
	if nu*unit > hi {
		hi = nu * unit
	}

	// fill array
	vals = make([]float64, ndiv)
	vals[0] = lo
	for i := 1; i < ndiv; i++ {
		vals[i] = vals[i-1] + unit
		if math.Abs(vals[i]) < roundingEps {
			vals[i] = 0.0
		}
	}
	return
}
