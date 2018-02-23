package main

import(
	"fmt"
	"math"
	"sort"
	"os"
	"bufio"
)

type Point struct{
	x, y float64
}

type byX []*Point
type byY []*Point

func main() {
	streamReader := bufio.NewReader(os.Stdin)
	for{
		var pointsNum int
		fmt.Fscanf(streamReader, "%d\n", &pointsNum)

		points := make([]*Point, 0, pointsNum)

		if pointsNum == 0{
			return			//Only way to exit loop
		}

		//clear points
		points = points[0:0]

		for i := 0; i < pointsNum; i++{
			var x, y float64
			fmt.Fscanf(streamReader, "%f %f\n", &x, &y)

			points = append(points, &Point{x, y})
		}

		sort.Sort(byX(points))
		a, b, _ := findMinDist(points)
		fmt.Println(a, b)
	}
}

func findMinDist(localPoints []*Point) (*Point, *Point, float64){
	localPointsLength := len(localPoints)
	if localPointsLength == 2{
		p1 := localPoints[0]
		p2 := localPoints[1]
		return p1, p2, calcSquareDist(p1, p2)
	}else if localPointsLength == 3{
		p1 := localPoints[0]
		p2 := localPoints[1]
		p3 := localPoints[2]

		dist1 := calcSquareDist(p1, p2)
		dist2 := calcSquareDist(p2, p3)
		dist3 := calcSquareDist(p3, p1)

		if dist1 < dist2 && dist1 < dist3{
			return p1, p2, dist1
		}else if dist2 < dist3{
			return p2, p3, dist2
		}else{
			return p3, p1, dist3
		}
	}

	median := localPointsLength / 2
	medianLocalPoint := localPoints[median]
	medianXVal := medianLocalPoint.x
	p1, p2, squareDist1 := findMinDist(localPoints[:median])	//Left Half
	p3, p4, squareDist2 := findMinDist(localPoints[median:])	//Right Half, middle point contained here

	var p1Ans *Point
	var p2Ans *Point
	var squareDistanceAns float64
	var delta float64
	if squareDist1 < squareDist2{
		p1Ans = p1
		p2Ans = p2
		squareDistanceAns = squareDist1
		delta = math.Sqrt(squareDist1)
	}else{
		p1Ans = p3
		p2Ans = p4
		squareDistanceAns = squareDist2
		delta = math.Sqrt(squareDist2)
	}

	// ================ Middle Slab ================
	//define middle slab bounds
	leftBound := medianXVal - delta
	rightBound := medianXVal + delta

	//search backwards through the local points, adds points to leftSlabPoints if within leftBound
	leftBoundIndex := median - 1
	for ; leftBoundIndex >= 0 && localPoints[leftBoundIndex].x > leftBound; leftBoundIndex-- {}
	leftBoundIndex++        //for loop above overshoots, need to increment index to include point within bounds
	leftSlabPoints := localPoints[leftBoundIndex:median]

	//search forward through the local points, adds points to rightSlabPoints if within rightBound
	rightBoundIndex := median
	for ; rightBoundIndex < localPointsLength && localPoints[rightBoundIndex].x < rightBound; rightBoundIndex++ {}
	rightSlabPoints := localPoints[median:rightBoundIndex]

	sort.Sort(byY(leftSlabPoints))
	sort.Sort(byY(rightSlabPoints))

	testPointDistance := func(left, right *Point) {
		if right.x - left.x < delta {
			sqDist := calcSquareDist(left, right)
			if sqDist < squareDistanceAns {
				delta = math.Sqrt(sqDist)
				squareDistanceAns = sqDist
				p1Ans = left
				p2Ans = right
			}
		}
	}
	leftIndex := 0
	rightIndex := 0
	for leftIndex < len(leftSlabPoints) && rightIndex < len(rightSlabPoints) {
		leftPoint := leftSlabPoints[leftIndex]
		rightPoint := rightSlabPoints[rightIndex]
		if leftPoint.y <= rightPoint.y {                        //Right point is above left point
			if rightPoint.y - leftPoint.y < delta {
				testPointDistance(leftPoint, rightPoint)        //Don't test distance if vertical distance is greater than delta
			}
			if nextRight := rightIndex+1; nextRight < len(rightSlabPoints) {
				testPointDistance(leftPoint, rightSlabPoints[nextRight])
			}
			leftIndex++
		} else {                                                //Left point is above right point
			if leftPoint.y - rightPoint.y < delta {
				testPointDistance(leftPoint, rightPoint)        //Don't test distance if vertical distance is greater than delta
			}
			if nextLeft := leftIndex+1; nextLeft < len(leftSlabPoints) {
				testPointDistance(leftSlabPoints[nextLeft], rightPoint)
			}
			rightIndex++
		}
	}

	return p1Ans, p2Ans, squareDistanceAns
}

//Computes the distance squared between two points
func calcSquareDist(a, b *Point) float64{
	return (a.x - b.x) * (a.x - b.x) + (a.y - b.y) * (a.y - b.y)
}

//for interfaces
func (a byX) Len() int{
	return len(a)
}

func (a byX) Less(i, j int) bool{
	return a[i].x < a[j].x
}

func (a byX) Swap(i, j int){
	a[i], a[j] = a[j], a[i]
}


func (a byY) Len() int{
	return len(a)
}

func (a byY) Less(i, j int) bool{
	return a[i].y < a[j].y
}

func (a byY) Swap(i, j int){
	a[i], a[j] = a[j], a[i]
}


func (p *Point) String() string {
	return fmt.Sprintf("%.2f %.2f", p.x, p.y)
}
