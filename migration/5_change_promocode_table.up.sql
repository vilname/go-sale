alter table promo_codes
    add machine_ids text;

alter table public.promo_codes
    drop constraint promo_codes_code_key;

create unique index promo_codes_code_organization_id_key
    on promo_codes (code, organization_id);
