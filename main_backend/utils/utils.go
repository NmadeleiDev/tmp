package utils

import (
	"backend/db/structuredDataStorage"
	"backend/types"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	types2 "go/types"
	"net/http"
	"reflect"
)

const (
	oneDayInSeconds = 86400
)

func RefreshRequestSessionKeyCookie(w http.ResponseWriter, user types.LoginData) bool {
	sessionKey, err := structuredDataStorage.IssueUserSessionKey(user)

	if err != nil {
		SendFailResponse(w, "incorrect user data")
		return false
	}

	SetCookie(w, "session_id", sessionKey)
	return true
}

func SetCookie(w http.ResponseWriter, cookieName, value string) {
	c := http.Cookie{
		Name:  cookieName,
		Value: value,
		Path:  "/",
		//SameSite: http.SameSiteNoneMode,
		MaxAge: oneDayInSeconds * 1}
	http.SetCookie(w, &c)
}

func GetCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println("Failed getting cookie", err)
		return ""
	} else {
		log.Println("Got cookie: ", cookie)
	}
	return cookie.Value
}

func SendFailResponse(w http.ResponseWriter, text string) {
	var packet []byte
	var err error

	response := &types.ResponseJson{Status: false, Data: text}
	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func SendSuccessResponse(w http.ResponseWriter) {
	var packet []byte
	var err error

	response := &types.ResponseJson{Status: true, Data: nil}
	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func SendDataResponse(w http.ResponseWriter, data interface{}) {
	var packet []byte
	var err error

	response := &types.ResponseJson{Status: true, Data: data}
	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func IdentifyUserBySession(r *http.Request) (string, bool) {
	session := GetCookieValue(r, "session_id")
	data, err := structuredDataStorage.Manager.GetUserIdBySession(session)
	if err != nil {
		return "", false
	}
	return data.Id, true
}

func ReflectInterface(from interface{}, to interface{}) {
	valFrom := reflect.ValueOf(from)
	typeTo := reflect.TypeOf(to)
	valTo := reflect.ValueOf(&to).Elem()
	for i := 0; i < typeTo.NumField(); i++ {
		name, ok := typeTo.Field(i).Tag.Lookup("ws")
		if !ok {
			continue
		}
		val := valFrom.FieldByName(name)
		if valTo.Field(i).CanSet() {
			valTo.Field(i).Set(val)
		} else {
			log.Infof("can't set field %v with val %v", typeTo.Field(i).Name, val)
		}
	}
	log.Infof("Reflected: %v", to)
}

func DoesArrayContain(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
