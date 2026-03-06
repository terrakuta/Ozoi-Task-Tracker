ALTER TABLE ozoi_tasks ADD COLUMN user_id UUID NOT NULL;

ALTER TABLE ozoi_tasks ADD CONSTRAINT fk_ozoi_user FOREIGN KEY (user_id) REFERENCES ozoi_users(id) ON DELETE CASCADE;
