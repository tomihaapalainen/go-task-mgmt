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

		log.Println("TASK IN::", taskIn)

		task := model.Task{
			ProjectID:  pID,
			AssigneeID: taskIn.AssigneeID,
			CreatorID:  user.ID,
			Title:      taskIn.Title,
			Content:    taskIn.Content,
			Status:     taskIn.Status,
		}
		log.Printf("TASK :: %+v", task)
		if err := task.Create(db); err != nil {
			log.Println("err creating task: ", err)
			return fmt.Errorf("error creating task")
		}

		return c.JSON(http.StatusOK, task)
	})
}