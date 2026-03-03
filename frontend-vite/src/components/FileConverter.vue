<template>
  <div class="max-w-2xl mx-auto space-y-5">

    <!-- ── 1. Drop Zone (no file loaded) ── -->
    <div
      v-if="!file"
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
        <p class="text-base font-medium">Video oder Audio hierher ziehen</p>
        <p class="text-sm text-muted-foreground">oder klicken zum Auswählen</p>
        <p class="text-xs text-muted-foreground/60 pt-1">MP4 · MKV · AVI · MOV · WebM · MP3 · WAV · AAC · FLAC …</p>
      </div>
      <input
        ref="fileInputRef"
        type="file"
        accept="video/*,audio/*"
        class="hidden"
        @change="handleFileInput"
      />
    </div>

    <!-- ── 2. File loaded ── -->
    <template v-else>

      <!-- File info bar -->
      <Card>
        <CardContent class="flex items-center gap-3 p-4">
          <div :class="['flex h-10 w-10 shrink-0 items-center justify-center rounded-lg', fileType === 'video' ? 'bg-blue-100 dark:bg-blue-950' : 'bg-purple-100 dark:bg-purple-950']">
            <Film v-if="fileType === 'video'" class="h-5 w-5 text-blue-600 dark:text-blue-400" />
            <Music2 v-else class="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{{ file.name }}</p>
            <p class="text-xs text-muted-foreground">
              {{ fileSizeFormatted }}
              <span v-if="fileDuration > 0"> · {{ durationFormatted }}</span>
              <span v-if="file.size > 500 * 1024 * 1024" class="text-amber-500 ml-1">⚠ Große Datei</span>
            </p>
          </div>
          <Button variant="ghost" size="icon" class="shrink-0 h-8 w-8" @click="clearFile">
            <X class="h-4 w-4" />
          </Button>
        </CardContent>
      </Card>

      <!-- Media preview -->
      <Card>
        <CardContent class="p-4">
          <video
            v-if="fileType === 'video'"
            :src="fileUrl"
            controls
            class="w-full rounded-lg max-h-64 bg-black"
            @loadedmetadata="handleMediaLoaded"
          />
          <audio
            v-else
            :src="fileUrl"
            controls
            class="w-full"
            @loadedmetadata="handleMediaLoaded"
          />
        </CardContent>
      </Card>

      <!-- ── 3. Options + Convert Button (idle / error) ── -->
      <Card v-if="appState === 'idle' || appState === 'error'">
        <CardHeader class="pb-3">
          <CardTitle class="text-base flex items-center gap-2">
            <Settings2 class="h-4 w-4" />
            Konvertierungsoptionen
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">

          <!-- Output format -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Ausgabe-Format</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="fmt in availableFormats"
                :key="fmt"
                @click="outputFormat = fmt"
                :class="[
                  'px-3 py-1.5 rounded-md text-sm font-medium border transition-colors',
                  outputFormat === fmt
                    ? 'bg-primary text-primary-foreground border-primary'
                    : 'border-border hover:border-primary/50 hover:bg-muted',
                ]"
              >{{ fmt.toUpperCase() }}</button>
            </div>
          </div>

          <Separator />

          <!-- Trim controls -->
          <div v-if="fileDuration > 0" class="space-y-3">
            <label class="text-sm font-medium flex items-center gap-2">
              <Scissors class="h-4 w-4" />
              Trimmen
              <span class="text-xs font-normal text-muted-foreground ml-1">
                {{ formatTime(trimStart) }} → {{ formatTime(trimEnd) }}
                ({{ formatTime(trimEnd - trimStart) }})
              </span>
            </label>
            <div class="space-y-2 px-1">
              <div class="flex items-center gap-3">
                <span class="text-xs text-muted-foreground w-14 shrink-0">Start</span>
                <input
                  type="range"
                  :min="0"
                  :max="fileDuration"
                  :step="0.5"
                  v-model.number="trimStart"
                  @input="clampTrimStart"
                  class="flex-1 accent-primary"
                />
                <span class="text-xs font-mono w-14 text-right shrink-0">{{ formatTime(trimStart) }}</span>
              </div>
              <div class="flex items-center gap-3">
                <span class="text-xs text-muted-foreground w-14 shrink-0">Ende</span>
                <input
                  type="range"
                  :min="0"
                  :max="fileDuration"
                  :step="0.5"
                  v-model.number="trimEnd"
                  @input="clampTrimEnd"
                  class="flex-1 accent-primary"
                />
                <span class="text-xs font-mono w-14 text-right shrink-0">{{ formatTime(trimEnd) }}</span>
              </div>
            </div>
          </div>

          <!-- Video-specific options -->
          <template v-if="showVideoOptions">
            <Separator />
            <div class="grid grid-cols-2 gap-4">
              <!-- Resize -->
              <div class="space-y-2">
                <label class="text-sm font-medium flex items-center gap-2">
                  <Monitor class="h-4 w-4" />
                  Auflösung
                </label>
                <select
                  v-model="resizePreset"
                  class="w-full h-9 px-3 text-sm border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option value="">Original</option>
                  <option value="1080">1080p</option>
                  <option value="720">720p</option>
                  <option value="480">480p</option>
                  <option value="360">360p</option>
                </select>
              </div>

              <!-- CRF (not for GIF) -->
              <div v-if="!isGifOutput" class="space-y-2">
                <label class="text-sm font-medium flex items-center gap-2">
                  <Gauge class="h-4 w-4" />
                  Qualität (CRF {{ crf }})
                </label>
                <input
                  type="range"
                  min="0" max="51" step="1"
                  v-model.number="crf"
                  class="w-full accent-primary mt-2"
                />
                <div class="flex justify-between text-xs text-muted-foreground">
                  <span>0 = beste</span>
                  <span>51 = kleinste</span>
                </div>
              </div>
            </div>
          </template>

          <!-- Audio bitrate (for audio output) -->
          <template v-if="isAudioOnlyOutput">
            <Separator />
            <div class="space-y-2">
              <label class="text-sm font-medium flex items-center gap-2">
                <Volume2 class="h-4 w-4" />
                Audio-Bitrate
              </label>
              <div class="flex gap-2">
                <button
                  v-for="br in ['128k', '192k', '256k', '320k'] as const"
                  :key="br"
                  @click="audioBitrate = br"
                  :class="[
                    'px-3 py-1.5 rounded-md text-sm border transition-colors',
                    audioBitrate === br
                      ? 'bg-primary text-primary-foreground border-primary'
                      : 'border-border hover:border-primary/50 hover:bg-muted',
                  ]"
                >{{ br }}</button>
              </div>
            </div>
          </template>

        </CardContent>
        <CardFooter class="flex flex-col gap-3 pt-0">
          <!-- Error message -->
          <div v-if="appState === 'error'" class="w-full flex items-start gap-2 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-sm text-destructive">
            <AlertCircle class="h-4 w-4 mt-0.5 shrink-0" />
            <span class="break-all">{{ errorMsg }}</span>
          </div>
          <!-- Convert button -->
          <Button class="w-full" size="lg" @click="convert">
            <Zap class="h-4 w-4 mr-2" />
            {{ isLoaded ? 'Jetzt konvertieren' : 'FFmpeg laden &amp; konvertieren' }}
          </Button>
        </CardFooter>
      </Card>

      <!-- ── 4. Progress (loading / converting) ── -->
      <Card v-if="appState === 'loading-ffmpeg' || appState === 'converting'">
        <CardContent class="p-6 space-y-4">
          <!-- Phase label -->
          <div class="flex items-center gap-3">
            <Loader2 class="h-5 w-5 animate-spin text-primary shrink-0" />
            <div class="flex-1">
              <p class="text-sm font-medium">
                <template v-if="appState === 'loading-ffmpeg'">FFmpeg wird geladen…</template>
                <template v-else-if="isGifOutput && gifPass > 0">
                  GIF Konvertierung – Pass {{ gifPass }}/2
                  <span class="text-muted-foreground font-normal ml-1">
                    ({{ gifPass === 1 ? 'Palette generieren' : 'GIF erstellen' }})
                  </span>
                </template>
                <template v-else>Konvertierung läuft…</template>
              </p>
              <p v-if="appState === 'converting'" class="text-xs text-muted-foreground mt-0.5">
                {{ outputFormat.toUpperCase() }}
                <span v-if="resizePreset && showVideoOptions"> · {{ resizePreset }}p</span>
                <span v-if="isAudioOnlyOutput"> · {{ audioBitrate }}</span>
              </p>
            </div>
            <span class="text-sm font-mono text-muted-foreground">{{ displayProgress }}%</span>
          </div>
          <Progress :model-value="displayProgress" class="h-2" />

          <!-- FFmpeg logs -->
          <div
            v-if="recentLogs.length > 0"
            class="rounded-md bg-muted/50 p-3 font-mono text-xs text-muted-foreground space-y-0.5 max-h-32 overflow-y-auto"
          >
            <div v-for="(line, i) in recentLogs" :key="i" class="leading-relaxed truncate">{{ line }}</div>
          </div>
        </CardContent>
      </Card>

      <!-- ── 5. Output (done) ── -->
      <Card v-if="appState === 'done'" class="border-green-500/30 bg-green-500/5">
        <CardHeader class="pb-3">
          <CardTitle class="text-base flex items-center gap-2 text-green-600 dark:text-green-400">
            <CheckCircle2 class="h-5 w-5" />
            Konvertierung abgeschlossen
          </CardTitle>
          <CardDescription>{{ outputFilename }} · {{ outputSizeFormatted }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <!-- Output preview -->
          <template v-if="outputUrl">
            <video
              v-if="outputMime.startsWith('video/')"
              :src="outputUrl"
              controls
              class="w-full rounded-lg max-h-64 bg-black"
            />
            <audio
              v-else-if="outputMime.startsWith('audio/')"
              :src="outputUrl"
              controls
              class="w-full"
            />
            <img
              v-else-if="outputMime === 'image/gif'"
              :src="outputUrl"
              alt="Konvertiertes GIF"
              class="max-w-full rounded-lg"
            />
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

    </template>
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
} from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'

