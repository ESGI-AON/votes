# API de vote

## Start the project
## Installation

- Clone project into your `$GOPATH`

```
git clone https://github.com/ESGI-AON/votes.git
```

### Install  dependencies

```
make install
```

### UP

```
make up
```
### RUN
```
make run
```
### DOWN

```make down```

### HOW TO TEST

You can login in POST /login as an admin with :
```json
{
  "email": "admin@admin.com",
  "pass": "toto"
}
```
