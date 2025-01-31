create table promo_codes
(
    id            serial       not null,
    code          varchar(255) not null unique,
    qty           integer,
    used          integer default 0,
    is_active     boolean      default true,
    is_selected   boolean      default false,
    discount_type varchar(255) default 'FIXED'
        constraint promo_codes_discount_type_check
            check ((discount_type)::text = ANY
        ((ARRAY [
        'PERCENT':: character varying,
        'FIXED':: character varying,
        'FREE':: character varying
        ])::text[])
) ,
    discount_amount integer      not null,
    period_from     timestamp(6),
    period_to       timestamp(6),
    organization_id integer      not null,
    created_id      varchar(255),
    updated_id      varchar(255),
    group_id        integer,
    description     text,
    created         timestamp(6) with time zone,
    updated         timestamp(6) with time zone,

    primary key (id)
);

create table schedules
(
    id            serial       not null,
    name          varchar(255) not null,
    period_from   varchar(255),
    period_to     varchar(255),
    weekday       varchar(255),
    promo_code_id integer      not null,
    created       timestamp(6) with time zone,
    updated       timestamp(6) with time zone,

    primary key (id)
);

alter table if exists schedules
    add constraint FK_schedules_promo_codes foreign key (promo_code_id) references promo_codes;

create table products
(
    id                 serial  not null,
    cell_purpose_id    integer,
    sport_pit_id       varchar(255),
    brand_id           varchar(255),
    ingredient_line_id varchar(255),
    promo_code_id      integer not null,
    created            timestamp(6) with time zone,
    updated            timestamp(6) with time zone,

    primary key (id)
);

alter table if exists products
    add constraint FK_products_promo_codes foreign key (promo_code_id) references promo_codes;

create table groups
(
    id              serial       not null,
    name            varchar(255) not null,
    organization_id integer      not null,
    created         timestamp(6) with time zone,
    updated         timestamp(6) with time zone,

    primary key (id)
);

alter table if exists promo_codes
    add constraint FK_promo_codes_groups foreign key (group_id) references groups;
