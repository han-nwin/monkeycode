package tui

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
    "github.com/charmbracelet/lipgloss"
)

var (
    CorrectRune int
)

// Styling variables
var (
	// Title style 
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F5C2E7")). // Pinkish-purple (Rosewater)
		Background(lipgloss.Color("#1E1E2E")). // Dark background (Base)
		Padding(1, 1)

	// Option style 
	optionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8AADF4")). // Soft blue (Blue)
		PaddingLeft(2)

	// Stats style 
	statsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6E3A1")). // Lime green (Green)
		Bold(true).
		PaddingLeft(1)

	// Instruction style 
	instructionStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#9399B2")). // Muted gray (Subtext1)
		PaddingLeft(1)

	// Correct character style 
	correctCharStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C9CBFF")). // Soft white (Text)
		Bold(true)

	// Incorrect character style 
	incorrectCharStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F38BA8")). // Red (Red)
		Bold(true)

	// Untyped character style 
	untypedCharStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#585B70")) // Dimmed gray (Surface1)

	// Cursor style (lavender)
	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CBA6F7")) // (Lavender)

	// WPM style
	wpmStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6E3A1")). // Lime green (Green)
		Bold(true)
)


func renderBackground(width, height int, color string) string {
	backgroundStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(color)).
		Width(width).
		Height(height)

	// Create a box with spaces filling the screen
	return backgroundStyle.Render(strings.Repeat(" ", width*height))
}

// Embed all files in the languages directory
//go:embed languages/*
var embeddedFiles embed.FS

//Cursor tick speed
type cursorTick struct{}

func cursorCmd(num time.Duration) tea.Cmd {
    return tea.Tick(time.Millisecond*num, func (t time.Time) tea.Msg {
        return cursorTick{}
    })
}

//Custom enum-like State struck
type State int

const (
	PreProgram State = iota
	TypingTUI
	Results
)
//TUI MODEL
type Model struct {
    State   State 
    Languages []string
    SelectedLang string
    PromptText string
    UserText *strings.Builder
    CursorVisible bool
    CursorPosition int
    Width int
    Height int
    StartTime time.Time
    FinalWPM int
    Accuracy float64
}

//Init model
func (m Model) Init() tea.Cmd {
    return cursorCmd(450)
}

//Update logic
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Always quit on Ctrl+C and Ctrl+Z
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if keyMsg.Type == tea.KeyCtrlC || keyMsg.Type == tea.KeyCtrlZ {
            return m, tea.Quit
        }
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:
		//======Handle PreProgram state
		if m.State == PreProgram {
			switch msg.String() {
			case "1":
				m.SelectedLang = "go"
			case "2":
				m.SelectedLang = "python"
			case "3":
				m.SelectedLang = "c++"
            case "q":
                return m, tea.Quit
			default:
				return m, nil
			}

			// Transition to TypingTUI state
			m.State = TypingTUI
			m.PromptText = LoadPrompt(m.SelectedLang)
			m.StartTime = time.Now()
			m.UserText = &strings.Builder{}
			return m, nil
		}

        //=======Handle TypingTUI State
        if m.State == TypingTUI {
            // Check if typing is complete
            if m.UserText.Len() == len(m.PromptText) {
                m.State = Results
                m.FinalWPM = CalculateWPM(m.UserText.String(), m.StartTime)
                return m, nil
            }

            //if user wanna restart
            if msg.Type == tea.KeyCtrlR {
                m.State = PreProgram
                return m, nil
            }

            //On every keystroke, add one character to userText
            if msg.Type == tea.KeyEnter { //type Enter
                m.UserText.WriteRune('\n')
                CorrectRune++

            } else if msg.Type == tea.KeyTab { //type Tab
                m.UserText.WriteString("    ")
                CorrectRune += 4;

            } else if msg.Type == tea.KeyBackspace { //type Backspace
                currentText := m.UserText.String()
                currentRunes := []rune(currentText)
                if len(currentRunes) > 0 {
                    lastTypedIndex := len(currentRunes) - 1

                    // Check if the last rune was correct and decrement `correctRune` if necessary
                    if lastTypedIndex < len(m.PromptText) && currentRunes[lastTypedIndex] == rune(m.PromptText[lastTypedIndex]) {
                        CorrectRune-- // Decrement if the last typed character was correct
                    }

                    // Remove the last rune from the user text
                    currentRunes = currentRunes[:lastTypedIndex]
                    m.UserText.Reset()
                    m.UserText.WriteString(string(currentRunes))
                }
            } else if len(msg.Runes) > 0 { //other keys
                currentIndex := m.UserText.Len()
                if currentIndex < len(m.PromptText) && rune(m.PromptText[currentIndex]) == msg.Runes[0]{
                    CorrectRune++ // Increment correctRune here
                }
                m.UserText.WriteRune(msg.Runes[0])

            }
            return m, nil
        }
        //===Handle Results state
        if m.State == Results {
            if len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
                return m, tea.Quit
            } else {
                return m, nil
            }
        }
    case cursorTick:
        //Toggle the cursor visible
        m.CursorVisible = !m.CursorVisible
        return m, cursorCmd(450)

    case tea.WindowSizeMsg:
        m.Width = msg.Width
        m.Height = msg.Height
    }

    return m, nil
}

