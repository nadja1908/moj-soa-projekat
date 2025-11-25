import React, { useState, useEffect } from 'react';
import { Card, Button, Spinner, Alert, Badge } from 'react-bootstrap';
import { followerApi, usersApi } from '../services/api';
import { useAuth } from '../context/AuthContext';

const Community = () => {
  const [recommendations, setRecommendations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [followingUsers, setFollowingUsers] = useState(new Set());
  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      loadRecommendations();
    }
  }, [user]);

  const loadRecommendations = async () => {
    try {
      setLoading(true);
      setError('');

      // Učitaj korisnike koje pratim
      const followingResponse = await followerApi.get('/following');
      const followingSet = new Set(followingResponse.data);
      setFollowingUsers(followingSet);

      console.log('Following users:', Array.from(followingSet));

      // Pozovi backend da dobije preporuke
      const response = await followerApi.get('/recommendations');
      console.log('Recommendations response:', response.data);

      setRecommendations(response.data);
    } catch (error) {
      console.error('Error loading recommendations:', error);
      setError('Greška pri učitavanju preporuka');
    } finally {
      setLoading(false);
    }
  };

  const handleFollow = async (userId) => {
    try {
      await followerApi.post('/follow', { followingId: userId });
      await loadRecommendations();
      setError('');
    } catch (error) {
      console.error('Error following user:', error);
      setError('Greška pri praćenju korisnika');
    }
  };

  if (loading) {
    return (
      <div className="text-center py-5">
        <Spinner animation="border" variant="primary" />
        <p className="mt-2 text-muted">Učitavanje preporuka...</p>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="display-5">👥 Zajednica</h1>
        <p className="lead text-muted">
          Otkrijte nove ljude koje možete pratiti
        </p>
      </div>

      {error && (
        <Alert variant="danger" dismissible onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {recommendations.length === 0 ? (
        <Card className="text-center py-5 border-0 shadow-sm">
          <Card.Body>
            <div className="feature-icon">🔍</div>
            <h3>Nema preporuka</h3>
            <p className="text-muted">
              {followingUsers.size === 0
                ? 'Počnite da pratite korisnike da biste dobili preporuke.'
                : 'Trenutno nema novih preporuka za vas.'}
            </p>
          </Card.Body>
        </Card>
      ) : (
        <div className="row">
          {recommendations.map((rec) => (
            <div key={rec.userId} className="col-md-6 col-lg-4 mb-4">
              <Card className="h-100 border-0 shadow-sm">
                <Card.Body>
                  <div className="d-flex align-items-center mb-3">
                    <div
                      className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center"
                      style={{ width: '50px', height: '50px', fontSize: '1.5rem' }}
                    >
                      {rec.username?.charAt(0).toUpperCase() || '?'}
                    </div>
                    <div className="ms-3">
                      <h5 className="mb-0">{rec.username}</h5>
                      <small className="text-muted">{rec.email}</small>
                    </div>
                  </div>

                  {rec.firstName && rec.lastName && (
                    <p className="mb-2">
                      <strong>Ime:</strong> {rec.firstName} {rec.lastName}
                    </p>
                  )}

                  {rec.role && (
                    <Badge bg="info" className="mb-2">
                      {rec.role}
                    </Badge>
                  )}

                  {rec.commonFollowers > 0 && (
                    <p className="text-muted small mb-3">
                      <i className="bi bi-people"></i> {rec.commonFollowers}{' '}
                      {rec.commonFollowers === 1 ? 'zajednički pratilac' : 'zajednička pratioca'}
                    </p>
                  )}

                  <Button
                    variant="primary"
                    size="sm"
                    className="w-100"
                    onClick={() => handleFollow(rec.userId)}
                  >
                    Prati
                  </Button>
                </Card.Body>
              </Card>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Community;
