package graphdraw

import (
	"fmt"
	"strings"
)

const padding = 2

type node struct {
	label string
	edges []edge
}

func (n *node) draw() string {
	horizontalWall := strings.Repeat("─", len(n.label)+(padding*2))
	innerPadding := strings.Repeat(" ", padding)

	topWall := fmt.Sprintf("┌%s┐", horizontalWall)
	middleWall := fmt.Sprintf("│%s%s%s│", innerPadding, n.label, innerPadding)
	bottomWall := fmt.Sprintf("└%s┘", horizontalWall)

	return strings.Join([]string{"", topWall, middleWall, bottomWall, ""}, "\n")
}

func (n *node) to(b *node) {

}

type edge struct {
}

func Draw(root node) string {
	return root.draw()
}
