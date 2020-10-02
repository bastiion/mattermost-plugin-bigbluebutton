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

      <div>
        {this.props.channelMembers.map((member, index) => (
          <Fragment>
            <Card>
              <Image src={Client4.getProfilePictureUrl(member.id, member.last_picture_updat)} wrapped ui={false} />
              <Card.Content>
                <Card.Header>Matthew</Card.Header>
                <Card.Meta>
                  <span className='date'>Joined in 2015</span>
                </Card.Meta>
                <Card.Description>
                  Matthew is a musician living in Nashville.
                </Card.Description>
              </Card.Content>
              <Card.Content extra>
                <a>
                  <Icon name='user' />
                  22 Friends
                </a>
              </Card.Content>
            </Card>

          <div key={index} className="card"
               style={{
                 width: "18rem",
                 display: "none",
                 margin: "1em"
               }}>
            <img src={Client4.getProfilePictureUrl(member.id, member.last_picture_updat)}                 className="card-img-top" alt=""/>
            <div className="card-body">
              <h5 className="card-title">Card title</h5>
              <p className="card-text">Some quick example text to build on the card title and make up the bulk of the
                card's content.</p>
              <a href="#" className="btn btn-primary">Go somewhere</a>
            </div>
          </div>
          </Fragment>
        )) }
      </div>

    )
  }


}


