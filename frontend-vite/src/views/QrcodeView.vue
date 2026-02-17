<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="text-center space-y-2">
      <h1 class="text-5xl font-bold tracking-tight">
        <span class="bg-gradient-to-r from-blue-600 to-cyan-500 bg-clip-text text-transparent">
          📱 QR Code Generator
        </span>
      </h1>
      <p class="text-muted-foreground text-lg">Erstelle QR Codes für URLs, WiFi, Kontakte und mehr</p>
    </div>

    <!-- Main Layout: Config Left + Preview Right -->
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">

      <!-- LEFT: Configuration -->
      <div class="lg:col-span-3 space-y-4">

        <!-- Content Type Tabs -->
        <Card>
          <CardHeader class="pb-3">
            <CardTitle class="flex items-center gap-2 text-base">
              <span>📝</span> Inhalt
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <!-- Type Selector -->
            <div class="grid grid-cols-4 sm:grid-cols-7 gap-1.5">
              <button
                v-for="type in contentTypes"
                :key="type.id"
                @click="currentType = type.id"
                class="flex flex-col items-center gap-1 p-2 rounded-lg border text-xs font-medium transition-all"
                :class="currentType === type.id
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'bg-card hover:bg-muted border-transparent'"
              >
                <span class="text-base">{{ type.icon }}</span>
                <span>{{ type.label }}</span>
              </button>
            </div>

            <Separator />

            <!-- Dynamic Form based on content type -->

            <!-- URL -->
            <div v-if="currentType === 'url'" class="space-y-3">
              <div class="space-y-1.5">
                <Label for="url-input">URL / Link</Label>
                <Input id="url-input" v-model="content.url" placeholder="https://beispiel.de" type="url" />
              </div>
            </div>

            <!-- Text -->
            <div v-else-if="currentType === 'text'" class="space-y-3">
              <div class="space-y-1.5">
                <Label>Text</Label>
                <textarea
                  v-model="content.text"
                  placeholder="Beliebiger Text..."
                  rows="4"
                  class="w-full px-3 py-2 rounded-md border bg-background text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
            </div>

            <!-- Email -->
            <div v-else-if="currentType === 'email'" class="space-y-3">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div class="space-y-1.5">
                  <Label>E-Mail Adresse</Label>
                  <Input v-model="content.email.to" placeholder="name@beispiel.de" type="email" />
                </div>
                <div class="space-y-1.5">
                  <Label>Betreff (optional)</Label>
                  <Input v-model="content.email.subject" placeholder="Betreff" />
                </div>
              </div>
              <div class="space-y-1.5">
                <Label>Nachricht (optional)</Label>
                <textarea
                  v-model="content.email.body"
                  placeholder="Nachricht..."
                  rows="3"
                  class="w-full px-3 py-2 rounded-md border bg-background text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
            </div>

            <!-- Phone -->
            <div v-else-if="currentType === 'phone'" class="space-y-3">
              <div class="space-y-1.5">
                <Label>Telefonnummer</Label>
                <Input v-model="content.phone" placeholder="+49 123 456789" type="tel" />
              </div>
            </div>

            <!-- SMS -->
            <div v-else-if="currentType === 'sms'" class="space-y-3">
              <div class="space-y-1.5">
                <Label>Telefonnummer</Label>
                <Input v-model="content.sms.number" placeholder="+49 123 456789" type="tel" />
              </div>
              <div class="space-y-1.5">
                <Label>Nachricht (optional)</Label>
                <textarea
                  v-model="content.sms.message"
                  placeholder="SMS Nachricht..."
                  rows="3"
                  class="w-full px-3 py-2 rounded-md border bg-background text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
            </div>

            <!-- WiFi -->
            <div v-else-if="currentType === 'wifi'" class="space-y-3">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div class="space-y-1.5">
                  <Label>Netzwerkname (SSID)</Label>
                  <Input v-model="content.wifi.ssid" placeholder="Mein WLAN" />
                </div>
                <div class="space-y-1.5">
                  <Label>Sicherheit</Label>
                  <Select v-model:model-value="content.wifi.security">
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="WPA">WPA / WPA2</SelectItem>
                      <SelectItem value="WEP">WEP</SelectItem>
                      <SelectItem value="nopass">Offen (kein Passwort)</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div v-if="content.wifi.security !== 'nopass'" class="space-y-1.5">
                <Label>Passwort</Label>
                <Input v-model="content.wifi.password" placeholder="WLAN Passwort" type="password" />
              </div>
              <div class="flex items-center space-x-2">
                <Checkbox id="wifi-hidden" v-model:checked="content.wifi.hidden" />
                <Label for="wifi-hidden" class="cursor-pointer font-normal">Verstecktes Netzwerk</Label>
              </div>
            </div>

            <!-- vCard -->
            <div v-else-if="currentType === 'vcard'" class="space-y-3">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div class="space-y-1.5">
                  <Label>Vorname</Label>
                  <Input v-model="content.vcard.firstName" placeholder="Max" />
                </div>
                <div class="space-y-1.5">
                  <Label>Nachname</Label>
                  <Input v-model="content.vcard.lastName" placeholder="Mustermann" />
                </div>
                <div class="space-y-1.5">
                  <Label>Telefon</Label>
                  <Input v-model="content.vcard.phone" placeholder="+49 123 456789" type="tel" />
                </div>
                <div class="space-y-1.5">
                  <Label>E-Mail</Label>
                  <Input v-model="content.vcard.email" placeholder="max@beispiel.de" type="email" />
                </div>
                <div class="space-y-1.5">
                  <Label>Unternehmen</Label>
                  <Input v-model="content.vcard.org" placeholder="Firma GmbH" />
                </div>
                <div class="space-y-1.5">
                  <Label>Webseite</Label>
                  <Input v-model="content.vcard.url" placeholder="https://beispiel.de" type="url" />
                </div>
              </div>
            </div>

            <!-- Geo -->
            <div v-else-if="currentType === 'geo'" class="space-y-3">
              <!-- Address Search -->
              <div class="flex gap-2">
                <Input
                  v-model="geoSearch"
                  placeholder="Adresse oder Ort suchen..."
                  @keyup.enter="searchAddress"
                  class="flex-1"
                />
                <Button @click="searchAddress" variant="outline" size="sm" :disabled="geoSearching">
                  {{ geoSearching ? '⏳' : '🔍' }}
                </Button>
                <Button @click="useMyLocation" variant="outline" size="sm" title="Meinen Standort verwenden">
                  📍
                </Button>
              </div>

              <!-- Interactive Map -->
              <div
                ref="mapContainer"
                class="h-72 rounded-lg border overflow-hidden"
                style="isolation: isolate;"
              />

              <!-- Coordinate Inputs (synced with map) -->
              <div class="grid grid-cols-2 gap-3">
                <div class="space-y-1.5">
                  <Label>Breitengrad</Label>
                  <Input
                    v-model="content.geo.lat"
                    placeholder="52.5200"
                    type="number"
                    step="any"
                    @change="updateMarkerFromInputs"
                  />
                </div>
                <div class="space-y-1.5">
                  <Label>Längengrad</Label>
                  <Input
                    v-model="content.geo.lng"
                    placeholder="13.4050"
                    type="number"
                    step="any"
                    @change="updateMarkerFromInputs"
                  />
                </div>
              </div>
              <p class="text-xs text-muted-foreground">💡 Auf die Karte klicken um einen Punkt zu setzen</p>
            </div>

            <!-- Generated Content Preview -->
            <div class="mt-2 p-2 rounded bg-muted text-xs font-mono text-muted-foreground break-all">
              {{ generatedContent || '(kein Inhalt)' }}
            </div>
          </CardContent>
        </Card>

        <!-- Options Card -->
        <Card>
          <CardHeader class="pb-3">
            <CardTitle class="flex items-center gap-2 text-base">
              <span>🎨</span> Anpassung
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-5">

            <!-- Size -->
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <Label class="flex items-center gap-1">
                  Größe
                  <Tooltip>
                    <TooltipTrigger as-child><span class="text-muted-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger>
                    <TooltipContent><p class="text-sm">Größe der Vorschau und des PNG-Downloads in Pixel.</p></TooltipContent>
                  </Tooltip>
                </Label>
                <span class="text-sm font-mono text-muted-foreground">{{ options.size }}px</span>
              </div>
              <input
                type="range"
                v-model.number="options.size"
                min="100"
                max="1000"
                step="50"
                class="w-full h-2 bg-muted rounded-full appearance-none cursor-pointer accent-primary"
              />
              <div class="flex justify-between text-xs text-muted-foreground">
                <span>100px</span>
                <span>1000px</span>
              </div>
            </div>

            <!-- Error Correction -->
            <div class="space-y-1.5">
              <Label class="flex items-center gap-1">
                Fehlerkorrektur
                <Tooltip>
                  <TooltipTrigger as-child><span class="text-muted-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger>
                  <TooltipContent class="max-w-xs">
                    <p class="text-sm">Je höher die Stufe, desto mehr vom QR Code kann beschädigt/überdeckt sein und trotzdem gelesen werden. H (30%) empfohlen bei Logo.</p>
                  </TooltipContent>
                </Tooltip>
              </Label>
              <div class="grid grid-cols-4 gap-2">
                <button
                  v-for="lvl in errorLevels"
                  :key="lvl.value"
                  @click="options.level = lvl.value"
                  class="flex flex-col items-center p-2 rounded-lg border text-xs transition-all"
                  :class="options.level === lvl.value ? 'bg-primary text-primary-foreground border-primary' : 'bg-card hover:bg-muted border-border'"
                >
                  <span class="font-bold">{{ lvl.value }}</span>
                  <span class="text-[10px] opacity-75">{{ lvl.pct }}</span>
                  <span class="text-[10px] opacity-60">{{ lvl.label }}</span>
                </button>
              </div>
            </div>

            <!-- Colors -->
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <Label>Vordergrundfarbe</Label>
                <div class="flex items-center gap-2 p-2 border rounded-md bg-card">
                  <input
                    type="color"
                    v-model="options.fgColor"
                    class="w-8 h-8 rounded cursor-pointer border-0 bg-transparent"
                    :disabled="options.transparent && false"
                  />
                  <span class="text-sm font-mono">{{ options.fgColor }}</span>
                </div>
              </div>
              <div class="space-y-1.5">
                <Label>Hintergrundfarbe</Label>
                <div class="flex items-center gap-2 p-2 border rounded-md" :class="options.transparent ? 'opacity-50' : 'bg-card'">
                  <input
                    type="color"
                    v-model="options.bgColor"
                    class="w-8 h-8 rounded cursor-pointer border-0 bg-transparent"
                    :disabled="options.transparent"
                  />
                  <span class="text-sm font-mono">{{ options.bgColor }}</span>
                </div>
              </div>
            </div>

            <!-- Transparent Background -->
            <div class="flex items-center space-x-2">
              <Checkbox id="transparent" v-model:checked="options.transparent" />
              <Label for="transparent" class="cursor-pointer font-normal flex items-center gap-1">
                Transparenter Hintergrund
                <Tooltip>
                  <TooltipTrigger as-child><span class="text-muted-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger>
                  <TooltipContent><p class="text-sm">Kein Hintergrund (nur PNG unterstützt Transparenz).</p></TooltipContent>
                </Tooltip>
              </Label>
            </div>

            <Separator />

            <!-- Logo Upload -->
            <div class="space-y-2">
              <Label class="flex items-center gap-1">
                Logo Overlay (optional)
                <Tooltip>
                  <TooltipTrigger as-child><span class="text-muted-foreground cursor-help text-xs">ⓘ</span></TooltipTrigger>
                  <TooltipContent class="max-w-xs"><p class="text-sm">Logo wird zentriert auf den QR Code gelegt. Fehlerkorrektur H empfohlen, damit der Code noch lesbar bleibt.</p></TooltipContent>
                </Tooltip>
              </Label>
              <div class="flex items-center gap-3">
                <label class="flex-1 flex items-center justify-center gap-2 p-3 border-2 border-dashed rounded-lg cursor-pointer hover:bg-muted/50 transition-colors text-sm text-muted-foreground">
                  <span>{{ logoName || '📎 Bild auswählen (PNG, JPG, SVG)' }}</span>
                  <input type="file" accept="image/*" @change="onLogoUpload" class="hidden" />
                </label>
                <Button v-if="logoDataUrl" @click="clearLogo" variant="outline" size="sm">✕</Button>
              </div>
              <div v-if="logoDataUrl" class="flex items-center gap-2">
                <Label class="text-xs">Logo Größe:</Label>
                <input
                  type="range"
                  v-model.number="options.logoSize"
                  min="10"
                  max="40"
                  step="1"
                  class="flex-1 h-1.5 bg-muted rounded-full appearance-none cursor-pointer accent-primary"
                />
                <span class="text-xs font-mono text-muted-foreground">{{ options.logoSize }}%</span>
              </div>
            </div>

          </CardContent>
        </Card>
      </div>

      <!-- RIGHT: Preview + Downloads -->
      <div class="lg:col-span-2 space-y-4">
        <Card class="sticky top-20">
          <CardHeader class="pb-3">
            <CardTitle class="flex items-center gap-2 text-base">
              <span>👁️</span> Vorschau
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">

            <!-- Canvas Preview -->
            <div class="flex justify-center items-center bg-muted/30 rounded-xl p-4 min-h-48">
              <div v-if="!generatedContent" class="text-center text-muted-foreground space-y-2">
                <div class="text-5xl">📱</div>
                <p class="text-sm">Inhalt eingeben um QR Code zu generieren</p>
              </div>
              <canvas
                v-show="generatedContent"
                ref="canvasRef"
                class="rounded-lg max-w-full"
                :style="{ width: Math.min(options.size, 320) + 'px', height: Math.min(options.size, 320) + 'px' }"
              />
            </div>

            <!-- Download Buttons -->
            <div v-if="generatedContent" class="space-y-2">
              <p class="text-xs text-muted-foreground text-center">Export</p>
              <div class="grid grid-cols-2 gap-2">
                <Button @click="downloadPNG" variant="default" class="bg-blue-600 hover:bg-blue-700 w-full">
                  📥 PNG
                </Button>
                <Button @click="downloadSVG" variant="outline" class="w-full">
                  📄 SVG
                </Button>
              </div>
              <Button @click="copyToClipboard" variant="outline" class="w-full">
                {{ copied ? '✅ Kopiert!' : '📋 In Zwischenablage' }}
              </Button>
              <Button @click="downloadHD" variant="outline" class="w-full text-purple-600 border-purple-300 hover:bg-purple-50 dark:hover:bg-purple-950">
                <span v-if="hdLoading">⏳ Generiere...</span>
                <span v-else>⚡ HD Download (2000px via Server)</span>
              </Button>
            </div>

            <!-- Stats -->
            <div v-if="generatedContent" class="grid grid-cols-2 gap-2 text-xs text-muted-foreground">
              <div class="bg-muted/30 rounded p-2 text-center">
                <div class="font-semibold text-foreground">{{ generatedContent.length }}</div>
                <div>Zeichen</div>
              </div>
              <div class="bg-muted/30 rounded p-2 text-center">
                <div class="font-semibold text-foreground">{{ options.level }}</div>
                <div>Fehlerkorrektur</div>
              </div>
            </div>

          </CardContent>
        </Card>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick, onUnmounted } from 'vue'
