package health

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type healthMonitor interface {
	Check() healthCheckData
}

// Echo bindings for the monitor
type HealthController struct {
	Monitor healthMonitor
}

func (controller HealthController) EchoHandler(c echo.Context) error {
	result := controller.Monitor.Check()
	if result.Status != StatusUp {
		return c.JSON(http.StatusServiceUnavailable, result)
	}

	return c.JSON(http.StatusOK, result)
}

func Loobpack(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]MonitoredStatus{"status": StatusUp})
}

// This is mostly some helper, which can be done manually by the user of this library
func EchoRegister(server *echo.Echo, controller HealthController, liveness string, readiness string) {
	server.GET(liveness, Loobpack)
	server.GET(readiness, controller.EchoHandler)
}
