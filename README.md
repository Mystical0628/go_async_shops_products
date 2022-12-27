***

## Configuration

__Configuration variables are located in the root directory in the environment file ```.env```__

### Database
```
DB_DRIVER               # DB driver, NOT USED YET. Default - mysql
DB_HOST                 # DB host. Default - 127.0.0.1
DB_PORT                 # DB port. Default - 3306
DB_DATABASE             # DB database name. Default - go_mysql_test
DB_USERNAME             # DB username. Default - root
DB_PASSWORD             # DB password. Default - root
```

### Migrations
```
MIGRATION_DIR           # Directory where placed migrations. Default - migrations
MIGRATION_EXT           # Migration files extension. Default - sql
MIGRATION_SEQ           # Sequence mode, otherwise using time. Default - true
MIGRATION_SEQ_DIGITS    # Number of digits in the sequence. Default - 6
```




___

## Migrate

__Database migrations written in Go using library [golang-migrate/migrate](https://github.com/golang-migrate/migrate).__

### Base usage
```bash
$ migrate -help
Usage: migrate OPTIONS COMMAND [arg...]
       migrate [ -help ]
Options:
  -help            Print usage
Commands:
  create NAME     Create a set of timestamped up/down migrations titled NAME
  delete V        Delete migration version V
  up [N] [-all]   Apply all or N up migrations
  down [N] [-all] Apply all or N down migrations
  drop            Drop everything inside database
  force V         Set version V but don`t run migration (ignores dirty state)
```

### Start migrations
So to start migrations type:
```bash
$ go run migration.go up -all # Apply all migrations
```

### Errors
If you get an error *"Dirty database version V. Fix and force version."* while use ```up``` or ```down```  Try to fix an error and type
```bash
$ go run migration.go force V
```
and then try to execute the command again.
***




___

## Sow

__Database seeder written in Go using library [bxcodec/faker](https://github.com/bxcodec/faker).__

#### Attention!
Due to the fact that the products table depends on shops table (by column shop_id), at the begining you should sow shops table and only then sow products table, otherwise you get an error. Range of shop_id values automatically sets by counts of shops rows

### Base usage
```bash
$ seeder -help
Usage: seeder OPTIONS SEED [arg...]
       seeder [ -help ]
Options:
  -help            Print usage
Seeds:
  product [N]     Sow N products
  shop [N]        Sow N shops
  truncate TABLE  Truncate the table TABPLE
```

### Start seeding
So to start seeding type:
```bash
$ go run seeder.go SEED [N] # Apply all migrations
```


## The end!