package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"jwt/model"
	"jwt/repository"

	jwt "github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	gubrak "github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

var APPLICATION_NAME = "My JWT App"
var LOGIN_EXPIRATION_TIME = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("the secret dion")

type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"Username"`
	Email    string `json:"Email"`
}

func main() {
	mux := new(CustomMux)
	mux.RegisterMiddleware(MiddlewareJWTAuthorization)

	mux.HandleFunc("/login", handlerLogin)
	mux.HandleFunc("/index", handlerIndex)
	mux.HandleFunc("/showMhs", handlerShowMhs)
	mux.HandleFunc("/createMhs", handlerCreateMhs)
	mux.HandleFunc("/updateMhs", handlerUpdateMhs)
	mux.HandleFunc("/deleteMhs", handlerDeleteMhs)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":8080"

	fmt.Println("Starting server at", server.Addr)
	server.ListenAndServe()
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	message := fmt.Sprintf("hello %s with email : %s", userInfo["Username"], userInfo["Email"])
	w.Write([]byte(message))
}

func handlerShowMhs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mhsRepo := repository.NewMhsRepository(db)
	mahasiswas, err := mhsRepo.FetchMhs()
	if err != nil {
		panic(err)
	}

	result, err := json.MarshalIndent(mahasiswas, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func handlerUpdateMhs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var mhs model.Mahasiswa
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprint("read body error: ", err.Error()), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &mhs)
	if err != nil {
		http.Error(w, fmt.Sprint("JSON encode error: ", err.Error()), http.StatusInternalServerError)
		return
	}
	err = repository.NewMhsRepository(db).UpdateMhs(mhs)
	if err != nil {
		http.Error(w, "Input Gagal!", http.StatusBadRequest)
	}
	w.Write([]byte("Data berhasil diubah!"))

}

func handlerDeleteMhs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = repository.NewMhsRepository(db).DeleteMhsByNIM(id)
	if err != nil {
		http.Error(w, "Delete Gagal!", http.StatusBadRequest)
	}
	w.Write([]byte("Data berhasil dihapus!"))

}

func handlerCreateMhs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var mhs model.Mahasiswa
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprint("read body error: ", err.Error()), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &mhs)
	if err != nil {
		http.Error(w, fmt.Sprint("JSON encode error: ", err.Error()), http.StatusInternalServerError)
		return
	}
	resultID, err := repository.NewMhsRepository(db).CreateMhs(mhs)
	if err != nil {
		http.Error(w, "Input Gagal!", http.StatusBadRequest)
	}
	concatenated := fmt.Sprintf("ID : %d berhasil ditambahkan!", resultID)
	w.Write([]byte(concatenated))

}
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Invalid username or password basic", http.StatusBadRequest)
		return
	}

	ok, userInfo := authenticateUser(username, password)
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	fmt.Println(userInfo)

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_TIME).Unix(),
		},
		Username: userInfo["Username"].(string),
		Email:    userInfo["Email"].(string),
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, _ := json.Marshal(M{"token": signedToken})

	w.Write([]byte(tokenString))
}

func authenticateUser(username, password string) (bool, M) {
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.FetchUser()
	if err != nil {
		panic(err)
	}

	data := make([]M, 0)
	inrec, _ := json.Marshal(user)
	json.Unmarshal(inrec, &data)

	res := gubrak.From(data).Find(func(each M) bool {
		return each["Username"] == username && each["Password"] == password
	}).Result()

	// panic(res != nil)
	if res != nil {
		resM := res.(M)
		delete(resM, "Password")
		return true, resM
	}

	return false, nil
}

func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing Method Invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("signing Method Invalid")
			}

			return JWT_SIGNATURE_KEY, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(context.Background(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
