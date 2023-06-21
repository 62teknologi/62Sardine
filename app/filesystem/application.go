package filesystem

import (
	"context"
	"fmt"

	"github.com/62teknologi/62sardine/config"

	"github.com/gookit/color"

	"github.com/goravel/framework/contracts/filesystem"
)

type Driver string

const (
	DriverLocal Driver = "local"
	DriverS3    Driver = "s3"
	DriverOss   Driver = "oss"
	DriverCos   Driver = "cos"
	DriverMinio Driver = "minio"
	DriverGCS   Driver = "gcs"
)

type Storage struct {
	filesystem.Driver
	drivers map[string]filesystem.Driver
}

func NewStorage(contentType string, visibility string) *Storage {

	defaultDisk, _ := config.ReadConfig("filesystems.default")
	if defaultDisk == "" {
		color.Redln("[filesystem] please set default disk")

		return nil
	}

	driver, err := NewDriver(defaultDisk, contentType, visibility == "public")
	if err != nil {
		color.Redf("[filesystem] %s\n", err)

		return nil
	}

	drivers := make(map[string]filesystem.Driver)
	drivers[defaultDisk] = driver
	return &Storage{
		Driver:  driver,
		drivers: drivers,
	}
}

func NewDriver(disk string, contentType string, isPublic bool) (filesystem.Driver, error) {
	ctx := context.Background()
	s, _ := config.ReadConfig(fmt.Sprintf("filesystems.disks.%s.driver", disk))
	driver := Driver(s)
	switch driver {
	case DriverLocal:
		return NewLocal(disk)
	case DriverOss:
		return NewOss(ctx, disk, contentType, isPublic)
	case DriverCos:
		return NewCos(ctx, disk)
	case DriverS3:
		return NewS3(ctx, disk, contentType, isPublic)
	case DriverMinio:
		return NewMinio(ctx, disk)
	case DriverGCS:
		return NewGCS(ctx, disk, contentType, isPublic)

	}

	return nil, fmt.Errorf("invalid driver: %s, only support local, s3, oss, cos, minio", driver)
}

// func (r *Storage) Disk(disk string) filesystem.Driver {
// 	if driver, exist := r.drivers[disk]; exist {
// 		return driver
// 	}

// 	driver, err := NewDriver(disk)
// 	if err != nil {
// 		panic(err)
// 	}

// 	r.drivers[disk] = driver

// 	return driver
// }
