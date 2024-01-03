package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func HandlePostCreateTask(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		user := c.Get("user").(model.User)

		projectID := c.Param("projectID")
		pID, err := strconv.Atoi(projectID)
		if err != nil {
			return fmt.Errorf("invalid project ID '%s'", projectID)
		}

		taskIn := schema.TaskIn{}
		if err := json.NewDecoder(c.Request().Body).Decode(&taskIn); err != nil {
			log.Println("err decoding body: ", err)
			return c.JSON(http.StatusBadRequest, schema.MessageResponse{Message: "invalid request body"})
		}

		task := model.Task{
			ProjectID:  pID,
			AssigneeID: taskIn.AssigneeID,
			CreatorID:  user.ID,
			Title:      taskIn.Title,
			Content:    taskIn.Content,
			Status:     taskIn.Status,
		}
		if err := task.Create(db); err != nil {
			log.Println("err creating task: ", err)
			return fmt.Errorf("error creating task")
		}

		return c.JSON(http.StatusOK, task)
	})
}

func HandleDeleteTask(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		projectID := c.Param("projectID")
		pID, err := strconv.Atoi(projectID)
		if err != nil {
			return fmt.Errorf("invalid project ID '%s'", projectID)
		}
		taskID := c.Param("id")
		tID, err := strconv.Atoi(taskID)
		if err != nil {
			return fmt.Errorf("invalid project ID '%s'", taskID)
		}

		task := model.Task{ID: tID, ProjectID: pID}
		if err := task.Delete(db); err != nil {
			return fmt.Errorf("unable to delete project '%d' task '%d'", pID, tID)
		}

		return c.JSON(
			http.StatusNoContent,
			schema.MessageResponse{Message: fmt.Sprintf("task '%d' deleted successfully", tID)},
		)
	})
}
