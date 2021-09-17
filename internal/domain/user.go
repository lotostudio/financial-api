package domain

type User struct {
	// Unique id
	ID int64 `json:"id" binding:"required" db:"id" example:"1"`
	// Unique email
	Email string `json:"email" binding:"required,email" db:"email" example:"sirius@gmail.com"`
	// First name
	FirstName string `json:"firstName" binding:"required,alpha" db:"first_name" example:"Sirius"`
	// Last name
	LastName string `json:"lastName" binding:"required,alpha" db:"last_name" example:"Sam"`
	// Secret password
	Password string `json:"-" binding:"omitempty,alphanum,min=8" db:"password" example:"qweqweqwe"`
} // @name User

type UserToCreate struct {
	// Unique email
	Email string `json:"email" binding:"required,email" example:"sirius@gmail.com"`
	// First name
	FirstName string `json:"firstName" binding:"required,alpha" example:"Sirius"`
	// Last name
	LastName string `json:"lastName" binding:"required,alpha" example:"Sam"`
	// Secret password
	Password string `json:"password" binding:"required,alphanum,min=8" example:"qweqweqwe"`
} // @name UserToCreate

type UserToUpdate struct {
	// First name
	FirstName *string `json:"firstName" binding:"omitempty,alpha" example:"Sirius"`
	// Last name
	LastName *string `json:"lastName" binding:"omitempty,alpha" example:"Sam"`
	// Secret password
	Password *string `json:"password" binding:"omitempty,alphanum,min=8" example:"qweqweqwe"`
} // @name UserToUpdate

type UserToLogin struct {
	// Unique email
	Email string `json:"email" binding:"required,email" example:"sirius@gmail.com"`
	// Secret password
	Password string `json:"password" binding:"required,alphanum,min=8" example:"qweqweqwe"`
} // @name UserToLogin
