import React, {Fragment} from 'react';
import PropTypes from "prop-types";
import {Client4} from 'mattermost-redux/client';
import { Card, Icon, Image } from 'semantic-ui-react'



//const PostUtils = window.PostUtils;


export default class PostTypeProfiles extends React.PureComponent {

  static propTypes = {
    post: PropTypes.object.isRequired,
    state: PropTypes.object.isRequired,
    creatorId: PropTypes.string.isRequired,
    currentUserId: PropTypes.string.isRequired,
    channelId: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    channel: PropTypes.object.isRequired,
    creatorName: PropTypes.string.isRequired,
    channelMembers: PropTypes.array.isRequired
  };




  render() {
    return (

      <Card.Group>
        {this.props.channelMembers.map((member, index) => (
            <Card>
              <Image src={Client4.getProfilePictureUrl(member.id, member.last_picture_updat)} wrapped ui={false} />
              <Card.Content>
                <Card.Header>{member.nickname || member.nickname.length > 0 ? member.nickname : member.username}</Card.Header>
                <Card.Meta>
                  <span className='date'>Joined in 2015</span>
                </Card.Meta>
                <Card.Description>
                  {member.position}
                </Card.Description>
              </Card.Content>
              <Card.Content extra>
                <a>
                  <Icon name='user' />
                  22 Stars
                </a>
              </Card.Content>
            </Card>
        )) }
      </Card.Group>

    )
  }


}


