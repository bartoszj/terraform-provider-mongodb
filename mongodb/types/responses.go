package types

// Responses
type Response struct {
	Ok bool `bson:"ok"`
}

type UsersInfoResponse struct {
	Response
	UserInfos []UserInfo `bson:"users"`
}
