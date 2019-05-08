package homlet

//go:generate enumer -type=HumidityStatus -output=humidity_gen.go -transform=snake -trimprefix=Humidity

// HumidityStatus represents humidity status
type HumidityStatus int

// Supported humidity statuses
const (
	HumidityNormal HumidityStatus = iota
	HumidityComfortable
	HumidityDry
	HumidityWet
)
