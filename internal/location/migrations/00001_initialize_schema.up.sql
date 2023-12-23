create table if not exists driver_location
(
    driver_id varchar(255)     not null
        constraint driver_location_pk
            primary key,
    lat       double precision not null,
    lng       double precision not null
);

alter table driver_location
    owner to location;
