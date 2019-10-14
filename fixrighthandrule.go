package geojsontools

import (
    "github.com/paulmach/go.geojson"
)

func FixRightHandRule(fc *geojson.FeatureCollection) *geojson.FeatureCollection {
    for i := 0; i < len(fc.Features) ; i++ {        
        if fc.Features[i].Geometry.Polygon == nil || len(fc.Features[i].Geometry.Polygon) == 0 {
            continue
        }
    
        if (!isCounterClockwise(fc.Features[i].Geometry.Polygon[0])) {
            reverse(&fc.Features[i].Geometry.Polygon[0])
        }
    
        for j := 1; j < len(fc.Features[i].Geometry.Polygon); j++ {
            if (isCounterClockwise(fc.Features[i].Geometry.Polygon[j])) {
                reverse(&fc.Features[i].Geometry.Polygon[j])
            }
        }
    }
    return fc
}
