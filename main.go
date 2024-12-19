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
    currentText := []rune(m.userText.String()) // Get the user text as runes for precise handling

    var renderedText string
    if m.cursorVisible {
        // Dynamically insert the cursor at the end of the text
        renderedText = string(append(currentText, '█'))
    } else {
        renderedText = string(currentText)
    }

    // Render the prompt text and the user text
    return wordwrap.String(m.promptText + "\n" + renderedText + "\n", m.width)
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
        fmt.Printf("monkeycode version 0.1.0\n");
        os.Exit(0)
    }
    
    fmt.Println(*language, *version);

    data, _ := os.ReadFile("test.txt")
    prompt := strings.Split(strings.TrimSpace(string(data)), "\n")

    for _, a := range prompt {
        println(a)
    }

    m := model{
        promptText:     prompt[1],
        userText:       &strings.Builder{},
        cursorVisible:  true,
    }

    program := tea.NewProgram(m, tea.WithAltScreen())
    //run program
    _, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }
}
