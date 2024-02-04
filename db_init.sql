--login using postgres user
--psql -U postgres postgres
--create user
create user db_user with password 'secret';

--create database for master data
create database kakeibo_user;
--connect to mf_master database
\c kakeibo_user
--grant all privileges to db_user for mf_master database using below command
GRANT ALL ON SCHEMA public to db_user;
