package project_test

import (
	"strings"
	"testing"

	"yellow-jersey/pkg/project"
)

func TestBase(t *testing.T) {
	path := "strava-api"
	if got := project.Root(); strings.Index(got, path) < 1 {
		t.Errorf("wrong project path %v", got)
	}
}
