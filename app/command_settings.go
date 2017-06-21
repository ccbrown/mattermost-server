// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/platform/model"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type SettingsProvider struct {
}

const (
	CMD_SETTINGS = "settings"
)

func init() {
	RegisterCommandProvider(&SettingsProvider{})
}

func (me *SettingsProvider) GetTrigger() string {
	return CMD_SETTINGS
}

func (me *SettingsProvider) GetCommand(T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_SETTINGS,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_settings.desc"),
		AutoCompleteHint: "",
		DisplayName:      T("api.command_settings.name"),
	}
}

func (me *SettingsProvider) DoCommand(args *model.CommandArgs, message string) *model.CommandResponse {
	// This command is handled client-side and shouldn't hit the server.
	return &model.CommandResponse{
		Text:         args.T("api.command_settings.unsupported.app_error"),
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	}
}