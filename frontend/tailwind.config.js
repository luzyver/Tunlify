/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        bg: '#ffffff',
        'bg-alt': '#fafafa',
        surface: '#f4f4f5',
        text: '#000000',
        'text-muted': '#52525b',
        'text-dim': '#a1a1aa',
        border: '#e4e4e7',
        'border-strong': '#d4d4d8',
        accent: '#ff5d1f',
        'accent-hover': '#e64a0e',
        'accent-soft': '#ffead9',
      },
      fontFamily: {
        display: ['Tiempos Headline', 'Georgia', 'serif'],
        sans: ['Geist', 'system-ui', 'sans-serif'],
        mono: ['Geist Mono', 'monospace'],
      },
      fontSize: {
        'xs': '12px',
        'sm': '14px',
        'base': '16px',
        'lg': '18px',
        'xl': '20px',
        '2xl': '24px',
        '3xl': '32px',
        '4xl': '48px',
        '5xl': '64px',
        '6xl': '88px',
      },
      spacing: {
        '1': '4px',
        '2': '8px',
        '3': '12px',
        '4': '16px',
        '6': '24px',
        '8': '32px',
        '12': '48px',
        '16': '64px',
        '24': '96px',
      },
      borderRadius: {
        DEFAULT: '6px',
        md: '6px',
        lg: '8px',
      },
      boxShadow: {
        popover: '0 2px 8px rgba(0,0,0,0.04)',
      },
      maxWidth: {
        reading: '720px',
        content: '1180px',
      },
    },
  },
  plugins: [],
}
