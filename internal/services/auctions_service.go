package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota

	SuccessfullyPlacedBid
	NewBidPlaced
	AuctionFinished
	FailedToPlaceBid
	InvalidJSON
)

type Message struct {
	Message string
	Kind    MessageKind
	UserID  uuid.UUID
	Amount  float64
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	ID         uuid.UUID
	Context    context.Context
	Clients    map[uuid.UUID]*Client
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client

	BidsService *BidsService
}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user Connected", "Client", c)
	r.Clients[c.UserID] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("New user Disconnected", "Client", c)
	delete(r.Clients, c.UserID)
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New message received", "RoomID", r.ID, "message", m.Message, "user_id", m.UserID)
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.ID, m.UserID, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidsTooLow) {
				if client, ok := r.Clients[m.UserID]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidsTooLow.Error(), UserID: m.UserID}
				}
				return
			}
		}

		if client, ok := r.Clients[m.UserID]; ok {
			client.Send <- Message{Kind: SuccessfullyPlacedBid, Message: "Your bid was Successfully placed.", UserID: m.UserID}
		}

		for id, client := range r.Clients {
			newBidMessage := Message{
				Kind:    NewBidPlaced,
				Message: "A new bid was placed",
				Amount:  bid.BidAmount,
				UserID:  m.UserID,
			}
			if id == m.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	case InvalidJSON:
		client, ok := r.Clients[m.UserID]
		if !ok {
			slog.Info("Client not found in hashmap", "user_id", m.UserID)
			return
		}
		client.Send <- m
	}
}

func (r *AuctionRoom) Run() {
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)

		case <-r.Context.Done():
			slog.Info("Auction has ended.", "auctionID", r.ID)
			for _, client := range r.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "auction has been finished"}
			}
			return
		}
	}
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, BidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		ID:          id,
		Context:     ctx,
		Clients:     make(map[uuid.UUID]*Client),
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		BidsService: &BidsService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserID uuid.UUID
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserID: userID,
	}
}
