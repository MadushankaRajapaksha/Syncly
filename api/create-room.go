package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserData struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type RoomData struct {
	RoomName     string `json:"roomName"`
	StreamUrl    string `json:"streamUrl"`
	SubtitileUrl string `json:"subtitleUrl"`
	IsRoomPrivet bool   `json:"isRoomPrivet"`
	Passcode     string `json:"passcode"`
}

type Room struct {
	UserData `json:"userData"`
	RoomData `json:"roomData"`
}

// type PgxPool pgxpool.Pool

type API struct {
	pool *pgxpool.Pool
}

func generateRoomId() (string, error) {
	// generate hex value with lenth 6 for room id
	b := make([]byte, 3)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	// a reuseble response wirter this can use write response with json
	w.Header().Set("Content-Type", "application/json") // content type in header
	w.WriteHeader(status)                              //status code
	json.NewEncoder(w).Encode(v)                       //content in json
}

func (a *API) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// cheack request method and if is not a POST return method not allowed
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method Not Allowed!")
	}

	var roomData Room

	err := json.NewDecoder(r.Body).Decode(&roomData)

	if err != nil {
		fmt.Fprint(w, "error")
		return
	}

	w.WriteHeader(http.StatusAccepted)

	ctx := context.Background()

	// -- add  user to database and get user id
	var userID int
	err = a.pool.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		roomData.UserName, roomData.Email,
	).Scan(&userID)
	// cheak is there have any errs
	if err != nil {
		fmt.Fprintf(w, "users has create faill : %v", err)
		log.Fatal("Users Insertion fail!")

	}

	// -- add room with admin as creatd user
	var roomId string
	// generet unique room id
	// somtimes functions generate room id that is already exsists id in db , itaratie untill unique room id
	for {
		roomIdDraft, _ := generateRoomId()

		err = a.pool.QueryRow(
			ctx,
			"INSERT INTO room (id,name,streamUrl,subtitleUrl,passcode,adminID) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
			roomIdDraft,
			roomData.RoomName,
			roomData.StreamUrl,
			roomData.SubtitileUrl,
			roomData.Passcode,
			userID,
		).Scan(&roomId)

		// cheack is have erros and its returning room id

		if err == nil {
			break // break when roomid is unique

		}

	}

	// make respose with room id an user id

	WriteJSON(w, http.StatusAccepted, map[string]any{
		"roomId": roomId,
		"userId": userID,
	})

}
