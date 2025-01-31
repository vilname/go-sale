alter table public.sales
    add name varchar(255);

update sales set promo_code_id = null where promo_code_id is not null;
update sales set name = 'Тестовое название' where name is null or name = '';

alter table public.sales
    alter column name set not null;

alter table public.sales
    rename column cup_volume to volume;

alter table public.sales
    alter column volume set not null;

alter table public.sales
    alter column price set not null;

alter table if exists sales
    add constraint FK_sales_promo_codes foreign key (promo_code_id) references promo_codes;

alter table sales
    add unit varchar(255)
        constraint sales_unit
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
                       ])::text[]));

update sales set unit = 'ML' where unit is null or unit = '';

alter table public.sales
    alter column unit set not null;
