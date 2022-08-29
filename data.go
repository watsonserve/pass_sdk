package pass_sdk

func (this *userService) getUserInfo(ticket string) (*UserData, error) {
    openId, err := GetOpenId(ticket, this.app, this.secret)
    if nil != err {
        return nil, err
    }
    userData, err := this.bao.Get(openId)
    if nil != err {
        userData, err = GetUserData(openId, this.app, this.secret)
        if nil == err {
            err = this.bao.Save(ticket, userData)
        }
    }
    return userData, err
}
