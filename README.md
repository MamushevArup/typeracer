# Typeracer Clone Game

Welcome to Typeracer Clon Game, a challenging typing game where users can improve their typing skills by racing against each other.
The project still under development with my teammate frontend developer.
To see the pretty design you can go to this repo {soon will added}

## Features

- **Single race mode:** Race against the clock to type as fast as possible.
- **Multiplayer mode:** Compete with up to 5 other racers in real-time.
- **Custom racetrack creation:** Create a racetrack and challenge your friends by sharing a unique link.
- **Real-time speed and place updates:** See your current typing speed and race position in real-time.
- **Automatic link expiration:** Racetrack links expire after one hour to ensure privacy and security.

## Warning!!!

Please pay attention that firslt you should run migration to the local postgres instance.
Go to dir `cd schema`
And if you use goose as migration tool use this command
`goose postgres "host=localhost password=<replace with password of your db> port=5432 user=postgres" up`
This run all migration files
If you use docker
`goose postgres "host=<docker-host> password=<replace with password of your db> port=<incoming port> user=postgres" up`

## Installation

Without Docker

1. Clone the repository: `git clone https://github.com/MamushevArup/typeracer.git`
2. Install dependencies: `go mod tidy`
3. Start the server: `go run cmd/main/app.go`

With Docker
1. Clone the repository: `git clone https://github.com/MamushevArup/typeracer.git`
2. Start the server `docker compose up --build`

## Usage

1. Navigate to the game directory.
2. Open your browser and go to `http://localhost:1001` to access the game.
3. Inspect all routes and without frontend use postman to test it
4. Start typing and enjoy the race!

