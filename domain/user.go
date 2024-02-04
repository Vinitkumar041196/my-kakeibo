package domain

type User struct {
	Email     string  `json:"email" gorm:"column:email;type:varchar(300);primaryKey"`
	FirstName string  `json:"firstName" gorm:"column:first_name;type:varchar(100)"`
	LastName  string  `json:"lastName" gorm:"column:last_name;type:varchar(100)"`
	Password  string  `json:"password" gorm:"column:password;type:varchar(200)"`
	CreatedAt *string `json:"-" gorm:"column:created_at;type:timestamp"`
	UpdatedAt *string `json:"-" gorm:"column:updated_at;type:timestamp"`
	IsActive  *bool   `json:"-" gorm:"column:is_active;default:true;type:boolean"`
}

func (User) TableName() string {
	return "users"
}

type SignUpRequest struct {
	User
	ConfirmPassword string `json:"confirmPassword"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserSuccessResponse struct {
	Message string `json:"message"`
}

type UserService interface {
	CreateUser(u *User) error
	UpdateUser(u *User) error
	GetUser(email string) (*User, error)
	DeleteUser(email string) error
}

type UserRepository interface {
	CreateUser(u *User) error
	UpdateUser(u *User) error
	GetUser(email string) (*User, error)
	SearchUserByKey(key, value string) (*User, error)
	DeleteUser(email string) error
}

func (u *User) MapMetaData(createdAt, updatedAt *string, isActive *bool) {
	u.UpdatedAt = updatedAt
	u.CreatedAt = createdAt
	u.IsActive = isActive
}
