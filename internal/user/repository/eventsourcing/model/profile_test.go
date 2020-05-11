package model

import (
	user_model "github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
	"testing"
)

func TestProfileChanges(t *testing.T) {
	type args struct {
		existing *Profile
		new      *Profile
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all attributes changed",
			args: args{
				existing: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
				new:      &Profile{UserName: "UserNameChanged", FirstName: "FirstNameChanged", LastName: "LastNameChanged", NickName: "NickNameChanged", DisplayName: "DisplayNameChanged", PreferredLanguage: language.English, Gender: int32(user_model.GENDER_MALE)},
			},
			res: res{
				changesLen: 6,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
				new:      &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
			},
			res: res{
				changesLen: 0,
			},
		},
		{
			name: "username changed",
			args: args{
				existing: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
				new:      &Profile{UserName: "UserNameChanged", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
			},
			res: res{
				changesLen: 0,
			},
		},
		{
			name: "empty names",
			args: args{
				existing: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
				new:      &Profile{UserName: "UserName", FirstName: "", LastName: "", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GENDER_FEMALE)},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
