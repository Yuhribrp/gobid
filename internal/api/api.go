package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/Yuhribrp/gobid/internal/services"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductsService
	BidsService    services.BidsService
	Sessions       *scs.SessionManager
	WsUpgrader     websocket.Upgrader
	AuctionLobby   *services.AuctionLobby
}
