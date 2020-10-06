import React, {useState, useEffect} from 'react';
import PropTypes from 'prop-types';
import {Card, Icon, Image} from "semantic-ui-react";
import {Client4} from "mattermost-redux/client";

ProfileCard.propTypes = {
  member: PropTypes.object.isRequired,
  handleAccordionClick: PropTypes.func.isRequired,
  getUserProfiles: PropTypes.func.isRequired,
  cardIndex: PropTypes.number.isRequired,
  isActive: PropTypes.bool
};

function ProfileCard({member: {id, last_picture_updat, nickname, username, position}, isActive, cardIndex, handleAccordionClick, getUserProfiles}) {


  const [{
    age,
    livingPlace,
    pronoun,
    magicPower,
    food,
    misc
  }, setUserProfile] = useState({
    age: '∞',
    livingPlace: 'Mutter Erde',
    pronoun: '*',
    magicPower: '',
    food:'',
    misc: ''
  });


  const _getUserProfiles = async (userids) => {
    const profiles = await getUserProfiles(userids);
    if (Array.isArray(profiles)) {
      for(let p of profiles) {
        if(p.user_id === id) {
          setUserProfile(p);
          return;
        }
      }
    }
  };

  useEffect(() => {
    try {

      _getUserProfiles([id]);
    } catch (e) {
    }
  }, [id])


  return (
    <Card>
      <Image src={Client4.getProfilePictureUrl(id, last_picture_updat)} wrapped ui={false} />
      <Card.Content>
        <Card.Header>{nickname || nickname.length > 0 ? nickname : username}</Card.Header>
        <Card.Meta>
          <span className='date'>{pronoun}</span>
        </Card.Meta>
        <Card.Description>
          <p>{position}</p>
          <p>{misc}</p>
        </Card.Description>
      </Card.Content>
      {magicPower && magicPower.length > 0 && (
        <Card.Content>
          <span className="strong" style={{fontWeight: 'bold'}}>Zauberkraft: </span>
          <span>{magicPower}</span>
        </Card.Content>
      )}
      {food && food.length > 0 && (
        <Card.Content>
          <span className="strong" style={{fontWeight: 'bold'}}>Lieblingsessen: </span>
          <span>{food}</span>
        </Card.Content>
      )}
      <Card.Content extra>
        <a>
          <Icon name='user'/>
          <span>{age && age.length > 0 ? age : '∞'} </span>
          <span>Jahre</span>
        </a>
        <br/>
        <a>
          <Icon name='world'/>
          {livingPlace && livingPlace.length > 0 ? livingPlace : 'Planet Earth'}
        </a>
      </Card.Content>
    </Card>

  );
}

export default ProfileCard;
