package main

import "proweb-backend/api"

const PORT = ":8000"
const URL = "http://localhost" + PORT

func main() {
	api.StartServer(URL, PORT)
}
