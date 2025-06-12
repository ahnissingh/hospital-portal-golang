package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// MockPatientRepository is a mock implementation of the PatientRepository interface
type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) FindByID(id uint) (*models.Patient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) UpdateMedicalNotes(id uint, medicalNotes string) error {
	args := m.Called(id, medicalNotes)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) List(page, limit int) ([]models.Patient, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Patient), args.Get(1).(int64), args.Error(2)
}

func (m *MockPatientRepository) Search(params models.PatientSearchRequest) ([]models.Patient, error) {
	args := m.Called(params)
	return args.Get(0).([]models.Patient), args.Error(1)
}

func (m *MockPatientRepository) ExistsByNameOrContact(name, contactInfo string) (bool, error) {
	args := m.Called(name, contactInfo)
	return args.Bool(0), args.Error(1)
}

func TestPatientService_Create_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patient
	patient := &models.Patient{
		Name:         "John Doe",
		ContactInfo:  "1234567890",
		Age:          30,
		Gender:       models.GenderMale,
		CreatedBy:    1,
	}

	// Set up expectations
	mockRepo.On("ExistsByNameOrContact", patient.Name, patient.ContactInfo).Return(false, nil)
	mockRepo.On("Create", patient).Return(nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	err := patientService.Create(patient)

	// Assert expectations
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_Create_DuplicatePatient(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patient
	patient := &models.Patient{
		Name:         "John Doe",
		ContactInfo:  "1234567890",
		Age:          30,
		Gender:       models.GenderMale,
		CreatedBy:    1,
	}

	// Set up expectations
	mockRepo.On("ExistsByNameOrContact", patient.Name, patient.ContactInfo).Return(true, nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	err := patientService.Create(patient)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "patient with this name or contact info already exists", err.Error())

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_GetByID_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patient
	patient := &models.Patient{
		Name:         "John Doe",
		ContactInfo:  "1234567890",
		Age:          30,
		Gender:       models.GenderMale,
		CreatedBy:    1,
	}

	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(patient, nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	result, err := patientService.GetByID(1)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, patient.Name, result.Name)
	assert.Equal(t, patient.ContactInfo, result.ContactInfo)
	assert.Equal(t, patient.Age, result.Age)
	assert.Equal(t, patient.Gender, result.Gender)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_GetByID_NotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	result, err := patientService.GetByID(1)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "not found", err.Error())

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_UpdateMedicalNotes_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patient
	patient := &models.Patient{
		Name:         "John Doe",
		MedicalNotes: "Initial notes",
		CreatedBy:    1,
	}

	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(patient, nil)
	mockRepo.On("UpdateMedicalNotes", uint(1), "Updated notes").Return(nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	err := patientService.UpdateMedicalNotes(1, "Updated notes")

	// Assert expectations
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_List_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patients
	patients := []models.Patient{
		{
			Name:         "John Doe",
			ContactInfo:  "1234567890",
			Age:          30,
			Gender:       models.GenderMale,
			CreatedBy:    1,
		},
		{
			Name:         "Jane Smith",
			ContactInfo:  "0987654321",
			Age:          25,
			Gender:       models.GenderFemale,
			CreatedBy:    1,
		},
	}

	// Set up expectations
	mockRepo.On("List", 1, 10).Return(patients, int64(2), nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	result, total, err := patientService.List(1, 10)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	assert.Equal(t, patients[0].Name, result[0].Name)
	assert.Equal(t, patients[1].Name, result[1].Name)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestPatientService_Search_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)

	// Create test patients
	patients := []models.Patient{
		{
			Name:         "John Doe",
			ContactInfo:  "1234567890",
			Age:          30,
			Gender:       models.GenderMale,
			CreatedBy:    1,
		},
	}

	// Create search parameters
	params := models.PatientSearchRequest{
		Name: "John",
	}

	// Set up expectations
	mockRepo.On("Search", params).Return(patients, nil)

	// Create patient service with mock repository
	patientService := services.NewPatientService(mockRepo)

	// Call the method being tested
	result, err := patientService.Search(params)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, patients[0].Name, result[0].Name)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}
