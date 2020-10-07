import React, {Component} from 'react';
import SpeeddatingModal from "./speeddating_modal";
import {createSpeeddatingRooms} from "../../actions";
import {getCurrentTeamId} from "mattermost-redux/selectors/entities/teams";
import {getAllChannels, getCurrentChannelId} from "mattermost-redux/selectors/entities/channels";
import {getUserStatuses} from "mattermost-redux/selectors/entities/users";
const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;

function mapStateToProps(state) {
  return {
    teamId: getCurrentTeamId(state),
    channelId: getCurrentChannelId(state),
    channels: getAllChannels(state),
    statuses: getUserStatuses(state)

  };
}

function mapDispatchToProps(dispatch) {
  return {
    actions: bindActionCreators({
      createSpeeddatingRooms,
    }, dispatch)
  };
}


export default connect(
  mapStateToProps,
  mapDispatchToProps
)(SpeeddatingModal);
