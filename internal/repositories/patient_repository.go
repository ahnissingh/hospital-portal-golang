package repositories

import (
	"errors"

	"gorm.io/gorm"

	"hospital-project/internal/models"
)

// PatientRepository interface defines methods for patient repository
type PatientRepository interface {
	Create(patient *models.Patient) error
	FindByID(id uint) (*models.Patient, error)
	Update(patient *models.Patient) error
	UpdateMedicalNotes(id uint, medicalNotes string) error
	Delete(id uint) error
	List() ([]models.Patient, error)
	Search(params models.PatientSearchRequest) ([]models.Patient, error)
	ExistsByNameOrContact(name string, contactInfo string) (bool, error)
}

// patientRepository implements PatientRepository interface
type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

// Create creates a new patient
func (r *patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

// FindByID finds a patient by ID
func (r *patientRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	result := r.db.First(&patient, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, result.Error
	}
	return &patient, nil
}

// Update updates a patient
func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

// UpdateMedicalNotes updates only the medical notes of a patient
func (r *patientRepository) UpdateMedicalNotes(id uint, medicalNotes string) error {
	result := r.db.Model(&models.Patient{}).Where("id = ?", id).Update("medical_notes", medicalNotes)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("patient not found")
	}
	return nil
}

// Delete deletes a patient
func (r *patientRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Patient{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("patient not found")
	}
	return nil
}

// List returns all patients
func (r *patientRepository) List() ([]models.Patient, error) {
	var patients []models.Patient
	result := r.db.Find(&patients)
	return patients, result.Error
}

// Search searches for patients based on search parameters
func (r *patientRepository) Search(params models.PatientSearchRequest) ([]models.Patient, error) {
	var patients []models.Patient
	query := r.db.Model(&models.Patient{})

	// Apply filters if provided
	if params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.AgeMin > 0 {
		query = query.Where("age >= ?", params.AgeMin)
	}
	if params.AgeMax > 0 {
		query = query.Where("age <= ?", params.AgeMax)
	}
	if params.Gender != "" {
		query = query.Where("gender = ?", params.Gender)
	}
	if params.ContactInfo != "" {
		query = query.Where("contact_info ILIKE ?", "%"+params.ContactInfo+"%")
	}

	result := query.Find(&patients)
	return patients, result.Error
}

// ExistsByNameOrContact checks if a patient with the given name or contact already exists
func (r *patientRepository) ExistsByNameOrContact(name string, contactInfo string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Patient{}).
		Where("name = ? OR contact_info = ?", name, contactInfo).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
