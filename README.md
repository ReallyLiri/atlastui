# AtlasTUI

Pronounced "atla-stewie", `atlastui` is an atlas text user interface.

## TODO

### Technicalities

- Handle non-tty terminals
- Verify terminal resize
- Verify Windows compatibility
- Handle width overflows

### Features

- Support non-json formats (hcl, mermaid)
- Support views, funcs, procedures.
- Support multiple schemas

### Run

-u "postgres://admin:1234@localhost:5432/import?sslmode=disable"
-f spec.json
