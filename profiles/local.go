package profiles

import (
    "encoding/json"
    "fmt"
    "os"
)

type LocalProfileBackend struct {
    ProfilesDir       string // Directory containing profile files
    LastActiveProfile string // Path to the metadata file for the last active profile
}

// NewLocalProfileBackend creates a new instance of LocalProfileBackend
func NewLocalProfileBackend(profilesDir, lastActiveProfile string) *LocalProfileBackend {
    return &LocalProfileBackend{
        ProfilesDir:       profilesDir,
        LastActiveProfile: lastActiveProfile,
    }
}

// LoadProfile retrieves a specific profile or the last active one if username is empty
func (l *LocalProfileBackend) LoadProfile(username string) (*Profile, error) {
    if username == "" {
        // Read the last active profile from the metadata file
        data, err := os.ReadFile(l.LastActiveProfile)
        if err != nil {
            return nil, fmt.Errorf("failed to load last active profile: %v", err)
        }
        username = string(data)
    }

    // Read the profile file
    profilePath := fmt.Sprintf("%s/%s.json", l.ProfilesDir, username)
    file, err := os.Open(profilePath)
    if os.IsNotExist(err) {
        return nil, fmt.Errorf("profile '%s' not found", username)
    } else if err != nil {
        return nil, err
    }
    defer file.Close()

    var profile Profile
    if err := json.NewDecoder(file).Decode(&profile); err != nil {
        return nil, err
    }
    return &profile, nil
}

// SaveProfile saves a single profile to its file
func (l *LocalProfileBackend) SaveProfile(profile *Profile) error {
    profilePath := fmt.Sprintf("%s/%s.json", l.ProfilesDir, profile.Username)
    file, err := os.Create(profilePath)
    if err != nil {
        return err
    }
    defer file.Close()

    return json.NewEncoder(file).Encode(profile)
}

// UpdateLastActiveProfile updates the last active profile metadata
func (l *LocalProfileBackend) UpdateLastActiveProfile(username string) error {
    return os.WriteFile(l.LastActiveProfile, []byte(username), 0644)
}
