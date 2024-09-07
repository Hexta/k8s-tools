package nodeutil

import (
	"fmt"
	"sort"
	"time"

	"github.com/rodaine/table"
	"github.com/sirupsen/logrus"
)

func FormatNodeInfo(nodes NodeInfoList) {
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Utilisation.CPU < nodes[j].Utilisation.CPU })

	tbl := table.New("Name", "CPU", "Memory", "Type", "Age").WithWriter(logrus.StandardLogger().Out)
	for _, node := range nodes {
		tbl.AddRow(
			node.Name,
			fmt.Sprintf("%.3f", node.Utilisation.CPU),
			fmt.Sprintf("%.3f", node.Utilisation.Memory),
			fmt.Sprintf("%v", node.InstanceType),
			node.Age.Truncate(time.Hour),
		)
	}

	tbl.Print()
}
