import React, { useState, useEffect } from 'react';
import { Modal, Button, Card, Badge, Alert, Spinner } from 'react-bootstrap';
import { tourApi } from '../services/api';

const TourDetailsModal = ({ show, onHide, tourId }) => {
  const [tour, setTour] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (show && tourId) {
      fetchTourDetails();
    }
  }, [show, tourId]);

  const fetchTourDetails = async () => {
    try {
      setLoading(true);
      setError('');
      
      // Koristimo tourist-specific endpoint koji vraća samo osnovne info + prvu ključnu tačku
      const response = await fetch(`http://localhost:8004/api/tours/public/${tourId}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch tour details');
      }
      
      const data = await response.json();
      if (data.success) {
        setTour(data.tour);
      } else {
        setError('Tura nije dostupna');
      }
    } catch (error) {
      console.error('Error fetching tour details:', error);
      setError('Greška pri učitavanju detalja ture');
    } finally {
      setLoading(false);
    }
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
      const parsed = JSON.parse(tagsString);
      return Array.isArray(parsed) ? parsed : [];
    } catch (error) {
      return tagsString.split(',').map(tag => tag.trim()).filter(tag => tag);
    }
  };

  const renderTags = (tagsString) => {
    const tags = parseTagsFromString(tagsString);
    
    if (!tags || tags.length === 0) {
      return <span className="text-muted">Nema tagova</span>;
    }
    
    return (
      <div className="d-flex flex-wrap gap-1">
        {tags.map((tag, index) => (
          <Badge key={index} bg="info" className="text-dark">
            {tag}
          </Badge>
        ))}
      </div>
    );
  };

  const formatDuration = (minutes) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours > 0) {
      return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
    }
    return `${minutes}m`;
  };

  const getTransportIcon = (transportType) => {
    switch (transportType) {
      case 'WALKING': return '🚶';
      case 'CYCLING': return '🚴';
      case 'CAR': return '🚗';
      case 'BUS': return '🚌';
      default: return '🚶';
    }
  };

  return (
    <Modal show={show} onHide={onHide} size="lg">
      <Modal.Header closeButton>
        <Modal.Title>
          {tour ? `🗺️ ${tour.tour.name}` : 'Detalji ture'}
        </Modal.Title>
      </Modal.Header>
      
      <Modal.Body>
        {loading && (
          <div className="text-center py-4">
            <Spinner animation="border" role="status" variant="primary">
              <span className="visually-hidden">Loading...</span>
            </Spinner>
          </div>
        )}

        {error && (
          <Alert variant="danger">
            {error}
          </Alert>
        )}

        {tour && (
          <div>
            {/* Osnovne informacije */}
            <Card className="mb-3">
              <Card.Header>
                <h5 className="mb-0">📋 Osnovne informacije</h5>
              </Card.Header>
              <Card.Body>
                <div className="row">
                  <div className="col-md-6">
                    <p><strong>Naziv:</strong> {tour.tour.name}</p>
                    <p><strong>Težina:</strong> {getDifficultyBadge(tour.tour.difficulty)}</p>
                    <p><strong>Cena:</strong> <span className="text-success fw-bold">€{tour.tour.price}</span></p>
                  </div>
                  <div className="col-md-6">
                    <p><strong>Rastojanje:</strong> {tour.tour.distanceKm || '0.0'} km</p>
                    <p><strong>Status:</strong> <Badge bg="success">Objavljena</Badge></p>
                  </div>
                </div>
                
                <div className="mt-3">
                  <p><strong>Opis:</strong></p>
                  <p className="text-muted">{tour.tour.description || 'Nema opisa'}</p>
                </div>

                <div className="mt-3">
                  <p><strong>Tagovi:</strong></p>
                  {renderTags(tour.tour.tags)}
                </div>
              </Card.Body>
            </Card>

            {/* Vremena trajanja */}
            {tour.durations && tour.durations.length > 0 && (
              <Card className="mb-3">
                <Card.Header>
                  <h5 className="mb-0">⏱️ Vremena trajanja</h5>
                </Card.Header>
                <Card.Body>
                  <div className="row">
                    {tour.durations.map((duration, index) => (
                      <div key={index} className="col-md-6 mb-2">
                        <Badge bg="secondary" className="text-white me-2">
                          {getTransportIcon(duration.transportType)} {formatDuration(duration.durationMinutes)}
                        </Badge>
                        <small className="text-muted">{duration.transportType.toLowerCase()}</small>
                      </div>
                    ))}\n                  </div>
                </Card.Body>
              </Card>
            )}

            {/* Prva ključna tačka */}
            {tour.keyPoints && tour.keyPoints.length > 0 && (
              <Card className="mb-3">
                <Card.Header>
                  <h5 className="mb-0">📍 Početna tačka</h5>
                </Card.Header>
                <Card.Body>
                  {tour.keyPoints.map((keyPoint, index) => (
                    <div key={index} className="mb-3">
                      <h6>{keyPoint.name}</h6>
                      <p className="text-muted mb-2">{keyPoint.description}</p>
                      <small className="text-info">
                        📍 Koordinate: {keyPoint.latitude.toFixed(6)}, {keyPoint.longitude.toFixed(6)}
                      </small>
                      {keyPoint.imageUrl && (
                        <div className="mt-2">
                          <img 
                            src={keyPoint.imageUrl} 
                            alt={keyPoint.name}
                            className="img-fluid rounded"
                            style={{ maxHeight: '200px' }}
                          />
                        </div>
                      )}
                    </div>
                  ))}
                  
                  <Alert variant="info" className="mt-3">
                    <small>
                      <i className="fas fa-info-circle me-1"></i>
                      Prikazuje se samo početna tačka. Ostatak rute će biti otkriven tokom ture.
                    </small>
                  </Alert>
                </Card.Body>
              </Card>
            )}

            {/* Napomena za rezervaciju */}
            <Alert variant="success">
              <h6><i className="fas fa-ticket-alt me-2"></i>Zainteresovani ste za ovu turu?</h6>
              <p className="mb-0">
                Kliknite na "Rezerviši" dugme u listi tura da rezervišete svoje mesto!
              </p>
            </Alert>
          </div>
        )}
      </Modal.Body>

      <Modal.Footer>
        <Button variant="secondary" onClick={onHide}>
          Zatvori
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default TourDetailsModal;