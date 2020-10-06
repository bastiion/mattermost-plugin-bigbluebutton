import {submitProfile, updateActiveSection} from "../../actions";

const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;
import SettingItemMin from "./setting_item_min.jsx";

const mapStateToProps = state => {
  return {
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({updateActiveSection}, dispatch);
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(SettingItemMin);
