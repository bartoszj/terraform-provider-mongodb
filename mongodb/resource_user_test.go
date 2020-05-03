package mongodb

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccMongoDBUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		//CheckDestroy: testAccMongoDBUserDestroy,
		Steps: []resource.TestStep{
			//resource.TestStep{
			//	Config: testAccMongoDBUserConfig,
			//	Check: resource.ComposeTestCheckFunc(
			//		//testCheckMongoDBUserExists("mongodb_user.user", t),
			//	),
			//},
			resource.TestStep{
				Config: testAccMongoDBUserRolesConfig,
				Check: resource.ComposeTestCheckFunc(
					//testCheckMongoDBUserExists("mongodb_user.user", t),
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
		name = "dbAdmin"
	}
}
`)
