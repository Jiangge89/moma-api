CREATE DATABASE mydatabase;

CREATE USER 'moma'@'localhost' IDENTIFIED BY 'duftee2023';

GRANT ALL PRIVILEGES ON moma_api.* TO 'moma'@'localhost';
FLUSH PRIVILEGES;
SHOW GRANTS FOR 'moma'@'localhost';
