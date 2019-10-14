package geojsontools

func isCounterClockwise(ring [][]float64) bool {
    toppos := 0
    for i := 1; i < len(ring); i++ {
        if ring[i][0] >= ring[toppos][0] {
            toppos = i
        }
    }

    var pointA,pointB,pointC []float64
    pointA = ring[(toppos + len(ring) - 1) % len(ring)]
    pointB = ring[toppos]
    pointC = ring[(toppos + len(ring) + 1) % len(ring)]
    cross := (pointA[0] - pointB[0]) * (pointC[1] - pointB[1]) - (pointA[1] - pointB[1]) * (pointC[0] - pointB[0])

    return (cross < 0)
}

func reverse(ring *[][]float64) {
    for left, right := 0, len(*ring)-1; left < right; left, right = left+1, right-1 {
        (*ring)[left], (*ring)[right] = (*ring)[right], (*ring)[left]
    }
}
