#  Simple bank example using GoLang

This project it's to understand the mainly tools used for language GoLang to API services, Migration and another stuffs.

##  Setup and Tools used

- GoLang 1.21
- MigrationDB
- Docker
    - Postgres Alpine 12

This list does not complete yet, will increase with the time.

##  Goals

- Understanding the tools most common used with GoLang

- Implementation using Docker to configuration in the environment

- Try to adopt the best practices for GoLang Projects

- Understand the tests with GoLang

## Using the project

- Open a terminal in the directory root this project and run:
```console
user@machine:~$ make postgres

```
- Create the database in the image of postgres
```console
user@machine:~$ make createdb

```
- Run the migrateDB to use the version of database until the moment
```console
user@machine:~$ make migratedb

```
This documentation it's not complete, I'm still working on this resource to automatization better

##  Database initial version
