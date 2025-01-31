alter table promo_codes
    alter column period_from type date using period_from::date;

alter table promo_codes
    alter column period_to type date using period_to::date;