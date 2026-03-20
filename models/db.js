const sqlite3 = require('sqlite3').verbose();
const path = require('path');

const dbPath = path.join(__dirname, '../data/accounting.db');

const db = new sqlite3.Database(dbPath, (err) => {
    if (err) {
        console.error('数据库连接失败:', err);
    } else {
        console.log('数据库连接成功');
    }
});

function initDatabase() {
    db.serialize(() => {
        // 创建一级分类表
        db.run(`CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            type TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`);

        // 创建二级分类表
        db.run(`CREATE TABLE IF NOT EXISTS subcategories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            category_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (category_id) REFERENCES categories(id)
        )`);

        // 创建记账记录表
        db.run(`CREATE TABLE IF NOT EXISTS records (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT NOT NULL,
            category_id INTEGER NOT NULL,
            subcategory_id INTEGER,
            amount REAL NOT NULL,
            date TEXT NOT NULL,
            note TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (category_id) REFERENCES categories(id),
            FOREIGN KEY (subcategory_id) REFERENCES subcategories(id)
        )`);

        // 初始化默认分类数据
        initDefaultData();
    });
}

function initDefaultData() {
    // 检查是否已有数据
    db.get('SELECT COUNT(*) as count FROM categories', (err, row) => {
        if (err) {
            console.error('检查分类数据失败:', err);
            return;
        }
        
        if (row.count === 0) {
            console.log('开始初始化默认分类数据...');
            
            // 插入默认支出分类
            const expenseCategories = [
                { name: '餐饮', type: 'expense' },
                { name: '交通', type: 'expense' },
                { name: '购物', type: 'expense' },
                { name: '娱乐', type: 'expense' },
                { name: '医疗', type: 'expense' },
                { name: '教育', type: 'expense' },
                { name: '住房', type: 'expense' },
                { name: '通讯', type: 'expense' },
                { name: '服饰', type: 'expense' },
                { name: '美容', type: 'expense' },
                { name: '运动', type: 'expense' },
                { name: '其他支出', type: 'expense' }
            ];

            // 插入默认收入分类
            const incomeCategories = [
                { name: '工资', type: 'income' },
                { name: '奖金', type: 'income' },
                { name: '投资收益', type: 'income' },
                { name: '兼职', type: 'income' },
                { name: '红包礼金', type: 'income' },
                { name: '其他收入', type: 'income' }
            ];

            // 插入默认转账分类
            const transferCategories = [
                { name: '银行卡转账', type: 'transfer' },
                { name: '支付宝转账', type: 'transfer' },
                { name: '微信转账', type: 'transfer' },
                { name: '信用卡还款', type: 'transfer' }
            ];

            // 插入默认贷款分类
            const loanCategories = [
                { name: '借入', type: 'loan' },
                { name: '借出', type: 'loan' },
                { name: '还款', type: 'loan' },
                { name: '收款', type: 'loan' }
            ];

            const allCategories = [...expenseCategories, ...incomeCategories, ...transferCategories, ...loanCategories];
            
            // 使用串行方式插入一级分类
            insertCategoriesSequentially(allCategories, 0, () => {
                console.log('一级分类插入完成，开始插入子分类...');
                insertAllSubcategories();
            });
        }
    });
}

// 串行插入分类（确保顺序和正确获取lastID）
function insertCategoriesSequentially(categories, index, callback) {
    if (index >= categories.length) {
        callback();
        return;
    }
    
    const cat = categories[index];
    db.run('INSERT INTO categories (name, type) VALUES (?, ?)', [cat.name, cat.type], function(err) {
        if (err) {
            console.error('插入分类失败:', err);
        }
        insertCategoriesSequentially(categories, index + 1, callback);
    });
}

