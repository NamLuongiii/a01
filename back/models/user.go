package models

import "gorm.io/gorm"

// User model
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
    Rooms []Room `json:"rooms" gorm:"many2many:user_rooms;"`
}

// UserRepository interface
type UserRepository interface {
    Create(user *User) error
    FindAll() ([]User, error)
    FindByID(id uint) (*User, error)
}

// userRepository implementation
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository creates new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]User, error) {
    var users []User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *userRepository) FindByID(id uint) (*User, error) {
    var user User
    err := r.db.First(&user, id).Error
    return &user, err
}