//
// SPDX-License-Identifier: BSD-3-Clause
//
package printers

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintTable(headers []string, data [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, strings.ToUpper(strings.Join(headers, "\t")))

	for _, row := range data {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
}
