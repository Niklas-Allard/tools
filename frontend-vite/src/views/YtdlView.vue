<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="text-center space-y-2">
      <h1 class="text-5xl font-bold tracking-tight">
        <span class="bg-gradient-to-r from-red-600 to-purple-600 bg-clip-text text-transparent">
          🎬 YouTube Downloader Pro
        </span>
      </h1>
      <p class="text-muted-foreground text-lg">Ultimate yt-dlp with 25+ Features</p>
    </div>

    <!-- Download Form Card -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <span class="text-2xl">🚀</span>
          Neuer Download
        </CardTitle>
        <CardDescription>
          YouTube Videos und Playlists mit erweiterten Optionen herunterladen
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-6">
        <!-- Mode Toggle -->
        <div class="flex gap-2">
          <Button
            @click="batchMode = false"
            :variant="!batchMode ? 'default' : 'outline'"
            size="sm"
          >🔗 Single URL</Button>
          <Button
            @click="batchMode = true"
            :variant="batchMode ? 'default' : 'outline'"
            size="sm"
          >📋 Batch URLs</Button>
        </div>

        <!-- Single URL Input -->
        <div v-if="!batchMode" class="space-y-2">
          <Label for="url">YouTube URL oder Playlist</Label>
          <Input
            id="url"
            v-model="form.url"
            type="text"
            placeholder="https://www.youtube.com/watch?v=..."
            class="text-base"
          />
        </div>

        <!-- Batch URL Input -->
        <div v-else class="space-y-2">
          <Label for="batch_urls" class="flex items-center gap-2">
            Batch URLs
            <span class="text-xs text-muted-foreground">(eine URL pro Zeile)</span>
          </Label>
          <textarea
            id="batch_urls"
            v-model="batchUrls"
            placeholder="https://www.youtube.com/watch?v=AAA&#10;https://www.youtube.com/watch?v=BBB&#10;https://www.youtube.com/watch?v=CCC"
            rows="5"
            class="w-full px-3 py-2 text-sm border rounded-md bg-background text-foreground resize-y focus:outline-none focus:ring-2 focus:ring-ring"
          />
          <p class="text-xs text-muted-foreground">
            {{ batchUrlCount }} URL{{ batchUrlCount !== 1 ? 's' : '' }} erkannt — alle Downloads landen im selben Ordner
          </p>
        </div>

        <Separator />

        <!-- 1️⃣ VIDEO/AUDIO Section -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <h3 class="text-lg font-semibold">🎥 Video & Audio</h3>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <!-- Audio Only -->
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="audio_only" v-model="form.audio_only" />
              <Label for="audio_only" class="cursor-pointer font-normal flex items-center gap-1">
                🎵 Nur Audio
                <Tooltip>
                  <TooltipTrigger as-child>
                    <span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span>
                  </TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">Nur Audio extrahieren (ignoriert Video-Streams). Praktisch für Podcasts oder Musik-Downloads.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
            </div>

            <!-- Audio Format -->
            <div v-if="form.audio_only" class="space-y-2">
              <Label class="flex items-center gap-1">
                Audio Format
                <Tooltip>
                  <TooltipTrigger as-child>
                    <span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span>
                  </TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">Zielformat für Audio-Dateien. Konvertiert automatisch mit ffmpeg.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
              <Select v-model:model-value="form.audio_format">
                <SelectTrigger>
                  <SelectValue placeholder="Format wählen" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="mp3">MP3</SelectItem>
                  <SelectItem value="m4a">M4A</SelectItem>
                  <SelectItem value="opus">Opus</SelectItem>
                  <SelectItem value="wav">WAV</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- Video Quality -->
            <div v-if="!form.audio_only" class="space-y-2">
              <Label class="flex items-center gap-1">
                Video Qualität
                <Tooltip>
                  <TooltipTrigger as-child>
                    <span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span>
                  </TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">Video-Höhe einschränken. yt-dlp wählt das beste verfügbare Format.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
              <Select v-model:model-value="form.quality">
                <SelectTrigger>
                  <SelectValue placeholder="Qualität wählen" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="best">Beste verfügbar</SelectItem>
                  <SelectItem value="2160">4K (2160p)</SelectItem>
                  <SelectItem value="1440">2K (1440p)</SelectItem>
                  <SelectItem value="1080">Full HD (1080p)</SelectItem>
                  <SelectItem value="720">HD (720p)</SelectItem>
                  <SelectItem value="480">SD (480p)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- FPS -->
            <div v-if="!form.audio_only" class="space-y-2">
              <Label class="flex items-center gap-1">
                FPS (optional)
                <Tooltip>
                  <TooltipTrigger as-child>
                    <span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span>
                  </TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">FPS-Filter für Videos. Beschränkt auf Streams mit exakt dieser Bildrate.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
              <Select v-model:model-value="form.fps">
                <SelectTrigger>
                  <SelectValue placeholder="Auto" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="auto">Auto</SelectItem>
                  <SelectItem value="60">60 FPS</SelectItem>
                  <SelectItem value="30">30 FPS</SelectItem>
                  <SelectItem value="24">24 FPS</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- Playlist -->
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="playlist" v-model="form.playlist" />
              <Label for="playlist" class="cursor-pointer font-normal flex items-center gap-1">
                📋 Ganze Playlist
                <Tooltip>
                  <TooltipTrigger as-child>
                    <span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span>
                  </TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">Gesamte Playlist herunterladen, nicht nur einzelnes Video.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
            </div>
          </div>
        </div>

        <Separator />

        <!-- 2️⃣ EMBED Section -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <h3 class="text-lg font-semibold">✨ Embedding</h3>
            <Badge variant="success" class="bg-green-500 hover:bg-green-600">
              ✅ Empfohlen
            </Badge>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="embed_chapters" v-model="form.embed_chapters" />
              <Label for="embed_chapters" class="cursor-pointer font-normal flex items-center gap-1">
                📑 Chapters
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Kapitelinformationen in Video einbetten.</p></TooltipContent></Tooltip>
              </Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="embed_metadata" v-model="form.embed_metadata" />
              <Label for="embed_metadata" class="cursor-pointer font-normal flex items-center gap-1">
                📝 Metadata
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Titel, Beschreibung, Upload-Datum als Metadaten einbetten.</p></TooltipContent></Tooltip>
              </Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="embed_thumbnail" v-model="form.embed_thumbnail" />
              <Label for="embed_thumbnail" class="cursor-pointer font-normal flex items-center gap-1">
                🖼️ Thumbnail
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Thumbnail direkt ins Video einbetten (Cover-Art).</p></TooltipContent></Tooltip>
              </Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="embed_subtitles" v-model="form.embed_subtitles" />
              <Label for="embed_subtitles" class="cursor-pointer font-normal flex items-center gap-1">
                💬 Subtitles
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Untertitel direkt ins Video einbetten.</p></TooltipContent></Tooltip>
              </Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="split_chapters" v-model="form.split_chapters" />
              <Label for="split_chapters" class="cursor-pointer font-normal flex items-center gap-1">
                ✂️ Split Chapters
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Video an Kapitelgrenzen in separate Dateien aufteilen.</p></TooltipContent></Tooltip>
              </Label>
            </div>
          </div>
        </div>

        <Separator />

        <!-- 3️⃣ SUBTITLES & FILES Section -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">💬 Untertitel & Dateien</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="space-y-2">
              <Label class="flex items-center gap-1">
                Untertitel Sprache
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Sprachen für Untertitel (Komma-separiert: 'de,en').</p></TooltipContent></Tooltip>
              </Label>
              <Select v-model:model-value="form.subtitles_lang">
                <SelectTrigger><SelectValue placeholder="Keine" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Keine</SelectItem>
                  <SelectItem value="de">🇩🇪 Deutsch</SelectItem>
                  <SelectItem value="en">🇬🇧 English</SelectItem>
                  <SelectItem value="de,en">🇩🇪 + 🇬🇧 Beide</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="write_info_json" v-model="form.write_info_json" />
              <Label for="write_info_json" class="cursor-pointer font-normal">📄 Info JSON</Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="write_description" v-model="form.write_description" />
              <Label for="write_description" class="cursor-pointer font-normal">📋 Description</Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="write_thumbnail" v-model="form.write_thumbnail" />
              <Label for="write_thumbnail" class="cursor-pointer font-normal">🖼️ Thumbnail Datei</Label>
            </div>
          </div>
        </div>

        <Separator />

        <!-- 4️⃣ LIMITS Section -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">⚙️ Limits & Einstellungen</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="space-y-2">
              <Label for="max_filesize" class="flex items-center gap-1">
                Max Dateigröße
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Maximale Dateigröße (z. B. '500M', '2G').</p></TooltipContent></Tooltip>
              </Label>
              <Input id="max_filesize" v-model="form.max_filesize" placeholder="z.B. 1G, 500M" />
            </div>
            <div class="space-y-2">
              <Label for="min_filesize" class="flex items-center gap-1">
                Min Dateigröße
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Minimale Dateigröße erzwingen.</p></TooltipContent></Tooltip>
              </Label>
              <Input id="min_filesize" v-model="form.min_filesize" placeholder="z.B. 100M, 10M" />
            </div>
            <div class="space-y-2">
              <Label for="max_downloads" class="flex items-center gap-1">
                Max Downloads (Playlist)
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Maximale Anzahl Videos aus Playlist herunterladen.</p></TooltipContent></Tooltip>
              </Label>
              <Input id="max_downloads" v-model.number="form.max_downloads" type="number" placeholder="z.B. 10" />
            </div>
            <div class="space-y-2">
              <Label for="rate_limit" class="flex items-center gap-1">
                Rate Limit
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Download-Geschwindigkeit begrenzen (z. B. '1M', '500K').</p></TooltipContent></Tooltip>
              </Label>
              <Input id="rate_limit" v-model="form.rate_limit" placeholder="z.B. 1M, 500K" />
            </div>
          </div>
        </div>

        <Separator />

        <!-- 5️⃣ PREMIUM Section -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">🔥 Premium Features</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div class="space-y-2">
              <Label class="flex items-center gap-1">
                Cookies from Browser
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Cookies aus Browser für age-gated/premium Inhalte.</p></TooltipContent></Tooltip>
              </Label>
              <Select v-model:model-value="form.cookies_from_browser">
                <SelectTrigger><SelectValue placeholder="Keine" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Keine</SelectItem>
                  <SelectItem value="chrome">Chrome</SelectItem>
                  <SelectItem value="firefox">Firefox</SelectItem>
                  <SelectItem value="edge">Edge</SelectItem>
                  <SelectItem value="safari">Safari</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border-2 border-purple-500 rounded-md bg-purple-50 dark:bg-purple-950">
              <Checkbox id="sponsorblock" v-model="form.sponsorblock" class="data-[state=checked]:bg-purple-600" />
              <Label for="sponsorblock" class="cursor-pointer font-semibold text-purple-700 dark:text-purple-400 flex items-center gap-1">
                ⚡ SponsorBlock
                <Tooltip><TooltipTrigger as-child><span class="text-purple-500 dark:text-purple-300 cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">SponsorBlock aktivieren – entfernt Sponsoren-Segmente automatisch.</p></TooltipContent></Tooltip>
              </Label>
            </div>
            <div class="flex items-center space-x-2 h-10 px-3 border rounded-md bg-card">
              <Checkbox id="audio_multistreams" v-model="form.audio_multistreams" />
              <Label for="audio_multistreams" class="cursor-pointer font-normal flex items-center gap-1">
                🎧 Audio Multistreams
                <Tooltip><TooltipTrigger as-child><span class="text-muted-foreground hover:text-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger><TooltipContent class="max-w-xs"><p class="text-sm">Alle verfügbaren Audio-Tracks beibehalten.</p></TooltipContent></Tooltip>
              </Label>
            </div>
          </div>
        </div>

        <Separator />

        <!-- Start Button -->
        <Button
          @click="batchMode ? startBatchDownload() : startDownload()"
          :disabled="(batchMode ? batchUrlCount === 0 : !form.url) || loading"
          class="w-full h-12 text-base font-semibold bg-gradient-to-r from-red-600 to-purple-600 hover:from-red-700 hover:to-purple-700"
          size="lg"
        >
          <span v-if="loading">⏳ Wird gestartet...</span>
          <span v-else-if="batchMode">🚀 Batch Download starten ({{ batchUrlCount }} URLs)</span>
          <span v-else>🚀 Download starten</span>
        </Button>
      </CardContent>
    </Card>

    <!-- Downloads List -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle class="flex items-center gap-2">
            <span class="text-2xl">📥</span>
            Downloads
          </CardTitle>
          <Badge variant="secondary" class="text-sm">{{ downloads.length }}</Badge>
        </div>
        <CardDescription>Übersicht über alle laufenden und abgeschlossenen Downloads</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="downloads.length === 0" class="text-center py-12">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-muted-foreground text-lg">Keine Downloads vorhanden</p>
        </div>
        <div v-else class="space-y-3">
          <template v-for="item in groupedDownloads" :key="item.type === 'single' ? item.download.id : item.batchId">

            <!-- Single Download -->
            <Card
              v-if="item.type === 'single'"
              class="overflow-hidden transition-all hover:shadow-md"
              :class="{
                'border-l-4 border-l-green-500': item.download.status === 'completed',
                'border-l-4 border-l-blue-500': item.download.status === 'downloading',
                'border-l-4 border-l-yellow-500': item.download.status === 'pending',
                'border-l-4 border-l-red-500': item.download.status === 'failed'
              }"
            >
              <CardContent class="p-4">
                <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                  <div class="flex-1 min-w-0 space-y-2">
                    <div class="flex items-center gap-2 flex-wrap">
                      <h3 class="font-semibold truncate">{{ item.download.title || item.download.url }}</h3>
                      <Badge
                        :variant="item.download.status === 'completed' ? 'success' : item.download.status === 'downloading' ? 'default' : item.download.status === 'pending' ? 'secondary' : 'destructive'"
                        class="shrink-0"
                      >
                        {{ statusText(item.download.status) }}
                      </Badge>
                    </div>
                    <p class="text-sm text-muted-foreground truncate">{{ item.download.url }}</p>
                    <p class="text-xs text-muted-foreground">{{ formatDate(item.download.created_at) }}</p>
                    <div v-if="item.download.status === 'downloading'" class="space-y-1">
                      <Progress :model-value="item.download.progress" class="h-2" />
                      <p class="text-xs text-muted-foreground">{{ item.download.progress }}%</p>
                    </div>
                  </div>
                  <div v-if="item.download.status === 'completed' || item.download.status === 'failed'" class="flex gap-2 flex-wrap shrink-0">
                    <Button
                      v-if="item.download.status === 'completed' && !item.download.is_playlist"
                      @click="downloadFile(item.download.job_id)"
                      variant="default" size="sm"
                      class="bg-blue-600 hover:bg-blue-700"
                    >📥 Datei</Button>
                    <Button
                      v-if="item.download.status === 'completed'"
                      @click="downloadFolder(item.download.job_id)"
                      variant="default" size="sm"
                      class="bg-purple-600 hover:bg-purple-700"
                    >📦 {{ item.download.is_playlist ? 'Playlist' : 'Ordner' }}</Button>
                    <AlertDialog>
                      <AlertDialogTrigger as-child>
                        <Button variant="destructive" size="sm">🗑️</Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Download wirklich löschen?</AlertDialogTitle>
                          <AlertDialogDescription>Diese Aktion kann nicht rückgängig gemacht werden.</AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Abbrechen</AlertDialogCancel>
                          <AlertDialogAction @click="deleteDownload(item.download.job_id)" class="bg-destructive text-destructive-foreground hover:bg-destructive/90">Löschen</AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Batch Group -->
            <Card v-else class="overflow-hidden border-l-4 border-l-orange-500">
              <CardContent class="p-4 space-y-3">
                <!-- Batch Header -->
                <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="text-lg font-semibold">🗂️ Batch</span>
                    <Badge variant="secondary" class="font-mono text-xs">{{ item.batchId.slice(-8) }}</Badge>
                    <Badge
                      :class="{
                        'bg-green-500 hover:bg-green-600': item.completed === item.total,
                        'bg-blue-500 hover:bg-blue-600': item.completed < item.total && item.failed === 0,
                        'bg-red-500 hover:bg-red-600': item.failed > 0 && item.completed + item.failed === item.total
                      }"
                    >
                      {{ item.completed }}/{{ item.total }} fertig
                      <span v-if="item.failed > 0" class="ml-1">({{ item.failed }} Fehler)</span>
                    </Badge>
                    <p class="text-xs text-muted-foreground">{{ formatDate(item.downloads[0].created_at) }}</p>
                  </div>
                  <div class="flex gap-2 flex-wrap shrink-0">
                    <Button
                      v-if="item.completed > 0"
                      @click="downloadBatch(item.batchId)"
                      variant="default" size="sm"
                      class="bg-orange-600 hover:bg-orange-700"
                    >📦 Batch ZIP ({{ item.completed }} Dateien)</Button>
                    <AlertDialog>
                      <AlertDialogTrigger as-child>
                        <Button variant="destructive" size="sm">🗑️ Batch</Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Ganzen Batch löschen?</AlertDialogTitle>
                          <AlertDialogDescription>Alle {{ item.total }} Downloads und Dateien werden permanent gelöscht.</AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Abbrechen</AlertDialogCancel>
                          <AlertDialogAction @click="deleteBatch(item.batchId)" class="bg-destructive text-destructive-foreground hover:bg-destructive/90">Löschen</AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>

                <!-- Batch Items -->
                <div class="space-y-2 pl-2 border-l-2 border-orange-200 dark:border-orange-900">
                  <div
                    v-for="dl in item.downloads"
                    :key="dl.id"
                    class="flex flex-col md:flex-row md:items-center md:justify-between gap-2 p-2 rounded-md bg-muted/30"
                  >
                    <div class="flex-1 min-w-0 space-y-1">
                      <div class="flex items-center gap-2 flex-wrap">
                        <span class="text-sm font-medium truncate">{{ dl.title || dl.url }}</span>
                        <Badge
                          :variant="dl.status === 'completed' ? 'success' : dl.status === 'downloading' ? 'default' : dl.status === 'pending' ? 'secondary' : 'destructive'"
                          class="shrink-0 text-xs"
                        >{{ statusText(dl.status) }}</Badge>
                      </div>
                      <p class="text-xs text-muted-foreground truncate">{{ dl.url }}</p>
                      <div v-if="dl.status === 'downloading'" class="space-y-1">
                        <Progress :model-value="dl.progress" class="h-1.5" />
                        <p class="text-xs text-muted-foreground">{{ dl.progress }}%</p>
                      </div>
                    </div>
                    <Button
                      v-if="dl.status === 'completed' && !dl.is_playlist"
                      @click="downloadFile(dl.job_id)"
                      variant="outline" size="sm"
                    >📥</Button>
                  </div>
                </div>
              </CardContent>
            </Card>

          </template>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Progress } from '@/components/ui/progress'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'