//View logic
func (m Model) View() string {
	switch m.State {
	case PreProgram:
        title := titleStyle.Render("Select a Language")
		options := lipgloss.JoinVertical(lipgloss.Left,
			optionStyle.Render("\n[1] Go"),
			optionStyle.Render("[2] Python"),
			optionStyle.Render("[3] C++"),
		)
        instructions := lipgloss.JoinVertical(lipgloss.Left,
            instructionStyle.Render("\nUse number keys to pick a language."),
            instructionStyle.Render("q: exit"),
        )

		return lipgloss.JoinVertical(lipgloss.Left, title, options, instructions)

    case TypingTUI:
        userInput := m.UserText.String()
        prompt := m.PromptText

        var renderedText strings.Builder
        // Inject cursor at the next untyped position
        cursorPosition := len(userInput)

        // Update final WPM only if typing is not complete
        m.FinalWPM = CalculateWPM(userInput, m.StartTime)
        //Calculate Accuracy
        m.Accuracy = CalculateAccuracy(m.UserText.String(), CorrectRune)

        for i, r := range prompt {
            if i == cursorPosition && m.CursorVisible {
                renderedText.WriteString(cursorStyle.Render("|")) // Add the styled cursor
            }
            if i < len(userInput) {
                if rune(userInput[i]) == r {
                    // Correct character
                    renderedText.WriteString(correctCharStyle.Render(string(r)))
                } else {
                    // Incorrect character
                    renderedText.WriteString(incorrectCharStyle.Render(string(userInput[i])))
                }
            } else {
                // Untyped character
                renderedText.WriteString(untypedCharStyle.Render(string(r)))
            }
        }
        renderedText.WriteString("\n\n")
        // Add WPM to the output
        renderedText.WriteString(wpmStyle.Render(fmt.Sprintf("WPM: %d", m.FinalWPM)))

        //renderedText.WriteString(fmt.Sprintf("\033[90m%s\033[0m","\n* Press any key to exit after finishing"))
        renderedText.WriteString(instructionStyle.Render("\n* Ctrl-R: restart"))

        return wordwrap.String(renderedText.String() , m.Width)

    case Results:
        // Render the styled elements using global styles
        title := titleStyle.Render("Typing Complete!")
        stats := lipgloss.JoinVertical(
            lipgloss.Left,
            statsStyle.Render(fmt.Sprintf("WPM: %d", m.FinalWPM)),
            statsStyle.Render(fmt.Sprintf("Accuracy: %.2f%%", m.Accuracy)),
        )
        instructions := lipgloss.JoinVertical(
            lipgloss.Left,
            instructionStyle.Render("Press q to exit."),
            instructionStyle.Render("Press r to restart."),
        )

        // Combine all the parts into a single view
        return lipgloss.JoinVertical(
            lipgloss.Center,
            title,
            stats,
            instructions,
        )
    }
    return ""
}

//function to count wpm
func CalculateWPM(userText string, startTime time.Time) int {
    elapsed := time.Since(startTime).Minutes() //time elapsed
    if elapsed == 0 {
        return 0
    }

    words := strings.Fields(userText) //split by white space
    length := len(words)

    return int(float64(length) / elapsed)
}

//function to calculate accuracy
func CalculateAccuracy(userText string, correctRune int) float64 {
    totalRunes := []rune(userText)

    return min(100.00,(float64(correctRune) / float64(len(totalRunes)))*100)
}

//Function to load prompt in main package
func LoadPrompt(language string) string {

    lang := language
    path := fmt.Sprintf("languages/%s", lang)

    // Read embedded directory
    entries, err := embeddedFiles.ReadDir(path)
    if err != nil {
        fmt.Printf("Cannot access 'languages/%s' directory: %v\n", lang, err)
        os.Exit(-1)
    }

    // Randomly select a file from the embedded directory
    randIdx := rand.Intn(len(entries))
    selectedFile := entries[randIdx]

    // Read the content of the randomly selected file
    filePath := fmt.Sprintf("%s/%s", path, selectedFile.Name())
    data, err := embeddedFiles.ReadFile(filePath)
    if err != nil {
        fmt.Printf("Error loading content file: %v\n", err)
        os.Exit(-1)
    }

    prompt := strings.TrimSpace(string(data))
    promptReplaceTab := strings.ReplaceAll(prompt, "\t", "    ")//replace tab into 4 spaces

    return promptReplaceTab
}

func NewProgram()  *tea.Program {
    m := Model{
        PromptText:     "",
        UserText:       &strings.Builder{},
        CursorVisible:  true,
        CursorPosition: 0,
        Width:          80,
        Height:         20,
        FinalWPM:       0,
    }

    program := tea.NewProgram(m, tea.WithAltScreen())
    return program
}
