package models

import (
	"time"

	"gorm.io/gorm"
)

// Gender type for patient gender
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Patient represents a patient in the system
type Patient struct {
	gorm.Model
	Name         string `gorm:"not null" json:"name" binding:"required"`
	Age          int    `gorm:"not null" json:"age" binding:"required,min=0,max=150"`
	Gender       Gender `gorm:"not null" json:"gender" binding:"required,oneof=male female other"`
	ContactInfo  string `gorm:"not null" json:"contact_info" binding:"required"`
	MedicalNotes string `json:"medical_notes"`
	CreatedBy    uint   `gorm:"not null" json:"created_by" binding:"required"`
}

// TableName overrides the table name
func (Patient) TableName() string {
	return "patients"
}

// PatientResponse is the DTO for patient responses
type PatientResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Age          int       `json:"age"`
	Gender       Gender    `json:"gender"`
	ContactInfo  string    `json:"contact_info"`
	MedicalNotes string    `json:"medical_notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse converts a Patient to a PatientResponse
func (p *Patient) ToResponse() PatientResponse {
	return PatientResponse{
		ID:           p.ID,
		Name:         p.Name,
		Age:          p.Age,
		Gender:       p.Gender,
		ContactInfo:  p.ContactInfo,
		MedicalNotes: p.MedicalNotes,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// CreatePatientRequest is the DTO for creating a patient
type CreatePatientRequest struct {
	Name         string `json:"name" binding:"required"`
	Age          int    `json:"age" binding:"required,min=0,max=150"`
	Gender       Gender `json:"gender" binding:"required"`
	ContactInfo  string `json:"contact_info" binding:"required"`
	MedicalNotes string `json:"medical_notes"`
}

// UpdatePatientRequest is the DTO for updating a patient
type UpdatePatientRequest struct {
	Name         string `json:"name"`
	Age          int    `json:"age" binding:"omitempty,min=0,max=150"`
	Gender       Gender `json:"gender"`
	ContactInfo  string `json:"contact_info"`
	MedicalNotes string `json:"medical_notes"`
}

// UpdateMedicalNotesRequest is the DTO for updating medical notes
type UpdateMedicalNotesRequest struct {
	MedicalNotes string `json:"medical_notes" binding:"required"`
}

// PatientSearchRequest is the DTO for searching patients
type PatientSearchRequest struct {
	Name        string `form:"name" binding:"omitempty"`
	AgeMin      int    `form:"age_min" binding:"omitempty,min=0,max=150"`
	AgeMax      int    `form:"age_max" binding:"omitempty,min=0,max=150,gtefield=AgeMin"`
	Gender      Gender `form:"gender" binding:"omitempty,oneof=male female other"`
	ContactInfo string `form:"contact_info" binding:"omitempty"`
}
