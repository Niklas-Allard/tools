<template>
  <div class="flex flex-col h-[calc(100vh-5rem)] gap-0">
    <!-- Toolbar -->
    <div class="flex items-center gap-3 px-4 py-2 border-b bg-card shrink-0 flex-wrap">
      <h1 class="font-semibold text-sm mr-2">📄 Markdown → PDF</h1>

      <div class="flex items-center gap-1.5">
        <Label class="text-xs text-muted-foreground">Theme</Label>
        <Select v-model:model-value="theme">
          <SelectTrigger class="h-7 text-xs w-28"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="github">GitHub</SelectItem>
            <SelectItem value="minimal">Minimal</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center gap-1.5">
        <Label class="text-xs text-muted-foreground">Format</Label>
        <Select v-model:model-value="pageSize">
          <SelectTrigger class="h-7 text-xs w-20"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="A4">A4</SelectItem>
            <SelectItem value="Letter">Letter</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="ml-auto flex items-center gap-2">
        <Button
          @click="loadExample"
          variant="outline"
          size="sm"
          class="h-7 text-xs"
        >
          Beispiel laden
        </Button>
        <Button
          @click="downloadPDF"
          size="sm"
          class="h-7 text-xs bg-emerald-600 hover:bg-emerald-700 text-white"
          :disabled="!markdown.trim() || pdfLoading"
        >
          <span v-if="pdfLoading">⏳ Generiere…</span>
          <span v-else>⬇ PDF herunterladen</span>
        </Button>
      </div>
    </div>

    <!-- Split pane -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Editor -->
      <div class="flex flex-col w-1/2 border-r">
        <div class="px-3 py-1 text-xs text-muted-foreground bg-muted/40 border-b font-medium">Editor</div>
        <textarea
          v-model="markdown"
          class="flex-1 font-mono text-sm p-4 resize-none focus:outline-none bg-background text-foreground overflow-auto"
          placeholder="Markdown hier eingeben…"
          spellcheck="false"
          wrap="off"
        />
      </div>

      <!-- Preview -->
      <div class="flex flex-col w-1/2 overflow-hidden">
        <div class="px-3 py-1 text-xs text-muted-foreground bg-muted/40 border-b font-medium">Vorschau</div>
        <div
          ref="previewEl"
          class="flex-1 overflow-auto px-8 py-6 prose prose-sm max-w-none preview-content"
          v-html="renderedHTML"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import taskLists from 'markdown-it-task-lists'
import githubAlerts from 'markdown-it-github-alerts'
import toc from 'markdown-it-toc-done-right'
import anchor from 'markdown-it-anchor'
import abbr from 'markdown-it-abbr'
import markdownItKatex from '@traptitech/markdown-it-katex'
import hljs from 'highlight.js'
import mermaid from 'mermaid'
import axios from 'axios'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import 'highlight.js/styles/github.css'
import 'katex/dist/katex.min.css'
import 'markdown-it-github-alerts/styles/github-base.css'
import 'markdown-it-github-alerts/styles/github-colors-light.css'

// ── Wikilink plugin ──────────────────────────────────────────────────────────
function wikilinkPlugin(md: MarkdownIt) {
  md.core.ruler.push('wikilinks', (state) => {
    for (const token of state.tokens) {
      if (token.type !== 'inline' || !token.children) continue
      for (const child of token.children) {
        if (child.type !== 'text') continue
        child.content = child.content.replace(/\[\[([^\]]+)\]\]/g, (_, label) => {
          return `<a href="#" class="wikilink">[[${label}]]</a>`
        })
        if (child.content !== child.content) {
          child.type = 'html_inline'
        }
      }
    }
  })
}

// ── markdown-it setup ────────────────────────────────────────────────────────
const escapeHtml = (s: string) => MarkdownIt().utils.escapeHtml(s)

const md: MarkdownIt = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight(code: string, lang: string): string {
    if (lang === 'mermaid') {
      // Emit a plain div; mermaid.run() picks it up after nextTick
      return `<div class="mermaid">${escapeHtml(code)}</div>`
    }
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(code, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch {}
    }
    return `<pre class="hljs"><code>${escapeHtml(code)}</code></pre>`
  },
})
  .use(taskLists, { enabled: true, label: true })
  .use(githubAlerts)
  .use(anchor, { permalink: anchor.permalink.headerLink() })
  .use(toc)
  .use(abbr)
  .use(markdownItKatex, { throwOnError: false })
  .use(wikilinkPlugin)

mermaid.initialize({ startOnLoad: false, theme: 'default' })

// ── State ────────────────────────────────────────────────────────────────────
const markdown = ref('')
const renderedHTML = ref('')
const previewEl = ref<HTMLElement | null>(null)
const theme = ref('github')
const pageSize = ref('A4')
const pdfLoading = ref(false)

// ── Render ───────────────────────────────────────────────────────────────────
let renderTimer: ReturnType<typeof setTimeout>

watch(markdown, () => {
  clearTimeout(renderTimer)
  renderTimer = setTimeout(renderPreview, 300)
}, { immediate: true })

async function renderPreview() {
  renderedHTML.value = md.render(markdown.value)
  await nextTick()
  if (!previewEl.value) return
  const nodes = previewEl.value.querySelectorAll<HTMLElement>('.mermaid')
  if (!nodes.length) return
  try {
    await mermaid.run({ nodes })
  } catch (e) {
    console.warn('Mermaid render error:', e)
  }
}

