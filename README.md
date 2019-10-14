# geojsontools
fix right-hand rule and merge polygons.
英語わかんないので基本日本語で書きます。

## Overview
geojsontoolsは、FeatureCollectionに含まれるFeatureを統合するツールです。
複数のPolygonを結合して、共有する辺を除去したPolygonへと変換できます。

## Important
geojsontoolsは、github.com/paulmach/go.geojson パッケージを前提として動作します。
github.com/paulmach/go.geojson が提供する構造体を入出力に用います。

## Install

	go get github.com/givemegohan/geojsontools

## Import

	import "github.com/givemegohan/geojsontools"

## Example

* ### Fix right-hand rule
右手系ではないFeatureCollectionを、右手系のFeatureCollectionへ修正します。

	fc, _ := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	
	fc = geojsontools.FixRightHandRule(fc)

* ### Merge Polygons
複数のPolygonを結合して、共有する辺を除去したPolygonへと変換します。
入力するFeatureCollectionは、右手系が守られている必要があります。

	fc, _ := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	
	fc = geojsontools.MergePolygons(fc)
