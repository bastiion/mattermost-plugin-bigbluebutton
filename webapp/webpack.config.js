module.exports = {
    entry: './index.js',
    output: {
        filename: 'dist/bigbluebutton_bundle.js'
    },
    externals: {
      react: 'React',
      redux: 'Redux',
      'react-redux': 'ReactRedux',
      'react-bootstrap': 'ReactBootstrap',

    },
    module: {
        loaders: [
            {
                test: /\.(js|jsx|tsx)?$/,
                loader: 'babel-loader',
                exclude: /(node_modules|non_npm_dependencies)/,
                query: {
                    presets: [
                        'react',
                        ['es2015', {modules: false}],
                        'stage-0'
                    ],
                    plugins: ['transform-runtime']
                }
            },
            {
                test: /\.[sac]ss$/i,
                use: [
                    // Creates `style` nodes from JS strings
                    'style-loader',
                    // Translates CSS into CommonJS
                    'css-loader',
                    // Compiles Sass to CSS
                    /*{
                        loader: 'sass-loader',
                        options: {
                            sassOptions: {
                                includePaths: ['node_modules/compass-mixins/lib'],
                            },
                        },
                    },*/
                ],
            },
        ]
    },
};
