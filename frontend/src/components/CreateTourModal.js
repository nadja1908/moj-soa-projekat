import React, { useState } from 'react';
import { Modal, Button, Form, Alert, Spinner } from 'react-bootstrap';
import { tourApi } from '../services/api';

const CreateTourModal = ({ show, onHide, onTourCreated }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    difficulty: 'EASY',
    price: 0,
    tags: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'number' ? parseFloat(value) || 0 : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validacija
    if (!formData.name.trim()) {
      setError('Naziv ture je obavezan');
      return;
    }
    
    if (!formData.description.trim()) {
      setError('Opis ture je obavezan');
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      const tourData = {
        ...formData,
        name: formData.name.trim(),
        description: formData.description.trim(),
        tags: formData.tags.trim() ? formData.tags.trim().split(',').map(tag => tag.trim()).filter(tag => tag) : []
      };
      
      console.log('Sending tour data:', tourData);
      const response = await tourApi.post('/', tourData);
      console.log('Tour creation response:', response);
      
      if (response.data.success) {
        console.log('Tour created successfully:', response.data.tour);
        onTourCreated(response.data.tour);
        
        // Reset form
        setFormData({
          name: '',
          description: '', 
          difficulty: 'EASY',
          price: 0,
          tags: ''
        });
      }
    } catch (error) {
      console.error('Error creating tour:', error);
      console.error('Error response:', error.response);
      setError(error.response?.data?.error || 'Greška pri kreiranju ture');
    } finally {
      setLoading(false);
    }
  };

  const handleModalHide = () => {
    if (!loading) {
      setError('');
      setFormData({
        name: '',
        description: '',
        difficulty: 'EASY', 
        price: 0,
        tags: ''
      });
      onHide();
    }
  };

  return (
    <Modal show={show} onHide={handleModalHide} size="lg" backdrop={loading ? 'static' : true}>
      <Modal.Header closeButton={!loading}>
        <Modal.Title>
          <i className="fas fa-plus-circle me-2 text-success"></i>
          Kreiraj novu turu
        </Modal.Title>
      </Modal.Header>

      <Form onSubmit={handleSubmit}>
        <Modal.Body>
          {error && (
            <Alert variant="danger" className="mb-3">
              <i className="fas fa-exclamation-triangle me-2"></i>
              {error}
            </Alert>
          )}

          <div className="row">
            <div className="col-md-8">
              <Form.Group className="mb-3">
                <Form.Label>
                  <i className="fas fa-route me-2"></i>
                  Naziv ture *
                </Form.Label>
                <Form.Control
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Unesite naziv ture..."
                  disabled={loading}
                  required
                />
              </Form.Group>
            </div>
            
            <div className="col-md-4">
              <Form.Group className="mb-3">
                <Form.Label>
                  <i className="fas fa-chart-line me-2"></i>
                  Težina
                </Form.Label>
                <Form.Select
                  name="difficulty"
                  value={formData.difficulty}
                  onChange={handleChange}
                  disabled={loading}
                >
                  <option value="EASY">Laka</option>
                  <option value="MODERATE">Umerena</option>
                  <option value="HARD">Teška</option>
                </Form.Select>
              </Form.Group>
            </div>
          </div>

          <Form.Group className="mb-3">
            <Form.Label>
              <i className="fas fa-align-left me-2"></i>
              Opis ture *
            </Form.Label>
            <Form.Control
              as="textarea"
              rows={4}
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Opišite šta tura obuhvata, koje lokacije se posećuju, što turisti mogu da očekuju..."
              disabled={loading}
              required
            />
          </Form.Group>

          <div className="row">
            <div className="col-md-6">
              <Form.Group className="mb-3">
                <Form.Label>
                  <i className="fas fa-euro-sign me-2"></i>
                  Cena (€)
                </Form.Label>
                <Form.Control
                  type="number"
                  step="0.01"
                  min="0"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  placeholder="0.00"
                  disabled={loading}
                />
              </Form.Group>
            </div>
            
            <div className="col-md-6">
              <Form.Group className="mb-3">
                <Form.Label>
                  <i className="fas fa-tags me-2"></i>
                  Tagovi (opciono)
                </Form.Label>
                <Form.Control
                  type="text"
                  name="tags"
                  value={formData.tags}
                  onChange={handleChange}
                  placeholder="kultura, istorija, priroda..."
                  disabled={loading}
                />
                <Form.Text className="text-muted">
                  Tagovi razdvojeni zarezom za lakše pretraganje
                </Form.Text>
              </Form.Group>
            </div>
          </div>

          <div className="alert alert-info">
            <i className="fas fa-info-circle me-2"></i>
            <strong>Napomena:</strong> Nakon kreiranja ture, moći ćete da dodate ključne tačke na mapi i postavite putanju. 
            Tura će biti u statusu "Nacrt" dok ne dodate sve potrebne informacije i objavite je.
          </div>
        </Modal.Body>

        <Modal.Footer>
          <Button 
            variant="secondary" 
            onClick={handleModalHide}
            disabled={loading}
          >
            Otkaži
          </Button>
          <Button 
            variant="success" 
            type="submit"
            disabled={loading}
          >
            {loading ? (
              <>
                <Spinner animation="border" size="sm" className="me-2" />
                Kreiranje...
              </>
            ) : (
              <>
                <i className="fas fa-save me-2"></i>
                Kreiraj turu
              </>
            )}
          </Button>
        </Modal.Footer>
      </Form>
    </Modal>
  );
};

export default CreateTourModal;