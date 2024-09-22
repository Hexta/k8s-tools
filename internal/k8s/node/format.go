package node

import (
	"fmt"
	"sort"
	"time"

	"github.com/rodaine/table"
	"github.com/sirupsen/logrus"
)

func FormatNodeInfo(nodes InfoList) {
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].CPUUtilisation < nodes[j].CPUUtilisation })

	tbl := table.New("Name", "CPU", "Memory", "Type", "Age").WithWriter(logrus.StandardLogger().Out)
	for _, node := range nodes {
		tbl.AddRow(
			node.Name,
			fmt.Sprintf("%.3f", node.CPUUtilisation),
			fmt.Sprintf("%.3f", node.MemoryUtilisation),
			fmt.Sprintf("%v", node.InstanceType),
			node.Age.Truncate(time.Hour),
		)
	}

	tbl.Print()
}
