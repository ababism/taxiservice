package metrics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)

func HandleFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Collect metrics
		metricsMap := globalRegistry.getMetrics()
		w := c.Writer

		// Sort map by key(namespace_name)
		keys := make([]string, 0, len(metricsMap))
		for k := range metricsMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// View
		lastNamespace := ""
		// Iter metrics by sorted keys
		for _, k := range keys {
			// if its new a namespace - view namespace header
			curNamespace := metricsMap[k][0].Namespace
			if lastNamespace != curNamespace {
				_, _ = fmt.Fprintln(w, "#", strings.ToUpper(curNamespace))
				_, _ = fmt.Fprintln(w, "")
				lastNamespace = curNamespace
			}

			// view metric header info
			_, _ = fmt.Fprintln(w, "# name:", metricsMap[k][0].Name)
			_, _ = fmt.Fprintln(w, "# description:", metricsMap[k][0].Description)
			_, _ = fmt.Fprintln(w, "# type:", metricsMap[k][0].Type)
			// view every metric
			for _, nM := range metricsMap[k] {
				var lv []string
				// collect labels fro metric
				if len(nM.Labels) > 0 {
					for n, l := range nM.Labels {
						lv = append(lv, fmt.Sprintf(`%s="%s"`, l, nM.LabelsValues[n]))
					}
				}
				// view
				_, _ = fmt.Fprintf(w, "%s{%s} %v\n", nM.Name, strings.Join(lv, ", "), nM.Value)
			}
			// \n after every metric
			_, _ = fmt.Fprintln(w, "")
		}
	}
}
