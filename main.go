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
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg: 
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
				s += "> "
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
	fileReader, err := os.Open(csvPath);
	if err != nil {
		log.Print("Occured during opening CSV: ", err.Error())
		return
	}

	 csvReader := csv.NewReader(fileReader)
	 csvData, err := csvReader.ReadAll()
	 if err != nil {
		 log.Print("Could not read CSV Data: ", err.Error())
		 return
	 }

	model := Model {
		cursorPos: CursorCoordination {
			row: 0,
			col: 0,
		},
		data: csvData,
	}
	
	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		log.Print("Running TUI: ", err.Error())
		os.Exit(1)
	}
}
