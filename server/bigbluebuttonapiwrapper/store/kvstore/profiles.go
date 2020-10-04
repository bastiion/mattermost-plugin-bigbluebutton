package kvstore

import (
	"encoding/json"
	"errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
)

type ProfilesStore struct {
	api plugin.API
}

const profilesPrefix = "bbb_profiles_"

func DecodeProfileFromByte(b []byte) *dataStructs.Profile {
	p := dataStructs.Profile{}
	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil
	}
	return &p
}

func EncodeProfileToByte(p *dataStructs.Profile) []byte {
	b, _ := json.Marshal(p)
	return b
}

func (s *ProfilesStore) Get(id string) (*dataStructs.Profile, error) {
	b, err := s.api.KVGet(profilesPrefix + id)
	if err != nil {
		return nil, err
	}

	profile := DecodeProfileFromByte(b)
	if profile == nil {
		return nil, errors.New("failed to decode profile")
	}

	return profile, nil
}

// Insert stores new a profiles in the KV Store.
func (s *ProfilesStore) Insert(profile *dataStructs.Profile) error {
	opt := model.PluginKVSetOptions{
		Atomic:   true,
		OldValue: nil,
	}
	ok, err := s.api.KVSetWithOptions(profilesPrefix+profile.ID, EncodeProfileToByte(profile), opt)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("profile already exists in database")
	}

	return nil
}

// Update updates an existing a profiles in the KV Store.
func (s *ProfilesStore) Update(prev *dataStructs.Profile, new *dataStructs.Profile) error {
	opt := model.PluginKVSetOptions{
		Atomic:   true,
		OldValue: EncodeProfileToByte(prev),
	}
	ok, err := s.api.KVSetWithOptions(profilesPrefix+prev.ID, EncodeProfileToByte(new), opt)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("profiles already exists in database")
	}

	return nil
}

// Delete deletes a profiles from the KV Store.
func (s *ProfilesStore) Delete(profile *dataStructs.Profile) error {
	if err := s.api.KVDelete(profilesPrefix + profile.ID); err != nil {
		return err
	}

	return nil
}
