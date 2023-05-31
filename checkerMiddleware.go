package main

import (
	"math"
	"strconv"
	"strings"
)

func checkAnswer(Ptype string, data interface{}) bool {
	if strings.ToLower(Ptype) == "sudoku" {
		var dt [][]int

		dt2, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt2 {
			dt3, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []int
			for _, e2 := range dt3 {
				a = append(a, int(e2.(float64)))
			}
			dt = append(dt, a)
		}

		sizeBox := int(math.Round(math.Sqrt(float64(len(dt)))))

		// Vertical Checking
		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				nOccur := false
				for k := 0; k < len(dt[i]); k++ {
					if dt[i][j] == dt[i][k] {
						if nOccur {
							return false
						} else {
							nOccur = true
						}
					}
				}
			}
		}

		// Horizontal Checking
		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				nOccur := false
				for k := 0; k < len(dt); k++ {
					if dt[i][j] == dt[k][j] {
						if nOccur {
							return false
						} else {
							nOccur = true
						}
					}
				}
			}
		}

		// Box checking
		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				nOccur := false
				xBox := int(math.Floor(float64(i) / float64(sizeBox)))
				yBox := int(math.Floor(float64(j) / float64(sizeBox)))

				for k := 0; k < sizeBox; k++ {
					for l := 0; l < sizeBox; l++ {
						if dt[i][j] == dt[k+xBox*sizeBox][l+yBox*sizeBox] {
							if nOccur {
								return false
							} else {
								nOccur = true
							}
						}
					}
				}
			}
		}

		return true
	} else if strings.ToLower(Ptype) == "akari" {
		var dt [][]int

		dt2, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt2 {
			dt3, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []int
			for _, e2 := range dt3 {
				a = append(a, int(e2.(float64)))
			}
			dt = append(dt, a)
		}

		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				if dt[i][j] > 0 {
					lightBulb := 0
					if i > 0 {
						if dt[i-1][j] == -2 {
							lightBulb++
						}
					}
					if j > 0 {
						if dt[i][j-1] == -2 {
							lightBulb++
						}
					}
					if i+1 < len(dt) {
						if dt[i+1][j] == -2 {
							lightBulb++
						}
					}
					if j+1 < len(dt[i]) {
						if dt[i][j+1] == -2 {
							lightBulb++
						}
					}
					if lightBulb != dt[i][j] {
						return false
					}

				}
			}
		}

		var tempBoard [][]bool
		for i := 0; i < len(dt); i++ {
			var a []bool
			for j := 0; j < len(dt[0]); j++ {
				a = append(a, false)
			}
			tempBoard = append(tempBoard, a)
		}

		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[0]); j++ {
				if dt[i][j] == -2 {
					tempBoard[i][j] = true
					for x := i - 1; x >= 0; x-- {
						if dt[x][j] == -2 {
							return false
						} else if dt[x][j] == 0 {
							tempBoard[x][j] = true
						} else {
							break
						}
					}
					for x := i + 1; x < len(dt); x++ {
						if dt[x][j] == -2 {
							return false
						} else if dt[x][j] == 0 {
							tempBoard[x][j] = true
						} else {
							break
						}
					}
					for y := j - 1; y >= 0; y-- {
						if dt[i][y] == -2 {
							return false
						} else if dt[i][y] == 0 {
							tempBoard[i][y] = true
						} else {
							break
						}
					}
					for y := j + 1; y < len(dt[i]); y++ {
						if dt[i][y] == -2 {
							return false
						} else if dt[i][y] == 0 {
							tempBoard[i][y] = true
						} else {
							break
						}
					}
				} else if dt[i][j] != 0 {
					tempBoard[i][j] = true
				}
			}
		}

		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				if !tempBoard[i][j] {
					return false
				}
			}
		}

		return true
	} else if strings.ToLower(Ptype) == "kakuro" {
		var boardDt [][]map[string]interface{}

		dt, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt {
			dt2, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []map[string]interface{}
			for _, e2 := range dt2 {
				dt3, ok := e2.(map[string]interface{})
				if !ok {
					return false
				}

				if dt3["v"] == nil {
					dt3["v"] = float64(-1)
				}

				a = append(a, dt3)
			}
			boardDt = append(boardDt, a)
		}

		for i := 0; i < len(boardDt); i++ {
			for j := 0; j < len(boardDt[i]); j++ {
				if boardDt[i][j]["v"].(float64) == 0 {
					return false
				} else if boardDt[i][j]["v"].(float64) == -2 {
					continue
				} else if boardDt[i][j]["v"].(float64) != -1 {
					for x := i - 1; x > 0; x-- {
						if boardDt[x][j]["v"].(float64) < 0 {
							break
						} else if boardDt[x][j]["v"] == boardDt[i][j]["v"] {
							return false
						}
					}

					for x := i + 1; x < len(boardDt); x++ {
						if boardDt[x][j]["v"].(float64) < 0 {
							break
						} else if boardDt[x][j]["v"] == boardDt[i][j]["v"] {
							return false
						}
					}

					for y := j - 1; y > 0; y-- {
						if boardDt[i][y]["v"].(float64) < 0 {
							break
						} else if boardDt[i][y]["v"] == boardDt[i][j]["v"] {
							return false
						}
					}

					for y := j + 1; y < len(boardDt[i]); y++ {
						if boardDt[i][y]["v"].(float64) < 0 {
							break
						} else if boardDt[i][y]["v"] == boardDt[i][j]["v"] {
							return false
						}
					}
				} else if boardDt[i][j]["v"].(float64) == -1 {
					if boardDt[i][j]["dt1"] != nil {
						vSum := float64(0)
						hSum := float64(0)
						target := boardDt[i][j]["dt1"].(float64)

						for x := i - 1; x > 0; x-- {
							if boardDt[x][j]["v"].(float64) < 0 {
								break
							} else {
								vSum += boardDt[x][j]["v"].(float64)
							}
						}

						for y := j + 1; y < len(boardDt[i]); y++ {
							if boardDt[i][y]["v"].(float64) < 0 {
								break
							} else {
								hSum += boardDt[i][y]["v"].(float64)
							}
						}

						if (target != float64(vSum)) && (target != float64(hSum)) {
							return false
						}
					}

					if boardDt[i][j]["dt2"] != nil {
						vSum := float64(0)
						hSum := float64(0)
						target := boardDt[i][j]["dt2"].(float64)

						for x := i + 1; x < len(boardDt); x++ {
							if boardDt[x][j]["v"].(float64) < 0 {
								break
							} else {
								vSum += boardDt[x][j]["v"].(float64)
							}
						}

						for y := j - 1; y > 0; y-- {
							if boardDt[i][y]["v"].(float64) < 0 {
								break
							} else {
								hSum += boardDt[i][y]["v"].(float64)
							}
						}

						if (target != float64(vSum)) && (target != float64(hSum)) {
							return false
						}
					}

				} else {
					return false
				}
			}
		}
		return true
	} else if strings.ToLower(Ptype) == "slitherlink" {
		var boardDt [][]int
		var hBarDt [][]int
		var vBarDt [][]int

		dt, ok := data.([]interface{})
		if !ok {
			return false
		}

		if len(dt) != 3 {
			return false
		}

		for i := 0; i < 3; i++ {
			dt2, ok := dt[i].([]interface{})
			if !ok {
				return false
			}

			for _, e := range dt2 {
				dt3, ok := e.([]interface{})
				if !ok {
					return false
				}

				var a []int
				for _, e2 := range dt3 {
					a = append(a, int(e2.(float64)))
				}

				if i == 0 {
					boardDt = append(boardDt, a)
				} else if i == 1 {
					vBarDt = append(vBarDt, a)
				} else if i == 2 {
					hBarDt = append(hBarDt, a)
				}
			}
		}

		if len(boardDt)+1 != len(hBarDt) {
			return false
		}
		if len(boardDt) != len(vBarDt) {
			return false
		}
		if len(boardDt[0]) != len(hBarDt[0]) {
			return false
		}
		if len(boardDt[0])+1 != len(vBarDt[0]) {
			return false
		}

		for i := 0; i < len(boardDt); i++ {
			for j := 0; j < len(boardDt[i]); j++ {
				if boardDt[i][j] != -1 {
					bars := 0
					if hBarDt[i][j] != 0 {
						bars++
					}
					if hBarDt[i+1][j] != 0 {
						bars++
					}
					if vBarDt[i][j] != 0 {
						bars++
					}
					if vBarDt[i][j+1] != 0 {
						bars++
					}
					if bars != boardDt[i][j] {
						return false
					}
				}
			}
		}

		for i := 0; i < len(hBarDt); i++ {
			for j := 0; j < len(hBarDt[i]); j++ {
				if hBarDt[i][j] != 0 {
					connection := [2]int{0, 0}

					if i < len(boardDt) {
						if vBarDt[i][j+1] != 0 {
							connection[1]++
						}
					}
					if j < len(boardDt[0])-1 {
						if hBarDt[i][j+1] != 0 {
							connection[1]++
						}
					}
					if i > 0 {
						if vBarDt[i-1][j+1] != 0 {
							connection[1]++
						}
					}

					if i < len(boardDt) {
						if vBarDt[i][j] != 0 {
							connection[0]++
						}
					}
					if j > 0 {
						if hBarDt[i][j-1] != 0 {
							connection[0]++
						}
					}
					if i > 0 {
						if vBarDt[i-1][j] != 0 {
							connection[0]++
						}
					}

					if (connection[0] != 1) || (connection[1] != 1) {
						return false
					}
				}
			}
		}

		for i := 0; i < len(vBarDt); i++ {
			for j := 0; j < len(vBarDt[i]); j++ {
				if vBarDt[i][j] != 0 {
					connection := [2]int{0, 0}
					if j < len(boardDt[i]) {
						if hBarDt[i+1][j] != 0 {
							connection[1]++
						}
					}
					if i < len(boardDt)-1 {
						if vBarDt[i+1][j] != 0 {
							connection[1]++
						}
					}
					if j > 0 {
						if hBarDt[i+1][j-1] != 0 {
							connection[1]++
						}
					}

					if j < len(boardDt[i]) {
						if hBarDt[i][j] != 0 {
							connection[0]++
						}
					}
					if i > 0 {
						if vBarDt[i-1][j] != 0 {
							connection[0]++
						}
					}
					if j > 0 {
						if hBarDt[i][j-1] != 0 {
							connection[0]++
						}
					}

					if (connection[0] != 1) || (connection[1] != 1) {
						return false
					}
				}
			}
		}
		return true
	} else if strings.ToLower(Ptype) == "hashi" {
		var dt [][]int

		dt2, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt2 {
			dt3, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []int
			for _, e2 := range dt3 {
				a = append(a, int(e2.(float64)))
			}
			dt = append(dt, a)
		}

		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				if dt[i][j] > 0 {
					/*--- CIRCLE CONNECTION CHECK ---*/
					connectedLines := 0
					if j > 0 {
						if dt[i][j-1] == -3 {
							connectedLines += 1
						} else if dt[i][j-1] == -4 {
							connectedLines += 2
						}
					}
					if j < len(dt[i])-1 {
						if dt[i][j+1] == -3 {
							connectedLines += 1
						} else if dt[i][j+1] == -4 {
							connectedLines += 2
						}
					}
					if i > 0 {
						if dt[i-1][j] == -1 {
							connectedLines += 1
						} else if dt[i-1][j] == -2 {
							connectedLines += 2
						}
					}
					if i < len(dt)-1 {
						if dt[i+1][j] == -1 {
							connectedLines += 1
						} else if dt[i+1][j] == -2 {
							connectedLines += 2
						}
					}

					if connectedLines != dt[i][j] {
						return false
					}
				} else if dt[i][j] < 0 {
					/*--- LINE CONNECTION CHECK ---*/
					ftx := dt[i][j]
					if ftx > -3 {
						/* Vertical Checking */
						for x := i - 1; x >= -1; x-- {
							if x == -1 {
								return false
							} else if dt[x][j] > 0 {
								break
							} else if dt[x][j] != ftx {
								return false
							}
						}
						for x := i + 1; x <= len(dt); x++ {
							if x == len(dt) {
								return false
							} else if dt[x][j] > 0 {
								break
							} else if dt[x][j] != ftx {
								return false
							}
						}

					} else {
						/* Horizontal Checking */
						for y := j - 1; y >= -1; y-- {
							if y == -1 {
								return false
							} else if dt[i][y] > 0 {
								break
							} else if dt[i][y] != ftx {
								return false
							}
						}
						for y := j + 1; y <= len(dt[i]); y++ {
							if y == len(dt[i]) {
								return false
							} else if dt[i][y] > 0 {
								break
							} else if dt[i][y] != ftx {
								return false
							}
						}
					}
				}
			}
		}

		return true
	} else if strings.ToLower(Ptype) == "futoshiki" {
		var boardDt [][]map[string]interface{}

		dt, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt {
			dt2, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []map[string]interface{}
			for _, e2 := range dt2 {
				dt3, ok := e2.(map[string]interface{})
				if !ok {
					return false
				}

				_, ok = dt3["v"].(float64)
				if !ok {
					return false
				}

				if dt3["u"] == nil {
					dt3["u"] = float64(0)
				} else if dt3["u"] == true {
					dt3["u"] = float64(1)
				} else {
					dt3["u"] = float64(-1)
				}

				if dt3["l"] == nil {
					dt3["l"] = float64(0)
				} else if dt3["l"] == true {
					dt3["l"] = float64(1)
				} else {
					dt3["l"] = float64(-1)
				}

				a = append(a, dt3)
			}
			boardDt = append(boardDt, a)
		}

		sizeBox := float64(len(boardDt))
		if float64(len(boardDt[0])) != sizeBox {
			return false
		}

		for i := 0; i < len(boardDt); i++ {
			for j := 0; j < len(boardDt[i]); j++ {
				if math.IsNaN(boardDt[i][j]["v"].(float64)) {
					return false
				} else {
					if boardDt[i][j]["v"].(float64) > sizeBox {
						return false
					} else if boardDt[i][j]["v"].(float64) < 1 {
						return false
					}

					for x := 0; x < len(boardDt); x++ {
						if x == i {
							continue
						}

						if boardDt[i][j]["v"] == boardDt[x][j]["v"] {
							return false
						}
					}

					for y := 0; y < len(boardDt[i]); y++ {
						if y == j {
							continue
						}

						if boardDt[i][j]["v"] == boardDt[i][y]["v"] {
							return false
						}
					}

					if i > 0 {
						if boardDt[i][j]["u"] != 0 {

							if boardDt[i][j]["u"].(float64) > 0 {
								if boardDt[i-1][j]["v"].(float64) > boardDt[i][j]["v"].(float64) {
									return false
								}
							} else if boardDt[i][j]["u"].(float64) < 0 {
								if boardDt[i-1][j]["v"].(float64) < boardDt[i][j]["v"].(float64) {
									return false
								}
							}
						}
					}

					if j > 0 {
						if boardDt[i][j]["l"] != 0 {

							if boardDt[i][j]["l"].(float64) > 0 {
								if boardDt[i][j-1]["v"].(float64) > boardDt[i][j]["v"].(float64) {
									return false
								}
							} else if boardDt[i][j]["l"].(float64) < 0 {
								if boardDt[i][j-1]["v"].(float64) < boardDt[i][j]["v"].(float64) {
									return false
								}
							}
						}
					}
				}
			}
		}
		return true
	} else if strings.ToLower(Ptype) == "nurikabe" {
		var dt [][]string

		dt2, ok := data.([]interface{})
		if !ok {
			return false
		}

		for _, e := range dt2 {
			dt3, ok := e.([]interface{})
			if !ok {
				return false
			}

			var a []string
			for _, e2 := range dt3 {
				a = append(a, e2.(string))
			}
			dt = append(dt, a)
		}

		seek := func(list [][]int, target []int) bool {
			for i := 0; i < len(list); i++ {
				if (list[i][0] == target[0]) && (list[i][1] == target[1]) {
					return true
				}
			}

			return false
		}

		//---   WHITE ROOM CHECKS   ---
		for i := 0; i < len(dt); i++ {
			for j := 0; j < len(dt[i]); j++ {
				check, err := strconv.Atoi(dt[i][j])
				if err == nil {
					testBox := []([]int){{i, j}}
					count := 0

					for k := 0; k < len(testBox); k++ {
						if dt[testBox[k][0]][testBox[k][1]] == "B" {
							testBox = append(testBox[:k], testBox[k+1:]...)
							k--
							continue
						} else {
							if (dt[testBox[k][0]][testBox[k][1]] != "") && (testBox[k][0] != i) && (testBox[k][1] != j) {
								return false
							}

							count++
							if (testBox[k][0] > 0) && (!seek(testBox, []int{testBox[k][0] - 1, testBox[k][1]})) {
								testBox = append(testBox, []int{testBox[k][0] - 1, testBox[k][1]})
							}
							if (testBox[k][0] < len(dt)-1) && (!seek(testBox, []int{testBox[k][0] + 1, testBox[k][1]})) {
								testBox = append(testBox, []int{testBox[k][0] + 1, testBox[k][1]})
							}
							if (testBox[k][1] > 0) && (!seek(testBox, []int{testBox[k][0], testBox[k][1] - 1})) {
								testBox = append(testBox, []int{testBox[k][0], testBox[k][1] - 1})
							}
							if (testBox[k][1] < len(dt[0])-1) && (!seek(testBox, []int{testBox[k][0], testBox[k][1] + 1})) {
								testBox = append(testBox, []int{testBox[k][0], testBox[k][1] + 1})
							}
						}
					}

					if count != check {
						return false
					}

				}
			}
		}

		//---   CLUSTERING CHECKS   ---
		for i := 0; i < len(dt)-1; i++ {
			for j := 0; j < len(dt[i])-1; j++ {
				if dt[i][j] == "B" {
					if (dt[i+1][j] == "B") && (dt[i][j+1] == "B") && (dt[i+1][j+1] == "B") {
						return false
					}
				}
			}
		}

		//---   WALL CONTINUITY CHECKS   ---
		var wallMap [][]bool
		for i := 0; i < len(dt); i++ {
			var a []bool
			for j := 0; j < len(dt[i]); j++ {
				a = append(a, (dt[i][j] != "B"))
			}
			wallMap = append(wallMap, a)
		}

		for i := 0; i < len(wallMap); i++ {
			for j := 0; j < len(wallMap[i]); j++ {
				if !wallMap[i][j] {

					testBox := []([]int){{i, j}}

					for k := 0; k < len(testBox); k++ {
						if wallMap[testBox[k][0]][testBox[k][1]] {
							continue
						} else {
							wallMap[testBox[k][0]][testBox[k][1]] = true
							if (testBox[k][0] > 0) && (!seek(testBox, []int{testBox[k][0] - 1, testBox[k][1]})) {
								testBox = append(testBox, []int{testBox[k][0] - 1, testBox[k][1]})
							}
							if (testBox[k][0] < len(dt)-1) && (!seek(testBox, []int{testBox[k][0] + 1, testBox[k][1]})) {
								testBox = append(testBox, []int{testBox[k][0] + 1, testBox[k][1]})
							}
							if (testBox[k][1] > 0) && (!seek(testBox, []int{testBox[k][0], testBox[k][1] - 1})) {
								testBox = append(testBox, []int{testBox[k][0], testBox[k][1] - 1})
							}
							if (testBox[k][1] < len(dt[0])-1) && (!seek(testBox, []int{testBox[k][0], testBox[k][1] + 1})) {
								testBox = append(testBox, []int{testBox[k][0], testBox[k][1] + 1})
							}
						}
					}

					for i := 0; i < len(wallMap); i++ {
						for j := 0; j < len(wallMap[i]); j++ {
							if !wallMap[i][j] {
								return false
							}
						}
					}
				}
			}
		}
		return true
	} else {
		return false
	}
}
