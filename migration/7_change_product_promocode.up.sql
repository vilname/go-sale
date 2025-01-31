alter table promo_codes
    add category_ids varchar(255) default '';

alter table promo_codes
    add view_ids varchar(255) default '';

alter table promo_codes
    add brand_ids varchar(255) default '';

alter table promo_codes
    add ingredient_line_ids varchar(255) default '';

alter table promo_codes
    add ingredient_ids text default '';