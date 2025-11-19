import React, { useState } from "react";
import { Form, Button, Card, Alert, Container, Row, Col } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { blogApi } from "../services/api";
import { useAuth } from "../context/AuthContext";

const CreatePost = () => {
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    content: "",
    imageUrls: ""
  });

  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { user } = useAuth();

  // samo vodič i admin
  if (!user || (user.role !== "guide" && user.role !== "administrator")) {
    return (
      <Container className="py-5">
        <Alert variant="warning" className="text-center shadow-sm">
          <Alert.Heading>⚠️ Nemate dozvolu</Alert.Heading>
          <p>Samo vodiči i administratori mogu kreirati blog postove.</p>
        </Alert>
      </Container>
    );
  }

  const handleChange = (e) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    // Validacija
    if (formData.title.length < 5) {
      setError("Naslov mora imati najmanje 5 karaktera.");
      setLoading(false);
      return;
    }
    if (formData.description.length < 10) {
      setError("Opis mora imati najmanje 10 karaktera.");
      setLoading(false);
      return;
    }
    if (formData.content.length < 20) {
      setError("Sadržaj mora imati najmanje 20 karaktera.");
      setLoading(false);
      return;
    }

    const images = formData.imageUrls
      .split("\n")
      .map((url) => url.trim())
      .filter((u) => u !== "");

    const payload = {
      title: formData.title,
      description: formData.description,
      content: formData.content,
      images
    };

    try {
      const token = localStorage.getItem("token");
      console.log("CreatePost token:", token);
      console.log("CreatePost payload:", payload);

      if (!token) {
        setError("Niste autentifikovani. Pokušajte ponovo da se prijavite.");
        setLoading(false);
        return;
      }

      // NE šaljemo ručno headers, blogApi interceptor već dodaje Authorization
      const res = await blogApi.post("/posts", payload);

      console.log("USPESNO!", res.data);
      navigate("/posts");
    } catch (err) {
      console.error("CreatePost error:", err);

      const status = err.response?.status;
      const serverMsg = err.response?.data?.error || err.response?.data?.message;

      setError(
        `Greška pri kreiranju posta. ` +
        (status ? `Status: ${status}. ` : "") +
        (serverMsg ? `Poruka servera: ${serverMsg}` : "")
      );
    }

    setLoading(false);
  };

  return (
    <Container className="py-5">
      <Row className="justify-content-center">
        <Col lg={10} xl={8}>
          <Card className="border-0 shadow-lg rounded-4">
            <Card.Body className="p-5">
              <div className="text-center mb-5">
                <h2 className="fw-bold text-primary">✍️ Kreiraj novi blog post</h2>
                <p className="text-muted">Dodaj naslov, opis, sadržaj i slike</p>
              </div>

              {error && <Alert variant="danger">{error}</Alert>}

              <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📝 Naslov</Form.Label>
                  <Form.Control
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleChange}
                    maxLength={200}
                    placeholder="Npr. Putovanje u Rim"
                    required
                  />
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📌 Kratak opis (Markdown)</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={3}
                    name="description"
                    value={formData.description}
                    onChange={handleChange}
                    placeholder="Kratak opis, može **markdown**"
                    required
                    maxLength={500}
                  />
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📖 Sadržaj (Markdown)</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={10}
                    name="content"
                    value={formData.content}
                    onChange={handleChange}
                    placeholder="# Dan 1\nPoseta Koloseumu..."
                    required
                    maxLength={5000}
                  />
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">🖼️ Slike (URL – po jedna u redu)</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={3}
                    name="imageUrls"
                    value={formData.imageUrls}
                    onChange={handleChange}
                    placeholder={
                      "https://example.com/slika1.jpg\nhttps://example.com/slika2.png"
                    }
                  />
                </Form.Group>

                <div className="d-flex gap-3">
                  <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    className="flex-grow-1"
                  >
                    {loading ? "Objavljivanje..." : "🚀 Objavi blog"}
                  </Button>

                  <Button
                    variant="outline-secondary"
                    size="lg"
                    onClick={() => navigate("/posts")}
                  >
                    Odustani
                  </Button>
                </div>
              </Form>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
};

export default CreatePost;
