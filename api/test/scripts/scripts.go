package scripts

import (
	"database/sql"
	"log"
)

var (
	dropPublicationsTable = "drop table if exists publications"
	dropFollowersTable    = "drop table if exists followers"
	dropUsersTable        = "drop table if exists users"

	createUserTable = `
		create table users (
			id int auto_increment primary key,
			name varchar(50) not null,
			nick varchar(50) not null unique,
			email varchar(50) not null unique,
			password varchar(200) not null,
			createdAt timestamp default current_timestamp() 
		) engine=innodb;
	`
	createFollowersTable = `
		create table followers (
			user_id int not null,
			follower_id int not null,
		
			foreign key (user_id) references users(id) on delete cascade,
			foreign key (follower_id) references users(id) on delete cascade,
			
			primary key (user_id, follower_id)
		) engine=innodb;
	`
	createPublicationsTable = `
		create table publications (
			id int auto_increment primary key,
			title varchar(50) not null,
			content varchar(50) not null,
			author_id int not null,
			likes int not null default 0,
			createdAt timestamp default current_timestamp() ,
		
			foreign key (author_id) references users(id) on delete cascade
		) engine=innodb;
	`
	// Password: alb1234
	insertUsers = `
		insert into users(name, nick, email, password) values
		("User 1", "user@1", "user1@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy"),
		("User 2", "user@2", "user2@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy"),
		("User 3", "user@3", "user3@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy")
	`
	insertFollowers = `
		insert into followers(user_id, follower_id) values
		(2, 1),
		(2, 3),
		(1, 3)
	`
	insertPublications = `
		insert into publications(title, content, author_id) values
		("Pub 1 - 1", "Pub of user 1", 1),
		("Pub 1 - 2", "Pub of user 2", 2),
		("Pub 1 - 3", "Pub of user 3", 3)
	`
)

func dropTables(db *sql.DB) {
	_, error := db.Exec(dropPublicationsTable)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(dropFollowersTable)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(dropUsersTable)
	if error != nil {
		log.Fatal(error)
	}
}

func createTables(db *sql.DB) {
	_, error := db.Exec(createUserTable)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(createFollowersTable)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(createPublicationsTable)
	if error != nil {
		log.Fatal(error)
	}
}

func insertTables(db *sql.DB) {
	_, error := db.Exec(insertUsers)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(insertFollowers)
	if error != nil {
		log.Fatal(error)
	}

	_, error = db.Exec(insertPublications)
	if error != nil {
		log.Fatal(error)
	}
}

func Run(db *sql.DB) {
	dropTables(db)
	createTables(db)
	insertTables(db)
}
