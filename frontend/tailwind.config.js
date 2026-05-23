
module.exports = {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  darkMode: 'class',
  theme: {
    fontSize: {
      '2xs': ['11px', '14px'],
      xs: ['13px', '18px'],
      sm: ['14px', '22px'],
      base: ['14px', '22px'],
      md: ['16px', '26px'],
      lg: ['18px', '26px'],
      xl: ['22px', '30px'],
      '2xl': ['28px', '34px'],
      '3xl': ['36px', '40px'],
      '4xl': ['48px', '52px'],
      '5xl': ['72px', '76px'],
    },
    extend: {
      colors: {
        bg: '#f6f5f2',
        'bg-alt': '#ebe8e0',
        surface: '#ffffff',
        'surface-alt': '#fbfaf6',

        text: '#2a2924',
        'text-muted': '#5a5548',
        'text-dim': '#8a8478',
        border: '#ded9ca',
        'border-strong': '#c9c3b3',

        accent: '#5266eb',
        'accent-hover': '#4255d4',
        'accent-soft': '#e2e6fb',

        success: '#2f7d57',
        warning: '#c98a42',
        danger: '#b54a3a',

        'bg-dark': '#15140f',
        'bg-dark-alt': '#1d1c16',
        'surface-dark': '#25241d',
        'text-dark': '#f6f5f2',
        'text-dark-muted': '#a8a394',
        'border-dark': '#34322a',
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'sans-serif'],

        mono: ['"Berkeley Mono"', '"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace'],

        display: ['"Arcadia Display"', '"Tiempos Headline"', 'Newsreader', '"Source Serif 4"', 'Georgia', 'serif'],
      },

      spacing: {
        '4.5': '18px',
        '13': '52px',
        '15': '60px',
      },
      borderRadius: {
        none: '0',
        sm: '4px',
        DEFAULT: '6px',
        md: '6px',
        lg: '8px',
        xl: '12px',
        full: '9999px',
      },
      boxShadow: {
        modal: '0 12px 32px rgba(42,41,36,0.10)',
        popover: '0 4px 12px rgba(42,41,36,0.06)',
        none: 'none',
      },
      maxWidth: {
        app: '1280px',
        marketing: '1080px',
        reading: '720px',
      },

      width: {
        sidebar: '240px',
      },
      ringWidth: {
        DEFAULT: '2px',
      },
      ringOffsetWidth: {
        DEFAULT: '2px',
      },
      transitionTimingFunction: {
        DEFAULT: 'cubic-bezier(0.2, 0, 0, 1)',
      },
    },
  },
  plugins: [],
}
