CREATE DATABASE airline_booking;
CREATE DATABASE kv_store;

CREATE USER 'airline_booking_user'@'%' IDENTIFIED BY 'pass1';
CREATE USER 'kv_store_user'@'%' IDENTIFIED BY 'pass2';

GRANT ALL PRIVILEGES ON airline_booking.* TO 'airline_booking_user'@'%';
GRANT ALL PRIVILEGES ON kv_store.* TO 'kv_store_user'@'%';