import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Modal, Button, Form, Alert, Card, Badge, ListGroup } from 'react-bootstrap';
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMapEvents } from 'react-leaflet';
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
const MapClickHandler = ({ onMapClick, isAddingMode, isEditMode }) => {
  useMapEvents({
    click: (e) => {
      console.log('🗺️ MapClickHandler received click:', { isAddingMode, isEditMode });
      if (isAddingMode || isEditMode) {
        console.log('🎯 Calling onMapClick with:', e.latlng);
        onMapClick(e.latlng);
      } else {
        console.log('❌ Click ignored - neither adding nor editing mode');
      }
    },
  });
  return null;
};

const KeyPointsMap = React.memo(({ show, onHide, tourId, tourName, onTourUpdate }) => {
  console.log('🗺️ KeyPointsMap RENDER:', { show, tourId, tourName, onTourUpdate: !!onTourUpdate });
  
  const [keyPoints, setKeyPoints] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [isAddingMode, setIsAddingMode] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [editingKeyPoint, setEditingKeyPoint] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [showEditForm, setShowEditForm] = useState(false);
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

  // Stabiliziraj tour podatke da sprečiš nepotrebne re-render-e
  // Debug logging za state promene
  useEffect(() => {
    console.log('🎯 isAddingMode changed to:', isAddingMode);
  }, [isAddingMode]);

  useEffect(() => {
    console.log('✏️ isEditMode changed to:', isEditMode);
  }, [isEditMode]);

  useEffect(() => {
    console.log('📍 showAddForm changed to:', showAddForm);
  }, [showAddForm]);

  useEffect(() => {
    console.log('📝 editingKeyPoint changed to:', editingKeyPoint);
  }, [editingKeyPoint]);

  const tourData = useMemo(() => {
    console.log('🔄 tourData useMemo:', { id: tourId, name: tourName });
    return { id: tourId, name: tourName };
  }, [tourId, tourName]);

  // Debugging lifecycle
  useEffect(() => {
    console.log('🟢 KeyPointsMap MOUNTED/UPDATED:', { show, tourData });
    return () => {
      console.log('🔴 KeyPointsMap UNMOUNTING');
    };
  }, [show, tourData]);

  console.log('KeyPointsMap render - keyPoints:', keyPoints, 'type:', typeof keyPoints, 'isArray:', Array.isArray(keyPoints));

  const fetchKeyPoints = useCallback(async () => {
    console.log('📡 fetchKeyPoints CALLED for tour:', tourData.id);
    try {
      setLoading(true);
      setError('');
      const response = await keypointsApi.get(`/tour/${tourData.id}`);
      console.log('✅ KeyPoints API response:', response.data);
      setKeyPoints(response.data.keyPoints || []);
      console.log('✅ KeyPoints fetched successfully - backend automatically calculates price');
    } catch (err) {
      console.error('❌ Error fetching key points:', err);
      setError('Greška pri učitavanju ključnih tačaka');
    } finally {
      setLoading(false);
    }
  }, [tourData.id]);

  useEffect(() => {
    if (show && tourData.id) {
      fetchKeyPoints();
    }
  }, [show, tourData.id, fetchKeyPoints]);

  const handleMapClick = (latlng) => {
    console.log('🗺️ Map clicked at:', latlng, 'isAddingMode:', isAddingMode, 'isEditMode:', isEditMode, 'showEditForm:', showEditForm);
    
    if (isAddingMode) {
      console.log('➕ Adding new keypoint at:', latlng);
      setSelectedCoordinates(latlng);
      setShowAddForm(true);
      setIsAddingMode(false);
    } else if (isEditMode && editingKeyPoint) {
      console.log('✏️ Moving keypoint to new position:', latlng);
      // Update position of editing keypoint
      setEditingKeyPoint({
        ...editingKeyPoint,
        latitude: latlng.lat,
        longitude: latlng.lng
      });
      setSelectedCoordinates(latlng);
      // Otvori edit modal TEK SADA
      setShowEditForm(true);
      setIsEditMode(false);
      console.log('✏️ Position updated, opening edit modal');
    } else {
      console.log('❌ Map click ignored - not in adding or edit mode');
    }
  };

  const handleAddKeyPoint = async (e) => {
    e.preventDefault();
    if (!selectedCoordinates) return;

    try {
      setLoading(true);
      
      // Kreiraj FormData objekat
      const formDataToSend = new FormData();
      formDataToSend.append('tourId', tourData.id);
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

  const handleEditKeyPoint = (keyPoint) => {
    console.log('✏️ handleEditKeyPoint called with:', keyPoint);
    setEditingKeyPoint(keyPoint);
    setIsEditMode(true);
    setIsAddingMode(false);
    setShowEditForm(false); // Zatvori modal ako je otvoren
    console.log('✏️ Edit mode activated - click on map to move keypoint:', keyPoint.id);
  };

  const handleUpdateKeyPoint = async (e) => {
    e.preventDefault();
    if (!editingKeyPoint) return;

    try {
      setLoading(true);
      
      // Koristi selectedCoordinates ako su postavljene, inače originalne koordinate
      const finalLatitude = selectedCoordinates ? selectedCoordinates.lat : editingKeyPoint.latitude;
      const finalLongitude = selectedCoordinates ? selectedCoordinates.lng : editingKeyPoint.longitude;
      
      const updateData = {
        name: editingKeyPoint.name,
        description: editingKeyPoint.description,
        latitude: finalLatitude,
        longitude: finalLongitude,
      };

      console.log('✏️ Updating keypoint:', editingKeyPoint.id);
      console.log('✏️ Original coordinates:', { lat: editingKeyPoint.latitude, lng: editingKeyPoint.longitude });
      console.log('✏️ Selected coordinates:', selectedCoordinates);
      console.log('✏️ Final update data:', updateData);
      
      const response = await keypointsApi.put(`/${editingKeyPoint.id}`, updateData);
      console.log('✅ Update response:', response);
      
      // Reset form and state
      setEditingKeyPoint(null);
      setSelectedCoordinates(null);
      setShowEditForm(false);
      setIsEditMode(false);
      
      // Refresh key points
      await fetchKeyPoints();
    } catch (err) {
      console.error('❌ Error updating key point:', err);
      console.error('❌ Error details:', err.response?.data);
      setError(err.response?.data?.error || err.message || 'Greška pri ažuriranju ključne tačke');
    } finally {
      setLoading(false);
    }
  };

  const handleEditFormChange = (e) => {
    setEditingKeyPoint({
      ...editingKeyPoint,
      [e.target.name]: e.target.value
    });
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
                <div>
              <div className="d-flex gap-2 mb-3">
                <Button
                  variant={isAddingMode ? "success" : "outline-primary"}
                  size="sm"
                  onClick={() => {
                    console.log('🎯 Dodaj tačku clicked, current isAddingMode:', isAddingMode);
                    setIsAddingMode(!isAddingMode);
                    setIsEditMode(false);
                    setEditingKeyPoint(null);
                    console.log('🎯 New isAddingMode will be:', !isAddingMode);
                  }}
                  title={isAddingMode ? 'Kliknite na mapu da dodate tačku' : 'Aktiviraj režim dodavanja'}
                  style={{ minWidth: '100px' }}
                >
                  <i className={isAddingMode ? "fas fa-map-marker-alt me-1" : "fas fa-plus me-1"}></i>
                  {isAddingMode ? 'Klik ovde' : 'Dodaj'}
                </Button>
                <Button
                  variant={isEditMode ? "warning" : "outline-secondary"}
                  size="sm"
                  onClick={() => {
                    console.log('✏️ Izmeni button clicked, current isEditMode:', isEditMode);
                    setIsEditMode(!isEditMode);
                    setIsAddingMode(false);
                    if (!isEditMode) {
                      console.log('✏️ Clearing editingKeyPoint');
                      setEditingKeyPoint(null);
                    }
                    console.log('✏️ New isEditMode will be:', !isEditMode);
                  }}
                  disabled={keyPoints.length === 0}
                  title={isEditMode ? 'Otkaži izmenu' : 'Aktiviraj režim izmene'}
                  style={{ minWidth: '85px' }}
                >
                  <i className={isEditMode ? "fas fa-times me-1" : "fas fa-edit me-1"}></i>
                  {isEditMode ? 'Otkaži' : 'Izmeni'}
                </Button>
              </div>
                </div>
              </div>

              {isEditMode && (
                <Alert variant="info" className="mb-3" style={{ padding: '0.5rem', fontSize: '0.85rem' }}>
                  <i className="fas fa-info-circle me-1"></i>
                  Kliknite na bilo koju ključnu tačku na mapi da je pomerite
                </Alert>
              )}
              
              {isAddingMode && (
                <Alert variant="success" className="mb-3" style={{ padding: '0.5rem', fontSize: '0.85rem' }}>
                  <i className="fas fa-map-marker-alt me-1"></i>
                  Kliknite na mapu da dodate novu ključnu tačku
                </Alert>
              )}

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
                        <div className="d-flex gap-2">
                          <Button
                            variant="outline-secondary"
                            size="sm"
                            onClick={() => {
                              console.log('📝 Edit button clicked for point:', point.id);
                              handleEditKeyPoint(point);
                            }}
                            title="Izmeni ključnu tačku"
                          >
                            <i className="fas fa-edit me-1"></i>
                            Izmeni
                          </Button>
                          <Button
                            variant="outline-danger"
                            size="sm"
                            onClick={() => handleDeleteKeyPoint(point.id)}
                            title="Obriši ključnu tačku"
                          >
                            <i className="fas fa-trash me-1"></i>
                            Obriši
                          </Button>
                        </div>
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
            {isEditMode && !editingKeyPoint && (
              <div className="position-absolute top-0 start-0 end-0 bg-warning text-dark text-center py-2 z-index-1000" style={{ zIndex: 1000 }}>
                <i className="fas fa-edit me-2"></i>
                Izaberite ključnu tačku koju želite da izmeníte
              </div>
            )}
            {isEditMode && editingKeyPoint && (
              <div className="position-absolute top-0 start-0 end-0 bg-warning text-dark text-center py-2 z-index-1000" style={{ zIndex: 1000 }}>
                <i className="fas fa-crosshairs me-2"></i>
                🎯 Kliknite na mapu da pomerite: <strong>{editingKeyPoint.name}</strong>
              </div>
            )}

            <MapContainer
              key={`map-${tourData.id}`}
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
                isEditMode={isEditMode}
              />              {/* Crtanje ture - povezivanje ključnih tačaka linijom */}
              {Array.isArray(keyPoints) && keyPoints.length > 1 && (
                <Polyline
                  positions={keyPoints.map(point => [point.latitude, point.longitude])}
                  color="#007bff"
                  weight={3}
                  opacity={0.7}
                  dashArray="5, 10"
                />
              )}

              {/* Marker za novu poziciju tokom editovanja */}
              {selectedCoordinates && isEditMode && editingKeyPoint && (
                <Marker
                  position={[selectedCoordinates.lat, selectedCoordinates.lng]}
                  icon={L.icon({
                    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
                    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
                    iconSize: [25, 41],
                    iconAnchor: [12, 41],
                    popupAnchor: [1, -34],
                    shadowSize: [41, 41]
                  })}
                >
                  <Popup>
                    <div className="text-center">
                      <div className="text-danger">
                        <i className="fas fa-crosshairs"></i> <strong>Nova pozicija</strong>
                      </div>
                      <div>za: {editingKeyPoint.name}</div>
                    </div>
                  </Popup>
                </Marker>
              )}

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
        <Button variant="secondary" size="sm" onClick={onHide}>
          <i className="fas fa-times me-1"></i>
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
          <Modal.Footer>
            <Button 
              variant="secondary" 
              size="sm"
              onClick={() => {
                setShowAddForm(false);
                setSelectedCoordinates(null);
                setFormData({ name: '', description: '', image: null, imagePreview: '' });
              }}
            >
              <i className="fas fa-times me-1"></i>
              Odustani
            </Button>
            <Button type="submit" variant="primary" size="sm" disabled={loading}>
              <i className={loading ? "fas fa-spinner fa-spin me-1" : "fas fa-save me-1"}></i>
              {loading ? 'Dodavanje...' : 'Dodaj'}
            </Button>
          </Modal.Footer>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Modal za izmenu ključne tačke */}
      <Modal show={showEditForm} onHide={() => {
        setShowEditForm(false);
        setEditingKeyPoint(null);
        setSelectedCoordinates(null);
      }}>
        <Modal.Header closeButton>
          <Modal.Title>Izmeni ključnu tačku</Modal.Title>
        </Modal.Header>
        {editingKeyPoint && (
          <Form onSubmit={handleUpdateKeyPoint}>
            <Modal.Body>
              {selectedCoordinates && (
                <Alert variant="info">
                  📍 Nova lokacija: {selectedCoordinates.lat.toFixed(6)}, {selectedCoordinates.lng.toFixed(6)}
                </Alert>
              )}

              <Form.Group className="mb-3">
                <Form.Label>Naziv ključne tačke *</Form.Label>
                <Form.Control
                  type="text"
                  value={editingKeyPoint.name}
                  onChange={(e) => setEditingKeyPoint({...editingKeyPoint, name: e.target.value})}
                  placeholder="npr. Muzej, Park, Spomenik..."
                  required
                />
              </Form.Group>

              <Form.Group className="mb-3">
                <Form.Label>Opis</Form.Label>
                <Form.Control
                  as="textarea"
                  rows={3}
                  value={editingKeyPoint.description || ''}
                  onChange={(e) => setEditingKeyPoint({...editingKeyPoint, description: e.target.value})}
                  placeholder="Opišite ovu ključnu tačku..."
                />
              </Form.Group>

              <Form.Group className="mb-3">
                <Form.Label>Pozicija</Form.Label>
                {selectedCoordinates ? (
                  <div>
                    <div className="text-success">
                      📍 <strong>Nova pozicija:</strong> {selectedCoordinates.lat.toFixed(6)}, {selectedCoordinates.lng.toFixed(6)}
                    </div>
                    <div className="text-muted small">
                      Stara pozicija: {editingKeyPoint.latitude.toFixed(6)}, {editingKeyPoint.longitude.toFixed(6)}
                    </div>
                  </div>
                ) : (
                  <div className="text-muted">
                    📍 {editingKeyPoint.latitude.toFixed(6)}, {editingKeyPoint.longitude.toFixed(6)}
                  </div>
                )}
              </Form.Group>
            </Modal.Body>

            <Modal.Footer>
            <Modal.Footer>
              <Button 
                variant="secondary" 
                size="sm"
                onClick={() => {
                  setShowEditForm(false);
                  setEditingKeyPoint(null);
                  setSelectedCoordinates(null);
                }}
              >
                <i className="fas fa-times me-1"></i>
                Odustani
              </Button>
              <Button type="submit" variant="warning" size="sm" disabled={loading}>
                <i className={loading ? "fas fa-spinner fa-spin me-1" : "fas fa-save me-1"}></i>
                {loading ? 'Ažuriranje...' : 'Sačuvaj'}
              </Button>
            </Modal.Footer>
            </Modal.Footer>
          </Form>
        )}
      </Modal>
    </Modal>
  );
});

// Memoization comparison za stabilnost
KeyPointsMap.displayName = 'KeyPointsMap';

export default KeyPointsMap;