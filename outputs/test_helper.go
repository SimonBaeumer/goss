package outputs

import (
	"github.com/SimonBaeumer/goss/resource"
	"time"
)

func GetExampleTestResult() []resource.TestResult {
	return []resource.TestResult{
		{
			Title:        "my title",
			Duration:     time.Duration(500),
			Successful:   true,
			ResourceType: "resource type",
			ResourceId:   "my resource id",
			Property:     "a property",
			Expected:     []string{"expected"},
		},
	}
}
