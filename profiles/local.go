package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

    "github.com/charmbracelet/lipgloss"
)

type LocalProfileBackend struct {
	ProfilesFile string // Path to the profiles.json file
}

//NOTE: Local simple database

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

//NOTE: PROFILES functionatlity

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

//NOTE: LEADERBOARD functionality

// Calculate average typing speed for the profile
func (p *Profile) AverageTypingSpeed() float64 {
	if len(p.TypingSpeed) == 0 {
		return 0
	}
	sum := 0
	for _, speed := range p.TypingSpeed {
		sum += speed
	}
	return float64(sum) / float64(len(p.TypingSpeed))
}

// Calculate average accuracy for the profile
func (p *Profile) AverageAccuracy() float64 {
	if len(p.Accuracy) == 0 {
		return 0
	}
	sum := 0.0
	for _, acc := range p.Accuracy {
		sum += acc
	}
	return sum / float64(len(p.Accuracy))
}

// Calculate the combined score based on average WPM and accuracy
func (p *Profile) Score() float64 {
	averageWPM := p.AverageTypingSpeed()
	averageAccuracy := p.AverageAccuracy()

	// Weighted formula: 70% WPM, 30% Accuracy
	return (averageWPM * 0.7) + (averageAccuracy * 0.3)
}

// Define Catppuccin theme colors
var (
	Flamingo  = lipgloss.Color("#F2CDCD")
	Pink      = lipgloss.Color("#F5C2E7")
	Mauve     = lipgloss.Color("#DDB6F2")
	Red       = lipgloss.Color("#F28FAD")
	Maroon    = lipgloss.Color("#E8A2AF")
	Peach     = lipgloss.Color("#F8BD96")
	Yellow    = lipgloss.Color("#FAE3B0")
	Green     = lipgloss.Color("#ABE9B3")
	Teal      = lipgloss.Color("#B5E8E0")
	Sky       = lipgloss.Color("#96CDFB")
	Blue      = lipgloss.Color("#89DCEB")
	Lavender  = lipgloss.Color("#C9CBFF")
	Base      = lipgloss.Color("#1E1E2E") // Background
	Surface   = lipgloss.Color("#313244") // Surface
	Text      = lipgloss.Color("#CDD6F4") // Text
	Subtext   = lipgloss.Color("#9399B2") // Subtext
)
// Define styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(Mauve).
		Background(Surface).
		Bold(true).
		Padding(1, 2).
		Align(lipgloss.Center)

	headerStyle = lipgloss.NewStyle().
			Foreground(Mauve).
			Background(Surface).
			Padding(0, 2).
			Bold(true)

	rowStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 2)

	usernameStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	wpmStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	accuracyStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)
)
//Display Leaderboard
//We can specify how many profiles to display
func DisplayLeaderboard(backend *LocalProfileBackend, topN int) error {
	// Load all profiles
	profiles, err := backend.LoadAllProfiles()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %v", err)
	}

	// Sort profiles by average WPM (descending order)
	sort.Slice(profiles.Users, func(i, j int) bool {
		return profiles.Users[i].Score() > profiles.Users[j].Score()
	})

    // Define column widths
	usernameWidth := 15
	wpmWidth := 15
	accuracyWidth := 20

    // Styled title
    totalWidth := usernameWidth + wpmWidth + accuracyWidth
	title := titleStyle.Width(totalWidth).Render(fmt.Sprintf("Leaderboard: Top %d", topN))

	// Styled header
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Width(usernameWidth).Render("Username"),
		headerStyle.Width(wpmWidth).Render("Average WPM"),
		headerStyle.Width(accuracyWidth).Render("Average Accuracy"),
	)

	// Styled rows
	var rows []string
	for i, user := range profiles.Users {
		if i >= topN {
			break
		}
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			usernameStyle.Width(usernameWidth).Render(user.Username),
			wpmStyle.Width(wpmWidth).Render(fmt.Sprintf("%.2f", user.AverageTypingSpeed())),
			accuracyStyle.Width(accuracyWidth).Render(fmt.Sprintf("%.2f", user.AverageAccuracy())),
		)
		rows = append(rows, rowStyle.Render(row))
	}

	// Render leaderboard
	fmt.Println(title)
	fmt.Println(header)
	for _, row := range rows {
		fmt.Println(row)
	}

	return nil
}
