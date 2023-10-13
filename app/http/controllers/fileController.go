package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/62teknologi/62sardine/app/filesystem"
	"github.com/62teknologi/62sardine/config"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	sfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type FileController struct {
	Adapter string
}

func ResizeImage(c *gin.Context, fileHeader multipart.FileHeader, width string, height string, tempFileName string) (*multipart.FileHeader, error) {
	// Open the uploaded file.
	srcFile, err := fileHeader.Open()
	if err != nil {
		if srcFile, err = os.Open(fileHeader.Filename); err != nil {
			return nil, err
		}
	}
	defer srcFile.Close()

	// Decode the uploaded image.
	srcImage, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, err
	}

	// Parse the width and height parameters.
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		widthInt = 0
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		heightInt = 0
	}

	// Resize the image.
	resizedImage := imaging.Resize(srcImage, widthInt, heightInt, imaging.Lanczos)

	// Create a new in-memory buffer to store the resized image.
	buf := new(bytes.Buffer)

	// Encode the resized image to the buffer in JPEG format.
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}

	// Write the resized image to disk.
	err = ioutil.WriteFile(tempFileName, buf.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	fileHeader.Header.Set("Content-Type", "image/jpeg")

	// Create a new multipart.FileHeader for the resized image.
	resizedFileHeader := &multipart.FileHeader{
		Filename: tempFileName,
		Size:     int64(buf.Len()),
		Header:   fileHeader.Header,
	}

	return resizedFileHeader, nil
}

func CompressImage(fileHeader *multipart.FileHeader, quality string, tempFileName string) (*multipart.FileHeader, error) {
	qualityInt, err := strconv.Atoi(quality)
	if err != nil {
		return nil, err
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		if srcFile, err = os.Open(fileHeader.Filename); err != nil {
			return nil, err
		}
	}
	defer srcFile.Close()

	// Decode the uploaded image.
	srcImage, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, err
	}

	// Create a new in-memory buffer to store the resized image.
	buf := new(bytes.Buffer)

	// Encode the image to the desired quality.
	err = jpeg.Encode(buf, srcImage, &jpeg.Options{
		Quality: qualityInt,
	})
	if err != nil {
		return nil, err
	}

	// Write the resized image to disk.
	err = ioutil.WriteFile(tempFileName, buf.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	fileHeader.Header.Set("Content-Type", "image/jpeg")

	// Return the compressed image.
	return &multipart.FileHeader{
		Filename: tempFileName,
		Size:     int64(buf.Len()),
		Header:   fileHeader.Header,
	}, nil
}

func (ctrl *FileController) TempUrl(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		ctrl.ResErr(ctx, errors.New("path is empty"))
		return
	}

	expiredInMinute := ctx.Query("expired_in_minute")
	expiredAt, err := strconv.Atoi(expiredInMinute)
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}
	if expiredInMinute == "" {
		expiredAt = 30
	}

	now := time.Now()

	fs := filesystem.NewStorage("", "")
	signedURL, err := fs.TemporaryUrl(path, now.Add(time.Duration(expiredAt)*time.Minute))
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"url": signedURL,
		},
	})
}

func (ctrl *FileController) FindAll(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		ctrl.ResErr(ctx, errors.New("path is empty"))
		return
	}

	fs := filesystem.NewStorage("", "")
	files, err := fs.AllFiles(path)
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": files,
	})
}

func (ctrl *FileController) Upload(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	moreInfo := make(map[string]any)
	contentType := file.Header.Get("Content-Type")
	if strings.Contains(contentType, "image") {
		// Open the uploaded file.
		srcFile, _ := file.Open()
		defer srcFile.Close()

		// Decode the uploaded image.
		img, _, _ := image.Decode(srcFile)
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()

		moreInfo["width"] = width
		moreInfo["height"] = height

		if ctx.PostForm("resize_width") != "" || ctx.PostForm("resize_height") != "" {
			name := str.Random(40) + ".jpg"
			file, err = ResizeImage(ctx, *file, ctx.PostForm("resize_width"), ctx.PostForm("resize_height"), name)
			if err != nil {
				ctrl.ResErr(ctx, err)
				return
			}
			defer os.Remove(name)

			if ctx.PostForm("resize_width") != "" {
				moreInfo["width"] = ctx.PostForm("resize_width")
			}

			if ctx.PostForm("resize_height") != "" {
				moreInfo["height"] = ctx.PostForm("resize_height")
			}
		}

		if ctx.PostForm("compress") != "" {
			name := str.Random(40) + ".jpg"
			file, err = CompressImage(file, ctx.PostForm("compress"), name)
			if err != nil {
				ctrl.ResErr(ctx, err)
				return
			}
			defer os.Remove(name)
		}
	}

	c, err := filesystem.NewFileFromRequest(file)
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	fs := filesystem.NewStorage(file.Header.Get("Content-Type"), ctx.PostForm("visibility"))

	folder, _ := config.ReadConfig("filesystems.default_folder")

	fileName := ctx.PostForm("file_name")
	var resultPath string

	path := folder + "/" + ctx.PostForm("folder")

	isRandom := false
	if fileName == "" {
		isRandom = true
		fileName = str.Random(40)
		fullPath, err := filesystem.GetFullPathOfFile(path, c, fileName)

		if err != nil {
			ctrl.ResErr(ctx, err)
			return
		}

		isExist := fs.Exists(fullPath)
		if isExist {
			ctrl.ResErr(ctx, errors.New("file already exist"))
			return
		}
	}

	resultPath, err = fs.PutFileAs(path, c, fileName)

	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	size, err := fs.Size(resultPath)
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	ext, err := c.Extension()
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	defaultDisk, _ := config.ReadConfig("filesystems.default")
	bucketName, _ := config.ReadConfig(fmt.Sprintf("filesystems.disks.%s.bucket", defaultDisk))

	if isRandom {
		extension, err := sfile.Extension(c.File(), true)
		if err != nil {
			ctrl.ResErr(ctx, err)
			return
		}

		fileName = fileName + "." + extension
	}

	responseData := gin.H{
		"data": map[string]any{
			"url":                       fs.Url(resultPath),
			"path":                      resultPath,
			"file_name":                 fileName,
			"size":                      size,
			"content_type":              file.Header.Get("Content-Type"),
			"extension":                 ext,
			"bucket":                    bucketName,
			"client_original_extention": c.GetClientOriginalExtension(),
			"client_original_name":      c.GetClientOriginalName(),
			"disk":                      defaultDisk,
			"more_info":                 moreInfo,
		},
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (ctrl *FileController) Delete(ctx *gin.Context) {
	fs := filesystem.NewStorage("", "")
	err := fs.Delete(ctx.QueryArray("path")...)
	if err != nil {
		ctrl.ResErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"success": true,
		},
	})
}

func (ctrl *FileController) ResErr(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
