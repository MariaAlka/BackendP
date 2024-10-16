package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"        json:"id"`
	Name       string    `db:"name"      json:"name"`
	Email      string    `db:"email"     json:"email"`
	Phone      string    `db:"phone"     json:"phone"`
	Img        *string   `db:"img"       json:"img"`
	Password   string    `db:"password"  json:"-"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Updated_at time.Time `db:"updated_at" json:"updated_at"`
	// Add roles field to retrieve associated roles
	Roles []Role `db:"-" json:"roles,omitempty"` // Not stored in 'users' table, but useful for response
}

type Role struct {
	ID         int       `db:"id"        json:"id"`
	Name       string    `db:"name"      json:"name"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Updated_at time.Time `db:"updated_at" json:"updated_at"`
}

type UserRole struct {
	UserID uuid.UUID `db:"user_id" json:"user_id"`
	RoleID int       `db:"role_id" json:"role_id"`
}

// Response struct for consistent API responses
type Response struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// Vendor struct to store vendor information
type Vendor struct {
	ID          uuid.UUID `db:"id"        json:"id"`
	Name        string    `db:"name"      json:"name"`
	Img         *string   `db:"img"       json:"img"`
	Description string    `db:"description" json:"description"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
	Updated_at  time.Time `db:"updated_at" json:"updated_at"`
}

// Item represents an item in the store
type Item struct {
	ID        uuid.UUID `db:"id" json:"id"`
	VendorID  uuid.UUID `db:"vendor_id" json:"vendor_id"`
	Name      string    `db:"name" json:"name"`
	Price     float64   `db:"price" json:"price"`
	Img       *string   `db:"img" json:"img,omitempty"` // optional
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type VendorAdmin struct {
	UserID   uuid.UUID `db:"user_id"`
	VendorID uuid.UUID `db:"vendor_id"`
}

type Table struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	VendorID       uuid.UUID  `db:"vendor_id" json:"vendor_id"`
	Name           string     `db:"name" json:"name"`
	IsAvailable    bool       `db:"is_available" json:"is_available"`
	CustomerID     *uuid.UUID `db:"customer_id" json:"customer_id,omitempty"` // Pointer to allow NULL value
	IsNeedsService bool       `db:"is_needs_service" json:"is_needs_service"`
}

// OrderStatus defines the possible statuses for an order
type OrderStatus string

const (
	Completed OrderStatus = "completed"
	Preparing OrderStatus = "preparing"
)

// Order represents the structure of the 'orders' database table
type Order struct {
	ID             uuid.UUID   `db:"id" json:"id"`
	TotalOrderCost float64     `db:"total_order_cost" json:"total_order_cost"`
	CustomerID     uuid.UUID   `db:"customer_id" json:"customer_id"`
	VendorID       uuid.UUID   `db:"vendor_id" json:"vendor_id"`
	Status         OrderStatus `db:"status" json:"status"`
	CreatedAt      time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ID       uuid.UUID `db:"id" json:"id`
	OrderID  uuid.UUID `db:"order_id" json:"order_id"`
	ItemID   uuid.UUID `db:"item_id" json:"item_id`
	Quantity int       `db:"quantity" json:"quantity"`
	Price    float64   `db:"price" json:"price"`
}

// Cart represents a shopping cart in the system
type Cart struct {
	ID         uuid.UUID `db:"id" json:"id"`
	TotalPrice float64   `db:"total_price" json:"total_price"`
	Quantity   int       `db:"quantity" json:"quantity"`
	VendorID   uuid.UUID `db:"vendor_id" json:"vendor_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CartItem struct {
	CartID   uuid.UUID `db:"cart_id" json:"cart_id"`
	ItemID   uuid.UUID `db:"item_id" json:"item_id"`
	Quantity int       `db:"quantity" json:"quantity"`
}
