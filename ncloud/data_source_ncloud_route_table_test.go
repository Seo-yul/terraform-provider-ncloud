package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudRouteTable_basic(t *testing.T) {
	name := fmt.Sprintf("test-table-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_route_table.foo"
	dataName := "data.ncloud_route_table.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTableConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "route_table_no", resourceName, "route_table_no"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "supported_subnet_type", resourceName, "supported_subnet_type"),
					resource.TestCheckResourceAttrPair(dataName, "is_default", resourceName, "is_default"),
					resource.TestCheckResourceAttrPair("data.ncloud_route_table.by_vpc_no", "name", resourceName, "name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudRouteTable_byVpcNo(t *testing.T) {
	dataName := "data.ncloud_route_table.by_vpc_no"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTableConfigVpcNo(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func testAccDataSourceNcloudRouteTableConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_route_table" "foo" {
	vpc_no                = ncloud_vpc.vpc.vpc_no
	name                  = "%[1]s"
	description           = "for test"
	supported_subnet_type = "PUBLIC"
}

data "ncloud_route_table" "by_id" {
	route_table_no        = ncloud_route_table.foo.id
}

data "ncloud_route_table" "by_vpc_no" {
	vpc_no                = ncloud_vpc.vpc.vpc_no
	name                  = "%[1]s"
}

`, name)
}

func testAccDataSourceNcloudRouteTableConfigVpcNo() string {
	return fmt.Sprintf(`
data "ncloud_route_table" "by_vpc_no" {
	vpc_no     = "815"
	is_default = false
	supported_subnet_type = "PUBLIC"
}
`)
}
