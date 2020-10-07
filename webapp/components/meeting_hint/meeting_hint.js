import React from 'react';
const {OverlayTrigger, Modal, Overlay, FormControl, FormGroup, ControlLabel, Checkbox} = window.ReactBootstrap

function MeetingHint({open, closeHintModal}) {
  const handleClose = () => {

  }
  return (
    <Modal show={open} onHide={handleClose}>
      <Modal.Header closeButton={true} style={style.header}>Kennlern Runde</Modal.Header>
      <Modal.Body>

      </Modal.Body>
      <Modal.Footer>
        <button type='button' className='btn btn-default' onClick={this.handleCloseSpeeddatingModal}>
          Close

        </button>

      </Modal.Footer>
    </Modal>
  );
}

export default MeetingHint;
