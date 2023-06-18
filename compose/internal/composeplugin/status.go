package composeplugin

import (
	"context"
	"encoding/json"
	"fmt"

	libstack "github.com/portainer/docker-compose-wrapper"
)

type publisher struct {
	URL           string
	TargetPort    int
	PublishedPort int
	Protocol      string
}

type service struct {
	ID         string
	Name       string
	Image      string
	Command    string
	Project    string
	Service    string
	Created    int64
	State      string
	Status     string
	Health     string
	ExitCode   int
	Publishers []publisher
}

// docker container state can be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
func getServiceStatus(service service) (libstack.Status, string) {
	switch service.State {
	case "created", "restarting", "paused":
		return libstack.StatusStarting, ""
	case "running":
		return libstack.StatusRunning, ""
	case "removing":
		return libstack.StatusRemoving, ""
	case "exited", "dead":
		if service.ExitCode != 0 {
			return libstack.StatusError, fmt.Sprintf("service %s exited with code %d", service.Name, service.ExitCode)
		}

		return libstack.StatusRemoved, ""
	default:
		return libstack.StatusUnknown, ""
	}
}

func aggregateStatuses(services []service) (libstack.Status, string) {
	statusCounts := make(map[libstack.Status]int)
	servicesCount := len(services)

	if servicesCount == 0 {
		return libstack.StatusRemoved, ""
	}

	errorMessage := ""
	for _, service := range services {
		status, serviceError := getServiceStatus(service)
		if serviceError != "" {
			errorMessage = serviceError
		}
		statusCounts[status]++
	}

	switch {
	case errorMessage != "":
		return libstack.StatusError, errorMessage
	case statusCounts[libstack.StatusStarting] > 0:
		return libstack.StatusStarting, ""
	case statusCounts[libstack.StatusRemoving] > 0:
		return libstack.StatusRemoving, ""
	case statusCounts[libstack.StatusRunning] == servicesCount:
		return libstack.StatusRunning, ""
	case statusCounts[libstack.StatusStopped] == servicesCount:
		return libstack.StatusStopped, ""
	case statusCounts[libstack.StatusRemoved] == servicesCount:
		return libstack.StatusRemoved, ""
	default:
		return libstack.StatusUnknown, ""
	}

}

func (wrapper *PluginWrapper) Status(ctx context.Context, projectName string) (libstack.Status, string, error) {
	output, err := wrapper.command(newCommand([]string{"ps", "-a", "--format", "json"}, nil), libstack.Options{
		ProjectName: projectName,
	})
	if len(output) == 0 || err != nil {
		return "", "", err
	}

	var services []service
	err = json.Unmarshal(output, &services)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse docker compose output: %w", err)
	}

	aggregateStatus, statusMessage := aggregateStatuses(services)
	return aggregateStatus, statusMessage, nil

}