const API_BASE = '/api/ytdl'

const form = ref({
  url: '',
  audio_only: false,
  audio_format: 'mp3',
  quality: 'best',
  fps: 'auto',
  playlist: false,
  embed_chapters: true,
  embed_metadata: true,
  embed_thumbnail: true,
  embed_subtitles: false,
  split_chapters: false,
  subtitles_lang: 'none',
  write_info_json: false,
  write_description: false,
  write_thumbnail: false,
  max_filesize: '',
  min_filesize: '',
  max_downloads: 0,
  rate_limit: '',
  cookies_from_browser: 'none',
  sponsorblock: false,
  audio_multistreams: false,
  remux_video: '',
  output_template: ''
})

const loading = ref(false)
const downloads = ref([])
const batchMode = ref(false)
const batchUrls = ref('')
let refreshInterval = null

const batchUrlCount = computed(() =>
  batchUrls.value.split('\n').map(u => u.trim()).filter(u => u.length > 0).length
)

const groupedDownloads = computed(() => {
  const result = []
  const batchMap = new Map()
  for (const dl of downloads.value) {
    if (dl.batch_id) {
      if (!batchMap.has(dl.batch_id)) {
        const group = { type: 'batch', batchId: dl.batch_id, downloads: [], completed: 0, failed: 0, total: 0 }
        batchMap.set(dl.batch_id, group)
        result.push(group)
      }
      const group = batchMap.get(dl.batch_id)
      group.downloads.push(dl)
      group.total++
      if (dl.status === 'completed') group.completed++
      if (dl.status === 'failed') group.failed++
    } else {
      result.push({ type: 'single', download: dl })
    }
  }
  return result
})

