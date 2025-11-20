package store

import (
	"database/sql"
	"fmt"

	"blog-service/internal/model"
)

// CreateBlogPost kreira novi blog post
func (s *Store) CreateBlogPost(post *model.BlogPost) error {
	query := `INSERT INTO blog_posts (user_id, title, description, content) 
			  VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, post.UserID, post.Title, post.Description, post.Content)
	if err != nil {
		return fmt.Errorf("failed to create blog post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	post.ID = id
	return nil
}

// GetAllBlogPosts vraća sve blog postove
func (s *Store) GetAllBlogPosts() ([]model.BlogPost, error) {
	query := `SELECT bp.id, bp.user_id, bp.title, bp.description, bp.content, 
			  bp.created_at, bp.updated_at,
			  COALESCE(COUNT(bl.id), 0) as likes_count
			  FROM blog_posts bp
			  LEFT JOIN blog_likes bl ON bp.id = bl.blog_id
			  GROUP BY bp.id
			  ORDER BY bp.created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog posts: %w", err)
	}
	defer rows.Close()

	var posts []model.BlogPost
	for rows.Next() {
		var post model.BlogPost
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Description,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.LikesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blog post: %w", err)
		}

		//  DODATO: dohvati slike za svaki post
		if images, err := s.GetImagesByBlogID(post.ID); err == nil {
			post.Images = images
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// GetBlogPostByID vraća blog post po ID-u
func (s *Store) GetBlogPostByID(id int64) (*model.BlogPost, error) {
	query := `SELECT bp.id, bp.user_id, bp.title, bp.description, bp.content, 
			  bp.created_at, bp.updated_at,
			  COALESCE(COUNT(bl.id), 0) as likes_count
			  FROM blog_posts bp
			  LEFT JOIN blog_likes bl ON bp.id = bl.blog_id
			  WHERE bp.id = ?
			  GROUP BY bp.id`

	var post model.BlogPost
	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Description,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.LikesCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	//  DODATO: slike i za jedan konkretan post
	if images, err := s.GetImagesByBlogID(post.ID); err == nil {
		post.Images = images
	}

	return &post, nil
}

// CreateComment kreira novi komentar
func (s *Store) CreateComment(comment *model.BlogComment) error {
	query := `INSERT INTO blog_comments (blog_id, user_id, comment_text) 
			  VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, comment.BlogID, comment.UserID, comment.CommentText)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	comment.ID = id
	return nil
}

// GetCommentsByBlogID vraća komentare za određeni blog post
func (s *Store) GetCommentsByBlogID(blogID int64) ([]model.BlogComment, error) {
	query := `SELECT id, blog_id, user_id, comment_text, created_at, updated_at 
			  FROM blog_comments 
			  WHERE blog_id = ? 
			  ORDER BY created_at ASC`

	rows, err := s.db.Query(query, blogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []model.BlogComment
	for rows.Next() {
		var comment model.BlogComment
		err := rows.Scan(
			&comment.ID,
			&comment.BlogID,
			&comment.UserID,
			&comment.CommentText,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// LikeBlogPost dodaje lajk na blog post
func (s *Store) LikeBlogPost(userID, blogID int64) error {
	query := `INSERT IGNORE INTO blog_likes (blog_id, user_id) VALUES (?, ?)`
	_, err := s.db.Exec(query, blogID, userID)
	if err != nil {
		return fmt.Errorf("failed to like blog post: %w", err)
	}
	return nil
}

// UnlikeBlogPost uklanja lajk sa blog posta
func (s *Store) UnlikeBlogPost(userID, blogID int64) error {
	query := `DELETE FROM blog_likes WHERE blog_id = ? AND user_id = ?`
	_, err := s.db.Exec(query, blogID, userID)
	if err != nil {
		return fmt.Errorf("failed to unlike blog post: %w", err)
	}
	return nil
}

// IsPostLikedByUser proverava da li je korisnik lajkovao post
func (s *Store) IsPostLikedByUser(userID, blogID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM blog_likes WHERE blog_id = ? AND user_id = ?`
	var count int
	err := s.db.QueryRow(query, blogID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if post is liked: %w", err)
	}
	return count > 0, nil
}

func (s *Store) CreateBlogImages(blogID int64, imageURLs []string) error {
	if len(imageURLs) == 0 {
		return nil
	}

	query := `INSERT INTO blog_images (blog_id, image_url) VALUES (?, ?)`

	for _, url := range imageURLs {
		_, err := s.db.Exec(query, blogID, url)
		if err != nil {
			return fmt.Errorf("failed to insert blog image: %w", err)
		}
	}

	return nil
}

func (s *Store) GetImagesByBlogID(blogID int64) ([]model.BlogImage, error) {
	rows, err := s.db.Query(
		`SELECT id, blog_id, image_url, created_at 
         FROM blog_images 
         WHERE blog_id = ?`,
		blogID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	defer rows.Close()

	var images []model.BlogImage
	for rows.Next() {
		var img model.BlogImage
		if err := rows.Scan(&img.ID, &img.BlogID, &img.ImageURL, &img.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		images = append(images, img)
	}
	return images, nil
}