// 插入所有子分类
function insertAllSubcategories() {
    // 先获取所有分类的ID映射
    db.all('SELECT id, name, type FROM categories', [], (err, cats) => {
        if (err) {
            console.error('获取分类ID失败:', err);
            return;
        }

        // 构建ID映射
        const categoryIdMap = {};
        cats.forEach(cat => {
            categoryIdMap[`${cat.type}_${cat.name}`] = cat.id;
        });

        // 定义默认二级分类
        const subcategories = [
            // 餐饮子类
            { category: { type: 'expense', name: '餐饮' }, name: '早餐' },
            { category: { type: 'expense', name: '餐饮' }, name: '午餐' },
            { category: { type: 'expense', name: '餐饮' }, name: '晚餐' },
            { category: { type: 'expense', name: '餐饮' }, name: '水果' },
            { category: { type: 'expense', name: '餐饮' }, name: '零食' },
            { category: { type: 'expense', name: '餐饮' }, name: '饮料' },
            { category: { type: 'expense', name: '餐饮' }, name: '聚餐' },
            // 交通子类
            { category: { type: 'expense', name: '交通' }, name: '地铁/公交' },
            { category: { type: 'expense', name: '交通' }, name: '出租车' },
            { category: { type: 'expense', name: '交通' }, name: '加油费' },
            { category: { type: 'expense', name: '交通' }, name: '停车费' },
            { category: { type: 'expense', name: '交通' }, name: '火车票' },
            { category: { type: 'expense', name: '交通' }, name: '机票' },
            // 购物子类
            { category: { type: 'expense', name: '购物' }, name: '日用品' },
            { category: { type: 'expense', name: '购物' }, name: '数码产品' },
            { category: { type: 'expense', name: '购物' }, name: '家居用品' },
            { category: { type: 'expense', name: '购物' }, name: '书籍' },
            // 娱乐子类
            { category: { type: 'expense', name: '娱乐' }, name: '电影' },
            { category: { type: 'expense', name: '娱乐' }, name: '游戏' },
            { category: { type: 'expense', name: '娱乐' }, name: '旅游' },
            { category: { type: 'expense', name: '娱乐' }, name: 'KTV' },
            // 住房子类
            { category: { type: 'expense', name: '住房' }, name: '房租' },
            { category: { type: 'expense', name: '住房' }, name: '水电煤' },
            { category: { type: 'expense', name: '住房' }, name: '物业费' },
            // 通讯子类
            { category: { type: 'expense', name: '通讯' }, name: '手机话费' },
            { category: { type: 'expense', name: '通讯' }, name: '宽带费' },
            // 服饰子类
            { category: { type: 'expense', name: '服饰' }, name: '衣服' },
            { category: { type: 'expense', name: '服饰' }, name: '鞋子' },
            { category: { type: 'expense', name: '服饰' }, name: '包包' },
            // 美容子类
            { category: { type: 'expense', name: '美容' }, name: '护肤品' },
            { category: { type: 'expense', name: '美容' }, name: '化妆品' },
            { category: { type: 'expense', name: '美容' }, name: '美发' },
            // 运动子类
            { category: { type: 'expense', name: '运动' }, name: '健身' },
            { category: { type: 'expense', name: '运动' }, name: '体育用品' },
            // 工资子类
            { category: { type: 'income', name: '工资' }, name: '基本工资' },
            { category: { type: 'income', name: '工资' }, name: '加班费' },
            { category: { type: 'income', name: '工资' }, name: '津贴补贴' },
            // 奖金子类
            { category: { type: 'income', name: '奖金' }, name: '年终奖' },
            { category: { type: 'income', name: '奖金' }, name: '绩效奖' },
            // 投资收益子类
            { category: { type: 'income', name: '投资收益' }, name: '股票' },
            { category: { type: 'income', name: '投资收益' }, name: '基金' },
            { category: { type: 'income', name: '投资收益' }, name: '利息' }
        ];

        // 插入二级分类
        const subStmt = db.prepare('INSERT INTO subcategories (category_id, name) VALUES (?, ?)');
        subcategories.forEach(sub => {
            const categoryKey = `${sub.category.type}_${sub.category.name}`;
            const categoryId = categoryIdMap[categoryKey];
            if (categoryId) {
                subStmt.run(categoryId, sub.name, (err) => {
                    if (err) {
                        console.error('插入子分类失败:', err, sub);
                    }
                });
            } else {
                console.warn('未找到分类ID:', categoryKey);
            }
        });
        subStmt.finalize((err) => {
            if (err) {
                console.error('finalize错误:', err);
            } else {
                console.log('默认分类和子分类数据初始化完成！');
            }
        });
    });
}

module.exports = {
    db,
    initDatabase
};
