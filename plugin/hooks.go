// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package plugin

import (
	"net/http"
)

// Hooks represents an object that handles events for a plugin. Methods are likely to be added over
// time, and plugins are not expected to implement all of them.
type Hooks interface {
	// OnActivate is invoked when the plugin is activated. Implementations will usually want to save
	// the api argument for later use. Loading configuration for the first time is also a commonly
	// done here.
	OnActivate(API) error

	// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to
	// use the API, and the plugin will be terminated shortly after this invocation.
	OnDeactivate() error

	// OnConfigurationChange is invoked when configuration changes may have been made.
	OnConfigurationChange() error

	// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for
	// the /plugins/{id} path will be routed to the plugin.
	//
	// The Mattermost-User-Id header will be present if (and only if) the request is by an
	// authenticated user.
	ServeHTTP(http.ResponseWriter, *http.Request)
}
