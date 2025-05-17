package handler

import (
	"net/http"
	"strconv"

	"dailyworkerroster/model"
	"dailyworkerroster/service"

	"github.com/gin-gonic/gin"
)

// ShiftHandler handles shift-related endpoints
type ShiftHandler struct {
	ShiftService service.ShiftServiceItf
}

// NewShiftHandler creates a new ShiftHandler
func NewShiftHandler(shiftService service.ShiftServiceItf) *ShiftHandler {
	return &ShiftHandler{ShiftService: shiftService}
}

// GetAssignedShifts godoc
// @Summary      Get assigned shifts for the current user
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.ListShiftDetail
// @Failure      500  {object}  map[string]string
// @Router       /worker/assigned [get]
func (h *ShiftHandler) GetAssignedShifts(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.ShiftService.GetAssignedShifts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAvailableShifts godoc
// @Summary      Get available shifts for a worker
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        workerID  path      int  true  "Worker ID"
// @Success      200  {array}   model.ShiftStatus
// @Failure      500  {object}  map[string]string
// @Router       /worker/{workerID}/available [get]
func (h *ShiftHandler) GetAvailableShifts(c *gin.Context) {
	workerID, _ := strconv.ParseInt(c.Param("workerID"), 10, 64)
	ctx := c.Request.Context()
	result, err := h.ShiftService.GetAvailableShifts(ctx, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// RequestShift godoc
// @Summary      Request a shift for a worker
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        shiftID   path      int  true  "Shift ID"
// @Param        workerID  path      int  true  "Worker ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /shift/{shiftID}/request/{workerID} [post]
func (h *ShiftHandler) RequestShift(c *gin.Context) {
	shiftID, _ := strconv.ParseInt(c.Param("shiftID"), 10, 64)
	workerID, _ := strconv.ParseInt(c.Param("workerID"), 10, 64)
	ctx := c.Request.Context()
	err := h.ShiftService.RequestShift(ctx, shiftID, workerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shift request submitted"})
}

// GetAllRequestedShifts godoc
// @Summary      Get all requested shifts for a worker
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        workerID  path      int  true  "Worker ID"
// @Success      200  {array}   model.ShiftStatus
// @Failure      500  {object}  map[string]string
// @Router       /worker/{workerID}/requests [get]
func (h *ShiftHandler) GetAllRequestedShifts(c *gin.Context) {
	workerID, _ := strconv.ParseInt(c.Param("workerID"), 10, 64)
	ctx := c.Request.Context()
	result, err := h.ShiftService.GetAllRequestedShift(ctx, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateShift godoc
// @Summary      Create a new shift
// @Tags         shifts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shift  body      model.Shift  true  "Shift"
// @Success      200  {object}  map[string]int64
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/shift [post]
func (h *ShiftHandler) CreateShift(c *gin.Context) {
	var shift model.Shift
	if err := c.ShouldBindJSON(&shift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	id, err := h.ShiftService.CreateShift(ctx, &shift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ApproveShiftRequest godoc
// @Summary      Approve a shift request for a worker
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        shiftID   path      int  true  "Shift ID"
// @Param        workerID  path      int  true  "Worker ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/shift/{shiftID}/approve/{workerID} [put]
func (h *ShiftHandler) ApproveShiftRequest(c *gin.Context) {
	shiftID, _ := strconv.ParseInt(c.Param("shiftID"), 10, 64)
	workerID, _ := strconv.ParseInt(c.Param("workerID"), 10, 64)
	ctx := c.Request.Context()
	err := h.ShiftService.ApproveShiftRequest(ctx, shiftID, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shift approved"})
}

// RejectShiftRequest godoc
// @Summary      Reject a shift request for a worker
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        shiftID   path      int  true  "Shift ID"
// @Param        workerID  path      int  true  "Worker ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/shift/{shiftID}/reject/{workerID} [put]
func (h *ShiftHandler) RejectShiftRequest(c *gin.Context) {
	shiftID, _ := strconv.ParseInt(c.Param("shiftID"), 10, 64)
	workerID, _ := strconv.ParseInt(c.Param("workerID"), 10, 64)
	ctx := c.Request.Context()
	err := h.ShiftService.RejectShiftRequest(ctx, shiftID, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shift rejected"})
}

// GetShiftsByDay godoc
// @Summary      Get all shifts by date
// @Tags         shifts
// @Produce      json
// @Security     BearerAuth
// @Param        date  query     string  true  "Date (YYYY-MM-DD)"
// @Success      200  {array}   model.ShiftStatus
// @Failure      500  {object}  map[string]string
// @Router       /admin/shifts/day [get]
func (h *ShiftHandler) GetShiftsByDay(c *gin.Context) {
	date := c.Query("date")
	ctx := c.Request.Context()
	result, err := h.ShiftService.GetShiftsByDay(ctx, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
