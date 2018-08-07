module.exports = {
  extends: ['standard', 'prettier', 'prettier/standard'],
  plugins: ['prettier', 'standard'],
  env: {
    node: true,
    es6: true,
  },
  rules: {
    'prettier/prettier': [
      'error',
      {
        printWidth: 80,
        tabWidth: 2,
        useTabs: false,
        semi: false,
        singleQuote: true,
        trailingComma: 'es5',
        bracketSpacing: true,
      },
    ],
  },
}
