package store

import "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"

type Store interface {
	Profiles() ProfilesStore
	System() SystemStore
}

type ProfilesStore interface {
	Get(id string) (*dataStructs.Profile, error)
	Insert(*dataStructs.Profile) error
	Update(prev *dataStructs.Profile, new *dataStructs.Profile) error
	Delete(*dataStructs.Profile) error
}
type SystemStore interface {
	GetVersion() (string, error)
	SaveVersion(version string) error
}
