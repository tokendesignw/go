package models

import "time"

type CryptoProject struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Status    string    `bson:"status"`
	Urgency   string    `bson:"urgency"`
	Deadline  time.Time `bson:"deadline"`
	CreatedAt time.Time `bson:"created_at"`
	URL       string    `bson:"url"` // Добавлено поле для ссылки
}
