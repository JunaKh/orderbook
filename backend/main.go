package main

import (
    "log"
    "net/http"
    "net/url"
    "strconv"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
    "encoding/json"
)

// OrderBook represents the order book structure
type OrderBook struct {
    LastUpdateID int64      `json:"lastUpdateId"`
    Bids         [][]string `json:"bids"`
    Asks         [][]string `json:"asks"`
}

// Client represents a WebSocket client
type Client struct {
    conn *websocket.Conn
}

// Server manages clients and broadcasts data
type Server struct {
    clients      map[*Client]bool
    addClient    chan *Client
    removeClient chan *Client
    broadcast    chan float64
    mu           sync.Mutex
}

// newServer creates a new server instance
func newServer() *Server {
    return &Server{
        clients:      make(map[*Client]bool),
        addClient:    make(chan *Client),
        removeClient: make(chan *Client),
        broadcast:    make(chan float64),
    }
}

// handleConnections manages WebSocket connections for clients
func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    client := &Client{conn: conn}

    s.addClient <- client

    defer func() {
        s.removeClient <- client
        client.conn.Close()
    }()

    // Listen for messages from the client (if needed)
    for {
        _, _, err := client.conn.ReadMessage()
        if err != nil {
            log.Println("Client disconnected:", err)
            break
        }
    }
}

// startWorkerPool starts multiple workers to handle broadcasting messages
func (s *Server) startWorkerPool(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        go s.worker()
    }
}

// worker handles broadcasting the messages to all connected clients
func (s *Server) worker() {
    for msg := range s.broadcast {
        s.mu.Lock()
        for client := range s.clients {
            client.conn.WriteJSON(msg)
        }
        s.mu.Unlock()
    }
}

// run starts the server to manage clients and broadcasting
func (s *Server) run() {
    for {
        select {
        case client := <-s.addClient:
            s.mu.Lock()
            s.clients[client] = true
            s.mu.Unlock()
        case client := <-s.removeClient:
            s.mu.Lock()
            delete(s.clients, client)
            s.mu.Unlock()
        case avgPrice := <-s.broadcast:
            s.mu.Lock()
            for client := range s.clients {
                err := client.conn.WriteJSON(avgPrice)
                if err != nil {
                    log.Println("Error sending to client:", err)
                    client.conn.Close()
                    delete(s.clients, client)
                }
            }
            s.mu.Unlock()
        }
    }
}

// connectToBinance connects to the Binance WebSocket and receives order book data
func (s *Server) connectToBinance() {
    u := url.URL{Scheme: "wss", Host: "stream.binance.com", Path: "/ws/btcusdt@depth20"}
    log.Printf("Connecting to %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("Error connecting to Binance WebSocket:", err)
    }
    defer c.Close()

    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            return
        }

        var orderBook OrderBook
        if err := json.Unmarshal(message, &orderBook); err != nil {
            log.Println("Error unmarshalling message:", err)
            continue
        }

        log.Printf("OrderBook: %+v", orderBook)

        avgPrice := calculateAveragePrice(orderBook)
        log.Printf("Average Price: %.2f", avgPrice)

        s.broadcast <- avgPrice
    }
}

// calculateAveragePrice calculates the average price from the bids and asks
func calculateAveragePrice(orderBook OrderBook) float64 {
    sum := 0.0
    count := 0.0

    // Sum the bid prices
    for _, bid := range orderBook.Bids {
        price, err := strconv.ParseFloat(bid[0], 64)
        if err != nil {
            log.Printf("Error parsing bid price: %v", err)
            continue
        }
        sum += price
        count++
    }

    // Sum the ask prices
    for _, ask := range orderBook.Asks {
        price, err := strconv.ParseFloat(ask[0], 64)
        if err != nil {
            log.Printf("Error parsing ask price: %v", err)
            continue
        }
        sum += price
        count++
    }

    // Avoid division by zero
    if count == 0 {
        log.Println("No valid bids or asks available for calculating average")
        return 0.0
    }

    return sum / count
}

func main() {
    server := newServer()

    go server.run()
    go server.connectToBinance()

    r := mux.NewRouter()
    r.HandleFunc("/ws", server.handleConnections)

    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
