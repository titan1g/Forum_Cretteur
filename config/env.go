package config

import uuid "github.com/gofrs/uuid"

//const for hashing
const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

// Initiliazing the port
const LocalhostPort = ":8080"

//basic struct
type Account struct {
	Name            string
	Password        string
	Email           string
	Uuid            uuid.UUID
	Profile_Picture string
	Role            string
}

type AllAccount struct {
	Data    []Account
	Account LoginYes
}

type LoginYes struct {
	Connected bool
	Account   Account
}

type TName struct {
	Id           int
	Title        string
	Desc         string
	Category     string
	CreationDate string
	Creator      string
	Pic          bool
	Like         int
	Liker        string
	Disliker     string
	Liked        int
}

type TContent struct {
	Id      int
	Uuid    string
	Name    string
	Text    string
	Written string
	Picture string
}

type UserActions struct {
	Commentaires []TContent
	Account      Account
	Login        LoginYes
}

//single topics
type Topics struct {
	Name     TName
	Content  []TContent
	Accounts []Account
	Login    LoginYes
}

//all topics
type AllTopics struct {
	Name  []TName
	Login LoginYes
}

type Erreur struct {
	Connected bool
	Miss      bool
}

type RealEmailResponse struct {
    Status string `json:"status"`
}