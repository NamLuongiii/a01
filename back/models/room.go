package models

import "gorm.io/gorm"

// Room model
type Room struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Users       []User `json:"users" gorm:"many2many:user_rooms;"`
}

// RoomRepository interface
type RoomRepository interface {
    Create(room *Room) error
    FindAll() ([]Room, error)
    FindByID(id uint) (*Room, error)
    AddUser(roomID uint, userID uint) error
}

// roomRepository implementation
type roomRepository struct {
    db *gorm.DB
}

// NewRoomRepository creates new room repository
func NewRoomRepository(db *gorm.DB) RoomRepository {
    return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *Room) error {
    return r.db.Create(room).Error
}

func (r *roomRepository) FindAll() ([]Room, error) {
    var rooms []Room
    err := r.db.Preload("Users").Find(&rooms).Error
    return rooms, err
}

func (r *roomRepository) FindByID(id uint) (*Room, error) {
    var room Room
    err := r.db.Preload("Users").First(&room, id).Error
    return &room, err
}

func (r *roomRepository) AddUser(roomID uint, userID uint) error {
    var room Room
    var user User
    
    if err := r.db.First(&room, roomID).Error; err != nil {
        return err
    }
    
    if err := r.db.First(&user, userID).Error; err != nil {
        return err
    }
    
    return r.db.Model(&room).Association("Users").Append(&user)
}