CREATE TABLE StoryBooks (
   `id` VARCHAR(64) NOT NULL,
   `title` TEXT NOT NULL,
   `character_name` TEXT NOT NULL,
    `language` TEXT NOT NULL,
    `theme` TEXT NOT NULL ,
    `title_image_path` TEXT,
    `story` JSON,
   PRIMARY KEY (id)
);