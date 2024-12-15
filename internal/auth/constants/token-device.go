package constants

type TokenDevice string

const (
	DeviceApi TokenDevice = "api"
)

func TokenDeviceValues() []TokenDevice {
	return []TokenDevice{DeviceApi}
}

func IsValidTokenDevice(device string) bool {
	switch TokenDevice(device) {
	case DeviceApi:
		return true
	default:
		return false
	}
}
