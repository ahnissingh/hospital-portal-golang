package services

import (
	"errors"

	"hospital-project/internal/models"
	"hospital-project/internal/repositories"
)

// PatientService interface defines methods for patient service
type PatientService interface {
	Create(patient *models.Patient) error
	GetByID(id uint) (*models.Patient, error)
	Update(patient *models.Patient) error
	UpdateMedicalNotes(id uint, medicalNotes string) error
	Delete(id uint) error
	List(page, limit int) ([]models.Patient, int64, error)
	Search(params models.PatientSearchRequest) ([]models.Patient, error)
}

// patientService implements PatientService interface
type patientService struct {
	patientRepo repositories.PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(patientRepo repositories.PatientRepository) PatientService {
	return &patientService{
		patientRepo: patientRepo,
	}
}

// Create creates a new patient
func (s *patientService) Create(patient *models.Patient) error {
	exists, err := s.patientRepo.ExistsByNameOrContact(patient.Name, patient.ContactInfo)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("patient with this name or contact info already exists")
	}
	// Save patient to database
	return s.patientRepo.Create(patient)
}

// GetByID gets a patient by ID
func (s *patientService) GetByID(id uint) (*models.Patient, error) {
	return s.patientRepo.FindByID(id)
}

// Update updates a patient
func (s *patientService) Update(patient *models.Patient) error {
	// Check if patient exists
	existingPatient, err := s.patientRepo.FindByID(patient.ID)
	if err != nil {
		return err
	}
	if existingPatient == nil {
		return errors.New("patient not found")
	}

	// Update patient in database
	return s.patientRepo.Update(patient)
}

// UpdateMedicalNotes updates only the medical notes of a patient
func (s *patientService) UpdateMedicalNotes(id uint, medicalNotes string) error {
	if id == 0 {
		return errors.New("invalid patient ID")
	}

	// Check if patient exists
	existingPatient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingPatient == nil {
		return errors.New("patient not found")
	}

	// Update medical notes in database
	return s.patientRepo.UpdateMedicalNotes(id, medicalNotes)
}

// Delete deletes a patient
func (s *patientService) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid patient ID")
	}
	return s.patientRepo.Delete(id)
}

// List returns all patients with pagination
func (s *patientService) List(page, limit int) ([]models.Patient, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.patientRepo.List(page, limit)
}

// Search searches for patients based on search parameters
func (s *patientService) Search(params models.PatientSearchRequest) ([]models.Patient, error) {
	return s.patientRepo.Search(params)
}
