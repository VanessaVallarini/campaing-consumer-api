insert into "user"(id,email,created_at,updated_at,active)
values ('c3eeb9b0-051c-4803-b0a4-f6060bcb40d9','vanessa@gmail.com',now(),now(),true);

insert into slug(id,user_id,created_at,updated_at,active,lat,long)
values ('f43e580b-ffb2-490d-aea1-b2f0435d624b','c3eeb9b0-051c-4803-b0a4-f6060bcb40d9',now(),now(),true,45.6085,-73.5493);

insert into merchant(id,user_id,slug_id,created_at,updated_at,name,active,lat,long)
values ('2ed8b772-1714-46de-98ab-c2653bb03d78','c3eeb9b0-051c-4803-b0a4-f6060bcb40d9','f43e580b-ffb2-490d-aea1-b2f0435d624b',now(),now(),'Se João',true,45.6085,-73.5493);