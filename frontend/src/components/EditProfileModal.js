import { useState, useEffect } from 'react';
import { Modal, Button, Form } from 'react-bootstrap';

const EditProfileModal = ({ show, onClose, profile, onSave }) => {
  const [editValues, setEditValues] = useState({
    firstName: '',
    lastName: '',
    biography: '',
    motto: '',
    profileImage: null,
  });

  useEffect(() => {
    if (profile) {
      setEditValues({
        firstName: profile.firstName,
        lastName: profile.lastName,
        biography: profile.biography,
        motto: profile.motto,
        profileImage: null,
      });
    }
  }, [profile]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditValues((prev) => ({ ...prev, [name]: value }));
  };

  const handleFileChange = (e) => {
    setEditValues((prev) => ({ ...prev, profileImage: e.target.files[0] }));
  };

  const handleSaveClick = () => {
    onSave(editValues);
  };

  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Izmeni profil</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Ime</Form.Label>
            <Form.Control
              name="firstName"
              value={editValues.firstName}
              onChange={handleInputChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Prezime</Form.Label>
            <Form.Control
              name="lastName"
              value={editValues.lastName}
              onChange={handleInputChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Biografija</Form.Label>
            <Form.Control
              as="textarea"
              rows={3}
              name="biography"
              value={editValues.biography}
              onChange={handleInputChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Moto (citat)</Form.Label>
            <Form.Control
              name="motto"
              value={editValues.motto}
              onChange={handleInputChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Slika</Form.Label>
            <Form.Control type="file" onChange={handleFileChange} />
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>Otkaži</Button>
        <Button variant="primary" onClick={handleSaveClick}>Sačuvaj</Button>
      </Modal.Footer>
    </Modal>
  );
};

export default EditProfileModal;
