package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type PasscodeBody struct {
	Passcode string `json:"passcode"`
}

// route " api/room/id" GET handler
func (a *API) RoomHandler(w http.ResponseWriter, r *http.Request) {
	// get room id

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	ctx := context.Background()

	// check room is privet and is privet required passcode
	var passcodeDB string
	var passcodeBody PasscodeBody
	err := a.pool.QueryRow(ctx, "SELECT passcode FROM room WHERE id = $1", roomId).Scan(&passcodeDB)
	if err != nil {
		WriteJSON(w, 500, map[string]any{
			"error": "room is not valid!",
		})
		return

	}

	// cheack is room privet
	if passcodeDB != "" {
		// get body data froms pass code
		err := json.NewDecoder(r.Body).Decode(&passcodeBody)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "passcode is not valid!",
			})
			return
		}

		// cheack passcord equils? db password

		if passcodeBody.Passcode != passcodeDB {
			WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "passcode is not valid!",
			})

			// debug
			// fmt.Printf("bodypassword : %v \ndbpassword : %v", passcodeBody.Passcode, passcodeDB)
			return

		}

	}

	// fetch data from db
	var roomData RoomData

	err = a.pool.QueryRow(ctx, "SELECT name,streamUrl,subtitleUrl FROM room WHERE id = $1 ", roomId).Scan(&roomData.RoomName, &roomData.StreamUrl, &roomData.SubtitileUrl)

	if err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]any{
			"error": "room is not valid!",
		})
		return

	}
	WriteJSON(w, 200, map[string]any{
		"id":           roomId,
		"name":         roomData.RoomName,
		"streamUrl":    roomData.StreamUrl,
		"subtitlesUrl": roomData.SubtitileUrl,
	})

}
