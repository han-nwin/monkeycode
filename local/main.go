package main

import (
    "fmt"
    "os"
    "github.com/han-nwin/monkeycode/tui"
)

func main() {
    program := tui.NewProgram()

    //run program
    _, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }
/**
    outsideModel := exited.(tui.Model)
    // Print the final WPM after the program exits
    if outsideModel.IsComplete {
        fmt.Printf("\nTyping Complete! WPM: %d\n", outsideModel.FinalWPM)
        fmt.Printf("Accuracy: %.2f%% \n", tui.CalculateAccuracy(outsideModel.UserText.String(), tui.CorrectRune))
    } else {
        fmt.Println("\nExited without completing.")
    }
*/
}
