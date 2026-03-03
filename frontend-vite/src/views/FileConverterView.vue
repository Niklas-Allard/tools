<template>
  <div class="space-y-8">
    <div class="text-center space-y-2">
      <h1 class="text-5xl font-bold tracking-tight">
        <span class="bg-gradient-to-r from-violet-600 to-blue-600 bg-clip-text text-transparent">
          🎞️ File Converter
        </span>
      </h1>
      <p class="text-muted-foreground text-lg">
        Video &amp; Audio direkt im Browser konvertieren — kein Server, 100% lokal.
      </p>
    </div>
    <FileConverter />
  </div>
</template>

<script setup lang="ts">
import FileConverter from '@/components/FileConverter.vue'
</script>