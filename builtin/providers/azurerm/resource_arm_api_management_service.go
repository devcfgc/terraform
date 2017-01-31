package azurerm

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"net/http"
)

func resourceArmApiManagementService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementServiceCreateUpdate,
		Read:   resourceArmApiManagementServiceRead,
		Update: resourceArmApiManagementServiceCreateUpdate,
		Delete: resourceArmApiManagementServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
						},
					},
				},
			},

			"publisher": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
						},

						"email": {
							Type:         schema.TypeString,
							Required:     true,
						},
					},
				},
			},
		},
	}
}

func resourceArmApiManagementServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiServicesClient

	log.Printf("[INFO] preparing arguments for Azure ARM API Management Service creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)

	sku := expandArmApiManagementServiceSku(d)
	properties := expandArmApiManagementServiceProperties(d)

	parameters := apimanagement.ServiceResource{
		Location: &location,
		Sku: &sku,
		ServiceProperties: &properties,
	}


	_, err := client.CreateOrUpdate(resGroup, name, parameters)
	if err != nil {
		return err
	}

	_, err = client.Get(resGroup, name)
	if err != nil {
		return err
	}


	//if read.ID == nil {
		//return fmt.Errorf("Cannot read EventHub %s (resource group %s) ID", name, resGroup)
	//}

	//d.SetId(*read.ID)
	d.SetId("carlostest")

	return resourceArmApiManagementServiceRead(d, meta)
}

func resourceArmApiManagementServiceRead(d *schema.ResourceData, meta interface{}) error {

	//flattenArmApiManagementServiceSku(d, read.)
	return nil
}

func resourceArmApiManagementServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["service"]

	resp, err := client.Delete(resGroup, name)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error issuing Azure ARM delete request of API Management Service '%s': %s", name, err)
	}

	return nil
}

func expandArmApiManagementServiceSku(d *schema.ResourceData) (apimanagement.ServiceSkuProperties) {
	skus := d.Get("sku").(*schema.Set).List()
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := int32(sku["capacity"].(int))

	return apimanagement.ServiceSkuProperties{
		Name: apimanagement.SkuType(name),
		Capacity: &capacity,
	}
}

func expandArmApiManagementServiceProperties(d *schema.ResourceData) apimanagement.ServiceProperties {
	publishers := d.Get("publisher").(*schema.Set).List()
	publisher := publishers[0].(map[string]interface{})

	publisherName := publisher["name"].(string)
	publisherEmail := publisher["email"].(string)

	return apimanagement.ServiceProperties{
		PublisherName: &publisherName,
		PublisherEmail: &publisherEmail,
		VpnType: apimanagement.None,
	}
}


func flattenArmApiManagementServiceSku(d *schema.ResourceData, properties *apimanagement.ServiceSkuProperties) {
	skus := make([]map[string]interface{}, 1)

	sku := map[string]interface{}{}
	sku["name"] = properties.Name
	sku["capacity"] = properties.Capacity
	skus = append(skus, sku)

	d.Set("sku", &skus)
}
