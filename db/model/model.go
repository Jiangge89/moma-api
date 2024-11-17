package model

type Rate struct {
	ID        int     `json:"id"`
	FromCode  string  `json:"from_code"`
	ToCode    string  `json:"to_code"`
	Rate      float32 `json:"rate"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

type Account struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
