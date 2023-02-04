package confluence

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfluenceAttachment_Created(t *testing.T) {
	rName := acctest.RandomWithPrefix("resource-attachment-test")
	resourceName := "confluence_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConfluenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConfluenceAttachmentConfigRequired(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfluenceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", "file.txt"),
					resource.TestCheckResourceAttr(resourceName, "data", rName),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
				),
			},
		},
	})
}

func testAccCheckConfluenceAttachmentConfigRequired(rName string) string {
	time.Sleep(time.Second)
	return fmt.Sprintf(`
resource confluence_content "default" {
	title = "%s"
	body  = "Original value"
}
resource confluence_attachment "default" {
  title = "file.txt"
	data  = "%s"
	page = confluence_content.default.id
}
`, rName, rName)
}
