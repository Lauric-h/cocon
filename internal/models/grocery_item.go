package models

import (
	"time"
)

type GroceryItem struct {
	ID        int
	Name      string
	Quantity  *int
	Checked   bool
	CreatedAt time.Time
}
