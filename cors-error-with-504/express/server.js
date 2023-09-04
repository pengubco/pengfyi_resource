const express = require('express');
const path = require('path');

const app = express();
const PORT = 8081;

app.use(express.static(path.join(__dirname, 'public')));

app.listen(PORT, () => {
    console.log(`Server is running at http://localhost:${PORT}`);
});
