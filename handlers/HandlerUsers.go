package handlers

import (
	"customer-profile-crud/connection"
	"customer-profile-crud/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Login(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var users structs.Users
	json.Unmarshal(payloads, &users)

	var user structs.Users
	connection.DB.Where("name = ?", users.Name).First(&user)

	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(users.Password))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong password"))
			return
		}
	}

	res := structs.Result{Data: user, Message: "Login success"}
	result, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var user structs.Users
	hash, _ := HashPassword(user.Password)

	json.Unmarshal(payloads, &user)

	if user.Name == "" || user.Age == 0 {
		http.Error(w, "Please enter a name and age", http.StatusBadRequest)
	} else {
		connection.DB.Create(&user)
		res := structs.Risk_profile{
			Userid: user.ID,
			Users: structs.Users{
				ID:       user.ID,
				Name:     user.Name,
				Age:      user.Age,
				Password: hash,
			},
			Mm_percent:    0,
			Bond_percent:  0,
			Stock_percent: 0,
			Total_percent: 0,
		}

		if user.Age >= 30 {
			res.Stock_percent = 72.5
			res.Bond_percent = 21.5
			res.Mm_percent = 100 - res.Stock_percent - res.Bond_percent
		} else if user.Age >= 20 {
			res.Stock_percent = 54.5
			res.Bond_percent = 25.5
			res.Mm_percent = 100 - res.Stock_percent - res.Bond_percent
		} else {
			res.Stock_percent = 34.5
			res.Bond_percent = 45.5
			res.Mm_percent = 100 - res.Stock_percent - res.Bond_percent
		}

		res.Total_percent = res.Stock_percent + res.Bond_percent + res.Mm_percent
		connection.DB.Create(&res)
		user := structs.Users{ID: user.ID, Name: user.Name, Age: user.Age, Password: user.Password}
		result, err := json.Marshal(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}
	}
}

func DetailUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user structs.Users
	connection.DB.First(&user, userID)

	res := structs.Risk_profile{
		Userid:        user.ID,
		Users:         structs.Users{ID: user.ID, Name: user.Name, Age: user.Age},
		Mm_percent:    0,
		Bond_percent:  0,
		Stock_percent: 0,
		Total_percent: 0,
	}

	if res.Userid == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	} else {
		connection.DB.Where("userid = ?", user.ID).First(&res)
		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(result)
		}
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	take := r.URL.Query().Get("take")

	users := []structs.Users{}
	connection.DB.
		Limit(take).
		Offset(page).
		Order(`id`).
		Find(&users)

	res := structs.Result{
		Data:    users,
		Message: "Data berhasil didapatkan Page = " + page + " Take = " + take,
	}

	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usersID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var usersUpdates structs.Users
	json.Unmarshal(payloads, &usersUpdates)
	var users structs.Users

	if usersUpdates.Name == "" || usersUpdates.Age == 0 {
		http.Error(w, "Please enter a name and age", http.StatusBadRequest)
	} else {
		connection.DB.First(&users, usersID)
		connection.DB.Model(&users).Updates(&usersUpdates)

		res := structs.Users{ID: users.ID, Name: usersUpdates.Name, Age: usersUpdates.Age}
		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user structs.Users
	connection.DB.Where("id = ?", userID).First(&user)

	if user.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	} else {
		connection.DB.Delete(&user)
		res := structs.Result{Data: user, Message: "Data berhasil dihapus"}
		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}
	}
}
