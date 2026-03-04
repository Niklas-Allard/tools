import { ref, type Ref } from 'vue'
import { FFmpeg } from '@ffmpeg/ffmpeg'
import { toBlobURL } from '@ffmpeg/util'

const CDN = 'https://unpkg.com/@ffmpeg/core-mt@0.12.10/dist/esm'

// ── Image conversion options ────────────────────────────────────────────────
export interface ImageConvertOptions {
  /** Perceptual quality 0–100 (100 = best). Mapped to codec-specific values.
   *  Ignored for lossless formats (PNG, BMP). */
  quality?: number
  /** Target width in pixels. 0 = auto (maintain aspect ratio). */
  width?: number
  /** Target height in pixels. 0 = auto (maintain aspect ratio). */
  height?: number
}

// ── Known output-extension → codec mapping (FFmpeg.wasm core-mt) ───────────
// Note: libaom-av1 (AVIF) and libx265 are NOT in the standard core-mt build.
// They are listed here for completeness but will produce an FFmpeg error.
const IMAGE_CODEC_MAP: Record<string, string> = {
  jpg: 'mjpeg',
  jpeg: 'mjpeg',
  png: 'png',
  webp: 'libwebp',
  bmp: 'bmp',
  gif: 'gif',
  avif: 'libaom-av1', // ⚠ experimental – may fail with core-mt
  tiff: 'tiff',
}

// ── Public interface ────────────────────────────────────────────────────────
export interface UseFFmpegReturn {
  isLoaded: Ref<boolean>
  isLoading: Ref<boolean>
  progress: Ref<number>
  logs: Ref<string[]>
  load: () => Promise<void>
  writeFile: (name: string, data: Uint8Array) => Promise<void>
  exec: (args: string[]) => Promise<number>
  readFile: (name: string) => Promise<Uint8Array>
  deleteFile: (name: string) => Promise<void>
  terminate: () => void
  /** Convert a still image that already lives in MEMFS (inputName) to another
   *  image format (determined by outputName's extension). Throws on failure. */
  convertImage: (inputName: string, outputName: string, options?: ImageConvertOptions) => Promise<void>
}

