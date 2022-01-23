module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['prettier', 'simple-import-sort', '@typescript-eslint'],
  parserOptions: {
    sourceType: 'module',
    ecmaVersion: 2015
  },
  extends: ['eslint:recommended',  'plugin:@typescript-eslint/recommended'],
  ignorePatterns: ['package.json', 'yarn.lock'],
  rules: {
    'prettier/prettier': [
      'error',
      {
        parser: 'typescript',
        printWidth: 100,
        trailingComma: 'none',
        singleQuote: true
      }
    ],
    'no-restricted-syntax': [
      'error',
      {
        selector: "ImportDeclaration[source.value='react'] > :matches(ImportDefaultSpecifier)",
        message: "Please use `import * as React from 'react'` instead."
      }
    ],
    'simple-import-sort/imports': 'error',
    'no-console': 'warn',
    'no-duplicate-imports': 'off',
    'func-style': ['off', 'expression'],
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/consistent-type-assertions': 'warn',
    '@typescript-eslint/explicit-member-accessibility': 'off',
    '@typescript-eslint/explicit-module-boundary-types': 'off',
    '@typescript-eslint/no-duplicate-imports': ['error'],
    '@typescript-eslint/no-empty-function': 'off',
    '@typescript-eslint/no-use-before-define': 'off',
    '@typescript-eslint/no-unused-vars': ['error', { ignoreRestSiblings: true }],
    '@typescript-eslint/no-non-null-assertion': 'off',
    '@typescript-eslint/ban-ts-comment': 'off',
    'react-hooks/exhaustive-deps': 'off', // issues with hocs: https://github.com/facebook/react/pull/16712
    'no-restricted-imports': [2, { patterns: ['@smartly/*/*'] }],
    'react/prop-types': 'off'
  },
  settings: {
    react: {
      version: 'detect'
    }
  }
};
