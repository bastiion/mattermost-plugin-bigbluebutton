import MeetingHint from "./meeting_hint";

const mapStateToProps = state => {
  return {};
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({}, dispatch);
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(MeetingHint);
