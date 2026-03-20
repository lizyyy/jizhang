const express = require('express');
const path = require('path');
const cors = require('cors');
const { initDatabase } = require('./models/db');
const { errorResponse } = require('./utils/response');

const app = express();
const PORT = process.env.PORT || 9001;

// 中间件
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(express.static(path.join(__dirname, 'public')));

// 初始化数据库
initDatabase();

// 路由
const recordRoutes = require('./routes/records');
const categoryRoutes = require('./routes/categories');

app.use('/api/records', recordRoutes);
app.use('/api/categories', categoryRoutes);

// 首页路由
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// 404 处理
app.use((req, res) => {
    res.status(404).json(errorResponse(404, '接口不存在', 'NOT_FOUND'));
});

// 错误处理
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json(errorResponse(500, '服务器内部错误', 'INTERNAL_ERROR'));
});

app.listen(PORT, () => {
    console.log(`服务器运行在 http://localhost:${PORT}`);
});

module.exports = app;