async function startDownload() {
  if (!form.value.url) return
  loading.value = true
  try {
    const response = await axios.post(`${API_BASE}/download`, {
      url: form.value.url,
      audio_only: form.value.audio_only,
      audio_format: form.value.audio_format,
      quality: form.value.quality === 'best' ? '' : form.value.quality,
      fps: form.value.fps === 'auto' ? '' : form.value.fps,
      playlist: form.value.playlist,
      embed_chapters: form.value.embed_chapters,
      embed_metadata: form.value.embed_metadata,
      embed_thumbnail: form.value.embed_thumbnail,
      embed_subtitles: form.value.embed_subtitles,
      split_chapters: form.value.split_chapters,
      subtitles_lang: form.value.subtitles_lang === 'none' ? '' : form.value.subtitles_lang,
      write_info_json: form.value.write_info_json,
      write_description: form.value.write_description,
      write_thumbnail: form.value.write_thumbnail,
      max_filesize: form.value.max_filesize,
      min_filesize: form.value.min_filesize,
      max_downloads: form.value.max_downloads,
      rate_limit: form.value.rate_limit,
      cookies_from_browser: form.value.cookies_from_browser === 'none' ? '' : form.value.cookies_from_browser,
      sponsorblock: form.value.sponsorblock,
      audio_multistreams: form.value.audio_multistreams,
      remux_video: form.value.remux_video,
      output_template: form.value.output_template
    })
    toast.success('✅ Download gestartet!', { description: `Job ID: ${response.data.job_id}` })
    form.value.url = ''
    await loadDownloads()
  } catch (error) {
    toast.error('❌ Fehler beim Download', { description: error.response?.data?.error || error.message })
  } finally {
    loading.value = false
  }
}

