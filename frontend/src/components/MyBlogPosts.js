import React, { useState, useEffect } from 'react';
import { Card, Badge, Button, Spinner, Alert } from 'react-bootstrap';
import { blogApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import ReactMarkdown from 'react-markdown';

const MyBlogPosts = () => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      fetchMyPosts();
    }
  }, [user]);

  const fetchMyPosts = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await blogApi.get('/posts');
      
      const postsData = Array.isArray(response.data) ? response.data : [];
      
      // Filter only my posts
      const myPosts = postsData.filter(post => post.author && post.author.id === user.id);
      setPosts(myPosts);
    } catch (error) {
      console.error('Error fetching my posts:', error);
      setError('Greška pri učitavanju vaših blog postova.');
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (postId) => {
    if (!window.confirm('Da li ste sigurni da želite da obrišete ovaj post?')) {
      return;
    }

    try {
      await blogApi.delete(`/posts/${postId}`);
      fetchMyPosts();
    } catch (error) {
      console.error('Error deleting post:', error);
      setError('Greška pri brisanju posta');
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('sr-RS', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const getFirstImageUrl = (post) => {
    if (!post) return null;

    const candidates = [];

    if (Array.isArray(post.images)) {
      for (const img of post.images) {
        if (typeof img === 'string') {
          candidates.push(img);
        } else if (img && typeof img === 'object') {
          if (img.imageUrl) candidates.push(img.imageUrl);
          if (img.url) candidates.push(img.url);
        }
      }
    }

    if (typeof post.imageUrls === 'string') {
      const parts = post.imageUrls
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean);
      candidates.push(...parts);
    }

    if (Array.isArray(post.blogImages)) {
      for (const img of post.blogImages) {
        if (!img) continue;
        if (img.imageUrl) candidates.push(img.imageUrl);
        if (img.url) candidates.push(img.url);
      }
    }

    const first = candidates.find((u) => !!u);
    return first || null;
  };

  if (!user) {
    return (
      <div className="container mt-5">
        <Alert variant="warning">
          Morate biti prijavljeni da biste videli svoje blogove.
        </Alert>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="text-center py-5">
        <Spinner animation="border" variant="primary" />
        <p className="mt-2 text-muted">Učitavanje vaših blog postova...</p>
      </div>
    );
  }

  return (
    <div className="blog-posts-container">
      <div className="container">
        <div className="d-flex justify-content-between align-items-center mb-5">
          <div>
            <h1 className="display-5">📝 Moji Blogovi</h1>
            <p className="lead text-muted">
              Pregledajte i upravljajte svojim blog postovima
            </p>
          </div>
          <Button variant="primary" href="/create-post" size="lg">
            ✍️ Novi Post
          </Button>
        </div>

        {error && (
          <Alert variant="danger" dismissible onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        {posts.length === 0 ? (
          <Card className="text-center py-5 border-0 shadow-sm">
            <Card.Body>
              <div className="feature-icon">📝</div>
              <h3>Nemate objavljenih blogova</h3>
              <p className="text-muted">Kreirajte svoj prvi blog post i podelite ga sa svetom!</p>
              <Button variant="primary" href="/create-post" className="mt-2">
                ✍️ Kreiraj novi blog
              </Button>
            </Card.Body>
          </Card>
        ) : (
          <div className="blog-posts-grid">
            {posts.map((post) => {
              const firstImageUrl = getFirstImageUrl(post);
              const createdAt = post.createdAt || post.created_at || post.CreatedAt || null;

              return (
                <Card key={post.id} className="blog-post-card border-0 shadow-sm">
                  {firstImageUrl && (
                    <Card.Img
                      variant="top"
                      src={firstImageUrl}
                      alt={post.title}
                      style={{
                        maxHeight: '220px',
                        objectFit: 'cover'
                      }}
                    />
                  )}

                  <Card.Body className="card-body-flex">
                    <div className="card-content">
                      <h5 className="card-title mb-2">{post.title}</h5>

                      {post.description && (
                        <div className="text-muted small mb-2 markdown-preview">
                          <ReactMarkdown>{post.description}</ReactMarkdown>
                        </div>
                      )}

                      {post.content && (
                        <div className="card-text text-muted small mb-3 markdown-preview">
                          <ReactMarkdown>
                            {post.content.length > 200
                              ? post.content.slice(0, 200) + '...'
                              : post.content}
                          </ReactMarkdown>
                        </div>
                      )}

                      <div className="mt-auto">
                        <div className="d-flex justify-content-between align-items-center mb-2">
                          <small className="text-muted">
                            📅 {formatDate(createdAt)}
                          </small>
                          <Badge bg="success" className="role-badge">
                            Tvoj post
                          </Badge>
                        </div>

                        <div className="d-flex justify-content-between align-items-center">
                          <div className="d-flex gap-2">
                            <span className="text-muted small">
                              ❤️ {post.likesCount || 0} lajkova
                            </span>
                          </div>
                          <div className="d-flex gap-2">
                            <Button
                              size="sm"
                              variant="outline-danger"
                              onClick={() => handleDelete(post.id)}
                            >
                              🗑️ Obriši
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </Card.Body>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      <style jsx="true">{`
        .blog-posts-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
          gap: 2rem;
          margin-bottom: 3rem;
        }

        .blog-post-card {
          transition: transform 0.2s, box-shadow 0.2s;
          height: 100%;
          display: flex;
          flex-direction: column;
        }

        .blog-post-card:hover {
          transform: translateY(-5px);
          box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15) !important;
        }

        .card-body-flex {
          display: flex;
          flex-direction: column;
          flex: 1;
        }

        .card-content {
          display: flex;
          flex-direction: column;
          flex: 1;
        }

        .feature-icon {
          font-size: 4rem;
          margin-bottom: 1rem;
        }

        .markdown-preview {
          line-height: 1.6;
        }

        .markdown-preview p {
          margin-bottom: 0.5rem;
        }

        .markdown-preview h1,
        .markdown-preview h2,
        .markdown-preview h3 {
          margin-top: 1rem;
          margin-bottom: 0.5rem;
        }

        .role-badge {
          font-size: 0.75rem;
          padding: 0.35rem 0.65rem;
        }
      `}</style>
    </div>
  );
};

export default MyBlogPosts;
