# Syncly

syncly is a lightweight watch-party web app i build so friends can watch videos in sync. you create a room, drop in a direct video link,share the room URL, and everyones player stays locked together in real time.

![alt text](image-1.png)

# What it does

- **Real-time playback sync** : Play,pause, seek , or swap the video URL , every action broadcasts over WebSockets so the whole room stays on the exact same second
-  **SRT and VTT subttile Suport** : Add `.srt` or `.vtt` link. the frontend converts `.srt` format to to webvtt in memory and renders a custom subtitle overlay that works even in fullscreen.
-  **Live Chat** : See who is currently in the room with their profile picture and chat side by side while watching
-  **Public or Privet rooms** : can u keep room with or witout locked
# How to use it

![alt text](image.png)
1. [Create a room](https://syncly-95r5.onrender.com/create-room) -  Complete your profile data and roomdata (add stream url as a direct video link and subtitle url also should be direct .srt or .vtt url) and make it privet or public with passcode
2. **Share the Url** - After create room go to the room and copy url and send it to ur friends.
3. **Play and Chat** - use controls its sync all friends browsers, and keep chat
# Tech Stack
- **Backend** : Go (`net/http`, `gorilla/mux`, `gorilla/websocket`)
- **Database** : PostgreSQL (`pgx/v5`)
- **Frontend** : HTML, CSS, js (no frontend frameworks )
  
# Running it locally [DEV]
## 1. Prerequisites
- Go 1.21+
- PostgressSQL

## 2. Clone and Config
```
gitclone
cd syncly
```
#### Config `.env`
create .env file and add those

```
DATABASE_URL=postgres://user:password@localhost:5432/syncly?sslmode=disable
PORT=3000

```

## 3. Run

```
go mod tidy
go run .
```
The server starts on http://localhost:3000.
