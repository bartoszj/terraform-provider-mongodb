[![Build Status](https://travis-ci.org/bartoszj/terraform-provider-mongodb.svg?branch=master)](https://travis-ci.org/bartoszj/terraform-provider-mongodb)

Terraform Provider MongoDB
==========================
The MongoDB Terraform Provider


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14

Using the provider
------------------

First, install the desired [plugin release](https://github.com/bartoszj/terraform-provider-mongodb/releases) following Terraform's [Third-party plugin docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

To configure the provider:
```hcl
provider "mongodb" {
    uri = "mongodb://localhost:27017"
}
```

Optionally you can configure Username and Password:

```hcl
provider "mongodb" {
    uri      = "mongodb://localhost:27017"
    username = "admin"
    password  = "admin"
}
```

You can use environment variables to set the provider properties instead. The following table shows all the config options, the corresponding environment variables, and their property defaults if you do not set them.

| Provider property   | Env variable       | Default if not set    | Use                         |
|:--------------------|:-------------------|:----------------------|:----------------------------|
| uri                 | MONGODB_URL        | not set               | The url of the MongoDB      |
| username            | MONGODB_USERNAME   | not set               | Username for the admin user |
| password            | MONGODB_PASSWORD   | not set               | Password for the admin user |

# Resources

## User

```hcl
resource "mongodb_user" "user" {
    database = "test-db"
    username = "user"
    password = "pass"

    role {
        name = "readWrite"
    }
    role {
        name     = "dbAdmin"
        database = "test-db"
    }
    role {
        name     = "read"
        database = "admin"
    }
}
```

To import a user:
```
terraform import mongodb_user.<user_identifier> <database>.<user>
```

| User property       | Description
|:--------------------|:---------------------------------------|
| database            | Database where user will be created    |
| username            | MongoDB username                       |
| password            | MongoDB user password                  |
| role.name           | Role name                              |
| role.database       | Database where role is stored          |

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
