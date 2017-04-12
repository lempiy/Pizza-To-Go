--user
CREATE TABLE person (
    id serial primary key,
    username varchar(100),
    password varchar(1000),
    email varchar(100),
    created_date timestamp
);

INSERT INTO person VALUES(1, 'lempiy', '1fc854110e5532480000542834f453de31936c2f', 'lempiy@mail.ru', '2017-03-19 15:30:59');

--category
CREATE TABLE pizza_category (
    id serial primary key,
    name varchar(1000) not null,
    description text,
    is_default integer
);


--pizza
CREATE TABLE pizza (
    id serial primary key,
    name varchar(1000) not null,
    user_id integer references person(id),
    category_id integer references pizza_category(id),
    size integer,
    deleted integer DEFAULT 0,
    accepted integer DEFAULT 0,
    price real,
    description text,
    img_url varchar(1000),
    created_date timestamp,
    updated_date timestamp
);

INSERT INTO pizza_category VALUES(1, 'Normal', 'Description is missing.', 1);
INSERT INTO pizza_category VALUES(2, 'Vegetarian', 'Description is missing.', 0);
INSERT INTO pizza_category VALUES(3, 'Cool', 'Description is missing.', 0);
INSERT INTO pizza_category VALUES(4, 'Hot', 'Description is missing.', 0);
INSERT INTO pizza_category VALUES(5, 'Meat', 'Description is missing.', 0);
INSERT INTO pizza_category VALUES(6, 'Wide', 'Description is missing.', 0);

--ingredient
CREATE TABLE ingredient (
  id serial primary key,
  name varchar(100),
  description text,
  image_url varchar(1000),
  price real,
  created_date timestamp,
  user_id integer references person(id)
);

INSERT INTO ingredient VALUES(1, 'pineapple', 'pineapple', 'assets/images/ananas.png', '2', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(2, 'eggplant', 'eggplant', 'assets/images/baklazhan.png', '3', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(3, 'bacon', 'bacon', 'assets/images/becone.png', '40', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(4, 'onion', 'onion', 'assets/images/cebulya.png', '1', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(5, 'mushrooms', 'mushrooms', 'assets/images/grib.png', '5', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(6, 'corn', 'corn', 'assets/images/kukurudza.png', '10', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(7, 'oleaceae', 'oleaceae', 'assets/images/maslina.png', '6', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(8, 'carrot', 'carrot', 'assets/images/morkva.png', '3', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(9, 'cucumber', 'cucumber', 'assets/images/ogirok.png', '2', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(10, 'pepper', 'pepper', 'assets/images/perec.png', '40', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(11, 'tomato', 'tomato', 'assets/images/pomidor.png', '1', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(12, 'meat-roll', 'meat-roll', 'assets/images/rulet.png', '12', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(13, 'cheese', 'cheese', 'assets/images/syr.png', '2', '2017-03-19 15:30:59', 1);
INSERT INTO ingredient VALUES(14, 'omelet', 'omelet', 'assets/images/yayco.png', '2', '2017-03-19 15:30:59', 1);

--used_ingredient
CREATE TABLE used_ingredient (
  id serial primary key,
  ingredient_id integer references ingredient(id),
  pizza_id integer references pizza(id)
);
