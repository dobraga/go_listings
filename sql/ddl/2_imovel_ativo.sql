CREATE TABLE imovel_ativo (
    url varchar NOT NULL,
    location_id varchar NOT NULL,
    business_type varchar NOT NULL,
    listing_type varchar NOT NULL,
    updated_date timestamp NULL,
    CONSTRAINT imovel_ativo_pkey PRIMARY KEY (url)
);

ALTER TABLE
    imovel_ativo
ADD
    CONSTRAINT imovel_ativo_url_fkey FOREIGN KEY (url) REFERENCES imovel(url);
