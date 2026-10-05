package pptx

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

// buildContentTypesXML generates [Content_Types].xml for the PPTX package.
func buildContentTypesXML(slideCount int, hasNotes []bool) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` + "\n")
	sb.WriteString(`  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` + "\n")
	sb.WriteString(`  <Default Extension="xml" ContentType="application/xml"/>` + "\n")
	sb.WriteString(`  <Default Extension="png" ContentType="image/png"/>` + "\n")
	sb.WriteString(`  <Default Extension="jpeg" ContentType="image/jpeg"/>` + "\n")
	sb.WriteString(`  <Default Extension="jpg" ContentType="image/jpeg"/>` + "\n")
	sb.WriteString(`  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` + "\n")

	for i := 1; i <= slideCount; i++ {
		sb.WriteString(fmt.Sprintf(`  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`+"\n", i))
		if i-1 < len(hasNotes) && hasNotes[i-1] {
			sb.WriteString(fmt.Sprintf(`  <Override PartName="/ppt/notesSlides/notesSlide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>`+"\n", i))
		}
	}
	sb.WriteString(`</Types>`)
	return []byte(sb.String())
}

// buildRootRelsXML generates _rels/.rels linking to ppt/presentation.xml.
func buildRootRelsXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`)
}

// buildPresentationRelsXML generates ppt/_rels/presentation.xml.rels linking each slide.
func buildPresentationRelsXML(slideCount int) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + "\n")
	for i := 1; i <= slideCount; i++ {
		sb.WriteString(fmt.Sprintf(`  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`+"\n", i, i))
	}
	sb.WriteString(`</Relationships>`)
	return []byte(sb.String())
}

// buildPresentationXML generates ppt/presentation.xml specifying 16:9 widescreen dimensions (12192000x6858000 EMU).
func buildPresentationXML(slideCount int) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` + "\n")
	sb.WriteString(`  <p:sldMasterIdLst/>` + "\n")
	sb.WriteString(`  <p:sldIdLst>` + "\n")
	for i := 1; i <= slideCount; i++ {
		sb.WriteString(fmt.Sprintf(`    <p:sldId id="%d" r:id="rId%d"/>`+"\n", 255+i, i))
	}
	sb.WriteString(`  </p:sldIdLst>` + "\n")
	// 16:9 widescreen: 12192000 x 6858000 EMU
	sb.WriteString(`  <p:sldSz cx="12192000" cy="6858000" type="screen16x9"/>` + "\n")
	sb.WriteString(`  <p:notesSz cx="6858000" cy="9144000"/>` + "\n")
	sb.WriteString(`</p:presentation>`)
	return []byte(sb.String())
}

// SlideLink defines an interactive hyperlink overlay geometry on a slide (in EMU).
type SlideLink struct {
	TargetURL string `json:"url"`
	X         int64  `json:"x"`
	Y         int64  `json:"y"`
	W         int64  `json:"w"`
	H         int64  `json:"h"`
}

// buildSlideXML generates ppt/slides/slideN.xml mapping fullscreen 16:9 image and clickable hyperlink overlays.
func buildSlideXML(links []SlideLink) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="0" cy="0"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="0" cy="0"/>
        </a:xfrm>
      </p:grpSpPr>
      <p:pic>
        <p:nvPicPr>
          <p:cNvPr id="2" name="Slide Image"/>
          <p:cNvPicPr>
            <a:picLocks noChangeAspect="1"/>
          </p:cNvPicPr>
          <p:nvPr/>
        </p:nvPicPr>
        <p:blipFill>
          <a:blip r:embed="rId1"/>
          <a:stretch>
            <a:fillRect/>
          </a:stretch>
        </p:blipFill>
        <p:spPr>
          <a:xfrm>
            <a:off x="0" y="0"/>
            <a:ext cx="12192000" cy="6858000"/>
          </a:xfrm>
          <a:prstGeom prst="rect">
            <a:avLst/>
          </a:prstGeom>
        </p:spPr>
      </p:pic>` + "\n")

	for idx, link := range links {
		spID := 10 + idx
		relID := fmt.Sprintf("rIdLink%d", idx+1)
		sb.WriteString(fmt.Sprintf(`      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="%d" name="Hyperlink Overlay %d">
            <a:hlinkClick r:id="%s"/>
          </p:cNvPr>
          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
          <p:nvPr/>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="%d" y="%d"/>
            <a:ext cx="%d" cy="%d"/>
          </a:xfrm>
          <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
          <a:solidFill>
            <a:srgbClr val="FFFFFF"><a:alpha val="0"/></a:srgbClr>
          </a:solidFill>
          <a:ln><a:noFill/></a:ln>
        </p:spPr>
      </p:sp>`+"\n", spID, idx+1, relID, link.X, link.Y, link.W, link.H))
	}

	sb.WriteString(`    </p:spTree>
  </p:cSld>
  <p:clrMapOvr>
    <a:masterClrMapping/>
  </p:clrMapOvr>
