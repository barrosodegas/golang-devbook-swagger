insert into users(name, nick, email, password) values
("User 1", "user@1", "user1@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy"), -- alb1234
("User 2", "user@2", "user2@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy"), -- alb1234
("User 3", "user@3", "user3@gmail.com", "$2a$10$iONPxhg4osPz6GxXm5uM/.0m/tGRByVOdNMIlqZrjrvtrgjcepLsy"); -- alb1234

insert into followers(user_id, follower_id) values
(2, 1),
(2, 3),
(1, 3);

insert into publications(title, content, author_id) values
("Teste P1", "Teste P1...", 1),
("Teste P1 - 2", "Teste P1 - 2...", 2),
("Teste P1 - 3", "Teste P1 - 3...", 3);
