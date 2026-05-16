/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        void: '#0a0a0f',
        surface: '#12121a',
        chrome: '#1c1c2e',
        border: '#2a2a3a',
        muted: '#6b7280',
        neon: '#00ff88',
        magenta: '#ff00ff',
        cyan: '#00d4ff',
        danger: '#ff3366',
      },
      fontFamily: {
        display: ['Orbitron', 'monospace'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      boxShadow: {
        neon: '0 0 5px #00ff88, 0 0 10px #00ff8840',
        'neon-lg': '0 0 10px #00ff88, 0 0 20px #00ff8860, 0 0 40px #00ff8830',
        'neon-magenta': '0 0 5px #ff00ff, 0 0 20px #ff00ff60',
        'neon-cyan': '0 0 5px #00d4ff, 0 0 20px #00d4ff60',
      },
      animation: {
        blink: 'blink 1s step-end infinite',
        glitch: 'glitch 3s infinite',
      },
      keyframes: {
        blink: { '50%': { opacity: '0' } },
        glitch: {
          '0%, 92%, 100%': { transform: 'translate(0)' },
          '93%': { transform: 'translate(-2px, 2px)' },
          '95%': { transform: 'translate(2px, -1px)' },
          '97%': { transform: 'translate(-1px, -1px)' },
        },
      },
    },
  },
  plugins: [],
}
