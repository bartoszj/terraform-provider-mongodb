package mongodb

import (
	"context"
	"github.com/bartoszj/terraform-provider-mongodb/mongodb/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBUserCreate,
		Read:   resourceMongoDBUserRead,
		Update: resourceMongoDBUserUpdate,
		Delete: resourceMongoDBUserDelete,

		Schema: map[string]*schema.Schema{
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  false,
				Sensitive: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeSet,
				Required: false,
				ForceNew: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"database": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config).client
	database := d.Get("database").(string)
	db := client.Database(database)

	// Create user
	var createUserResponse *types.Response
	createUserRequest := createUserRequestFromResourceData(d)
	if err := db.RunCommand(context.Background(), createUserRequest).Decode(&createUserResponse); err != nil {
		return err
	}

	// Read data
	var usersInfoResponse *types.UsersInfoResponse
	userInfoRequest := userInfoRequestFromResourceData(d)
	if err := db.RunCommand(context.Background(), userInfoRequest).Decode(&usersInfoResponse); err != nil {
		return err
	}

	d.SetId(*usersInfoResponse.UserInfos[0].Id)

	return resourceMongoDBUserRead(d, meta)
}

func resourceMongoDBUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config).client
	database := d.Get("database").(string)
	db := client.Database(database)

	// Read data
	var usersInfoResponse *types.UsersInfoResponse
	userInfoRequest := userInfoRequestFromResourceData(d)
	if err := db.RunCommand(context.Background(), userInfoRequest).Decode(&usersInfoResponse); err != nil {
		return err
	}

	if len(usersInfoResponse.UserInfos) == 0 {
		d.SetId("")
	} else {
		user := usersInfoResponse.UserInfos[0]
		d.Set("database", *user.Database)
		d.Set("username", *user.User)
	}

	return nil
}

func resourceMongoDBUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMongoDBUserRead(d, meta)
}

func resourceMongoDBUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config).client
	database := d.Get("database").(string)
	db := client.Database(database)

	// Drop user
	var response *types.Response
	dropUserRequest := dropUserRequestFromResourceData(d)
	if err := db.RunCommand(context.Background(), dropUserRequest).Decode(&response); err != nil {
		return err
	}

	return nil
}

func createUserRequestFromResourceData(d *schema.ResourceData) *types.CreateUserRequest {
	c := &types.CreateUserRequest{
		User:     d.Get("username").(string),
		Password: d.Get("password").(string),
		Roles:    getMongoDBUserRoles(d.Get("role").(*schema.Set), d.Get("database").(string)),
	}

	return c
}

func userInfoRequestFromResourceData(d *schema.ResourceData) *types.UserInfoRequest {
	u := &types.UserInfoRequest{
		User: d.Get("username").(string),
	}
	return u
}

func dropUserRequestFromResourceData(d *schema.ResourceData) *types.DropUserRequest {
	dr := &types.DropUserRequest{
		User: d.Get("username").(string),
	}
	return dr
}

func getMongoDBUserRoles(roles *schema.Set, defaultDatabase string) []*types.Role {
	r := make([]*types.Role, roles.Len())

	for i, role := range roles.List() {
		rm := role.(map[string]interface{})
		roleName := rm["name"].(string)
		database := rm["database"].(string)
		if len(database) == 0 {
			database = defaultDatabase
		}
		b := &types.Role{
			Role:     roleName,
			Database: database,
		}
		r[i] = b
	}
	return r
}
