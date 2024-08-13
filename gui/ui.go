/* TODO(anton2920):
 * BUGS:
 *	1. If active widget is not called, it's impossible to deactivate it.
 */

package gui

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/gui/color"
	"github.com/anton2920/gofa/gui/fonts"
	"github.com/anton2920/gofa/gui/gr"
	"github.com/anton2920/gofa/trace"
	"github.com/anton2920/gofa/util"
)

type Point struct {
	X, Y int
}

type ID unsafe.Pointer

func ID2(id ID) ID {
	return ID(uintptr(id) + 1)
}

func ID3(id ID) ID {
	return ID(uintptr(id) + 2)
}

func ID4(id ID) ID {
	return ID(uintptr(id) + 3)
}

type WidgetID struct {
	ID     ID
	Index  int
	Parent ID
}

type UI struct {
	Renderer Renderer
	Layout   Layout
	Font     gr.Font

	/* This is set by whichever widget is hot/active; you can watch for it to check. */
	IsHot, IsActive bool

	LeftUp, LeftDown     bool
	MiddleUp, MiddleDown bool
	RightUp, RightDown   bool

	MouseX int
	MouseY int

	Active  WidgetID
	Current WidgetID
	Hot     WidgetID
	HotToBe WidgetID

	dragX, dragY int
}

func NewUI(renderer Renderer) *UI {
	var ui UI

	ui.Renderer = renderer
	ui.Layout = DefaultLayout()
	ui.Font = gr.DecompressFont(fonts.Font21)

	return &ui
}

func (ui *UI) active(id ID) bool {
	if ui.Current.ID != nil {
		id = ui.Current.ID
	}
	return idEquals(id, ui.Current.Index, ui.Current.Parent, ui.Active)
}

func (ui *UI) anyActive() bool {
	return ui.Active.ID != nil
}

func (ui *UI) Begin() {
	ui.Hot = ui.HotToBe
	ui.HotToBe = WidgetID{}

	ui.IsActive = false
	ui.IsHot = false

	ui.Layout.CurrentX = ui.Layout.SpacingWidth
	ui.Layout.CurrentY = ui.Layout.SpacingHeight
}

func (ui *UI) ButtonLogic(id ID, over bool) bool {
	t := trace.Begin("")

	var result bool

	/* NOTE(anton2920): this logic happens correctly for button down then up in one frame, but not up then down. */
	if !ui.anyActive() {
		if over {
			ui.setHot(id)
		}
		if (ui.hot(id)) && (ui.LeftDown) {
			ui.setActive(id)
		}
	}

	if ui.active(id) {
		ui.IsActive = true
		if over {
			ui.setHot(id)
		}
		if ui.LeftUp {
			if ui.hot(id) {
				result = true
			}
			ui.ClearActive()
		}
	}

	if ui.hot(id) {
		ui.IsHot = true
	}

	trace.End(t)
	return result
}

func (ui *UI) ButtonLogicDown(id ID, over bool) bool {
	t := trace.Begin("")

	var result bool

	/* NOTE(anton2920): this logic happens correctly for button down then up in one frame, but not up then down. */
	if !ui.anyActive() {
		if over {
			ui.setHot(id)
		}
		if (ui.hot(id)) && (ui.LeftDown) {
			ui.setActive(id)
			result = true
		}
	}

	if ui.active(id) {
		ui.IsActive = true
		if over {
			ui.setHot(id)
		}
		if ui.LeftUp {
			ui.ClearActive()
		}
	}

	if ui.hot(id) {
		ui.IsHot = true
	}

	trace.End(t)
	return result
}

func (ui *UI) ButtonLogicRect(id ID, x, y, width, height int) bool {
	t := trace.Begin("")

	result := ui.ButtonLogic(id, ui.inRect(x, y, width, height))

	trace.End(t)
	return result
}

func (ui *UI) Button(id ID, label string) bool {
	t := trace.Begin("")

	result := ui.ButtonW(id, label, ui.Font.TextWidth(label)+ui.Layout.ButtonPaddingWidth*2)

	trace.End(t)
	return result
}

func (ui *UI) ButtonToggle(labelUnchecked, labelChecked string, checked *bool) bool {
	t := trace.Begin("")

	widthUnchecked := ui.Font.TextWidth(labelUnchecked)
	widthChecked := ui.Font.TextWidth(labelChecked)

	var label string
	if *checked {
		label = labelChecked
	} else {
		label = labelUnchecked
	}

	result := ui.ButtonW(ID(checked), label, util.Max(widthUnchecked, widthChecked)+ui.Layout.ButtonPaddingWidth*2)
	if result {
		*checked = !*checked
	}

	trace.End(t)
	return result
}

