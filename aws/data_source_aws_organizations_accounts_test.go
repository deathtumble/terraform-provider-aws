package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSAvailabilityZones_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsOrganizationAccountsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrganizationAccountsState("data.organization_accounts.accounts"),
				),
			},
		},
	})
}

func testAccCheckOrganizationAccountsCheck(n string) resource.TestCheckFunc {
	return func(S *terraformState) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find accounts resource: %s", n)
		}
	}
}
