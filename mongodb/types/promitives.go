package types

// Types
type UserInfo struct {
	Id       *string `bson:"_id"`
	User     *string `bson:"user"`
	Database *string `bson:"db"`
	Roles    []Role  `bson:"roles"`
}

type Role struct {
	Role     string `bson:"role"`
	Database string `bson:"db"`
}
