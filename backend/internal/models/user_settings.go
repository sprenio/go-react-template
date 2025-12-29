package models

import (
	"time"
)

const (
	APP_OPTION_A AppOption2 = "OPT_A"
	APP_OPTION_B AppOption2 = "OPT_B"
	APP_RADIO_A  AppOption3 = "RADIO_A"
	APP_RADIO_B  AppOption3 = "RADIO_B"
)

type UserFlags struct {
	Flag1 bool `json:"flag_1"`
	Flag2 bool `json:"flag_2"`
	Flag3 bool `json:"flag_3"`
}
type AppFlags struct {
	FlagA bool `json:"flag_a"`
	FlagB bool `json:"flag_b"`
}

type AppOption2 string
type AppOption3 string

type UserSettings struct {
	Id        uint       `db:"id" json:"id"`
	UserId    uint       `db:"user_id" json:"user_id"`
	LangId    uint8      `db:"lang_id" json:"lang_id"`
	UserFlags uint64     `db:"user_flags" json:"user_flags"`
	AppFlags  uint64     `db:"app_flags" json:"app_flags"`
	AppOpt1   string     `db:"app_opt_1" json:"app_opt_1"`
	AppOpt2   AppOption2 `db:"app_opt_2" json:"app_opt_2"`
	AppOpt3   AppOption3 `db:"app_opt_3" json:"app_opt_3"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type UserSettingsData struct {
	AppFlags  *AppFlags   `json:"app_flags,omitempty"`
	AppOpt1   *string     `json:"app_opt_1,omitempty"`
	AppOpt2   *AppOption2 `json:"app_opt_2,omitempty"`
	AppOpt3   *AppOption3 `json:"app_opt_3,omitempty"`
	Language  *string     `json:"language,omitempty"`
	UserFlags *UserFlags  `json:"user_flags,omitempty"`
}


func (us *UserSettings) GetUserFlags() *UserFlags {
	return &UserFlags{
		Flag1: us.UserFlags&(1<<0) != 0,
		Flag2: us.UserFlags&(1<<1) != 0,
		Flag3: us.UserFlags&(1<<2) != 0,
	}
}

func (us *UserSettings) GetAppFlags() *AppFlags {
	return &AppFlags{
		FlagA: us.AppFlags&(1<<0) != 0,
		FlagB: us.AppFlags&(1<<1) != 0,
	}
}
func (us *UserSettings) SetUserFlags(flags UserFlags) {
	var userFlags uint64
	userFlags |= boolToUint64(flags.Flag1) << 0
	userFlags |= boolToUint64(flags.Flag2) << 1
	userFlags |= boolToUint64(flags.Flag3) << 2
	us.UserFlags = userFlags
}
func (us *UserSettings) SetAppFlags(flags AppFlags) {
	var appFlags uint64
	appFlags |= boolToUint64(flags.FlagA) << 0
	appFlags |= boolToUint64(flags.FlagB) << 1
	us.AppFlags = appFlags
}

func boolToUint64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}
