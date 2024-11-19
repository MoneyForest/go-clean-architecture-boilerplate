package dto

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type ListUsersQueryParams struct {
	Page     int    `query:"page"`
	PageSize int    `query:"pageSize"`
	SortBy   string `query:"sortBy"`
	Email    string `query:"email"`
}

type ListUsersResponse struct {
	Users     []UserResponse `json:"users"`
	Total     int            `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"pageSize"`
	TotalPage int            `json:"totalPage"`
}

type CreateUserRequestBody struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserResponse UserResponse

type GetUserParams struct {
	ID string `param:"id"`
}

type GetUserResponse UserResponse

type UpdateUserParams struct {
	ID string `param:"id"`
}

type UpdateUserRequestBody struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserResponse UserResponse

type DeleteUserParams struct {
	ID string `param:"id"`
}

type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
