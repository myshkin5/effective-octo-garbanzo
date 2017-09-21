create table garbanzo (
  id serial primary key,
  api_uuid uuid not null unique,
  first_name varchar(40) not null,
  last_name varchar(40) not null
);
