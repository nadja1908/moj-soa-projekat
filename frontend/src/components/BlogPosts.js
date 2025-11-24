import React, { useState, useEffect } from 'react';
import { Card, Badge, Button, Spinner, Alert, Form, ListGroup } from 'react-bootstrap';
import { blogApi, followerApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import ReactMarkdown from 'react-markdown';

const BlogPosts = () => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [followingUsers, setFollowingUsers] = useState(new Set());
  const [expandedPost, setExpandedPost] = useState(null);
  const [comments, setComments] = useState({});
  const [newComment, setNewComment] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    const loadData = async () => {
      if (user) {
        // Get following users first
        const followingSet = await fetchFollowingUsers();
        // Then fetch and filter posts with the following users
        await fetchPosts(followingSet);
      } else {
        await fetchPosts(new Set());
      }
    };
    loadData();
  }, [user]);

  const fetchFollowingUsers = async () => {
    try {
      const response = await followerApi.get('/following');
      const followingSet = new Set(response.data);
      setFollowingUsers(followingSet);
      return followingSet;
    } catch (error) {
      console.error('Error fetching following users:', error);
      const emptySet = new Set();
      setFollowingUsers(emptySet);
      return emptySet;
    }
  };

  const fetchPosts = async (followingSet) => {
    try {
      setLoading(true);
      setError('');
      const response = await blogApi.get('/posts');
      console.log('BlogPosts API response:', response);

      const postsData = response.data;
      let allPosts = [];
      
      if (Array.isArray(postsData)) {
        allPosts = postsData;
      } else if (postsData && postsData.posts && Array.isArray(postsData.posts)) {
        allPosts = postsData.posts;
      } else {
        console.warn('API did not return array, setting empty array');
        allPosts = [];
      }

      // Filtriraj postove: prikaži samo postove od korisnika koje pratiš (bez tvojih postova)
      if (user) {
        console.log('Filtering posts. Following users:', Array.from(followingSet));
        console.log('User ID:', user.id);
        const filteredPosts = allPosts.filter(post => {
          // NE prikazuj svoj post
          if (post.author && post.author.id === user.id) {
            console.log(`Excluding own post: ${post.id}`);
            return false;
          }
          // Prikaži post od korisnika koje pratiš
          if (post.author && followingSet.has(post.author.id)) {
            console.log(`Including post ${post.id} from followed user ${post.author.id}`);
            return true;
          }
          console.log(`Excluding post ${post.id} from unfollowed user ${post.author.id}`);
          return false;
        });
        console.log('Filtered posts count:', filteredPosts.length);
        setPosts(filteredPosts);
        // Učitaj broj komentara za sve postove
        await fetchAllCommentsCount(filteredPosts);
      } else {
        setPosts(allPosts);
        await fetchAllCommentsCount(allPosts);
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
      fetchPosts(followingUsers);
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
      fetchPosts(followingUsers);
    } catch (error) {
      console.error('Error unliking post:', error);
      setError('Greška pri uklanjanju označavanja');
    }
  };

  const handleFollow = async (userId) => {
    try {
      console.log('Following user:', userId);
      const response = await followerApi.post('/follow', { followingId: userId });
      console.log('Follow response:', response);
      const updatedFollowing = await fetchFollowingUsers();
      await fetchPosts(updatedFollowing);
      setError('');
    } catch (error) {
      console.error('Error following user:', error);
      console.error('Error details:', error.response?.data);
      setError(error.response?.data?.message || 'Greška pri praćenju korisnika');
    }
  };

  const handleUnfollow = async (userId) => {
    try {
      console.log('Unfollowing user:', userId);
      const response = await followerApi.delete(`/unfollow/${userId}`);
      console.log('Unfollow response:', response);
      const updatedFollowing = await fetchFollowingUsers();
      console.log('Updated following:', Array.from(updatedFollowing));
      await fetchPosts(updatedFollowing);
      setError('');
    } catch (error) {
      console.error('Error unfollowing user:', error);
      console.error('Error response:', error.response);
      console.error('Error details:', error.response?.data);
      setError(error.response?.data?.message || 'Greška pri otpraćivanju korisnika');
    }
  };

  const fetchComments = async (postId) => {
    try {
      const response = await blogApi.get(`/posts/${postId}`);
      if (response.data && response.data.comments) {
        setComments(prev => ({ ...prev, [postId]: response.data.comments }));
      }
    } catch (error) {
      console.error('Error fetching comments:', error);
    }
  };

  const fetchAllCommentsCount = async (postsArray) => {
    const commentsData = {};
    for (const post of postsArray) {
      try {
        const response = await blogApi.get(`/posts/${post.id}`);
        if (response.data && response.data.comments) {
          commentsData[post.id] = response.data.comments;
        }
      } catch (error) {
        console.error(`Error fetching comments for post ${post.id}:`, error);
      }
    }
    setComments(commentsData);
  };

  const handleAddComment = async (postId) => {
    if (!newComment.trim()) return;

    try {
      await blogApi.post(`/posts/${postId}/comments`, {
        commentText: newComment
      });
      setNewComment('');
      await fetchComments(postId);
      setError('');
    } catch (error) {
      console.error('Error adding comment:', error);
      if (error.response?.status === 403) {
        setError('Morate pratiti autora da biste mogli komentarisati njegov blog');
      } else {
        setError('Greška pri dodavanju komentara');
      }
    }
  };

  const togglePostDetails = async (postId) => {
    if (expandedPost === postId) {
      setExpandedPost(null);
    } else {
      setExpandedPost(postId);
      await fetchComments(postId);
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
                            <div className="d-flex align-items-center gap-2">
                              <Badge bg="info" className="role-badge">
                                {post.author?.username || 'Nepoznat autor'}
                              </Badge>
                              {user && post.author && post.author.id !== user.id && (
                                <Button
                                  size="sm"
                                  variant={followingUsers.has(post.author.id) ? 'secondary' : 'primary'}
                                  onClick={() =>
                                    followingUsers.has(post.author.id)
                                      ? handleUnfollow(post.author.id)
                                      : handleFollow(post.author.id)
                                  }
                                >
                                  {followingUsers.has(post.author.id) ? 'Otprati' : 'Prati'}
                                </Button>
                              )}
                            </div>
                          </div>

                          <div className="d-flex justify-content-between align-items-center mb-2">
                            <Button
                              size="sm"
                              variant="link"
                              className="p-0 text-decoration-none"
                              onClick={() => togglePostDetails(post.id)}
                            >
                              💬 {comments[post.id]?.length || 0} komentara
                            </Button>
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
                                      ? 'Ukloni lajk'
                                      : 'Lajkuj post'
                                    : 'Prijavite se da lajkujete post'
                                }
                                style={{
                                  background: 'none',
                                  border: 'none',
                                  cursor: user ? 'pointer' : 'not-allowed',
                                  fontSize: '1.5rem',
                                  padding: 0,
                                  transition: 'transform 0.2s'
                                }}
                                onMouseEnter={(e) => e.currentTarget.style.transform = 'scale(1.2)'}
                                onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
                              >
                                {post.isLiked ? '❤️' : '🤍'}
                              </button>
                              <small className="text-muted">
                                {post.likesCount || 0}
                              </small>
                            </div>
                          </div>

                          {/* Comments Section */}
                          {expandedPost === post.id && (
                            <div className="mt-3 border-top pt-3">
                              <h6 className="mb-3">Komentari</h6>
                              
                              {/* Add Comment Form */}
                              {user && (
                                <Form className="mb-3">
                                  <Form.Group>
                                    <Form.Control
                                      as="textarea"
                                      rows={2}
                                      placeholder="Napišite komentar..."
                                      value={newComment}
                                      onChange={(e) => setNewComment(e.target.value)}
                                    />
                                  </Form.Group>
                                  <Button
                                    size="sm"
                                    variant="primary"
                                    className="mt-2"
                                    onClick={() => handleAddComment(post.id)}
                                    disabled={!newComment.trim()}
                                  >
                                    Dodaj komentar
                                  </Button>
                                </Form>
                              )}

                              {/* Comments List */}
                              {comments[post.id] && comments[post.id].length > 0 ? (
                                <ListGroup variant="flush">
                                  {comments[post.id].map((comment) => (
                                    <ListGroup.Item key={comment.id} className="px-0">
                                      <div className="d-flex justify-content-between">
                                        <strong className="text-primary">Korisnik {comment.userId}</strong>
                                        <small className="text-muted">
                                          {formatDate(comment.createdAt)}
                                        </small>
                                      </div>
                                      <p className="mb-0 mt-1">{comment.commentText}</p>
                                    </ListGroup.Item>
                                  ))}
                                </ListGroup>
                              ) : (
                                <p className="text-muted text-center">Nema komentara</p>
                              )}
                            </div>
                          )}
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
