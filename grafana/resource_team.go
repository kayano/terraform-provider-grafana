package grafana

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gapi "github.com/nytm/go-grafana-api"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: createTeam,
		Read:   readTeam,
		Delete: deleteTeam,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func createTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	name := d.Get("name").(string)
	email := d.Get("email").(string)

	err := client.AddTeam(name, email)
	if err != nil {
		return err
	}

	result, err := client.SearchTeam(name)
	if err != nil {
		return err
	}

	var team *gapi.Team
	for _, t := range result.Teams {
		if t.Name == name {
			team = t
			break
		}
	}

	d.SetId(strconv.FormatInt(team.Id, 10))

	return readTeam(d, meta)
}

func readTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	team, err := client.Team(id)
	if err != nil {
		if err.Error() == "404 Not Found" {
			log.Printf("[WARN] removing team %d from state because it no longer exists in grafana", id)
			d.SetId("")
			return nil
		}

		return err
	}
	d.SetId(strconv.FormatInt(team.Id, 10))
	d.Set("name", team.Name)
	d.Set("email", team.Email)

	return nil
}

func deleteTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	return client.DeleteTeam(id)
}
