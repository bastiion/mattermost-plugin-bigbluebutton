import React, {Fragment} from 'react';
import PropTypes from 'prop-types';
import {Button, Header} from 'semantic-ui-react'
import {Modal, Image} from "react-bootstrap";

function RemindProfileFiller({open, openAccountSettings}) {
  const [closed, setClosed] = React.useState(false);

  const handleOpen = e => {
    e.preventDefault();
    setClosed(true);
  }

  const _openAccountSettings = e => {
    e.preventDefault();
    setClosed(true);
    openAccountSettings();
  }


  return (
    <Fragment>
      <Modal
        show={open && ! closed}
        size='tiny'
      >
        <Modal.Header>Dein Profil ist noch so leer...</Modal.Header>
        <Modal.Body image>
          <Image circle src='https://media.giphy.com/media/toWxnNj8hcQJS3oV2y/source.gif' wrapped/>
          <Header>Kannst du uns noch etwas über dich erzählen?</Header>
          <p>
            Dein Profil scheint noch nicht vollständig ausgefüllt worden zu sein. Um deine Sichtbarkeit zu erhöhen
            und die Vernetzung zu erleichtern bitten wir dich, dir einen Moment Zeit zu nehmen und für
            die Konferenz ein paar Stichpunke wie Alter, Wohnort, Zauberkräfte und Tätigkeiten auszufüllen.
          </p>
          <p>Es ist jederzeit über deine <a onClick={_openAccountSettings}>Kontoeinstellungen</a> änderbar.</p>
          <p>Jetzt Profil ausfüllen?</p>
        </Modal.Body>
        <Modal.Footer>
          <Button color='black' onClick={() => setClosed(true)}>
            Nöö
          </Button>
          <Button
            content="Ja, bitte!"
            labelPosition='right'
            icon='checkmark'
            onClick={_openAccountSettings}
            positive
          />
        </Modal.Footer>
      </Modal>
    </Fragment>
  )
}


RemindProfileFiller.propTypes = {
  openAccountSettings: PropTypes.func.isRequired
};

export default RemindProfileFiller;
