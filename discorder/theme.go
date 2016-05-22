package discorder

import (
	"encoding/json"
	"github.com/jonas747/discorder/ui"
	"github.com/nsf/termbox-go"
	"log"
)

type Theme struct {
	Name      string             `json:"name"`
	Author    string             `json:"author"`
	Comment   string             `json:"comment"`
	ColorMode termbox.OutputMode `json:"color_mode"` // see termbox.OutputMode for info

	Theme map[string]ThemeAttribPair `json:"theme"`
}

func (t *Theme) GetAttribute(key string, fg bool) (attrib termbox.Attribute, ok bool) {
	var pair ThemeAttribPair
	pair, ok = t.Theme[key]
	if !ok {
		return
	}

	if fg {
		attrib = pair.FG.Attribute()
	} else {
		attrib = pair.BG.Attribute()
	}
	return
}

func (app *App) GetThemeAttribPair(key string) ThemeAttribPair {
	if app.userTheme != nil {
		pair, ok := app.userTheme.Theme[key]
		if ok {
			return pair
		}
	}

	pair, _ := app.defaultTheme.Theme[key]
	return pair
}

func (app *App) GetThemeAttribute(key string, fg bool) termbox.Attribute {
	if app.userTheme != nil {
		userAttrib, ok := app.userTheme.GetAttribute(key, fg)
		if ok {
			return userAttrib
		}
	}
	defaultAttrib, _ := app.defaultTheme.GetAttribute(key, fg)
	return defaultAttrib
}

func (app *App) ApplyThemeToList(list *ui.ListWindow) {
	list.NormalFG = app.GetThemeAttribute("element_normal", true)
	list.NormalBG = app.GetThemeAttribute("element_normal", false)
	list.MarkedFG = app.GetThemeAttribute("element_marked", true)
	list.MarkedBG = app.GetThemeAttribute("element_marked", false)
	list.SelectedFG = app.GetThemeAttribute("element_selected", true)
	list.SelectedBG = app.GetThemeAttribute("element_selected", false)
	list.MarkedSelectedFG = app.GetThemeAttribute("element_selected_marked", true)
	list.MarkedSelectedBG = app.GetThemeAttribute("element_selected_marked", false)

	app.ApplyThemeToWindow(list.Window)
}

func (app *App) ApplyThemeToWindow(window *ui.Window) {
	window.BorderFG = app.GetThemeAttribute("window_border", true)
	window.BorderBG = app.GetThemeAttribute("window_border", false)
	window.FillBG = app.GetThemeAttribute("window_fill", false)
}

func (app *App) ApplyThemeToText(text *ui.Text, key string) {
	pair := app.GetThemeAttribPair(key)
	text.BG = pair.BG.Attribute()
	text.FG = pair.FG.Attribute()
}

func (t *Theme) Read() ([]byte, error) {
	out, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ThemeAttribute struct {
	Color     Color `json:"color"`
	Bold      bool  `json:"bold"`
	Underline bool  `json:"underline"`
	Reverse   bool  `json:"reverse"`
}

func (t *ThemeAttribute) Attribute() termbox.Attribute {
	attr := termbox.Attribute(t.Color)
	if t.Bold {
		attr |= termbox.AttrBold
	}
	if t.Underline {
		attr |= termbox.AttrUnderline
	}
	if t.Reverse {
		attr |= termbox.AttrReverse
	}
	return attr
}

type ThemeAttribPair struct {
	FG ThemeAttribute `json:"fg"`
	BG ThemeAttribute `json:"bg"`
}

func (t ThemeAttribPair) AttribPair() ui.AttribPair {
	return ui.AttribPair{
		FG: t.FG.Attribute(),
		BG: t.BG.Attribute(),
	}
}

type Color uint8

func (c *Color) UnmarshallJSON(data []byte) error {
	log.Println("Unmarshalling color")
	var raw interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	switch t := raw.(type) {
	case string:
		col, ok := Colors[t]
		if !ok {
			log.Println("Color not found", t)
		}
		*c = Color(col)
	case float64:
		intColor := uint8(t)
		*c = Color(intColor)
	}
	return nil
}

var Colors = map[string]uint8{
	"default": 0,
	"black":   1,
	"red":     2,
	"green":   3,
	"yellow":  4,
	"blue":    5,
	"magenta": 6,
	"cyan":    7,
	"white":   8,
}
