package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := VersionHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Add more assertions here if needed
}
