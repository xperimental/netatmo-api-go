package netatmo

const (
	baseURL  = "https://api.netatmo.net/"
	authURL  = baseURL + "oauth2/authorize"
	tokenURL = baseURL + "oauth2/token"

	deviceURL = baseURL + "/api/getstationsdata"

	homesURL      = baseURL + "api/homesdata"
	homeStatusURL = baseURL + "api/homestatus"
)
