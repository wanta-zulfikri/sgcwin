package helper 

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	depedency "github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/errorr"
	"github.com/golang-jwt/jwt"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GetUid(token *jwt.Token) int {
	parse := token.Claims.(jwt.MapClaims)
	id := int (parse["id"]. (float64))

	return id 
}

func GetRole(token *jwt.Token) string {
	parse := token.Claims.(jwt.MapClaims)
	return parse["role"].(string)
}
func GetStatus(token *jwt.Token) string {
	parse := token.Claims.(jwt.MapClaims)
	return parse["verified"].(string)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(passhash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password))
}

func GenerateJWT(id int, role string, is_verified string, dp depedency.Depend) string {
	var informasi = jwt.MapClaims{}
	informasi["id"] = id
	informasi["role"] = role
	informasi["verified"] = is_verified
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, informasi)
	resultToken, err := rawToken.SignedString([]byte(dp.Config.JwtSecret))
	if err != nil {
		log.Println("generate jwt error ", err.Error())
		return ""
	}
	return resultToken
}

func GenerateEndTime(timee string, duration float32) string {
	t, err := time.Parse("2006-01-02 15:04:05", strings.Replace(timee, "T", " ", 1))
	if err != nil {
		log.Printf("error when generate endtime : %v", err)
		return ""
	}
	minute := duration * 60
	return t.Add(time.Minute * time.Duration(int(minute))).Format("2006-01-02 15:04:05")
}
func GenerateExpiretime(timee string, duration int) string {
	t, err := time.Parse("2006-01-02 15:04:05", timee)
	if err != nil {
		return ""
	}
	return t.Add(time.Minute * time.Duration(duration)).Format("2006-01-02 15:04:05")
}
func GenerateInvoice(schoolid int, userid int) string {
	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Intn(9999) + 1000
	return fmt.Sprintf("INV-%d%d%d", userid, schoolid, randomNum)

}

var store = base64Captcha.DefaultMemStore

func GenerateCaptcha() (string, string, error) {
	DriverString := &base64Captcha.DriverString{
		Height:          60,
		Width:           240,
		ShowLineOptions: 0,
		NoiseCount:      0,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		Length:          7,
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	var driver base64Captcha.Driver
	driver = DriverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	return id, b64s, err
}

func VerifyCaptcha(captcha string, value string) bool {
	return store.Verify(captcha, value, true)
}

func CheckNPSN(npsn string, log *logrus.Logger) error {
	client := http.Client{}
	link := fmt.Sprintf("https://referensi.data.kemdikbud.go.id/tabs.php?npsn=%s", npsn)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATING HTTP REQUEST, Error : %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("[ERROR]WHEN GETTING DATA NPSN, Error : %v", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("[ERROR]WHEN READING DATA HTML, Error : %v", err)
	}
	if !doc.Find("div").HasClass("tabby-content") {
		return errorr.NewBad("NPSN not registered")
	}
	return nil
}
func IsValidPhone(number string) bool {
	re := regexp.MustCompile(`^0\d{11,12}$`)
	if re.MatchString(number) {
		return true
	}
	return false
}
