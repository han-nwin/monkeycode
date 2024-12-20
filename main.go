package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	//"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

var (
    language *string
    list_languages *bool
    version *bool

    correctRune int
)

//Cursor tick speed
type cursorTick struct{}

func cursorCmd(num time.Duration) tea.Cmd {
    return tea.Tick(time.Millisecond*num, func (t time.Time) tea.Msg {
        return cursorTick{}
    })
}


//TUI MODEL
type model struct {
    promptText string
    userText *strings.Builder
    cursorVisible bool
    cursorPosition int
    width int
    height int
    startTime time.Time
    isComplete bool
    finalWPM int
}

//Init model
func (m model) Init() tea.Cmd {
    return cursorCmd(450)
}

//Update logic
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Always handle Ctrl+C and Ctrl+Z
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if keyMsg.Type == tea.KeyCtrlC || keyMsg.Type == tea.KeyCtrlZ {
            return m, tea.Quit
        }
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Check if typing is complete
        if m.userText.Len() == len(m.promptText) {
            m.isComplete = true
            m.finalWPM = calculateWPM(m.userText.String(), m.startTime)
            return m, tea.Quit
        }

        //On every keystroke, add one character to userText
        if msg.Type == tea.KeyEnter { //type Enter
            m.userText.WriteRune('\n')
            correctRune++

        } else if msg.Type == tea.KeyTab { //type Tab
            m.userText.WriteString("    ")
            correctRune += 4;

        } else if msg.Type == tea.KeyBackspace { //type Backspace
            currentText := m.userText.String()
            currentRunes := []rune(currentText)
            if len(currentRunes) > 0 {
                lastTypedIndex := len(currentRunes) - 1

                // Check if the last rune was correct and decrement `correctRune` if necessary
                if lastTypedIndex < len(m.promptText) && currentRunes[lastTypedIndex] == rune(m.promptText[lastTypedIndex]) {
                    correctRune-- // Decrement if the last typed character was correct
                }

                // Remove the last rune from the user text
                currentRunes = currentRunes[:lastTypedIndex]
                m.userText.Reset()
                m.userText.WriteString(string(currentRunes))
            }
        } else if len(msg.Runes) > 0 { //other keys
            currentIndex := m.userText.Len()
            if currentIndex < len(m.promptText) && rune(m.promptText[currentIndex]) == msg.Runes[0]{
                correctRune++ // Increment correctRune here
            }
            m.userText.WriteRune(msg.Runes[0])

        }
        return m, nil

    case cursorTick:
        //Toggle the cursor visible
        m.cursorVisible = !m.cursorVisible
        return m, cursorCmd(450)

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }

    return m, nil
}

//View logic
func (m model) View() string {
    userInput := m.userText.String()
    prompt := m.promptText

    var renderedText strings.Builder
    // Inject cursor at the next untyped position
    cursorPosition := len(userInput)

    // Update final WPM only if typing is not complete
    m.finalWPM = calculateWPM(userInput, m.startTime)

    for i, r := range prompt {
        if i == cursorPosition && m.cursorVisible {
            renderedText.WriteString("|") // Add the cursor here
        }
        if i < len(userInput) {
            if rune(userInput[i]) == r {
                // Correct character (white)
                renderedText.WriteString(fmt.Sprintf("\033[97m%c\033[0m", r)) // White
            } else {
                // Incorrect character (red)
                renderedText.WriteString(fmt.Sprintf("\033[31m%c\033[0m", userInput[i])) // Red
            }
        } else {
            // Unentered character (gray prompt text)
            renderedText.WriteString(fmt.Sprintf("\033[90m%c\033[0m", r)) // Gray
        }
    }
    renderedText.WriteString("\n\n")
    // Add WPM to the output
    renderedText.WriteString(fmt.Sprintf("\033[97m%s\033[0m", fmt.Sprintf("WPM: %d", m.finalWPM)))

    renderedText.WriteString(fmt.Sprintf("\033[90m%s\033[0m","\n* Press any key to exit after finishing"))
    renderedText.WriteString(fmt.Sprintf("\033[90m%s\033[0m","\n* Ctrl-C or Ctrl-Z: exit any time\n"))

    return wordwrap.String(renderedText.String() , m.width)
}

//function to count wpm
func calculateWPM(userText string, startTime time.Time) int {
    elapsed := time.Since(startTime).Minutes() //time elapsed
    if elapsed == 0 {
        return 0
    }

    words := strings.Fields(userText) //split by white space
    length := len(words)

    return int(float64(length) / elapsed)
}

//function to calculate accuracy
func calculateAccuracy(userText string, correctRune int) float64 {
    totalRunes := []rune(userText)

    return (float64(correctRune) / float64(len(totalRunes)))*100
}

//Initilize utility flag
func init() {
    version = flag.Bool("version", false, "Check program version")
    language = flag.String("lang", "cpp", "Specify a language")
    list_languages = flag.Bool("listlang", false, "List all available languages")
}


func main() {
    //Call the flag and it will handle args
    flag.Parse()

    if *version {
        fmt.Printf("monkeycode version v0.1.0\n");
        os.Exit(0)
    }

    data, _ := os.ReadFile("test.txt")
    prompt := strings.TrimSpace(string(data))
    promptReplaceTab := strings.ReplaceAll(prompt, "\t", "    ")//replace tab into 4 spaces

    m := model{
        promptText:     promptReplaceTab,
        userText:       &strings.Builder{},
        cursorVisible:  true,
        cursorPosition: 0,
        startTime:      time.Now(),
        width:          80,
        height:         20,
        isComplete:     false,
        finalWPM:       0,
    }

    program := tea.NewProgram(m, tea.WithAltScreen())
    //run program
    exited, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }

    outsideModel := exited.(model)
    // Print the final WPM after the program exits
    if outsideModel.isComplete {
        fmt.Printf("\nTyping Complete! Final WPM: %d\n", outsideModel.finalWPM)
        fmt.Printf("Accuracy: %.2f%% \n", calculateAccuracy(outsideModel.userText.String(), correctRune))
    } else {
        fmt.Println("\nExited without completing.")
    }
}
