package handlers

import (
	"customer-profile-crud/connection"
	"customer-profile-crud/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var profile structs.Users
	json.Unmarshal(payloads, &profile)

	json.Unmarshal(payloads, &profile)

	if profile.Name == "" || profile.Age == 0 {
		http.Error(w, "Please enter a name and age", http.StatusBadRequest)
	} else {
		connection.DB.Create(&profile)
		res := structs.Risk_profile{
			Userid:        profile.ID,
			Users:         structs.Users{},
			Mm_percent:    0,
			Bond_percent:  0,
			Stock_percent: 0,
			Total_percent: 0,
		}

		if profile.Age >= 30 {
			res.Stock_percent = 72.5
			res.Bond_percent = 21.5
			res.Mm_percent = 100 - res.Stock_percent - res.Bond_percent
		} else if profile.Age >= 20 {
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
		profile := structs.Users{ID: profile.ID, Name: profile.Name, Age: profile.Age}
		result, err := json.Marshal(profile)

		if err != nil {
			panic(err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(result)
		}
	}
}

func DetailUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["id"]

	var profile structs.Users
	connection.DB.First(&profile, profileID)

	res := structs.Risk_profile{
		Userid:        profile.ID,
		Users:         structs.Users{ID: profile.ID, Name: profile.Name, Age: profile.Age},
		Mm_percent:    0,
		Bond_percent:  0,
		Stock_percent: 0,
		Total_percent: 0,
	}

	if res.Userid == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	} else {
		connection.DB.Where("userid = ?", profile.ID).First(&res)
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
		ListUser: users,
		Message:  "Data berhasil didapatkan Page = " + page + " Take = " + take,
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

	users := []structs.Users{}
	connection.DB.First(&users, userID)
	connection.DB.Delete(&users)

	res := structs.Result{
		Message: "Success delete user",
	}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
