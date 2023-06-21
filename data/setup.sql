create table users(
                      id serial PRIMARY KEY,
                      uuid varchar(64) not null UNIQUE,
                      name varchar(255),
                      email varchar(255) not null UNIQUE,
                      user_name varchar(255) not null UNIQUE,
                      PASSWORD varchar(255) not null,
                      created_at timestamp not null

);

CREATE table sessions (
                          id serial PRIMARY KEY,
                          uuid varchar(64) not null UNIQUE,
                          name varchar(255),
                          email varchar(255),
                          user_name varchar(255),
                          user_id integer references users(id) ,
                          created_at timestamp not null
);

create table threads (
                         id serial primary key,
                         uuid varchar(64) not null unique,
                         topic text,
                         user_id integer references users(id),
                         created_atetimestamp not null
);

create table posts(
                      id serial primary key,
                      uuid varchar(64) not null unique,
                      body text,
                      user_id integer references users(id),
                      thread_id integer references threads(id),
                      created_at timestamp not null
)











