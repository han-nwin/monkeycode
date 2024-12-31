package profiles


//== Profile structure ==//
type Profile struct {
    Username string `json:"username"`
    TypingSpeed []int `json:"typingspeed"`
    Accuracy []float64 `json:"accuracy"`
}

// ALL Profiles 
type Profiles struct {
    Users []Profile `json:"users"`
    LastActiveProfile Profile `json:"last_active_profile"`
}

//== Backend Interface for both Local and SSH ==//
type ProfilesBackEnd interface {
    SaveProfiles(*Profiles) error // Save Profiles struct
    LoadProfile(userName string) (*Profile, error) // Load 0 or 1 Profile
    LoadALLProfiles() //Use for leaderboard functionality
}
