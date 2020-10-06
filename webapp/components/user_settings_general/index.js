import {getActiveSection} from "../../selectors";

const {bindActionCreators} = window.Redux;
const {connect} = window.ReactRedux;
import UserSettingsBpBProfile from "./UserSettingsBpBProfile";
import {getUserProfile, resetActiveSection, submitProfile} from "../../actions";

const mapStateToProps = state => {
  return {
    activeSection: getActiveSection( state )
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({
    getUserProfile,
    resetActiveSection,
    submitProfile
  }, dispatch);
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(UserSettingsBpBProfile);
