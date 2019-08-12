-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- User.
create table user (
  id integer primary key autoincrement,
  email varchar(64) not null unique,
  name varchar(64) not null,
  password blob not null,
  -- ID card.
  rg varchar (64) default "",
  cpf varchar(64) default "",
  -- Mobile number.
  mobile varchar(64) default "",
  createdAt date not null,
  updatedAt date not null,
  permission integer default 0,
  saved boolean default true
);
create index idx_user_email on user(email);

-- Email confirmation.
create table email_confirmation (
  email varchar(64) primary key,
  uuid varchar(64) not null,
  -- Name must be empty when used for change email, instead to create a new user.
  name varchar(64) not null, 
  password blob not null,
  createdAt date not null
);

-- Password reset.
create table password_reset (
  user_email varchar(64) primary key,
  uuid varchar(64) not null,
  createdAt date not null,
  foreign key(user_email) references user(email)
);

-- Session uuid.
create table sessionUUID (
  uuid varchar(64) primary key,
  user_id varchar(64) not null,
  createdAt date not null,
  foreign key(user_id) references user(id)
);

-- Student.
create table student (
  id integer primary key autoincrement,
  name varchar(64) not null,
  -- mobile number
  mobile varchar(64) null,
  email varchar(64) not null unique,
  createdAt date not null
);
create index idx_student_name on student(name);

-- create table parent(a primary key, b unique, c, d, e, f);
-- create unique index i1 on parent(c, d);
-- create index i2 on parent(e);
-- create unique index i3 on parent(f collate nocase);
