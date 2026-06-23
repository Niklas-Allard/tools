package mdpdf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ── AST node ─────────────────────────────────────────────────────────────────

var KindAlert = ast.NewNodeKind("GithubAlert")

type alertNode struct {
	ast.BaseBlock
	alertType string // "note" | "tip" | "warning" | "caution" | "important"
}

func (n *alertNode) Kind() ast.NodeKind { return KindAlert }
func (n *alertNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{"type": n.alertType}, nil)
}

// ── AST transformer ───────────────────────────────────────────────────────────

var alertHeadRE = regexp.MustCompile(`(?i)^\[!(NOTE|TIP|WARNING|CAUTION|IMPORTANT)\]\s*`)

type alertTransformer struct{}

func (t *alertTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	type replacement struct{ old, neu ast.Node }
	var pending []replacement

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindBlockquote {
			return ast.WalkContinue, nil
		}

		// First child of blockquote must be a paragraph
		firstPara := n.FirstChild()
		if firstPara == nil || firstPara.Kind() != ast.KindParagraph {
			return ast.WalkSkipChildren, nil
		}

		// First child of paragraph must be a text node matching [!TYPE]
		firstText := firstPara.FirstChild()
		if firstText == nil || firstText.Kind() != ast.KindText {
			return ast.WalkSkipChildren, nil
		}
		seg := firstText.(*ast.Text).Segment
		line := string(reader.Value(seg))
		m := alertHeadRE.FindStringSubmatch(line)
		if m == nil {
			return ast.WalkSkipChildren, nil
		}

		typ := strings.ToLower(m[1])
		alert := &alertNode{alertType: typ}

		// Strip [!TYPE] text node from the paragraph
		firstPara.RemoveChild(firstPara, firstText)

		// Drop the paragraph entirely if it is now empty
		if firstPara.FirstChild() == nil {
			n.RemoveChild(n, firstPara)
		}

		// Move all remaining children from blockquote to alert node
		for child := n.FirstChild(); child != nil; {
			sib := child.NextSibling()
			n.RemoveChild(n, child)
			alert.AppendChild(alert, child)
			child = sib
		}

		pending = append(pending, replacement{old: n, neu: alert})
		return ast.WalkSkipChildren, nil
	})

	for _, r := range pending {
		if p := r.old.Parent(); p != nil {
			p.ReplaceChild(p, r.old, r.neu)
		}
	}
}

// ── HTML renderer ─────────────────────────────────────────────────────────────

var alertIcons = map[string]string{
	"note":      "ℹ️",
	"tip":       "💡",
	"warning":   "⚠️",
	"caution":   "🛑",
	"important": "❗",
}

type alertRenderer struct{}

func (r *alertRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindAlert, r.render)
}

func (r *alertRenderer) render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*alertNode)
	if entering {
		icon := alertIcons[n.alertType]
		title := strings.ToUpper(n.alertType[:1]) + n.alertType[1:]
		fmt.Fprintf(w, "<div class=\"admonition %s\"><p class=\"admonition-title\">%s %s</p>\n",
			n.alertType, icon, title)
	} else {
		fmt.Fprint(w, "</div>\n")
	}
	return ast.WalkContinue, nil
}

// ── Extension ─────────────────────────────────────────────────────────────────

type githubAlertsExt struct{}

func (e *githubAlertsExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(util.Prioritized(&alertTransformer{}, 500)),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(util.Prioritized(&alertRenderer{}, 500)),
	)
}

// GithubAlerts is a goldmark extension that renders > [!NOTE] / > [!TIP] etc.
var GithubAlerts goldmark.Extender = &githubAlertsExt{}
