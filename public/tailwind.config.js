/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      './views/**/*.templ',
    ],
    darkMode: 'class',
    corePlugins: {
      preflight: true,
    }
  }