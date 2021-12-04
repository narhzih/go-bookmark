package service

import (
	"bytes"
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type FileUploadInformation struct {
	Ctx           *gin.Context
	Logger        zerolog.Logger
	FileInputName string
	Type          string
}

func randSeq(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func saveFileToBuffer(f FileUploadInformation) (interface{}, string, *bytes.Buffer, error) {
	var buff bytes.Buffer
	file, header, err := f.Ctx.Request.FormFile(f.FileInputName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", "", &buff, http.ErrMissingFile
		}
		f.Logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
		return "", "", &buff, err
	}
	defer file.Close()
	//fileName := strings.Split(header.Filename, ".")[0]
	//fileExt := strings.Split(header.Filename, ".")[1]

	//io.Copy(&buff, file)
	//content := buff.String()
	//buff.Reset()
	return file, header.Filename, &buff, nil
}

func saveFileToLocalStorage(f FileUploadInformation) (string, error) {
	file, header, err := f.Ctx.Request.FormFile(f.FileInputName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", http.ErrMissingFile
		}
		f.Logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
		return "", err
	}
	fileName := header.Filename
	fileName = randSeq(20) + "_" + f.Type + "_cover_photo_" + fileName
	out, err := os.Create(fileName)
	if err != nil {
		f.Logger.Err(err).Msg(fmt.Sprintf("Error occurred while trying to save file %+v ", err.Error()))
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		f.Logger.Err(err).Msg(err.Error())
		return "", err
	}
	return fileName, nil
}

func UploadToCloudinary(f FileUploadInformation) (string, error) {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	if err != nil {
		f.Logger.Err(err).Msg(fmt.Sprintf("Cloudinary connection error: %+v", err.Error()))

		return "", err
	}
	//fileName, err := saveFileToLocalStorage(f)
	file, fileNameWithExt, buffer, err := saveFileToBuffer(f)
	if err != nil {
		return "", err
	}
	f.FileInputName = randSeq(20) + "_cover_photo_" + strings.Split(fileNameWithExt, ".")[0]
	resp, err := cld.Upload.Upload(f.Ctx, file, uploader.UploadParams{PublicID: f.FileInputName, FilenameOverride: fileNameWithExt})
	f.Logger.Info().Msg("file name  is " + resp.SecureURL)
	if err != nil {
		f.Logger.Err(err).Msg(err.Error())
	}
	buffer.Reset()

	//err = os.Remove(fileName)
	//if err != nil {
	//	f.Logger.Err(err).Msg(err.Error())
	//}

	return resp.SecureURL, nil
}
