package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func ValidateSSN(ctx workflow.Context, input types.ValidateSSNInput) (types.ValidateSSNResult, error) {
	return mocks.ValidateSSNResults[input], nil
}