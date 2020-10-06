import {combineReducers} from 'redux';
import {
  STATUS_CHANGE,
  OPEN_ROOT_MODAL,
  CLOSE_ROOT_MODAL,
  UPDATE_ACTIVE_SECTION,
  RESET_ACTIVE_SECTION
} from './action_types';

const enabled = (state = false, action) => {
    switch (action.type) {
    case STATUS_CHANGE:
        return action.data;

    default:
        return state;
    }
};

const rootModalVisible = (state = false, action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return true;
    case CLOSE_ROOT_MODAL:
        return false;
    default:
        return state;
    }
};


const activeSection = (state = '', action) => {
  switch (action.type) {
    case UPDATE_ACTIVE_SECTION: {
      return action.data;
    }
    case RESET_ACTIVE_SECTION: {
      return '';
    }
    default:
      return state;
  }
}


const userProfiles = (state = {}, action) => {
  switch (action.type) {
    case 't': {

    }
    default:
      return state;
  }
}

export default combineReducers({
    enabled,
    rootModalVisible,
  activeSection
});
