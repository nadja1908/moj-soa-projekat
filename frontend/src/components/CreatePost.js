import React, { useState } from 'react';
import { Form, Button, Card, Alert, Container } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { blogApi } from '../services/api';
import { useAuth } from '../context/AuthContext';

const CreatePost = () => {
  const [formData, setFormData] = useState({
    title: '',
    content: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { user } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (formData.title.trim().length < 5) {
      setError('Naslov mora imati najmanje 5 karaktera');
      setLoading(false);
      return;
    }

    if (formData.content.trim().length < 20) {
      setError('Sadržaj mora imati najmanje 20 karaktera');
      setLoading(false);
      return;
    }

    try {
      await blogApi.post('/posts', {
        title: formData.title.trim(),
        content: formData.content.trim()
      });
      
      navigate('/posts');
    } catch (error) {
      console.error('Error creating post:', error);
      setError(error.response?.data?.error || 'Greška pri kreiranju posta');
    }
    
    setLoading(false);
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  if (!user || (user.role !== 'guide' && user.role !== 'administrator')) {
    return (
      <Container>
        <Alert variant="warning" className="text-center">
          <Alert.Heading>⚠️ Nemate dozvolu</Alert.Heading>
          <p>Samo vodiči i administratori mogu kreirati blog postove.</p>
        </Alert>
      </Container>
    );
  }

  return (
    <div className="container">
      <div className="row justify-content-center">
        <div className="col-lg-8">
          <Card className="border-0 shadow-sm">
            <Card.Body className="p-4">
              <div className="text-center mb-4">
                <h2 className="display-6 text-primary">✍️ Kreiraj novi blog post</h2>
                <p className="lead text-muted">Podelite svoja iskustva sa zajednicom</p>
              </div>

              {error && <Alert variant="danger">{error}</Alert>}

              <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📝 Naslov posta</Form.Label>
                  <Form.Control
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleChange}
                    placeholder="Unesite naslov vašeg blog posta..."
                    required
                    maxLength={200}
                    className="form-control-lg"
                  />
                  <Form.Text className="text-muted">
                    {formData.title.length}/200 karaktera
                  </Form.Text>
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📖 Sadržaj posta</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={12}
                    name="content"
                    value={formData.content}
                    onChange={handleChange}
                    placeholder="Napišite sadržaj vašeg blog posta...&#10;&#10;Možete pisati o:&#10;• Turističkim destinacijama koje ste posetili&#10;• Savetima za putovanja&#10;• Lokalnim specijalitetima&#10;• Kulturnim atrakcijama&#10;• Vašim iskustvima kao vodič/turista"
                    required
                    maxLength={5000}
                  />
                  <Form.Text className="text-muted">
                    {formData.content.length}/5000 karaktera
                  </Form.Text>
                </Form.Group>

                <div className="d-flex gap-3 mb-4">
                  <Button 
                    variant="primary" 
                    type="submit" 
                    disabled={loading}
                    className="flex-grow-1"
                    size="lg"
                  >
                    {loading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        Objavljivanje...
                      </>
                    ) : (
                      '🚀 Objavi post'
                    )}
                  </Button>
                  
                  <Button 
                    variant="outline-secondary" 
                    onClick={() => navigate('/posts')}
                    disabled={loading}
                    size="lg"
                  >
                    ❌ Odustani
                  </Button>
                </div>
              </Form>

              <div className="p-4 bg-light rounded">
                <h6 className="mb-3">💡 Saveti za pisanje dobrog blog posta:</h6>
                <div className="row">
                  <div className="col-md-6">
                    <ul className="small text-muted mb-0">
                      <li>Koristite zanimljiv i deskriptivan naslov</li>
                      <li>Pišite jasno i razumljivo</li>
                      <li>Podelite lična iskustva i savete</li>
                    </ul>
                  </div>
                  <div className="col-md-6">
                    <ul className="small text-muted mb-0">
                      <li>Dodajte korisne informacije za čitaoce</li>
                      <li>Budite poštovani prema svim korisnicima</li>
                      <li>Strukturirajte tekst u logičke celine</li>
                    </ul>
                  </div>
                </div>
              </div>
            </Card.Body>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default CreatePost;