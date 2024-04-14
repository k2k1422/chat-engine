package Model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
}
