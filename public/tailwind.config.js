/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      'views/**/*.templ',
    ],
    darkMode: 'class',
    theme: {
      extend: {
        fontFamily: {
          mono: ['Courier Prime', 'monospace'],
        }
      },
    },
    corePlugins: {
      preflight: true,
    }
  }