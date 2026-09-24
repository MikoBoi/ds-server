package server

import (
	"context"

	"discord/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetServers(ctx context.Context, userID int64) ([]model.Server, error) {
	// rows, err := r.db.Query(ctx, `SELECT id, name, icon_url FROM servers ORDER BY id`)
	rows, err := r.db.Query(ctx, `SELECT s.id, s.name, s.icon_url FROM servers s, server_members sm where s.id = sm.server_id and user_id = $1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []model.Server

	for rows.Next() {
		var server model.Server

		err := rows.Scan(
			&server.Id,
			&server.Name,
			&server.IconURL,
		)
		if err != nil {
			return nil, err
		}

		channelRows, err := r.db.Query(ctx, `SELECT id, name FROM channels WHERE server_id = $1 ORDER BY id`, server.Id)
		if err != nil {
			return nil, err
		}
		defer channelRows.Close()

		var channels []model.Channel

		for channelRows.Next() {
			var channel model.Channel

			err := channelRows.Scan(
				&channel.Id,
				&channel.Name,
			)
			if err != nil {
				return nil, err
			}

			channelTitleRows, err := r.db.Query(ctx, `SELECT id, name FROM channel_titles WHERE server_id = $1 AND channel_id = $2 ORDER BY id`, server.Id, channel.Id)
			if err != nil {
				return nil, err
			}
			defer channelTitleRows.Close()

			var channelTitles []model.ChannelTitle

			for channelTitleRows.Next() {
				var channelTitle model.ChannelTitle

				err := channelTitleRows.Scan(
					&channelTitle.Id,
					&channelTitle.Name,
				)
				if err != nil {
					return nil, err
				}
				channelTitles = append(channelTitles, channelTitle)
			}

			channel.ChannelTitles = channelTitles
			channels = append(channels, channel)
		}

		server.Channels = channels
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *Repository) GetMessages(ctx context.Context, server_id, channel_id, chat_id int64) ([]model.Message, error) {
	rows, err := r.db.Query(ctx, `SELECT author_id, content from messages where server_id = $1 and channel_id = $2 and chat_id = $3`, server_id, channel_id, chat_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var message model.Message

		err := rows.Scan(
			&message.AuthorId,
			&message.Content,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, username, email, avatar_url FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `SELECT id, username, email, avatar_url FROM users where id = $1`, userID).Scan(&user.Id, &user.Username, &user.Email, &user.AvatarURL)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Authorization(ctx context.Context, email, password string) (string, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `SELECT id, email, password_hash from users where email = $1`, email).Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", err
	}
	token := generateToken(user.Id)

	return token, nil
}

func (r *Repository) CreateMessage(ctx context.Context, server_id, channel_id, chat_id, userID int64, content string) error {
	query := `
	INSERT INTO messages (server_id, channel_id, chat_id, author_id, content)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, server_id, channel_id, chat_id, userID, content)
	if err != nil {
		return err
	}

	return nil
}