func (ui *UI) ButtonW(id ID, label string, width int) bool {
	t := trace.Begin("")

	height := ui.Font.CharHeight('g') + ui.Layout.ButtonPaddingHeight*2
	p := ui.Layout.Put(&width, &height)

	ui.rectOutlined(p.X, p.Y, width, height, ui.color(id), ui.Layout.BackgroundDark)

	centerX := p.X + width/2 - ui.Font.TextWidth(label)/2
	p.X += ui.Layout.ButtonPaddingWidth
	p.Y += ui.Layout.ButtonPaddingHeight
	ui.Renderer.RenderText(label, ui.Font, centerX, p.Y, ui.Layout.Foreground)

	result := ui.ButtonLogic(id, ui.inRect(p.X, p.Y, width, height))

	trace.End(t)
	return result
}

func (ui *UI) Clear() {
	ui.LeftDown = false
	ui.LeftUp = false

	ui.MiddleDown = false
	ui.MiddleUp = false

	ui.RightDown = false
	ui.RightUp = false
}

func (ui *UI) ClearActive() {
	ui.Active.ID = nil
	ui.Active.Index = 0
	ui.Active.Parent = nil

	/* Mark all UI for this frame as processed. */
	ui.Clear()
}

func (ui *UI) color(id ID) color.Color {
	var clr color.Color
	if ui.hot(id) {
		clr = ui.Layout.BackgroundLite
	} else {
		clr = ui.Layout.Background
	}
	return clr
}

