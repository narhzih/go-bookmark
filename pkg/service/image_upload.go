package service

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	fileName, err := saveFileToLocalStorage(f)
	if err != nil {
		//f.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		//	"message": "An error occurred",
		//	"err":     err.Error(),
		//})
		//return
	}
	f.FileInputName = fileName
	resp, err := cld.Upload.Upload(f.Ctx, f.FileInputName, uploader.UploadParams{PublicID: f.FileInputName})
	if err != nil {
		f.Logger.Err(err).Msg(err.Error())
	}

	err = os.Remove(fileName)
	if err != nil {
		f.Logger.Err(err).Msg(err.Error())
	}

	return resp.SecureURL, nil
}
