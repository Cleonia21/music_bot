package user

import (
	"MusicBot/user/utils"
	"testing"
)

func Test_hostUser_join(t *testing.T) {
	var host hostUser

	host.init(utils.UserID{Username: "host"}, nil, nil)

	type sender struct {
		id     utils.UserID
		getCh  chan msgBetweenUsers
		sendCh chan<- msgBetweenUsers
	}
	senders := []sender{
		{
			id:    utils.UserID{Username: "sender1"},
			getCh: make(chan msgBetweenUsers),
		},
		{
			id:    utils.UserID{Username: "sender2"},
			getCh: make(chan msgBetweenUsers),
		},
		{
			id:    utils.UserID{Username: "sender3"},
			getCh: make(chan msgBetweenUsers),
		},
	}

	for i := 0; i < len(senders); i++ {
		senders[i].sendCh = host.join(senders[i].id, senders[i].getCh)
	}

	/*
		type fields struct {
			userFather     userFather
			pass           string
			playList       playList.PlayList
			audio          *Audio.Audio
			getFromUsersCh chan msgBetweenUsers
			usersMapMutex  sync.RWMutex
			sendToUsersChs map[utils.UserID]chan<- msgBetweenUsers
		}
		type args struct {
			id         utils.UserID
			senderInCh chan<- msgBetweenUsers
		}
		tests := []struct {
			name         string
			fields       fields
			args         args
			wantHostInCh chan<- msgBetweenUsers
		}{
			// TODO: Add test cases.
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := &hostUser{
					userFather:     tt.fields.userFather,
					pass:           tt.fields.pass,
					playList:       tt.fields.playList,
					audio:          tt.fields.audio,
					getFromUsersCh: tt.fields.getFromUsersCh,
					usersMapMutex:  tt.fields.usersMapMutex,
					sendToUsersChs: tt.fields.sendToUsersChs,
				}
				if gotHostInCh := h.join(tt.args.id, tt.args.senderInCh); !reflect.DeepEqual(gotHostInCh, tt.wantHostInCh) {
					t.Errorf("join() = %v, want %v", gotHostInCh, tt.wantHostInCh)
				}
			})
		}
	*/
}
