<template>
  <div class="max-w-3xl mx-auto space-y-5">

    <!-- ── 1. Drop Zone ──────────────────────────────────────────────────── -->
    <div
      v-if="!file && batchFiles.length === 0"
      @dragover.prevent="isDragOver = true"
      @dragleave.prevent="isDragOver = false"
      @drop.prevent="handleDrop"
      :class="[
        'relative flex flex-col items-center justify-center gap-4 rounded-xl border-2 border-dashed p-12 text-center transition-colors cursor-pointer',
        isDragOver
          ? 'border-primary bg-primary/5'
          : 'border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/30',
      ]"
      @click="fileInputRef?.click()"
    >
      <div class="flex h-14 w-14 items-center justify-center rounded-full bg-muted">
        <Upload class="h-7 w-7 text-muted-foreground" />
      </div>
      <div class="space-y-1">
        <p class="text-base font-medium">Datei hierher ziehen</p>
        <p class="text-sm text-muted-foreground">oder klicken zum Auswählen</p>
        <p class="text-xs text-muted-foreground/60 pt-1">
          Video: MP4 · MKV · AVI · MOV · WebM · OGV &nbsp;|&nbsp;
          Audio: MP3 · WAV · OGG · FLAC &nbsp;|&nbsp;
          Bild: PNG · JPG · WEBP · BMP · GIF
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click.stop="toggleBatchMode"
          :class="[
            'flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border transition-colors',
            batchMode
              ? 'bg-primary text-primary-foreground border-primary'
              : 'border-border hover:border-primary/50 text-muted-foreground',
          ]"
        >
          <Layers class="h-3 w-3" />
          Batch-Modus
        </button>
      </div>
      <input
        ref="fileInputRef"
        type="file"
        :multiple="batchMode"
        accept="video/*,audio/*,image/*,.mkv,.flac,.opus,.ogg,.webm,.avif,.bmp,.tiff"
        class="hidden"
        @change="handleFileInput"
      />
    </div>

    <!-- ── 2a. Batch file list ────────────────────────────────────────────── -->
    <template v-if="batchFiles.length > 0">
      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-base flex items-center justify-between">
            <span class="flex items-center gap-2">
              <Layers class="h-4 w-4" />
              Batch-Warteschlange ({{ batchFiles.length }} Datei{{ batchFiles.length !== 1 ? 'en' : '' }})
            </span>
            <button
              class="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1"
              @click="addMoreFiles"
            >
              <Plus class="h-3 w-3" /> Hinzufügen
            </button>
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-1.5 pt-0">
          <div
            v-for="item in batchFiles"
            :key="item.id"
            class="flex items-center gap-2 rounded-md px-2 py-1.5 hover:bg-muted/40"
          >
            <!-- Category icon -->
            <div :class="['flex h-7 w-7 shrink-0 items-center justify-center rounded', categoryIconClass(item.category)]">
              <Film v-if="item.category === 'video'" class="h-3.5 w-3.5" />
              <Music2 v-else-if="item.category === 'audio'" class="h-3.5 w-3.5" />
              <ImageIcon v-else class="h-3.5 w-3.5" />
            </div>
            <!-- Name -->
            <p class="flex-1 text-sm truncate">{{ item.file.name }}</p>
            <span class="text-xs text-muted-foreground shrink-0">{{ fmtBytes(item.file.size) }}</span>
            <!-- Status -->
            <span :class="['text-xs px-1.5 py-0.5 rounded-full shrink-0', batchStatusClass(item.status)]">
              {{ batchStatusLabel(item.status) }}
            </span>
            <!-- Download (done) -->
            <button v-if="item.status === 'done' && item.outputBlob" @click="downloadBatchItem(item)" class="shrink-0 text-muted-foreground hover:text-foreground">
              <Download class="h-3.5 w-3.5" />
            </button>
            <!-- Remove (pending) -->
            <button v-if="item.status === 'pending'" @click="removeBatchItem(item.id)" class="shrink-0 text-muted-foreground hover:text-destructive">
              <X class="h-3.5 w-3.5" />
            </button>
          </div>
        </CardContent>
      </Card>
      <input ref="fileInputMoreRef" type="file" multiple accept="video/*,audio/*,image/*,.mkv,.flac,.opus,.ogg,.webm,.avif,.bmp,.tiff" class="hidden" @change="handleMoreFilesInput" />
    </template>

    <!-- ── 2b. Single file view ───────────────────────────────────────────── -->
    <template v-else-if="file">
      <!-- File info bar -->
      <Card>
        <CardContent class="flex items-center gap-3 p-4">
          <div :class="['flex h-10 w-10 shrink-0 items-center justify-center rounded-lg', categoryIconClass(fileCategory)]">
            <Film v-if="fileCategory === 'video'" class="h-5 w-5" />
            <Music2 v-else-if="fileCategory === 'audio'" class="h-5 w-5" />
            <ImageIcon v-else class="h-5 w-5" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{{ file.name }}</p>
            <p class="text-xs text-muted-foreground">
              {{ fmtBytes(file.size) }}
              <span v-if="fileDuration > 0"> · {{ formatTime(fileDuration) }}</span>
              <span v-if="file.size > 500 * 1024 * 1024" class="text-amber-500 ml-1">⚠ Große Datei</span>
            </p>
          </div>
          <Button variant="ghost" size="icon" class="shrink-0 h-8 w-8" @click="clearFile">
            <X class="h-4 w-4" />
          </Button>
        </CardContent>
      </Card>

      <!-- Media / Image preview -->
      <Card>
        <CardContent class="p-4">
          <video v-if="fileCategory === 'video'" :src="fileUrl" controls class="w-full rounded-lg max-h-64 bg-black" @loadedmetadata="handleMediaLoaded" />
          <audio v-else-if="fileCategory === 'audio'" :src="fileUrl" controls class="w-full" @loadedmetadata="handleMediaLoaded" />
          <div v-else class="flex items-center justify-center rounded-lg bg-muted/30 p-4 min-h-32">
            <img :src="fileUrl" alt="Vorschau" class="max-w-full max-h-64 rounded-lg object-contain" />
          </div>
        </CardContent>
      </Card>

      <!-- SVG warning -->
      <div v-if="isSvgInput" class="flex items-start gap-2 rounded-lg bg-amber-500/10 border border-amber-500/20 p-3 text-sm text-amber-700 dark:text-amber-400">
        <AlertCircle class="h-4 w-4 mt-0.5 shrink-0" />
        <span>SVG-Dateien können von FFmpeg.wasm nicht direkt verarbeitet werden – nur Vektorgrafik-Informationen werden gespeichert. Konvertierung wird wahrscheinlich fehlschlagen.</span>
      </div>
    </template>

    <!-- ── 3. Options + Convert ───────────────────────────────────────────── -->
    <Card v-if="(file || batchFiles.length > 0) && (appState === 'idle' || appState === 'error')">
      <CardHeader class="pb-3">
        <CardTitle class="text-base flex items-center gap-2">
          <Settings2 class="h-4 w-4" />
          Konvertierungsoptionen
        </CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">

        <!-- ── Output format group tabs ── -->
        <div class="space-y-3">
          <label class="text-sm font-medium">Ausgabe-Kategorie</label>
          <div class="flex gap-1 rounded-lg bg-muted p-1 w-fit">
            <button
              v-for="cat in formatCategories"
              :key="cat.id"
              @click="setFormatCategory(cat.id)"
              :class="[
                'flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                outputFormatCategory === cat.id
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground',
              ]"
            >
              <component :is="cat.icon" class="h-3.5 w-3.5" />
              {{ cat.label }}
            </button>
          </div>

          <!-- Format pills within selected category -->
          <div class="flex flex-wrap gap-2">
            <button
              v-for="fmt in formatsInCategory"
              :key="fmt.ext"
              @click="selectFormat(fmt.ext)"
              :class="[
                'group relative px-3 py-1.5 rounded-md text-sm font-medium border transition-colors',
                outputFormat === fmt.ext
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'border-border hover:border-primary/50 hover:bg-muted',
              ]"
            >
              {{ fmt.label }}
              <span v-if="fmt.note" class="ml-1 text-[10px] opacity-60">{{ fmt.note }}</span>
            </button>
          </div>
        </div>

        <!-- ── Codec selector ── -->
        <div v-if="currentFormatDef?.altCodecs?.length" class="space-y-2">
          <label class="text-sm font-medium flex items-center gap-2">
            <Cpu class="h-4 w-4" />
            Codec
          </label>
          <select
            v-model="selectedCodec"
            class="w-full h-9 px-3 text-sm border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          >
            <option value="">{{ currentFormatDef.label }} Standard ({{ currentFormatDef.codec }})</option>
            <option v-for="ac in currentFormatDef.altCodecs" :key="ac.id" :value="ac.id">
              {{ ac.label }} ({{ ac.id }})
            </option>
          </select>
        </div>

        <Separator />

        <!-- ── Video output options ── -->
        <template v-if="outputFormatCategory === 'video'">
          <div class="grid grid-cols-2 gap-4">
            <!-- Resize -->
            <div class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <Monitor class="h-4 w-4" />
                Auflösung
              </label>
              <select v-model="resizePreset" class="w-full h-9 px-3 text-sm border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring">
                <option value="">Original</option>
                <option value="2160">4K (2160p)</option>
                <option value="1440">1440p (2K)</option>
                <option value="1080">1080p (Full HD)</option>
                <option value="720">720p (HD)</option>
                <option value="480">480p</option>
                <option value="360">360p</option>
              </select>
            </div>
            <!-- FPS -->
            <div class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <Timer class="h-4 w-4" />
                Framerate (FPS)
              </label>
              <select v-model="fpsPreset" class="w-full h-9 px-3 text-sm border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring">
                <option value="">Original</option>
                <option value="60">60 fps</option>
                <option value="30">30 fps</option>
                <option value="25">25 fps (PAL)</option>
                <option value="24">24 fps (Cinema)</option>
                <option value="15">15 fps</option>
              </select>
            </div>
          </div>

          <!-- CRF (quality) — hidden for GIF output -->
          <div v-if="outputFormat !== 'gif'" class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Gauge class="h-4 w-4" />
              Qualität — CRF {{ crf }}
              <span class="text-xs font-normal text-muted-foreground">(niedrig = besser)</span>
            </label>
            <input type="range" min="0" max="51" step="1" v-model.number="crf" class="w-full accent-primary" />
            <div class="flex justify-between text-xs text-muted-foreground">
              <span>0 = verlustfrei</span><span>↑ beste</span><span>28 = Standard</span><span>↑ kleinste</span><span>51 = schlecht</span>
            </div>
          </div>

          <!-- Audio strip -->
          <div class="flex items-center justify-between rounded-lg border px-4 py-3">
            <div>
              <p class="text-sm font-medium">Audio entfernen</p>
              <p class="text-xs text-muted-foreground">Nur Video, kein Ton (-an Flag)</p>
            </div>
            <button
              role="switch"
              :aria-checked="stripAudio"
              @click="stripAudio = !stripAudio"
              :class="['relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus-visible:outline-none focus-visible:ring-2', stripAudio ? 'bg-primary' : 'bg-input']"
            >
              <span :class="['inline-block h-3.5 w-3.5 rounded-full bg-white shadow-sm transition-transform', stripAudio ? 'translate-x-4' : 'translate-x-1']" />
            </button>
          </div>

          <!-- Audio bitrate (when audio not stripped) -->
          <template v-if="!stripAudio && fileCategory !== 'image'">
            <div class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <Volume2 class="h-4 w-4" />
                Audio-Bitrate
              </label>
              <div class="flex gap-2 flex-wrap">
                <button
                  v-for="br in AUDIO_BITRATES"
                  :key="br"
                  @click="audioBitrate = br"
                  :class="['px-3 py-1.5 rounded-md text-sm border transition-colors', audioBitrate === br ? 'bg-primary text-primary-foreground border-primary' : 'border-border hover:border-primary/50 hover:bg-muted']"
                >{{ br }}</button>
              </div>
            </div>
          </template>
        </template>

        <!-- ── Audio output options ── -->
        <template v-if="outputFormatCategory === 'audio'">
          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Volume2 class="h-4 w-4" />
              Audio-Bitrate
              <span class="text-xs font-normal text-muted-foreground">(nicht für WAV/FLAC)</span>
            </label>
            <div class="flex gap-2 flex-wrap">
              <button
                v-for="br in AUDIO_BITRATES"
                :key="br"
                @click="audioBitrate = br"
                :class="['px-3 py-1.5 rounded-md text-sm border transition-colors', audioBitrate === br ? 'bg-primary text-primary-foreground border-primary' : 'border-border hover:border-primary/50 hover:bg-muted']"
              >{{ br }}</button>
            </div>
          </div>
        </template>

        <!-- ── Image output options ── -->
        <template v-if="outputFormatCategory === 'image'">
          <!-- Quality (hidden for lossless formats) -->
          <div v-if="!LOSSLESS_IMAGE_FORMATS.includes(outputFormat)" class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Gauge class="h-4 w-4" />
              Bildqualität — {{ imageQuality }}%
            </label>
            <input type="range" min="1" max="100" step="1" v-model.number="imageQuality" class="w-full accent-primary" />
            <div class="flex justify-between text-xs text-muted-foreground">
              <span>1 = kleinste Datei</span><span>85 = empfohlen</span><span>100 = beste</span>
            </div>
          </div>
          <div v-else class="text-xs text-muted-foreground bg-muted/40 rounded-md px-3 py-2">
            {{ outputFormat.toUpperCase() }} ist verlustfrei — Qualitätsstufe nicht anwendbar.
          </div>
          <!-- Resize for image output -->
          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Monitor class="h-4 w-4" />
              Skalierung (Höhe)
            </label>
            <select v-model="resizePreset" class="w-full h-9 px-3 text-sm border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring">
              <option value="">Original</option>
              <option value="2160">4K (2160px)</option>
              <option value="1080">1080px</option>
              <option value="720">720px</option>
              <option value="480">480px</option>
              <option value="240">240px</option>
            </select>
          </div>
        </template>

        <!-- ── Trim (video/audio input only, not batch) ── -->
        <template v-if="fileDuration > 0 && !batchMode">
          <Separator />
          <div class="space-y-3">
            <label class="text-sm font-medium flex items-center gap-2">
              <Scissors class="h-4 w-4" />
              Trimmen
              <span class="text-xs font-normal text-muted-foreground ml-1">
                {{ formatTime(trimStart) }} → {{ formatTime(trimEnd) }}
                ({{ formatTime(Math.max(0, trimEnd - trimStart)) }})
              </span>
            </label>
            <div class="space-y-2 px-1">
              <div class="flex items-center gap-3">
                <span class="text-xs text-muted-foreground w-14 shrink-0">Start</span>
                <input type="range" :min="0" :max="fileDuration" :step="0.5" v-model.number="trimStart" @input="clampTrimStart" class="flex-1 accent-primary" />
                <span class="text-xs font-mono w-14 text-right shrink-0">{{ formatTime(trimStart) }}</span>
              </div>
              <div class="flex items-center gap-3">
                <span class="text-xs text-muted-foreground w-14 shrink-0">Ende</span>
                <input type="range" :min="0" :max="fileDuration" :step="0.5" v-model.number="trimEnd" @input="clampTrimEnd" class="flex-1 accent-primary" />
                <span class="text-xs font-mono w-14 text-right shrink-0">{{ formatTime(trimEnd) }}</span>
              </div>
            </div>
          </div>
        </template>

      </CardContent>
      <CardFooter class="flex flex-col gap-3 pt-0">
        <div v-if="appState === 'error'" class="w-full flex items-start gap-2 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-sm text-destructive">
          <AlertCircle class="h-4 w-4 mt-0.5 shrink-0" />
          <span class="break-all">{{ errorMsg }}</span>
        </div>
        <Button class="w-full" size="lg" @click="batchFiles.length > 0 ? convertBatch() : convert()">
          <Zap class="h-4 w-4 mr-2" />
          {{ isLoaded ? (batchFiles.length > 0 ? `Alle ${batchFiles.length} konvertieren` : 'Jetzt konvertieren') : 'FFmpeg laden &amp; konvertieren' }}
        </Button>
        <Button v-if="batchFiles.length > 0" variant="outline" size="sm" class="w-full" @click="clearBatch">
          <X class="h-3 w-3 mr-2" />
          Batch leeren
        </Button>
      </CardFooter>
    </Card>

    <!-- ── 4. Progress ────────────────────────────────────────────────────── -->
    <Card v-if="appState === 'loading-ffmpeg' || appState === 'converting'">
      <CardContent class="p-6 space-y-4">
        <div class="flex items-center gap-3">
          <Loader2 class="h-5 w-5 animate-spin text-primary shrink-0" />
          <div class="flex-1">
            <p class="text-sm font-medium">
              <template v-if="appState === 'loading-ffmpeg'">FFmpeg wird geladen…</template>
              <template v-else-if="gifPass > 0">
                GIF – Pass {{ gifPass }}/2
                <span class="text-muted-foreground font-normal ml-1">({{ gifPass === 1 ? 'Palette' : 'Rendern' }})</span>
              </template>
              <template v-else-if="batchProgress.total > 1">
                Batch: {{ batchProgress.current }}/{{ batchProgress.total }} — {{ batchProgress.name }}
              </template>
              <template v-else>Konvertierung läuft…</template>
            </p>
            <p v-if="appState === 'converting'" class="text-xs text-muted-foreground mt-0.5">
              {{ outputFormat.toUpperCase() }}
              <span v-if="selectedCodec"> · {{ selectedCodec }}</span>
              <span v-else-if="currentFormatDef"> · {{ currentFormatDef.codec }}</span>
              <span v-if="resizePreset && outputFormatCategory === 'video'"> · {{ resizePreset }}p</span>
              <span v-if="outputFormatCategory === 'audio'"> · {{ audioBitrate }}</span>
              <span v-if="outputFormatCategory === 'image' && !LOSSLESS_IMAGE_FORMATS.includes(outputFormat)"> · {{ imageQuality }}%</span>
            </p>
          </div>
          <span class="text-sm font-mono text-muted-foreground">{{ displayProgress }}%</span>
        </div>
        <Progress :model-value="displayProgress" class="h-2" />
        <div v-if="recentLogs.length > 0" class="rounded-md bg-muted/50 p-3 font-mono text-xs text-muted-foreground space-y-0.5 max-h-32 overflow-y-auto">
          <div v-for="(line, i) in recentLogs" :key="i" class="leading-relaxed truncate">{{ line }}</div>
        </div>
      </CardContent>
    </Card>

    <!-- ── 5. Output (done) ───────────────────────────────────────────────── -->
    <Card v-if="appState === 'done' && !batchMode" class="border-green-500/30 bg-green-500/5">
      <CardHeader class="pb-3">
        <CardTitle class="text-base flex items-center gap-2 text-green-600 dark:text-green-400">
          <CheckCircle2 class="h-5 w-5" />
          Konvertierung abgeschlossen
        </CardTitle>
        <CardDescription>{{ outputFilename }} · {{ fmtBytes(outputSize) }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <!-- Output preview -->
        <template v-if="outputUrl">
          <video v-if="outputMime.startsWith('video/')" :src="outputUrl" controls class="w-full rounded-lg max-h-64 bg-black" />
          <audio v-else-if="outputMime.startsWith('audio/')" :src="outputUrl" controls class="w-full" />
          <div v-else-if="outputMime.startsWith('image/')" class="flex items-center justify-center rounded-lg bg-muted/30 p-4">
            <img :src="outputUrl" alt="Ergebnis" class="max-w-full max-h-64 rounded-lg object-contain" />
          </div>
        </template>
        <div class="flex gap-3">
          <Button class="flex-1" @click="downloadOutput">
            <Download class="h-4 w-4 mr-2" />
            Herunterladen
          </Button>
          <Button variant="outline" class="flex-1" @click="resetToConvert">
            <RefreshCw class="h-4 w-4 mr-2" />
            Neu konvertieren
          </Button>
        </div>
      </CardContent>
    </Card>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { fetchFile } from '@ffmpeg/util'
import { toast } from 'vue-sonner'
import { useFFmpeg } from '@/composables/useFFmpeg'
import {
  Upload, Film, Music2, X, Download, RefreshCw, Settings2, Scissors,
  Monitor, Gauge, Volume2, Zap, Loader2, CheckCircle2, AlertCircle,
  Layers, Plus, Cpu, Timer, Image as ImageIcon,
} from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'

// ── Types ────────────────────────────────────────────────────────────────────
type AppState = 'idle' | 'loading-ffmpeg' | 'converting' | 'done' | 'error'
type FormatCategory = 'video' | 'audio' | 'image'
type FileCategory = 'video' | 'audio' | 'image' | 'unknown'
type AudioBitrate = '64k' | '128k' | '192k' | '256k' | '320k'

interface FormatDef {
  ext: string
  label: string
  category: FormatCategory
  mime: string
  /** Default codec */
  codec: string
  /** Extra codecs the user can choose */
  altCodecs?: { id: string; label: string }[]
  /** Short warning shown in pill */
  note?: string
}

interface BatchItem {
  id: string
  file: File
  url: string
  category: FileCategory
  status: 'pending' | 'converting' | 'done' | 'error'
  outputBlob?: Blob
  outputFilename?: string
  error?: string
}

// ── Format registry ─────────────────────────────────────────────────────────
// Codecs verified available in @ffmpeg/core-mt@0.12.x:
//   libx264 ✓  libvpx ✓  libvpx-vp9 ✓  libtheora ✓  mpeg4 ✓
//   libmp3lame ✓  aac ✓  libopus ✓  libvorbis ✓  pcm_s16le ✓  flac ✓
//   libwebp ✓  mjpeg ✓  png ✓  bmp ✓  gif ✓
//   libaom-av1 ✗ (not in standard build)  libx265 ✗ (not compiled in)
const FORMATS: FormatDef[] = [
  // ── Video
  {
    ext: 'mp4', label: 'MP4', category: 'video', mime: 'video/mp4',
    codec: 'libx264', altCodecs: [{ id: 'mpeg4', label: 'MPEG-4' }],
  },
  {
    ext: 'webm', label: 'WebM', category: 'video', mime: 'video/webm',
    codec: 'libvpx-vp9', altCodecs: [{ id: 'libvpx', label: 'VP8' }],
  },
  {
    ext: 'ogv', label: 'OGV', category: 'video', mime: 'video/ogg',
    codec: 'libtheora',
  },
  {
    ext: 'avi', label: 'AVI', category: 'video', mime: 'video/x-msvideo',
    codec: 'mpeg4', altCodecs: [{ id: 'libx264', label: 'H.264' }],
  },
  {
    ext: 'mkv', label: 'MKV', category: 'video', mime: 'video/x-matroska',
    codec: 'libx264', altCodecs: [{ id: 'libvpx-vp9', label: 'VP9' }, { id: 'libtheora', label: 'Theora' }],
  },
  // ── Audio
  { ext: 'mp3',  label: 'MP3',  category: 'audio', mime: 'audio/mpeg', codec: 'libmp3lame' },
  { ext: 'wav',  label: 'WAV',  category: 'audio', mime: 'audio/wav',  codec: 'pcm_s16le' },
  { ext: 'ogg',  label: 'OGG',  category: 'audio', mime: 'audio/ogg',  codec: 'libvorbis' },
  { ext: 'opus', label: 'OPUS', category: 'audio', mime: 'audio/ogg',  codec: 'libopus' },
  { ext: 'flac', label: 'FLAC', category: 'audio', mime: 'audio/flac', codec: 'flac' },
  { ext: 'aac',  label: 'AAC',  category: 'audio', mime: 'audio/aac',  codec: 'aac' },
  // ── Image
  { ext: 'jpg',  label: 'JPEG', category: 'image', mime: 'image/jpeg', codec: 'mjpeg' },
  { ext: 'png',  label: 'PNG',  category: 'image', mime: 'image/png',  codec: 'png' },
  { ext: 'webp', label: 'WebP', category: 'image', mime: 'image/webp', codec: 'libwebp' },
  { ext: 'bmp',  label: 'BMP',  category: 'image', mime: 'image/bmp',  codec: 'bmp' },
  { ext: 'gif',  label: 'GIF',  category: 'image', mime: 'image/gif',  codec: 'gif' },
  { ext: 'avif', label: 'AVIF', category: 'image', mime: 'image/avif', codec: 'libaom-av1', note: '⚠ exp.' },
]

const LOSSLESS_IMAGE_FORMATS = ['png', 'bmp', 'gif']  // quality slider hidden
const AUDIO_BITRATES: AudioBitrate[] = ['64k', '128k', '192k', '256k', '320k']

// ── FFmpeg ──────────────────────────────────────────────────────────────────
const {
  isLoaded, progress: ffmpegProgress, logs, load,
  writeFile, exec, readFile, deleteFile, terminate, convertImage,
} = useFFmpeg()

// ── App State ───────────────────────────────────────────────────────────────
const appState = ref<AppState>('idle')
const errorMsg = ref('')
const gifPass = ref(0)
const batchProgress = ref({ current: 0, total: 0, name: '' })

// ── Single file ─────────────────────────────────────────────────────────────
const fileInputRef = ref<HTMLInputElement | null>(null)
const file = ref<File | null>(null)
const fileUrl = ref('')
const fileCategory = ref<FileCategory>('unknown')
const fileDuration = ref(0)
const isDragOver = ref(false)

// ── Batch mode ───────────────────────────────────────────────────────────────
const batchMode = ref(false)
const batchFiles = ref<BatchItem[]>([])
const fileInputMoreRef = ref<HTMLInputElement | null>(null)

// ── Options ─────────────────────────────────────────────────────────────────
const outputFormatCategory = ref<FormatCategory>('video')
const outputFormat = ref('mp4')
const selectedCodec = ref('')
const crf = ref(23)
const imageQuality = ref(85)
const resizePreset = ref('')
const fpsPreset = ref('')
const stripAudio = ref(false)
const audioBitrate = ref<AudioBitrate>('192k')
const trimStart = ref(0)
const trimEnd = ref(0)

// ── Output ───────────────────────────────────────────────────────────────────
const outputUrl = ref('')
const outputFilename = ref('')
const outputMime = ref('')
const outputSize = ref(0)

// ── Format computed ──────────────────────────────────────────────────────────
const formatsInCategory = computed(() => FORMATS.filter(f => f.category === outputFormatCategory.value))
const currentFormatDef = computed(() => FORMATS.find(f => f.ext === outputFormat.value))

const formatCategories = [
  { id: 'video' as FormatCategory, label: 'Video', icon: Film },
  { id: 'audio' as FormatCategory, label: 'Audio', icon: Music2 },
  { id: 'image' as FormatCategory, label: 'Bild', icon: ImageIcon },
]

const isSvgInput = computed(() =>
  file.value ? file.value.name.toLowerCase().endsWith('.svg') || file.value.type === 'image/svg+xml' : false,
)

// ── Progress display ─────────────────────────────────────────────────────────
const displayProgress = computed(() => {
  if (appState.value === 'loading-ffmpeg') return Math.min(ffmpegProgress.value, 99)
  if (gifPass.value === 1) return Math.round(ffmpegProgress.value * 0.5)
  if (gifPass.value === 2) return Math.round(50 + ffmpegProgress.value * 0.5)
  return ffmpegProgress.value
})

const recentLogs = computed(() => logs.value.slice(-10))

// ── Watchers ─────────────────────────────────────────────────────────────────
watch(outputFormat, () => { selectedCodec.value = '' })
watch(fileDuration, (d) => { if (d > 0) trimEnd.value = Math.floor(d) })

// ── Helpers ───────────────────────────────────────────────────────────────────
function fmtBytes(b: number): string {
  const mb = b / 1024 / 1024
  return mb < 1 ? `${(mb * 1024).toFixed(0)} KB` : `${mb.toFixed(1)} MB`
}

function formatTime(s: number): string {
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const sec = Math.floor(s % 60)
  if (h > 0) return `${h}:${m.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}`
  return `${m}:${sec.toString().padStart(2, '0')}`
}

function detectCategory(f: File): FileCategory {
  if (f.type.startsWith('video/') || /\.(mp4|mkv|avi|mov|webm|ogv|flv|wmv|m4v|3gp)$/i.test(f.name)) return 'video'
  if (f.type.startsWith('audio/') || /\.(mp3|wav|aac|flac|ogg|opus|m4a|wma)$/i.test(f.name)) return 'audio'
  if (f.type.startsWith('image/') || /\.(png|jpg|jpeg|webp|gif|bmp|avif|tiff?|svg)$/i.test(f.name)) return 'image'
  return 'unknown'
}

function categoryIconClass(cat: FileCategory): string {
  if (cat === 'video') return 'bg-blue-100 dark:bg-blue-950 text-blue-600 dark:text-blue-400'
  if (cat === 'audio') return 'bg-purple-100 dark:bg-purple-950 text-purple-600 dark:text-purple-400'
  if (cat === 'image') return 'bg-green-100 dark:bg-green-950 text-green-600 dark:text-green-400'
  return 'bg-muted text-muted-foreground'
}

function setFormatCategory(cat: FormatCategory) {
  outputFormatCategory.value = cat
  const first = FORMATS.find(f => f.category === cat)
  if (first) outputFormat.value = first.ext
  selectedCodec.value = ''
  // Reset strip-audio for non-video output
  if (cat !== 'video') stripAudio.value = false
}

function selectFormat(ext: string) {
  outputFormat.value = ext
  selectedCodec.value = ''
}

function setDefaultFormatForInput(cat: FileCategory) {
  if (cat === 'audio') {
    outputFormatCategory.value = 'audio'
    outputFormat.value = 'mp3'
  } else if (cat === 'image') {
    outputFormatCategory.value = 'image'
    outputFormat.value = 'webp'
  } else {
    outputFormatCategory.value = 'video'
    outputFormat.value = 'mp4'
  }
  selectedCodec.value = ''
}

// ── Trim helpers ──────────────────────────────────────────────────────────────
function clampTrimStart() {
  if (trimStart.value >= trimEnd.value) trimStart.value = Math.max(0, trimEnd.value - 0.5)
}
function clampTrimEnd() {
  if (trimEnd.value <= trimStart.value) trimEnd.value = Math.min(fileDuration.value, trimStart.value + 0.5)
}
function trimArgs(): { pre: string[]; post: string[] } {
  const hasTrim = trimStart.value > 0.5 || trimEnd.value < fileDuration.value - 0.5
  if (!hasTrim || fileDuration.value === 0) return { pre: [], post: [] }
  const dur = Math.max(0.5, trimEnd.value - trimStart.value)
  return { pre: ['-ss', trimStart.value.toFixed(3)], post: ['-t', dur.toFixed(3)] }
}

// ── FFmpeg arg builders ───────────────────────────────────────────────────────
function buildVideoArgs(inputName: string, outputName: string): string[] {
  const { pre, post } = trimArgs()
  const codec = selectedCodec.value || currentFormatDef.value?.codec || 'libx264'
  const args: string[] = [...pre, '-i', inputName]

  // Video codec + quality
  switch (codec) {
    case 'libx264':
      args.push('-c:v', 'libx264', '-crf', String(crf.value), '-preset', 'medium', '-movflags', '+faststart')
      break
    case 'libvpx-vp9':
      args.push('-c:v', 'libvpx-vp9', '-crf', String(crf.value), '-b:v', '0')
      break
    case 'libvpx':
      args.push('-c:v', 'libvpx', '-crf', String(crf.value), '-b:v', '1M')
      break
    case 'libtheora': {
      // Theora quality: 0-10 (10=best). Map CRF 0-51 → q 10-1
      const q = Math.max(1, Math.round(10 - (crf.value / 51) * 9))
      args.push('-c:v', 'libtheora', '-q:v', String(q))
      break
    }
    case 'mpeg4': {
      const q = Math.max(1, Math.round(1 + (crf.value / 51) * 30))
      args.push('-c:v', 'mpeg4', '-q:v', String(q))
      break
    }
    default:
      args.push('-c:v', codec, '-crf', String(crf.value))
  }

  // Video filter chain: scale + fps
  const vfParts: string[] = []
  if (resizePreset.value) vfParts.push(`scale=-2:${resizePreset.value}`)
  if (fpsPreset.value)    vfParts.push(`fps=${fpsPreset.value}`)
  if (vfParts.length)     args.push('-vf', vfParts.join(','))

  // Audio
  if (stripAudio.value) {
    args.push('-an')
  } else {
    const audioCodec =
      outputFormat.value === 'webm' || outputFormat.value === 'mkv'
        ? 'libopus'
        : outputFormat.value === 'ogv'
          ? 'libvorbis'
          : 'aac'
    args.push('-c:a', audioCodec, '-b:a', audioBitrate.value)
  }

  return [...args, ...post, '-y', outputName]
}

function buildAudioArgs(inputName: string, outputName: string): string[] {
  const { pre, post } = trimArgs()
  const codec = selectedCodec.value || currentFormatDef.value?.codec || 'libmp3lame'
  const args: string[] = [...pre, '-i', inputName, '-vn']

  switch (codec) {
    case 'libmp3lame':
      args.push('-c:a', 'libmp3lame', '-b:a', audioBitrate.value, '-q:a', '2')
      break
    case 'pcm_s16le':
      args.push('-c:a', 'pcm_s16le')
      break
    case 'libvorbis':
      args.push('-c:a', 'libvorbis', '-b:a', audioBitrate.value)
      break
    case 'libopus':
      args.push('-c:a', 'libopus', '-b:a', audioBitrate.value)
      if (outputFormat.value === 'opus') args.push('-f', 'opus')
      break
    case 'flac':
      args.push('-c:a', 'flac')
      break
    case 'aac':
      args.push('-c:a', 'aac', '-b:a', audioBitrate.value)
      break
    default:
      args.push('-c:a', codec)
  }

  return [...args, ...post, '-y', outputName]
}

// Two-pass animated GIF (video input)
async function convertToGif(inputName: string, outName: string): Promise<void> {
  const palette = 'palette.png'
  const scale = resizePreset.value ? `scale=-2:${resizePreset.value}` : 'scale=480:-1'
  const fps = fpsPreset.value || '10'
  const { pre, post } = trimArgs()

  gifPass.value = 1
  const code1 = await exec([
    ...pre, '-i', inputName, ...post,
    '-vf', `fps=${fps},${scale}:flags=lanczos,palettegen=stats_mode=diff`,
    '-y', palette,
  ])
  if (code1 !== 0) throw new Error(`GIF Palettegen fehlgeschlagen (Exit ${code1})`)

  gifPass.value = 2
  const code2 = await exec([
    ...pre, '-i', inputName, '-i', palette, ...post,
    '-filter_complex', `[0:v]fps=${fps},${scale}:flags=lanczos[x];[x][1:v]paletteuse=dither=bayer:bayer_scale=5:diff_mode=rectangle`,
    '-loop', '0', '-y', outName,
  ])
  if (code2 !== 0) throw new Error(`GIF Paletteuse fehlgeschlagen (Exit ${code2})`)

  await deleteFile(palette)
  gifPass.value = 0
}

// ── File handling ─────────────────────────────────────────────────────────────
function processFile(f: File): void {
  const cat = detectCategory(f)
  if (cat === 'unknown') {
    toast.error('Format unbekannt', { description: `${f.name} — FFmpeg wird es trotzdem versuchen.` })
  }
  if (f.size > 500 * 1024 * 1024) {
    toast.warning('Große Datei', { description: `${fmtBytes(f.size)} — Konvertierung kann sehr lange dauern.` })
  }

  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)

  file.value = f
  fileUrl.value = URL.createObjectURL(f)
  fileCategory.value = cat
  fileDuration.value = 0
  trimStart.value = 0
  trimEnd.value = 0
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
  setDefaultFormatForInput(cat)
}

