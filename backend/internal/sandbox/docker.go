package sandbox

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type SandboxManager struct {
	cli *client.Client
}

func NewSandboxManager() (*SandboxManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &SandboxManager{cli: cli}, nil
}

// StartSandbox creates and starts a new container for the agent
func (s *SandboxManager) StartSandbox(ctx context.Context, imageName, workspaceHostPath string) (string, error) {
	log.Printf("Starting sandbox with image %s", imageName)

	resp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"tail", "-f", "/dev/null"},
		Tty:   false,
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/workspace", workspaceHostPath)},
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// ExecuteCommand runs a command inside the container and streams output
func (s *SandboxManager) ExecuteCommand(ctx context.Context, containerID string, cmd []string) (io.Reader, error) {
	execOpts := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		WorkingDir:   "/workspace",
	}

	execResp, err := s.cli.ContainerExecCreate(ctx, containerID, execOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec: %w", err)
	}

	attachResp, err := s.cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach to exec: %w", err)
	}

	return attachResp.Reader, nil
}

func (s *SandboxManager) StopSandbox(ctx context.Context, containerID string) error {
	return s.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}
