package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom. The objects are always displayed
// at their vertical MinSize. Use a different layout if the objects are intended
// to be larger than their vertical MinSize.
func NewVBoxLayout() fyne.Layout {
	return vBoxLayout{
		paddingFunc: theme.Padding,
	}
}

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right. The objects are always displayed
// at their horizontal MinSize. Use a different layout if the objects are intended
// to be larger than their horizontal MinSize.
func NewHBoxLayout() fyne.Layout {
	return hBoxLayout{
		paddingFunc: theme.Padding,
	}
}

// NewCustomPaddedHBoxLayout returns a layout similar to HBoxLayout that uses a custom
// amount of padding in between objects instead of the theme.Padding value.
//
// Since: 2.5
func NewCustomPaddedHBoxLayout(padding float32) fyne.Layout {
	return hBoxLayout{
		paddingFunc: func() float32 { return padding },
	}
}

// NewCustomPaddedVBoxLayout returns a layout similar to VBoxLayout that uses a custom
// amount of padding in between objects instead of the theme.Padding value.
//
// Since: 2.5
func NewCustomPaddedVBoxLayout(padding float32) fyne.Layout {
	return vBoxLayout{
		paddingFunc: func() float32 { return padding },
	}
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*vBoxLayout)(nil)

type vBoxLayout struct {
	paddingFunc func() float32
}

// Layout is called to pack all child objects into a specified size.
// This will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (v vBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	visibleObjects := 0
	// Size taken up by visible objects
	totalMinHights := float32(0)
	growFactors := float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isVerticalSpacer(child) {
			spacers++
			continue
		}

		visibleObjects++

		minH := child.MinSize().Height
		totalMinHights += minH
		if growable, ok := child.(fyne.Growable); ok {
			if minH < 1.0 {
				minH = 1.0
			}
			growFactors += growable.GrowFactor() * minH
		} else {
			child.Resize(fyne.NewSize(size.Width, minH))
		}
	}

	padding := v.paddingFunc()

	// Amount of space not taken up by visible objects and inter-object padding
	extra := size.Height - totalMinHights - (padding * float32(visibleObjects-1))
	hPerFactor := float32(0)
	if growFactors > 0.0 {
		hPerFactor = extra / growFactors
	}

	// Spacers split extra space equally
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x, y := float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isVerticalSpacer(child) {
			y += spacerSize
			continue
		}
		child.Move(fyne.NewPos(x, y))

		height := child.MinSize().Height
		if growable, ok := child.(fyne.Growable); ok {
			height += hPerFactor * growable.GrowFactor() * height
			child.Resize(fyne.NewSize(size.Width, height))
		}
		y += padding + height
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of all children combined with padding between each.
func (v vBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	padding := v.paddingFunc()
	for _, child := range objects {
		if !child.Visible() || isVerticalSpacer(child) {
			continue
		}

		childMin := child.MinSize()
		minSize.Width = fyne.Max(childMin.Width, minSize.Width)
		minSize.Height += childMin.Height
		if addPadding {
			minSize.Height += padding
		}
		addPadding = true
	}
	return minSize
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*hBoxLayout)(nil)

type hBoxLayout struct {
	paddingFunc func() float32
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g hBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	visibleObjects := 0
	// Size taken up by visible objects
	totalMinWidths := float32(0)
	growFactors := float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isHorizontalSpacer(child) {
			spacers++
			continue
		}

		visibleObjects++

		minW := child.MinSize().Width
		totalMinWidths += minW
		if growable, ok := child.(fyne.Growable); ok {
			if minW < 1.0 {
				minW = 1.0
			}
			growFactors += growable.GrowFactor() * minW
		} else {
			child.Resize(fyne.NewSize(minW, size.Height))
		}
	}

	padding := g.paddingFunc()

	// Amount of space not taken up by visible objects and inter-object padding
	extra := size.Width - totalMinWidths - (padding * float32(visibleObjects-1))
	wPerFactor := float32(0.0)
	if growFactors > 0.0 {
		wPerFactor = extra / growFactors
	}

	// Spacers split extra space equally
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x, y := float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isHorizontalSpacer(child) {
			x += spacerSize
			continue
		}
		child.Move(fyne.NewPos(x, y))

		width := child.MinSize().Width
		if growable, ok := child.(fyne.Growable); ok {
			width += wPerFactor * growable.GrowFactor() * width
			child.Resize(fyne.NewSize(width, size.Height))
		}
		x += padding + width
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of all children combined with padding between each.
func (g hBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	padding := g.paddingFunc()
	for _, child := range objects {
		if !child.Visible() || isHorizontalSpacer(child) {
			continue
		}

		childMin := child.MinSize()
		minSize.Height = fyne.Max(childMin.Height, minSize.Height)
		minSize.Width += childMin.Width
		if addPadding {
			minSize.Width += padding
		}
		addPadding = true
	}
	return minSize
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandVertical()
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandHorizontal()
}
