import React, { useEffect, useState, useRef } from 'react';
import {
  Modal,
  Button,
  Form,
  Alert,
  Spinner,
  Badge,
  Row,
  Col,
  Card
} from 'react-bootstrap';

const API_BASE = 'http://localhost:8004'; // direktno Tour-Service

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const d = new Date(dateString);

  // Ako nije validan datum ili je neka nerealna godina (0001 itd.)
  if (isNaN(d.getTime()) || d.getFullYear() < 2000) {
    return 'N/A';
  }

  return d.toLocaleDateString('sr-RS', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  });
};

const TourReviewsModal = ({ show, onHide, tour, user }) => {
  const [reviews, setReviews] = useState([]);
  const [loadingReviews, setLoadingReviews] = useState(false);
  const [error, setError] = useState('');

  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState('');
  const [dateVisited, setDateVisited] = useState('');

  // slike – sada radimo sa nizom imageUrls (može URL, može base64)
  const [imageUrls, setImageUrls] = useState([]);
  const [newImageUrl, setNewImageUrl] = useState('');

  const [submitting, setSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const todayStr = new Date().toISOString().split('T')[0];

  // ref za skriveni <input type="file">
  const fileInputRef = useRef(null);

  // Učitavanje recenzija za izabranu turu
  const fetchReviews = async () => {
    if (!tour) return;
    try {
      setLoadingReviews(true);
      setError('');

      const res = await fetch(`${API_BASE}/api/reviews/tour/${tour.id}`, {
        method: 'GET'
      });

      if (!res.ok) {
        throw new Error(`Greška ${res.status}`);
      }

      const data = await res.json();
      setReviews(data.reviews || []);
    } catch (err) {
      console.error('Greška pri učitavanju recenzija:', err);
      setError('Greška pri učitavanju recenzija za ovu turu.');
    } finally {
      setLoadingReviews(false);
    }
  };

  // Kad se otvori modal ili promeni tura → povuci recenzije
  useEffect(() => {
    if (show && tour) {
      setSuccessMessage('');
      setError('');
      fetchReviews();
    }
  }, [show, tour]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!tour) return;

    if (!user) {
      setError('Morate biti prijavljeni da biste ostavili recenziju.');
      return;
    }

    if (!dateVisited) {
      setError('Molimo izaberite datum kada ste bili na turi.');
      return;
    }

    // zabrana budućeg datuma i na nivou validacije
    if (dateVisited > todayStr) {
      setError('Ne možete izabrati datum u budućnosti.');
      return;
    }

    if (rating < 1 || rating > 5) {
      setError('Ocena mora biti između 1 i 5.');
      return;
    }

    setSubmitting(true);
    setError('');
    setSuccessMessage('');

    // koristimo niz imageUrls (URL-ovi + base64 slike)
    const cleanImageUrls = imageUrls
      .map((u) => u.trim())
      .filter((u) => u.length > 0);

    try {
      const res = await fetch(`${API_BASE}/api/reviews`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-User-ID': String(user.id || user.userId),
          'X-User-Role': user.role || 'tourist'
        },
        body: JSON.stringify({
          tourId: tour.id,
          rating: Number(rating),
          comment,
          dateVisited, // format: YYYY-MM-DD
          imageUrls: cleanImageUrls
        })
      });

      if (res.status === 401) {
        setError('Morate biti prijavljeni kao turista da biste ostavili recenziju.');
        return;
      }

      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || `Greška ${res.status}`);
      }

      await res.json();

      setSuccessMessage('Uspešno ste ostavili recenziju. Hvala! ❤️');

      // Očisti formu
      setRating(5);
      setComment('');
      setDateVisited('');
      setImageUrls([]);
      setNewImageUrl('');

      // Ponovo učitaj recenzije da se nova odmah pojavi
      fetchReviews();
    } catch (err) {
      console.error('Greška pri čuvanju recenzije:', err);
      setError('Došlo je do greške pri čuvanju recenzije.');
    } finally {
      setSubmitting(false);
    }
  };

  const renderRatingStars = (value) => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(
        <span key={i} style={{ color: i <= value ? '#f5c518' : '#ddd' }}>
          ★
        </span>
      );
    }
    return <span>{stars}</span>;
  };

  // --- Obrada fajlova (slike sa laptopa) → pretvaramo u base64 data URL ---
  const handleFiles = (files) => {
    const fileArray = Array.from(files);
    fileArray.forEach((file) => {
      if (!file.type.startsWith('image/')) return;
      const reader = new FileReader();
      reader.onload = (event) => {
        const dataUrl = event.target.result;
        setImageUrls((prev) => (prev.includes(dataUrl) ? prev : [...prev, dataUrl]));
      };
      reader.readAsDataURL(file);
    });
  };

  // --- Drag & Drop za slike (fajlovi ili URL-ovi) ---

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();

    // 1) ako su prevučeni fajlovi sa računara
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFiles(e.dataTransfer.files);
      e.dataTransfer.clearData();
      return;
    }

    // 2) ako su prevučeni URL-ovi iz browsera
    const uriList = e.dataTransfer.getData('text/uri-list');
    const text = e.dataTransfer.getData('text');

    let urls = [];

    if (uriList) {
      urls = uriList.split('\n').filter((u) => u.trim().length > 0);
    } else if (text) {
      urls = text.split(/\s+/).filter((u) => u.startsWith('http'));
    }

    if (urls.length > 0) {
      setImageUrls((prev) => {
        const merged = [...prev];
        urls.forEach((u) => {
          if (!merged.includes(u)) merged.push(u);
        });
        return merged;
      });
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleAddImageUrl = () => {
    const trimmed = newImageUrl.trim();
    if (!trimmed) return;
    setImageUrls((prev) => (prev.includes(trimmed) ? prev : [...prev, trimmed]));
    setNewImageUrl('');
  };

  const handleRemoveImage = (urlToRemove) => {
    setImageUrls((prev) => prev.filter((u) => u !== urlToRemove));
  };

  return (
    <Modal show={show} onHide={onHide} size="lg" centered>
      <Modal.Header closeButton>
        <Modal.Title>
          ⭐ Recenzije za turu: <span className="text-primary">{tour?.name}</span>
        </Modal.Title>
      </Modal.Header>

      <Modal.Body>
        {error && (
          <Alert variant="danger" onClose={() => setError('')} dismissible>
            {error}
          </Alert>
        )}

        {successMessage && (
          <Alert variant="success" onClose={() => setSuccessMessage('')} dismissible>
            {successMessage}
          </Alert>
        )}

        {/* Lista recenzija */}
        <h5 className="mb-3">Postojeće recenzije</h5>

        {loadingReviews ? (
          <div className="d-flex justify-content-center my-3">
            <Spinner animation="border" />
          </div>
        ) : reviews.length === 0 ? (
          <p className="text-muted">Još uvek nema recenzija za ovu turu. Budite prvi! 🎉</p>
        ) : (
          <div className="mb-4" style={{ maxHeight: '250px', overflowY: 'auto' }}>
            {reviews.map((rev) => (
              <Card key={rev.id} className="mb-2">
                <Card.Body>
                  <Row>
                    <Col xs={12} md={8}>
                      <div className="d-flex align-items-center mb-2">
                        <div
                          style={{
                            width: 40,
                            height: 40,
                            borderRadius: '50%',
                            overflow: 'hidden',
                            marginRight: 10,
                            backgroundColor: '#eee',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center'
                          }}
                        >
                          {rev.touristAvatar && (
                            <img
                              src={rev.touristAvatar}
                              alt={rev.touristUsername}
                              style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                            />
                          )}
                        </div>
                        <div>
                          <strong>{rev.touristUsername}</strong>
                          <div style={{ fontSize: '0.8rem' }}>
                            Posetio/la: {formatDate(rev.dateVisited)} <br />
                            Komentar ostavljen: {formatDate(rev.dateCreated)}
                          </div>
                        </div>
                      </div>
                      <div className="mb-1">{renderRatingStars(rev.rating)}</div>
                      <p className="mb-0" style={{ whiteSpace: 'pre-wrap' }}>
                        {rev.comment}
                      </p>
                    </Col>
                    <Col xs={12} md={4} className="mt-2 mt-md-0">
        {rev.imageUrls && rev.imageUrls.length > 0 && (
  <>
    <small className="text-muted d-block mb-1">Slike:</small>
    <div className="d-flex flex-wrap gap-2">
      {rev.imageUrls.map((url, idx) => (
        <div
          key={idx}
          style={{
            width: 80,
            height: 80,
            borderRadius: 8,
            overflow: 'hidden',
            border: '1px solid #ddd',
            position: 'relative',
            cursor: 'pointer'
          }}
          onClick={() => window.open(url, '_blank')}
        >
          <img
            src={url}
            alt={`review-${idx}`}
            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
            onError={(e) => {
              // ako URL nije dobar, da ne ruži UI
              e.target.style.display = 'none';
            }}
          />
        </div>
      ))}
    </div>
  </>
)}

                    </Col>
                  </Row>
                </Card.Body>
              </Card>
            ))}
          </div>
        )}

        <hr />

        {/* Forma za novu recenziju */}
        <h5 className="mb-3">Ostavite svoju recenziju</h5>

        <Form onSubmit={handleSubmit}>
          <Row className="mb-3">
            <Col md={4}>
              <Form.Group controlId="rating">
                <Form.Label>Ocena (1–5)</Form.Label>
                <Form.Select
                  value={rating}
                  onChange={(e) => setRating(Number(e.target.value))}
                  required
                >
                  {[5, 4, 3, 2, 1].map((val) => (
                    <option key={val} value={val}>
                      {val}
                    </option>
                  ))}
                </Form.Select>
              </Form.Group>
            </Col>
            <Col md={8}>
              <Form.Group controlId="dateVisited">
                <Form.Label>Datum kada ste bili na turi</Form.Label>
                <Form.Control
                  type="date"
                  value={dateVisited}
                  max={todayStr}
                  onChange={(e) => setDateVisited(e.target.value)}
                  required
                />
              </Form.Group>
            </Col>
          </Row>

          <Form.Group className="mb-3" controlId="comment">
            <Form.Label>Komentar</Form.Label>
            <Form.Control
              as="textarea"
              rows={3}
              maxLength={500}
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Kako vam se dopala tura, vodič, atmosfera..."
              required
            />
            <Form.Text muted>{comment.length}/500 karaktera</Form.Text>
          </Form.Group>

          {/* DRAG & DROP + FILE PICKER za slike */}
          <Form.Group className="mb-3">
            <Form.Label>Slike sa ture</Form.Label>

            {/* skriveni input za biranje fajlova sa računara */}
            <input
              type="file"
              accept="image/*"
              multiple
              ref={fileInputRef}
              style={{ display: 'none' }}
              onChange={(e) => {
                if (e.target.files && e.target.files.length > 0) {
                  handleFiles(e.target.files);
                  e.target.value = '';
                }
              }}
            />

            <div
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onClick={() => fileInputRef.current && fileInputRef.current.click()}
              style={{
                border: '2px dashed #ccc',
                borderRadius: '8px',
                padding: '16px',
                textAlign: 'center',
                cursor: 'pointer',
                marginBottom: '10px'
              }}
            >
              <p className="mb-1">
                Kliknite ovde da izaberete slike sa računara ili prevucite slike / linkove slika.
              </p>
              <small className="text-muted">
                Podržano: fajlovi tipa slika (jpg, png...) ili javni URL-ovi (https://...).
              </small>
            </div>

            <div className="d-flex gap-2 mb-2">
              <Form.Control
                type="url"
                placeholder="https://example.com/slika.jpg"
                value={newImageUrl}
                onChange={(e) => setNewImageUrl(e.target.value)}
              />
              <Button variant="outline-primary" onClick={handleAddImageUrl}>
                Dodaj URL
              </Button>
            </div>

            {imageUrls.length > 0 && (
              <div>
                <small className="text-muted d-block mb-1">Pregled slika:</small>
                <div className="d-flex flex-wrap gap-2">
                  {imageUrls.map((url) => (
                    <div
                      key={url}
                      style={{
                        width: 80,
                        height: 80,
                        borderRadius: 8,
                        overflow: 'hidden',
                        border: '1px solid #ddd',
                        position: 'relative'
                      }}
                    >
                      <img
                        src={url}
                        alt="preview"
                        style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                      />
                      <Button
                        size="sm"
                        variant="danger"
                        onClick={() => handleRemoveImage(url)}
                        style={{
                          position: 'absolute',
                          top: 2,
                          right: 2,
                          padding: '0 4px',
                          lineHeight: 1
                        }}
                      >
                        ×
                      </Button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </Form.Group>

          <div className="d-flex justify-content-end">
            <Button variant="secondary" className="me-2" onClick={onHide}>
              Zatvori
            </Button>
            <Button type="submit" variant="primary" disabled={submitting}>
              {submitting ? 'Slanje...' : 'Sačuvaj recenziju'}
            </Button>
          </div>
        </Form>
      </Modal.Body>
    </Modal>
  );
};

export default TourReviewsModal;
