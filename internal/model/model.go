package model

type Server struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	IconURL  string    `json:"logo_link"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	ServerId      int            `json:"-"`
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	ChannelTitles []ChannelTitle `json:"chats"`
}

type ChannelTitle struct {
	ServerId  int    `json:"-"`
	ChannelId int    `json:"-"`
	Id        int    `json:"id"`
	Name      string `json:"title"`
}

type ServerMember struct {
	ServerId int `json:"server_id"`
	UserId   int `json:"user_id"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	AvatarURL    string `json:"avatar_url"`
}

type Message struct {
	Id        int    `json:"-"`
	ServerId  int    `json:"server_id"`
	ChannelId int    `json:"channel_id"`
	ChatId    int    `json:"chat_id"`
	AuthorId  int    `json:"author_id"`
	Content   string `json:"content"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetMessagesRequest struct {
	ServerId  int `json:"server_id"`
	ChannelId int `json:"channel_id"`
	ChatId    int `json:"chat_id"`
}

type ClientMessage struct {
	Type      string `json:"type"`
	ServerID  int64  `json:"server_id"`
	ChannelID int64  `json:"channel_id"`
	ChatID    int64  `json:"chat_id"`
}

type MessageEvent struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}
