import { createApp } from 'vue'
import { createPinia } from 'pinia'

import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'
import './assets/main.css'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'bitcoinDefi',
    themes: {
      bitcoinDefi: {
        dark: true,
        colors: {
          background: '#030304',
          surface: '#0F1115',
          'surface-variant': '#0A0C10',
          primary: '#F7931A',
          secondary: '#EA580C',
          accent: '#FFD600',
          error: '#EF4444',
          info: '#F7931A',
          success: '#FFD600',
          warning: '#F7931A',
          'on-background': '#FFFFFF',
          'on-surface': '#FFFFFF',
          'on-surface-variant': '#94A3B8',
        },
      },
    },
  },
  defaults: {
    VCard: {
      variant: 'outlined',
      flat: true,
      color: 'surface',
    },
    VTextField: {
      variant: 'outlined',
      density: 'compact',
      hideDetails: 'auto',
      color: 'primary',
      baseColor: 'grey',
    },
    VBtn: {
      variant: 'tonal',
      density: 'compact',
      color: 'primary',
    },
    VSelect: {
      variant: 'outlined',
      density: 'compact',
      hideDetails: 'auto',
    },
    VTextarea: {
      variant: 'outlined',
      density: 'compact',
      hideDetails: 'auto',
    },
    VChip: {
      size: 'small',
    },
    VAlert: {
      variant: 'tonal',
      density: 'compact',
    },
    VNavigationDrawer: {
      color: 'background',
    },
  },
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(vuetify)
app.mount('#app')
