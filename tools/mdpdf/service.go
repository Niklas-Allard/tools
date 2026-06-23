package mdpdf

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"os"
	"os/exec"
	"regexp"
	"text/template"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/toc"
	"go.abhg.dev/goldmark/wikilink"
)

var abbrPattern = regexp.MustCompile(`(?m)^\*\[([^\]]+)\]:\s*(.+)$`)

func extractAbbreviations(md string) (string, map[string]string) {
	abbrs := map[string]string{}
	cleaned := abbrPattern.ReplaceAllStringFunc(md, func(match string) string {
		parts := abbrPattern.FindStringSubmatch(match)
		if len(parts) == 3 {
			abbrs[parts[1]] = parts[2]
		}
		return ""
	})
	return cleaned, abbrs
}

func applyAbbreviations(htmlBody string, abbrs map[string]string) string {
	for abbr, title := range abbrs {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(abbr) + `\b`)
		escaped := html.EscapeString(title)
		htmlBody = re.ReplaceAllString(htmlBody, `<abbr title="`+escaped+`">`+abbr+`</abbr>`)
	}
	return htmlBody
}

func markdownToHTML(req ConvertRequest) (string, error) {
	mdSrc, abbrs := extractAbbreviations(req.Markdown)

	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			extension.Footnote,
			extension.Typographer,
			extension.DefinitionList,
			GithubAlerts,
			&toc.Extender{},
			&wikilink.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
		),
	)

	if err := md.Convert([]byte(mdSrc), &buf); err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}

	return applyAbbreviations(buf.String(), abbrs), nil
}

