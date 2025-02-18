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
    usernameFlag = flag.String("u", "", "Specify a user name")
    leaderboardFlag = flag.Bool("l", false, "Display leaderboard.")
    currentUserFlag = flag.Bool("c" , false, "Display current user.")
)


func main() {
    flag.Parse()

    // Declare variables
    var (
        allProfiles *profiles.Profiles
        profile *profiles.Profile
        err     error
    )
    
    // Initialize the backend
    backend := profiles.NewLocalProfileBackend()

	// Display leaderboard if the flag is set
	if *leaderboardFlag {
		err := profiles.DisplayLeaderboard(backend, 100) //Display 100 profiles
		if err != nil {
			fmt.Printf("Error displaying leaderboard: %v\n", err)
		}
		return
	}

    //Show last active profile if currentUserFlag is set
    if *currentUserFlag {
        allProfiles, err = backend.LoadAllProfiles();
        lastActiveUserName := allProfiles.LastActiveProfile;
        fmt.Printf("Current(last active) user: %s \n", lastActiveUserName)
        
        profile,_ = backend.LoadProfile(lastActiveUserName)
        fmt.Printf("Average WPM: %.2f \n", profile.AverageTypingSpeed())
        fmt.Printf("Average Accuracy: %.2f%%\n", profile.AverageAccuracy())

        return
    }

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
        allProfiles, err = backend.LoadAllProfiles()
        if err != nil {
            fmt.Printf("Error loading profiles: %v\n", err)
            return
        }

        if allProfiles.LastActiveProfile == "" {
            fmt.Println("No last active profile found. Please specify a new user with: \nmonkeycode -u <username>")
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
        fmt.Printf("Saved WPM: %d\n", finalModel.SavedWPM)
        fmt.Printf("Saved Accuracy: %.2f%%\n", finalModel.SavedAccuracy)
        fmt.Println();

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
