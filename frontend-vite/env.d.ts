/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'markdown-it-task-lists'
declare module 'markdown-it-abbr'
declare module 'markdown-it-github-alerts'
declare module 'markdown-it-toc-done-right'
declare module '@traptitech/markdown-it-katex'
