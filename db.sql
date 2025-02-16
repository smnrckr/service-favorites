create table
    favorite_list (
        id serial primary key,
        user_id integer not null,
        name varchar(255) not null
    );

create table
    favorite (
        id serial primary key,
        user_id integer not null,
        product_id integer not null,
        list_id integer not null constraint favorite_list_id___fk references favorite_list,
        constraint favorite_product_id_list_id_unique_key_pk unique (product_id, list_id)
    );