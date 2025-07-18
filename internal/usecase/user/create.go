package user

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateUser struct {
	DB *gorm.DB
}

func (uc *CreateUser) Process(entity *domain.User) error {
	model := persistence.UserPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