async function startBatchDownload() {
  const urls = batchUrls.value.split('\n').map(u => u.trim()).filter(u => u.length > 0)
  if (urls.length === 0) return
  loading.value = true
  try {
    const response = await axios.post(`${API_BASE}/batch-download`, {
      urls,
      audio_only: form.value.audio_only,
      audio_format: form.value.audio_format,
      quality: form.value.quality === 'best' ? '' : form.value.quality,
      fps: form.value.fps === 'auto' ? '' : form.value.fps,
      playlist: form.value.playlist,
      embed_chapters: form.value.embed_chapters,
      embed_metadata: form.value.embed_metadata,
      embed_thumbnail: form.value.embed_thumbnail,
      embed_subtitles: form.value.embed_subtitles,
      split_chapters: form.value.split_chapters,
      subtitles_lang: form.value.subtitles_lang === 'none' ? '' : form.value.subtitles_lang,
      write_info_json: form.value.write_info_json,
      write_description: form.value.write_description,
      write_thumbnail: form.value.write_thumbnail,
      max_filesize: form.value.max_filesize,
      min_filesize: form.value.min_filesize,
      max_downloads: form.value.max_downloads,
      rate_limit: form.value.rate_limit,
      cookies_from_browser: form.value.cookies_from_browser === 'none' ? '' : form.value.cookies_from_browser,
      sponsorblock: form.value.sponsorblock,
      audio_multistreams: form.value.audio_multistreams,
      remux_video: form.value.remux_video,
      output_template: form.value.output_template
    })
    toast.success('✅ Batch Download gestartet!', { description: response.data.message })
    batchUrls.value = ''
    await loadDownloads()
  } catch (error) {
    toast.error('❌ Fehler beim Batch Download', { description: error.response?.data?.error || error.message })
  } finally {
    loading.value = false
  }
}

