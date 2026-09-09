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
4. (Optional) Run the event generator: `go run cmd/generator/main.go`
