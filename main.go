package main

import (
    "fmt"
    "flag"
    "os"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/muesli/reflow/wordwrap"
)

var (
    language *string
    list_languages *bool
    version *bool
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
}

//Init model
func (m model) Init() tea.Cmd {
    return cursorCmd(450)
}

//Update logic
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        //Exit con Ctrl+C
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
        
        //On every keystroke, add one character to userText
        if msg.Type == tea.KeyEnter { //type Enter
            m.userText.WriteRune('\n')

        } else if msg.Type == tea.KeyTab { //type Tab
            m.userText.WriteString("    ")

        } else if msg.Type == tea.KeyBackspace { //type Backspace
            currentText := m.userText.String()
            runes := []rune(currentText)
            if len(runes) > 0 {
                runes = runes[:len(runes)-1]
            }
            m.userText.Reset()
            m.userText.WriteString(string(runes))

        } else if len(msg.Runes) > 0 { //other keys
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
    //
    // Add the cursor at the position corresponding to userText length
    if len(userInput) < len(prompt) {

        // Inject cursor at the next untyped position
        cursorPosition := len(userInput)
        renderedText.Reset() // Rebuild the string with the cursor inserted

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
        if cursorPosition > len(prompt) {
            renderedText.WriteString("|")
        }
    }
    return wordwrap.String(renderedText.String() , m.width)
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
    
    //data, _ := os.ReadFile("test.txt")
    //prompt := strings.TrimSpace(string(data))

    m := model{
        promptText:     "printf(\"Hello World\")",
        userText:       &strings.Builder{},
        cursorVisible:  true,
        cursorPosition: 0,
    }

    program := tea.NewProgram(m, tea.WithAltScreen())
    //run program
    _, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }
}
