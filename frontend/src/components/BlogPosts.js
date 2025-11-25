import React, { useState, useEffect } from 'react';
import { Card, Badge, Button, Spinner, Alert } from 'react-bootstrap';
import { blogApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import ReactMarkdown from 'react-markdown';

const BlogPosts = () => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await blogApi.get('/posts');
      console.log('BlogPosts API response:', response);

      const postsData = response.data;
      if (Array.isArray(postsData)) {
        setPosts(postsData);
      } else if (postsData && postsData.posts && Array.isArray(postsData.posts)) {
        setPosts(postsData.posts);
      } else {
        console.warn('API did not return array, setting empty array');
        setPosts([]);
      }
    } catch (error) {
      console.error('Error fetching posts:', error);
      setError(
        'Greška pri učitavanju blog postova. Proverite da li je blog servis pokrenut.'
      );
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  const handleLike = async (postId) => {
    if (!user) {
      setError('Morate biti prijavljeni da biste označili post');
      return;
    }

    try {
      await blogApi.post(`/posts/${postId}/like`);
      fetchPosts();
    } catch (error) {
      console.error('Error liking post:', error);
      setError('Greška pri označavanju posta');
    }
  };

  const handleUnlike = async (postId) => {
    if (!user) {
      return;
    }

    try {
      await blogApi.delete(`/posts/${postId}/like`);
      fetchPosts();
    } catch (error) {
      console.error('Error unliking post:', error);
      setError('Greška pri uklanjanju označavanja');
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'Nepoznat datum';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'Nepoznat datum';

    return date.toLocaleDateString('sr-RS', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  //  ŠtA god backend pošalje, pokušavamo da izvučemo prvi URL slike
  const getFirstImageUrl = (post) => {
    if (!post) return null;

    const candidates = [];

    // 1) images: ["http://...", "http://..."]
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

    // 2) imageUrls: "url1\nurl2" – ako backend tako šalje
    if (typeof post.imageUrls === 'string') {
      const parts = post.imageUrls
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean);
      candidates.push(...parts);
    }

    // 3) još jedna varijanta: post.blogImages = [{url: ...}, ...]
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

  if (loading) {
    return (
      <div className="text-center py-5">
        <Spinner animation="border" variant="primary" />
        <p className="mt-2 text-muted">Učitavanje blog postova...</p>
      </div>
    );
  }

  return (
    <div className="blog-posts-container">
      <div className="container">
        <div className="d-flex justify-content-between align-items-center mb-5">
          <div>
            <h1 className="display-5">📰 Blog Postovi</h1>
            <p className="lead text-muted">
              Istražite najnovije blog postove o turističkim destinacijama
            </p>
          </div>
          {user && (
            <Button variant="primary" href="/create-post" size="lg">
              ✍️ Novi Post
            </Button>
          )}
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
              <h3>Nema dostupnih blogova</h3>
              <p className="text-muted">Trenutno nema objavljenih blog postova.</p>
              {user && (
                <Button variant="primary" href="/create-post" className="mt-2">
                  ✍️ Dodaj novi blog
                </Button>
              )}
              {!user && (
                <p className="text-muted mt-2">
                  <small>Prijavite se da biste mogli da kreirate blog postove</small>
                </p>
              )}
            </Card.Body>
          </Card>
        ) : (
          <div className="blog-posts-grid">
            {Array.isArray(posts) &&
              posts.map((post) => {
                const firstImageUrl = getFirstImageUrl(post);
                const createdAt =
                  post.createdAt || post.created_at || post.CreatedAt || null;

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

                        {/* opis u Markdown-u */}
                        {post.description && (
                          <div className="text-muted small mb-2 markdown-preview">
                            <ReactMarkdown>{post.description}</ReactMarkdown>
                          </div>
                        )}

                        {/* sadržaj u Markdown-u (skraćeno) */}
                        {post.content && (
                          <div className="card-text text-muted small mb-3 markdown-preview">
                            <ReactMarkdown>
                              {post.content.length > 400
                                ? post.content.slice(0, 400) + '...'
                                : post.content}
                            </ReactMarkdown>
                          </div>
                        )}

                        <div className="mt-auto">
                          <div className="d-flex justify-content-between align-items-center mb-2">
                            <small className="text-muted">
                              📅 {formatDate(createdAt)}
                            </small>
                            <Badge bg="info" className="role-badge">
                              {post.author?.username || 'Nepoznat autor'}
                            </Badge>
                          </div>

                          <div className="d-flex justify-content-between align-items-center">
                            <small className="text-muted">
                              💬 {post.comments?.length || 0} komentara
                            </small>
                            <div className="d-flex align-items-center">
                              <button
                                className={`like-button me-2 ${
                                  post.isLiked ? 'liked' : ''
                                }`}
                                onClick={() =>
                                  post.isLiked
                                    ? handleUnlike(post.id)
                                    : handleLike(post.id)
                                }
                                disabled={!user}
                                title={
                                  user
                                    ? post.isLiked
                                      ? 'Ukloni označavanje'
                                      : 'Označi post'
                                    : 'Prijavite se da označite post'
                                }
                              >
                                ❤️
                              </button>
                              <small className="text-muted">
                                {post.likes || 0}
                              </small>
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

        {!user && (
          <Alert variant="info" className="mt-5 border-0 shadow-sm">
            <Alert.Heading>💡 Tip</Alert.Heading>
            <p>
              Prijavite se da biste mogli da označavate postove kao omiljene i ostavljate
              komentare!
            </p>
          </Alert>
        )}
      </div>
    </div>
  );
};

export default BlogPosts;
