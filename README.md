# template-go

> Template for scaffolding your next Go http server project

## Project structure

```
.
├── cmd
│   └── server                // index entry point
├── data
│   ├── migrations            // database migration
│   └── seed                  // initial data seed
├── mocks                     // generated mocks for testing
└── pkg
    ├── config                // app configs, singleton get via config.GetConfig()
    ├── consts                // app constant
    ├── entities              // app bushiness logic
    ├── handlers              // app handles
    │   └── v1                // versioning v1
    ├── middleware            // route middleware: auth
    ├── model                 // api modal, shared across the app
    ├── monitoring            // logging package
    ├── repo
    │   └── user
    │       └── testdata      // contains seed data for repo testing
    ├── routes                // all rest api routes
    │   └── v1.go             // versioning v1
    └── util                  // contains shared utilities
        └── testutil          // helper for db testing
```

## Usage

### Available commands

**Dev**

```
# only need run once
make setup

# start local server
make dev


# build binary
make build


# run all unit test
make test

```

**Migration**

```bash
# create a new migration
make migration-new name=example-alter-table

# apply new migrations
make migrate-up

# rollback a migration version
make migrate-down

```

### Monitoring

- Get logger instance

```go
import "github.com/dwarvesf/go-template/pkg/monitoring"

m := monitoring.FromContext(ctx)
```

- Do logging in handlers and entities package

- Log string format `[package.Function] invokingMethod(param1, param2=%v)`

```go
// example
m.Errorf(err, "[entity.LoginUser] GetUserByEmail(ctx, email=%v)", email)

m.Infof("[entity.LoginUser] GetUserByEmail(ctx, email=%v)", email)

```

### Testing

- Use mockery to generate mock for testing
- Command: `make mock`
- Repository testing (DB testing) required an docker DB instance running
- check example in `repo/user` for db integration testing
  - sql seed file in testdata
  - the seed file loaded by testutils

## License

MIT &copy; [dwarvesf](github.com/dwarvesf)
