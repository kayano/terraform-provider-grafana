package grafana

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	gapi "github.com/nytm/go-grafana-api"
)

func TestAccTeam_basic(t *testing.T) {
	var team gapi.Team
	resourceName := "grafana_team.test_team"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTeamCheckDestroy(&team),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccTeamCheckExists(resourceName, &team),
					resource.TestMatchResourceAttr(
						resourceName, "id", regexp.MustCompile(`\d+`),
					),
					resource.TestCheckResourceAttr(resourceName, "name", "test team"),
				),
			},
		},
	})
}

func testAccTeamCheckExists(rn string, team *gapi.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s\n %#v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*gapi.Client)
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		if id == 0 {
			return fmt.Errorf("got a team id of 0")
		}

		gotTeam, err := client.Team(id)
		if err != nil {
			return fmt.Errorf("error getting team: %v", err)
		}

		*team = *gotTeam

		return nil
	}
}

func testAccTeamCheckDestroy(team *gapi.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*gapi.Client)
		_, err := client.Team(team.Id)
		if err == nil {
			return fmt.Errorf("team still exists")
		}
		return nil
	}
}

const testAccTeamConfig_basic = `
resource "grafana_team" "test_team" {
    name = "test team"
}
`
