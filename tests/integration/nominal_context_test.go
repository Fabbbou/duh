package integration_test

import (
	"duh/internal/application/contexts"
	"duh/internal/domain/service"
	"os"
	"testing"
)

func TestNominalContextInitialization(t *testing.T) {
	tmpDir := "test_context"
	defer os.RemoveAll(tmpDir)
	os.Mkdir("test_context", 0755)
	cliService, err := contexts.InitializeContexts(&tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize nominal context: %v", err)
	}
	if cliService == (service.CliService{}) {
		t.Fatalf("Initialized CliService is empty")
	}

	inj, err := cliService.Inject()
	if err != nil {
		t.Fatalf("Failed to inject CLI service: %v", err)
	}
	expectedInjection := "" // No additional repositories activated yet
	if inj != expectedInjection {
		t.Fatalf("Unexpected injection output. Got: %s, Expected: %s", inj, expectedInjection)
	}
}
