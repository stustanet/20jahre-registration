CREATE TABLE registration (
    id SERIAL PRIMARY KEY,
    name VARCHAR (60) NOT NULL,
    attending BOOLEAN,
    adults SMALLINT NOT NULL,
    adults_veg SMALLINT NOT NULL,
    children SMALLINT NOT NULL,
    children_veg SMALLINT NOT NULL,
    comment VARCHAR (500)
);
