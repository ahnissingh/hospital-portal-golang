package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hospital-project/internal/models"
	"hospital-project/internal/repositories"
)

func setupPatientTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get container host and port
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Connect to database
	dsn := "host=" + host + " port=" + port.Port() + " user=test password=test dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&models.Patient{})
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close()
		container.Terminate(ctx)
	}

	return db, cleanup
}

func TestPatientRepository_CRUD(t *testing.T) {
	db, cleanup := setupPatientTestDB(t)
	defer cleanup()

	repo := repositories.NewPatientRepository(db)

	// Test Create
	patient := &models.Patient{
		Name:        "John Doe",
		Age:         30,
		Gender:      models.GenderMale,
		ContactInfo: "1234567890",
		CreatedBy:   1,
	}

	err := repo.Create(patient)
	assert.NoError(t, err)
	assert.NotZero(t, patient.ID)

	// Test FindByID
	found, err := repo.FindByID(patient.ID)
	assert.NoError(t, err)
	assert.Equal(t, patient.Name, found.Name)
	assert.Equal(t, patient.Age, found.Age)
	assert.Equal(t, patient.Gender, found.Gender)
	assert.Equal(t, patient.ContactInfo, found.ContactInfo)

	// Test Update
	patient.Name = "Jane Doe"
	err = repo.Update(patient)
	assert.NoError(t, err)

	updated, err := repo.FindByID(patient.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", updated.Name)

	// Test UpdateMedicalNotes
	err = repo.UpdateMedicalNotes(patient.ID, "Updated notes")
	assert.NoError(t, err)
	updated, err = repo.FindByID(patient.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated notes", updated.MedicalNotes)

	// Test ExistsByNameOrContact
	exists, err := repo.ExistsByNameOrContact("Jane Doe", "1234567890")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test List
	patients, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, patients, 1)

	// Test Search
	searchParams := models.PatientSearchRequest{Name: "Jane"}
	results, err := repo.Search(searchParams)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Jane Doe", results[0].Name)

	// Test Delete
	err = repo.Delete(patient.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(patient.ID)
	assert.Error(t, err)
}

func TestPatientRepository_Concurrent(t *testing.T) {
	db, cleanup := setupPatientTestDB(t)
	defer cleanup()

	repo := repositories.NewPatientRepository(db)

	// Create multiple patients concurrently
	patients := make([]*models.Patient, 5)
	for i := 0; i < 5; i++ {
		patients[i] = &models.Patient{
			Name:        "Patient " + string(rune('0'+i)),
			Age:         30 + i,
			Gender:      models.GenderMale,
			ContactInfo: "123456789" + string(rune('0'+i)),
			CreatedBy:   1,
		}
	}

	// Create patients concurrently
	for _, patient := range patients {
		go func(p *models.Patient) {
			err := repo.Create(p)
			assert.NoError(t, err)
		}(patient)
	}

	// Wait for all goroutines to complete
	time.Sleep(100 * time.Millisecond)

	// Verify all patients were created
	for _, patient := range patients {
		found, err := repo.FindByID(patient.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	}
}

func TestPatientRepository_ListPagination(t *testing.T) {
	db, cleanup := setupPatientTestDB(t)
	defer cleanup()

	repo := repositories.NewPatientRepository(db)

	// Create 15 patients
	for i := 0; i < 15; i++ {
		patient := &models.Patient{
			Name:        "Patient " + string(rune('0'+i)),
			Age:         30 + i,
			Gender:      models.GenderMale,
			ContactInfo: "123456789" + string(rune('0'+i)),
			CreatedBy:   1,
		}
		err := repo.Create(patient)
		assert.NoError(t, err)
	}

	// Test first page
	patients, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, patients, 10)

	// Test second page
	patients, total, err = repo.List(2, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, patients, 5)
}

func TestPatientRepository_Search(t *testing.T) {
	db, cleanup := setupPatientTestDB(t)
	defer cleanup()

	repo := repositories.NewPatientRepository(db)

	// Create test patients
	patients := []*models.Patient{
		{
			Name:        "John Smith",
			Age:         30,
			Gender:      models.GenderMale,
			ContactInfo: "1234567890",
			CreatedBy:   1,
		},
		{
			Name:        "Jane Smith",
			Age:         25,
			Gender:      models.GenderFemale,
			ContactInfo: "0987654321",
			CreatedBy:   1,
		},
		{
			Name:        "Bob Johnson",
			Age:         40,
			Gender:      models.GenderMale,
			ContactInfo: "5555555555",
			CreatedBy:   1,
		},
	}

	for _, patient := range patients {
		err := repo.Create(patient)
		assert.NoError(t, err)
	}

	// Test search by name
	results, err := repo.Search(models.PatientSearchRequest{Name: "Smith"})
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Test search by gender
	results, err = repo.Search(models.PatientSearchRequest{Gender: models.GenderFemale})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Jane Smith", results[0].Name)

	// Test search by age range
	results, err = repo.Search(models.PatientSearchRequest{AgeMin: 35})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Bob Johnson", results[0].Name)
}
