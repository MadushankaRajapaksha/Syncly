package api

import (
	"encoding/json"
	"net/http"
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
