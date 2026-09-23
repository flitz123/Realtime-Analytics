# Real-Time Analytics & Event Processing Platform

A real-time analytics platform built with Go that processes and visualizes events as they happen.

## Features

- **Event Ingestion API**: REST endpoint to receive events
- **WebSocket Support**: Real-time bidirectional communication
- **In-Memory Processing**: Fast event processing and analytics
- **Real-Time Dashboard**: Live metrics and event visualization
- **Event Simulation**: Automatic event generation for testing

## Architecture

- **Event Store**: In-memory storage for recent events
- **Analytics Processor**: Calculates real-time metrics every 5 seconds
- **WebSocket Hub**: Manages client connections and broadcasts
- **REST API**: HTTP endpoints for event ingestion and queries

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Web browser with JavaScript enabled


## Running the Project

1. Run `go mod tidy` to download dependencies
2. Start the server: `go run cmd/server/main.go`
3. Open `http://localhost:8080/dashboard.html` in your browser

## Deployment

This project is deployable as a traditional Go service on a host that supports
long-running processes, such as Render, Railway, Fly.io, or a VM. Set the
`PORT` environment variable in production and bind the HTTP server to that
port.

### Render Backend

Create a Render **Web Service** connected to this repository with:

- **Runtime:** Go
- **Build command:** `go build -o app ./cmd/server`
- **Start command:** `./app`
- **Health check path:** `/api/stats`

Render supplies `PORT` automatically. If the frontend is hosted on Vercel,
add a Render environment variable named `FRONTEND_URL` with the Vercel URL,
for example `https://your-dashboard.vercel.app`.

Before deploying the frontend, set `window.BACKEND_URL` in `web/config.js` to
the Render service URL, for example `https://your-service.onrender.com`.

It is not deployable as the complete application on Vercel as currently
implemented. The dashboard can be hosted as static files on Vercel, but the
Go server uses in-memory state, a background metrics ticker, and persistent
WebSocket connections. Vercel Functions are request-based and do not provide
that long-running WebSocket process. For a Vercel frontend, deploy the Go
backend separately and configure the dashboard to use its API and WebSocket
URL.
