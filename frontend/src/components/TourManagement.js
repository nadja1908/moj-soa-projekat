import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Card, Button, Table, Badge, Alert, Spinner, Modal, Row, Col } from 'react-bootstrap';
import { tourApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import CreateTourModal from './CreateTourModal';
import KeyPointsMap from './KeyPointsMap';

const TourManagement = () => {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showMapModal, setShowMapModal] = useState(false);
  const [selectedTour, setSelectedTour] = useState(null);
  const [showPublishModal, setShowPublishModal] = useState(false);
  const [publishPrice, setPublishPrice] = useState('');
  const [publishLoading, setPublishLoading] = useState(false);
  const [tourDurations, setTourDurations] = useState({});
  const [showAllTours, setShowAllTours] = useState(false);

  // Debug logging za state promene
  useEffect(() => {
    console.log('📍 showMapModal changed:', showMapModal);
  }, [showMapModal]);

  useEffect(() => {
    console.log('📍 selectedTour changed:', selectedTour);
  }, [selectedTour]);
  const { user } = useAuth();

  const fetchMyTours = useCallback(async () => {
    try {
      setLoading(true);
      const response = await tourApi.get('/my');
      const tours = response.data.tours || [];
      setTours(tours);
      
      // Fetch durations for each tour
      const durationsMap = {};
      for (const tour of tours) {
        try {
          const durationResponse = await fetch(`http://localhost:8004/api/durations/tour/${tour.id}`);
          if (durationResponse.ok) {
            const durationData = await durationResponse.json();
            durationsMap[tour.id] = durationData.durations || [];
          }
        } catch (err) {
          console.warn(`Failed to fetch durations for tour ${tour.id}:`, err);
        }
      }
      setTourDurations(durationsMap);
    } catch (error) {
      console.error('Error fetching tours:', error);
      setError('Greška pri učitavanju tura');
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchAllTours = useCallback(async () => {
    try {
      setLoading(true);
      const response = await tourApi.get('/published');
      const tours = response.data.tours || [];
      setTours(tours);
      
      // Fetch durations for each tour
      const durationsMap = {};
      for (const tour of tours) {
        try {
          const durationResponse = await fetch(`http://localhost:8004/api/durations/tour/${tour.id}`);
          if (durationResponse.ok) {
            const durationData = await durationResponse.json();
            durationsMap[tour.id] = durationData.durations || [];
          }
        } catch (err) {
          console.warn(`Failed to fetch durations for tour ${tour.id}:`, err);
        }
      }
      setTourDurations(durationsMap);
    } catch (error) {
      console.error('Error fetching all tours:', error);
      setError('Greška pri učitavanju dostupnih tura');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    console.log('TourManagement useEffect - user:', user);
    console.log('User role:', user?.role);
    
    if (user && user.role === 'guide') {
      console.log(showAllTours ? 'Fetching all tours...' : 'Fetching my tours...');
      if (showAllTours) {
        fetchAllTours();
      } else {
        fetchMyTours();
      }
    } else {
      console.log('User not authenticated or not a guide');
    }
  }, [user, fetchMyTours, fetchAllTours, showAllTours]);

  const handleCreateTour = () => {
    setShowCreateModal(true);
  };

  const handleTourCreated = (newTour) => {
    setTours(prevTours => [newTour, ...prevTours]);
    setShowCreateModal(false);
  };

  const handleKeyPointsClick = (tour) => {
    console.log('🗺️ handleKeyPointsClick called with tour:', tour);
    setSelectedTour(tour);
    
    // Koristimo setTimeout da osiguramo da se selectedTour setter završi
    setTimeout(() => {
      setShowMapModal(true);
      console.log('🗺️ Modal should be opening...');
    }, 50);
  };

  const handleMapModalClose = () => {
    console.log('🔴 Closing map modal...');
    setShowMapModal(false);
    // NE čisti selectedTour da se spreči unmounting KeyPointsMap komponente
    console.log('🔴 Map modal closed - preserving selectedTour for component stability');
  };

  const handlePublishClick = (tour) => {
    console.log('📤 Publish clicked for tour:', tour.id);
    setSelectedTour(tour);
    setPublishPrice(tour.price || '');
    setShowPublishModal(true);
  };

  const handlePublishTour = async () => {
    if (!selectedTour || !publishPrice) {
      setError('Molimo unesite cenu ture');
      return;
    }

    try {
      setPublishLoading(true);
      setError('');
      
      const response = await tourApi.post(`/${selectedTour?.id}/publish`, {
        price: parseFloat(publishPrice)
      });

      if (response.data.success) {
        // Refresh lista tura
        if (showAllTours) {
          await fetchAllTours();
        } else {
          await fetchMyTours();
        }
        setShowPublishModal(false);
        setPublishPrice('');
        setSelectedTour(null);
      }
    } catch (error) {
      console.error('Error publishing tour:', error);
      setError(error.response?.data?.error || 'Greška pri objavljivanju ture');
    } finally {
      setPublishLoading(false);
    }
  };

  const handleArchiveTour = async (tour) => {
    if (!window.confirm(`Da li ste sigurni da želite da arhivirate turu "${tour.name}"?`)) {
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      const response = await tourApi.post(`/${tour.id}/archive`);

      if (response.data.success) {
        // Refresh lista tura
        if (showAllTours) {
          await fetchAllTours();
        } else {
          await fetchMyTours();
        }
      }
    } catch (error) {
      console.error('Error archiving tour:', error);
      setError(error.response?.data?.error || 'Greška pri arhiviranju ture');
    } finally {
      setLoading(false);
    }
  };

  const handleReactivateTour = async (tour) => {
    if (!window.confirm(`Da li ste sigurni da želite da reaktivirate turu "${tour.name}"?`)) {
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      const response = await tourApi.post(`/${tour.id}/reactivate`);

      if (response.data.success) {
        // Refresh lista tura
        if (showAllTours) {
          await fetchAllTours();
        } else {
          await fetchMyTours();
        }
      }
    } catch (error) {
      console.error('Error reactivating tour:', error);
      setError(error.response?.data?.error || 'Greška pri reaktivaciji ture');
    } finally {
      setLoading(false);
    }
  };

  // Stabiliziraj selectedTour podatke
  const selectedTourData = useMemo(() => ({
    id: selectedTour?.id,
    name: selectedTour?.name
  }), [selectedTour?.id, selectedTour?.name]);

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

  const renderDuration = (tour) => {
    const durations = tourDurations[tour.id];
    
    if (!durations || durations.length === 0) {
      return <span className="text-muted">-</span>;
    }
    
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
      <div className="d-flex flex-wrap gap-1">
        {durations.slice(0, 2).map((duration, index) => (
          <Badge key={index} bg="secondary" className="text-white" style={{ fontSize: '0.65rem' }}>
            {getTransportIcon(duration.transportType)} {formatDuration(duration.durationMinutes)}
          </Badge>
        ))}
        {durations.length > 2 && (
          <Badge bg="light" className="text-dark" style={{ fontSize: '0.65rem' }}>
            +{durations.length - 2}
          </Badge>
        )}
      </div>
    );
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    
    try {
      const date = new Date(dateString);
      // Proverava da li je datum valjan
      if (isNaN(date.getTime())) {
        return 'N/A';
      }
      return date.toLocaleDateString('sr-RS', {
        day: '2-digit',
        month: '2-digit', 
        year: 'numeric'
      });
    } catch (error) {
      console.warn('Error formatting date:', dateString, error);
      return 'N/A';
    }
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
              <div className="d-flex align-items-center gap-3">
                <h4 className="mb-0">
                  {showAllTours ? 'Sve dostupne ture' : 'Moje ture'}
                </h4>
                <div className="btn-group" role="group">
                  <button
                    type="button"
                    className={`btn btn-sm ${!showAllTours ? 'btn-primary' : 'btn-outline-primary'}`}
                    onClick={() => setShowAllTours(false)}
                  >
                    🏠 Moje ture
                  </button>
                  <button
                    type="button"
                    className={`btn btn-sm ${showAllTours ? 'btn-primary' : 'btn-outline-primary'}`}
                    onClick={() => setShowAllTours(true)}
                  >
                    🌍 Sve ture
                  </button>
                </div>
              </div>
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
                      <th>Vreme</th>
                      <th>Kreirana</th>
                      <th>Objavljena</th>
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
                        <td>{renderDuration(tour)}</td>
                        <td>
                          {formatDate(tour.createdAt)}
                        </td>
                        <td>
                          {tour.status === 'PUBLISHED' || tour.status === 'ARCHIVED' 
                            ? formatDate(tour.publishedAt) 
                            : <span className="text-muted">-</span>
                          }
                        </td>
                        <td>
                          <div className="d-flex gap-1">
                            <Button
                              variant="outline-primary"
                              size="sm"
                              title="Izmeni turu"
                              disabled={tour.status === 'PUBLISHED' || tour.status === 'ARCHIVED'}
                            >
                              ✏️
                            </Button>
                            <Button
                              variant="outline-info" 
                              size="sm"
                              title="Ključne tačke"
                              onClick={() => handleKeyPointsClick(tour)}
                            >
                              🗺️
                            </Button>
                            
                            {/* Status-specific buttons */}
                            {tour.status === 'DRAFT' && (
                              <Button
                                variant="outline-success"
                                size="sm"
                                title="Objavi turu"
                                onClick={() => handlePublishClick(tour)}
                              >
                                📤
                              </Button>
                            )}
                            
                            {tour.status === 'PUBLISHED' && (
                              <Button
                                variant="outline-warning"
                                size="sm"
                                title="Arhiviraj turu"
                                onClick={() => handleArchiveTour(tour)}
                              >
                                📦
                              </Button>
                            )}
                            
                            {tour.status === 'ARCHIVED' && (
                              <Button
                                variant="outline-success"
                                size="sm"
                                title="Reaktiviraj turu"
                                onClick={() => handleReactivateTour(tour)}
                              >
                                🔄
                              </Button>
                            )}
                            
                            {/* Delete only for drafts */}
                            {tour.status === 'DRAFT' && (
                              <Button
                                variant="outline-danger"
                                size="sm"
                                title="Obriši turu"
                              >
                                🗑️
                              </Button>
                            )}
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

      {/* Modal za upravljanje ključnim tačkama - uvek render-ovan za očuvanje state-a */}
      {selectedTourData && selectedTourData.id && (
        <KeyPointsMap 
          show={showMapModal}
          onHide={handleMapModalClose}
          tourId={selectedTourData.id}
          tourName={selectedTourData.name}
          onTourUpdate={fetchMyTours}
        />
      )}

      {/* Modal za objavljivanje ture */}
      <Modal show={showPublishModal} onHide={() => !publishLoading && setShowPublishModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Objavi turu</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedTour && (
            <div>
              <h6>{selectedTour?.name}</h6>
              <p className="text-muted">{selectedTour?.description}</p>
              
              <div className="mb-3">
                <label htmlFor="tourPrice" className="form-label">Cena ture (€)</label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  className="form-control"
                  id="tourPrice"
                  value={publishPrice}
                  onChange={(e) => setPublishPrice(e.target.value)}
                  placeholder="Unesite cenu ture"
                  disabled={publishLoading}
                />
                <small className="text-muted">Trenutno ima {selectedTour?.keyPointsCount || 0} ključnih tačaka</small>
              </div>

              {error && (
                <Alert variant="danger" className="mb-3">
                  {error}
                </Alert>
              )}
            </div>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button 
            variant="secondary" 
            onClick={() => setShowPublishModal(false)}
            disabled={publishLoading}
          >
            Odustani
          </Button>
          <Button 
            variant="success" 
            onClick={handlePublishTour}
            disabled={publishLoading || !publishPrice}
          >
            {publishLoading ? (
              <>
                <Spinner size="sm" className="me-2" />
                Objavljivanje...
              </>
            ) : (
              'Objavi turu'
            )}
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default TourManagement;