// ── Types ──────────────────────────────────────────────────────────────────
type AppState = 'idle' | 'loading-ffmpeg' | 'converting' | 'done' | 'error'
type OutputFormat = 'mp4' | 'webm' | 'mp3' | 'wav' | 'gif'
type ResizePreset = '' | '1080' | '720' | '480' | '360'
type AudioBitrate = '128k' | '192k' | '256k' | '320k'

// ── FFmpeg ─────────────────────────────────────────────────────────────────
const { isLoaded, progress: ffmpegProgress, logs, load, writeFile, exec, readFile, deleteFile, terminate } = useFFmpeg()

// ── State ──────────────────────────────────────────────────────────────────
const appState = ref<AppState>('idle')
const errorMsg = ref('')
const gifPass = ref(0) // 1 or 2 during GIF two-pass

// ── File ───────────────────────────────────────────────────────────────────
const fileInputRef = ref<HTMLInputElement | null>(null)
const file = ref<File | null>(null)
const fileUrl = ref('')
const fileType = ref<'video' | 'audio' | null>(null)
const fileDuration = ref(0)
const isDragOver = ref(false)

// ── Options ────────────────────────────────────────────────────────────────
const outputFormat = ref<OutputFormat>('mp4')
const trimStart = ref(0)
const trimEnd = ref(0)
const resizePreset = ref<ResizePreset>('720')
const crf = ref(23)
const audioBitrate = ref<AudioBitrate>('192k')

