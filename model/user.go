package model

type User struct {
	Base
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique_index"`
	PasswordHash string `gorm:"not null"`
	TokenHash    string `gorm:"not null;unique_index"`
}

func (u *User) Equals(other *User) bool {
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

func (u *User) Update(update *UpdateView, passwordHash string, tokenHash string) {
	u.Name = update.Name
	u.Email = update.Email
	u.PasswordHash = passwordHash
	u.TokenHash = tokenHash
}
