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
    defaultTheme: 'tunlify',
    themes: {
      tunlify: {
        dark: false,
        colors: {
          primary: '#5266eb',
          secondary: '#5a5548',
          accent: '#5266eb',
          error: '#b54a3a',
          info: '#5266eb',
          success: '#2f7d57',
          warning: '#c98a42',
          background: '#f6f5f2',
          surface: '#ffffff',
          'surface-variant': '#fbfaf6',
          'on-surface': '#2a2924',
          'on-surface-variant': '#5a5548',
        },
      },
    },
  },
  defaults: {
    VCard: {
      variant: 'outlined',
      flat: true,
    },
    VTextField: {
      variant: 'outlined',
      density: 'compact',
      hideDetails: 'auto',
    },
    VBtn: {
      variant: 'tonal',
      density: 'compact',
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
  },
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(vuetify)
app.mount('#app')
