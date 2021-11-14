const production = !process.env.ROLLUP_WATCH;
module.exports = {
  mode: 'jit',
  future: {
    purgeLayersByDefault: true,
    removeDeprecatedGapUtilities: true,
  },
  plugins: [
    require('daisyui'),
  ],
  // https://daisyui.com/docs/config
  daisyui: {
    styled: true,
    themes: true,
    themes: [
      'wireframe',
      'light', // first one will be the default theme
      'dark',
      'forest',
      'retro',
      'cyberpunk',
      'black',
      'dracula',
    ],
    base: true,
    utils: true,
    logs: true,
    rtl: false,
  },
  purge: {
    content: [
      "./public/index.html",
      "./src/**/*.svelte",
    ],
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  enabled: production,
}