var htmlTemplate = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="de">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Document</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
<style>
  :root {
    --color-note: #0969da; --color-tip: #1a7f37; --color-warning: #9a6700;
    --color-caution: #cf222e; --color-important: #8250df;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 16px; line-height: 1.6; color: #1f2328;
    max-width: 900px; margin: 0 auto; padding: 32px 40px;
  }
  h1,h2,h3,h4,h5,h6 { font-weight: 600; line-height: 1.25; margin: 24px 0 16px; }
  h1 { font-size: 2em; padding-bottom: .3em; border-bottom: 1px solid #d1d9e0; }
  h2 { font-size: 1.5em; padding-bottom: .3em; border-bottom: 1px solid #d1d9e0; }
  h3 { font-size: 1.25em; } h4 { font-size: 1em; }
  p { margin: 0 0 16px; }
  a { color: #0969da; text-decoration: none; }
  a:hover { text-decoration: underline; }
  code {
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
    font-size: .85em; background: #f6f8fa; border: 1px solid #d1d9e0;
    border-radius: 6px; padding: .2em .4em;
  }
  pre { background: #f6f8fa; border: 1px solid #d1d9e0; border-radius: 6px; padding: 16px; overflow: auto; margin: 0 0 16px; }
  pre code { background: none; border: none; padding: 0; font-size: .85em; }
  blockquote { padding: 0 1em; color: #57606a; border-left: .25em solid #d1d9e0; margin: 0 0 16px; }
  table { border-collapse: collapse; width: 100%; margin: 0 0 16px; }
  th, td { border: 1px solid #d1d9e0; padding: 6px 13px; }
  th { background: #f6f8fa; font-weight: 600; }
  tr:nth-child(even) { background: #f6f8fa; }
  img { max-width: 100%; }
  hr { border: none; border-top: 1px solid #d1d9e0; margin: 24px 0; }
  ul, ol { padding-left: 2em; margin: 0 0 16px; }
  li + li { margin-top: .25em; }
  dl dt { font-weight: 600; margin-top: 8px; }
  dl dd { margin-left: 2em; }
  abbr { cursor: help; border-bottom: 1px dotted #57606a; text-decoration: none; }
  sup { font-size: .75em; }
  .task-list-item { list-style: none; }
  .task-list-item input { margin-right: .5em; }
  .admonition { border: 1px solid; border-radius: 6px; padding: 12px 16px; margin: 0 0 16px; }
  .admonition-title { font-weight: 600; margin-bottom: 6px; display: flex; align-items: center; gap: 6px; }
  .admonition-title::before { font-style: normal; margin-right: 2px; }
  .admonition.note      { border-color: var(--color-note);      background: #ddf4ff; }
  .admonition.note      .admonition-title { color: var(--color-note); }
  .admonition.note      .admonition-title::before { content: "ℹ️"; }
  .admonition.tip       { border-color: var(--color-tip);       background: #dafbe1; }
  .admonition.tip       .admonition-title { color: var(--color-tip); }
  .admonition.tip       .admonition-title::before { content: "💡"; }
  .admonition.warning   { border-color: var(--color-warning);   background: #fff8c5; }
  .admonition.warning   .admonition-title { color: var(--color-warning); }
  .admonition.warning   .admonition-title::before { content: "⚠️"; }
  .admonition.caution   { border-color: var(--color-caution);   background: #ffebe9; }
  .admonition.caution   .admonition-title { color: var(--color-caution); }
  .admonition.caution   .admonition-title::before { content: "🛑"; }
  .admonition.important { border-color: var(--color-important); background: #fbefff; }
  .admonition.important .admonition-title { color: var(--color-important); }
  .admonition.important .admonition-title::before { content: "❗"; }
  nav#TableOfContents { background: #f6f8fa; border: 1px solid #d1d9e0; border-radius: 6px; padding: 16px 24px; margin: 0 0 24px; display: inline-block; min-width: 200px; }
  nav#TableOfContents::before { content: "Inhalt"; font-weight: 600; display: block; margin-bottom: 8px; }
  nav#TableOfContents ul { margin: 0; padding-left: 1.2em; }
  nav#TableOfContents li { margin: 2px 0; }
  .mermaid { text-align: center; margin: 0 0 16px; }
  @page { size: {{.PageSize}}; margin: 20mm; }
  @media print {
    body { padding: 0; -webkit-print-color-adjust: exact; print-color-adjust: exact; }
    pre { white-space: pre-wrap; }
    a { color: inherit; }
  }
</style>
</head>
<body>
{{.Body}}
<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js"></script>
<script>
renderMathInElement(document.body, {
  delimiters: [
    {left:"$$",right:"$$",display:true},
    {left:"$",right:"$",display:false},
    {left:"\\(",right:"\\)",display:false},
    {left:"\\[",right:"\\]",display:true}
  ]
});
</script>
<script src="https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.min.js"></script>
<script>
// Convert goldmark's <pre><code class="language-mermaid"> to <div class="mermaid">
document.querySelectorAll('code.language-mermaid').forEach(function(el) {
  var div = document.createElement('div');
  div.className = 'mermaid';
  div.textContent = el.textContent;
  var pre = el.closest('pre') || el;
  pre.parentNode.replaceChild(div, pre);
});
mermaid.initialize({ startOnLoad: false, theme: 'default' });
window.__mdReady = false;
mermaid.run().then(function() { window.__mdReady = true; }).catch(function() { window.__mdReady = true; });
</script>
</body>
</html>`))

type templateData struct {
	Body     string
	PageSize string
}

func buildHTML(req ConvertRequest, body string) (string, error) {
	ps := req.PageSize
	if ps == "" {
		ps = "A4"
	}
	var buf bytes.Buffer
	if err := htmlTemplate.Execute(&buf, templateData{Body: body, PageSize: ps}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderHTML returns the full styled HTML document (used by the frontend preview/print path).
func RenderHTML(req ConvertRequest) (string, error) {
	body, err := markdownToHTML(req)
	if err != nil {
		return "", err
	}
	return buildHTML(req, body)
}

func findChromium() (string, error) {
	for _, name := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("chromium not found; please install chromium or google-chrome")
}

// ConvertToPDF renders the Markdown to PDF via a headless Chromium instance.
func ConvertToPDF(req ConvertRequest) ([]byte, error) {
	body, err := markdownToHTML(req)
	if err != nil {
		return nil, err
	}
	fullHTML, err := buildHTML(req, body)
	if err != nil {
		return nil, err
	}

	// Write to a temp file; file:// lets CDN scripts load without CORS issues.
	tmp, err := os.CreateTemp("", "mdpdf-*.html")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err = tmp.WriteString(fullHTML); err != nil {
		return nil, err
	}
	tmp.Close()

	chromePath, err := findChromium()
	if err != nil {
		return nil, err
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("allow-file-access-from-files", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var pdfBuf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+tmp.Name()),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		// Poll until mermaid.run() resolves (max 10 s)
		chromedp.ActionFunc(func(ctx context.Context) error {
			deadline := time.Now().Add(10 * time.Second)
			for time.Now().Before(deadline) {
				var ready bool
				if err := chromedp.Evaluate(`!!window.__mdReady`, &ready).Do(ctx); err != nil {
					return err
				}
				if ready {
					return nil
				}
				time.Sleep(150 * time.Millisecond)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}

	return pdfBuf, nil
}