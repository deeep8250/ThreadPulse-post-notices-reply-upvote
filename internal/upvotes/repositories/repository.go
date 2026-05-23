package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/threadpulse/internal/config"
)

type UpvotesRepository struct {
	db *sqlx.DB
}

func NewUpvotesRepository() *UpvotesRepository {
	return &UpvotesRepository{
		db: config.PostgisDB,
	}
}

func (r *UpvotesRepository) CreateUpvote(postID, userID int) error {
	query := `insert into upvotes(post_id,user_id) values($1,$2)`
	_, err := r.db.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UpvotesRepository) GetUpvotes(postId int) (int, error) {
	var UpvotesCount int
	query := `select count(*) from upvotes where post_id=$1`
	err := r.db.Get(&UpvotesCount, query, postId)
	if err != nil {
		return 0, err
	}
	return UpvotesCount, nil
}

func (r *UpvotesRepository) CheckUpvote(postID, userID int) (bool, error) {
	var checkPostsExists int
	query1 := `select count(*) from posts where id=$1`
	err := r.db.Get(&checkPostsExists, query1, postID)
	if err != nil {
		return false, err
	} else if checkPostsExists == 0 {
		return false, errors.New("post id invalid")
	}

	var count int
	query := `SELECT COUNT(*) FROM upvotes WHERE post_id=$1 AND user_id=$2`
	err = r.db.Get(&count, query, postID, userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
