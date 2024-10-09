-- SHOW CREATE TABLE matches;
-- Delete a foreign key

-- ALTER TABLE matches DROP FOREIGN KEY fk_matches_owner;
-- ALTER TABLE matches DROP FOREIGN KEY fk_matches_matched_user;
-- ALTER TABLE matches DROP FOREIGN KEY fk_matches_owner_book;
-- ALTER TABLE matches DROP FOREIGN KEY fk_matches_matched_book;

-- DROP INDEX idx_matches_owner_book_id ON matches;

-- ALTER TABLE matches 
-- ADD CONSTRAINT fk_matches_owner FOREIGN KEY (owner_id) REFERENCES users(user_id);

-- ALTER TABLE matches 
-- ADD CONSTRAINT fk_matches_matched_user FOREIGN KEY (matched_user_id) REFERENCES users(user_id);

-- ALTER TABLE matches 
-- ADD CONSTRAINT fk_matches_owner_book FOREIGN KEY (owner_book_id) REFERENCES books(book_id);

-- ALTER TABLE matches 
-- ADD CONSTRAINT fk_matches_matched_book FOREIGN KEY (matched_book_id) REFERENCES books(book_id);

-- SELECT * FROM matches;
