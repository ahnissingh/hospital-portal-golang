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
	List() ([]models.Patient, error)
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
	// Validate patient data
	if patient.Name == "" {
		return errors.New("patient name is required")
	}
	if patient.Age < 0 || patient.Age > 150 {
		return errors.New("invalid patient age")
	}
	if patient.Gender == "" {
		return errors.New("patient gender is required")
	}
	if patient.ContactInfo == "" {
		return errors.New("patient contact information is required")
	}
	if patient.CreatedBy == 0 {
		return errors.New("creator ID is required")
	}

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
	if id == 0 {
		return nil, errors.New("invalid patient ID")
	}
	return s.patientRepo.FindByID(id)
}

// Update updates a patient
func (s *patientService) Update(patient *models.Patient) error {
	// Validate patient data
	if patient.ID == 0 {
		return errors.New("invalid patient ID")
	}
	if patient.Name == "" {
		return errors.New("patient name is required")
	}
	if patient.Age < 0 || patient.Age > 150 {
		return errors.New("invalid patient age")
	}
	if patient.Gender == "" {
		return errors.New("patient gender is required")
	}
	if patient.ContactInfo == "" {
		return errors.New("patient contact information is required")
	}

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

// List returns all patients
func (s *patientService) List() ([]models.Patient, error) {
	return s.patientRepo.List()
}

// Search searches for patients based on search parameters
func (s *patientService) Search(params models.PatientSearchRequest) ([]models.Patient, error) {
	// Validate search parameters
	if params.AgeMin < 0 {
		return nil, errors.New("minimum age cannot be negative")
	}
	if params.AgeMax > 150 {
		return nil, errors.New("maximum age cannot exceed 150")
	}
	if params.AgeMin > params.AgeMax && params.AgeMax > 0 {
		return nil, errors.New("minimum age cannot be greater than maximum age")
	}

	return s.patientRepo.Search(params)
}
