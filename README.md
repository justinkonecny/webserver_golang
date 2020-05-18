# Calendays Web Server
Web Server written in Go.

## Session Authentication
 - `POST /login`: creates a session for the user with the given UUID
   - (Required) Header: `FirebaseUUID`


## Routes
 - `GET /networks`: returns a list of networks the user belongs to

 - `GET /events`: returns a list of events for all networks the user belongs to