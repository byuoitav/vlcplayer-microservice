package helpers

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

// GetStream returns the url of the srteam currently playing

func GetStream(player *vlc.Player) (string, error) {
	var stream string
	media, err := player.Media()
	if err != nil {
		log.Printf("error getting media: %s", err)
		return "", err
	}
	stream, err = media.Location()
	if err != nil {
		log.Printf("error getting stream: %s", err)
		return "", err
	}

	return stream, nil
}

// GetPlaybackStatus returns the status of the player

func GetPlaybackStatus(player *vlc.Player) (string, error) {

	playing := player.IsPlaying()
	state := player.WillPlay()

	if playing {
		return "playing", nil
	} else if state {
		return "media not finished or in error", nil
	} else {
		return "error playing media", nil
	}
}

// StopStream quits the vlc player

func StopStream(player *vlc.Player) error {

	err := player.Stop()
	if err != nil {
		log.Printf("error stopping player: %s", err)
		return err
	}
	player.Release()
	return nil
}

// SwitchStream switches the player output to a new stream

func SwitchStream(player *vlc.Player, stream string) (*vlc.Player, error) {

	// Initialize libVLC. Additional command line arguments can be passed in
	// to the libVLC by specifying them in the Init function.
	player.Release()

	if err := vlc.Init("-vvv"); err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer vlc.Release()

	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	media, err := player.LoadMediaFromURL(stream)
	if err != nil {
		log.Printf("error loading media: %s", err)
		return nil, err
	}
	defer media.Release()

	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	err = player.Play()
	if err != nil {
		log.Printf("error playing media: %s", err)
		return nil, err
	}
	return player, nil
}

// Volume returns the current volume level

func Volume(player *vlc.Player) (int, error) {

	volume, err := player.Volume()
	if err != nil {
		log.Printf("error getting volume: %s", err)
		return 0, err
	}

	return volume, nil
}

// SetVolume sets the volume level
func SetVolume(player *vlc.Player, volume int) error {

	err := player.SetVolume(volume)
	if err != nil {
		log.Printf("error setting volume: %s", err)
		return err
	}
	return nil
}

// Mute mutes the current player output

func Mute(player *vlc.Player) error {

	err := player.SetMute(true)
	if err != nil {
		log.Printf("error muting player: %s", err)
		return err
	}

	return nil
}

// Unmute unmutes the current player output

func Unmute(player *vlc.Player) error {

	err := player.SetMute(false)
	if err != nil {
		log.Printf("error unmuting player: %s", err)
		return err
	}

	return nil
}
