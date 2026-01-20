const express = require('express');
const app = express();

app.use(express.json());

app.get('/', (req, res) => {
  res.json({
    service: 'frontend',
    status: 'running',
    message: 'Hello from Frontend (React Dashboard)'
  });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`Frontend service listening on port ${PORT}`);
});



