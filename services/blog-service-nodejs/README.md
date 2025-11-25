# Node.js Blog Service

This is the Node.js implementation of the blog microservice, providing the same functionality as the Go version.

## Features

- **Public Endpoints:**
  - `GET /posts` - Get all blog posts
  - `GET /posts/:id` - Get single blog post with comments
  - `GET /health` - Health check

- **Protected Endpoints (requires authentication):**
  - `POST /posts` - Create new blog post
  - `POST /posts/:id/comments` - Add comment to blog post
  - `POST /posts/:id/like` - Like a blog post
  - `DELETE /posts/:id/like` - Unlike a blog post

## Environment Variables

- `PORT` - Server port (default: 8002)
- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database username (default: user)
- `DB_PASS` - Database password (default: password)
- `DB_NAME` - Database name (default: blog_db)
- `AUTH_SERVICE_URL` - Auth service URL (default: http://auth-service:8003)

## Running

```bash
npm install
npm start
```

## Development

```bash
npm run dev
```