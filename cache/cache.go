package cache

import (
	"github.com/byuoitav/vlcplayer-microservice/data"
	bolt "go.etcd.io/bbolt"
)

type configService struct {
	configService data.ConfigService
	db            *bolt.DB
}

/*func (c *configService) GetStreamConfig(ctx context.Context, streamURL string) (data.ConfigService, error) {

	db, err := bolt.Open(streamURL, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("stream"))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("device"))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create bucket: %w", err)
	}

}*/