// ── Output ─────────────────────────────────────────────────────────────────
const outputUrl = ref('')
const outputFilename = ref('')
const outputMime = ref('')
const outputSize = ref(0)

// ── Computed ───────────────────────────────────────────────────────────────
const isAudioOnlyOutput = computed(() => outputFormat.value === 'mp3' || outputFormat.value === 'wav')
const isGifOutput = computed(() => outputFormat.value === 'gif')
const showVideoOptions = computed(() => fileType.value === 'video' && !isAudioOnlyOutput.value)

const availableFormats = computed((): OutputFormat[] =>
  fileType.value === 'audio' ? ['mp3', 'wav'] : ['mp4', 'webm', 'mp3', 'wav', 'gif'],
)

const displayProgress = computed(() => {
  if (appState.value === 'loading-ffmpeg') return Math.min(ffmpegProgress.value, 99)
  if (gifPass.value === 1) return Math.round(ffmpegProgress.value * 0.5)
  if (gifPass.value === 2) return Math.round(50 + ffmpegProgress.value * 0.5)
  return ffmpegProgress.value
})

const recentLogs = computed(() => logs.value.slice(-10))

const fileSizeFormatted = computed(() => {
  if (!file.value) return ''
  const mb = file.value.size / 1024 / 1024
  return mb < 1 ? `${(mb * 1024).toFixed(0)} KB` : `${mb.toFixed(1)} MB`
})

const outputSizeFormatted = computed(() => {
  const mb = outputSize.value / 1024 / 1024
  return mb < 1 ? `${(mb * 1024).toFixed(0)} KB` : `${mb.toFixed(1)} MB`
})

const durationFormatted = computed(() => formatTime(fileDuration.value))