import QRCode from 'qrcode'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

const API_BASE = '/api/qrcode'

// ── Content Types ────────────────────────────────────────────────────────────
const contentTypes = [
  { id: 'url', icon: '🔗', label: 'URL' },
  { id: 'text', icon: '📝', label: 'Text' },
  { id: 'email', icon: '✉️', label: 'E-Mail' },
  { id: 'phone', icon: '📞', label: 'Telefon' },
  { id: 'sms', icon: '💬', label: 'SMS' },
  { id: 'wifi', icon: '📶', label: 'WiFi' },
  { id: 'vcard', icon: '👤', label: 'vCard' },
  { id: 'geo', icon: '📍', label: 'Standort' },
]

const errorLevels = [
  { value: 'L', pct: '7%', label: 'Niedrig' },
  { value: 'M', pct: '15%', label: 'Mittel' },
  { value: 'Q', pct: '25%', label: 'Hoch' },
  { value: 'H', pct: '30%', label: 'Max' },
]

// ── State ─────────────────────────────────────────────────────────────────────
const currentType = ref('url')
const canvasRef = ref(null)
const copied = ref(false)
const hdLoading = ref(false)
const logoDataUrl = ref(null)
const logoName = ref('')

// ── Map (Geo) ─────────────────────────────────────────────────────────────────
const mapContainer = ref(null)
const geoSearch = ref('')
const geoSearching = ref(false)
let leafletMap = null
let leafletMarker = null

