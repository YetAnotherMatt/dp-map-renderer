package renderer_test

import (
	"bytes"
	"testing"

	"fmt"

	"strings"

	. "github.com/ONSdigital/dp-map-renderer/htmlutil"
	"github.com/ONSdigital/dp-map-renderer/models"
	"github.com/ONSdigital/dp-map-renderer/renderer"
	"github.com/ONSdigital/dp-map-renderer/testdata"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestRenderHTMLWithSVG(t *testing.T) {

	Convey("Successfully render an html map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")
		So(GetAttribute(container, "id"), ShouldEqual, renderRequest.Filename+"-map")

		svg := FindNode(container, atom.Svg)
		So(svg, ShouldNotBeNil)

		// the footer - source
		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		// footnotes
		notes := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__notes"})
		So(notes, ShouldNotBeNil)
		So(notes.FirstChild.Data, ShouldResemble, "Notes")
		footnotes := FindNodes(footer, atom.Li)
		So(len(footnotes), ShouldEqual, len(renderRequest.Footnotes))

	})
}

func TestRenderHTMLWithPNG(t *testing.T) {

	Convey("Successfully render a png image of the map", t, func() {

		renderer.UsePNGConverter(pngConverter)

		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}

		container, html := invokeRenderHTMLWithPNG(renderRequest)

		fmt.Println(html)
		So(GetAttribute(container, "class"), ShouldEqual, "figure")
		So(GetAttribute(container, "id"), ShouldEqual, renderRequest.Filename+"-map")

		svg := FindNode(container, atom.Svg)
		So(svg, ShouldBeNil)

		img := FindNode(container, atom.Img)
		So(img, ShouldNotBeNil)
		So(len(GetAttribute(img, "width")), ShouldBeGreaterThan, 0)
		So(len(GetAttribute(img, "height")), ShouldBeGreaterThan, 0)

		// the footer - source
		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		// footnotes
		notes := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__notes"})
		So(notes, ShouldNotBeNil)
		So(notes.FirstChild.Data, ShouldResemble, "Notes")
		footnotes := FindNodes(footer, atom.Li)
		So(len(footnotes), ShouldEqual, len(renderRequest.Footnotes))

	})
}

func TestRenderHTMLWithPNG_ConverterNotAvailable(t *testing.T) {

	Convey("Return the svg version when a png converter is not available", t, func() {

		renderer.UsePNGConverter(nil)

		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.IncludeFallbackPng = false

		container, _ := invokeRenderHTMLWithPNG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")
		So(GetAttribute(container, "id"), ShouldEqual, renderRequest.Filename+"-map")

		svg := FindNode(container, atom.Svg)
		So(svg, ShouldNotBeNil)

		img := FindNode(container, atom.Image)
		So(img, ShouldBeNil)

	})
}

func TestRenderHTML_HorizontalLegend(t *testing.T) {

	Convey("Should render a horizontal legend before the map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.Choropleth.HorizontalLegendPosition = models.LegendPositionBefore
		renderRequest.Choropleth.VerticalLegendPosition = ""

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the legend
		keys := findNodesWithClass(container, atom.Div, "map_key")
		So(len(keys), ShouldEqual, 1)
		key := keys[0]
		So(GetAttribute(key, "class"), ShouldContainSubstring, "horizontal")
		So(GetAttribute(key.NextSibling, "class"), ShouldEqual, "map")

	})

	Convey("Should render a horizontal legend after the map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.Choropleth.HorizontalLegendPosition = models.LegendPositionAfter
		renderRequest.Choropleth.VerticalLegendPosition = ""

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the legend
		keys := findNodesWithClass(container, atom.Div, "map_key")
		So(len(keys), ShouldEqual, 1)
		key := keys[0]
		So(GetAttribute(key, "class"), ShouldContainSubstring, "horizontal")
		So(GetAttribute(key.PrevSibling, "class"), ShouldEqual, "map")

	})
}

