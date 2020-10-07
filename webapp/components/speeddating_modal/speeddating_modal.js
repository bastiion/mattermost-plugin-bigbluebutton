import React, {useState, useEffect} from 'react';
import PropTypes from 'prop-types';
import {Client4} from "mattermost-redux/client";
import {getChannelByName} from "mattermost-redux/utils/channel_utils";

const {OverlayTrigger, Modal, Overlay, FormControl, FormGroup, ControlLabel, Checkbox} = window.ReactBootstrap

SpeeddatingModal.propTypes = {};

function SpeeddatingModal({teamId, channelId,channels, statuses, actions: {createSpeeddatingRooms}}) {

  const [inlcudedUsers, setIncludedUsers] = useState([])
  const [excludedUsers, setExcludedUsers] = useState([])
  const [excludeChannel, setExcludeChannel] = useState("orga-team");
  const [onlyOnline, setOnlyOnline] = useState(false);
  const [onlyChannelMembers, setOnlyChannelMembers] = useState(false);
  const [duration, setDuration] = useState(20);
  const [usersPerRoom, setUsersPerRoom] = useState(6);
  const [displayName, setDisplayName] = useState("Kennenlernen")


  const filterOffline = (users) => users.filter( u => (statuses[u.id] && statuses[u.id] !== "offline"))

  const _useProfilesInTeam = async () => {
    //const ps = await getProfilesInTeam(teamId,0, 200);
    const ps = await Client4.getProfilesInTeam(teamId, 0, 1000);
    setIncludedUsers(onlyOnline ? filterOffline(ps) : ps);

  }
  const _useMembersInChannels = async () => {
    //const ps = await getProfilesInTeam(teamId,0, 200);
    const ps = await Client4.getProfilesInChannel(channelId);
    setIncludedUsers(onlyOnline ? filterOffline(ps) : ps);

  }


  const _exlcudeMembersOfChannel = async () => {
    const c = await getChannelByName(channels, excludeChannel);
    if(c) {
      const ps = await Client4.getProfilesInChannel(c.id);
      setExcludedUsers(ps)
    } else {
      setExcludedUsers([])
    }
  }

  useEffect(() => {
    try {
      if(onlyChannelMembers) {
        _useMembersInChannels();
      } else {
        _useProfilesInTeam();
      }
    } catch (e) {
    }
  }, [teamId, onlyChannelMembers, onlyOnline]);

  useEffect(() => {
    try {
      _exlcudeMembersOfChannel();
    } catch (e) {
    }
  }, [channels, excludeChannel]);

  const _startSpeeddating = async () => {
    console.log("start");
    const resp = await createSpeeddatingRooms({
      user_ids: inlcudedUsers.map(m => m.id),
      excluded_user_ids: excludedUsers.map(m => m.id),
      users_per_room: usersPerRoom,
      room_display_name: displayName,
      duration
    });
    console.log(resp);
  };

  const _setUsersPerRoom = (e) => {
    const num = parseInt(e.target.value);
    if (!isNaN(num)) {
      setUsersPerRoom(num);
    }
  };
  const _setDuration = (e) => {
    const num = parseInt(e.target.value);
    if (!isNaN(num)) {
      setDuration(num);
    }
  };

  const _setDisplayName = (e) => setDisplayName(e.target.value);

  const _setExcludeChannel = (e) => setExcludeChannel(e.target.value);
  return (
    <Modal.Body>
      <div>
        <h3>Nutzer für Kennenlernen</h3>
        <p>
          {inlcudedUsers.map(m => m.username).join(", ")}
        </p>
        <h3>Ausgeschlossene Nutzer</h3>
        <p>
          {excludedUsers.map(m => m.username).join(", ")}
        </p>
        <form>
          <FormGroup>
            <Checkbox checked={onlyChannelMembers} onChange={() => setOnlyChannelMembers(!onlyChannelMembers)}> nur Nutzer_innen dieses Kanals?</Checkbox>
            <Checkbox checked={onlyOnline} onChange={() => setOnlyOnline(!onlyOnline)}> nur Online Nutzer_innen (away, online)?</Checkbox>
            <ControlLabel>
              Dauer des Kennlern-Meetings
            </ControlLabel>
            <FormControl
              type="number"
              placeholder="Dauer in Minuten"
              value={duration}
              onChange={_setDuration}
            />
          </FormGroup>
          <FormGroup>
            <ControlLabel>
              Wieviele Nutzer_innen pro Raum
            </ControlLabel>
            <FormControl
              type="number"
              placeholder="Anzahl Nutzer_innen pro Raum"
              value={usersPerRoom}
              onChange={_setUsersPerRoom}
            />
          </FormGroup>
          <FormGroup>
            <ControlLabel>
              Wie heißen die privaten Kanäle
            </ControlLabel>
            <FormControl
              type="text"
              placeholder="Anzeigename"
              value={displayName}
              onChange={_setDisplayName}
            />
          </FormGroup>
          <FormGroup>
            <ControlLabel>
              Nutzer_inen von folgendem Kanals ausschließen
            </ControlLabel>
            <FormControl
              type="text"
              placeholder="Kanalname"
              value={excludeChannel}
              onChange={_setExcludeChannel}
            />
          </FormGroup>
        </form>

        <button type='button' className='btn btn-primary pull-left' onClick={_startSpeeddating}>
          Starten
        </button>
      </div>
    </Modal.Body>
  );
}

export default SpeeddatingModal;
