module.exports = {
    title: 'Gilfoyle',
    tagline: 'Cloud-native solution to embed media streaming in any application at any scale.',
    url: 'https://docs.gilfoyle.dreamvo.com',
    baseUrl: '/',
    onBrokenLinks: 'throw',
    onBrokenMarkdownLinks: 'warn',
    favicon: 'img/favicon.ico',
    organizationName: 'dreamvo', // Usually your GitHub org/user name.
    projectName: 'gilfoyle', // Usually your repo name.
    themeConfig: {
        navbar: {
            title: 'Gilfoyle',
            // logo: {
            //     alt: 'Gilfoyle logo',
            //     src: 'img/logo.svg',
            // },
            items: [
                {
                    to: 'docs/',
                    activeBasePath: 'docs',
                    label: 'Get started',
                    position: 'left',
                },
                {
                    to: 'docs/',
                    activeBasePath: 'docs',
                    label: 'Use cases',
                    position: 'left',
                },
                {
                    to: 'docs/',
                    activeBasePath: 'docs',
                    label: 'Concepts',
                    position: 'left',
                },
                {
                    to: 'docs/',
                    activeBasePath: 'docs',
                    label: 'Migration guides',
                    position: 'left',
                },
                {
                    to: 'docs/',
                    activeBasePath: 'docs',
                    label: 'References',
                    position: 'left',
                },
                {
                    to: 'https://github.com/dreamvo/gilfoyle/releases',
                    activeBasePath: 'docs',
                    label: 'Releases',
                    position: 'right',
                },
                {
                    to: 'blog', label: 'Blog', position: 'right'
                },
                {
                    href: 'https://github.com/dreamvo/gilfoyle',
                    label: 'GitHub',
                    position: 'right',
                },
            ],
        },
        footer: {
            style: 'dark',
            links: [
                {
                    title: 'Learn',
                    items: [
                        {
                            label: 'Getting started',
                            to: 'docs/',
                        },
                        {
                            label: 'Deploy on Kubernetes',
                            to: 'docs/doc2/',
                        },
                        {
                            label: 'Manage a worker pool',
                            to: 'docs/doc2/',
                        },
                        {
                            label: 'OpenAPI specs',
                            to: 'https://petstore.swagger.io/?url=https://raw.githubusercontent.com/dreamvo/gilfoyle/master/api/docs/swagger.json',
                        },
                    ],
                },
                {
                    title: 'Community',
                    items: [
                        {
                            label: 'Discuss on GitHub',
                            href: 'https://stackoverflow.com/questions/tagged/docusaurus',
                        },
                        {
                            label: 'Report a bug',
                            href: 'https://discordapp.com/invite/docusaurus',
                        },
                        {
                            label: 'Follow us on Twitter',
                            href: 'https://twitter.com/dreamvoapp',
                        },
                    ],
                },
                {
                    title: 'More',
                    items: [
                        {
                            label: 'Blog',
                            to: 'blog',
                        },
                        {
                            label: 'GitHub',
                            href: 'https://github.com/dreamvo/gilfoyle',
                        },
                    ],
                },
            ],
            copyright: `Copyright © ${new Date().getFullYear()} Dreamvo, Inc.`,
        },
    },
    presets: [
        [
            '@docusaurus/preset-classic',
            {
                docs: {
                    sidebarPath: require.resolve('./sidebars.js'),
                    // Please change this to your repo.
                    editUrl:
                        'https://github.com/dreamvo/gilfoyle/edit/master/website/',
                },
                blog: {
                    showReadingTime: true,
                    // Please change this to your repo.
                    editUrl:
                        'https://github.com/dreamvo/gilfoyle/edit/master/website/blog/',
                },
                theme: {
                    customCss: require.resolve('./src/css/custom.css'),
                },
            },
        ],
    ],
};