function addToBatch(files: FileList): void {
  Array.from(files).forEach(f => {
    batchFiles.value.push({
      id: `${Date.now()}-${Math.random()}`,
      file: f,
      url: URL.createObjectURL(f),
      category: detectCategory(f),
      status: 'pending',
    })
  })
}

function handleDrop(e: DragEvent): void {
  isDragOver.value = false
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return
  if (batchMode.value) {
    addToBatch(files)
  } else if (files.length === 1) {
    processFile(files[0])
  } else {
    // Multiple files dropped in non-batch mode → auto-enable batch
    batchMode.value = true
    addToBatch(files)
  }
}

function handleFileInput(e: Event): void {
  const input = e.target as HTMLInputElement
  const files = input.files
  if (!files || files.length === 0) return
  if (batchMode.value || files.length > 1) {
    batchMode.value = true
    addToBatch(files)
  } else {
    processFile(files[0])
  }
  input.value = ''
}

function handleMoreFilesInput(e: Event): void {
  const input = e.target as HTMLInputElement
  if (input.files) addToBatch(input.files)
  input.value = ''
}

function clearFile(): void {
  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  file.value = null
  fileUrl.value = ''
  fileCategory.value = 'unknown'
  fileDuration.value = 0
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
}

function toggleBatchMode(): void {
  batchMode.value = !batchMode.value
  if (batchMode.value && file.value) {
    addToBatch(makeFileList(file.value))
    clearFile()
  }
}

