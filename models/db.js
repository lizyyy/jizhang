const sqlite3 = require('sqlite3').verbose();
const path = require('path');
const fs = require('fs');

const dataDir = path.join(__dirname, '../data');
if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
}

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
            // 定义默认分类数据（包含子分类）
            const defaultCategories = [
                // 支出分类
                {
                    name: '餐饮',
                    type: 'expense',
                    subcategories: ['早餐', '午餐', '晚餐', '零食', '饮料', '聚餐']
                },
                {
                    name: '交通',
                    type: 'expense',
                    subcategories: ['公交地铁', '出租车', '加油', '停车费', '过路费', '保养维修']
                },
                {
                    name: '购物',
                    type: 'expense',
                    subcategories: ['服装', '鞋帽', '箱包', '化妆品', '日用品', '电子产品', '书籍']
                },
                {
                    name: '居住',
                    type: 'expense',
                    subcategories: ['房租', '房贷', '水电煤', '物业费', '维修费']
                },
                {
                    name: '娱乐',
                    type: 'expense',
                    subcategories: ['电影', '游戏', '旅游', '运动健身', '兴趣爱好']
                },
                {
                    name: '医疗',
                    type: 'expense',
                    subcategories: ['挂号费', '药品', '检查费', '住院费', '保健品']
                },
                {
                    name: '教育',
                    type: 'expense',
                    subcategories: ['学费', '培训费', '书籍资料', '考试费']
                },
                {
                    name: '通讯',
                    type: 'expense',
                    subcategories: ['手机费', '宽带费', '视频会员', '音乐会员']
                },
                {
                    name: '人情',
                    type: 'expense',
                    subcategories: ['送礼', '红包', '请客', '份子钱']
                },
                {
                    name: '其他支出',
                    type: 'expense',
                    subcategories: ['捐赠', '罚款', '其他']
                },
                // 收入分类
                {
                    name: '工资',
                    type: 'income',
                    subcategories: ['基本工资', '绩效工资', '加班费', '年终奖']
                },
                {
                    name: '奖金',
                    type: 'income',
                    subcategories: ['项目奖金', '全勤奖', '节日奖金', '优秀员工奖']
                },
                {
                    name: '投资',
                    type: 'income',
                    subcategories: ['股票', '基金', '理财', '房产', '利息']
                },
                {
                    name: '兼职',
                    type: 'income',
                    subcategories: ['写作', '设计', '翻译', '代驾', '其他兼职']
                },
                {
                    name: '其他收入',
                    type: 'income',
                    subcategories: ['红包', '退款', '二手出售', '其他']
                },
                // 转账分类
                {
                    name: '银行卡转账',
                    type: 'transfer',
                    subcategories: ['同行转账', '跨行转账']
                },
                {
                    name: '支付宝转账',
                    type: 'transfer',
                    subcategories: ['余额转账', '银行卡转账']
                },
                {
                    name: '微信转账',
                    type: 'transfer',
                    subcategories: ['零钱转账', '银行卡转账']
                },
                // 贷款分类
                {
                    name: '借款',
                    type: 'loan',
                    subcategories: ['银行贷款', '网贷', '亲友借款']
                },
                {
                    name: '还款',
                    type: 'loan',
                    subcategories: ['本金', '利息', '逾期费用']
                }
            ];

            // 使用递归方式顺序插入数据
            let catIndex = 0;
            
            function insertNextCategory() {
                if (catIndex >= defaultCategories.length) {
                    console.log('默认分类数据已初始化（包含二级分类）');
                    return;
                }
                
                const cat = defaultCategories[catIndex];
                catIndex++;
                
                db.run('INSERT INTO categories (name, type) VALUES (?, ?)', [cat.name, cat.type], function(err) {
                    if (err) {
                        console.error('插入分类失败:', err);
                        insertNextCategory();
                        return;
                    }
                    
                    const categoryId = this.lastID;
                    let subIndex = 0;
                    
                    function insertNextSubcategory() {
                        if (subIndex >= cat.subcategories.length) {
                            insertNextCategory();
                            return;
                        }
                        
                        const subName = cat.subcategories[subIndex];
                        subIndex++;
                        
                        db.run('INSERT INTO subcategories (category_id, name) VALUES (?, ?)', 
                            [categoryId, subName], function(err) {
                                if (err) {
                                    console.error('插入子分类失败:', err);
                                }
                                insertNextSubcategory();
                            });
                    }
                    
                    insertNextSubcategory();
                });
            }
            
            insertNextCategory();
        }
    });
}

module.exports = {
    db,
    initDatabase
};
