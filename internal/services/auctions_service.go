package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota + 1
)

type Message struct {
	Message string
	Kind    MessageKind
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
