/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js}", "./dist/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        "black": "#242323",
        "gray": "#bdb4b4",
        "blue": "#257dcf",
      },
    },
  },
  plugins: [],
}

