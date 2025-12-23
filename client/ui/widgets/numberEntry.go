package widgets

import (
	"strconv"

	"fyne.io/fyne/v2/widget"
)

type NumberEntry struct {
	widget.Entry
	Value    int
	MinValue int
	MaxValue int
}

func NewNumberEntry() *NumberEntry {
	ne := &NumberEntry{}
	ne.ExtendBaseWidget(ne)
	ne.OnChanged = func(s string) {
		// Empty input -> set to MinValue
		if s == "" {
			ne.Value = ne.MinValue
			d := strconv.Itoa(ne.Value)
			if s != d {
				ne.SetText(d)
			}
			return
		}

		if _, err := strconv.Atoi(s); err != nil {
			if len(s) > 0 {
				ne.SetText(s[:len(s)-1])
			}
			return
		}

		// Parse and clamp
		ne.Value, _ = strconv.Atoi(s)
		if ne.Value < ne.MinValue {
			ne.Value = ne.MinValue
		}
		if ne.MaxValue != 0 && ne.Value > ne.MaxValue {
			ne.Value = ne.MaxValue
		}

		// If we adjusted the value, update the displayed text
		d := strconv.Itoa(ne.Value)
		if s != d {
			ne.SetText(d)
		}
	}
	ne.Text = strconv.Itoa(ne.Value)
	return ne
}

func (ne *NumberEntry) ChangeValue(newValue int) {
	ne.Value = newValue
	if ne.Value < ne.MinValue {
		ne.Value = ne.MinValue
	}
	if ne.MaxValue != 0 && ne.Value > ne.MaxValue {
		ne.Value = ne.MaxValue
	}
	ne.SetText(strconv.Itoa(ne.Value))
}