// ── Composable ─────────────────────────────────────────────────────────────
export function useFFmpeg(): UseFFmpegReturn {
  const ffmpeg = new FFmpeg()
  const isLoaded = ref(false)
  const isLoading = ref(false)
  const progress = ref(0)
  const logs = ref<string[]>([])

  // ── Load ─────────────────────────────────────────────────────────────────
  const load = async (): Promise<void> => {
    if (isLoaded.value) return

    // SharedArrayBuffer requires cross-origin isolation.
    // The Vite dev server must be running with the COOP/COEP headers
    // configured in vite.config.ts — restart it if you just added them.
    if (!self.crossOriginIsolated) {
      throw new Error(
        'Seite ist nicht cross-origin isoliert. Bitte den Dev-Server neu starten (die COOP/COEP-Header aus vite.config.ts werden erst nach einem Neustart aktiv).',
      )
    }

    isLoading.value = true
    logs.value = []

    ffmpeg.on('log', ({ message }) => {
      logs.value.push(message)
      if (logs.value.length > 300) logs.value.splice(0, 100)
    })

    ffmpeg.on('progress', ({ progress: p }) => {
      progress.value = Math.max(0, Math.min(100, Math.round(p * 100)))
    })

    try {
      await ffmpeg.load({
        coreURL: await toBlobURL(`${CDN}/ffmpeg-core.js`, 'text/javascript'),
        wasmURL: await toBlobURL(`${CDN}/ffmpeg-core.wasm`, 'application/wasm'),
        workerURL: await toBlobURL(`${CDN}/ffmpeg-core.worker.js`, 'text/javascript'),
      })
      isLoaded.value = true
    } finally {
      isLoading.value = false
    }
  }

  // ── Primitives ────────────────────────────────────────────────────────────
  const writeFile = async (name: string, data: Uint8Array): Promise<void> => {
    await ffmpeg.writeFile(name, data)
  }

  const exec = async (args: string[]): Promise<number> => {
    return await ffmpeg.exec(args)
  }

  const readFile = async (name: string): Promise<Uint8Array> => {
    const data = await ffmpeg.readFile(name)
    return data as Uint8Array
  }

  const deleteFile = async (name: string): Promise<void> => {
    try {
      await ffmpeg.deleteFile(name)
    } catch {
      // file may not exist
    }
  }

  const terminate = (): void => {
    try {
      ffmpeg.terminate()
    } catch {
      // already terminated
    }
    isLoaded.value = false
    progress.value = 0
  }

  // ── Image helper ──────────────────────────────────────────────────────────
  /**
   * Convert a still image in MEMFS. Detects codec from outputName extension.
   *
   * Supported codecs (core-mt build):
   *   jpg/jpeg → mjpeg    (quality: q:v 2-31, lower = better)
   *   png      → png      (lossless; quality controls compression_level)
   *   webp     → libwebp  (quality: q:v 0-100)
   *   bmp      → bmp      (lossless)
   *   gif      → gif      (use component's two-pass helper for animated GIF)
   *   avif     → libaom-av1  ⚠ may not be available in core-mt
   *   tiff     → tiff     (lossless)
   */
  const convertImage = async (
    inputName: string,
    outputName: string,
    options: ImageConvertOptions = {},
  ): Promise<void> => {
    const { quality = 85, width = 0, height = 0 } = options
    const outExt = outputName.split('.').pop()?.toLowerCase() ?? ''
    const codec = IMAGE_CODEC_MAP[outExt]

    if (!codec) throw new Error(`Unbekanntes Bildformat: .${outExt}`)

    const args: string[] = ['-i', inputName]

    // Resize filter
    const scaleW = width > 0 ? String(width) : '-1'
    const scaleH = height > 0 ? String(height) : '-1'
    const needsScale = width > 0 || height > 0
    // For images we always need to emit exactly 1 frame
    const vfParts: string[] = []
    if (needsScale) vfParts.push(`scale=${scaleW}:${scaleH}`)
    if (vfParts.length) args.push('-vf', vfParts.join(','))

    // Force single frame output (essential for still images)
    args.push('-frames:v', '1')

    // Codec
    args.push('-c:v', codec)

    // Quality mapping per codec
    switch (codec) {
      case 'mjpeg': {
        // q:v range: 2 (best) – 31 (worst)
        const qv = Math.round(2 + ((100 - quality) / 100) * 29)
        args.push('-q:v', String(Math.max(2, Math.min(31, qv))))
        break
      }
      case 'libwebp': {
        // q:v range: 0 (worst) – 100 (best)
        args.push('-q:v', String(Math.max(0, Math.min(100, quality))))
        break
      }
      case 'png': {
        // compression_level: 0 (none) – 9 (max); quality 100 → level 0
        const lvl = Math.round((1 - quality / 100) * 9)
        args.push('-compression_level', String(lvl))
        break
      }
      case 'libaom-av1': {
        // For AVIF: crf 0 (lossless) – 63 (worst)
        const crf = Math.round(63 - (quality / 100) * 63)
        args.push('-crf', String(crf), '-b:v', '0', '-still-picture', '1')
        break
      }
      // bmp, gif, tiff: no quality param
    }

    // Strip audio stream (images have none, but prevents "audio stream not found" warnings)
    args.push('-an', '-y', outputName)

    const exitCode = await ffmpeg.exec(args)
    if (exitCode !== 0) {
      throw new Error(`Bild-Konvertierung fehlgeschlagen (Exit-Code ${exitCode})`)
    }
  }

  return {
    isLoaded,
    isLoading,
    progress,
    logs,
    load,
    writeFile,
    exec,
    readFile,
    deleteFile,
    terminate,
    convertImage,
  }
}
