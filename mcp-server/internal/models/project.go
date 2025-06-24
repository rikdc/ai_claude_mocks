package models

import (
	"time"

	"github.com/example/mockery-mcp-server/internal/types"
)

// MockeryProject represents a project with mockery configuration
type MockeryProject struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Path        string                 `json:"path"`
	Config      types.MockeryConfig    `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Interfaces  []types.InterfaceDefinition `json:"interfaces"`
}

// GeneratedMock represents a generated mock file
type GeneratedMock struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	InterfaceName string    `json:"interface_name"`
	PackagePath   string    `json:"package_path"`
	FilePath      string    `json:"file_path"`
	GeneratedAt   time.Time `json:"generated_at"`
	MockeryVersion string   `json:"mockery_version"`
	Hash          string    `json:"hash"` // Hash of the generated content for change detection
}

// InterfaceRegistry manages discovered interfaces for a project
type InterfaceRegistry struct {
	ProjectID   string                      `json:"project_id"`
	Interfaces  []types.InterfaceDefinition `json:"interfaces"`
	LastScanned time.Time                   `json:"last_scanned"`
	ScanResults ScanResults                 `json:"scan_results"`
}

// ScanResults holds statistics about interface scanning
type ScanResults struct {
	FilesScanned    int           `json:"files_scanned"`
	InterfacesFound int           `json:"interfaces_found"`
	ScanDuration    time.Duration `json:"scan_duration"`
	Errors          []string      `json:"errors,omitempty"`
}

// MockGenerationJob represents a mock generation job
type MockGenerationJob struct {
	ID          string                        `json:"id"`
	ProjectID   string                        `json:"project_id"`
	Request     types.MockGenerationRequest   `json:"request"`
	Result      *types.MockGenerationResult   `json:"result,omitempty"`
	Status      JobStatus                     `json:"status"`
	CreatedAt   time.Time                     `json:"created_at"`
	StartedAt   *time.Time                    `json:"started_at,omitempty"`
	CompletedAt *time.Time                    `json:"completed_at,omitempty"`
}

// JobStatus represents the status of a generation job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusRunning    JobStatus = "running"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

// ProjectManager manages mockery projects
type ProjectManager struct {
	projects map[string]*MockeryProject
	mocks    map[string]*GeneratedMock
	jobs     map[string]*MockGenerationJob
}

// NewProjectManager creates a new project manager
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		projects: make(map[string]*MockeryProject),
		mocks:    make(map[string]*GeneratedMock),
		jobs:     make(map[string]*MockGenerationJob),
	}
}

// GetProject retrieves a project by ID
func (pm *ProjectManager) GetProject(id string) (*MockeryProject, bool) {
	project, exists := pm.projects[id]
	return project, exists
}

// CreateProject creates a new project
func (pm *ProjectManager) CreateProject(name, path string) *MockeryProject {
	project := &MockeryProject{
		ID:        generateID(),
		Name:      name,
		Path:      path,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	pm.projects[project.ID] = project
	return project
}

// AddGeneratedMock records a generated mock
func (pm *ProjectManager) AddGeneratedMock(mock *GeneratedMock) {
	pm.mocks[mock.ID] = mock
}

// GetGeneratedMocks returns all mocks for a project
func (pm *ProjectManager) GetGeneratedMocks(projectID string) []*GeneratedMock {
	var mocks []*GeneratedMock
	for _, mock := range pm.mocks {
		if mock.ProjectID == projectID {
			mocks = append(mocks, mock)
		}
	}
	return mocks
}

// CreateJob creates a new mock generation job
func (pm *ProjectManager) CreateJob(projectID string, request types.MockGenerationRequest) *MockGenerationJob {
	job := &MockGenerationJob{
		ID:        generateID(),
		ProjectID: projectID,
		Request:   request,
		Status:    JobStatusPending,
		CreatedAt: time.Now(),
	}
	pm.jobs[job.ID] = job
	return job
}

// GetJob retrieves a job by ID
func (pm *ProjectManager) GetJob(id string) (*MockGenerationJob, bool) {
	job, exists := pm.jobs[id]
	return job, exists
}

// UpdateJobStatus updates the status of a job
func (pm *ProjectManager) UpdateJobStatus(jobID string, status JobStatus) {
	if job, exists := pm.jobs[jobID]; exists {
		job.Status = status
		now := time.Now()
		
		switch status {
		case JobStatusRunning:
			job.StartedAt = &now
		case JobStatusCompleted, JobStatusFailed, JobStatusCancelled:
			job.CompletedAt = &now
		}
	}
}

// Helper function to generate unique IDs
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}