async function loadDownloads() {
  try {
    const res = await axios.get(`${API_BASE}/downloads`)
    downloads.value = res.data || []
  } catch (error) {
    console.error('Failed to load downloads:', error)
  }
}

function downloadFile(jobId) {
  window.location.href = `${API_BASE}/download-file/${jobId}`
}

function downloadFolder(jobId) {
  window.location.href = `${API_BASE}/download-folder/${jobId}`
}

function downloadBatch(batchId) {
  window.location.href = `${API_BASE}/download-batch/${batchId}`
}

async function deleteDownload(jobId) {
  try {
    await axios.delete(`${API_BASE}/download/${jobId}`)
    toast.success('🗑️ Download gelöscht')
    await loadDownloads()
  } catch (error) {
    toast.error('❌ Fehler beim Löschen', { description: error.message })
  }
}

async function deleteBatch(batchId) {
  try {
    await axios.delete(`${API_BASE}/batch/${batchId}`)
    toast.success('🗑️ Batch gelöscht')
    await loadDownloads()
  } catch (error) {
    toast.error('❌ Fehler beim Löschen', { description: error.message })
  }
}

function statusText(status) {
  const texts = { pending: '⏳ Wartend', downloading: '⬇️ Lädt...', completed: '✅ Fertig', failed: '❌ Fehler' }
  return texts[status] || status
}

function formatDate(date) {
  return new Date(date).toLocaleString('de-DE', {
    day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit'
  })
}

onMounted(() => {
  loadDownloads()
  refreshInterval = setInterval(loadDownloads, 3000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>