package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LocalProfileBackend struct {
	ProfilesFile string // Path to the profiles.json file
}

//Helper function to get/create profiles.json in user home
// ~/.config/monkeycode/profiles.json
func getProfilesFile() string {
    homeDir, err := os.UserHomeDir() //return $HOME on UNIX. %USERPROFILE% on Windows
	if err != nil {
		panic(fmt.Sprintf("failed to get user config directory: %v", err))
	}

	// Construct the path to profiles.json
	return filepath.Join(homeDir, ".config", "monkeycode", "profiles.json")
}

//Create new instant of LocalProfileBackend
//Check if the profileDir exist -> if not create it
func NewLocalProfileBackend() *LocalProfileBackend{
    profilesFile := getProfilesFile();

    // Ensure the parent directory exists
	dir := filepath.Dir(profilesFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755) // Create the directory
		if err != nil {
			panic(fmt.Sprintf("failed to create directory: %v", err))
		}
	}

	// Ensure the profiles.json file exists
	if _, err := os.Stat(profilesFile); os.IsNotExist(err) {
		// Create the file with default content
		defaultProfiles := Profiles{
			Users:             []Profile{},
			LastActiveProfile: "",
		}

		file, err := os.Create(profilesFile)
		if err != nil {
			panic(fmt.Sprintf("failed to create profiles file: %v", err))
		}
		defer file.Close()

		// Write the default content to the file
		err = json.NewEncoder(file).Encode(defaultProfiles)
		if err != nil {
			panic(fmt.Sprintf("failed to write default profiles: %v", err))
		}
	}

    return &LocalProfileBackend{
        ProfilesFile: profilesFile,
    }
}


//Load ALL profiles from the file
func (l *LocalProfileBackend) LoadAllProfiles() (*Profiles, error) {
	file, err := os.Open(l.ProfilesFile)
	if os.IsNotExist(err) {
		// Return an empty Profiles struct if the file doesn't exist
		return &Profiles{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles Profiles
	if err := json.NewDecoder(file).Decode(&profiles); err != nil {
		return nil, err
	}
	return &profiles, nil
}

//Load a single profile
func (l *LocalProfileBackend) LoadProfile(username string) (*Profile, error) {
	// Load all profiles
	profiles, err := l.LoadAllProfiles()
	if err != nil {
		return nil, err
	}

	// If no username is provided, use the last active profile
	if username == "" {
		username = profiles.LastActiveProfile
		if username == "" {
			return nil, fmt.Errorf("no profile name provided and no last active profile set")
		}
	}

	// Find the requested profile
	for _, profile := range profiles.Users {
		if profile.Username == username {
			return &profile, nil
		}
	}

	return nil, fmt.Errorf("profile '%s' not found", username)
}

//Save ALL profiles to the file
func (l *LocalProfileBackend) SaveAllProfiles(profiles *Profiles) error {
	file, err := os.Create(l.ProfilesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(profiles)
}

//Save a single profile
func (l *LocalProfileBackend) SaveProfile(profile *Profile) error {
	// Load all profiles
	profiles, err := l.LoadAllProfiles()
	if err != nil {
		return err
	}

	// Check if the profile already exists
	found := false
	for i, p := range profiles.Users {
		if p.Username == profile.Username {
			profiles.Users[i] = *profile // Update existing profile
			found = true
			break
		}
	}

	if !found {
		// Add new profile
		profiles.Users = append(profiles.Users, *profile)
	}

	// Save all profiles back to the file
	return l.SaveAllProfiles(profiles)
}


// SetLastActiveProfile updates the LastActiveProfile in profiles.json
func (l *LocalProfileBackend) SetLastActiveProfile(username string) error {
	// Load all profiles
	profiles, err := l.LoadAllProfiles()
	if err != nil {
		return err
	}

	// Update the LastActiveProfile field
	profiles.LastActiveProfile = username

	// Save the updated profiles back to the file
	return l.SaveAllProfiles(profiles)
}
