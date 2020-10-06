import React, {useState, useEffect, Fragment} from 'react';
import SettingItemMin from "../setting_item_min";
import SettingItemMax from "../setting_item_min/setting_item_max.jsx";

function UserSettingsBpBProfile(props) {


  const profileFieldConfigs = [
    {
      field: 'age',
      label: 'Alter',
      extraInfo: 'Alter in Jahren',
      type: 'number',
      maxLength: 3,
      describe: 'sameAsValue'

    },
    {
      field: 'livingPlace',
      label: 'Wohnort',
      extraInfo: 'Wohnort oder Hauptaufenthaltsort',
      type: 'text',
      maxLength: 50,
      describe: 'sameAsValue'
    },
    {
      field: 'pronoun',
      label: 'Pronomen',
      extraInfo: 'Freitext zur Geschlechteridentität / gewünschtes Pronomen',
      type: 'text',
      maxLength: 25,
      describe: 'sameAsValue'
    },
    {
      field: 'magicPower',
      label: 'Zauberkraft',
      extraInfo: 'Welche Zauberkraft oder Superkraft besitzt du?',
      type: 'text',
      maxLength: 200,
      describe: 'sameAsValue'

    },
    {
      field: 'misc',
      label: 'Über mich',
      extraInfo: 'Was ihr über mich wissen solltet... Freitext',
      type: 'text'

    }
  ];
  const profileFields = profileFieldConfigs.map(({field}) => field); //['age', 'livingPlace', 'pronoun', 'magicPower', 'food', 'misc'];
  const fieldStates = {}

  for (let field of profileFields) {
    const [state, setState] = useState('');
    fieldStates[field] = {state, setState};
  }

  const updateFieldCallback = (fieldName) => {
    return (e) => {
      fieldStates[fieldName].setState(e.target.value);
    }
  }


  const _getUserProfile = async () => {
    const profile = await props.getUserProfile();
    if (profile) {
      for (let field of profileFields) {
        if (profile[field]) {
          fieldStates[field].setState(profile[field])
        }
      }
    }
  };

  useEffect(() => {
    try {
      _getUserProfile()
    } catch (e) {
      console.error('ERROR!', e)
    }
  }, []);

  const _resetActiveSection = async () => {
    await props.resetActiveSection();
    await _getUserProfile();
  };

  const submitProfileCallback = (fieldName) => {
    return async (e) => {

      const resp = await props.submitProfile(fieldName, fieldStates[fieldName].state);
      if (!resp || !resp.error) {
        await _resetActiveSection();
      }
    }
  };

  return (
    <div>
      {profileFieldConfigs.map(({field, label, extraInfo, type, maxLength, describe}) =>
        (<Fragment>
            <div className='divider-light'/>
            {props.activeSection === field + 'Section' ? (
              <SettingItemMax
                title={label}
                inputs={
                  <div
                    key={`${field}Setting`}
                    className='form-group'
                  >
                    <label className='col-sm-5 control-label'>{label}</label>
                    <div className='col-sm-7'>
                      <input
                        id={field}
                        autoFocus={true}
                        className='form-control'
                        type={type}
                        onChange={updateFieldCallback(field)}
                        value={fieldStates[field].state}
                        maxLength={maxLength}
                        autoCapitalize='off'
                        aria-label={label}
                      />
                    </div>
                  </div>
                }
                extraInfo={extraInfo}
                section={`${field}Section`}
                updateSection={_resetActiveSection}
                submit={submitProfileCallback(field)}
              />
            ) : (
              <SettingItemMin
                title={label}
                section={`${field}Section`}
                describe={
                  describe === 'sameAsValue' && fieldStates[field].state && fieldStates[field].state.length > 0
                    ? fieldStates[field].state
                    : extraInfo}
              />
            )}
          </Fragment>
        ))}
    </div>

  );
}

export default UserSettingsBpBProfile;
