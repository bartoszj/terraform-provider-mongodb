package mongodb

import (
	"context"
	"fmt"
	"github.com/bartoszj/terraform-provider-mongodb/mongodb/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

func TestAccMongoDBUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMongoDBUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoDBUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBUserExists("mongodb_user.user"),
					testAccCheckMongoDBUserRoles("mongodb_user.user", []types.Role{}),
					resource.TestCheckNoResourceAttr("mongodb_user.user", "role"),
				),
			},
		},
	})
}

func TestAccMongoDBUser_roles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMongoDBUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoDBUserRolesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBUserExists("mongodb_user.user"),
					testAccCheckMongoDBUserRoles("mongodb_user.user", testAccMongoDBUserRoles),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.0.name", "readWrite"),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.0.database", ""),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.1.name", "dbAdmin"),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.1.database", "testing"),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.2.name", "read"),
					resource.TestCheckResourceAttr("mongodb_user.user", "role.2.database", "admin"),
				),
			},
		},
	})
}

var testAccMongoDBUserConfig = fmt.Sprintf(`
resource "mongodb_user" "user" {
	database = "testing"
    username = "user"
    password = "pass"
}
`)

var testAccMongoDBUserRolesConfig = fmt.Sprintf(`
resource "mongodb_user" "user" {
	database = "testing"
    username = "user"
    password = "pass"
	role {
		name = "readWrite"
	}
	role {
		name     = "dbAdmin"
		database = "testing"
	}
	role {
		name     = "read"
		database = "admin"
	}
}
`)

var testAccMongoDBUserRoles = []types.Role{{Role: "readWrite", Database: "testing"}, {Role: "dbAdmin", Database: "testing"}, {Role: "read", Database: "admin"}}

func testAccCheckMongoDBUserExists(resourceKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).client
		var ids = strings.Split(rs.Primary.ID, ".")
		database := ids[0]
		username := ids[1]
		db := client.Database(database)

		// Read data
		var usersInfoResponse types.UsersInfoResponse
		userInfoRequest := userInfoRequestFromResourceData(username)
		if err := db.RunCommand(context.Background(), userInfoRequest).Decode(&usersInfoResponse); err != nil {
			return err
		}

		if len(usersInfoResponse.UserInfos) == 0 {
			return fmt.Errorf("user %v not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckMongoDBUserRoles(resourceKey string, roles []types.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).client
		var ids = strings.Split(rs.Primary.ID, ".")
		database := ids[0]
		username := ids[1]
		db := client.Database(database)

		// Read data
		var usersInfoResponse types.UsersInfoResponse
		userInfoRequest := userInfoRequestFromResourceData(username)
		if err := db.RunCommand(context.Background(), userInfoRequest).Decode(&usersInfoResponse); err != nil {
			return err
		}

		if len(usersInfoResponse.UserInfos) == 0 {
			return fmt.Errorf("user %v not found", rs.Primary.ID)
		}

		user := usersInfoResponse.UserInfos[0]
		if len(user.Roles) != len(roles) {
			return fmt.Errorf("user have %d roles frome %d expected", len(user.Roles), len(roles))
		}

		for i, _ := range user.Roles {
			mongoRole := user.Roles[i]
			expectedRole := roles[i]

			if mongoRole.Role != expectedRole.Role {
				return fmt.Errorf("mongodb role %s doesn't match expected %s", mongoRole.Role, expectedRole.Role)
			}

			if mongoRole.Database != expectedRole.Database {
				return fmt.Errorf("mongodb database %s doesn't match expected %s", mongoRole.Database, expectedRole.Database)
			}
 		}

		return nil
	}
}

func testAccMongoDBUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*config).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodb_user" {
			continue
		}

		username := rs.Primary.Attributes["username"]
		database := rs.Primary.Attributes["database"]

		db := client.Database(database)
		// Read data
		var usersInfoResponse types.UsersInfoResponse
		userInfoRequest := userInfoRequestFromResourceData(username)
		if err := db.RunCommand(context.Background(), userInfoRequest).Decode(&usersInfoResponse); err != nil {
			return err
		}
	}
	return nil
}
