package repository

import (
	"api/src/model"
	"database/sql"
)

// Publications represents a connection to the publications table.
type Publications struct {
	db *sql.DB
}

// NewPublicationsRepository generates a connection to the publications table.
// returns a connection to the publications table.
func NewPublicationsRepository(db *sql.DB) *Publications {
	return &Publications{db}
}

// CreatePublication create a publication.
// returns the ID of the created publication or an error.
func (repository Publications) CreatePublication(publication model.Publication) (uint64, error) {

	statement, error := repository.db.Prepare(
		"insert into publications(title, content, author_id) values(?, ?, ?)",
	)
	if error != nil {
		return 0, error
	}

	defer statement.Close()

	result, error := statement.Exec(publication.Title, publication.Content, publication.AuthorId)
	if error != nil {
		return 0, error
	}

	lastInsertId, error := result.LastInsertId()
	if error != nil {
		return 0, error
	}

	return uint64(lastInsertId), nil
}

// ListMyAndFollowPublications lists the publications of the logged in user and the publications they follow.
// returns a list containing the posts of the logged in user and the posts of whom he follows or an error.
func (repository Publications) ListMyAndFollowPublications(userId uint64) ([]model.Publication, error) {

	results, error := repository.db.Query(
		`
			select 
				distinct p.*, u.nick 
			from 
				publications p
				join users u on u.id = p.author_id
				join followers f on f.user_id = p.author_id
			where 
				u.id = ? or f.follower_id = ?
			order by p.id desc
		`,
		userId,
		userId,
	)
	if error != nil {
		return nil, error
	}

	defer results.Close()

	var publications []model.Publication

	for results.Next() {
		var publication model.Publication

		if error = results.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); error != nil {
			return nil, error
		}
		publications = append(publications, publication)
	}

	return publications, nil
}

// FindPublicationById search for a publication by the given ID.
// returns the publication or an error if unable to create the publication.
func (repository Publications) FindPublicationById(publicationId uint64) (model.Publication, error) {

	results, error := repository.db.Query(
		`
			select 
				p.*, u.nick 
			from 
				publications p
				join users u on u.id = p.author_id
			where p.id = ?
		`,
		publicationId,
	)
	if error != nil {
		return model.Publication{}, error
	}

	defer results.Close()

	var publication model.Publication

	if results.Next() {
		if error = results.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); error != nil {
			return model.Publication{}, error
		}
	}

	return publication, nil
}

// UpdatePublicationById updates a publication by the given ID.
// returns an error if unable to update the publication.
func (repository Publications) UpdatePublicationById(publicationId uint64, publication model.Publication) error {

	statement, error := repository.db.Prepare("update publications set title = ?, content = ? where id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(publication.Title, publication.Content, publicationId); error != nil {
		return error
	}

	return nil
}

// DeletePublicationById deletes a publication by the given ID.
// returns an error if unable to delete the post.
func (repository Publications) DeletePublicationById(publicationId uint64) error {

	statement, error := repository.db.Prepare("delete from publications where id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(publicationId); error != nil {
		return error
	}

	return nil
}

// ListPublicationsByUserId lists a user's publications.
// returns a list of publications or an error.
func (repository Publications) ListPublicationsByUserId(userId uint64) ([]model.Publication, error) {

	results, error := repository.db.Query(
		`
			select 
				distinct p.*, u.nick 
			from 
				publications p
				join users u on u.id = p.author_id
			where 
				p.author_id = ?
			order by p.id desc
		`,
		userId,
	)
	if error != nil {
		return nil, error
	}

	defer results.Close()

	var publications []model.Publication

	for results.Next() {
		var publication model.Publication

		if error = results.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); error != nil {
			return nil, error
		}
		publications = append(publications, publication)
	}

	return publications, nil
}

// LikePublicationById likes a publication by the given ID
// returns an error if unable to like the publication.
func (repository Publications) LikePublicationById(publicationId uint64) error {

	statement, error := repository.db.Prepare("update publications set likes = likes + 1 where id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(publicationId); error != nil {
		return error
	}

	return nil
}

// UnlikePublicationById unlike a publication by the given ID
// returns an error if unable to unlike the publication.
func (repository Publications) UnLikePublicationById(publicationId uint64) error {

	statement, error := repository.db.Prepare(
		`
			update 
				publications 
			set 
				likes = case when likes > 0 then likes - 1 else 0 end
			where id = ?
		`,
	)
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(publicationId); error != nil {
		return error
	}

	return nil
}
