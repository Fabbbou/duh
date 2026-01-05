package usecase

import (
	"duh/internal/domain/service"
)

type AliasUsecase struct {
	aliasService *service.AliasService
}

func NewAliasUsecase(aliasService *service.AliasService) *AliasUsecase {
	return &AliasUsecase{
		aliasService: aliasService,
	}
}

func (a *AliasUsecase) SetAlias(aliasName, value string) error {
	// Application layer: Input validation can be added here
	if err := a.aliasService.ValidateAliasName(aliasName); err != nil {
		return err
	}

	// Delegate to domain service for business logic
	return a.aliasService.SetAlias(aliasName, value)
}

func (a *AliasUsecase) UnsetAlias(aliasName string) error {
	// Delegate to domain service for business logic
	return a.aliasService.UnsetAlias(aliasName)
}

func (a *AliasUsecase) ListAliases() (map[string]string, error) {
	// Delegate to domain service for business logic
	return a.aliasService.GetMergedAliases()
}
