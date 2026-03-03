<template>
  <div class="min-h-screen bg-background">
    <Toaster position="top-right" rich-colors expand :duration="4000" />
    <TooltipProvider>

      <!-- Navigation Bar -->
      <nav class="border-b bg-card/50 backdrop-blur-sm sticky top-0 z-50">
        <div class="container mx-auto px-4 max-w-7xl">
          <div class="flex items-center h-14 gap-4">
            <!-- Logo / Home -->
            <RouterLink
              to="/"
              class="flex items-center gap-2 font-bold text-lg hover:opacity-80 transition-opacity"
            >
              <span class="text-xl">🛠️</span>
              <span class="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                Dev Tools
              </span>
            </RouterLink>

            <!-- Breadcrumb -->
            <template v-if="currentTool">
              <span class="text-muted-foreground">/</span>
              <span class="text-sm font-medium">{{ currentTool.name }}</span>
            </template>

            <div class="flex-1" />

            <!-- Back to Home -->
            <RouterLink
              v-if="currentTool"
              to="/"
              class="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              ← Alle Tools
            </RouterLink>
          </div>
        </div>
      </nav>

      <!-- Main Content -->
      <main class="container mx-auto px-4 py-8 max-w-7xl">
        <RouterView />
      </main>

    </TooltipProvider>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Toaster } from 'vue-sonner'
import { TooltipProvider } from '@/components/ui/tooltip'

const route = useRoute()

const tools: Record<string, { name: string }> = {
  ytdl: { name: 'YouTube Downloader' },
  qrcode: { name: 'QR Code Generator' },
  'file-converter': { name: 'File Converter' },
}

const currentTool = computed(() => {
  const name = route.name as string
  return tools[name] ?? null
})
</script>
