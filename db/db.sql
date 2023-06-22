CREATE EXTENSION IF NOT EXISTS CITEXT; -- eliminate calls to lower

CREATE UNLOGGED TABLE users
(
    Nickname   CITEXT PRIMARY KEY,
    FullName   TEXT NOT NULL,
    About      TEXT NOT NULL DEFAULT '',
    Email      CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Title    TEXT   NOT NULL,
    "user"   CITEXT,
    Slug     CITEXT PRIMARY KEY,
    Posts    INT    DEFAULT 0,
    Threads  INT    DEFAULT 0
);

CREATE UNLOGGED TABLE thread
(
    Id      SERIAL    PRIMARY KEY,
    Title   TEXT      NOT NULL,
    Author  CITEXT    REFERENCES "users"(Nickname),
    Forum   CITEXT    REFERENCES "forum"(Slug),
    Message TEXT      NOT NULL,
    Votes   INT       DEFAULT 0,
    Slug    CITEXT,
    Created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE post
(
    Id        SERIAL      PRIMARY KEY,
    Author    CITEXT,
    Created   TIMESTAMP   WITH TIME ZONE DEFAULT now(),
    Forum     CITEXT,
    IsEdited  BOOLEAN     DEFAULT FALSE,
    Message   CITEXT      NOT NULL,
    Parent    INT         DEFAULT 0,
    Thread    INT,
    Path      INTEGER[],
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    FOREIGN KEY (author) REFERENCES "users"  (nickname)
);

CREATE UNLOGGED TABLE vote
(
    ID       SERIAL PRIMARY KEY,
    Author   CITEXT    REFERENCES "users" (Nickname),
    Voice    INT       NOT NULL,
    Thread   INT,
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    UNIQUE (Author, Thread)
);


CREATE UNLOGGED TABLE users_forum
(
    Nickname  CITEXT  NOT NULL,
    FullName  TEXT    NOT NULL,
    About     TEXT,
    Email     CITEXT,
    Slug      CITEXT  NOT NULL,
    FOREIGN KEY (Nickname) REFERENCES "users" (Nickname),
    FOREIGN KEY (Slug) REFERENCES "forum" (Slug),
    UNIQUE (Nickname, Slug)
);

--     Update vote in thread

CREATE OR REPLACE FUNCTION addUserFirstVote() RETURNS TRIGGER AS
$$
BEGIN
UPDATE thread SET Votes=(Votes+New.Voice) WHERE Id = NEW.Thread;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_vote
    AFTER INSERT ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE addUserFirstVote();


CREATE OR REPLACE FUNCTION changeVoteOnThread() RETURNS TRIGGER AS
$$
BEGIN
UPDATE thread SET Votes=(Votes+2*New.Voice) WHERE Id = NEW.Thread;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_update_vote
    AFTER UPDATE ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE changeVoteOnThread();

--     Update users_forum

CREATE OR REPLACE FUNCTION PostUpdateUserForum() RETURNS TRIGGER AS
$$
DECLARE
authorFullName TEXT;
   authorAbout    TEXT;
   authorEmail    CITEXT;
BEGIN
SELECT FullName, About, Email FROM users WHERE Nickname = NEW.Author INTO authorFullName, authorAbout, authorEmail;
INSERT INTO users_forum (Nickname, FullName, About, Email, Slug)
VALUES (NEW.Author, authorFullName, authorAbout, authorEmail, NEW.Forum)
    ON CONFLICT DO NOTHING;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER post_update_user_forum
    AFTER INSERT ON post
    FOR EACH ROW
    EXECUTE PROCEDURE PostUpdateUserForum();


CREATE OR REPLACE FUNCTION ThreadUpdateUserForum() RETURNS TRIGGER AS
$$
DECLARE
authorFullName CITEXT;
   authorAbout    CITEXT;
   authorEmail    CITEXT;
BEGIN
SELECT FullName, About, Email FROM users WHERE Nickname = NEW.Author INTO authorFullName, authorAbout, authorEmail;
INSERT INTO users_forum (Nickname, FullName, About, Email, Slug)
VALUES (NEW.Author, authorFullName, authorAbout, authorEmail, NEW.Forum)
    ON CONFLICT DO NOTHING;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER thread_update_users_forum
    AFTER INSERT ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE ThreadUpdateUserForum();

--     Update thread in forum

CREATE OR REPLACE FUNCTION addThreadInForum() RETURNS TRIGGER AS
$$
BEGIN
UPDATE forum SET Threads=(Threads + 1) WHERE Slug = NEW.Forum;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER new_thread_in_forum
    AFTER INSERT ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE addThreadInForum();

--     Update posts in forum

CREATE OR REPLACE FUNCTION addPostInForum() RETURNS TRIGGER AS
$$
BEGIN
UPDATE forum SET Posts=(Posts + 1) WHERE Slug = NEW.Forum;
return NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER new_post_in_forum
    AFTER INSERT ON post
    FOR EACH ROW
    EXECUTE PROCEDURE addPostInForum();