</p:sld>`)
	return []byte(sb.String())
}

// buildSlideRelsXML generates ppt/slides/_rels/slideN.xml.rels linking media image, optional notes, and hyperlinks.
func buildSlideRelsXML(slideIdx int, hasNotes bool, links []SlideLink) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + "\n")
	sb.WriteString(fmt.Sprintf(`  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/slide%d.png"/>`+"\n", slideIdx))
	if hasNotes {
		sb.WriteString(fmt.Sprintf(`  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" Target="../notesSlides/notesSlide%d.xml"/>`+"\n", slideIdx))
	}
	for idx, link := range links {
		escapedTarget := html.EscapeString(link.TargetURL)
		sb.WriteString(fmt.Sprintf(`  <Relationship Id="rIdLink%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="%s" TargetMode="External"/>`+"\n", idx+1, escapedTarget))
	}
	sb.WriteString(`</Relationships>`)
	return []byte(sb.String())
}

// buildNotesSlideXML generates ppt/notesSlides/notesSlideN.xml preserving speaker notes as text.
func buildNotesSlideXML(notesText string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	buf.WriteString(`<p:notes xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` + "\n")
	buf.WriteString(`  <p:cSld>` + "\n")
	buf.WriteString(`    <p:spTree>` + "\n")
	buf.WriteString(`      <p:nvGrpSpPr>` + "\n")
	buf.WriteString(`        <p:cNvPr id="1" name=""/>` + "\n")
	buf.WriteString(`        <p:cNvGrpSpPr/>` + "\n")
	buf.WriteString(`        <p:nvPr/>` + "\n")
	buf.WriteString(`      </p:nvGrpSpPr>` + "\n")
	buf.WriteString(`      <p:grpSpPr>` + "\n")
	buf.WriteString(`        <a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/><a:chOff x="0" y="0"/><a:chExt cx="0" cy="0"/></a:xfrm>` + "\n")
	buf.WriteString(`      </p:grpSpPr>` + "\n")
	buf.WriteString(`      <p:sp>` + "\n")
	buf.WriteString(`        <p:nvSpPr>` + "\n")
	buf.WriteString(`          <p:cNvPr id="2" name="Notes Text"/>` + "\n")
	buf.WriteString(`          <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>` + "\n")
	buf.WriteString(`          <p:nvPr><p:ph type="body" idx="1"/></p:nvPr>` + "\n")
	buf.WriteString(`        </p:nvSpPr>` + "\n")
	buf.WriteString(`        <p:spPr/>` + "\n")
	buf.WriteString(`        <p:txBody>` + "\n")
	buf.WriteString(`          <a:bodyPr/>` + "\n")
	buf.WriteString(`          <a:lstStyle/>` + "\n")

	lines := strings.Split(notesText, "\n")
	for _, line := range lines {
		escaped := html.EscapeString(line)
		buf.WriteString(fmt.Sprintf(`          <a:p><a:r><a:rPr lang="ko-KR" altLang="en-US" dirty="0"/><a:t>%s</a:t></a:r></a:p>`+"\n", escaped))
	}

	buf.WriteString(`        </p:txBody>` + "\n")
	buf.WriteString(`      </p:sp>` + "\n")
	buf.WriteString(`    </p:spTree>` + "\n")
	buf.WriteString(`  </p:cSld>` + "\n")
	buf.WriteString(`  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>` + "\n")
	buf.WriteString(`</p:notes>`)
	return buf.Bytes()
}

// buildNotesSlideRelsXML generates ppt/notesSlides/_rels/notesSlideN.xml.rels linking back to slideN.xml.
func buildNotesSlideRelsXML(slideIdx int) []byte {
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="../slides/slide%d.xml"/>
</Relationships>`, slideIdx))
}
