package mongodb

import (
	"bytes"
	"context"
	"github.com/bartoszj/terraform-provider-mongodb/mongodb/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
	"time"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBUserCreate,
		Read:   resourceMongoDBUserRead,
		Update: resourceMongoDBUserUpdate,
		Delete: resourceMongoDBUserDelete,
		Exists: resourceMongoDBUserExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

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
	username := d.Get("username").(string)
	db := client.Database(database)

	// Create user
	var createUserResponse types.Response
	createUserRequest := createUserRequestFromResourceData(d)
	ctx, _ := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	if err := db.RunCommand(ctx, createUserRequest).Decode(&createUserResponse); err != nil {
		return err
	}

	var id bytes.Buffer
	id.WriteString(database)
	id.WriteString(".")
	id.WriteString(username)

	d.SetId(id.String())

	return resourceMongoDBUserRead(d, meta)
}

func resourceMongoDBUserRead(d *schema.ResourceData, meta interface{}) error {
	var ids = strings.Split(d.Id(), ".")
	database := ids[0]
	username := ids[1]

	client := meta.(*config).client
	db := client.Database(database)

	// Read data
	var usersInfoResponse types.UsersInfoResponse
	userInfoRequest := userInfoRequestFromResourceData(username)
	ctx, _ := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	if err := db.RunCommand(ctx, userInfoRequest).Decode(&usersInfoResponse); err != nil {
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
	client := meta.(*config).client
	database := d.Get("database").(string)
	db := client.Database(database)

	// Update user
	var updateUserResponse types.Response
	updateUserRequest := updateUserRequestFromResourceData(d)
	ctx, _ := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	if err := db.RunCommand(ctx, updateUserRequest).Decode(&updateUserResponse); err != nil {
		return err
	}

	return resourceMongoDBUserRead(d, meta)
}

func resourceMongoDBUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config).client
	database := d.Get("database").(string)
	db := client.Database(database)

	// Drop user
	var response types.Response
	dropUserRequest := dropUserRequestFromResourceData(d)
	ctx, _ := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	if err := db.RunCommand(ctx, dropUserRequest).Decode(&response); err != nil {
		return err
	}

	return nil
}

func resourceMongoDBUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := resourceMongoDBUserRead(d, meta)
	return err == nil, nil
}

func createUserRequestFromResourceData(d *schema.ResourceData) *types.CreateUserRequest {
	c := &types.CreateUserRequest{
		User:     d.Get("username").(string),
		Password: d.Get("password").(string),
		Roles:    getMongoDBUserRoles(d.Get("role").(*schema.Set), d.Get("database").(string)),
	}

	return c
}

func userInfoRequestFromResourceData(username string) *types.UserInfoRequest {
	u := &types.UserInfoRequest{
		User: username,
	}
	return u
}

func updateUserRequestFromResourceData(d *schema.ResourceData) *types.UpdateUserRequest {
	u := &types.UpdateUserRequest{
		User:     d.Get("username").(string),
		Password: d.Get("password").(string),
		Roles:    getMongoDBUserRoles(d.Get("role").(*schema.Set), d.Get("database").(string)),
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