function makeFileList(f: File): FileList {
  const dt = new DataTransfer()
  dt.items.add(f)
  return dt.files
}

function removeBatchItem(id: string): void {
  const idx = batchFiles.value.findIndex(b => b.id === id)
  if (idx >= 0) {
    URL.revokeObjectURL(batchFiles.value[idx].url)
    batchFiles.value.splice(idx, 1)
  }
}

function clearBatch(): void {
  batchFiles.value.forEach(b => URL.revokeObjectURL(b.url))
  batchFiles.value = []
  batchMode.value = false
  appState.value = 'idle'
}

function addMoreFiles(): void {
  fileInputMoreRef.value?.click()
}

function handleMediaLoaded(e: Event): void {
  const el = e.target as HTMLVideoElement
  if (el.duration && isFinite(el.duration)) {
    fileDuration.value = el.duration
    trimEnd.value = Math.floor(el.duration)
  }
}

function batchStatusLabel(s: BatchItem['status']): string {
  return { pending: 'Wartend', converting: 'Läuft…', done: 'Fertig', error: 'Fehler' }[s]
}
function batchStatusClass(s: BatchItem['status']): string {
  return {
    pending:    'bg-muted text-muted-foreground',
    converting: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',
    done:       'bg-green-500/10 text-green-600 dark:text-green-400',
    error:      'bg-red-500/10 text-red-600 dark:text-red-400',
  }[s]
}

