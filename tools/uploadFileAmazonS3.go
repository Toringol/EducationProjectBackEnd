package tools

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadFileToS3 saves a file to aws bucket and returns the url to // the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := "avatars/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	// config settings
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("educationprojectstorage"),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, err
}

// LoadAvatar - check size of file if ok -> create new AWS session
// and upload avatar with unique name to AWS S3 bucket
// return fileName relative to storagePath
func LoadAvatar(ctx echo.Context) (string, error) {
	maxSize := int64(1024000) // allow only 1MB of file size

	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		return "", err
	}

	file, fileHeader, err := ctx.Request().FormFile("profileAvatar")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(
			viper.GetString("accessKeyID"),     // id
			viper.GetString("secretAccessKey"), // secret
			""),                                // token can be left blank for now
	})
	if err != nil {
		return "", err
	}

	fileName, err := UploadFileToS3(s, file, fileHeader)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
