package domain

type User struct {
	// Unique id
	ID int64 `json:"id" db:"id" example:"1"`
	// Unique email
	Email string `json:"email" db:"email" example:"sirius@gmail.com"`
	// First name
	FirstName string `json:"firstName" db:"first_name" example:"Sirius"`
	// Last name
	LastName string `json:"lastName" db:"last_name" example:"Sam"`
	// Secret password
	Password string `json:"-" db:"password" example:"qweqweqwe"`
} // @name User

type UserToCreate struct {
	// Unique email
	Email string `json:"email" example:"sirius@gmail.com"`
	// First name
	FirstName string `json:"firstName" example:"Sirius"`
	// Last name
	LastName string `json:"lastName" example:"Sam"`
	// Secret password
	Password string `json:"password" example:"qweqweqwe"`
} // @name UserToCreate

type UserToUpdate struct {
	// First name
	FirstName *string `json:"firstName" example:"Sirius"`
	// Last name
	LastName *string `json:"lastName" example:"Sam"`
	// Secret password
	Password *string `json:"password" example:"qweqweqwe"`
} // @name UserToUpdate

type UserToLogin struct {
	// Unique email
	Email string `json:"email" example:"sirius@gmail.com"`
	// Secret password
	Password string `json:"password" example:"qweqweqwe"`
} // @name UserToLogin

type Tokens struct {
	// Token used for accessing operations and/or resources
	AccessToken string `json:"accessToken" example:"access token"`
} // @name Tokens
