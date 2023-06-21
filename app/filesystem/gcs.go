package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/62teknologi/62sardine/config"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/support/str"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCS struct {
	ctx        context.Context
	instance   *storage.Client
	disk       string
	bucketName string
	isPublic   bool
	contenType string
	url        string
	Email      string
	PrivateKey string
}

type GCSJSON struct {
	Email      string `json:"client_email"`
	PrivateKey string `json:"private_key"`
}

func NewGCS(ctx context.Context, disk string, contenType string, isPublic bool) (*GCS, error) {
	path, _ := config.ReadConfig(fmt.Sprintf("filesystems.disks.%s.path", disk))
	bucketName, _ := config.ReadConfig(fmt.Sprintf("filesystems.disks.%s.bucket", disk))
	url, _ := config.ReadConfig(fmt.Sprintf("filesystems.disks.%s.url", disk))

	creds := option.WithCredentialsFile(path)

	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		return &GCS{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return &GCS{}, err
	}
	defer file.Close()

	var js GCSJSON
	err = json.NewDecoder(file).Decode(&js)
	if err != nil {
		return &GCS{}, err
	}

	return &GCS{
		ctx:        ctx,
		instance:   client,
		bucketName: bucketName,
		disk:       disk,
		contenType: contenType,
		isPublic:   isPublic,
		url:        url,
		Email:      js.Email,
		PrivateKey: js.PrivateKey,
	}, err

}

func (r *GCS) AllDirectories(path string) ([]string, error) {
	var dirs []string

	return dirs, nil
}

func (r *GCS) AllFiles(path string) ([]string, error) {
	var files []string

	bucket := r.instance.Bucket(r.bucketName)

	it := bucket.Objects(r.ctx, &storage.Query{Prefix: path})

	for {
		objectAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}
		files = append(files, objectAttrs.Name)
	}

	return files, nil

}

func (r *GCS) Copy(oldFile, newFile string) error {

	return errors.New("not implemneted yet")
}

func (r *GCS) Delete(file ...string) error {
	bucket := r.instance.Bucket(r.bucketName)
	for _, fileName := range file {
		obj := bucket.Object(fileName)
		err := obj.Delete(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GCS) DeleteDirectory(directory string) error {
	return errors.New("not implemneted yet")
}

func (r *GCS) Directories(path string) ([]string, error) {
	var dirs []string
	return dirs, errors.New("not implemneted yet")
}

func (r *GCS) Exists(file string) bool {
	bucket := r.instance.Bucket(r.bucketName)
	_, err := bucket.Object(file).Attrs(r.ctx)

	if err == storage.ErrObjectNotExist {
		return false
	} else if err != nil {
		return false
	}

	return true
}

func (r *GCS) Files(path string) ([]string, error) {
	var files []string
	return files, errors.New("not implemneted yet")
}

func (r *GCS) Get(file string) (string, error) {
	return "", errors.New("not implemneted yet")
}

func (r *GCS) MakeDirectory(directory string) error {
	return errors.New("not implemneted yet")
}

func (r *GCS) Missing(file string) bool {
	return !r.Exists(file)
}

func (r *GCS) Move(oldFile, newFile string) error {
	return errors.New("not implemneted yet")
}

func (r *GCS) Path(file string) string {
	return file
}

func (r *GCS) Put(file, content string) error {
	bucket := r.instance.Bucket(r.bucketName)
	obj := bucket.Object(file)

	w := obj.NewWriter(r.ctx)
	if _, err := io.Copy(w, strings.NewReader(content)); err != nil {
		return err
	}
	err := w.Close()
	if err := w.Close(); err != nil {
		return err
	}

	if r.isPublic {
		if err := obj.ACL().Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}
	}

	return err

}

func (r *GCS) PutFile(filePath string, source filesystem.File) (string, error) {
	return r.PutFileAs(filePath, source, str.Random(40))
}

func (r *GCS) PutFileAs(filePath string, source filesystem.File, name string) (string, error) {
	fullPath, err := fullPathOfFile(filePath, source, name)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(source.File())
	if err != nil {
		return "", err
	}

	if err := r.Put(fullPath, string(data)); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (r *GCS) Size(file string) (int64, error) {
	bucket := r.instance.Bucket(r.bucketName)
	obj := bucket.Object(file)
	attrs, err := obj.Attrs(context.Background())
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return attrs.Size, err
}

func (r *GCS) TemporaryUrl(file string, time time.Time) (string, error) {

	url, err := storage.SignedURL(r.bucketName, file, &storage.SignedURLOptions{
		GoogleAccessID: r.Email,
		PrivateKey:     []byte(r.PrivateKey),
		Method:         "GET",
		Expires:        time,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return url, nil
}

func (r *GCS) WithContext(ctx context.Context) filesystem.Driver {
	driver, err := NewGCS(ctx, r.disk, "", false)
	if err != nil {
		fmt.Errorf("init %s disk fail: %+v", r.disk, err)
	}

	return driver
}

func (r *GCS) Url(file string) string {
	return strings.TrimSuffix(r.url, "/") + "/" + strings.TrimPrefix(file, "/")
}