const mapIcon = L.divIcon({
  html: '<div style="font-size:28px;line-height:1;filter:drop-shadow(0 2px 3px rgba(0,0,0,.4))">📍</div>',
  iconSize: [28, 28],
  iconAnchor: [14, 28],
  className: '',
})

watch(currentType, async (type) => {
  if (type === 'geo') {
    await nextTick()
    initMap()
  } else {
    destroyMap()
  }
})

onUnmounted(destroyMap)

function initMap() {
  if (!mapContainer.value || leafletMap) return

  const lat = parseFloat(content.geo.lat) || 51.1657
  const lng = parseFloat(content.geo.lng) || 10.4515
  const zoom = content.geo.lat ? 13 : 5

  leafletMap = L.map(mapContainer.value, { zoomControl: true }).setView([lat, lng], zoom)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="https://openstreetmap.org">OpenStreetMap</a>',
    maxZoom: 19,
  }).addTo(leafletMap)

  leafletMap.on('click', (e) => {
    const { lat: clickLat, lng: clickLng } = e.latlng
    content.geo.lat = clickLat.toFixed(6)
    content.geo.lng = clickLng.toFixed(6)
    placeMarker(clickLat, clickLng)
  })

  if (content.geo.lat && content.geo.lng) {
    placeMarker(lat, lng)
  }
}

