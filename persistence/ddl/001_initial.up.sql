create table octo (
  id   serial      primary key,
  name varchar(40) unique
);

create table garbanzo_type (
  id   smallint    primary key,
  name varchar(20) unique
);

insert into garbanzo_type (id, name) values
  (1001, 'DESI'),
  (1002, 'KABULI');

create table garbanzo (
  id               serial   primary key,
  api_uuid         uuid     not null unique,
  garbanzo_type_id smallint not null references garbanzo_type(id),
  octo_id          integer  not null references octo(id),
  diameter_mm      float    not null
);
