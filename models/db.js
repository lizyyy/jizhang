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
    db.get('SELECT COUNT(*) as count FROM categories', (err, row) => {
        if (err) {
            console.error('检查分类数据失败:', err);
            return;
        }
        
        if (row.count === 0) {
            const categoriesWithSubcategories = {
                expense: [
                    { name: '餐饮', subcategories: ['早餐', '午餐', '晚餐', '零食饮料', '外卖', '聚餐'] },
                    { name: '交通', subcategories: ['公交地铁', '打车', '加油', '停车费', '过路费', '汽车保养'] },
                    { name: '购物', subcategories: ['日用品', '服饰', '数码产品', '家电', '化妆品', '网购'] },
                    { name: '娱乐', subcategories: ['电影', '游戏', 'KTV', '旅游', '运动健身', '兴趣爱好'] },
                    { name: '医疗', subcategories: ['药品', '门诊', '住院', '体检', '保健品'] },
                    { name: '教育', subcategories: ['学费', '培训', '书籍', '文具', '网课'] },
                    { name: '住房', subcategories: ['房租', '水电费', '物业费', '装修', '维修'] },
                    { name: '其他支出', subcategories: ['红包', '礼物', '捐赠', '杂项'] }
                ],
                income: [
                    { name: '工资', subcategories: ['基本工资', '绩效奖金', '年终奖', '加班费'] },
                    { name: '奖金', subcategories: ['项目奖金', '季度奖金', '节日奖金', '其他奖金'] },
                    { name: '投资', subcategories: ['股票', '基金', '理财', '分红', '利息'] },
                    { name: '兼职', subcategories: ['自由职业', '副业', '临时工作'] },
                    { name: '其他收入', subcategories: ['红包', '退款', '报销', '意外收入'] }
                ],
                transfer: [
                    { name: '银行卡转账', subcategories: ['同行转账', '跨行转账', '定期存款'] },
                    { name: '支付宝转账', subcategories: ['转账', '充值', '提现'] },
                    { name: '微信转账', subcategories: ['转账', '充值', '提现'] }
                ],
                loan: [
                    { name: '借款', subcategories: ['借出', '借入'] },
                    { name: '还款', subcategories: ['还本金', '还利息', '提前还款'] }
                ]
            };

            db.serialize(() => {
                const stmt = db.prepare('INSERT INTO categories (name, type) VALUES (?, ?)');
                
                Object.entries(categoriesWithSubcategories).forEach(([type, cats]) => {
                    cats.forEach(cat => {
                        stmt.run(cat.name, type);
                    });
                });
                stmt.finalize();

                db.each('SELECT id, name, type FROM categories', [], (err, category) => {
                    if (err) return;
                    
                    const typeData = categoriesWithSubcategories[category.type];
                    if (typeData) {
                        const catData = typeData.find(c => c.name === category.name);
                        if (catData && catData.subcategories) {
                            const subStmt = db.prepare('INSERT INTO subcategories (category_id, name) VALUES (?, ?)');
                            catData.subcategories.forEach(subName => {
                                subStmt.run(category.id, subName);
                            });
                            subStmt.finalize();
                        }
                    }
                }, () => {
                    console.log('默认分类和子分类数据已初始化');
                });
            });
        }
    });
}

module.exports = {
    db,
    initDatabase
};
