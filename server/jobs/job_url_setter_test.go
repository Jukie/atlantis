package jobs_test

import (
	"errors"
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/runatlantis/atlantis/server/events/command"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/jobs"
	"github.com/runatlantis/atlantis/server/jobs/mocks"
	"github.com/runatlantis/atlantis/server/jobs/mocks/matchers"
	. "github.com/runatlantis/atlantis/testing"
	"github.com/stretchr/testify/assert"
)

func TestJobURLSetter(t *testing.T) {
	ctx := createTestProjectCmdContext(t)

	t.Run("update project status with project jobs url", func(t *testing.T) {
		RegisterMockTestingT(t)
		projectStatusUpdater := mocks.NewMockProjectStatusUpdater()
		projectJobURLGenerator := mocks.NewMockProjectJobURLGenerator()
		url := "url-to-project-jobs"
		jobURLSetter := jobs.NewJobURLSetter(projectJobURLGenerator, projectStatusUpdater)
		result := &command.ProjectResult{}

		When(projectJobURLGenerator.GenerateProjectJobURL(matchers.EqCommandProjectContext(ctx))).ThenReturn(url, nil)
		When(projectStatusUpdater.UpdateProject(ctx, command.Plan, models.PendingCommitStatus, url, nil)).ThenReturn(nil)
		err := jobURLSetter.SetJobURLWithStatus(ctx, command.Plan, models.PendingCommitStatus, result)
		Ok(t, err)

		projectStatusUpdater.VerifyWasCalledOnce().UpdateProject(ctx, command.Plan, models.PendingCommitStatus, "url-to-project-jobs", result)
	})

	t.Run("update project status with project jobs url error", func(t *testing.T) {
		RegisterMockTestingT(t)
		projectStatusUpdater := mocks.NewMockProjectStatusUpdater()
		projectJobURLGenerator := mocks.NewMockProjectJobURLGenerator()
		jobURLSetter := jobs.NewJobURLSetter(projectJobURLGenerator, projectStatusUpdater)

		When(projectJobURLGenerator.GenerateProjectJobURL(matchers.EqCommandProjectContext(ctx))).ThenReturn("url-to-project-jobs", errors.New("some error"))
		err := jobURLSetter.SetJobURLWithStatus(ctx, command.Plan, models.PendingCommitStatus, nil)
		assert.Error(t, err)
	})
}
