package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func TestBackgroundCheckWorkflowStandard(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := activities.Activities{SMTPStub: true, HTTPStub: true}

	env.RegisterWorkflow(workflows.Accept)
	env.RegisterActivity(a.SendAcceptEmail)
	env.RegisterWorkflow(workflows.SSNTrace)
	env.RegisterActivity(a.SSNTrace)
	env.RegisterWorkflow(workflows.FederalCriminalSearch)
	env.RegisterActivity(a.FederalCriminalSearch)
	env.RegisterActivity(a.SendReportEmail)

	details := types.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.SetOnChildWorkflowStartedListener(func(workflowInfo *workflow.Info, ctx workflow.Context, args converter.EncodedValues) {
		if workflowInfo.WorkflowExecution.ID == workflows.AcceptWorkflowID("john@example.com") {
			env.SignalWorkflowByID(
				workflows.AcceptWorkflowID("john@example.com"),
				workflows.AcceptSubmissionSignal,
				types.AcceptSubmissionSignal{Accepted: true, CandidateDetails: details},
			)
		}
	})

	env.ExecuteWorkflow(workflows.BackgroundCheck, &types.BackgroundCheckWorkflowInput{Email: "john@example.com", Package: "standard"})

	var result types.BackgroundCheckWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
	assert.Empty(t, result.CheckErrors)
}

func TestBackgroundCheckWorkflowFull(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := activities.Activities{SMTPStub: true, HTTPStub: true}

	env.RegisterWorkflow(workflows.Accept)
	env.RegisterActivity(a.SendAcceptEmail)
	env.RegisterWorkflow(workflows.SSNTrace)
	env.RegisterActivity(a.SSNTrace)
	env.RegisterWorkflow(workflows.FederalCriminalSearch)
	env.RegisterActivity(a.FederalCriminalSearch)
	env.RegisterWorkflow(workflows.StateCriminalSearch)
	env.RegisterActivity(a.StateCriminalSearch)
	env.RegisterWorkflow(workflows.MotorVehicleIncidentSearch)
	env.RegisterWorkflow(workflows.EmploymentVerification)
	env.RegisterActivity(a.SendEmploymentVerificationRequestEmail)
	env.RegisterActivity(a.SendReportEmail)

	details := types.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.SetOnChildWorkflowStartedListener(func(workflowInfo *workflow.Info, ctx workflow.Context, args converter.EncodedValues) {
		if workflowInfo.WorkflowExecution.ID == workflows.AcceptWorkflowID("john@example.com") {
			env.SignalWorkflowByID(
				workflows.AcceptWorkflowID("john@example.com"),
				workflows.AcceptSubmissionSignal,
				types.AcceptSubmissionSignal{Accepted: true, CandidateDetails: details},
			)
		}
	})

	env.ExecuteWorkflow(workflows.BackgroundCheck, &types.BackgroundCheckWorkflowInput{Email: "john@example.com", Package: "full"})

	var result types.BackgroundCheckWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
	assert.Empty(t, result.CheckErrors)
}
