package model

type User struct {
	Base
	Name         string
	Email        string `gorm:"not null;unique_index"`
	PasswordHash string `gorm:"not null"`
	TokenHash    string `gorm:"not null;unique_index"`
}

func NewUserFromRegister(register *RegisterForm, passwordHash string, tokenHash string) *User {
	return &User{
		Name:         register.Name,
		Email:        register.Email,
		PasswordHash: passwordHash,
		TokenHash:    tokenHash,
	}
}

func NewUserFromUpdate(update *UpdateForm, passwordHash string, tokenHash string) *User {
	return &User{
		Name:         update.Name,
		Email:        update.Email,
		PasswordHash: passwordHash,
		TokenHash:    tokenHash,
	}
}

func (u *User) DeepEqual(other *User) bool {
	if u.ID != other.ID {
		return false
	}

	if u.CreatedAt != other.CreatedAt {
		return false
	}

	if u.UpdatedAt != other.UpdatedAt {
		return false
	}

	if u.DeletedAt != other.DeletedAt {
		return false
	}

	if u.Name != other.Name {
		return false
	}

	if u.Email != other.Email {
		return false
	}
	return true
}
