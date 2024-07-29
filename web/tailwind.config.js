/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./static/public/assets/css/**/*.css",
    "./static/public/templates/**/*.html",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('flowbite/plugin')
  ],
};
