package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMApiManagementService_developer(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMApiManagementService_developer, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementServiceExists("azurerm_api_management_service.test"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementService_standard(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMApiManagementService_standard, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementServiceExists("azurerm_api_management_service.test"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementService_premium(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMApiManagementService_premium, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementServiceExists("azurerm_api_management_service.test"),
				),
			},
		},
	})
}

func testCheckAzureRMApiManagementServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for API Management Service: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).apiServicesClient

		resp, err := conn.Get(resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on apiServicesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: API Management Service %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("API Management Service still exists:\n%#v", resp)
		}
	}

	return nil
}

var testAccAzureRMApiManagementService_developer = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_api_management_service" "test" {
    name = "acctestavset-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"

    sku {
      type = "Standard"
      capacity = 1
    }

    publisher {
      name  = "Terraform Acceptance Tests"
      email = "tfacctests@somedomain.com"
    }
}
`

var testAccAzureRMApiManagementService_standard = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_api_management_service" "test" {
    name = "acctestavset-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"

    sku {
      type = "Standard"
      capacity = 1
    }

    publisher {
      name  = "Terraform Acceptance Tests"
      email = "tfacctests@somedomain.com"
    }
}
`

var testAccAzureRMApiManagementService_premium = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_api_management_service" "test" {
    name = "acctestavset-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"

    sku {
      type = "Premium"
      capacity = 2
    }

    publisher {
      name  = "Terraform Acceptance Tests"
      email = "tfacctests@somedomain.com"
    }
}
`
