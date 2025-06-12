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

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
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
	err = db.AutoMigrate(&models.User{})
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

func TestUserRepository_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	// Test Create
	user := &models.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         models.RoleReceptionist,
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Test FindByID
	found, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.PasswordHash, found.PasswordHash)
	assert.Equal(t, user.Role, found.Role)

	// Test FindByUsername
	foundByUsername, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundByUsername.ID)

	// Test Update
	user.Username = "updateduser"
	err = repo.Update(user)
	assert.NoError(t, err)

	updated, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updated.Username)

	// Test List
	users, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, users, 1)

	// Test Delete
	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(user.ID)
	assert.Error(t, err)
}

func TestUserRepository_Concurrent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	// Create multiple users concurrently
	users := make([]*models.User, 5)
	for i := 0; i < 5; i++ {
		users[i] = &models.User{
			Username:     "testuser" + string(rune('0'+i)),
			PasswordHash: "hashedpassword",
			Role:         models.RoleReceptionist,
		}
	}

	// Create users concurrently
	for _, user := range users {
		go func(u *models.User) {
			err := repo.Create(u)
			assert.NoError(t, err)
		}(user)
	}

	// Wait for all goroutines to complete
	time.Sleep(100 * time.Millisecond)

	// Verify all users were created
	for _, user := range users {
		found, err := repo.FindByUsername(user.Username)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	}
}

func TestUserRepository_List(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	// Create 15 users
	for i := 0; i < 15; i++ {
		user := &models.User{
			Username:     "testuser" + string(rune('0'+i)),
			PasswordHash: "hashedpassword",
			Role:         models.RoleReceptionist,
		}
		err := repo.Create(user)
		assert.NoError(t, err)
	}

	// Test List
	users, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, users, 15)
}