// ── Watchers ───────────────────────────────────────────────────────────────
watch(outputFormat, (fmt) => {
  if ((fmt === 'mp3' || fmt === 'wav') && resizePreset.value) resizePreset.value = ''
  if (fmt !== 'mp3' && fmt !== 'wav' && !resizePreset.value) resizePreset.value = '720'
})

watch(fileDuration, (d) => {
  if (d > 0) trimEnd.value = Math.floor(d)
})

// ── Helpers ────────────────────────────────────────────────────────────────
function formatTime(s: number): string {
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const sec = Math.floor(s % 60)
  if (h > 0) return `${h}:${m.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}`
  return `${m}:${sec.toString().padStart(2, '0')}`
}

const EXT_MAP: Record<OutputFormat, string> = { mp4: 'mp4', webm: 'webm', mp3: 'mp3', wav: 'wav', gif: 'gif' }
const MIME_MAP: Record<OutputFormat, string> = {
  mp4: 'video/mp4', webm: 'video/webm', mp3: 'audio/mpeg', wav: 'audio/wav', gif: 'image/gif',
}

function trimArgs(): { pre: string[]; post: string[] } {
  const hasTrim = trimStart.value > 0.5 || trimEnd.value < fileDuration.value - 0.5
  if (!hasTrim || fileDuration.value === 0) return { pre: [], post: [] }
  const duration = Math.max(0.5, trimEnd.value - trimStart.value)
  return {
    pre: ['-ss', trimStart.value.toFixed(3)],
    post: ['-t', duration.toFixed(3)],
  }
}

function buildArgs(inputName: string, outputName: string): string[] {
  const { pre, post } = trimArgs()
  const args: string[] = [...pre, '-i', inputName]

  switch (outputFormat.value) {
    case 'mp4':
      args.push('-c:v', 'libx264', '-crf', String(crf.value), '-preset', 'medium', '-movflags', '+faststart')
      if (resizePreset.value) args.push('-vf', `scale=-2:${resizePreset.value}`)
      args.push('-c:a', 'aac', '-b:a', '128k')
      break
    case 'webm':
      args.push('-c:v', 'libvpx-vp9', '-crf', String(crf.value), '-b:v', '0')
      if (resizePreset.value) args.push('-vf', `scale=-2:${resizePreset.value}`)
      args.push('-c:a', 'libopus', '-b:a', '128k')
      break
    case 'mp3':
      args.push('-vn', '-c:a', 'libmp3lame', '-b:a', audioBitrate.value, '-q:a', '2')
      break
    case 'wav':
      args.push('-vn', '-c:a', 'pcm_s16le')
      break
  }

  return [...args, ...post, '-y', outputName]
}

// ── File handling ──────────────────────────────────────────────────────────
function processFile(f: File) {
  const isVideo = f.type.startsWith('video/') || /\.(mp4|mkv|avi|mov|webm|flv|wmv|m4v)$/i.test(f.name)
  const isAudio = f.type.startsWith('audio/') || /\.(mp3|wav|aac|flac|ogg|m4a|wma)$/i.test(f.name)

  if (!isVideo && !isAudio) {
    toast.error('Nicht unterstützt', { description: 'Nur Video- und Audiodateien werden akzeptiert.' })
    return
  }

  if (f.size > 500 * 1024 * 1024) {
    toast.warning('Große Datei', {
      description: `${(f.size / 1024 / 1024).toFixed(0)} MB — Konvertierung kann sehr lange dauern.`,
    })
  }

  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)

  file.value = f
  fileUrl.value = URL.createObjectURL(f)
  fileType.value = isVideo ? 'video' : 'audio'
  fileDuration.value = 0
  trimStart.value = 0
  trimEnd.value = 0
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
  outputFormat.value = isAudio ? 'mp3' : 'mp4'
}

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  const f = e.dataTransfer?.files?.[0]
  if (f) processFile(f)
}

function handleFileInput(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files?.[0]
  if (f) processFile(f)
  input.value = ''
}

function clearFile() {
  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  file.value = null
  fileUrl.value = ''
  fileType.value = null
  fileDuration.value = 0
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
}

function handleMediaLoaded(e: Event) {
  const el = e.target as HTMLVideoElement
  if (el.duration && isFinite(el.duration)) {
    fileDuration.value = el.duration
    trimEnd.value = Math.floor(el.duration)
  }
}

function clampTrimStart() {
  if (trimStart.value >= trimEnd.value) trimStart.value = Math.max(0, trimEnd.value - 0.5)
}

