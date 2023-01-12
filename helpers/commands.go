package helpers

const ()

// GetStream returns the url of the srteam currently playing

func GetStream() (string, error) {

	return "", nil
}

// GetPlaybackStatus returns the status of the player

func GetPlaybackStatus() (string, error) {

	return "", nil
}

// StopStream quits the vlc player

func StopStream() error {

	return nil
}

// SwitchStream switches the player output to a new stream

func SwitchStream(stream string) error {

	return nil
}

// VolumeControl always returns the current volume and optionally can change the volume

func VolumeControl(volume int) (int, error) {

	return 0, nil
}

// Mute mutes the current player output

func Mute() error {

	return nil
}

// Unmute unmutes the current player output

func Unmute() error {

	return nil
}
