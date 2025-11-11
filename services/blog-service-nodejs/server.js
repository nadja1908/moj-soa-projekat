const express = require('express');
const mysql = require('mysql2/promise');
const cors = require('cors');
const axios = require('axios');

const app = express();
const PORT = process.env.PORT || 8002;

// Database configuration
const dbConfig = {
  host: process.env.DB_HOST || 'localhost',
  user: process.env.DB_USER || 'user', 
  password: process.env.DB_PASS || 'password',
  database: process.env.DB_NAME || 'blog_db',
  waitForConnections: true,
  connectionLimit: 10,
  queueLimit: 0
};

let pool;

// Initialize database connection
async function initDatabase() {
  try {
    pool = mysql.createPool(dbConfig);
    console.log('Connected to MySQL database');
  } catch (error) {
    console.error('Database connection failed:', error);
    process.exit(1);
  }
}

// Middleware
app.use(cors());
app.use(express.json());

// Auth middleware
async function authMiddleware(req, res, next) {
  const authHeader = req.headers.authorization;
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Authorization header required' });
  }

  const token = authHeader.substring(7);
  
  try {
    // Validate token with auth service
    const authResponse = await axios.get(`${process.env.AUTH_SERVICE_URL || 'http://auth-service:8003'}/validate`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    req.user = {
      userID: authResponse.data.userId,
      role: authResponse.data.role
    };
    
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', service: 'blog-service-nodejs' });
});

// Get all blog posts (public)
app.get('/posts', async (req, res) => {
  try {
    const [rows] = await pool.execute(`
      SELECT id, user_id, title, description, content, created_at, updated_at
      FROM blog_posts
      ORDER BY created_at DESC
    `);
    
    res.json({ posts: rows });
  } catch (error) {
    console.error('Error fetching blog posts:', error);
    res.status(500).json({ error: 'Failed to get blog posts' });
  }
});

// Get single blog post with comments (public)
app.get('/posts/:id', async (req, res) => {
  const blogId = parseInt(req.params.id);
  
  if (isNaN(blogId)) {
    return res.status(400).json({ error: 'Invalid blog ID' });
  }

  try {
    // Get blog post
    const [postRows] = await pool.execute(`
      SELECT id, user_id, title, description, content, created_at, updated_at
      FROM blog_posts
      WHERE id = ?
    `, [blogId]);

    if (postRows.length === 0) {
      return res.status(404).json({ error: 'Blog post not found' });
    }

    // Get comments
    const [commentRows] = await pool.execute(`
      SELECT id, blog_id, user_id, content, created_at, updated_at
      FROM comments
      WHERE blog_id = ? 
      ORDER BY created_at ASC
    `, [blogId]);

    res.json({ 
      post: postRows[0], 
      comments: commentRows 
    });
  } catch (error) {
    console.error('Error fetching blog post:', error);
    res.status(500).json({ error: 'Failed to get blog post' });
  }
});

// Create new blog post (protected)
app.post('/posts', authMiddleware, async (req, res) => {
  const { title, description, content } = req.body;
  
  if (!title || !content) {
    return res.status(400).json({ error: 'Title and content are required' });
  }

  try {
    const [result] = await pool.execute(`
      INSERT INTO blog_posts (user_id, title, description, content, created_at, updated_at) 
      VALUES (?, ?, ?, ?, NOW(), NOW())
    `, [req.user.userID, title, description || '', content]);

    const [newPost] = await pool.execute(`
      SELECT id, user_id, title, description, content, created_at, updated_at
      FROM blog_posts
      WHERE id = ?
    `, [result.insertId]);

    res.status(201).json({ 
      message: 'Blog post created successfully',
      post: newPost[0]
    });
  } catch (error) {
    console.error('Error creating blog post:', error);
    res.status(500).json({ error: 'Failed to create blog post' });
  }
});

// Create comment on blog post (protected)
app.post('/posts/:id/comments', authMiddleware, async (req, res) => {
  const blogId = parseInt(req.params.id);
  const { content } = req.body;
  
  if (isNaN(blogId)) {
    return res.status(400).json({ error: 'Invalid blog ID' });
  }
  
  if (!content) {
    return res.status(400).json({ error: 'Comment content is required' });
  }

  try {
    // Check if blog post exists
    const [blogCheck] = await pool.execute('SELECT id FROM blog_posts WHERE id = ?', [blogId]);
    if (blogCheck.length === 0) {
      return res.status(404).json({ error: 'Blog post not found' });
    }

    const [result] = await pool.execute(`
      INSERT INTO comments (blog_id, user_id, content, created_at, updated_at) 
      VALUES (?, ?, ?, NOW(), NOW())
    `, [blogId, req.user.userID, content]);

    const [newComment] = await pool.execute(`
      SELECT id, blog_id, user_id, content, created_at, updated_at
      FROM comments
      WHERE id = ?
    `, [result.insertId]);

    res.status(201).json({ 
      message: 'Comment created successfully',
      comment: newComment[0]
    });
  } catch (error) {
    console.error('Error creating comment:', error);
    res.status(500).json({ error: 'Failed to create comment' });
  }
});

// Like blog post (protected)
app.post('/posts/:id/like', authMiddleware, async (req, res) => {
  const blogId = parseInt(req.params.id);
  
  if (isNaN(blogId)) {
    return res.status(400).json({ error: 'Invalid blog ID' });
  }

  try {
    // Check if already liked
    const [existing] = await pool.execute(`
      SELECT id FROM likes WHERE blog_id = ? AND user_id = ?
    `, [blogId, req.user.userID]);

    if (existing.length > 0) {
      return res.status(400).json({ error: 'Post already liked' });
    }

    await pool.execute(`
      INSERT INTO likes (blog_id, user_id, created_at) 
      VALUES (?, ?, NOW())
    `, [blogId, req.user.userID]);

    res.json({ message: 'Post liked successfully' });
  } catch (error) {
    console.error('Error liking post:', error);
    res.status(500).json({ error: 'Failed to like post' });
  }
});

// Unlike blog post (protected)
app.delete('/posts/:id/like', authMiddleware, async (req, res) => {
  const blogId = parseInt(req.params.id);
  
  if (isNaN(blogId)) {
    return res.status(400).json({ error: 'Invalid blog ID' });
  }

  try {
    const [result] = await pool.execute(`
      DELETE FROM likes WHERE blog_id = ? AND user_id = ?
    `, [blogId, req.user.userID]);

    if (result.affectedRows === 0) {
      return res.status(400).json({ error: 'Post not liked' });
    }

    res.json({ message: 'Post unliked successfully' });
  } catch (error) {
    console.error('Error unliking post:', error);
    res.status(500).json({ error: 'Failed to unlike post' });
  }
});

// Start server
async function startServer() {
  try {
    await initDatabase();
    
    app.listen(PORT, () => {
      console.log(`Blog service (Node.js) starting on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

// Graceful shutdown
process.on('SIGINT', async () => {
  console.log('Shutting down gracefully...');
  if (pool) {
    await pool.end();
  }
  process.exit(0);
});

startServer();