create table todo (
  id bigserial primary key
  , description text not null
  , status text not null
);

create table action (
  id bigserial primary key
  , todo_id bigint not null
  , category text not null
  , performed_at timestamp with time zone not null
  , foreign key (todo_id) references todo (id)
);
