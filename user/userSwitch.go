package user

import "github.com/mymmrac/telego"

func (a *Admin) unregToHost(unreg *unregUser, host *hostUser) users {
	host.init(
		a.tg,
		a.logger,
		unreg.id,
	)
	return host
}

func (a *Admin) unregToSending(unreg *unregUser, sending *sendingUser) users {
	host, ok := a.searchUser(telego.ChatID{Username: unreg.url}).(*hostUser)
	if ok && host.validatePass(unreg.pass) {
		sending.init(
			a.tg,
			a.logger,
			unreg.id,
			host,
			a.audio,
		)
		host.join(sending)
		return sending
	} else {
		unreg.notValidate()
		return unreg
	}
}

func (a *Admin) hostToUnreg(host *hostUser, unreg *unregUser) users {
	unreg.Init(
		a.tg,
		a.logger,
		host.id,
	)
	return unreg
}

func (a *Admin) sendingToUnreg(sending *sendingUser, unreg *unregUser) users {
	unreg.Init(
		a.tg,
		a.logger,
		sending.id,
	)
	return unreg
}

func (a *Admin) userSwitch(from, to users) users {
	var unreg *unregUser
	var sending *sendingUser
	var host *hostUser
	ok := false

	if unreg, ok = from.(*unregUser); ok {
		if host, ok = to.(*hostUser); ok {
			return a.unregToHost(unreg, host)
		} else if sending, ok = to.(*sendingUser); ok {
			return a.unregToSending(unreg, sending)
		}
	} else if unreg, ok = to.(*unregUser); ok {
		if host, ok = from.(*hostUser); ok {
			return a.hostToUnreg(host, unreg)
		} else if sending, ok = from.(*sendingUser); ok {
			return a.sendingToUnreg(sending, unreg)
		}
	}
	return nil
}

func (a *Admin) searchUser(id telego.ChatID) users {
	for _, user := range a.users {
		if id.ID == user.getID().ID ||
			id.Username == user.getID().Username {
			return user
		}
	}
	return nil
}