func TestRenderHTML_VerticalLegend(t *testing.T) {

	Convey("Should render a vertical legend before the map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.Choropleth.HorizontalLegendPosition = ""
		renderRequest.Choropleth.VerticalLegendPosition = models.LegendPositionBefore

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the legend
		keys := findNodesWithClass(container, atom.Div, "map_key")
		So(len(keys), ShouldEqual, 1)
		key := keys[0]
		So(GetAttribute(key, "class"), ShouldContainSubstring, "vertical")
		So(GetAttribute(key.NextSibling, "class"), ShouldEqual, "map")

	})

	Convey("Should render a vertical legend after the map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.Choropleth.HorizontalLegendPosition = ""
		renderRequest.Choropleth.VerticalLegendPosition = models.LegendPositionAfter

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the legend
		keys := findNodesWithClass(container, atom.Div, "map_key")
		So(len(keys), ShouldEqual, 1)
		key := keys[0]
		So(GetAttribute(key, "class"), ShouldContainSubstring, "vertical")
		So(GetAttribute(key.PrevSibling, "class"), ShouldEqual, "map")

	})
}

func TestRenderHTML_BothLegends(t *testing.T) {

	Convey("Should render a vertical and horizontal legend before the map", t, func() {
		reader := bytes.NewReader(testdata.LoadExampleRequest(t))
		renderRequest, err := models.CreateRenderRequest(reader)
		if err != nil {
			t.Fatal(err)
		}
		renderRequest.Choropleth.HorizontalLegendPosition = models.LegendPositionBefore
		renderRequest.Choropleth.VerticalLegendPosition = models.LegendPositionBefore

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the legend
		keys := findNodesWithClass(container, atom.Div, "map_key")
		So(len(keys), ShouldEqual, 2)
		key := keys[0]
		So(GetAttribute(key, "class"), ShouldContainSubstring, "horizontal")
		So(GetAttribute(key.NextSibling, "class"), ShouldContainSubstring, "vertical")
		So(GetAttribute(key.NextSibling.NextSibling, "class"), ShouldEqual, "map")

	})
}

func TestRenderHTMLWithNoSVG(t *testing.T) {

	Convey("Successfully render an html response when no geography provided", t, func() {
		renderRequest := &models.RenderRequest{
			Filename:  "testname",
			Source:    "source text",
			Footnotes: []string{"Note1", "Note2"},
		}

		container, _ := invokeRenderHTMLWithSVG(renderRequest)

		So(GetAttribute(container, "class"), ShouldEqual, "figure")

		// the footer - source
		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		// source
		source := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__source"})
		So(source, ShouldNotBeNil)
		// footnotes
		notes := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__notes"})
		So(notes, ShouldNotBeNil)
		So(notes.FirstChild.Data, ShouldResemble, "Notes")
		footnotes := FindNodes(footer, atom.Li)
		So(len(footnotes), ShouldEqual, len(renderRequest.Footnotes))

	})
}

