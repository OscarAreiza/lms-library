package usecase

import (
	"context"

	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/circulation"
	"github.com/OscarAreiza/lms-library/circulation-service/internal/domain/service"
)

// ReturnLoan implements HU-07's acceptance criteria — a thin wrapper around
// ReturnRegistrationService.RegisterReturn, which also triggers the HU-08
// suspension policy when the return is late.
type ReturnLoan struct {
	service *service.ReturnRegistrationService
}

func NewReturnLoan(svc *service.ReturnRegistrationService) *ReturnLoan {
	return &ReturnLoan{service: svc}
}

func (uc *ReturnLoan) Execute(ctx context.Context, loanID string) (*circulation.Loan, error) {
	return uc.service.RegisterReturn(ctx, loanID)
}
