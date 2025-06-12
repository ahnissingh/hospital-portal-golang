package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hospital-project/internal/middleware"
	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// PatientController handles patient requests
type PatientController struct {
	patientService services.PatientService
	authMiddleware *middleware.AuthMiddleware
}

// NewPatientController creates a new patient controller
func NewPatientController(patientService services.PatientService, authMiddleware *middleware.AuthMiddleware) *PatientController {
	return &PatientController{
		patientService: patientService,
		authMiddleware: authMiddleware,
	}
}

// @Summary Create patient
// @Description Create a new patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Param request body models.CreatePatientRequest true "Create Patient Request"
// @Success 201 {object} models.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients [post]
// @Security Bearer
func (c *PatientController) CreatePatient(ctx *gin.Context) {
	var request models.CreatePatientRequest

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user
	currentUser, ok := middleware.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Map request to patient model
	patient := &models.Patient{
		Name:         request.Name,
		Age:          request.Age,
		Gender:       request.Gender,
		ContactInfo:  request.ContactInfo,
		MedicalNotes: request.MedicalNotes,
		CreatedBy:    currentUser.ID,
	}

	err := c.patientService.Create(patient)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, patient.ToResponse())
}

// @Summary Get patient by ID
// @Description Get a patient by ID (Both Receptionist and Doctor)
// @Tags patients
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} models.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients/{id} [get]
// @Security Bearer
func (c *PatientController) GetPatient(ctx *gin.Context) {
	// Get ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Get patient
	patient, err := c.patientService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	ctx.JSON(http.StatusOK, patient.ToResponse())
}

// @Summary Update patient
// @Description Update a patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param request body models.UpdatePatientRequest true "Update Patient Request"
// @Success 200 {object} models.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients/{id} [put]
// @Security Bearer
func (c *PatientController) UpdatePatient(ctx *gin.Context) {
	// Get ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var patient models.Patient

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set ID from path
	patient.ID = uint(id)

	// Update patient
	err = c.patientService.Update(&patient)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, patient.ToResponse())
}

// @Summary Update medical notes
// @Description Update a patient's medical notes (Doctor only)
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param request body models.UpdateMedicalNotesRequest true "Update Medical Notes Request"
// @Success 200 {object} models.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients/{id}/medical-notes [put]
// @Security Bearer
func (c *PatientController) UpdateMedicalNotes(ctx *gin.Context) {
	// Get ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var request models.UpdateMedicalNotesRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update medical notes
	err = c.patientService.UpdateMedicalNotes(uint(id), request.MedicalNotes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated patient
	patient, err := c.patientService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated patient"})
		return
	}

	ctx.JSON(http.StatusOK, patient.ToResponse())
}

// @Summary Delete patient
// @Description Delete a patient (Receptionist only)
// @Tags patients
// @Param id path int true "Patient ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients/{id} [delete]
// @Security Bearer
func (c *PatientController) DeletePatient(ctx *gin.Context) {
	// Get ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Delete patient
	err = c.patientService.Delete(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary List patients
// @Description List all patients (Both Receptionist and Doctor)
// @Tags patients
// @Produce json
// @Success 200 {array} models.PatientResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients [get]
// @Security Bearer
func (c *PatientController) ListPatients(ctx *gin.Context) {
	// Get patients
	patients, err := c.patientService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patients"})
		return
	}

	// Convert to response
	var response []models.PatientResponse
	for _, patient := range patients {
		response = append(response, patient.ToResponse())
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Search patients
// @Description Search for patients (Receptionist only)
// @Tags patients
// @Produce json
// @Param name query string false "Patient name"
// @Param age_min query int false "Minimum age"
// @Param age_max query int false "Maximum age"
// @Param gender query string false "Patient gender"
// @Param contact_info query string false "Contact information"
// @Success 200 {array} models.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/patients/search [get]
// @Security Bearer
func (c *PatientController) SearchPatients(ctx *gin.Context) {
	var request models.PatientSearchRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameters"})
		return
	}

	// Search patients
	patients, err := c.patientService.Search(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response
	var response []models.PatientResponse
	for _, patient := range patients {
		response = append(response, patient.ToResponse())
	}

	ctx.JSON(http.StatusOK, response)
}

// RegisterRoutes registers the patient routes
func (c *PatientController) RegisterRoutes(router *gin.Engine) {
	patients := router.Group("/api/patients")
	patients.Use(c.authMiddleware.Authenticate())
	{
		// Routes for both receptionist and doctor
		patients.GET("", c.ListPatients)
		patients.GET("/:id", c.GetPatient)

		// Routes for receptionist only
		receptionistRoutes := patients.Group("")
		receptionistRoutes.Use(c.authMiddleware.RequireRole(models.RoleReceptionist))
		{
			receptionistRoutes.POST("", c.CreatePatient)
			receptionistRoutes.PUT("/:id", c.UpdatePatient)
			receptionistRoutes.DELETE("/:id", c.DeletePatient)
			receptionistRoutes.GET("/search", c.SearchPatients)
		}

		// Routes for doctor only
		doctorRoutes := patients.Group("")
		doctorRoutes.Use(c.authMiddleware.RequireRole(models.RoleDoctor))
		{
			doctorRoutes.PUT("/:id/medical-notes", c.UpdateMedicalNotes)
		}
	}
}