func (ui *UI) DragX(id ID, x *int, width int, y0, y1 int) bool {
	t := trace.Begin("")

	if y1 < y0 {
		y0, y1 = y1, y0
	}

	/* TODO(anton2920): copy below '(*UI) ButtonDownLogic' drag offseting code. */
	ui.ButtonLogic(id, ui.inRect(*x-width/2, util.Min(y0, y1), width, y1-y0))
	if ui.active(id) {
		if ui.MouseX != *x {
			*x = ui.MouseX

			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (ui *UI) DragY(id ID, y *int, height int, x0, x1 int) bool {
	t := trace.Begin("")

	if x1 < x0 {
		x0, x1 = x1, x0
	}

	/* TODO(anton2920): copy below '(*UI) ButtonDownLogic' drag offseting code. */
	ui.ButtonLogic(id, ui.inRect(util.Min(x0, x1), *y-height/2, x1-x0, height))
	if ui.active(id) {
		if ui.MouseY != *y {
			*y = ui.MouseY

			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

/* DragXY is a generic draggable rectangle... If you want its position clamped, do so yourself. */
func (ui *UI) DragXY(id ID, x *int, width int, y *int, height int) bool {
	t := trace.Begin("")

	if ui.ButtonLogicDown(id, ui.inRect(*x-width/2, *y-height/2, width, height)) {
		ui.dragX = *x - ui.MouseX
		ui.dragY = *y - ui.MouseY
	}

	if ui.active(id) {
		if (ui.MouseX+ui.dragX != *x) || (ui.MouseY+ui.dragY != *y) {
			*x = ui.MouseX + ui.dragX
			*y = ui.MouseY + ui.dragY

			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (ui *UI) End() {
	ui.Clear()
}

func (ui *UI) hot(id ID) bool {
	if ui.Current.ID != nil {
		id = ui.Current.ID
	}
	return idEquals(id, ui.Current.Index, ui.Current.Parent, ui.Hot)
}

func (ui *UI) inRect(x, y, width, height int) bool {
	return (ui.MouseX >= x) && (ui.MouseX <= x+width) && (ui.MouseY >= y) && (ui.MouseY <= y+height)
}

func (ui *UI) inRectPlus(x, y, width, height, plus int) bool {
	return (ui.MouseX >= x-plus) && (ui.MouseX <= x+width+plus) && (ui.MouseY >= y-plus) && (ui.MouseY <= y+height+plus)
}

func (ui *UI) MouseMove(x, y int) {
	ui.MouseX = x
	ui.MouseY = y
}

func (ui *UI) MousePress(x, y int, button MouseButton) {
	ui.MouseMove(x, y)

	switch button {
	case Button1:
		ui.LeftDown = true
		ui.LeftUp = false
	case Button2:
		ui.MiddleDown = true
		ui.MiddleUp = false
	case Button3:
		ui.RightDown = true
		ui.RightUp = false
	}
}

func (ui *UI) MouseRelease(x, y int, button MouseButton) {
	ui.MouseMove(x, y)

	switch button {
	case Button1:
		ui.LeftDown = false
		ui.LeftUp = true
	case Button2:
		ui.MiddleDown = false
		ui.MiddleUp = true
	case Button3:
		ui.RightDown = false
		ui.RightUp = true
	}
}

func (ui *UI) rectOutlined(x, y, width, height int, bg, fg color.Color) {
	ui.Renderer.RenderSolidRectWH(x, y, width, height, bg)
	ui.Renderer.RenderRectWH(x, y, width, height, fg)
}

func (ui *UI) setActive(id ID) {
	ui.Active.ID = id
	ui.Active.Index = ui.Current.Index
	ui.Active.Parent = ui.Current.Parent
}

func (ui *UI) setHot(id ID) {
	ui.HotToBe.ID = id
	ui.HotToBe.Index = ui.Current.Index
	ui.HotToBe.Parent = ui.Current.Parent
}

func (ui *UI) Slider(label string, valueMin, valueMax float32, value *float32) bool {
	t := trace.Begin("")

	result := ui.SliderRaw(ID(value), label, valueMin, valueMax, value, false)

	trace.End(t)
	return result
}

func (ui *UI) SliderDisplay(label string, valueMin, valueMax float32, value *float32, display bool) bool {
	t := trace.Begin("")

	if display {
		label = fmt.Sprintf("%s = %g", label, *value)
	}
	result := ui.Slider(label, valueMin, valueMax, value)

	trace.End(t)
	return result
}

func (ui *UI) SliderInt(label string, valueMin, valueMax int, value *int) bool {
	t := trace.Begin("")

	oldValue := *value
	z := float32(*value)

	if ui.SliderRaw(ID(value), label, float32(valueMin), float32(valueMax), &z, true) {
		*value = int(z) + util.Bool2Int((z-float32(int(z))) >= 0.5)

		trace.End(t)
		return oldValue != *value
	}

	trace.End(t)
	return false
}

func (ui *UI) SliderIntDisplay(label string, valueMin, valueMax int, value *int, display bool) bool {
	t := trace.Begin("")

	if display {
		label = fmt.Sprintf("%s = %d", label, *value)
	}
	result := ui.SliderInt(label, valueMin, valueMax, value)

	trace.End(t)
	return result
}

func (ui *UI) SliderRaw(id ID, label string, valueMin, valueMax float32, value *float32, drawDots bool) bool {
	t := trace.Begin("")

	var labelWidth int
	var p Point

	if valueMax < valueMin {
		valueMin, valueMax = valueMax, valueMin
	}

	sliderWidth := ui.Layout.WidthForAll
	sliderHeight := ui.Layout.SliderSlotHeight

	tabWidth := ui.Layout.SliderTabWidth
	tabHeight := ui.Layout.SliderTabHeight

	const labelWidthAdjustment = 4
	if label != "" {
		labelHeight := ui.Font.CharHeight('g') - ui.Layout.SpacingHeight + 1
		labelWidth = ui.Font.TextWidth(label) + labelWidthAdjustment
		p = ui.Layout.Put(&labelWidth, &labelHeight)
		ui.Renderer.RenderText(label, ui.Font, p.X, p.Y, ui.Layout.Foreground)
	}

	if sliderWidth == WidthAuto {
		if label != "" {
			sliderWidth = labelWidth + labelWidthAdjustment
		} else {
			sliderWidth = ui.Layout.SliderSlotDefaultWidth
		}
		if sliderWidth < 50 {
			sliderWidth = 50
		}
	}

	p = ui.Layout.Put(&sliderWidth, &tabHeight)

	/* Compute location of left edge of tab. */
	pos := int(float32(p.X) + (*value-valueMin)/(valueMax-valueMin)*float32(sliderWidth) - float32(tabWidth)/2)

	ui.rectOutlined(p.X, p.Y+(tabHeight-sliderHeight)*3/4, sliderWidth, sliderHeight, ui.color(id), ui.Layout.BackgroundDark)

	if drawDots {
		n := int(valueMax - valueMin + 1)
		if sliderWidth >= ui.Layout.SliderDotSpacing*n {
			for i := 0; i < n; i++ {
				pos := p.X + i*int(float32(sliderWidth)/(valueMax-valueMin))
				ui.Renderer.RenderPoint(pos, p.Y+(tabHeight-sliderHeight)/4, 1, ui.Layout.BackgroundDark)
			}
		}
	}

	ui.rectOutlined(pos, p.Y, tabWidth, tabHeight, ui.Layout.BackgroundLite, ui.Layout.BackgroundDark)

	ui.ButtonLogic(id, ui.inRectPlus(p.X, p.Y+(tabHeight-sliderHeight)/2, sliderWidth, sliderHeight+(tabHeight-sliderHeight)/4, 2) || ui.inRectPlus(pos, p.Y, tabWidth, tabHeight, 1))

	if ui.active(id) {
		oldValue := *value
		z := float32(ui.MouseX-p.X)*(valueMax-valueMin)/float32(sliderWidth) + valueMin
		if z < valueMin {
			z = valueMin
		} else if z > valueMax {
			z = valueMax
		}
		*value = z

		ui.setHot(id) /* sliders are always hot while active. */

		trace.End(t)
		return *value != oldValue
	}

	trace.End(t)
	return false
}

func idEquals(id ID, index int, parent ID, dest WidgetID) bool {
	return (id == dest.ID) && (index == dest.Index) && (parent == dest.Parent)
}
