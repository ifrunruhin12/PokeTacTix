# PokeTacTix Server

This is the backend server component of PokeTacTix, built with Go Fiber. It serves the web API and renders templates from the client directory.

## Project Structure

The server is part of a larger project structure:
```
PokeTacTix/
├── client/                 # Frontend code (HTML, CSS, templates)
      ├── public/           # Conatins CSS and JS files
      ├── views/            # Conatins HTML files
├── server/                 # Backend API (Fiber, templates)
│   ├── main.go            # Server entry point
│   └── README.md          # This file
├── pokemon/                # Shared Pokemon logic
├── game/                   # Game logic
└── go.mod                  # Go dependency
```

## Prerequisites

- Go 1.24.3 or later
- The client directory must be at the same level as the server directory

## Getting Started

1. Make sure you're in the server directory:
   ```bash
   cd server
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

4. The server will start at:
   ```
   http://localhost:3000
   ```

## API Endpoints

- `GET /` - Home page
- `GET /pokemon?name=<pokemon_name>` - Get Pokemon details
  - Returns 400 if name is missing
  - Returns 404 if Pokemon not found
  - Returns Pokemon card data and renders the template on success

## Static Files

The server serves static files from the `../client/public` directory and templates from `../client/views`. Make sure these directories exist and contain the necessary files.

## Development

- The server uses Go Fiber for routing and API handling
- Templates are rendered using the `html/template` engine
- CORS is enabled for all origins (can be configured in main.go)
- The server reuses the Pokemon types and functions from the `pokemon` package

## Contributing

Feel free to submit issues and enhancement requests! 
