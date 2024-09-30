package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
}

// model repository hanya bisa diakses repo
type repository struct {
	db *gorm.DB
}

//Bikin instance repository
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
