package grafana

import (
	"crypto/rand"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	gapi "github.com/nytm/go-grafana-api"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: createUser,
		Read:   readUser,
		Delete: deleteUser,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	password, err := genPassword()
	if err != nil {
		return err
	}

	u := gapi.User{
		Email:    d.Get("email").(string),
		Password: password,
	}
	id, err := client.CreateUser(u)

	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(id, 10))

	return readUser(d, meta)
}

func readUser(d *schema.ResourceData, meta interface{}) error {
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

func deleteUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	return client.DeleteUser(id)
}

func genPassword() (string, error) {
	n := 64
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes[:n]), nil
}
