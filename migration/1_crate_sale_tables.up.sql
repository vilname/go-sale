create table sales
(
    id                 serial not null,
    serial_number      varchar(255),
    cup_volume         integer,
    price              real,
    discount_id        integer,
    promo_code_id      integer,
    machine_product_id integer,
    date_sale          timestamp(6),
    created            timestamp(6) with time zone,
    updated            timestamp(6) with time zone,

    primary key (id)
);

create table write_offs
(
    id          serial not null,
    cell_number integer,
    volume      integer,
    product_id  integer,
    sale_id     integer,
    unit        varchar(255)
        constraint write_offs_unit_check
            check ((unit)::text = ANY
        ((ARRAY [
        'ML':: character varying,
        'G':: character varying,
        'MG':: character varying,
        'KCAL':: character varying,
        'FLOZ':: character varying,
        'KG':: character varying,
        'MCG':: character varying,
        'OZ':: character varying,
        'KJ':: character varying
        ])::text[])
),
    created          timestamp(6) with time zone,
    updated          timestamp(6) with time zone,

    primary key (id)
);

create table payments
(
    id      serial not null,
    price   real,
    sale_id integer,
    method  varchar(255)
        constraint payments_method_check
            check ((method)::text = ANY
        ((ARRAY [
        'CASH':: character varying,
        'CARD':: character varying,
        'QR_CODE':: character varying,
        'RF_ID':: character varying
        ])::text[])
),
    created timestamp(6) with time zone,
    updated timestamp(6) with time zone,

    primary key (id)
);



alter table if exists write_offs
    add constraint FK_write_offs_sales foreign key (sale_id) references sales;

alter table if exists payments
    add constraint FK_payments_sales foreign key (sale_id) references sales;
