package router

import (
	"net/http"
	"net/http/httptest"
	"teamflow/internal/controller"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTaskRoutesRequireJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := SetupRouter(
		&controller.UserController{},
		&controller.TeamController{},
		&controller.ProjectController{},
		&controller.TaskController{},
	)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?project_id=1", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestLegacyPublicTaskRouteIsNotRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := SetupRouter(
		&controller.UserController{},
		&controller.TeamController{},
		&controller.ProjectController{},
		&controller.TaskController{},
	)

	request := httptest.NewRequest(http.MethodGet, "/tasks?project_id=1", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
