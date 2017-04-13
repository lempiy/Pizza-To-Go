package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	cldry "github.com/b4t3ou/cloudinary-go"
)

//Keys JSON struct stores secret keys and DB init vars
var Keys keys
var cloudinary *cldry.Cloudinary

type keys struct {
	SessionsKey string `json:"sessions_key"`
	TokenKey    string `json:"token_key"`
	UsernameDB  string `json:"username_db"`
	PasswordDB  string `json:"password_db"`
	NameDB      string `json:"name_db"`
}

// CloudinaryURL is a url of image remote storage for Heroku
var CloudinaryURL string

func init() {
	CloudinaryURL = os.Getenv("CLOUDINARY_URL")

	file, err := ioutil.ReadFile("./keys.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		panic(err)
	}
	log.Println(CloudinaryURL)
	if CloudinaryURL != "" {
		rgxp := regexp.MustCompile(`^[a-z]+(?::\/\/)(?P<Public>\d*):(?P<Secret>\w*)@(?P<Name>\w*$)`)
		slices := rgxp.FindAllStringSubmatch(CloudinaryURL, -1)
		cloudinary = cldry.Create(slices[0][1], slices[0][2], slices[0][3])
	}

	json.Unmarshal(file, &Keys)
}

//Itob - converts int to bool
func Itob(n int) bool {
	if n > 0 {
		return true
	}
	return false
}

// SaveEncodedImage saves base64 encoded image to file or Cloudinary
func SaveEncodedImage(imageCode string) (string, error) {
	// If host is Heroku use Cloudinary
	if CloudinaryURL != "" {
		option := cldry.Option{}

		uploaded, err := cloudinary.Upload(imageCode, option)

		if err != nil {
			log.Println(err)
			return "", err
		}
		return uploaded.SecureUrl, err
	}
	// Else upload on server
	randomFileName := md5.New()
	var fullFileName string
	io.WriteString(randomFileName, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(randomFileName, "pizza")
	token := fmt.Sprintf("%x", randomFileName.Sum(nil))
	pattern := regexp.MustCompile("^data:image/(png|jpeg);base64,")

	imgBase64 := pattern.ReplaceAllString(imageCode, "")

	imageReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgBase64))
	pngImage, _, err := image.Decode(imageReader)
	fullFileName = token + ".png"
	log.Println("Open file " + "./upload/" + fullFileName)
	imgFile, err := os.OpenFile("./upload/"+fullFileName, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println("Open file file error")
		log.Println(err)
		return "", err
	}

	err = png.Encode(imgFile, pngImage)

	if err != nil {
		log.Println(err)
		return "", err
	}
	defer imgFile.Close()
	return fullFileName, err
}

//RemoveUnusedImg deletes useless file from drive
func RemoveUnusedImg(imageURL string) error {
	var err error
	if !strings.HasPrefix(imageURL, "https://") {
		err = os.Remove("./" + imageURL)
	}
	return err
}

//EncryptPassword makes password encryption with SHA1 algorithm
func EncryptPassword(password string) string {
	h := sha1.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
