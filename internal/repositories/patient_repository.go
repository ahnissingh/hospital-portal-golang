package repositories

import (
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
	List(page, limit int) ([]models.Patient, int64, error)
	Search(params models.PatientSearchRequest) ([]models.Patient, error)
	ExistsByNameOrContact(name, contactInfo string) (bool, error)
}

// patientRepository implements PatientRepository interface
type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{
		db: db,
	}
}

// Create creates a new patient
func (r *patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

// FindByID finds a patient by ID
func (r *patientRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.First(&patient, id).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// Update updates a patient
func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

// UpdateMedicalNotes updates medical notes
func (r *patientRepository) UpdateMedicalNotes(id uint, medicalNotes string) error {
	return r.db.Model(&models.Patient{}).Where("id = ?", id).Update("medical_notes", medicalNotes).Error
}

// Delete deletes a patient
func (r *patientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

// List returns all patients with pagination
func (r *patientRepository) List(page, limit int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	// Count total records
	if err := r.db.Model(&models.Patient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Find(&patients).Error
	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// Search searches for patients
func (r *patientRepository) Search(params models.PatientSearchRequest) ([]models.Patient, error) {
	var patients []models.Patient
	query := r.db.Model(&models.Patient{})

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

	err := query.Find(&patients).Error
	return patients, err
}

// ExistsByNameOrContact checks if a patient exists with the given name or contact info
func (r *patientRepository) ExistsByNameOrContact(name, contactInfo string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Patient{}).
		Where("name = ? OR contact_info = ?", name, contactInfo).
		Count(&count).Error
	return count > 0, err
}
