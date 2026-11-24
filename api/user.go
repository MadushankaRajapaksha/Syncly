package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

//	const user = {
//	          userId : data.userId,
//	          username,
//	          email,
//	          profilePicture
//	        };
type UserInDb struct {
	UserId         any    `json:"userId"`
	UserName       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

type UserBodyCreate struct {
	UserName       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

func (a *API) ValidateUser(w http.ResponseWriter, r *http.Request) {
	// get user data from body in reques
	var userBody UserInDb
	err := json.NewDecoder(r.Body).Decode(&userBody)
	// cheack body data is have and is not return err
	if err != nil {
		WriteJSON(w, 500, map[string]string{
			"error": "Erro, required data are missing!",
		})
		return

	}

	defer r.Body.Close()

	// using database and compaire data for frontend data and database data are equel

	var userDB UserInDb
	// ctx := context.Background()
	err = a.pool.QueryRow(r.Context(), "SELECT id,name,email FROM users WHERE id = $1", userBody.UserId).Scan(&userDB.UserId, &userDB.UserName, &userDB.Email)

	if err != nil {
		WriteJSON(w, 500, map[string]any{
			"err": err,
		})
		return

	}

	// compare with selction
	if userBody.Email == userDB.Email {
		WriteJSON(w, 200, map[string]any{
			"id": userDB.UserId,
		})
		return
	}

	// is not valid user return emty string as id
	WriteJSON(w, 404, map[string]any{
		"id": 0,
	})

}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	// get body data
	var user UserBodyCreate
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, 500, map[string]string{
			"error": "Erro, required data are missing!",
		})
		return

	}

	defer r.Body.Close()
	// cheack  body data
	// cheack user data is emty
	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.UserName) == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "missing required fields",
		})
		return
	}

	ctx := context.Background()

	var userDataDb UserInDb

	// cheack db user data is alredy have
	err = a.pool.QueryRow(ctx, "SELECT id, name, email, profilepicture FROM users WHERE email = $1", user.Email).Scan(&userDataDb.UserId, &userDataDb.UserName, &userDataDb.Email, &userDataDb.ProfilePicture)

	if errors.Is(err, pgx.ErrNoRows) { //user is not found
		// create user and return data
		var userReturn UserInDb

		err = a.pool.QueryRow(ctx,
			"INSERT INTO users (name, email, profilepicture) VALUES ($1, $2, $3) RETURNING id, name, email, profilepicture",
			user.UserName, user.Email, user.ProfilePicture,
		).Scan(
			&userReturn.UserId,
			&userReturn.UserName,
			&userReturn.Email,
			&userReturn.ProfilePicture,
		)
		// cheak is there have any errs
		if err != nil {
			fmt.Fprintf(w, "users has create faill : %v", err)
			log.Fatal("Users Insertion fail!")
			WriteJSON(w, 500, map[string]any{
				"erro": "Internal Server Error",
			})
			return

		}

		WriteJSON(w, 200, userReturn)
		return

	}
	if err != nil {
		log.Printf("Database query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Internal Server Error",
		})
		return
	}

	WriteJSON(w, 200, userDataDb)

}
