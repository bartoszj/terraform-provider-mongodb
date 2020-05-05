package types

// Requests
type CreateUserRequest struct {
	User     string  `bson:"createUser"`
	Password string  `bson:"pwd"`
	Roles    []*Role `bson:"roles"`
}

type UserInfoRequest struct {
	User string `bson:"usersInfo"`
}

type UpdateUserRequest struct {
	User     string  `bson:"updateUser"`
	Password string  `bson:"pwd"`
	Roles    []*Role `bson:"roles"`
}

type DropUserRequest struct {
	User string `bson:"dropUser"`
}
