package model

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Status    string `json:"status"`
	UserID    uint
	// User      User      `json:"users,omitempty"`
	Visitors []Visitor `json:"visitors,omitempty"`
}

func (s *Staff) PublicStaff() *Staff {
	staff := &Staff{
		Firstname: s.Firstname,
		Lastname:  s.Lastname,
		Email:     s.Email,
		Phone:     s.Phone,
		Address:   s.Address,
		UserID:    s.UserID,
		Status:    s.Status,
	}
	staff.ID = s.ID
	staff.CreatedAt = s.CreatedAt
	return staff
}

func (Staff) TableName() string {
	return "staff"
}
