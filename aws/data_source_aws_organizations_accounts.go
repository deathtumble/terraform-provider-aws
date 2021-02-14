package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func dataSourceAwsOrganizationsAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsOrganzationAccountsRead,

		Schema: map[string]*schema.Schema{
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAwsOrganzationAccountsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	var outputAccounts []*organizations.Account

	conn.ListAccountsPages(&organizations.ListAccountsInput{}, func(page *organizations.ListAccountsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		outputAccounts = append(outputAccounts, page.Accounts...)

		return !lastPage
	})

	var accounts []map[string]interface{}
	for _, outputAccount := range outputAccounts {
		tagsOutput, err := keyvaluetags.OrganizationsListTags(conn, *outputAccount.Id)

		if err != nil {
			return fmt.Errorf("error listing tags for Account (%s): %s", *outputAccount.Id, err)
		}

		tags := tagsOutput.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()

		accounts = append(accounts, map[string]interface{}{
			"arn":    aws.StringValue(outputAccount.Arn),
			"email":  aws.StringValue(outputAccount.Email),
			"id":     aws.StringValue(outputAccount.Id),
			"name":   aws.StringValue(outputAccount.Name),
			"status": aws.StringValue(outputAccount.Status),
			"tags":   tags,
		})
	}

	if err := d.Set("accounts", accounts); err != nil {
		return fmt.Errorf("error setting accounts: %w", err)
	}

	return nil
}
