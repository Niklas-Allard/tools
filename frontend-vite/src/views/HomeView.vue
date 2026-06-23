<template>
  <div class="space-y-10">
    <!-- Hero Header -->
    <div class="text-center space-y-3">
      <h1 class="text-5xl font-bold tracking-tight">
        <span class="bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 bg-clip-text text-transparent">
          Dev Tools
        </span>
      </h1>
      <p class="text-muted-foreground text-lg max-w-xl mx-auto">
        Eine Sammlung nützlicher Tools für den Alltag.
      </p>
    </div>

    <!-- Tools Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <button
        v-for="tool in tools"
        :key="tool.id"
        @click="router.push({ name: tool.id })"
        class="group text-left p-5 rounded-xl border bg-card hover:border-primary/50 hover:shadow-lg transition-all duration-200 hover:-translate-y-0.5"
      >
        <div class="flex items-start gap-4">
          <div
            class="text-3xl p-2.5 rounded-lg shrink-0"
            :class="tool.iconBg"
          >
            {{ tool.icon }}
          </div>
          <div class="flex-1 min-w-0 space-y-1">
            <div class="flex items-center gap-2">
              <h3 class="font-semibold text-sm group-hover:text-primary transition-colors">
                {{ tool.name }}
              </h3>
              <span
                v-if="tool.badge"
                class="text-xs px-1.5 py-0.5 rounded-full font-medium"
                :class="tool.badgeClass"
              >
                {{ tool.badge }}
              </span>
            </div>
            <p class="text-xs text-muted-foreground leading-relaxed">
              {{ tool.description }}
            </p>
          </div>
        </div>
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const tools = [
  {
    id: 'ytdl',
    icon: '🎬',
    iconBg: 'bg-red-100 dark:bg-red-950',
    name: 'YouTube Downloader',
    description: 'Videos & Playlists herunterladen mit 25+ yt-dlp Optionen, SponsorBlock und mehr.',
    badge: 'Pro',
    badgeClass: 'bg-red-100 text-red-700 dark:bg-red-950 dark:text-red-400',
  },
  {
    id: 'qrcode',
    icon: '📱',
    iconBg: 'bg-blue-100 dark:bg-blue-950',
    name: 'QR Code Generator',
    description: 'QR Codes aus URLs, Text oder Kontakten generieren.',
    badge: 'Neu',
    badgeClass: 'bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-400',
  },
  {
    id: 'file-converter',
    icon: '🎞️',
    iconBg: 'bg-violet-100 dark:bg-violet-950',
    name: 'File Converter',
    description: 'Video & Audio im Browser konvertieren: MP4, WebM, GIF, MP3, WAV — 100% lokal per FFmpeg.wasm.',
    badge: 'Neu',
    badgeClass: 'bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-400',
  },
  {
    id: 'md-pdf',
    icon: '📄',
    iconBg: 'bg-emerald-100 dark:bg-emerald-950',
    name: 'Markdown → PDF',
    description: 'MD mit Mermaid-Diagrammen, KaTeX/LaTeX, Admonitions, TOC, Wikilinks und Syntax-Highlighting zu PDF exportieren.',
    badge: 'Neu',
    badgeClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400',
  },
  // Weitere Tools hier eintragen:
  // {
  //   id: 'qrcode',
  //   icon: '📱',
  //   iconBg: 'bg-blue-100 dark:bg-blue-950',
  //   name: 'QR Code Generator',
  //   description: 'QR Codes aus URLs, Text oder Kontakten generieren.',
  //   badge: 'Neu',
  //   badgeClass: 'bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-400',
  // },
]
</script>