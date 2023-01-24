package helpers

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

// StartVLC starts the vlc player

func StartVLC() (*vlc.Player, error) {

	// Initialize libVLC. Additional command line arguments can be passed in
	// to the libVLC by specifying them in the Init function.

	if err := vlc.Init(); err != nil {
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

	return player, nil
}
