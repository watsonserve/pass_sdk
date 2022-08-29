package pass_sdk

func GetOpenId(ticket string, app string, secret string) (string, error) {
	// @TODO
	return ticket, nil
}

func GetUserData(openId string, app string, secret string) (*UserData, error) {
	// @TODO
	return &UserData{
		UserId: openId,
	}, nil
}
