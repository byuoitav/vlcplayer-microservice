package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/byuoitav/vlcplayer-microservice/data"
	bolt "go.etcd.io/bbolt"
)

// Cache is a cache of data retrieved from the couch database storing the sreams and devices

const (
	_streamBucket = "streams"
	_deviceBucket = "devices"
)

type configService struct {
	configService data.ConfigService
	db            *bolt.DB
}

func New(cs data.ConfigService, path string) (data.ConfigService, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("uabble to open cache at %s: %w", path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(_streamBucket))
		if err != nil {
			return fmt.Errorf("unable to create stream bucket: %w", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(_deviceBucket))
		if err != nil {
			return fmt.Errorf("unable to create device bucket: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to initialize cachle:%w", err)
	}

	return &configService{
		configService: cs,
		db:            db,
	}, nil
}

func (c *configService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, err := c.configService.GetStreamConfig(ctx, streamURL)
	if err != nil {
		stream, cacheErr := c.streamConfigFromCache(ctx, streamURL)
		if cacheErr != nil {
			fmt.Errorf("unable to get stream from cache: %w", cacheErr)
			return stream, err
		}

		return stream, nil
	}

	if err := c.cacheStream(ctx, streamURL, stream); err != nil {
		fmt.Errorf("unable to cache stream: %w", err)
	}

	return stream, nil
}

func (c *configService) cacheStream(ctx context.Context, streamURL string, stream data.Stream) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_streamBucket))
		if b == nil {
			return fmt.Errorf("unable to get stream bucket (does not exsist)")
		}

		buf, err := json.Marshal(stream)
		if err != nil {
			return fmt.Errorf("unable to marshal stream: %w", err)
		}

		err = b.Put([]byte(streamURL), buf)
		if err != nil {
			return fmt.Errorf("unable to put stream in cache: %w", err)
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

func (c *configService) streamConfigFromCache(ctx context.Context, streamURL string) (data.Stream, error) {
	var stream data.Stream

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_streamBucket))
		if b == nil {
			return fmt.Errorf("unable to get stream bucket (does not exsist)")
		}

		buf := b.Get([]byte(streamURL))
		if buf == nil {
			return fmt.Errorf("unable to get stream from cache (stream not in cache)")
		}

		err := json.Unmarshal(buf, &stream)
		if err != nil {
			return fmt.Errorf("unable to unmarshal stream: %w", err)
		}

		return nil
	})
	if err != nil {
		return stream, err
	}

	return stream, nil
}

func (c *configService) GetDeviceConfig(ctx context.Context, hostname string) (data.Device, error) {
	device, err := c.configService.GetDeviceConfig(ctx, hostname)
	if err != nil {
		device, cacheErr := c.deviceConfigFromCache(ctx, hostname)
		if cacheErr != nil {
			fmt.Errorf("unable to get device from cache: %w", cacheErr)
			return device, err
		}

		return device, nil
	}

	if err := c.cacheDevice(ctx, hostname, device); err != nil {
		fmt.Errorf("unable to cache device %s: %w", hostname, err)
	}

	return device, nil
}

func (c *configService) cacheDevice(ctx context.Context, hostname string, device data.Device) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_deviceBucket))
		if b == nil {
			return fmt.Errorf("unable to get device bucket (does not exsist)")
		}

		buf, err := json.Marshal(device)
		if err != nil {
			return fmt.Errorf("unable to marshal device: %w", err)
		}

		err = b.Put([]byte(hostname), buf)
		if err != nil {
			return fmt.Errorf("unable to put device in cache: %w", err)
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

func (c *configService) deviceConfigFromCache(ctx context.Context, hostname string) (data.Device, error) {
	var device data.Device

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_deviceBucket))
		if b == nil {
			return fmt.Errorf("unable to get device bucket (does not exsist)")
		}

		buf := b.Get([]byte(hostname))
		if buf == nil {
			return fmt.Errorf("unable to get device from cache (device not in cache)")
		}

		err := json.Unmarshal(buf, &device)
		if err != nil {
			return fmt.Errorf("unable to unmarshal device: %w", err)
		}

		return nil
	})
	if err != nil {
		return device, err
	}

	return device, nil
}
