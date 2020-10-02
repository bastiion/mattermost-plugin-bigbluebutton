import {displayUsernameForUser} from "../../utils/user_utils";

const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;
import {getCurrentUserId, getProfilesInTeam, makeGetProfilesInChannel} from 'mattermost-redux/selectors/entities/users';
import PostTypeProfiles from "./post_type_profiles.jsx";
import {getChannelMembersInChannels} from "mattermost-redux/selectors/entities/channels";
import {getCurrentTeamId} from "mattermost-redux/selectors/entities/teams";


function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
  const user = state.entities.users.profiles[post.user_id] || {};
  let channelId = state.entities.channels.currentChannelId;
  const channel = state.entities.channels.channels[channelId]
  const userid = getCurrentUserId(state) || {};
  const teamid = getCurrentTeamId(state);
  const doGetChannelMembers = makeGetProfilesInChannel();
  const channelMembers = state.entities
    && state.entities.users
    && state.entities.users.profiles
    && Object.keys(state.entities.users.profiles).map(k => state.entities.users.profiles[k])
    || []; //getProfilesInTeam(state, teamid); //doGetChannelMembers(state, channelId, true);
  return {
    channelId,
    channel,
    state,
    ...ownProps,
    channelMembers,
    doGetChannelMembers,
    currentUserId: userid,
    creatorId: user.id,
    username: user.username,
    creatorName: displayUsernameForUser(user, state.entities.general.config),
  };
}

function mapDispatchToProps(dispatch) {
  return {
    actions: bindActionCreators({
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(PostTypeProfiles);
