
import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/han-nwin/monkeycode/tui"
)

const profilesFile = "profiles.json"

// Profile represents a user profile with optional progress tracking
type Profile struct {
	Username      string   `json:"username"`
	TypingSpeeds  []int    `json:"typing_speeds"`
	AccuracyRates []float64 `json:"accuracy_rates"`
}

// Profiles contains all user profiles and metadata
type Profiles struct {
	LastActiveProfile string    `json:"last_active_profile"`
	Users             []Profile `json:"users"`
}

var (
	usernameFlag = flag.String("user", "", "Specify the username to use or create a new profile")
)

// LoadProfiles loads profiles from the profiles.json file
func LoadProfiles(filepath string) (*Profiles, error) {
	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		// Return an empty profiles structure if the file doesn't exist
		return &Profiles{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles Profiles
	err = json.NewDecoder(file).Decode(&profiles)
	if err != nil {
		return nil, err
	}

	return &profiles, nil
}

// SaveProfiles saves profiles to the profiles.json file
func SaveProfiles(filepath string, profiles *Profiles) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(profiles)
}

// GetProfile retrieves an existing profile or creates a new one
func GetProfile(profiles *Profiles, username string) *Profile {
	for i := range profiles.Users {
		if profiles.Users[i].Username == username {
			return &profiles.Users[i]
		}
	}

	// Create a new profile if not found
	newProfile := Profile{
		Username:      username,
		TypingSpeeds:  []int{},
		AccuracyRates: []float64{},
	}
	profiles.Users = append(profiles.Users, newProfile)
	return &profiles.Users[len(profiles.Users)-1]
}

// GetLastActiveProfile retrieves the last active profile if available
func GetLastActiveProfile(profiles *Profiles) *Profile {
	if profiles.LastActiveProfile == "" {
		return nil
	}

	for i := range profiles.Users {
		if profiles.Users[i].Username == profiles.LastActiveProfile {
			return &profiles.Users[i]
		}
	}

	return nil
}

func main() {
	flag.Parse()

	// Load profiles
	profiles, err := LoadProfiles(profilesFile)
	if err != nil {
		fmt.Printf("Error loading profiles: %v\n", err)
		os.Exit(1)
	}

	var profile *Profile

	// Determine the active profile
	if *usernameFlag != "" {
		// Use the profile specified by the flag
		profile = GetProfile(profiles, *usernameFlag)
		profiles.LastActiveProfile = profile.Username
	} else {
		// Use the last active profile
		profile = GetLastActiveProfile(profiles)
		if profile == nil {
			fmt.Println("No active profile found. Use the -user flag to create or select a profile.")
			os.Exit(1)
		}
	}

	fmt.Printf("Welcome, %s!\n", profile.Username)

	// Run the typing test program
	program := tui.NewProgram()
	_, err = program.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Example: Update profile with dummy data after the program
	profile.TypingSpeeds = append(profile.TypingSpeeds, 60) // Dummy speed
	profile.AccuracyRates = append(profile.AccuracyRates, 98.5) // Dummy accuracy

	// Save profiles
	err = SaveProfiles(profilesFile, profiles)
	if err != nil {
		fmt.Printf("Error saving profiles: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Profile updated successfully!")
}
