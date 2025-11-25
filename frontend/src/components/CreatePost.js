import React, { useState } from "react";
import { Form, Button, Card, Alert, Container, Row, Col } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { blogApi } from "../services/api";
import { useAuth } from "../context/AuthContext";
import ReactMarkdown from "react-markdown";

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

  if (!user) {
    return (
      <Container className="py-5">
        <Alert variant="warning" className="text-center shadow-sm">
          <Alert.Heading>⚠️ Nemate dozvolu</Alert.Heading>
          <p>Morate se prijaviti da biste mogli da kreirate blog postove.</p>
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

  // helper: upload jedne slike
  const uploadSingleImage = async (file) => {
    const formDataImage = new FormData();
    formDataImage.append("image", file);

    try {
      const uploadRes = await blogApi.post("/uploads", formDataImage, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      // backend sada vraća FULL URL, npr:
      // "http://localhost:8000/api/blog/uploads/..."
      return uploadRes.data.url;
    } catch (err) {
      console.error("Upload error:", err);
      setError("Greška pri upload-u slike");
      return null;
    }
  };

  // uklanjanje pojedinačne slike iz liste
  const handleRemoveImage = (indexToRemove) => {
    setFormData((prev) => {
      const urls = prev.imageUrls
        .split("\n")
        .map((u) => u.trim())
        .filter((u) => u !== "");

      const filtered = urls.filter((_, idx) => idx !== indexToRemove);

      return {
        ...prev,
        imageUrls: filtered.join("\n"),
      };
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

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
              <div className="text-center mb-4">
                <h2 className="fw-bold text-primary">✍️ Kreiraj novi blog post</h2>
                <p className="text-muted mb-1">
                  Dodaj naslov, opis, sadržaj i slike
                </p>
                <small className="text-muted">
                  🧩 <strong>Napomena:</strong> Opis i sadržaj podržavaju{" "}
                  <code>Markdown</code> formatiranje (bold, naslovi, liste…)
                </small>
              </div>

              <Card className="mb-4 border-0 bg-light">
                <Card.Body className="py-3">
                  <h6 className="fw-bold mb-2">ℹ️ Kako da koristiš Markdown?</h6>
                  <ul className="mb-2 small text-muted">
                    <li><code>**tekst**</code> ➝ <strong>podebljano</strong></li>
                    <li><code>*tekst*</code> ➝ <em>iskošeno</em></li>
                    <li><code># Naslov 1</code>, <code>## Naslov 2</code> ➝ naslovi</li>
                    <li><code>- stavka</code> ➝ lista</li>
                    <li><code>[link](https://primer.com)</code> ➝ link</li>
                    <li><code>---</code> ➝ horizontalna linija</li>
                  </ul>
                  <p className="small mb-0 text-muted">
                    Sve što napišeš u polju za opis i sadržaj biće prikazano na sajtu
                    sa ovim formatiranjem.
                  </p>
                </Card.Body>
              </Card>

              {error && <Alert variant="danger">{error}</Alert>}

              <Form onSubmit={handleSubmit}>
                {/* Naslov */}
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

                {/* Opis (Markdown) */}
                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📌 Kratak opis (Markdown)</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={3}
                    name="description"
                    value={formData.description}
                    onChange={handleChange}
                    placeholder="Kratak opis, može **podebljano**, *iskošeno* ili [link](https://...)"
                    required
                    maxLength={500}
                  />
                  <Form.Text className="text-muted">
                    Ovaj opis će se pojaviti kao uvod na stranici bloga.
                  </Form.Text>

                  <Card className="mt-2">
                    <Card.Body className="py-2">
                      <small className="text-muted d-block mb-1">
                        🔍 Pregled opisa (Markdown)
                      </small>
                      <div className="small">
                        {formData.description.trim() ? (
                          <ReactMarkdown>{formData.description}</ReactMarkdown>
                        ) : (
                          <span className="text-muted">
                            Ovde će se prikazati kako opis izgleda sa Markdown formatiranjem.
                          </span>
                        )}
                      </div>
                    </Card.Body>
                  </Card>
                </Form.Group>

                {/* Sadržaj */}
                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">📖 Sadržaj (Markdown)</Form.Label>
                  <Form.Control
                    as="textarea"
                    rows={10}
                    name="content"
                    value={formData.content}
                    onChange={handleChange}
                    placeholder={
                      "# Dan 1\n" +
                      "Dolazak u grad...\n\n" +
                      "## Šta smo videli\n" +
                      "- Muzej\n- Trg\n- Park\n\n" +
                      "**Preporuka:** probajte lokalnu hranu!"
                    }
                    required
                    maxLength={5000}
                  />
                  <Form.Text className="text-muted">
                    Ovde pišeš ceo blog post. Naslovi, liste, bold i linkovi se pišu u Markdown formatu.
                  </Form.Text>

                  <Card className="mt-3">
                    <Card.Header className="py-2">
                      <small className="fw-bold">👀 Pregled sadržaja (Markdown)</small>
                    </Card.Header>
                    <Card.Body style={{ maxHeight: "300px", overflowY: "auto" }}>
                      {formData.content.trim() ? (
                        <ReactMarkdown>{formData.content}</ReactMarkdown>
                      ) : (
                        <p className="text-muted mb-0 small">
                          Kako budeš kucala tekst, ovde ćeš videti kako će izgledati na sajtu.
                        </p>
                      )}
                    </Card.Body>
                  </Card>
                </Form.Group>

                {/* Slike - Drag & Drop */}
                <Form.Group className="mb-4">
                  <Form.Label className="fw-bold">🖼️ Slike (prevuci ili izaberi)</Form.Label>

                  <div
                    onDragOver={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                    }}
                    onDrop={async (e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      const file = e.dataTransfer.files[0];
                      if (!file) return;

                      const url = await uploadSingleImage(file);
                      if (!url) return;

                      setFormData((prev) => ({
                        ...prev,
                        imageUrls:
                          prev.imageUrls + (prev.imageUrls ? "\n" : "") + url,
                      }));
                    }}
                    onClick={() => {
                      const input = document.createElement("input");
                      input.type = "file";
                      input.accept = "image/*";

                      input.onchange = async (e) => {
                        const file = e.target.files[0];
                        if (!file) return;

                        const url = await uploadSingleImage(file);
                        if (!url) return;

                        setFormData((prev) => ({
                          ...prev,
                          imageUrls:
                            prev.imageUrls + (prev.imageUrls ? "\n" : "") + url,
                        }));
                      };

                      input.click();
                    }}
                    style={{
                      border: "2px dashed #888",
                      borderRadius: "10px",
                      padding: "20px",
                      textAlign: "center",
                      cursor: "pointer",
                      background: "#fafafa",
                    }}
                  >
                    <p className="text-muted mb-0">
                      Prevuci sliku ili klikni da izabereš sa računara
                    </p>
                  </div>

                  {/* Preview uploadovanih slika sa X za brisanje */}
                  {formData.imageUrls && (
                    <div className="mt-3 d-flex flex-wrap gap-3">
                      {formData.imageUrls.split("\n").map((url, i) => {
                        const trimmed = url.trim();
                        if (!trimmed) return null;

                        return (
                          <div
                            key={i}
                            style={{
                              position: "relative",
                              width: "120px",
                              height: "120px",
                            }}
                          >
                            <img
                              src={trimmed} // backend već vratio full URL
                              alt="preview"
                              style={{
                                width: "100%",
                                height: "100%",
                                objectFit: "cover",
                                borderRadius: "10px",
                                border: "1px solid #ddd",
                              }}
                            />
                            <button
                              type="button"
                              onClick={() => handleRemoveImage(i)}
                              style={{
                                position: "absolute",
                                top: "4px",
                                right: "4px",
                                borderRadius: "50%",
                                border: "none",
                                width: "24px",
                                height: "24px",
                                fontSize: "14px",
                                lineHeight: "24px",
                                textAlign: "center",
                                background: "rgba(0,0,0,0.6)",
                                color: "#fff",
                                cursor: "pointer",
                              }}
                              title="Ukloni ovu sliku"
                            >
                              ×
                            </button>
                          </div>
                        );
                      })}
                    </div>
                  )}

                  <Form.Control
                    as="textarea"
                    rows={3}
                    className="mt-3"
                    name="imageUrls"
                    value={formData.imageUrls}
                    onChange={handleChange}
                    placeholder="Linkovi slika će se automatski dodati ovde nakon upload-a"
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
