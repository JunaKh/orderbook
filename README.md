# Order Book Real-Time WebSocket Application

This project is a full-stack real-time order book service that connects to Binance's WebSocket API to fetch and display the current average price of BTC/USDT. 
It uses a Go backend to handle WebSocket connections with Binance and clients, and a SvelteKit-based frontend to display the data in an interactive chart.

## Table of Contents

- [Project Overview](#project-overview)
- [How It Works](#how-it-works)
- [Technologies Used](#technologies-used)
- [Features](#features)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
    - [Backend (Go)](#backend-go)
    - [Frontend (SvelteKit)](#frontend-sveltekit)
- [Improvements & Future Features](#improvements--future-features)
- [Customization](#customization)
- [Contributing](#contributing)


## Project Overview

This service implements real-time WebSocket connections both for fetching data from Binance's WebSocket API and broadcasting it to multiple clients connected via WebSockets. The frontend visualizes this data in a live chart using **Chart.js**.

The application is designed to be responsive and scalable, making it capable of handling multiple client connections efficiently.

## How It Works

1. **Backend**:
    - Connects to the Binance WebSocket stream and receives order book updates for BTC/USDT.
    - Calculates the average price based on bid and ask prices.
    - Broadcasts the calculated average price to all connected WebSocket clients.

2. **Frontend**:
    - Listens for WebSocket messages from the backend.
    - Processes the received average price and displays it in a real-time chart using **Chart.js**.

## Technologies Used

### Backend:
- **Go** for backend logic and concurrency.
- **Gorilla WebSocket** for WebSocket communication with both Binance and clients.

### Frontend:
- **SvelteKit** for the frontend framework, providing fast reactivity and component-based architecture.
- **Chart.js** for rendering dynamic, real-time charts.

## Features

- Real-time WebSocket connection with Binance to fetch BTC/USDT order book data.
- Calculates and displays the average price from the order book.
- Broadcasts data to multiple connected WebSocket clients.
- Interactive chart powered by Chart.js for real-time data visualization.
- Fully responsive UI with dynamic updates and error handling.

## Getting Started

### Prerequisites

- Go 1.23 or higher installed on your machine.
- Node.js (v16.x or higher) and npm/yarn for frontend development.
- Binance WebSocket API access.

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/your-repository/orderbook.git
   cd orderbook
   ```

2. **Backend Setup**:

    - Navigate to the `backend` folder and run the following:

      ```bash
      cd backend
      go mod tidy
      ```

3. **Frontend Setup**:

    - Navigate to the `frontend` folder and install dependencies:

      ```bash
      cd frontend
      npm install
      ```

### Running the Application

1. **Start the Backend**:

    - Run the Go backend to start the WebSocket server:

      ```bash
      cd backend
      go run main.go
      ```

    - The server will be available on port `8080`.

2. **Start the Frontend**:

    - Run the frontend development server:

      ```bash
      cd frontend
      npm run dev
      ```

    - The frontend will be accessible on `http://localhost:3000`.

## Project Structure

```plaintext
orderbook/
│
├── backend/                      # Go WebSocket server
│   ├── go.mod
│   ├── go.sum
│   ├── main.go                   # Main Go file handling Binance connection and WebSocket clients
│
├── frontend/
│   ├── src/
│   │   ├── lib/                  # Helper functions for frontend
│   │   │   └── index.ts          # Utilities
│   │   ├── routes/               # Svelte routes
│   │   │   └── +page.svelte      # Main page, WebSocket connection, and Chart.js logic
│   │   ├── app.d.ts              # TypeScript definitions
│   │   └── app.html              # HTML template
│   ├── public/
│   ├── package.json
│   ├── svelte.config.js
│   ├── tsconfig.json
│   ├── vite.config.ts
│
└── README.md
```

### Backend (Go)

The backend is responsible for:

- Establishing a WebSocket connection with Binance to fetch BTC/USDT order book data.
- Calculating the average price from the order book bids and asks.
- Broadcasting the average price to all connected WebSocket clients.
- Managing client connections, handling messages, and disconnects.

### Frontend (SvelteKit)

The frontend:

- Connects to the backend WebSocket to receive real-time average price updates.
- Displays the real-time data using a Chart.js line chart.
- Implements error handling for WebSocket connection issues.
- Provides a responsive layout for different screen sizes.

Key files:
- **`+page.svelte`**: Contains the main logic for establishing WebSocket connections and rendering the Chart.js chart.
- **`index.ts`**: Utility functions used across the frontend.

## Improvements & Future Features

1. **Error Handling**:  
   Enhance error handling in the WebSocket logic to reconnect automatically if the connection is dropped.

2. **Deployment**:  
   Add a production-ready setup, such as using Docker containers for both the frontend and backend.

3. **Testing**:  
   Add unit tests for backend logic (e.g., calculating average prices) to ensure correctness and improve maintainability.

## Customization

### Changing WebSocket Endpoints

If you need to point to a different WebSocket endpoint, update the URL in the frontend:

```javascript
const ws = new WebSocket("ws://localhost:8080/ws");
```

And in the Go backend's `main.go` file, you can update the Binance WebSocket URL:

```go
u := url.URL{Scheme: "wss", Host: "stream.binance.com", Path: "/ws/btcusdt@depth20"}
```

### Chart.js Customization

Chart.js options can be customized in the `+page.svelte` file. For example, to change the chart colors, update the dataset configurations:

```javascript
borderColor: "rgba(75, 192, 192, 1)",
backgroundColor: "rgba(75, 192, 192, 0.2)",
```

## Contributing

Feel free to submit a pull request or open an issue if you'd like to contribute or find any bugs.

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.
