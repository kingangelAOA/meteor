module.exports = {
  '/api': {
    target: 'http://127.0.0.1:9090/api',
    changeOrigin: true,
    pathRewrite: {
      '^/api': ''
    },
    logLevel: 'debug'
  }
}
