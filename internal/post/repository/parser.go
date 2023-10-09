package repository

import (
	"bytes"
	"io"

	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/post"
	"github.com/zanz1n/go-htmx/internal/utils"
	"golang.org/x/net/html"
)

func NewHtmlParser(removeScripts bool, removeStyles bool) *HtmlParser {
	return &HtmlParser{
		removeScripts: removeScripts,
		removeStyles:  removeStyles,
	}
}

type HtmlParser struct {
	removeScripts bool
	removeStyles  bool
}

func (p *HtmlParser) RenderTree(ele *html.Node) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 128))

	if err := html.Render(buf, ele); err != nil {
		return nil, errors.ErrInvalidHtmlFragment
	}

	return buf.Bytes(), nil
}

func (p *HtmlParser) SanitizePostFragment(r io.Reader) ([]byte, []post.PostTopic, error) {
	frag, err := p.ParseFragment(r)
	if err != nil {
		return nil, nil, err
	}
	frag = p.SanitizeTree(frag)

	buf, err := p.RenderTree(frag)
	if err != nil {
		return nil, nil, err
	}

	pt := p.ExtractPostTitles(frag.FirstChild)

	return buf, pt, nil
}

func (p *HtmlParser) ParseFragment(r io.Reader) (*html.Node, error) {
	tree, err := html.Parse(r)
	if err != nil {
		return nil, errors.ErrInvalidHtmlFragment
	}

	tree = p.SanitizeTree(tree)

	if tree == nil || tree.FirstChild == nil || tree.FirstChild.FirstChild == nil {
		return nil, errors.ErrInvalidHtmlFragment
	}

	var body *html.Node = tree.FirstChild.FirstChild
	for body != nil {
		if body.Type == html.ElementNode && body.Data == "body" {
			break
		}
		body = body.NextSibling
	}

	if body == nil {
		return nil, errors.ErrInvalidHtmlFragment
	}

	body.Attr = nil
	body.Data = ""
	body.DataAtom = 0
	body.Namespace = ""
	body.Parent = nil
	body.Type = html.DocumentNode

	return body, nil
}

func (p *HtmlParser) SanitizeTree(frag *html.Node) *html.Node {
	if frag.Type == html.ElementNode {
		switch frag.Data {
		case "script":
			if p.removeScripts {
				return nil
			}
		case "style":
			if p.removeStyles {
				return nil
			}
		case "img", "audio", "video":
			frag.Attr = p.sanitizeNodeAttributes(frag.Attr)
		default:
			frag.Attr = nil
		}
	} else if frag.Type == html.TextNode {
		frag.Data = p.sanitizeText(frag.Data)
	}

	if frag.FirstChild != nil {
		frag.FirstChild = p.SanitizeTree(frag.FirstChild)
	}
	if frag.NextSibling != nil {
		frag.NextSibling = p.SanitizeTree(frag.NextSibling)
	}

	return frag
}

func (p *HtmlParser) ExtractPostTitles(frag *html.Node) []post.PostTopic {
	ts := []post.PostTopic{}

	var current = frag

	for current != nil {
		if current.Type == html.ElementNode {
			switch current.Data {
			case "h2", "h3", "h4", "h5":
				if current.FirstChild != nil && current.FirstChild.Type == html.TextNode {
					ts = append(ts, post.PostTopic{
						Kind:  post.PostTopicKind(current.Data),
						Title: current.FirstChild.Data,
					})
				}
			default:
			}
		}

		current = current.NextSibling
	}

	return ts
}

func (p *HtmlParser) sanitizeText(text string) string {
	b := []byte(text)

	s := bytes.Split(b, []byte{'\n'})
	for i, line := range s {
		if len(line) > 0 {
			if line[0] == ' ' || line[0] == '\t' {
				newline := []byte{}
				identEnd := false
				for _, char := range line {
					if !identEnd {
						if char == '\t' || char == ' ' {
							continue
						} else {
							identEnd = true
						}
					}
					newline = append(newline, char)
				}
				s[i] = newline
			}
		}
	}

	return utils.B2S(bytes.Join(s, []byte{}))
}

func (p *HtmlParser) sanitizeNodeAttributes(attrs []html.Attribute) []html.Attribute {
	nattrs := []html.Attribute{}

	for _, attr := range attrs {
		switch attr.Key {
		case "style", "class", "id":
			if p.removeStyles {
				continue
			}
		case "src", "alt":
		default:
			continue
		}

		nattrs = append(nattrs, attr)
	}

	return nattrs
}
