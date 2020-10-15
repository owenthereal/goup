// Package MultipleChoice is intended for use in CLI applications
package MultipleChoice

import (
	"fmt"
	"strings"

	"github.com/gosuri/uilive"
	"github.com/pkg/term"
)

// Selection allows the program to ask a question
// and give the user several options to chose from,
// then returning the selected option.
func Selection(question string, options []string) string {
	return MultipleChoice(question, options, false)[0]
}

func MultiSelection(question string, options []string) []string {
	return MultipleChoice(question, options, true)
}

const multiSelectInstruction = "Use space to toggle selection"

type option struct {
	text     string `json:"value"`
	selected bool   `json:"selected"` // For checkboxes.
}

type multipleChoice struct {
	writer      *uilive.Writer
	question    string
	position    int
	multiSelect bool
	options     []*option
}

func MultipleChoice(question string, options []string, multiSelect bool) []string {
	m := &multipleChoice{
		writer:      uilive.New(),
		question:    question,
		options:     make([]*option, 0),
		position:    0,
		multiSelect: multiSelect,
	}
	for i := 0; i < len(options); i++ {
		m.options = append(m.options, &option{text: options[i], selected: false})
	}
	return m.start()
}

func (m *multipleChoice) toggle() {
	m.options[m.position].selected = !m.options[m.position].selected
	m.draw()
}

func (m *multipleChoice) down() {
	m.position = (m.position + 1) % len(m.options)
	m.draw()
}

func (m *multipleChoice) up() {
	m.position = (m.position + len(m.options) - 1) % len(m.options)
	m.draw()
}

func (m *multipleChoice) draw() {
	m.writer.Start()
	output := []string{}
	output = append(output, fmt.Sprintf("%s", strings.TrimSpace(m.question)))
	for i := 0; i < len(m.options); i++ {
		if m.multiSelect {
			prefix := " "
			if m.position == i {
				prefix = ">"
			}
			selected := "-"
			if m.options[i].selected {
				selected = "√"
			}
			output = append(output, fmt.Sprintf("%s %s %s", prefix, selected, strings.TrimSpace(m.options[i].text)))
		} else {
			prefix := "-"
			if m.position == i {
				prefix = ">"
			}
			output = append(output, fmt.Sprintf("%s %s", prefix, strings.TrimSpace(m.options[i].text)))
		}
	}
	if m.multiSelect {
		output = append(output, fmt.Sprintf("%s", multiSelectInstruction))
	}
	fmt.Fprint(m.writer, strings.Join(output, "\n")+"\n")
	m.writer.Stop()
}

func (m *multipleChoice) start() []string {
	// start listening for updates and render
	m.draw()

	for {

		ascii, keyCode, _ := getChar()

		if ascii == 13 {
			selected := []string{}
			if m.multiSelect {
				for i := 0; i < len(m.options); i++ {
					if m.options[i].selected {
						selected = append(selected, m.options[i].text)
					}
				}
			} else {
				selected = append(selected, m.options[m.position].text)
			}
			return selected
		} else if ascii == 32 {
			if m.multiSelect {
				m.toggle()
			}
		}

		if keyCode == 38 {
			m.up()
		} else if keyCode == 40 {
			m.down()
		}
	}
	// return ""
}

// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// up
			keyCode = 38
		} else if bytes[2] == 66 {
			// down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}
