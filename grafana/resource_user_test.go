package grafana

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	gapi "github.com/nytm/go-grafana-api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUser_basic(t *testing.T) {
	var user gapi.User
	resourceName := "grafana_user.test_user"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccUserCheckDestroy(&user),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccUserCheckExists(resourceName, &user),
					resource.TestMatchResourceAttr(
						resourceName, "id", regexp.MustCompile(`\d+`),
					),
				),
			},
		},
	})
}

func testAccUserCheckExists(rn string, user *gapi.User) resource.TestCheckFunc {
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
			return fmt.Errorf("got a user id of 0")
		}

		users, err := client.Users()
		if err != nil {
			return fmt.Errorf("error getting users: %v", err)
		}

		var gotUser *gapi.User
		for _, u := range users {
			if u.Id == id {
				gotUser = &u
				break
			}
		}

		*user = *gotUser

		return nil
	}
}

func testAccUserCheckDestroy(user *gapi.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*gapi.Client)
		_, err := client.UserByEmail(user.Email)
		if err == nil {
			return fmt.Errorf("user still exists")
		}
		return nil
	}
}

const testAccUserConfig_basic = `
resource "grafana_user" "test_user" {
    email = "test.user@example.com"
}
`
