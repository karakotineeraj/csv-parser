package main

import (
	"encoding/csv"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorCoordination struct {
	row int;
	col int;
}

type Model struct {
	cursorPos CursorCoordination;
	data [][]string;
	editMode bool;
	csvWriter *csv.Writer;
	fileWriter *os.File;
	editCursor int;
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editMode {
			switch msg.Type {
			case tea.KeyEnter: 
				m.editMode = false
			case tea.KeyLeft:
				if m.editCursor > 0 {
					m.editCursor--
				}
			case tea.KeyRight:
				if m.editCursor < len(m.data[m.cursorPos.row][m.cursorPos.col]) {
					m.editCursor++
				}
			case tea.KeyBackspace:
				if m.editCursor > 0 {
					m.data[m.cursorPos.row][m.cursorPos.col] = m.data[m.cursorPos.row][m.cursorPos.col][:m.editCursor-1] + m.data[m.cursorPos.row][m.cursorPos.col][m.editCursor:]
					m.editCursor--
				}
			case tea.KeyRunes:
				m.data[m.cursorPos.row][m.cursorPos.col] += msg.String()
				m.editCursor++
			}

			return m, nil
		}

		switch msg.String() {
		case "j":
			if m.cursorPos.row < len(m.data) - 1 {
				m.cursorPos.row++
			}
		case "k":
			if m.cursorPos.row > 0 {
				m.cursorPos.row--
			}
			case "h": 
			if m.cursorPos.col > 0 {
				m.cursorPos.col--
			}
		case "l":
			if m.cursorPos.col < len(m.data[m.cursorPos.row]) - 1 {
				m.cursorPos.col++
			}
		case "w":
			m.fileWriter.Truncate(0);
			m.fileWriter.Seek(0, 0);
			if err := m.csvWriter.WriteAll(m.data); err != nil {
				return m, nil
			}
		case "enter":
			m.editMode = true
			m.editCursor = len(m.data[m.cursorPos.row][m.cursorPos.col])
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := ""

	for i, rowVal := range m.data {
		for j, val := range rowVal {
			if i == m.cursorPos.row && j == m.cursorPos.col {
				if m.editMode {
					s += ">" + val[:m.editCursor] + "|" + val[m.editCursor:] + "\t"
					continue;
				} else {
					s += ">"
				}
			}

			s += val + "\t"
		}

		s += "\n"
	}

	return s
}

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	csvPath := os.Args[1]
	fileReader, err := os.OpenFile(csvPath, os.O_RDWR, 0646);
	if err != nil {
		log.Print("Occured during opening CSV: ", err.Error())
		return
	}

	csvReader := csv.NewReader(fileReader)
	csvReader.TrimLeadingSpace = true
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Print("Could not read CSV Data: ", err.Error())
		return
	}

	csvWriter := csv.NewWriter(fileReader)

	model := Model {
		cursorPos: CursorCoordination {
			row: 0,
			col: 0,
		},
		data: csvData,
		editMode: false,
		csvWriter: csvWriter,
		fileWriter: fileReader,
	}

	p := tea.NewProgram(model, tea.WithInput(os.Stdin))
	if _, err = p.Run(); err != nil {
		log.Print("Running TUI: ", err.Error())
		os.Exit(1)
	}
}