func TestRenderHTML_Source(t *testing.T) {

	Convey("A renderRequest without a source should not have a source paragraph", t, func() {
		request := models.RenderRequest{Filename: "myId"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		So(FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__source"}), ShouldBeNil)
	})

	Convey("A renderRequest with a source should have a source paragraph", t, func() {
		request := models.RenderRequest{Filename: "myId", Source: "mySource"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		source := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__source"})
		So(source, ShouldNotBeNil)
		So(source.FirstChild.Data, ShouldResemble, "Source: "+request.Source)
	})

	Convey("A renderRequest with a source link should have a source paragraph with anchor link", t, func() {
		request := models.RenderRequest{Filename: "myId", Source: "mySource", SourceLink: "http://foo/bar"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		source := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__source"})
		So(source, ShouldNotBeNil)
		link := FindNodeWithAttributes(source, atom.A, map[string]string{"href": "http://foo/bar"})
		So(link, ShouldNotBeNil)
		So(link.FirstChild.Data, ShouldResemble, request.Source)
	})
}

func TestRenderHTML_Licence(t *testing.T) {

	Convey("A renderRequest without a licence should not have a licence paragraph", t, func() {
		request := models.RenderRequest{Filename: "myId"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		So(FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__licence"}), ShouldBeNil)
	})

	Convey("A renderRequest with a licence should have a licence paragraph", t, func() {
		request := models.RenderRequest{Filename: "myId", Licence: "© Crown copyright 2015"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		licence := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__licence"})
		So(licence, ShouldNotBeNil)
		So(licence.FirstChild.Data, ShouldResemble, request.Licence)
	})
}

func TestRenderHTML_Footer(t *testing.T) {
	Convey("A renderRequest without footnotes should not have notes paragraph", t, func() {
		request := models.RenderRequest{Filename: "myId"}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)
		So(GetAttribute(footer, "class"), ShouldEqual, "figure__footer")
		So(FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__notes"}), ShouldBeNil)
		So(len(FindNodes(footer, atom.Li)), ShouldBeZeroValue)
	})

	Convey("Footnotes should render as li elements with id", t, func() {
		request := models.RenderRequest{Filename: "myId", Footnotes: []string{"Note1", "Note2"}}
		container, _ := invokeRenderHTMLWithSVG(&request)

		footer := FindNode(container, atom.Footer)
		So(footer, ShouldNotBeNil)

		p := FindNodeWithAttributes(footer, atom.P, map[string]string{"class": "figure__notes"})
		So(p, ShouldNotBeNil)
		So(p.FirstChild.Data, ShouldResemble, "Notes")

		list := FindNode(footer, atom.Ol)
		So(list, ShouldNotBeNil)
		So(GetAttribute(list, "class"), ShouldEqual, "figure__footnotes")
		notes := FindNodes(list, atom.Li)
		So(len(notes), ShouldEqual, len(request.Footnotes))
		for i, note := range request.Footnotes {
			So(GetAttribute(notes[i], "id"), ShouldEqual, fmt.Sprintf("map-%s-note-%d", request.Filename, i+1))
			So(GetAttribute(notes[i], "class"), ShouldEqual, "figure__footnote-item")
			So(strings.Trim(notes[i].FirstChild.Data, " "), ShouldResemble, note)
		}
	})

	Convey("Footnotes should be properly parsed", t, func() {
		request := models.RenderRequest{Filename: "myId", Footnotes: []string{"Note1", "Note2\nOn Two Lines"}}
		_, result := invokeRenderHTMLWithSVG(&request)

		So(result, ShouldContainSubstring, "Note2<br/>On Two Lines")
	})
}

func invokeRenderHTMLWithSVG(renderRequest *models.RenderRequest) (*html.Node, string) {
	response, err := renderer.RenderHTMLWithSVG(renderRequest)
	So(err, ShouldBeNil)
	nodes, err := html.ParseFragment(bytes.NewReader([]byte(response)), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	So(err, ShouldBeNil)
	So(len(nodes), ShouldBeGreaterThanOrEqualTo, 1)
	// the containing container
	node := nodes[0]
	So(node.DataAtom, ShouldEqual, atom.Figure)
	return node, string(response)
}

func invokeRenderHTMLWithPNG(renderRequest *models.RenderRequest) (*html.Node, string) {
	response, err := renderer.RenderHTMLWithPNG(renderRequest)
	So(err, ShouldBeNil)
	nodes, err := html.ParseFragment(bytes.NewReader([]byte(response)), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	So(err, ShouldBeNil)
	So(len(nodes), ShouldBeGreaterThanOrEqualTo, 1)
	// the containing container
	node := nodes[0]
	So(node.DataAtom, ShouldEqual, atom.Figure)
	return node, string(response)
}

func findNodesWithClass(parent *html.Node, a atom.Atom, class string) []*html.Node {
	nodes := FindNodes(parent, a)
	result := make([]*html.Node, 0)
	for _, n := range nodes {
		classAttr := GetAttribute(n, "class")
		for _, c := range strings.Split(classAttr, " ") {
			if c == class {
				result = append(result, n)
			}
		}
	}
	return result
}
