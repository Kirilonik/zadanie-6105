CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Тестовое заполнение таблицы user
INSERT INTO "user" (username, first_name, last_name)
VALUES 
('user1', 'Иван', 'Иванов'),
('user2', 'Петр', 'Петров'),
('user3', 'Анна', 'Аннова');

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Тестовое заполнение таблицы employee
INSERT INTO employee (username, first_name, last_name)
VALUES 
('user1', 'Иван', 'Иванов'),
('user2', 'Петр', 'Петров'),
('user3', 'Анна', 'Аннова');

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Тестовое заполнение таблицы organization
INSERT INTO organization (name, description, type)
VALUES 
('ООО "Строительная компания"', 'Компания занимается строительством жилых и коммерческих зданий', 'LLC'),
('ИП "Мебельный Мастер"', 'Индивидуальный предприниматель занимается производством офисной мебели', 'IE'),
('АО "Транспортная Компания"', 'Компания занимается транспортировкой грузов по всей стране', 'JSC');

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (user_id) REFERENCES employee(id)
);

-- Тестовое заполнение таблицы organization_responsible
INSERT INTO organization_responsible (organization_id, user_id)
VALUES 
((SELECT id FROM organization LIMIT 1), (SELECT id FROM employee LIMIT 1)),
((SELECT id FROM organization OFFSET 1 LIMIT 1), (SELECT id FROM employee OFFSET 1 LIMIT 1)),
((SELECT id FROM organization OFFSET 2 LIMIT 1), (SELECT id FROM employee OFFSET 2 LIMIT 1));

-- Создание таблицы tender
CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    service_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    organization_id UUID NOT NULL,
    creator_username VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (creator_username) REFERENCES employee(username)
);

-- Тестовое заполнение таблицы tender
INSERT INTO tender (name, description, service_type, status, organization_id, creator_username)
VALUES 
('Тендер на строительство офиса', 'Требуется построить новый офис в центре города', 'Construction', 'Created', (SELECT id FROM organization LIMIT 1), 'user1'),
('Тендер на доставку оборудования', 'Нужно доставить оборудование из Москвы в Санкт-Петербург', 'Delivery', 'Published', (SELECT id FROM organization OFFSET 1 LIMIT 1), 'user2'),
('Тендер на производство мебели', 'Требуется изготовить партию офисной мебели', 'Manufacture', 'Closed', (SELECT id FROM organization OFFSET 2 LIMIT 1), 'user3');

-- Создание таблицы bid
CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    tender_id UUID NOT NULL,
    author_type VARCHAR(20) NOT NULL,
    author_id UUID NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tender_id) REFERENCES tender(id)
);

-- Тестовое заполнение таблицы bid
INSERT INTO bid (name, description, status, tender_id, author_type, author_id)
VALUES 
('Предложение по строительству', 'Готовы построить офис за 6 месяцев', 'Created', (SELECT id FROM tender WHERE name LIKE '%строительство%'), 'Organization', (SELECT id FROM organization LIMIT 1)),
('Быстрая доставка', 'Доставим оборудование за 3 дня', 'Published', (SELECT id FROM tender WHERE name LIKE '%доставку%'), 'User', (SELECT id FROM employee WHERE username = 'user2')),
('Эксклюзивная мебель', 'Изготовим мебель из экологичных материалов', 'Canceled', (SELECT id FROM tender WHERE name LIKE '%мебели%'), 'Organization', (SELECT id FROM organization OFFSET 2 LIMIT 1));

-- Создание таблицы bid_review
CREATE TABLE bid_review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bid_id) REFERENCES bid(id)
);

-- Тестовое заполнение таблицы bid_review
INSERT INTO bid_review (bid_id, description)
VALUES 
((SELECT id FROM bid WHERE name = 'Предложение по строительству'), 'Хорошее предложение, но слишком долгий срок'),
((SELECT id FROM bid WHERE name = 'Быстрая доставка'), 'Отличные условия, принимаем'),
((SELECT id FROM bid WHERE name = 'Эксклюзивная мебель'), 'Интересное предложение, но слишком дорого');


-- Создание таблицы bid_history
CREATE TABLE bid_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    tender_id UUID NOT NULL,
    author_type VARCHAR(20) NOT NULL,
    author_id UUID NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bid_id) REFERENCES bid(id)
);

-- Создание таблицы tender_history
CREATE TABLE tender_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    service_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    organization_id UUID NOT NULL,
    creator_username VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (tender_id) REFERENCES tender(id)
);