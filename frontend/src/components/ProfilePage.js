import React, { useEffect, useState } from 'react';
import { Card, Spinner, Alert, Row, Col, Badge, Button } from 'react-bootstrap';
import { useAuth } from '../context/AuthContext';
import { usersApi } from '../services/api';
import EditProfileModal from './EditProfileModal'; // import modal

const ProfilePage = () => {
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showEditModal, setShowEditModal] = useState(false);

  useEffect(() => {
    if (!user) return;

    const fetchProfile = async () => {
      try {
        setLoading(true);
        const res = await usersApi.get('/profile');
        setProfile(res.data);
      } catch (err) {
        console.error('Error fetching profile:', err);
        setError(err.response?.data?.error || 'Greška pri učitavanju profila');
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [user]);

  const handleSaveProfile = async (editValues) => {
    try {
      const formData = new FormData();
      formData.append('firstName', editValues.firstName);
      formData.append('lastName', editValues.lastName);
      formData.append('biography', editValues.biography);
      formData.append('motto', editValues.motto);
      if (editValues.profileImage) {
        formData.append('profileImage', editValues.profileImage);
      }

      await usersApi.put('/profile', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      // Osveži local state
      setProfile((prev) => ({
        ...prev,
        firstName: editValues.firstName,
        lastName: editValues.lastName,
        biography: editValues.biography,
        motto: editValues.motto,
        profileImageUrl: editValues.profileImage
          ? URL.createObjectURL(editValues.profileImage)
          : prev.profileImageUrl,
      }));

      setShowEditModal(false);
    } catch (err) {
      console.error('Error updating profile:', err);
      alert('Greška pri čuvanju profila');
    }
  };

  if (loading) return <Spinner animation="border" role="status" />;

  if (error) return <Alert variant="danger">{error}</Alert>;
  if (!profile) return <Alert variant="warning">Profil nije pronađen</Alert>;

  return (
    <div className="container-fluid mt-4">
      <Row>
        <Col md={{ span: 8, offset: 2 }}>
          <Card>
            <Card.Header className="d-flex justify-content-between align-items-center">
              <h4 className="mb-0">Detalji profila</h4>
              <Button variant="primary" size="sm" onClick={() => setShowEditModal(true)}>Izmeni</Button>
            </Card.Header>
            <Card.Body>
              <Row className="align-items-center">
                <Col md={4} className="text-center mb-3 mb-md-0">
                  {profile.profileImageUrl ? (
                    <img
                      src={profile.profileImageUrl}
                      alt={`${profile.firstName} ${profile.lastName}`}
                      className="img-fluid rounded-circle"
                      style={{ width: '150px', height: '150px', objectFit: 'cover' }}
                    />
                  ) : (
                    <div className="bg-secondary rounded-circle" style={{ width: '150px', height: '150px' }} />
                  )}
                </Col>
                <Col md={8}>
                  <h2>{profile.firstName} {profile.lastName}</h2>
                  {profile.motto && <p className="fst-italic text-muted">"{profile.motto}"</p>}
                  {profile.biography && <p>{profile.biography}</p>}
                  {profile.role && <Badge bg="info" className="text-dark">{profile.role.charAt(0).toUpperCase() + profile.role.slice(1)}</Badge>}
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <EditProfileModal
        show={showEditModal}
        onClose={() => setShowEditModal(false)}
        profile={profile}
        onSave={handleSaveProfile}
      />
    </div>
  );
};

export default ProfilePage;