// ── Core conversion logic ─────────────────────────────────────────────────────
async function ensureFFmpeg(): Promise<boolean> {
  if (isLoaded.value) return true
  appState.value = 'loading-ffmpeg'
  try {
    await load()
    return true
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    appState.value = 'error'
    errorMsg.value = `FFmpeg konnte nicht geladen werden: ${msg}`
    toast.error('Ladefehler', { description: errorMsg.value })
    return false
  }
}

async function convertFile(
  f: File,
  inputCat: FileCategory,
): Promise<{ blob: Blob; filename: string; mime: string }> {
  const outFmt = currentFormatDef.value!
  const inputExt = f.name.split('.').pop() ?? 'bin'
  const inputName = `input.${inputExt}`
  const outExt = outputFormat.value
  const outName = `output.${outExt}`
  const baseName = f.name.replace(/\.[^.]+$/, '')

  const inputData = await fetchFile(f)
  await writeFile(inputName, inputData)

  try {
    if (outFmt.category === 'image') {
      if (outExt === 'gif' && inputCat === 'video') {
        await convertToGif(inputName, outName)
      } else {
        // Still image → image (or video → thumbnail via -frames:v 1)
        await convertImage(inputName, outName, {
          quality: imageQuality.value,
          height: resizePreset.value ? Number(resizePreset.value) : 0,
        })
      }
    } else if (outFmt.category === 'audio') {
      const args = buildAudioArgs(inputName, outName)
      console.log('[FFmpeg]', args.join(' '))
      const code = await exec(args)
      if (code !== 0) throw new Error(`FFmpeg Exit-Code ${code}`)
    } else {
      // video output
      if (inputCat === 'audio') throw new Error('Audio-Datei kann nicht in ein Video-Format konvertiert werden.')
      const args = buildVideoArgs(inputName, outName)
      console.log('[FFmpeg]', args.join(' '))
      const code = await exec(args)
      if (code !== 0) throw new Error(`FFmpeg Exit-Code ${code}`)
    }

    const outData = await readFile(outName)
    const blob = new Blob([outData.slice(0)], { type: outFmt.mime })

    return { blob, filename: `${baseName}_converted.${outExt}`, mime: outFmt.mime }
  } finally {
    await deleteFile(inputName)
    await deleteFile(outName)
  }
}

