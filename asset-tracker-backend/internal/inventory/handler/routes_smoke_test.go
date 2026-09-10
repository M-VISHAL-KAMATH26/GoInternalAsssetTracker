package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// Registering a static child (/assets/catalog) next to a parameterized
// sibling (/assets/:id) must not panic at startup.
func TestRegisterAssetRoutesDoesNotConflict(t *testing.T) {
	router := gin.New()
	RegisterAssetRoutes(router, NewAssetHandler(&fakeAssetRepo{}))

	want := map[string]string{
		"GET /assets/catalog":      "",
		"GET /assets/:id":          "",
		"POST /assets":             "",
		"PATCH /assets/:id/retire": "",
	}
	for _, route := range router.Routes() {
		delete(want, route.Method+" "+route.Path)
	}
	if len(want) != 0 {
		t.Fatalf("missing routes: %v", want)
	}
}
