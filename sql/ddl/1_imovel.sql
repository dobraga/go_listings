CREATE TABLE imovel (
    url varchar NOT NULL,
    neighborhood varchar NOT NULL,
    location_id varchar NOT NULL,
    state varchar NOT NULL,
    city varchar NOT NULL,
    "zone" varchar NOT NULL,
    business_type varchar NOT NULL,
    listing_type varchar NOT NULL,
    raw json NOT NULL,
    title varchar NULL,
    usable_area float8 NULL,
    floors int4 NULL,
    type_unit varchar NULL,
    bedrooms int4 NULL,
    bathrooms int4 NULL,
    suites int4 NULL,
    parking_spaces int4 NULL,
    amenities _varchar NULL,
    address_lat float8 NULL,
    address_lon float8 NULL,
    price float8 NULL,
    condo_fee float8 NULL,
    total_fee float8 NULL,
    linha varchar NULL,
    estacao varchar NULL,
    distance float8 NULL,
    created_date timestamp NULL,
    updated_date timestamp NULL,
    total_fee_predict float8 NULL,
    images _varchar NULL,
    address varchar NULL,
    lat_metro float8 NULL,
    lon_metro float8 NULL,
    created timestamp NULL DEFAULT now(),
    updated timestamp NULL DEFAULT now(),
    CONSTRAINT imovel_pkey PRIMARY KEY (url)
);