function destroyMap() {
  if (leafletMap) {
    leafletMap.remove()
    leafletMap = null
    leafletMarker = null
  }
}

function placeMarker(lat, lng) {
  if (!leafletMap) return
  if (leafletMarker) {
    leafletMarker.setLatLng([lat, lng])
  } else {
    leafletMarker = L.marker([lat, lng], { icon: mapIcon }).addTo(leafletMap)
  }
}

function updateMarkerFromInputs() {
  const lat = parseFloat(content.geo.lat)
  const lng = parseFloat(content.geo.lng)
  if (!isNaN(lat) && !isNaN(lng) && leafletMap) {
    placeMarker(lat, lng)
    leafletMap.setView([lat, lng], leafletMap.getZoom() < 10 ? 13 : leafletMap.getZoom())
  }
}

async function searchAddress() {
  if (!geoSearch.value.trim()) return
  geoSearching.value = true
  try {
    const res = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(geoSearch.value)}&limit=1`,
      { headers: { 'Accept-Language': 'de' } }
    )
    const data = await res.json()
    if (data.length > 0) {
      const { lat, lon, display_name } = data[0]
      content.geo.lat = parseFloat(lat).toFixed(6)
      content.geo.lng = parseFloat(lon).toFixed(6)
      if (leafletMap) {
        leafletMap.setView([lat, lon], 15)
        placeMarker(parseFloat(lat), parseFloat(lon))
      }
      geoSearch.value = display_name.split(',').slice(0, 3).join(',').trim()
      toast.success('📍 Ort gefunden')
    } else {
      toast.error('❌ Kein Ergebnis', { description: 'Adresse nicht gefunden.' })
    }
  } catch {
    toast.error('❌ Suche fehlgeschlagen')
  } finally {
    geoSearching.value = false
  }
}

function useMyLocation() {
  if (!navigator.geolocation) {
    toast.error('❌ Geolocation nicht unterstützt')
    return
  }
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      const lat = pos.coords.latitude
      const lng = pos.coords.longitude
      content.geo.lat = lat.toFixed(6)
      content.geo.lng = lng.toFixed(6)
      if (leafletMap) {
        leafletMap.setView([lat, lng], 15)
        placeMarker(lat, lng)
      }
      toast.success('📍 Standort erkannt')
    },
    () => toast.error('❌ Standort nicht verfügbar')
  )
}

const content = reactive({
  url: '',
  text: '',
  email: { to: '', subject: '', body: '' },
  phone: '',
  sms: { number: '', message: '' },
  wifi: { ssid: '', password: '', security: 'WPA', hidden: false },
  vcard: { firstName: '', lastName: '', phone: '', email: '', org: '', url: '' },
  geo: { lat: '', lng: '' },
})

const options = reactive({
  size: 300,
  level: 'M',
  fgColor: '#000000',
  bgColor: '#ffffff',
  transparent: false,
  logoSize: 20,
})

// ── Content Builder ───────────────────────────────────────────────────────────
const generatedContent = computed(() => {
  switch (currentType.value) {
    case 'url':
      return content.url.trim()

    case 'text':
      return content.text.trim()

    case 'email': {
      if (!content.email.to) return ''
      let mailto = `mailto:${content.email.to}`
      const params = []
      if (content.email.subject) params.push(`subject=${encodeURIComponent(content.email.subject)}`)
      if (content.email.body) params.push(`body=${encodeURIComponent(content.email.body)}`)
      if (params.length) mailto += '?' + params.join('&')
      return mailto
    }

    case 'phone':
      return content.phone ? `tel:${content.phone}` : ''

    case 'sms':
      if (!content.sms.number) return ''
      return content.sms.message
        ? `SMSTO:${content.sms.number}:${content.sms.message}`
        : `sms:${content.sms.number}`

    case 'wifi': {
      if (!content.wifi.ssid) return ''
      const sec = content.wifi.security
      const pass = content.wifi.password
      const hidden = content.wifi.hidden ? 'true' : 'false'
      return `WIFI:T:${sec};S:${content.wifi.ssid};P:${pass};H:${hidden};;`
    }

    case 'vcard': {
      const v = content.vcard
      if (!v.firstName && !v.lastName) return ''
      const lines = [
        'BEGIN:VCARD',
        'VERSION:3.0',
        `FN:${v.firstName} ${v.lastName}`.trim(),
        `N:${v.lastName};${v.firstName};;;`,
      ]
      if (v.phone) lines.push(`TEL:${v.phone}`)
      if (v.email) lines.push(`EMAIL:${v.email}`)
      if (v.org) lines.push(`ORG:${v.org}`)
      if (v.url) lines.push(`URL:${v.url}`)
      lines.push('END:VCARD')
      return lines.join('\n')
    }

    case 'geo':
      return content.geo.lat && content.geo.lng
        ? `geo:${content.geo.lat},${content.geo.lng}`
        : ''

    default:
      return ''
  }
})

// ── QR Generation ─────────────────────────────────────────────────────────────
let generateTimer = null

watch(
  [generatedContent, () => options.size, () => options.level, () => options.fgColor, () => options.bgColor, () => options.transparent, logoDataUrl, () => options.logoSize],
  () => {
    clearTimeout(generateTimer)
    generateTimer = setTimeout(renderQR, 150)
  },
  { immediate: true }
)

async function renderQR() {
  if (!generatedContent.value) return
  await nextTick()
  if (!canvasRef.value) return

  const qrOpts = {
    width: Math.min(options.size, 320),
    errorCorrectionLevel: options.level,
    color: {
      dark: options.fgColor,
      light: options.transparent ? '#ffffff00' : options.bgColor,
    },
    margin: 2,
  }

  try {
    await QRCode.toCanvas(canvasRef.value, generatedContent.value, qrOpts)

    // Logo overlay
    if (logoDataUrl.value) {
      await drawLogo(canvasRef.value)
    }
  } catch (e) {
    console.error('QR render error:', e)
  }
}

function drawLogo(canvas) {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      const ctx = canvas.getContext('2d')
      const logoW = canvas.width * (options.logoSize / 100)
      const logoH = (img.height / img.width) * logoW
      const x = (canvas.width - logoW) / 2
      const y = (canvas.height - logoH) / 2

      // White background behind logo
      ctx.fillStyle = options.transparent ? 'rgba(255,255,255,0.85)' : options.bgColor
      const pad = 6
      ctx.fillRect(x - pad, y - pad, logoW + pad * 2, logoH + pad * 2)
      ctx.drawImage(img, x, y, logoW, logoH)
      resolve()
    }
    img.src = logoDataUrl.value
  })
}

// ── Logo Upload ───────────────────────────────────────────────────────────────
function onLogoUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return
  logoName.value = file.name
  const reader = new FileReader()
  reader.onload = (e) => {
    logoDataUrl.value = e.target.result
    // Recommend H error correction for logo
    if (options.level !== 'H') {
      options.level = 'H'
      toast.info('💡 Fehlerkorrektur auf H gesetzt', { description: 'Empfohlen bei Logo-Overlay für bessere Lesbarkeit.' })
    }
  }
  reader.readAsDataURL(file)
}

function clearLogo() {
  logoDataUrl.value = null
  logoName.value = ''
}

// ── Export ────────────────────────────────────────────────────────────────────
async function downloadPNG() {
  if (!canvasRef.value || !generatedContent.value) return

  // Render at full requested size
  const offscreen = document.createElement('canvas')
  const qrOpts = {
    width: options.size,
    errorCorrectionLevel: options.level,
    color: {
      dark: options.fgColor,
      light: options.transparent ? '#ffffff00' : options.bgColor,
    },
    margin: 2,
  }

  await QRCode.toCanvas(offscreen, generatedContent.value, qrOpts)
  if (logoDataUrl.value) {
    await drawLogoOnCanvas(offscreen)
  }

  offscreen.toBlob((blob) => {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `qrcode-${options.size}px.png`
    a.click()
    URL.revokeObjectURL(url)
    toast.success('📥 PNG heruntergeladen')
  }, 'image/png')
}

function drawLogoOnCanvas(canvas) {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      const ctx = canvas.getContext('2d')
      const logoW = canvas.width * (options.logoSize / 100)
      const logoH = (img.height / img.width) * logoW
      const x = (canvas.width - logoW) / 2
      const y = (canvas.height - logoH) / 2
      const pad = Math.round(canvas.width * 0.02)
      ctx.fillStyle = options.transparent ? 'rgba(255,255,255,0.85)' : options.bgColor
      ctx.fillRect(x - pad, y - pad, logoW + pad * 2, logoH + pad * 2)
      ctx.drawImage(img, x, y, logoW, logoH)
      resolve()
    }
    img.src = logoDataUrl.value
  })
}

async function downloadSVG() {
  if (!generatedContent.value) return
  try {
    const svg = await QRCode.toString(generatedContent.value, {
      type: 'svg',
      errorCorrectionLevel: options.level,
      color: { dark: options.fgColor, light: options.bgColor },
      margin: 2,
    })
    const blob = new Blob([svg], { type: 'image/svg+xml' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'qrcode.svg'
    a.click()
    URL.revokeObjectURL(url)
    toast.success('📄 SVG heruntergeladen')
  } catch (e) {
    toast.error('❌ SVG Fehler', { description: e.message })
  }
}

async function copyToClipboard() {
  if (!canvasRef.value) return
  try {
    const blob = await new Promise((resolve) => canvasRef.value.toBlob(resolve, 'image/png'))
    await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })])
    copied.value = true
    toast.success('📋 In Zwischenablage kopiert')
    setTimeout(() => { copied.value = false }, 2000)
  } catch (e) {
    toast.error('❌ Kopieren fehlgeschlagen', { description: e.message })
  }
}

async function downloadHD() {
  if (!generatedContent.value) return
  hdLoading.value = true
  try {
    const response = await axios.post(`${API_BASE}/generate`, {
      content: generatedContent.value,
      size: 2000,
      level: options.level,
      fg_color: options.fgColor,
      bg_color: options.transparent ? '#ffffff' : options.bgColor,
    }, { responseType: 'arraybuffer' })

    const blob = new Blob([response.data], { type: 'image/png' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'qrcode-2000px-HD.png'
    a.click()
    URL.revokeObjectURL(url)
    toast.success('⚡ HD PNG heruntergeladen (2000px)')
  } catch (e) {
    toast.error('❌ HD Download fehlgeschlagen', { description: e.message })
  } finally {
    hdLoading.value = false
  }
}
</script>