package aws

import (
	"fmt"
	"time"

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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"joinedMethod": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"joinedtimestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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

	accountsInput := &organizations.ListAccountsInput{}

	var outputAccounts []*organizations.Account

	conn.ListAccountsPages(accountsInput, func(page *organizations.ListAccountsOutput, lastPage bool) bool {
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

		accounts = append(accounts, map[string]interface{}{
			"arn":             aws.StringValue(outputAccount.Arn),
			"id":              aws.StringValue(outputAccount.Id),
			"email":           aws.StringValue(outputAccount.Email),
			"joinedmethod":    aws.StringValue(outputAccount.JoinedMethod),
			"joinedtimestamp": aws.TimeValue(outputAccount.JoinedTimestamp).Format(time.RFC3339),
			"name":            aws.StringValue(outputAccount.Name),
			"status":          aws.StringValue(outputAccount.Status),
		})
		if err := d.Set("tags", tagsOutput.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
			return fmt.Errorf("error setting tags: %s", err)
		}
	}

	if err := d.Set("accounts", accounts); err != nil {
		return fmt.Errorf("error setting ids: %w", err)
	}

	return nil
}
