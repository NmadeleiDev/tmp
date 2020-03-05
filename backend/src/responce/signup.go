package responce

import (
	"Matcha/postgres"
	"Matcha/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

func SignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/signup/")
}

func Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("passwd1")
	passwordConfirm := r.FormValue("passwd2")
	if strings.Compare(password, passwordConfirm) != 0 {
		w.Write([]byte("User: " + username + "  Password: " + password))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	uniqueKey := utils.GetMD5(time.Now().String() + username + password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(" InsertUserData error. Err: " + err.Error())
	}
	queryResult := postgres.InsertUserData(username, email, string(hashedPassword), uniqueKey)
	if queryResult == false {
		fmt.Println("Account " + username + " already exists")
		w.Write([]byte("Already exists"))
		w.WriteHeader(http.StatusAccepted)
		return
	} else {
		fmt.Println("Account " + username + " created")
		sendVerifEmail(email, uniqueKey)
		w.Write([]byte("Okay"))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func sendVerifEmail(email, vkey string) {
	auth := smtp.PlainAuth("", "saveencrypteddata@gmail.com", "LYwu4>wT", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Follow this link to verify your account:\r\n" +
		"\r\n" + "http://localhost:8080/verify?key=" + vkey + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "saveencrypteddata@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func VerifyEmail(w http.ResponseWriter, r *http.Request)  {
	vkey, ok := r.URL.Query()["key"]
	if !ok || len(vkey[0]) < 10 {
		fmt.Println("Some shit with vkey")
		return
	}
	fmt.Println("Got vkey: " + vkey[0])
	result := postgres.VerifyAccount(vkey[0])
	if result {
		fmt.Println("Validated")
		http.Redirect(w, r, "/", 200)
	} else {
		fmt.Println("Vkey is invalid")
		http.Redirect(w, r, "/signup", 200)
	}
}