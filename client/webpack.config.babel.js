import toml from 'toml-require';
import path from 'path';
import webpack from 'webpack';
import { CleanWebpackPlugin } from 'clean-webpack-plugin';
import CopyPlugin from 'copy-webpack-plugin';
import HtmlWebpackPlugin from 'html-webpack-plugin';
import PreloadWebpackPlugin from 'preload-webpack-plugin';
import VueLoaderPlugin from 'vue-loader/lib/plugin';
import WebpackPwaManifest from 'webpack-pwa-manifest';

import { InjectManifest } from 'workbox-webpack-plugin';

const BASE_URL = 'http://127.0.0.1:60850';
const DEV = process.env.NODE_ENV !== 'production';
const API_URL = 'http://127.0.0.1:60850/api';
const PUBLIC_PATH = 'http://127.0.0.1:60850/';
const WS_URL = 'ws://127.0.0.1:60850/ws';

toml.install();
const TOML = require('../config.toml');

const iconDstPath = 'images/app/';
const iconSrcPath = path.resolve(__dirname, 'assets/icons/');

export default {
  mode: 'development',

  entry: [
    './src/index.js',
  ],

  target: 'web',

  output: {
    filename: 'js/client/[name].[contenthash].js',
    library: 'SkyStreamer',
    path: path.resolve(__dirname, '../public/'),
    publicPath: TOML.PUBLIC_PATH || PUBLIC_PATH,
  },

  resolve: {
    alias: {
      vue$: 'vue/dist/vue.esm.js',
    },
    modules: [path.resolve('./', 'src'), 'node_modules'],
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        include: path.resolve(__dirname, 'src'),
        loader: 'babel-loader',
      },
      {
        test: /\.vue$/,
        use: 'vue-loader',
      },
      {
        test: /\.(css|less|scss)$/,
        use: [
          'vue-style-loader',
          'css-loader',
          'sass-loader',
        ],
      },
    ],
  },

  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(DEV || 'production'),
      API_URL: JSON.stringify(TOML.API_URL || API_URL),
      BASE_URL: JSON.stringify(TOML.BASE_URL || BASE_URL),
      PUBLIC_PATH: JSON.stringify(TOML.PUBLIC_PATH || PUBLIC_PATH),
      WS_URL: JSON.stringify(TOML.WS_URL || WS_URL),
    }),
    new webpack.HashedModuleIdsPlugin(),
    new CleanWebpackPlugin({
      cleanStaleWebpackAssets: false,
      dry: false,
      verbose: true,
    }),
    new CopyPlugin({
      patterns: [
        { from: 'images/icon-20.png', to: 'images/static/' },
        { from: 'images/logo.png', to: 'images/static/' },
      ],
    }),
    new HtmlWebpackPlugin({
      favicon: 'assets/icons/favicon.ico',
      minify: {
        collapseWhitespace: true,
        minifyCSS: true,
        minifyJS: true,
        removeComments: true,
        removeRedundantAttributes: true,
        removeScriptTypeAttributes: true,
        removeStyleLinkTypeAttributes: true,
        useShortDoctype: true,
      },
      template: 'templates/index.gohtml',
      title: 'SkyStreamer',
      meta: {
        viewport: 'width=device-width, initial-scale=1',
      },
      filename: 'index.gohtml',
    }),
    new InjectManifest({
      importWorkboxFrom: 'disabled',
      swDest: 'service-worker.js',
      swSrc: path.join('src', 'sw.js'),
    }),
    new PreloadWebpackPlugin({
      rel: 'preload',
      include: 'initial',
    }),
    new VueLoaderPlugin(),
    new WebpackPwaManifest({
      name: 'SkyStreamer',
      short_name: 'skystream',
      description: 'WebApp based WebRTC Receiver & Control Sender',
      background_color: '#ffffff',
      theme_color: '#ffffff',
      crossorigin: 'use-credentials', // can be null, use-credentials or anonymous
      icons: [
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-96.png'),
          size: '96x96', // you can also use the specifications pattern
        },
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-128.png'),
          size: '128x128', // you can also use the specifications pattern
        },
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-192.png'),
          size: '192x192', // you can also use the specifications pattern
        },
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-256.png'),
          size: '256x256', // you can also use the specifications pattern
        },
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-512.png'),
          size: '512x512', // multiple sizes
        },
        {
          destination: iconDstPath,
          src: path.join(iconSrcPath, 'icon-1024.png'),
          size: '1024x1024', // you can also use the specifications pattern
        },
      ],
      fingerprints: true,
      start_url: '/',
    }),
  ],
};
