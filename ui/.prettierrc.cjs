module.exports = {
    pluginSearchDirs: false,
    plugins: [require('prettier-plugin-svelte')],
    overrides: [{ files: '*.svelte', options: { parser: 'svelte' } }],
};
