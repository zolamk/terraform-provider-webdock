package api

import (
	webdock "github.com/webdock-io/go-sdk"
)

// ClientInterface defines the methods used from webdock-io/go-sdk.
// This interface allows generating mocks for testing the Terraform provider.
type ClientInterface interface {
	CreatePublicKey(opts webdock.CreatePublicKeyOptions) (webdock.PublicKey, error)
	DeletePublicKey(options webdock.DeletePublicOptions) error
	ListAccountPublicKeys(options webdock.ListAccountPublicKeysOptions) ([]webdock.PublicKey, error)

	ListEvents(options webdock.ListEventsOptions) (webdock.ListEventsResponse, error)
	ListOSImages(options webdock.ListOSImagesOptions) ([]webdock.Image, error)
	ListLocations(options webdock.ListLocationsOptions) ([]webdock.Location, error)
	ListPossibleProfilesInLocation(opts webdock.ListPossibleProfilesInLocationOptions) ([]webdock.Profile, error)

	ListServer(options webdock.ListServerOptions) (webdock.ListServers, error)
	CreateServerFromImage(ops webdock.CreateServerFromImageOptions) (webdock.CreatedServer, error)
	DeleteServerBySlug(options webdock.DeleteServerBySlugOptions) error
	GetServerBySlug(options webdock.GetServerBySlugOptions) (webdock.Server, error)
	UpdateServer(opts webdock.UpdateServerOptions) (webdock.Server, error)
	ReinstallServer(options webdock.ReinstallServerOptions) (string, error)
	ResizeServer(options webdock.ResizeServersOptions) (string, error)
	DryRunResizeServer(options webdock.DryRunResizeServerOptions) (webdock.ResizeDryRunResponse, error)

	ListServerShellUser(options webdock.ListServerShellUserOptions) ([]webdock.ShellUser, error)
	CreateServerShellUser(opts webdock.CreateServerShellUserOptions) (webdock.CreatedShellUser, error)
	DeleteShellUser(options webdock.DeleteShellUserOptions) (webdock.DeleteShellUserResponse, error)
	UpdateServerShellUser(opts webdock.UpdateServerShellUserOptions) (webdock.ShellUser, error)
}
