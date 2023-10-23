module.exports = {
    plugins: [require('prettier-plugin-svelte')],
    overrides: [{ files: '*.svelte', options: { parser: 'svelte' } }],
};
