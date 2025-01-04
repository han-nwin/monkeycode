package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/han-nwin/monkeycode/tui"
	"github.com/han-nwin/monkeycode/profiles"
)

var (
    usernameFlag = flag.String("user", "", "Specify a user name")
)


func main() {
    flag.Parse()

    // Initialize the backend
    backend := profiles.NewLocalProfileBackend()

    // Declare variables
    var (
        profile *profiles.Profile
        err     error
    )

    // Check if a username is provided
    if *usernameFlag != "" {
        profile, err = backend.LoadProfile(*usernameFlag)
        if err != nil {
            fmt.Printf("Creating '%s' profile.\n", *usernameFlag)
            profile = &profiles.Profile{
                Username:    *usernameFlag,
                TypingSpeed: []int{},
                Accuracy:    []float64{},
            }
            err = backend.SaveProfile(profile)
            if err != nil {
                fmt.Printf("Error saving profile: %v\n", err)
                return
            }
            fmt.Println("New profile created.")
        } else {
            fmt.Printf("Loaded profile: '%v'\n", profile.Username)
        }

        err = backend.SetLastActiveProfile(profile.Username)
        if err != nil {
            fmt.Printf("Error setting last active profile: %v\n", err)
            return
        }

    } else {
        var allProfiles *profiles.Profiles
        allProfiles, err = backend.LoadAllProfiles()
        if err != nil {
            fmt.Printf("Error loading profiles: %v\n", err)
            return
        }

        if allProfiles.LastActiveProfile == "" {
            fmt.Println("No last active profile found. Please specify a new user with: \nmonkeycode -user <username>")
            return
        }

        profile, err = backend.LoadProfile("")
        if err != nil {
            fmt.Printf("Error loading last active profile: %v\n", err)
            return
        }

        fmt.Printf("Last active profile '%v' loaded\n", profile.Username)
    }

    // Run TUI Program
    program := tui.NewProgram()
    exitModel, err := program.Run()
    if err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(-1)
    }

    // Program Exit Logic
    if finalModel, ok := exitModel.(tui.Model); ok {
        fmt.Printf("Final WPM: %d\n", finalModel.WPM)
        fmt.Printf("Accuracy: %.2f%%\n", finalModel.Accuracy)
        fmt.Printf("Saved WPM: %d\n", finalModel.SavedWPM)
        fmt.Printf("Saved Accuracy: %.2f%%\n", finalModel.SavedAccuracy)

        // Append WPM and Accuracy to the profile
        profile.TypingSpeed = append(profile.TypingSpeed, finalModel.SavedWPM...)
        profile.Accuracy = append(profile.Accuracy, finalModel.SavedAccuracy...)

        // Save the updated profile
        err = backend.SaveProfile(profile)
        if err != nil {
            fmt.Printf("Error saving profile: %v\n", err)
            os.Exit(-1)
        }

        fmt.Printf("Profile '%v' updated successfully!\n", profile.Username)
    } else {
        fmt.Println("Failed to assert the returned model to the custom Model type.")
        os.Exit(-1)
    }
}