function clampTrimEnd() {
  if (trimEnd.value <= trimStart.value) trimEnd.value = Math.min(fileDuration.value, trimStart.value + 0.5)
}

// ── GIF two-pass conversion ────────────────────────────────────────────────
async function convertToGif(inputName: string, outName: string): Promise<void> {
  const paletteName = 'palette.png'
  const scale = resizePreset.value ? `scale=-2:${resizePreset.value}` : 'scale=480:-1'
  const { pre, post } = trimArgs()

  // Pass 1: palettegen
  gifPass.value = 1
  const exitCode1 = await exec([
    ...pre, '-i', inputName, ...post,
    '-vf', `fps=10,${scale}:flags=lanczos,palettegen=stats_mode=diff`,
    '-y', paletteName,
  ])
  if (exitCode1 !== 0) throw new Error(`GIF Palettegen fehlgeschlagen (Exit-Code ${exitCode1})`)

  // Pass 2: paletteuse
  gifPass.value = 2
  const exitCode2 = await exec([
    ...pre, '-i', inputName,
    '-i', paletteName,
    ...post,
    '-filter_complex', `[0:v]fps=10,${scale}:flags=lanczos[x];[x][1:v]paletteuse=dither=bayer:bayer_scale=5:diff_mode=rectangle`,
    '-loop', '0',
    '-y', outName,
  ])
  if (exitCode2 !== 0) throw new Error(`GIF Paletteuse fehlgeschlagen (Exit-Code ${exitCode2})`)

  await deleteFile(paletteName)
  gifPass.value = 0
}

// ── Main convert function ──────────────────────────────────────────────────
async function convert() {
  if (!file.value) return

  // Load FFmpeg on first use
  if (!isLoaded.value) {
    appState.value = 'loading-ffmpeg'
    try {
      await load()
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      appState.value = 'error'
      errorMsg.value = `FFmpeg konnte nicht geladen werden: ${msg}`
      toast.error('Ladefehler', { description: errorMsg.value })
      return
    }
  }

  appState.value = 'converting'
  gifPass.value = 0

  const inputExt = file.value.name.split('.').pop() ?? 'bin'
  const inputName = `input.${inputExt}`
  const outExt = EXT_MAP[outputFormat.value]
  const baseName = file.value.name.replace(/\.[^.]+$/, '')
  const outName = `output.${outExt}`

  try {
    // Write input to MEMFS
    const inputData = await fetchFile(file.value)
    await writeFile(inputName, inputData)

    // Exec conversion
    if (isGifOutput.value) {
      await convertToGif(inputName, outName)
    } else {
      const args = buildArgs(inputName, outName)
      console.log('[FFmpeg] exec:', args.join(' '))
      const exitCode = await exec(args)
      if (exitCode !== 0) throw new Error(`FFmpeg beendet mit Exit-Code ${exitCode}`)
    }

    // Read output
    const outData = await readFile(outName)
    // .slice(0) produces Uint8Array<ArrayBuffer> (avoids SharedArrayBuffer TS error in Blob ctor)
    const blob = new Blob([outData.slice(0)], { type: MIME_MAP[outputFormat.value] })

    // Cleanup MEMFS
    await deleteFile(inputName)
    await deleteFile(outName)

    if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
    outputUrl.value = URL.createObjectURL(blob)
    outputFilename.value = `${baseName}_converted.${outExt}`
    outputMime.value = MIME_MAP[outputFormat.value]
    outputSize.value = blob.size

    appState.value = 'done'
    toast.success('Fertig!', { description: `${outputFilename.value} (${outputSizeFormatted.value})` })
    downloadOutput()
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    errorMsg.value = msg
    appState.value = 'error'
    toast.error('Konvertierungsfehler', { description: msg })
    console.error('[FileConverter]', err)
  }
}

// ── Output actions ─────────────────────────────────────────────────────────
function downloadOutput() {
  if (!outputUrl.value) return
  const a = document.createElement('a')
  a.href = outputUrl.value
  a.download = outputFilename.value
  a.click()
}

function resetToConvert() {
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  outputUrl.value = ''
  outputFilename.value = ''
  appState.value = 'idle'
  errorMsg.value = ''
}

// ── Cleanup ────────────────────────────────────────────────────────────────
onUnmounted(() => {
  if (fileUrl.value) URL.revokeObjectURL(fileUrl.value)
  if (outputUrl.value) URL.revokeObjectURL(outputUrl.value)
  terminate()
})
</script>