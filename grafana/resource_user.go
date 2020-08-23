package grafana

import (
	"crypto/rand"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gapi "github.com/nytm/go-grafana-api"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Read:   ReadUser,
		Update: ReadUser,
		Delete: ReadUser,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	n := 64
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return err
	}
	pass := string(bytes[:n])
	u := gapi.User{
		Login:    d.Get("login").(string),
		Email:    d.Get("email").(string),
		Password: pass,
	}
	id, err := client.CreateUser(u)

	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(id, 10))

	return ReadUser(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	email := d.Get("email").(string)
	user, err := client.UserByEmail(email)
	if err != nil {
		if err.Error() == "404 Not Found" {
			log.Printf("[WARN] removing user %s from state because it no longer exists in grafana", email)
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(strconv.FormatInt(user.Id, 10))

	return nil
}
