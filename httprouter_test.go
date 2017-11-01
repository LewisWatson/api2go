// +build !gingonic,!gorillamux,!echo

package api2go

import (
	"log"

	"github.com/LewisWatson/api2go/routing"
)

func newTestRouter() routing.Routeable {
	return routing.NewHTTPRouter(testPrefix, &notAllowedHandler{})
}

func init() {
	log.Println("Testing with default router")
}
