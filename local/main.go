package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/han-nwin/monkeycode/tui"
	"github.com/han-nwin/monkeycode/profiles"
)

var (
    usernameFlag = flag.String("user", "", "Specify a user name")
)

type Profile struct {
    Username string `json:"username"`
    TypingSpeed []int `json:"typingspeed"`
    Accuracy []float64 `json:"accuracy"`
}

type Profiles struct {
    Users []Profile `json:"users"`
    LastActiveProfile Profile `json:"last_active_profile"`
}

func main() {
    flag.Parse()
    
    //TODO: If no -user provided, check if there's a lastActiveProfile
    if *usernameFlag != "" {
        profiles.
    } else { //if user flag provided load or create new profile

    }

    time.Sleep(0*time.Second) //Sleep for 2 sec before running TUI
    //CREATE and RUN TUI program
    program := tui.NewProgram()
    //run program
    exitModel, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }
    //== Program EXIT
    // Use type assertion to access the custom fields
	if finalModel, ok := exitModel.(tui.Model); ok {
		fmt.Printf("Final WPM: %d\n", finalModel.WPM)
		fmt.Printf("Accuracy: %.2f%%\n", finalModel.Accuracy)
		fmt.Printf("Saved WPM: %d\n", finalModel.SavedWPM)
		fmt.Printf("Saved Accuracy: %.2f%%\n", finalModel.SavedAccuracy)
	} else {
		fmt.Println("Failed to assert the returned model to the custom Model type.")
		os.Exit(-1)
	}
    
}
