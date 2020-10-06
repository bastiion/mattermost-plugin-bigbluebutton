const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;
import PostTypeProfiles from "./post_type_profiles.jsx";
import {getCurrentTeamId} from "mattermost-redux/selectors/entities/teams";
import {getUserProfile, getUserProfiles, openModal} from "../../actions";


function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
  const teamId = getCurrentTeamId(state)
  return {
    teamId,
    post,
    state,
    ...ownProps,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    actions: bindActionCreators({
      getUserProfiles,
      getOwnUserProfile: getUserProfile,
      openModal
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(PostTypeProfiles);
