const express = require('express')
const app = express()
const port = 3000

app.get('/', (req, res) => {
  res.send('Hello World! This webpage is hosted live on <a href="hostahack.com">hostahack.com</a>')
})

app.listen(port, () => {
  console.log(`Example listening on port ${port}`)
})
