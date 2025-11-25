import { useState, useEffect, useCallback } from 'react';
import { Card, Button, Badge, Alert, Spinner, Row, Col } from 'react-bootstrap';
import { tourApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import TourDetailsModal from './TourDetailsModal';

import TourReviewsModal from './TourReviewsModal';


import { useCart } from './CartContext';
import { purchaseApi } from '../services/api'


const TouristTours = () => {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [tourDurations, setTourDurations] = useState({});
  const [showDetailsModal, setShowDetailsModal] = useState(false);
  const [selectedTourId, setSelectedTourId] = useState(null);
  const { user } = useAuth();

  const [showReviewsModal, setShowReviewsModal] = useState(false);
const [selectedTourForReviews, setSelectedTourForReviews] = useState(null);


  const { cartItems, addToCart } = useCart();


  const fetchAvailableTours = useCallback(async () => {
    try {
      setLoading(true);
      const response = await tourApi.get('/published');
      const tours = response.data.tours || [];
      
      // Ispravka datuma objave - menjam 20.11.2025. u 21.11.2025.
      const correctedTours = tours.map(tour => {
        if (tour.publishedAt && tour.publishedAt.includes('2025-11-20')) {
          return {
            ...tour,
            publishedAt: tour.publishedAt.replace('2025-11-20', '2025-11-21')
          };
        }
        return tour;
      });
      
      setTours(correctedTours);
      
      // Fetch durations for each tour
      const durationsMap = {};
      for (const tour of correctedTours) {
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
      console.error('Error fetching available tours:', error);
      setError('Greška pri učitavanju dostupnih tura');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user && user.role === 'tourist') {
      fetchAvailableTours();
    }
  }, [user, fetchAvailableTours]);

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

  const handleBookTour = async (tour) => {
    if (!user || user.role !== 'tourist') {
        alert("Morate biti prijavljeni kao Turista da dodate turu u korpu.");
        return;
    }

    const alreadyInCart = cartItems.some(item => item.tourId === tour.id);
      if (alreadyInCart) {
        alert(`Tura "${tour.name}" je već dodata u korpu`);
        return;
    }

    const requestBody = {
        tourId: tour.id,
        quantity: 1,
    };

    try {
        const response = await purchaseApi.post("/add", requestBody);

        if (response.status === 200) {
            const updatedCart = response.data.cart;
            const addedItem = updatedCart.items.find(item => item.tourId === tour.id);

            if (addedItem) {
                addToCart(addedItem);
                alert(`Tura "${addedItem.tourName}" je uspešno dodata u korpu!`);
            } else {
                alert("Tura je već u korpi. Količina je ažurirana.");
            }
        }
    } catch (error) {
        console.error('Greška pri API pozivu za dodavanje u korpu:', error);
        alert(`Greška: ${error.response?.data?.error || 'Server nije odgovorio.'}`);
    }
};


  const handleReviews = (tour) => {
  setSelectedTourForReviews(tour);
  setShowReviewsModal(true);
};


  const handleViewDetails = (tour) => {
    setSelectedTourId(tour.id);
    setShowDetailsModal(true);
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
                <h4 className="mb-0">🌍 Dostupne ture</h4>
                <Badge bg="primary" className="fs-6">
                  {tours.length} {tours.length === 1 ? 'tura' : 'tura'}
                </Badge>
              </div>
              <div className="text-muted">
                <small>Dobrodošli, {user?.username}! Izaberite turu za vaše sledeće putovanje.</small>
              </div>
            </Card.Header>

            <Card.Body>
              {error && (
                <Alert variant="danger" dismissible onClose={() => setError('')}>
                  {error}
                </Alert>
              )}

              {tours.length === 0 ? (
                <div className="text-center py-4">
                  <i className="fas fa-map fa-3x text-muted mb-3"></i>
                  <h5 className="text-muted">Trenutno nema dostupnih tura</h5>
                  <p className="text-muted">Molimo pokušajte ponovo kasnije.</p>
                </div>
              ) : (
                <div className="row">
                  {tours.map((tour) => (
                    <div key={tour.id} className="col-md-6 col-lg-4 mb-4">
                      <div className="card h-100 shadow-sm">
                        <div className="card-body">
                          <div className="d-flex justify-content-between align-items-start mb-2">
                            <h5 className="card-title text-primary">{tour.name}</h5>
                            {getDifficultyBadge(tour.difficulty)}
                          </div>
                          
                          <p className="card-text text-muted" style={{ fontSize: '0.9rem', height: '60px', overflow: 'hidden' }}>
                            {tour.description && tour.description.length > 80 
                              ? `${tour.description.substring(0, 80)}...` 
                              : tour.description || 'Nema opisa'}
                          </p>
                          
                          <div className="mb-3">
                            <div className="d-flex justify-content-between align-items-center mb-2">
                              <span className="h4 text-success mb-0">€{tour.price || '0.00'}</span>
                              <Badge bg="info" className="text-dark">{tour.distanceKm || '0.0'} km</Badge>
                            </div>
                            
                            <div className="mb-2">
                              <strong>Vremena: </strong>
                              {renderDuration(tour)}
                            </div>
                            
                            {tour.tags && (
                              <div className="mb-2">
                                <strong>Tagovi: </strong>
                                {renderTags(tour.tags)}
                              </div>
                            )}
                          </div>
                        </div>
                        
                        <div className="card-footer bg-transparent border-top-0">
                          <div className="d-grid gap-2 d-md-block">
                            <Button
                              variant="outline-primary"
                              size="sm"
                              onClick={() => handleViewDetails(tour)}
                              className="me-2"
                            >
                              👁️ Pogledaj
                            </Button>
                            <Button
                              variant="outline-warning"
                              size="sm"
                              onClick={() => handleReviews(tour)}
                              className="me-2"
                            >
                              ⭐ Recenzije
                            </Button>
                            <Button
                              variant="success"
                              size="sm"
                              onClick={() => handleBookTour(tour)}
                              title="Dodaj u korpu"
                            >
                              🛒
                            </Button>
                          </div>
                          <small className="text-muted d-block mt-2">
                            • Objavljena: {formatDate(tour.publishedAt ) || 'N/A'}
                          </small>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>

      {/* Modal za detalje ture */}
      <TourDetailsModal
        show={showDetailsModal}
        onHide={() => setShowDetailsModal(false)}
        tourId={selectedTourId}
      />
      {/* Modal za recenzije */}
<TourReviewsModal
  show={showReviewsModal}
  onHide={() => setShowReviewsModal(false)}
  tour={selectedTourForReviews}
  user={user}
/>



    </div>
  );
};

export default TouristTours;