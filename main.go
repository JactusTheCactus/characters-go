package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"sort"
	"strconv"
	"strings"
)

func cap(s string) string {
	return cases.Title(language.English, cases.NoLower).String(s)
}
func uni(u string) string {
	var err error
	u, err = strconv.Unquote(`"` + u + `"`)
	if err != nil {
		fmt.Println(err)
	}
	return u
}

var u = map[string]string{
	"dot":    uni(`\u00B7`),
	"male":   uni(`\u2642`),
	"female": uni(`\u2640`),
	"neuter": uni(`\u26A5`),
	"acute":  uni(`\u0301`),
	"grave":  uni(`\u0300`),
	"omega":  uni(`\u03C9`),
	"schwa":  uni(`\u0259`),
}

type Char struct {
	name    []string
	pron    []string
	species []string
	sex     string
	extra   []string
}

func (c Char) CharGet(key string) any {
	switch key {
	case "name":
		return c.name
	case "pron":
		return c.pron
	case "species":
		return c.species
	case "sex":
		return c.sex
	case "extra":
		return c.extra
	default:
		return nil
	}
}
func (c Char) exists() bool {
	e := false
	for _, i := range []string{
		"name",
		"pron",
		"species",
		"sex",
		"extra",
	} {
		stat := c.CharGet(i)
		switch s := stat.(type) {
		case []string:
			if len(s) > 0 {
				e = true
			}
		case string:
			if len(s) > 0 {
				e = true
			}
		}
	}
	return e
}
func (c Char) getSex() string {
	s, ok := u[strings.ToLower(c.sex)]
	if !ok {
		s = u["neuter"]
	}
	return s
}

type model struct {
	chosen   bool
	selected int
	cursor   int
	char     Char
}

var chars = []Char{
	{
		[]string{
			"Hound",
		},
		[]string{
			"haund",
		},
		[]string{
			"Changeling",
		},
		"Female",
		[]string{
			"Shapeshifts into a large, black Wolf",
		},
	},
	{
		[]string{
			"Morrigan",
		},
		[]string{
			fmt.Sprint("M", u["omega"], u["acute"], "r", u["schwa"], "gy", u["grave"], "n"),
		},
		[]string{
			"Reaper",
		},
		"Female",
		[]string{
			"Wields a scythe",
			"Killing Touch",
		},
	},
}

func initialModel() model {
	return model{
		chosen: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q",
			"esc":
			return m, tea.Quit
		case "up",
			"w":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(chars) - 1
			}
		case "down",
			"s":
			if m.cursor < len(chars)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter",
			" ",
			"left",
			"a",
			"right",
			"d":
			if m.selected == m.cursor && m.char.exists() {
				m.chosen = !m.chosen
				m.char = Char{}
			} else {
				m.chosen = true
				m.selected = m.cursor
				m.char = chars[m.selected]
			}
		}
	}
	title := ""
	switch m.chosen {
	case true:
		title = fmt.Sprintf(
			"Characters | %s%s",
			m.char.getSex(),
			strings.Join(m.char.name, " "),
		)
	case false:
		title = "Characters"
	}
	return m, tea.Batch(
		tea.SetWindowTitle(title),
	)
}
func log(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprintf("\t%s\n", format),
		a...,
	)
}
func (c Char) getStat(s string) string {
	stat := c.CharGet(s)
	o := log("%s:", cap(s))
	switch s {
	case "species",
		"extra":
		sort.Strings(stat.([]string))
	}
	switch st := stat.(type) {
	case string:
		o += log("\t- %s", s)
	case []string:
		for _, s := range st {
			o += log("\t- %s", s)
		}
	}
	return o
}
func getChar(c Char) string {
	var o string
	sex := c.getSex()
	o += log("%s%s", sex, strings.Join(c.name, " "))
	o += log("\t<%s>", strings.Join(c.pron, u["dot"]))
	o += c.getStat("species")
	o += c.getStat("extra")
	return o
}
func (m model) View() string {
	s := "Pick a character:\n"
	for i, char := range chars {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, strings.Join(char.name, " "))
	}
	s += "\n"
	if m.chosen {
		s += getChar(m.char)
	}
	return s
}
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Print("Uh-Oh:", err)
		os.Exit(1)
	}
}