async function convert(): Promise<void> {
  if (!file.value) return
  if (!await ensureFFmpeg()) return

  appState.value = 'converting'
  gifPass.value = 0

  try {
    const { blob, filename, mime } = await convertFile(file.value, fileCategory.value)
    if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
    outputUrl.value = URL.createObjectURL(blob)
    outputFilename.value = filename
    outputMime.value = mime
    outputSize.value = blob.size
    appState.value = 'done'
    toast.success('Fertig!', { description: `${filename} (${fmtBytes(blob.size)})` })
    downloadOutput()
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    errorMsg.value = msg
    appState.value = 'error'
    toast.error('Konvertierungsfehler', { description: msg })
    console.error('[FileConverter]', err)
    // A RuntimeError means the wasm instance crashed; reset it so the next
    // attempt can reload FFmpeg from scratch instead of reusing a broken instance.
    if (err instanceof Error && err.message.includes('RuntimeError')) {
      terminate()
    }
  }
}

async function convertBatch(): Promise<void> {
  const pending = batchFiles.value.filter(b => b.status === 'pending')
  if (pending.length === 0) return
  if (!await ensureFFmpeg()) return

  appState.value = 'converting'
  gifPass.value = 0
  batchProgress.value = { current: 0, total: pending.length, name: '' }

  for (const item of pending) {
    batchProgress.value.current++
    batchProgress.value.name = item.file.name
    item.status = 'converting'

    // Re-load FFmpeg if a prior wasm crash caused a terminate()
    if (!isLoaded.value) {
      const ok = await ensureFFmpeg()
      if (!ok) {
        item.status = 'error'
        item.error = 'FFmpeg konnte nicht neu geladen werden.'
        continue
      }
    }

    try {
      const { blob, filename } = await convertFile(item.file, item.category)
      item.outputBlob = blob
      item.outputFilename = filename
      item.status = 'done'
      // Auto-download each completed file
      downloadBatchItem(item)
      toast.success(filename, { description: fmtBytes(blob.size) })
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      item.error = msg
      item.status = 'error'
      toast.error(`Fehler: ${item.file.name}`, { description: msg })
      // Wasm crash – reset so the next batch item doesn't reuse a broken instance.
      if (err instanceof Error && err.message.includes('RuntimeError')) {
        terminate()
      }
    }
  }

  const doneCount = batchFiles.value.filter(b => b.status === 'done').length
  const errCount  = batchFiles.value.filter(b => b.status === 'error').length
  toast.info(`Batch abgeschlossen: ${doneCount} ✓, ${errCount} ✗`)
  appState.value = 'idle'
}

// ── Output ────────────────────────────────────────────────────────────────────
function downloadOutput(): void {
  if (!outputUrl.value) return
  const a = document.createElement('a')
  a.href = outputUrl.value
  a.download = outputFilename.value
  a.click()
}

function downloadBatchItem(item: BatchItem): void {
  if (!item.outputBlob || !item.outputFilename) return
  const url = URL.createObjectURL(item.outputBlob)
  const a = document.createElement('a')
  a.href = url
  a.download = item.outputFilename
  a.click()
  // Short delay before revoking so the download can start
  setTimeout(() => URL.revokeObjectURL(url), 2000)
}

function resetToConvert(): void {
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
}

// ── Cleanup ───────────────────────────────────────────────────────────────────
onUnmounted(() => {
  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  batchFiles.value.forEach(b => URL.revokeObjectURL(b.url))
  terminate()
})
</script>
