package handler

import (
	"net/http"
	"strconv"

	"dailyworkerroster/model"
	"dailyworkerroster/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	UserService service.UserServiceItf
}

func NewUserHandler(userService service.UserServiceItf) *UserHandler {
	return &UserHandler{UserService: userService}
}

// SignUp godoc
// @Summary      Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "User"
// @Success      200   {object}  map[string]int64
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.UserService.SignUp(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Login godoc
// @Summary      Login a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body  object{identifier=string,password=string}  true  "Login credentials"
// @Success      200   {object}  model.User
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.UserService.Login(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAllWorkers godoc
// @Summary      Get all workers
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200   {array}   model.User
// @Failure      500   {object}  map[string]string
// @Router       /workers [get]
func (h *UserHandler) GetAllWorkers(c *gin.Context) {
	users, err := h.UserService.GetAllWorkers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetWorkerByID godoc
// @Summary      Get worker by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Worker ID"
// @Success      200  {object}  model.User
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /worker/{id} [get]
func (h *UserHandler) GetWorkerByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}
	user, err := h.UserService.GetWorkerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
