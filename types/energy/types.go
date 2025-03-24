package energy

type HomesDataResponse struct {
	Status string `json:"status"`
	Body   struct {
		Homes []Home `json:"homes,omitempty"`
		User  User   `json:"user,omitempty"`
	} `json:"body,omitempty"`
}

type HomeStatusResponse struct {
	Status string `json:"status"`
	Body   struct {
		Home   Home          `json:"home,omitempty"`
		Errors []ModuleError `json:"errors,omitempty"`
	}
}

type Home struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Altitude    uint32    `json:"altitude,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Country     string    `json:"country,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	Rooms       []Room    `json:"rooms,omitempty"`
	Modules     []Module  `json:"modules,omitempty"`
}

type Room struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      string   `json:"type,omitempty"`
	ModuleIDs []string `json:"module_ids,omitempty"`
}

type Module struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	SetupDate int64  `json:"setup_date,omitempty"`
	RoomID    string `json:"room_id,omitempty"`
}

type User struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	Language          string `json:"language"`
	Locale            string `json:"locale"`
	FeelLikeAlgorithm uint8  `json:"feel_like_algorithm"`
	UnitPressure      uint8  `json:"unit_pressure"`
	UnitSystem        uint8  `json:"unit_system"`
	UnitWind          uint8  `json:"unit_wind"`
}

type ModuleError struct {
	ID   string `json:"id,omitempty"`
	Code int64  `json:"code,omitempty"`
}
