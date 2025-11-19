import React, { useState, useEffect, useRef } from 'react';
import { Modal, Button, Form, Alert, Card, Badge, ListGroup } from 'react-bootstrap';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import { tourApi, keypointsApi } from '../services/api';
import 'leaflet/dist/leaflet.css';

// Leaflet ikone fix
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

// Komponenta za hvatanje klikova na mapu
const MapClickHandler = ({ onMapClick, isAddingMode }) => {
  useMapEvents({
    click: (e) => {
      if (isAddingMode) {
        onMapClick(e.latlng);
      }
    },
  });
  return null;
};

const KeyPointsMap = ({ show, onHide, tourId, tourName }) => {
  const [keyPoints, setKeyPoints] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [isAddingMode, setIsAddingMode] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [selectedCoordinates, setSelectedCoordinates] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    image: null,
    imagePreview: ''
  });

  // Default centar mape (Beograd)
  const [mapCenter] = useState([44.8176, 20.4633]);
  const mapRef = useRef();

  console.log('KeyPointsMap render - keyPoints:', keyPoints, 'type:', typeof keyPoints, 'isArray:', Array.isArray(keyPoints));

  useEffect(() => {
    if (show && tourId) {
      fetchKeyPoints();
    }
  }, [show, tourId]);

  const fetchKeyPoints = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await keypointsApi.get(`/tour/${tourId}`);
      console.log('KeyPoints API response:', response.data);
      setKeyPoints(response.data.keyPoints || []);
    } catch (err) {
      console.error('Error fetching key points:', err);
      setError('Greška pri učitavanju ključnih tačaka');
    } finally {
      setLoading(false);
    }
  };

  const handleMapClick = (latlng) => {
    setSelectedCoordinates(latlng);
    setShowAddForm(true);
    setIsAddingMode(false);
  };

  const handleAddKeyPoint = async (e) => {
    e.preventDefault();
    if (!selectedCoordinates) return;

    try {
      setLoading(true);
      
      // Kreiraj FormData objekat
      const formDataToSend = new FormData();
      formDataToSend.append('tourId', tourId);
      formDataToSend.append('name', formData.name);
      formDataToSend.append('description', formData.description);
      formDataToSend.append('latitude', selectedCoordinates.lat);
      formDataToSend.append('longitude', selectedCoordinates.lng);
      formDataToSend.append('orderIndex', keyPoints.length);
      
      // Dodaj sliku ako je selektovana
      if (formData.image) {
        formDataToSend.append('image', formData.image);
      }

      await keypointsApi.post(``, formDataToSend, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      // Reset form
      setFormData({ name: '', description: '', image: null, imagePreview: '' });
      setSelectedCoordinates(null);
      setShowAddForm(false);
      
      // Refresh key points
      await fetchKeyPoints();
    } catch (err) {
      console.error('Error adding key point:', err);
      setError(err.response?.data?.error || 'Greška pri dodavanju ključne tačke');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteKeyPoint = async (keyPointId) => {
    if (!window.confirm('Da li ste sigurni da želite da obrišete ovu ključnu tačku?')) {
      return;
    }

    try {
      setLoading(true);
      await keypointsApi.delete(`/${keyPointId}`);
      await fetchKeyPoints();
    } catch (err) {
      console.error('Error deleting key point:', err);
      setError('Greška pri brisanju ključne tačke');
    } finally {
      setLoading(false);
    }
  };

  const handleFormChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setFormData({
        ...formData,
        image: file,
        imagePreview: URL.createObjectURL(file)
      });
    }
  };

  return (
    <Modal show={show} onHide={onHide} size="xl" centered>
      <Modal.Header closeButton>
        <Modal.Title>
          📍 Ključne tačke ture: {tourName}
        </Modal.Title>
      </Modal.Header>

      <Modal.Body style={{ height: '70vh', padding: 0 }}>
        {error && (
          <Alert variant="danger" className="m-3">
            {error}
          </Alert>
        )}

        <div className="d-flex h-100">
          {/* Leva strana - Lista ključnih tačaka */}
          <div className="col-md-4 border-end">
            <div className="p-3">
              <div className="d-flex justify-content-between align-items-center mb-3">
                <h6 className="mb-0">Ključne tačke ({keyPoints.length})</h6>
                <Button
                  variant={isAddingMode ? "success" : "outline-primary"}
                  size="sm"
                  onClick={() => setIsAddingMode(!isAddingMode)}
                >
                  {isAddingMode ? '✓ Klik na mapu' : '+ Dodaj tačku'}
                </Button>
              </div>

              {keyPoints.length === 0 ? (
                <div className="text-center text-muted py-4">
                  <i className="fas fa-map-marker-alt fa-2x mb-2"></i>
                  <p>Nema dodanih ključnih tačaka</p>
                  <small>Kliknite "Dodaj tačku" pa kliknite na mapu</small>
                </div>
              ) : (
                <ListGroup variant="flush">
                  {Array.isArray(keyPoints) ? keyPoints.map((point, index) => (
                    <ListGroup.Item key={point.id} className="px-0">
                      <div className="d-flex justify-content-between align-items-start">
                        <div className="flex-grow-1">
                          <h6 className="mb-1">
                            <Badge bg="primary" className="me-2">{index + 1}</Badge>
                            {point.name}
                          </h6>
                          <p className="mb-1 text-muted small">
                            {point.description}
                          </p>
                          <small className="text-muted">
                            📍 {point.latitude.toFixed(4)}, {point.longitude.toFixed(4)}
                          </small>
                        </div>
                        <Button
                          variant="outline-danger"
                          size="sm"
                          onClick={() => handleDeleteKeyPoint(point.id)}
                        >
                          <i className="fas fa-trash"></i>
                        </Button>
                      </div>
                    </ListGroup.Item>
                  )) : null}
                </ListGroup>
              )}
            </div>
          </div>

          {/* Desna strana - Mapa */}
          <div className="col-md-8 position-relative">
            {isAddingMode && (
              <div className="position-absolute top-0 start-0 end-0 bg-success text-white text-center py-2 z-index-1000" style={{ zIndex: 1000 }}>
                <i className="fas fa-mouse-pointer me-2"></i>
                Kliknite na mapu da dodate ključnu tačku
              </div>
            )}

            <MapContainer
              center={mapCenter}
              zoom={13}
              style={{ height: '100%', width: '100%' }}
              ref={mapRef}
            >
              <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              />
              
              <MapClickHandler 
                onMapClick={handleMapClick}
                isAddingMode={isAddingMode}
              />

              {/* Markeri za postojeće ključne tačke */}
              {Array.isArray(keyPoints) ? keyPoints.map((point, index) => (
                <Marker
                  key={point.id}
                  position={[point.latitude, point.longitude]}
                >
                  <Popup>
                    <div>
                      <h6>
                        <Badge bg="primary" className="me-2">{index + 1}</Badge>
                        {point.name}
                      </h6>
                      <p className="mb-2">{point.description}</p>
                      {point.imageUrl && (
                        <img 
                          src={point.imageUrl.startsWith('/') ? `http://localhost:8004${point.imageUrl}` : point.imageUrl} 
                          alt={point.name}
                          style={{ width: '100%', maxWidth: '200px', height: 'auto' }}
                        />
                      )}
                      <div className="mt-2">
                        <small className="text-muted">
                          📍 {point.latitude.toFixed(4)}, {point.longitude.toFixed(4)}
                        </small>
                      </div>
                    </div>
                  </Popup>
                </Marker>
              )) : null}

              {/* Temporarni marker za novu tačku */}
              {selectedCoordinates && (
                <Marker position={[selectedCoordinates.lat, selectedCoordinates.lng]}>
                  <Popup>
                    <div className="text-center">
                      <h6>Nova ključna tačka</h6>
                      <p>📍 {selectedCoordinates.lat.toFixed(4)}, {selectedCoordinates.lng.toFixed(4)}</p>
                    </div>
                  </Popup>
                </Marker>
              )}
            </MapContainer>
          </div>
        </div>
      </Modal.Body>

      <Modal.Footer>
        <Button variant="secondary" onClick={onHide}>
          Zatvori
        </Button>
      </Modal.Footer>

      {/* Modal za dodavanje nove ključne tačke */}
      <Modal show={showAddForm} onHide={() => setShowAddForm(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Dodaj ključnu tačku</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleAddKeyPoint}>
          <Modal.Body>
            {selectedCoordinates && (
              <Alert variant="info">
                📍 Lokacija: {selectedCoordinates.lat.toFixed(6)}, {selectedCoordinates.lng.toFixed(6)}
              </Alert>
            )}

            <Form.Group className="mb-3">
              <Form.Label>Naziv ključne tačke *</Form.Label>
              <Form.Control
                type="text"
                name="name"
                value={formData.name}
                onChange={handleFormChange}
                placeholder="npr. Muzej, Park, Spomenik..."
                required
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Opis</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                name="description"
                value={formData.description}
                onChange={handleFormChange}
                placeholder="Opišite ovu ključnu tačku..."
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Slika (opciono)</Form.Label>
              <Form.Control
                type="file"
                accept="image/*"
                onChange={handleFileChange}
              />
              {formData.imagePreview && (
                <div className="mt-2">
                  <img 
                    src={formData.imagePreview} 
                    alt="Preview" 
                    style={{ maxWidth: '200px', maxHeight: '150px', objectFit: 'cover' }}
                    className="rounded"
                  />
                </div>
              )}
            </Form.Group>
          </Modal.Body>

          <Modal.Footer>
            <Button 
              variant="secondary" 
              onClick={() => {
                setShowAddForm(false);
                setSelectedCoordinates(null);
                setFormData({ name: '', description: '', image: null, imagePreview: '' });
              }}
            >
              Odustani
            </Button>
            <Button type="submit" variant="primary" disabled={loading}>
              {loading ? 'Dodavanje...' : 'Dodaj ključnu tačku'}
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </Modal>
  );
};

export default KeyPointsMap;