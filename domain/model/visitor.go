package model

import "gorm.io/gorm"

type Visitor struct {
	gorm.Model
	Firstname  string `json:"first_name"`
	Lastname   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Occupation string `json:"occupation"`
	Note       string `json:"note"`
	StaffID    uint   `json:"staff_id"`
	// Staff      Staff  `json:"staff,omitempty"`
	UserID uint
	// User       User `json:"user,omitempty"`
}

func (Visitor) TableName() string {
	return "visitors"
}

func (v Visitor) PublicVisitor() *Visitor {
	visitor := &Visitor{
		Firstname:  v.Firstname,
		Lastname:   v.Lastname,
		Email:      v.Email,
		Phone:      v.Phone,
		Occupation: v.Occupation,
		Note:       v.Note,
		StaffID:    v.StaffID,
		UserID:     v.UserID,
	}
	visitor.ID = v.ID
	visitor.CreatedAt = v.CreatedAt

	return visitor
}
