package gui

import (
	"github.com/anton2920/gofa/gui/color"
)

type Layout struct {
	Foreground     color.Color
	Background     color.Color
	BackgroundDark color.Color
	BackgroundLite color.Color

	WidthForAll  int
	WidthForNext int

	SpacingWidth  int
	SpacingHeight int

	ButtonPaddingWidth  int
	ButtonPaddingHeight int

	SliderTabWidth         int
	SliderTabHeight        int
	SliderSlotHeight       int
	SliderSlotDefaultWidth int /* only used if width can't be computed otherwise. */

	SliderDotSpacing int

	CurrentX, CurrentY int
}

func DefaultLayout() Layout {
	var layout Layout

	layout.Foreground = color.Black
	layout.Background = color.Grey(220)
	layout.BackgroundDark = color.Dark(color.Grey(220))
	layout.BackgroundLite = color.Lite(color.Grey(220))

	layout.SpacingWidth = 20
	layout.SpacingHeight = 4

	layout.WidthForAll = WidthAuto
	layout.WidthForNext = WidthAuto

	layout.ButtonPaddingWidth = 4
	layout.ButtonPaddingHeight = 4

	layout.SliderSlotHeight = 4
	layout.SliderTabWidth = 8
	layout.SliderTabHeight = 16
	layout.SliderSlotDefaultWidth = 100

	layout.SliderDotSpacing = 12

	layout.SetCurrentPosition(0, 0)

	return layout
}

const WidthAuto = 0

func (layout *Layout) Put(width, height *int) Point {
	p := Point{layout.CurrentX, layout.CurrentY}

	if layout.WidthForAll != WidthAuto {
		*width = layout.WidthForAll
	}

	if layout.WidthForNext != WidthAuto {
		*width = layout.WidthForNext
		layout.WidthForNext = WidthAuto
	}

	layout.CurrentY += *height + layout.SpacingHeight

	return p
}

func (layout *Layout) SetCurrentPosition(x, y int) {
	layout.CurrentX = x
	layout.CurrentY = y
}
