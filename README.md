# whim_assigment

## Description

Whim Assigment is a repository for Whim Assigment Test Assignment. Its a web-service app for currency conversion. It's using **Go** as its programming language.

## Onboarding and Development Guide

### Documentation

- Blueprint
  For API documentation, you can see it [here](https://rbpermadi.github.io/docs/whim_assignment_blueprint.html)

- Database diagram
  ```
  ----------------------------                    ----------------------------
  |        Currencies        |                    |        Conversions       |
  ----------------------------                    ----------------------------
  | id unsigned bigint (pk)  |---------|          | id unsigned bigint (pk)  |
  | name varchar(50)         |         |---------<| currency_id_from bigint  |
  | created_at datetime      |         |---------<| currency_id_to bigint    |
  | updated_at datetime      |                    | rate float               |
  ----------------------------                    | created_at datetime      |
                                                  | updated_at datetime      |
                                                  ----------------------------
  ```

### Prequisites

* [**Git**](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* [**Go (1.9.7 or later)**](https://golang.org/doc/install)
* [**Go Dep 0.5 or later**](https://golang.github.io/dep/docs/installation.html)
* [**MySQL**](https://www.mysql.com/downloads/)
* [**Docker**](https://docs.docker.com/get-docker/)

### Setup

- Please install/clone the [Prequisites](#prequisites) and make sure it works perfectly on your local machine.

- Clone this repo in your local at $GOPATH/src/github.com/rbpermadi If you have not set your GOPATH, set it using [this](https://golang.org/doc/code.html#GOPATH) guide. If you don't have directory src, github.com, or rbpermadi in your GOPATH, please make them.

    ```
    cd $GOPATH/src/github.com/rbpermadi
    git clone git@github.com:rbpermadi/whim_assigment.git
    ```

- Go to Whim Assigment directory and sync the vendor file

    ```
    cd $GOPATH/src/github.com/rbpermadi/whim_assigment
    make dep
    ```

- Copy and edit(optional) `env.sample`

    ```
    cp env.sample .env
    ```

### Running the app with docker

If you want to docker-compose to run **Whim Assigment**, you can use the command below. But you must stop mysql service on your PC since docker image is also run mysql service.

- Compile code, its required if its your first time trying to run the app. if not, you can ignore it.
```
> make compile
> make build
```

- Running the app.
```
> docker-compose up
```
To kill the server you just need to run `docker-compose down`

### Running the app without docker

- To prepare database, you can use your own mysql_client to import db/whim_development.sql.

Finally, run **Whim Assigment** in your local machines.

```
> make run
```

To kill the server you just need to hold `Ctrl + C`


### Test

- Run all tests

  ```
  make test
  ```

- Test coverage

  ```sh
  make coverage
  ```

### Contributing

1. Make new branch with descriptive name about your change(s) and checkout to that branch

   ````
   git checkout -b branch_name
   ````


2. Commit and push your change to upstream

   ````
   git commit -m "message"
   git push [remote_name] [branch_name]
   ````

3. Open pull request in `Github`

4. Ask someone to review your code.

5. If your code is approved, the pull request can be merged.

## FAQ

> Not available yet
