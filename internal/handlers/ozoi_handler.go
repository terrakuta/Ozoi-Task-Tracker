package handlers

import (
	"Ozoi/internal/dto"
	"Ozoi/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTaskHandler godoc
// @Summary      Create a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateOzoiInput  true  "Task data"
// @Success      201    {object}  models.OzoiTask
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     CookieAuth
// @Router       /ozoi [post]
func CreateTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in gin context"})
			return
		}

		userId := userInterface.(string)

		var input dto.CreateOzoiInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return

		}

		if err := input.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task, err := repository.CreateTask(pool, input.Title, input.Completed, input.Description, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

// GetAllTasksHandler godoc
// @Summary      Get all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   models.OzoiTask
// @Failure      500  {object}  map[string]string
// @Security     CookieAuth
// @Router       /ozoi [get]
func GetAllTasksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in gin context"})
			return
		}

		userId := userInterface.(string)

		tasks, err := repository.GetAllTasks(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

// GetTaskByIDHandler godoc
// @Summary      Get task by ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.OzoiTask
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     CookieAuth
// @Router       /ozoi/{id} [get]
func GetTaskByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in gin context"})
			return
		}

		userId := userInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
			return
		}
		task, err := repository.GetTaskByID(pool, id, userId)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"Error": "Row is not found"})
			}

			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}
		c.JSON(http.StatusOK, task)
	}
}

// UpdateTaskByIDHandler godoc
// @Summary      Update task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id     path      int                  true  "Task ID"
// @Param        input  body      dto.UpdateOzoiInput  true  "Updated task data"
// @Success      200    {object}  models.OzoiTask
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     CookieAuth
// @Router       /ozoi/{id} [put]
func UpdateTaskByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in gin context"})
			return
		}

		userId := userInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
			return
		}

		var input dto.UpdateOzoiInput

		if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": bindErr.Error()})
			return
		}

		if validateErr := input.Validate(); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		descriptionValue := ""
		if input.Description != nil {
			descriptionValue = *input.Description
		}

		task, err := repository.UpdateTaskByID(pool, id, input.Title, descriptionValue, input.Completed, userId)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"Error": "Task not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}
		c.JSON(http.StatusOK, task)
	}
}

// DeleteTaskByIDHandler godoc
// @Summary      Delete task by ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     CookieAuth
// @Router       /ozoi/{id} [delete]
func DeleteTaskByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in gin context"})
			return
		}

		userId := userInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
			return
		}

		err = repository.DeleteTaskByID(pool, id, userId)

		if err != nil {
			if err.Error() == "Task "+idString+" is not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task is not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task deletion success"})
	}
}
