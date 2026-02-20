package sudoku

import (
	"math/rand"
	"time"
)

// Grid 代表 9x9 的数独矩阵
type Grid [9][9]int

const N = 9

// Generate 这是一个对外公开的方法，用于生成一个新的游戏
// difficulty: 挖掉多少个数字 (例如: 简单=30, 困难=50)
// 返回: puzzle(题目), solution(答案)
func Generate(difficulty int) (Grid, Grid) {
	// 设置随机种子，保证每次都不一样
	rand.Seed(time.Now().UnixNano())

	var solution Grid
	// 1. 填充对角线上的 3x3 宫格 (这三个宫格相互独立，随机填不会冲突，能加速后续求解)
	solution.fillDiagonal()

	// 2. 使用回溯法填充剩余格子，生成完整的解
	solution.solve(0, 0)

	// 3. 复制一份作为题目
	puzzle := solution

	// 4. 随机挖洞
	puzzle.removeDigits(difficulty)

	return puzzle, solution
}

// 辅助方法：填充对角线上的 3 个 3x3 宫
func (g *Grid) fillDiagonal() {
	for i := 0; i < N; i = i + 3 {
		g.fillBox(i, i)
	}
}

// 辅助方法：在指定的 3x3 宫内填入随机数字
func (g *Grid) fillBox(rowStart, colStart int) {
	num := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for {
				num = rand.Intn(N) + 1
				if g.isSafeInBox(rowStart, colStart, num) {
					g[rowStart+i][colStart+j] = num
					break
				}
			}
		}
	}
}

// 辅助方法：检查 3x3 宫内是否有重复
func (g *Grid) isSafeInBox(rowStart, colStart, num int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g[rowStart+i][colStart+j] == num {
				return false
			}
		}
	}
	return true
}

// 辅助方法：检查某个位置放入 num 是否合法
func (g *Grid) isSafe(i, j, num int) bool {
	// 检查行
	for c := 0; c < N; c++ {
		if g[i][c] == num {
			return false
		}
	}
	// 检查列
	for r := 0; r < N; r++ {
		if g[r][j] == num {
			return false
		}
	}
	// 检查 3x3 宫
	rowStart := i - i%3
	colStart := j - j%3
	return g.isSafeInBox(rowStart, colStart, num)
}

// 核心算法：递归回溯求解/填充剩余格子
func (g *Grid) solve(row, col int) bool {
	if row == N-1 && col == N {
		return true
	}
	if col == N {
		row++
		col = 0
	}
	if g[row][col] != 0 {
		return g.solve(row, col+1)
	}

	for num := 1; num <= N; num++ {
		if g.isSafe(row, col, num) {
			g[row][col] = num
			if g.solve(row, col+1) {
				return true
			}
			g[row][col] = 0 // 回溯
		}
	}
	return false
}

// 挖洞：随机移除 count 个数字
func (g *Grid) removeDigits(count int) {
	removed := 0
	for removed < count {
		i := rand.Intn(N)
		j := rand.Intn(N)
		if g[i][j] != 0 {
			g[i][j] = 0
			removed++
		}
	}
}