// ── PDF Download ─────────────────────────────────────────────────────────────
async function downloadPDF() {
  if (!markdown.value.trim()) return
  pdfLoading.value = true
  const payload = { markdown: markdown.value, theme: theme.value, page_size: pageSize.value }
  try {
    // Primary: server-side PDF via headless Chromium
    const res = await axios.post('/api/mdpdf/convert', payload, { responseType: 'blob' })
    const url = URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }))
    const a = document.createElement('a')
    a.href = url
    a.download = 'document.pdf'
    a.click()
    URL.revokeObjectURL(url)
    toast.success('PDF heruntergeladen')
  } catch (err: any) {
    // Fallback: browser print (works when Chromium is not installed on server)
    const serverError: string = err.response
      ? await (err.response.data as Blob).text()
      : ''
    const errObj = (() => { try { return JSON.parse(serverError) } catch { return {} } })()
    if (errObj.error?.includes('chromium not found')) {
      await printViaBrowser(payload)
    } else {
      toast.error('PDF-Fehler', { description: errObj.error ?? err.message })
    }
  } finally {
    pdfLoading.value = false
  }
}

async function printViaBrowser(payload: { markdown: string; theme: string; page_size: string }) {
  try {
    const res = await axios.post('/api/mdpdf/render', payload)
    const iframe = document.createElement('iframe')
    iframe.style.cssText = 'position:fixed;left:-9999px;top:-9999px;width:1px;height:1px;border:none;'
    document.body.appendChild(iframe)
    const doc = iframe.contentDocument!
    doc.open()
    doc.write(res.data)
    doc.close()

    // Wait for Mermaid + KaTeX (max 10 s)
    const deadline = Date.now() + 10_000
    await new Promise<void>((resolve) => {
      const check = () => {
        if ((iframe.contentWindow as any)?.__mdReady || Date.now() > deadline) return resolve()
        setTimeout(check, 150)
      }
      setTimeout(check, 400)
    })

    iframe.contentWindow!.print()
    setTimeout(() => document.body.removeChild(iframe), 2000)
    toast.info('Druckdialog geöffnet', { description: 'Chromium nicht gefunden – bitte "Als PDF speichern" wählen.' })
  } catch (e: any) {
    toast.error('Fehler', { description: e.message })
  }
}

// ── Example ──────────────────────────────────────────────────────────────────
function loadExample() {
  markdown.value = `# Markdown → PDF Demo

[[TOC]]

## Textformatierung

**Fett**, *kursiv*, ~~durchgestrichen~~, \`inline code\`

> [!NOTE]
> Das ist eine **Note**-Admonition.

> [!WARNING]
> Vorsicht! Das ist eine Warnung.

> [!TIP]
> Nützlicher Tipp hier.

## Mathematik (KaTeX)

Inline: $E = mc^2$ und $\\alpha + \\beta = \\gamma$

Block:
$$
\\int_0^\\infty e^{-x^2} dx = \\frac{\\sqrt{\\pi}}{2}
$$

## Mermaid-Diagramm

\`\`\`mermaid
graph TD
    A[Start] --> B{Entscheidung}
    B -->|Ja| C[Aktion A]
    B -->|Nein| D[Aktion B]
    C --> E[Ende]
    D --> E
\`\`\`

## Tabelle

| Spalte 1 | Spalte 2 | Spalte 3 |
|----------|----------|----------|
| A        | B        | C        |
| 1        | 2        | 3        |

## Task-Liste

- [x] Feature implementiert
- [x] Tests geschrieben
- [ ] Dokumentation aktualisieren

## Syntax-Highlighting

\`\`\`typescript
interface User {
  id: number
  name: string
}

function greet(user: User): string {
  return \`Hello, \${user.name}!\`
}
\`\`\`

## Wikilinks

Verweis auf [[Andere Seite]] oder [[Dokumentation]].

## Abkürzungen

HTML und CSS sind Webstandards.

*[HTML]: HyperText Markup Language
*[CSS]: Cascading Style Sheets

## Fußnote

Das ist ein Text mit Fußnote[^1].

[^1]: Das ist die Fußnote.
`
}
</script>

<style scoped>
.preview-content :deep(pre) {
  background: #f6f8fa;
  border: 1px solid #d1d9e0;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
  margin-bottom: 16px;
}
.preview-content :deep(code) {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 0.85em;
}
.preview-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 16px;
}
.preview-content :deep(th),
.preview-content :deep(td) {
  border: 1px solid #d1d9e0;
  padding: 6px 13px;
}
.preview-content :deep(th) {
  background: #f6f8fa;
  font-weight: 600;
}
.preview-content :deep(.mermaid) {
  text-align: center;
  margin-bottom: 16px;
}
.preview-content :deep(abbr) {
  cursor: help;
  border-bottom: 1px dotted #57606a;
  text-decoration: none;
}
.preview-content :deep(.wikilink) {
  color: #6f42c1;
  font-style: italic;
}
.preview-content :deep(a.header-anchor) {
  display: none;
}
.preview-content :deep(.task-list-item) {
  list-style: none;
}
.preview-content :deep(blockquote) {
  border-left: 4px solid #d1d9e0;
  padding: 0 1em;
  color: #57606a;
  margin: 0 0 16px;
}
</style>
