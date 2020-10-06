import React, {Fragment, useState, useEffect} from 'react';
import PropTypes from "prop-types";
import {Client4} from 'mattermost-redux/client';
import {Card} from 'semantic-ui-react'
import ProfileCard from "./profile_card.jsx";
import {Button} from "react-bootstrap";
import RemindProfileFiller from "./remind_profile_filler.jsx";
import Cookies from 'js-cookie';


//const PostUtils = window.PostUtils;


function PostTypeProfiles({teamId, actions: {getOwnUserProfile, getUserProfiles, openModal}}) {
  const [activeIndex, setActiveIndex] = useState(0);
  const [reminderOpen, setReminderOpen] = React.useState(false);
  const [channelMembers, setChannelMembers] = useState([]);

  const handleClick = (e, titleProps) => {
    const {index} = titleProps;
    const newIndex = activeIndex === index ? -1 : index;

    setActiveIndex(newIndex)
  };


  const checkOwnProfile = async () => {
    const ownProfile = await getOwnUserProfile();
    if (ownProfile) {
      const {livingPlace, pronoun} = ownProfile;
      if (livingPlace.length == 0 || pronoun.length == 0) {
        setReminderOpen(true);
      }
    }
  }

  const _getProfilesInTeam = async () => {
    //const ps = await getProfilesInTeam(teamId,0, 200);
    const ps = await Client4.getProfilesInTeam(teamId, 0, 1000);
    setChannelMembers(ps);

  }
  const setCSRFFromCookie = () => {
    const csrf = Cookies.get('MMCSRF');
    Client4.setCSRF(csrf);
  }

  useEffect(() => {
    try {
      //setCSRFFromCookie();
      checkOwnProfile();
      _getProfilesInTeam();

    } catch (e) {
    }
  }, [teamId]);


  const handleEditAccountSettings = () => {
    if (window.UserSettingsModal) {
      openModal({ModalId: 'user_settings', dialogType: window.UserSettingsModal});
    }
  }


  const createGroupChat = async () => {
    const userIds = channelMembers.map(m => m.id);
    const href = await Client4.createGroupChannel(userIds);
    console.log(href);
    window.location.href = href;
  }


  return (
    <Fragment>
      <Button onClick={createGroupChat}>create group chat</Button><br/><br/>
      <Button onClick={handleEditAccountSettings}>eigenes Profil bearbeiten</Button><br/><br/>
      <RemindProfileFiller open={reminderOpen} openAccountSettings={handleEditAccountSettings}/>
      <Card.Group>
        {channelMembers.map((member, index) => (
          <ProfileCard getUserProfiles={getUserProfiles}
                       key={`profile_${index}`}
                       member={member}
                       isActive={activeIndex === index}
                       cardIndex={index}
                       handleAccordionClick={handleClick}/>
        ))}
      </Card.Group>
    </Fragment>

  );


}


PostTypeProfiles.propTypes = {
  post: PropTypes.object.isRequired,
  state: PropTypes.object.isRequired,
  teamId: PropTypes.number.isRequired,
  actions: {

    getUserProfiles: PropTypes.func.isRequired,
    getOwnUserProfile: PropTypes.func.isRequired,
    openModal: PropTypes.func.isRequired,
  }
};


export default PostTypeProfiles
