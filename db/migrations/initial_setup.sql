-- migrate:up
create table role_definitions
(
    id        varchar(128) primary key,
    priority  integer,
    transient boolean,
    color     varchar(6)
);

create table role_permissions
(
    role_id    varchar(128),
    permission varchar(256),
    primary key (role_id, permission)
);

create table user_roles
(
    user_account_id uuid,
    role_id         varchar(128),
    primary key (user_account_id, role_id)
);

create table user_permissions
(
    user_account_id uuid,
    permission      varchar(128),
    primary key (user_account_id, permission)
);

-- migrate:down
-- drop table role_definitions;
-- drop table role_permissions;
