-- migrate:up
create table role_definitions
(
    id        uuid primary key,
    name      varchar(128),
    type      varchar(64),
    priority  integer,
    transient boolean,
    color     varchar(6)
);

create table role_permissions
(
    role_id    uuid,
    permission varchar(256),
    primary key (role_id, permission),
    foreign key (role_id) references role_definitions (id)
);

create table user_roles
(
    user_account_id uuid,
    role_id         uuid,
    primary key (user_account_id, role_id),
    foreign key (role_id) references role_definitions (id)
);

create table user_permissions
(
    user_account_id uuid,
    permission      varchar(128),
    primary key (user_account_id, permission)
);

-- migrate:down
drop table role_definitions;
drop table role_permissions;
drop table user_roles;
drop table user_permissions;

