package mapper

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"auth/internal/infrastructure/db/gorm/model"
)

func CredentialEntityToModel(e *entity.Credential) *model.Credential {
	if e == nil {
		return nil
	}

	return &model.Credential{
		ID:             e.ID().Value(),
		Username:       e.Username(),
		HashedPassword: e.HashedPassword().Value(),
		CreatedAt:      e.CreatedAt(),
		UpdatedAt:      e.UpdatedAt(),
	}
}

func CredentialModelToEntity(m *model.Credential) (*entity.Credential, error) {
	if m == nil {
		return nil, nil
	}

	credentialID, err := vo.NewCredentialID(m.ID)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := vo.NewHashedPassword(m.HashedPassword)
	if err != nil {
		return nil, err
	}

	return entity.RehydrateCredential(credentialID, m.Username, hashedPassword, m.CreatedAt, m.UpdatedAt)
}
