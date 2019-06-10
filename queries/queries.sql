--name: create-post
INSERT INTO posts(title, content) VALUES(%s, %s);

--name: update-post
UPDATE posts SET title=%s, content=%s WHERE id=%s;

--name: delete-post
DELETE FROM posts WHERE id=%s;

--name: select-post
SELECT id AS ID, content AS Content, title AS Title, created_at  FROM posts WHERE id=%s;

--name: select-all
SELECT id AS ID, content AS Content, title AS Title, created_at FROM posts;