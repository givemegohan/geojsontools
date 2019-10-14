package geojsontools

import (
    "github.com/paulmach/go.geojson"
)

type point struct {
    lat float64
    lng float64
}

type edgekey struct {
    p1 point
    p2 point
}

type polygon struct {
    circuits [][][]float64
    outerCircuitId int
}
func MergePolygons(fc *geojson.FeatureCollection) *geojson.FeatureCollection {
    nodes := make(map[point](map[point]struct{}))
    edges := make(map[edgekey](int))

    divCount := 0
    mergeCount := 0
    mergeDict := make(map[int]int)

    /* record node and edge */
    for i := 0; i < len(fc.Features); i++ {
        divId := divCount
        divCount++

        polygon := fc.Features[i].Geometry.Polygon

        mergeDict[divId] = mergeCount
        mergeCount++

        for LRid := 0; LRid < len(polygon); LRid++ {
            p1 := point{polygon[LRid][0][0], polygon[LRid][0][1]}
            for posid := 1; posid < len(polygon[LRid]); posid++ {
                p2 := point{polygon[LRid][posid][0], polygon[LRid][posid][1]}

                revedgekey := edgekey{p2, p1}
                revdivname, ok := edges[revedgekey]
                if ok {
                    /* duplicated edge */

                    edgeMergeId, _ := mergeDict[divId]
                    revedgeMergeId, _ := mergeDict[revdivname]
                    if edgeMergeId != revedgeMergeId {
                        /* merge */
                        mergeDict[divId] = revedgeMergeId
                        for mdivname, mnum := range(mergeDict) {
                            if mnum == edgeMergeId {
                                mergeDict[mdivname] = revedgeMergeId
                            }
                        }
                        if edgeMergeId == (mergeCount - 1) {
                            mergeCount--
                        }
                    }

                    /* delete p2 -> p1 */
                    delete(edges, revedgekey)
                    p2node := nodes[p2]
                    delete(p2node, p1)
                    if len(p2node) == 0 {
                        delete(nodes, p2)
                    } else {
                        nodes[p2] = p2node
                    }

                } else {
                    /* add p1 -> p2 */
                    p1node, ok := nodes[p1]
                    if ok == false {
                        p1node = make(map[point]struct{})
                    }
                    p1node[p2] = struct{}{}
                    nodes[p1] = p1node
    
                    edgekey := edgekey{p1, p2}
                    edges[edgekey] = divId
                }
                
                p1 = p2
            }
        }
    }

    /* create circuit */
    polygons := make(map[int](polygon))
    for ;len(edges) > 0 ; {
        var workedge edgekey
        var circuitNum int

        /* get first node */
        for edgekey, divname := range(edges) {
            workedge = edgekey
            circuitNum = mergeDict[divname]
            break
        }

        polygon, ok := polygons[circuitNum]
        if ok == false {
            polygon.circuits = make([][][]float64,0)
            polygon.outerCircuitId = -1
        }

        var firstpoint = workedge.p1

        circuit := make([][]float64,0)
        circuit = append(circuit, []float64{workedge.p1.lat, workedge.p1.lng})
        for ;; {
            delete(edges, workedge)
            circuit = append(circuit, []float64{workedge.p2.lat, workedge.p2.lng})

            if workedge.p2 == firstpoint {
                break
            }

            /* search edge */
            findflag := false
            for nextpos, _ := range(nodes[workedge.p2]) {
                nextedge := edgekey{workedge.p2, nextpos}
                nextmergeNum := mergeDict[edges[nextedge]]
                if nextmergeNum == circuitNum {
                    findflag = true
                    workedge = nextedge
                    break
                }
            }
            if findflag == false {
                panic("circuit error")
            }
        }

        polygon.circuits = append(polygon.circuits, circuit)
        polygons[circuitNum] = polygon
    }
   
    /* classify circuit as ring or hole */
    for circuitNum,polygon := range(polygons) {
        if len(polygon.circuits) == 1 {
            polygon.outerCircuitId = 0
        } else {
            for i := 0 ; i < len(polygon.circuits) ; i++ {
                if (isCounterClockwise(polygon.circuits[i])) {
                    polygon.outerCircuitId = i
                    break
                }
            }
        }
        if polygon.outerCircuitId == -1 {
            panic("outerCircuit is not found")
        }
        polygons[circuitNum] = polygon
    }

    /* make FeatureCollection */
    fc.Features = make([]*geojson.Feature,0)
    for _,polygon := range(polygons) {
        convertpolygon := make([][][]float64,0)

        /* ring */
        convertcircuit := make([][]float64,0)
        for j := 0; j < len(polygon.circuits[polygon.outerCircuitId]); j++ {
            convertcircuit = append(convertcircuit, polygon.circuits[polygon.outerCircuitId][j])
        }
        convertpolygon = append(convertpolygon, convertcircuit)

        /* hole */
        for i := 0 ; i < len(polygon.circuits) ; i++ {
            if i == polygon.outerCircuitId {
                continue
            }
            convertcircuit := make([][]float64,0)
            for j := 0; j < len(polygon.circuits[i]); j++ {
                convertcircuit = append(convertcircuit, polygon.circuits[i][j])
            }
            convertpolygon = append(convertpolygon, convertcircuit)
        }
        fc.AddFeature(geojson.NewPolygonFeature(convertpolygon))
    }

    return fc
}
