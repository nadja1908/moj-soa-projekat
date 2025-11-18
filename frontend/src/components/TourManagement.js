import React, { useState, useEffect } from 'react';
import { Card, Button, Table, Badge, Alert, Spinner, Modal, Row, Col } from 'react-bootstrap';
import { tourApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import CreateTourModal from './CreateTourModal';

const TourManagement = () => {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    if (user && user.role === 'guide') {
      fetchMyTours();
    }
  }, [user]);

  const fetchMyTours = async () => {
    try {
      setLoading(true);
      const response = await tourApi.get('/my');
      setTours(response.data.tours || []);
    } catch (error) {
      console.error('Error fetching tours:', error);
      setError('Greška pri učitavanju tura');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTour = () => {
    setShowCreateModal(true);
  };

  const handleTourCreated = (newTour) => {
    setTours(prevTours => [newTour, ...prevTours]);
    setShowCreateModal(false);
  };

  const getStatusBadge = (status) => {
    const variants = {
      DRAFT: 'secondary',
      PUBLISHED: 'success', 
      ARCHIVED: 'warning'
    };
    
    const labels = {
      DRAFT: 'Nacrt',
      PUBLISHED: 'Objavljena',
      ARCHIVED: 'Arhivirana'
    };

    return (
      <Badge bg={variants[status] || 'secondary'}>
        {labels[status] || status}
      </Badge>
    );
  };

  const getDifficultyBadge = (difficulty) => {
    const variants = {
      EASY: 'success',
      MODERATE: 'warning',
      HARD: 'danger'
    };
    
    const labels = {
      EASY: 'Laka',
      MODERATE: 'Umerena',
      HARD: 'Teška'
    };

    return (
      <Badge bg={variants[difficulty] || 'secondary'}>
        {labels[difficulty] || difficulty}
      </Badge>
    );
  };

  const parseTagsFromString = (tagsString) => {
    if (!tagsString) return [];
    
    try {
      // Ako je već niz
      if (Array.isArray(tagsString)) {
        return tagsString;
      }
      
      // Ako je JSON string
      const parsed = JSON.parse(tagsString);
      return Array.isArray(parsed) ? parsed : [];
    } catch (error) {
      // Ako nije validan JSON, podeli po zarez
      return tagsString.split(',').map(tag => tag.trim()).filter(tag => tag);
    }
  };

  const renderTags = (tagsString) => {
    const tags = parseTagsFromString(tagsString);
    
    if (!tags || tags.length === 0) {
      return <span className="text-muted">-</span>;
    }
    
    return (
      <div className="d-flex flex-wrap gap-1">
        {tags.slice(0, 3).map((tag, index) => (
          <Badge key={index} bg="info" className="text-dark" style={{ fontSize: '0.7rem' }}>
            {tag}
          </Badge>
        ))}
        {tags.length > 3 && (
          <Badge bg="light" className="text-dark" style={{ fontSize: '0.7rem' }}>
            +{tags.length - 3}
          </Badge>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ height: '400px' }}>
        <Spinner animation="border" role="status" variant="primary">
          <span className="visually-hidden">Loading...</span>
        </Spinner>
      </div>
    );
  }

  return (
    <div className="container-fluid mt-4">
      <Row>
        <Col>
          <Card>
            <Card.Header className="d-flex justify-content-between align-items-center">
              <h4 className="mb-0">Upravljanje turama</h4>
              <Button 
                variant="success" 
                onClick={handleCreateTour}
                size="sm"
              >
                <i className="fas fa-plus me-2"></i>
                Kreiraj novu turu
              </Button>
            </Card.Header>

            <Card.Body>
              {error && (
                <Alert variant="danger" dismissible onClose={() => setError('')}>
                  {error}
                </Alert>
              )}

              {tours.length === 0 ? (
                <div className="text-center py-4">
                  <i className="fas fa-route fa-3x text-muted mb-3"></i>
                  <h5 className="text-muted">Nema kreирanih tura</h5>
                  <p className="text-muted">Kliknite na "Kreiraj novu turu" da počnete.</p>
                </div>
              ) : (
                <Table responsive hover>
                  <thead className="table-light">
                    <tr>
                      <th>Naziv</th>
                      <th>Opis</th>
                      <th>Težina</th>
                      <th>Tagovi</th>
                      <th>Status</th>
                      <th>Cena (€)</th>
                      <th>Rastojanje (km)</th>
                      <th>Kreirana</th>
                      <th>Akcije</th>
                    </tr>
                  </thead>
                  <tbody>
                    {tours.map((tour) => (
                      <tr key={tour.id}>
                        <td>
                          <strong>{tour.name}</strong>
                        </td>
                        <td>
                          <div style={{ maxWidth: '200px' }}>
                            {tour.description && tour.description.length > 80 
                              ? `${tour.description.substring(0, 80)}...` 
                              : tour.description || 'Nema opisa'}
                          </div>
                        </td>
                        <td>{getDifficultyBadge(tour.difficulty)}</td>
                        <td>{renderTags(tour.tags)}</td>
                        <td>{getStatusBadge(tour.status)}</td>
                        <td>€{tour.price || '0.00'}</td>
                        <td>{tour.distanceKm || '0.0'} km</td>
                        <td>
                          {tour.createdAt 
                            ? new Date(tour.createdAt).toLocaleDateString('sr-RS')
                            : 'N/A'
                          }
                        </td>
                        <td>
                          <div className="d-flex gap-1">
                            <Button
                              variant="outline-primary"
                              size="sm"
                              title="Izmeni turu"
                            >
                              <i className="fas fa-edit"></i>
                            </Button>
                            <Button
                              variant="outline-info" 
                              size="sm"
                              title="Ključne tačke"
                            >
                              <i className="fas fa-map-marker-alt"></i>
                            </Button>
                            {tour.status === 'DRAFT' && (
                              <Button
                                variant="outline-success"
                                size="sm"
                                title="Objavi turu"
                              >
                                <i className="fas fa-upload"></i>
                              </Button>
                            )}
                            <Button
                              variant="outline-danger"
                              size="sm"
                              title="Obriši turu"
                            >
                              <i className="fas fa-trash"></i>
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>

      {/* Modal za kreiranje nove ture */}
      <CreateTourModal 
        show={showCreateModal}
        onHide={() => setShowCreateModal(false)}
        onTourCreated={handleTourCreated}
      />
    </div>
  );
};

export default TourManagement;