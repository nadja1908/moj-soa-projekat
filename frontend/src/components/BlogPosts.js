import React, { useState, useEffect } from 'react';
import { Card, Badge, Button, Spinner, Alert } from 'react-bootstrap';
import { blogApi } from '../services/api';
import { useAuth } from '../context/AuthContext';

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
      const response = await blogApi.get('/posts');
      setPosts(response.data || []);
    } catch (error) {
      console.error('Error fetching posts:', error);
      setError('Greška pri učitavanju blog postova');
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
      // Refresh posts to get updated like count
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
      // Refresh posts to get updated like count
      fetchPosts();
    } catch (error) {
      console.error('Error unliking post:', error);
      setError('Greška pri uklanjanju označavanja');
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('sr-RS', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
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
            <p className="lead text-muted">Istražite najnovije blog postove o turističkim destinacijama</p>
          </div>
          {user && (user.role === 'guide' || user.role === 'administrator') && (
            <Button variant="primary" href="/create-post" size="lg">
              ✍️ Novi Post
            </Button>
          )}
        </div>

        {error && <Alert variant="danger" dismissible onClose={() => setError('')}>{error}</Alert>}

        {posts.length === 0 ? (
          <Card className="text-center py-5 border-0 shadow-sm">
            <Card.Body>
              <div className="feature-icon">📝</div>
              <h3>Nema dostupnih postova</h3>
              <p className="text-muted">Trenutno nema objavljenih blog postova.</p>
              {user && (user.role === 'guide' || user.role === 'administrator') && (
                <Button variant="primary" href="/create-post" className="mt-2">
                  Budi prvi koji će objaviti post!
                </Button>
              )}
            </Card.Body>
          </Card>
        ) : (
          <div className="blog-posts-grid">
            {posts.map((post) => (
              <Card key={post.id} className="blog-post-card border-0 shadow-sm">
                <Card.Body className="card-body-flex">
                  <div className="card-content">
                    <h5 className="card-title mb-3">{post.title}</h5>
                    <p className="card-text text-muted card-text-truncate">
                      {post.content}
                    </p>
                    
                    <div className="mt-auto">
                      <div className="d-flex justify-content-between align-items-center mb-3">
                        <small className="text-muted">
                          📅 {formatDate(post.createdAt)}
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
                            className={`like-button me-2 ${post.isLiked ? 'liked' : ''}`}
                            onClick={() => post.isLiked ? handleUnlike(post.id) : handleLike(post.id)}
                            disabled={!user}
                            title={user ? (post.isLiked ? 'Ukloni označavanje' : 'Označi post') : 'Prijavite se da označite post'}
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
            ))}
          </div>
        )}

        {!user && (
          <Alert variant="info" className="mt-5 border-0 shadow-sm">
            <Alert.Heading>💡 Tip</Alert.Heading>
            <p>
              Prijavite se da biste mogli da označavate postove kao omiljene i ostavljate komentare!
            </p>
          </Alert>
        )}
      </div>
    </div>
  );
};

export default BlogPosts;