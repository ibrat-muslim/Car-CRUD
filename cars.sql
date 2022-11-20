create table cars (
	id serial primary key,
	name varchar(255) not null,
	price numeric(18, 2) check(price > 1000) not null,
	color varchar(255) not null,
	year smallint check(year >= 0) not null,
	image_url varchar(255) not null,
	created_at timestamp(0) with time zone not null default current_timestamp
);

create table car_images (
	id serial primary key,
	car_id integer not null references cars(id) on delete restrict,
	image_url varchar(255) not null,
	sequence_number integer not null
);

create index cars_name_idx on cars(name);

create index cars_price_idx on cars(price);

create index cars_color_idx on cars(color);

create index cars_year_idx on cars(year);
