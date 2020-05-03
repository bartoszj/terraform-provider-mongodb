package mongodb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_URL", ""),
				Description: "The MongoDB url",
			},
			"username": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_USERNAME", ""),
				Description: "The MongoDB username",
			},
			"password": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_PASSWORD", ""),
				Description: "The MongoDB password",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mongodb_user": resourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	uri := d.Get("uri").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	opts := options.Client()
	if len(uri) > 0 {
		opts.ApplyURI(uri)
	}
	if len(username) > 0 && len(password) > 0 {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	config := &config{
		client: client,
	}

	return config, nil
}
