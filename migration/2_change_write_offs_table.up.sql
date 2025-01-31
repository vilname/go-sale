alter table write_offs drop column product_id;

alter table write_offs
    add ingredient_id integer;

update write_offs set ingredient_id = 7 where ingredient_id is null;

alter table write_offs
    alter column ingredient_id set not null;

alter table write_offs
    alter column sale_id set not null;


alter table write_offs
    add cell_type varchar(255)
        constraint write_offs_cell_type
            check ((cell_type)::text = ANY
        ((ARRAY [
        'INGREDIENT':: character varying,
        'CUP':: character varying,
        'WATER':: character varying,
        'DISPOSABLE':: character varying
        ])::text[]));

alter table sales
    drop column machine_product_id;

