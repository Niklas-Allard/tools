import { ref, type Ref } from 'vue'
import { FFmpeg } from '@ffmpeg/ffmpeg'
import { toBlobURL } from '@ffmpeg/util'

const CDN = 'https://unpkg.com/@ffmpeg/core-mt@0.12.10/dist/esm'

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
}

export function useFFmpeg(): UseFFmpegReturn {
  const ffmpeg = new FFmpeg()
  const isLoaded = ref(false)
  const isLoading = ref(false)
  const progress = ref(0)
  const logs = ref<string[]>([])

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

  return { isLoaded, isLoading, progress, logs, load, writeFile, exec, readFile, deleteFile, terminate }
}