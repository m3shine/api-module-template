package aws_s3

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"west.garden/template/common/config"
	"west.garden/template/common/log"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
)

type S3Session struct {
	Sess   *session.Session
	Bucket string
}

var s3Sess *S3Session

func Init(cfg *config.S3Config) error {
	s3Sess = new(S3Session)
	creds := credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	s, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.UsEast2RegionID),
	})
	if err != nil {
		log.Log.Error("aws_s3 Init:", err)
		return err
	}
	s3Sess.Sess = s
	s3Sess.Bucket = cfg.Bucket
	return nil
}

var (
	allowFileExt = map[string]int{".png": 1, ".PNG": 1, ".jpg": 1, ".JPG": 1, ".jpeg": 1, ".JPEG": 1, ".gif": 1, ".GIF": 1, ".svg": 1, ".json": 1, ".zip": 1}
	NotAllowExt  = errors.New("not allow ext")
)

func UploadFile(file multipart.File, fileHeader *multipart.FileHeader, dir string) (string, error) {
	originFilename := filepath.Base(fileHeader.Filename)
	ext := path.Ext(originFilename)
	if _, ok := allowFileExt[ext]; !ok {
		return "", NotAllowExt
	}
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	sh := md5.New()
	sh.Write(buffer)
	imageNameHash := hex.EncodeToString(sh.Sum([]byte("")))
	s := s3Sess
	fileName := fmt.Sprintf("%s/%s%s", dir, imageNameHash, ext)
	_, err := s3.New(s.Sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(buffer),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func GetFileByKey(key string) ([]byte, error) {
	out, err := s3.New(s3Sess.Sess).GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Sess.Bucket), Key: aws.String(key)})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(out.Body)
}

func UploadByteFile(data []byte, dir, ext string) (string, error) {
	sh := md5.New()
	sh.Write(data)
	imageNameHash := hex.EncodeToString(sh.Sum([]byte("")))

	s := s3Sess
	fileName := fmt.Sprintf("file/%s/%s%s", dir, imageNameHash, ext)
	_, err := s3.New(s.Sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(http.DetectContentType(data)),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		return "", err
	}
	return fileName, nil
}
